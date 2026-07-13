package helm

import (
	"context"
	"encoding/base64"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

// The provisioner creates a random dir under XDG_RUNTIME_DIR and rewrites the
// command line, so it can't go through plugintest.TestProvisioner's strict
// output equality. Test it directly, including executing the self-destruct
// wrapper to prove the kubeconfig is removed once the command exits.
func TestKubeconfigProvisioner(t *testing.T) {
	rawConfig := plugintest.LoadFixture(t, "config")
	encodedConfig := base64.StdEncoding.EncodeToString([]byte(rawConfig))

	runtimeDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", runtimeDir)

	in := sdk.ProvisionInput{
		ItemFields: map[sdk.FieldName]string{fieldname.Credential: encodedConfig},
		HomeDir:    "~",
		TempDir:    t.TempDir(),
	}
	out := sdk.ProvisionOutput{
		Environment: make(map[string]string),
		Files:       make(map[string]sdk.OutputFile),
		CommandLine: []string{"true"},
	}

	(&helmKubeconfigProvisioner{}).Provision(context.Background(), in, &out)

	if len(out.Diagnostics.Errors) > 0 {
		t.Fatalf("Provision returned errors: %v", out.Diagnostics.Errors)
	}

	configPath := out.Environment["KUBECONFIG"]
	if configPath == "" {
		t.Fatal("KUBECONFIG not set")
	}
	if !strings.HasPrefix(configPath, runtimeDir+string(os.PathSeparator)) {
		t.Fatalf("kubeconfig %q not under XDG_RUNTIME_DIR %q", configPath, runtimeDir)
	}
	contents, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("reading provisioned kubeconfig: %v", err)
	}
	if string(contents) != rawConfig {
		t.Fatal("provisioned kubeconfig contents mismatch")
	}
	info, _ := os.Stat(configPath)
	if info.Mode().Perm() != 0600 {
		t.Fatalf("kubeconfig mode = %v, want 0600", info.Mode().Perm())
	}

	if len(out.CommandLine) != 5 || out.CommandLine[0] != "/bin/sh" || out.CommandLine[4] != "true" {
		t.Fatalf("unexpected wrapped command line: %v", out.CommandLine)
	}

	// Execute the wrapped command line and verify the self-destruct trap
	// removes the kubeconfig dir while preserving the exit status.
	if _, err := os.Stat("/bin/sh"); err != nil {
		t.Skip("/bin/sh not available")
	}
	cmd := exec.Command(out.CommandLine[0], out.CommandLine[1:]...)
	if err := cmd.Run(); err != nil {
		t.Fatalf("wrapped command failed: %v", err)
	}
	if _, err := os.Stat(filepath.Dir(configPath)); !os.IsNotExist(err) {
		t.Fatalf("kubeconfig dir still exists after wrapped command exited")
	}
}

func TestKubeconfigImporter(t *testing.T) {
	rawConfig := plugintest.LoadFixture(t, "config")
	encodedConfig := base64.StdEncoding.EncodeToString([]byte(rawConfig))

	plugintest.TestImporter(t, Kubeconfig().Importer, map[string]plugintest.ImportCase{
		"kubeconfig file": {
			Files: map[string]string{
				"~/.kube/config": rawConfig,
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Credential: encodedConfig,
					},
					NameHint: "",
				},
			},
		},
	})
}
