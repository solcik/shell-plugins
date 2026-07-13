package azure

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

// azureConfigDirProvisioner authenticates the Azure CLI by running
// `az login --service-principal` against a temporary AZURE_CONFIG_DIR.
// The az CLI has no environment-variable credential support, so login state
// must live in a config dir; using a throwaway one keeps the user's ~/.azure
// untouched.
//
// Cleanup cannot rely on Deprovision: op never calls it for local plugins.
// Instead the config dir is created under XDG_RUNTIME_DIR (tmpfs, 0700,
// wiped at logout) and the executed command line is wrapped in a shell trap
// that removes the dir the moment the CLI exits, so the provisioned tokens
// never outlive the invocation.
type azureConfigDirProvisioner struct{}

func (p *azureConfigDirProvisioner) Description() string {
	return "Provision a temporary AZURE_CONFIG_DIR logged in with a service principal"
}

func (p *azureConfigDirProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	// Prefer XDG_RUNTIME_DIR over op's TempDir under /tmp: op does not clean
	// TempDir for local plugins, while the runtime dir is tmpfs and dies with
	// the session even if the self-destruct wrapper is skipped (SIGKILL).
	parent := os.Getenv("XDG_RUNTIME_DIR")
	if parent == "" {
		parent = in.TempDir
	}
	configDir, err := os.MkdirTemp(parent, "op-azure-")
	if err != nil {
		out.AddError(err)
		return
	}
	// Telemetry runs as a detached child that outlives az and re-creates the
	// config dir after the self-destruct trap has removed it — disable it.
	env := append(os.Environ(),
		"AZURE_CONFIG_DIR="+configDir,
		"AZURE_CORE_COLLECT_TELEMETRY=0")

	login := exec.CommandContext(ctx, "az", "login", "--service-principal",
		"--username", in.ItemFields[fieldname.ClientID],
		"--password", in.ItemFields[fieldname.ClientSecret],
		"--tenant", in.ItemFields[fieldname.TenantID],
		"--output", "none", "--only-show-errors")
	login.Env = env
	if output, err := login.CombinedOutput(); err != nil {
		os.RemoveAll(configDir)
		out.AddError(fmt.Errorf("az login failed: %w: %s", err, output))
		return
	}

	if subscription := in.ItemFields[fieldname.Subscription]; subscription != "" {
		set := exec.CommandContext(ctx, "az", "account", "set",
			"--subscription", subscription, "--only-show-errors")
		set.Env = env
		if output, err := set.CombinedOutput(); err != nil {
			os.RemoveAll(configDir)
			out.AddError(fmt.Errorf("az account set failed: %w: %s", err, output))
			return
		}
	}

	out.AddEnvVar("AZURE_CONFIG_DIR", configDir)
	out.AddEnvVar("AZURE_CORE_COLLECT_TELEMETRY", "0")
	// The trap wrapper gives instant cleanup, but op drops CommandLine
	// rewrites when merging multi-credential outputs, so the reaper is the
	// guaranteed path.
	out.CommandLine = selfDestructCommandLine(out.CommandLine, configDir)
	spawnReaper(configDir, os.Getppid())
}

func (p *azureConfigDirProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
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

// spawnReaper starts a detached watcher that removes dir once the process
// with the given pid (the op CLI — this plugin's parent during Provision)
// exits. It does not depend on op honoring CommandLine or calling
// Deprovision, neither of which is guaranteed for local plugins.
func spawnReaper(dir string, pid int) {
	if !filepath.IsAbs(dir) || filepath.Clean(dir) == "/" || pid <= 1 {
		return
	}
	cmd := exec.Command("/bin/sh", "-c",
		`while kill -0 "$1" 2>/dev/null; do sleep 0.2; done; rm -rf -- "$0"`,
		dir, strconv.Itoa(pid))
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	// No Wait: the reaper must outlive this plugin process; once orphaned it
	// is reparented to init, which reaps it.
	_ = cmd.Start()
}
