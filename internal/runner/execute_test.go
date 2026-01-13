package runner

import (
	"strings"
	"testing"
)

func TestBuildArgs_Basic(t *testing.T) {
	cfg := RunConfig{
		Workflow: "deploy.yml",
	}

	args := BuildArgs(cfg)

	if len(args) != 3 {
		t.Errorf("expected 3 args, got %d: %v", len(args), args)
	}

	if args[0] != "workflow" || args[1] != "run" || args[2] != "deploy.yml" {
		t.Errorf("unexpected args: %v", args)
	}
}

func TestBuildArgs_WithBranch(t *testing.T) {
	cfg := RunConfig{
		Workflow: "deploy.yml",
		Branch:   "main",
	}

	args := BuildArgs(cfg)

	hasRef := false
	for i, arg := range args {
		if arg == "--ref" && i+1 < len(args) && args[i+1] == "main" {
			hasRef = true
			break
		}
	}

	if !hasRef {
		t.Errorf("expected --ref main in args: %v", args)
	}
}

func TestBuildArgs_WithInputs(t *testing.T) {
	cfg := RunConfig{
		Workflow: "deploy.yml",
		Inputs: map[string]string{
			"environment": "production",
			"dry_run":     "false",
		},
	}

	args := BuildArgs(cfg)

	hasEnv := false
	for i, arg := range args {
		if arg == "-f" && i+1 < len(args) && args[i+1] == "environment=production" {
			hasEnv = true
			break
		}
	}

	if !hasEnv {
		t.Errorf("expected -f environment=production in args: %v", args)
	}
}

func TestBuildArgs_EmptyInputsSkipped(t *testing.T) {
	cfg := RunConfig{
		Workflow: "deploy.yml",
		Inputs: map[string]string{
			"environment": "production",
			"empty":       "",
		},
	}

	args := BuildArgs(cfg)

	for _, arg := range args {
		if strings.Contains(arg, "empty=") {
			t.Errorf("empty input should be skipped: %v", args)
		}
	}
}

func TestFormatCommand(t *testing.T) {
	args := []string{"workflow", "run", "deploy.yml", "--ref", "main", "-f", "env=prod"}

	cmd := FormatCommand(args)

	if !strings.HasPrefix(cmd, "gh workflow run") {
		t.Errorf("expected command to start with 'gh workflow run': %s", cmd)
	}

	if !strings.Contains(cmd, "deploy.yml") {
		t.Errorf("expected command to contain 'deploy.yml': %s", cmd)
	}
}

func TestFormatCommand_QuotesSpecialChars(t *testing.T) {
	args := []string{"workflow", "run", "deploy.yml", "-f", "message=hello world"}

	cmd := FormatCommand(args)

	if !strings.Contains(cmd, "\"message=hello world\"") {
		t.Errorf("expected quoted value with spaces: %s", cmd)
	}
}

func TestDryRun(t *testing.T) {
	cfg := RunConfig{
		Workflow: "deploy.yml",
		Branch:   "main",
		Inputs:   map[string]string{"env": "prod"},
	}

	cmd := DryRun(cfg)

	if cmd == "" {
		t.Error("expected non-empty dry run output")
	}

	if !strings.Contains(cmd, "deploy.yml") {
		t.Errorf("expected command to contain workflow: %s", cmd)
	}
}
