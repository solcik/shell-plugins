package helm

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

// helmKubeconfigProvisioner writes the kubeconfig as a real file (not via
// out.AddSecretFile, which op materializes as a one-read FIFO — Helm reads
// the kubeconfig multiple times and blocks on the second read).
//
// Cleanup cannot rely on Deprovision: op never calls it for local plugins.
// Instead the kubeconfig dir is created under XDG_RUNTIME_DIR (tmpfs, 0700,
// wiped at logout) and the executed command line is wrapped in a shell trap
// that removes the dir the moment the CLI exits, so the credential never
// outlives the invocation.
type helmKubeconfigProvisioner struct{}

func (p *helmKubeconfigProvisioner) Description() string {
	return "Provision kubeconfig file for Helm"
}

func (p *helmKubeconfigProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	// Decode base64 kubeconfig
	encoded := in.ItemFields[fieldname.Credential]
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		out.AddError(err)
		return
	}

	// Prefer XDG_RUNTIME_DIR over op's TempDir under /tmp: op does not clean
	// TempDir for local plugins, while the runtime dir is tmpfs and dies with
	// the session even if the self-destruct wrapper is skipped (SIGKILL).
	parent := os.Getenv("XDG_RUNTIME_DIR")
	if parent == "" {
		parent = in.TempDir
	}
	dir, err := os.MkdirTemp(parent, "op-helm-")
	if err != nil {
		out.AddError(err)
		return
	}

	configPath := filepath.Join(dir, "config")
	if err := os.WriteFile(configPath, decoded, 0600); err != nil {
		os.RemoveAll(dir)
		out.AddError(err)
		return
	}
	out.AddEnvVar("KUBECONFIG", configPath)
	out.CommandLine = selfDestructCommandLine(out.CommandLine, dir)
}

func (p *helmKubeconfigProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Never called by op for local plugins; cleanup is handled by the
	// self-destruct command wrapper set in Provision.
}

// selfDestructCommandLine wraps the command op is about to execute in a shell
// that removes dir as soon as the command exits, preserving its exit status.
// The dir is passed as $0 rather than interpolated into the script, so no
// path byte (quote, $, backtick, space) can break quoting or inject shell
// syntax into the trap.
func selfDestructCommandLine(commandLine []string, dir string) []string {
	if len(commandLine) == 0 {
		return commandLine
	}
	// MkdirTemp guarantees an absolute, non-root path; guard anyway so the
	// trap can never be wired up against "/" or a relative path.
	if !filepath.IsAbs(dir) || filepath.Clean(dir) == "/" {
		return commandLine
	}
	return append([]string{"/bin/sh", "-c", `trap 'rm -rf -- "$0"' EXIT; "$@"`, dir}, commandLine...)
}
