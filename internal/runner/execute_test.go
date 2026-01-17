package runner

import (
	"strings"
	"testing"
)

func TestBuildArgs(t *testing.T) {
	tests := []struct {
		name         string
		cfg          RunConfig
		wantContains []string
		wantExcludes []string
		wantLen      int
	}{
		{
			name: "basic",
			cfg: RunConfig{
				Workflow: "deploy.yml",
			},
			wantContains: []string{"workflow", "run", "deploy.yml"},
			wantLen:      3,
		},
		{
			name: "with branch",
			cfg: RunConfig{
				Workflow: "deploy.yml",
				Branch:   "main",
			},
			wantContains: []string{"workflow", "run", "deploy.yml", "--ref", "main"},
		},
		{
			name: "with inputs",
			cfg: RunConfig{
				Workflow: "deploy.yml",
				Inputs: map[string]string{
					"environment": "production",
					"dry_run":     "false",
				},
			},
			wantContains: []string{"workflow", "run", "deploy.yml", "-f", "environment=production", "-f", "dry_run=false"},
		},
		{
			name: "empty inputs skipped",
			cfg: RunConfig{
				Workflow: "deploy.yml",
				Inputs: map[string]string{
					"environment": "production",
					"empty":       "",
				},
			},
			wantContains: []string{"environment=production"},
			wantExcludes: []string{"empty="},
		},
		{
			name: "all options",
			cfg: RunConfig{
				Workflow: "ci.yml",
				Branch:   "feature/test",
				Inputs: map[string]string{
					"env":     "staging",
					"verbose": "true",
				},
			},
			wantContains: []string{"workflow", "run", "ci.yml", "--ref", "feature/test", "-f", "env=staging", "-f", "verbose=true"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := BuildArgs(tt.cfg)

			if tt.wantLen > 0 && len(args) != tt.wantLen {
				t.Errorf("BuildArgs() length = %d, want %d: %v", len(args), tt.wantLen, args)
			}

			argsStr := strings.Join(args, " ")
			for _, want := range tt.wantContains {
				if !strings.Contains(argsStr, want) {
					t.Errorf("BuildArgs() missing %q in: %v", want, args)
				}
			}

			for _, exclude := range tt.wantExcludes {
				if strings.Contains(argsStr, exclude) {
					t.Errorf("BuildArgs() should not contain %q in: %v", exclude, args)
				}
			}
		})
	}
}

func TestFormatCommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		wantPrefix   string
		wantContains []string
	}{
		{
			name:         "basic command",
			args:         []string{"workflow", "run", "deploy.yml", "--ref", "main", "-f", "env=prod"},
			wantPrefix:   "gh workflow run",
			wantContains: []string{"deploy.yml", "main", "env=prod"},
		},
		{
			name:         "quotes special chars",
			args:         []string{"workflow", "run", "deploy.yml", "-f", "message=hello world"},
			wantPrefix:   "gh workflow run",
			wantContains: []string{"\"message=hello world\""},
		},
		{
			name:         "quotes args with equals",
			args:         []string{"workflow", "run", "test.yml", "-f", "key=value"},
			wantPrefix:   "gh workflow run",
			wantContains: []string{"\"key=value\""},
		},
		{
			name:         "no special chars",
			args:         []string{"workflow", "run", "simple.yml"},
			wantPrefix:   "gh workflow run",
			wantContains: []string{"simple.yml"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := FormatCommand(tt.args)

			if !strings.HasPrefix(cmd, tt.wantPrefix) {
				t.Errorf("FormatCommand() prefix = %q, want %q", cmd, tt.wantPrefix)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(cmd, want) {
					t.Errorf("FormatCommand() missing %q in: %s", want, cmd)
				}
			}
		})
	}
}

func TestDryRun(t *testing.T) {
	tests := []struct {
		name         string
		cfg          RunConfig
		wantContains []string
	}{
		{
			name: "full config",
			cfg: RunConfig{
				Workflow: "deploy.yml",
				Branch:   "main",
				Inputs:   map[string]string{"env": "prod"},
			},
			wantContains: []string{"gh workflow run", "deploy.yml", "main", "env=prod"},
		},
		{
			name: "minimal config",
			cfg: RunConfig{
				Workflow: "test.yml",
			},
			wantContains: []string{"gh workflow run", "test.yml"},
		},
		{
			name: "with multiple inputs",
			cfg: RunConfig{
				Workflow: "ci.yml",
				Branch:   "feature",
				Inputs: map[string]string{
					"debug":   "true",
					"verbose": "1",
				},
			},
			wantContains: []string{"gh workflow run", "ci.yml", "feature"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := DryRun(tt.cfg)

			if cmd == "" {
				t.Error("DryRun() returned empty string")
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(cmd, want) {
					t.Errorf("DryRun() missing %q in: %s", want, cmd)
				}
			}
		})
	}
}
