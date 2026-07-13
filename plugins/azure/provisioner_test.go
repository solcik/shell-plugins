package azure

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
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
	if len(wrapped) != 5 || wrapped[0] != "/bin/sh" || wrapped[4] != "true" {
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

	// Empty command lines are passed through untouched.
	if got := selfDestructCommandLine(nil, dir); got != nil {
		t.Fatalf("expected nil passthrough, got %v", got)
	}
}
