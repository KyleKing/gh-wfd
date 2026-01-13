package runner

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunConfig holds the configuration for running a workflow.
type RunConfig struct {
	Workflow string
	Branch   string
	Inputs   map[string]string
	Watch    bool
}

// BuildArgs constructs the gh workflow run arguments.
func BuildArgs(cfg RunConfig) []string {
	args := []string{"workflow", "run", cfg.Workflow}

	if cfg.Branch != "" {
		args = append(args, "--ref", cfg.Branch)
	}

	for k, v := range cfg.Inputs {
		if v != "" {
			args = append(args, "-f", k+"="+v)
		}
	}

	return args
}

// FormatCommand returns a human-readable command string.
func FormatCommand(args []string) string {
	quoted := make([]string, len(args))
	for i, arg := range args {
		if strings.Contains(arg, " ") || strings.Contains(arg, "=") {
			quoted[i] = fmt.Sprintf("%q", arg)
		} else {
			quoted[i] = arg
		}
	}
	return "gh " + strings.Join(quoted, " ")
}

// Execute runs the workflow using gh CLI.
// It prints the command being run (like lazygit) then executes it.
func Execute(cfg RunConfig) error {
	args := BuildArgs(cfg)

	fmt.Println()
	fmt.Println("Running command:")
	fmt.Println("  " + FormatCommand(args))
	fmt.Println()

	cmd := exec.Command("gh", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gh workflow run failed: %w", err)
	}

	if cfg.Watch {
		return watchLatestRun(cfg.Workflow)
	}

	return nil
}

func watchLatestRun(workflow string) error {
	fmt.Println()
	fmt.Println("Watching run...")
	fmt.Println()

	cmd := exec.Command("gh", "run", "watch")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// DryRun prints the command that would be executed without running it.
func DryRun(cfg RunConfig) string {
	args := BuildArgs(cfg)
	return FormatCommand(args)
}
