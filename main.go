package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyleking/gh-workflow-runner/internal/app"
	"github.com/kyleking/gh-workflow-runner/internal/frecency"
	"github.com/kyleking/gh-workflow-runner/internal/runner"
	"github.com/kyleking/gh-workflow-runner/internal/workflow"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	workflows, err := workflow.Discover(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering workflows: %v\n", err)
		os.Exit(1)
	}

	if len(workflows) == 0 {
		fmt.Println("No dispatchable workflows found in .github/workflows/")
		os.Exit(0)
	}

	repo, err := runner.DetectRepo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not detect repository: %v\n", err)
		repo = "unknown/unknown"
	}

	history, err := frecency.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not load history: %v\n", err)
		history = frecency.NewStore()
	}

	model := app.New(workflows, history, repo)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
