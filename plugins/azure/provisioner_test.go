package azure

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// Provision itself execs `az login`, so it can't run in tests; the
// self-destruct wrapper is the security-critical part and is tested here,
// including executing it to prove the dir is removed once the command exits.
func TestSelfDestructCommandLine(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "op-azure-test")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "msal_token_cache.json"), []byte("secret"), 0600); err != nil {
		t.Fatal(err)
	}

	wrapped := selfDestructCommandLine([]string{"true"}, dir)
	if len(wrapped) != 5 || wrapped[0] != "/bin/sh" || wrapped[3] != dir || wrapped[4] != "true" {
		t.Fatalf("unexpected wrapped command line: %v", wrapped)
	}

	if _, err := os.Stat("/bin/sh"); err != nil {
		t.Skip("/bin/sh not available")
	}
	if err := exec.Command(wrapped[0], wrapped[1:]...).Run(); err != nil {
		t.Fatalf("wrapped command failed: %v", err)
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatal("config dir still exists after wrapped command exited")
	}

	// Exit status of the wrapped command must pass through the trap.
	err := exec.Command("/bin/sh", append([]string{}, selfDestructCommandLine([]string{"sh", "-c", "exit 42"}, t.TempDir())[1:]...)...).Run()
	exitErr, ok := err.(*exec.ExitError)
	if !ok || exitErr.ExitCode() != 42 {
		t.Fatalf("exit status not preserved, got %v", err)
	}

	// Hostile path bytes must neither break quoting nor execute: the dir is
	// passed as $0, never interpolated into the script.
	marker := filepath.Join(t.TempDir(), "pwned")
	hostile := filepath.Join(t.TempDir(), `evil '$(touch `+marker+`) "dir`)
	if err := os.MkdirAll(hostile, 0700); err != nil {
		t.Fatal(err)
	}
	if err := exec.Command("/bin/sh", selfDestructCommandLine([]string{"true"}, hostile)[1:]...).Run(); err != nil {
		t.Fatalf("wrapped command with hostile dir failed: %v", err)
	}
	if _, err := os.Stat(hostile); !os.IsNotExist(err) {
		t.Fatal("hostile-named dir still exists after wrapped command exited")
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("shell injection executed: marker file created")
	}

	// Empty command lines and unsafe dirs are passed through untouched.
	if got := selfDestructCommandLine(nil, dir); got != nil {
		t.Fatalf("expected nil passthrough, got %v", got)
	}
	for _, unsafe := range []string{"/", "relative/path", ""} {
		if got := selfDestructCommandLine([]string{"true"}, unsafe); len(got) != 1 || got[0] != "true" {
			t.Fatalf("dir %q: expected unwrapped passthrough, got %v", unsafe, got)
		}
	}
}

// The reaper must remove the dir shortly after the watched process exits,
// without any cooperation from the watched process itself.
func TestSpawnReaper(t *testing.T) {
	if _, err := os.Stat("/bin/sh"); err != nil {
		t.Skip("/bin/sh not available")
	}

	dir := filepath.Join(t.TempDir(), "op-azure-reap")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}

	// Stand-in for the op CLI: a short-lived process the reaper watches.
	fakeOp := exec.Command("sleep", "0.3")
	if err := fakeOp.Start(); err != nil {
		t.Fatal(err)
	}

	spawnReaper(dir, fakeOp.Process.Pid)
	_ = fakeOp.Wait()

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("reaper did not remove dir after watched process exited")
}
