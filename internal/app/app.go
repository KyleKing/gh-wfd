package app

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyleking/gh-workflow-runner/internal/frecency"
	"github.com/kyleking/gh-workflow-runner/internal/ui"
	"github.com/kyleking/gh-workflow-runner/internal/workflow"
)

// FocusedPane represents which pane currently has focus.
type FocusedPane int

const (
	PaneWorkflows FocusedPane = iota
	PaneHistory
	PaneConfig
)

// Model is the root bubbletea model for the application.
type Model struct {
	focused   FocusedPane
	workflows []workflow.WorkflowFile
	history   *frecency.Store
	repo      string

	selectedWorkflow int
	selectedHistory  int
	branch           string
	inputs           map[string]string
	watchRun         bool

	width  int
	height int
	keys   KeyMap
}

// New creates a new application model.
func New(workflows []workflow.WorkflowFile, history *frecency.Store, repo string) Model {
	m := Model{
		focused:   PaneWorkflows,
		workflows: workflows,
		history:   history,
		repo:      repo,
		inputs:    make(map[string]string),
		keys:      DefaultKeyMap(),
	}

	if len(workflows) > 0 {
		m.initializeInputs(workflows[0])
	}

	return m
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Tab):
			m.focused = (m.focused + 1) % 3
			return m, nil

		case key.Matches(msg, m.keys.ShiftTab):
			m.focused = (m.focused + 2) % 3
			return m, nil

		case key.Matches(msg, m.keys.Up):
			m.handleUp()
			return m, nil

		case key.Matches(msg, m.keys.Down):
			m.handleDown()
			return m, nil

		case key.Matches(msg, m.keys.Enter):
			return m.handleEnter()

		case key.Matches(msg, m.keys.Watch):
			m.watchRun = !m.watchRun
			return m, nil
		}
	}

	return m, nil
}

func (m *Model) handleUp() {
	switch m.focused {
	case PaneWorkflows:
		if m.selectedWorkflow > 0 {
			m.selectedWorkflow--
			m.initializeInputs(m.workflows[m.selectedWorkflow])
		}
	case PaneHistory:
		if m.selectedHistory > 0 {
			m.selectedHistory--
		}
	}
}

func (m *Model) handleDown() {
	switch m.focused {
	case PaneWorkflows:
		if m.selectedWorkflow < len(m.workflows)-1 {
			m.selectedWorkflow++
			m.initializeInputs(m.workflows[m.selectedWorkflow])
		}
	case PaneHistory:
		entries := m.currentHistoryEntries()
		if m.selectedHistory < len(entries)-1 {
			m.selectedHistory++
		}
	}
}

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.focused {
	case PaneHistory:
		entries := m.currentHistoryEntries()
		if m.selectedHistory < len(entries) {
			entry := entries[m.selectedHistory]
			m.branch = entry.Branch
			m.inputs = entry.Inputs
			if m.inputs == nil {
				m.inputs = make(map[string]string)
			}
		}
	case PaneConfig:
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) initializeInputs(wf workflow.WorkflowFile) {
	m.inputs = make(map[string]string)
	for name, input := range wf.GetInputs() {
		m.inputs[name] = input.Default
	}
	m.selectedHistory = 0
}

func (m Model) currentHistoryEntries() []frecency.HistoryEntry {
	if m.history == nil {
		return nil
	}
	var workflowFilter string
	if m.selectedWorkflow < len(m.workflows) {
		workflowFilter = m.workflows[m.selectedWorkflow].Filename
	}
	return m.history.TopForRepo(m.repo, workflowFilter, 10)
}

// SelectedWorkflow returns the currently selected workflow.
func (m Model) SelectedWorkflow() *workflow.WorkflowFile {
	if m.selectedWorkflow >= len(m.workflows) {
		return nil
	}
	return &m.workflows[m.selectedWorkflow]
}

// View implements tea.Model.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	topHeight := m.height / 2
	bottomHeight := m.height - topHeight

	leftWidth := m.width / 3
	rightWidth := m.width - leftWidth

	workflowPane := m.viewWorkflowPane(leftWidth, topHeight)
	historyPane := m.viewHistoryPane(rightWidth, topHeight)
	configPane := m.viewConfigPane(m.width, bottomHeight)

	top := lipgloss.JoinHorizontal(lipgloss.Top, workflowPane, historyPane)
	return lipgloss.JoinVertical(lipgloss.Left, top, configPane)
}

func (m Model) viewWorkflowPane(width, height int) string {
	style := ui.PaneStyle(width, height, m.focused == PaneWorkflows)

	title := ui.TitleStyle.Render("Workflows")
	var content string
	for i, wf := range m.workflows {
		name := wf.Name
		if name == "" {
			name = wf.Filename
		}
		line := wf.Filename
		if name != wf.Filename {
			line = name + " (" + wf.Filename + ")"
		}
		if i == m.selectedWorkflow {
			content += ui.SelectedStyle.Render("> " + line)
		} else {
			content += ui.NormalStyle.Render("  " + line)
		}
		if i < len(m.workflows)-1 {
			content += "\n"
		}
	}

	return style.Render(title + "\n" + content)
}

func (m Model) viewHistoryPane(width, height int) string {
	style := ui.PaneStyle(width, height, m.focused == PaneHistory)

	var workflowName string
	if m.selectedWorkflow < len(m.workflows) {
		workflowName = m.workflows[m.selectedWorkflow].Filename
	}
	title := ui.TitleStyle.Render("Recent Runs (" + workflowName + ")")

	entries := m.currentHistoryEntries()
	var content string
	if len(entries) == 0 {
		content = ui.SubtitleStyle.Render("No history")
	} else {
		for i, e := range entries {
			line := e.Branch
			if i == m.selectedHistory {
				content += ui.SelectedStyle.Render("> " + line)
			} else {
				content += ui.NormalStyle.Render("  " + line)
			}
			if i < len(entries)-1 {
				content += "\n"
			}
		}
	}

	return style.Render(title + "\n" + content)
}

func (m Model) viewConfigPane(width, height int) string {
	style := ui.PaneStyle(width, height, m.focused == PaneConfig)

	title := ui.TitleStyle.Render("Configuration")

	var workflowLine, branchLine, inputsLine string

	if m.selectedWorkflow < len(m.workflows) {
		wf := m.workflows[m.selectedWorkflow]
		workflowLine = "Workflow: " + wf.Filename

		branch := m.branch
		if branch == "" {
			branch = "(not set)"
		}
		branchLine = "Branch: [b] " + branch

		inputs := wf.GetInputs()
		if len(inputs) > 0 {
			inputsLine = "\nInputs:"
			i := 1
			for name, input := range inputs {
				val := m.inputs[name]
				if val == "" {
					val = input.Default
				}
				if val == "" {
					val = "(empty)"
				}
				inputsLine += "\n  [" + string(rune('0'+i)) + "] " + name + ": " + val
				i++
				if i > 9 {
					break
				}
			}
		}
	}

	watchStatus := ""
	if m.watchRun {
		watchStatus = " [w] watch: on"
	} else {
		watchStatus = " [w] watch: off"
	}

	content := workflowLine + "\n" + branchLine + watchStatus + inputsLine

	helpLine := "\n\n" + ui.HelpStyle.Render("[Tab] switch pane  [Enter] run  [b] branch  [1-9] edit input  [q] quit")

	return style.Render(title + "\n" + content + helpLine)
}
