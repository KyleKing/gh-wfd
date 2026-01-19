package panes

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyleking/gh-lazydispatch/internal/github"
	"github.com/kyleking/gh-lazydispatch/internal/ui"
	"github.com/kyleking/gh-lazydispatch/internal/watcher"
)

// LiveRunsModel manages the live runs display.
type LiveRunsModel struct {
	runs          []watcher.WatchedRun
	selectedIndex int
	width         int
	height        int
	focused       bool
}

// NewLiveRunsModel creates a new live runs model.
func NewLiveRunsModel() LiveRunsModel {
	return LiveRunsModel{selectedIndex: 0}
}

// SetRuns updates the list of watched runs.
func (m *LiveRunsModel) SetRuns(runs []watcher.WatchedRun) {
	m.runs = runs
	if m.selectedIndex >= len(runs) && len(runs) > 0 {
		m.selectedIndex = len(runs) - 1
	}
}

// SetSize updates the pane dimensions.
func (m *LiveRunsModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetFocused updates the focus state.
func (m *LiveRunsModel) SetFocused(focused bool) {
	m.focused = focused
}

// MoveUp moves selection up.
func (m *LiveRunsModel) MoveUp() {
	if m.selectedIndex > 0 {
		m.selectedIndex--
	}
}

// MoveDown moves selection down.
func (m *LiveRunsModel) MoveDown() {
	if m.selectedIndex < len(m.runs)-1 {
		m.selectedIndex++
	}
}

// SelectedRun returns the currently selected run.
func (m LiveRunsModel) SelectedRun() (watcher.WatchedRun, bool) {
	if len(m.runs) == 0 || m.selectedIndex >= len(m.runs) {
		return watcher.WatchedRun{}, false
	}
	return m.runs[m.selectedIndex], true
}

// SelectedIndex returns the current selection index.
func (m LiveRunsModel) SelectedIndex() int {
	return m.selectedIndex
}

// RunCount returns the number of runs.
func (m LiveRunsModel) RunCount() int {
	return len(m.runs)
}

// Update handles messages for the live runs model.
func (m LiveRunsModel) Update(msg tea.Msg) (LiveRunsModel, tea.Cmd) {
	return m, nil
}

// ViewContent renders the live runs content without the pane border.
func (m LiveRunsModel) ViewContent() string {
	if len(m.runs) == 0 {
		var content strings.Builder
		content.WriteString(ui.SubtitleStyle.Render("No active runs"))
		content.WriteString("\n\n")
		content.WriteString(ui.HelpStyle.Render("Runs appear here when"))
		content.WriteString("\n")
		content.WriteString(ui.HelpStyle.Render("Watch is enabled"))
		return content.String()
	}

	var content strings.Builder
	for i, run := range m.runs {
		icon := runStatusIcon(run.Status, run.Conclusion)
		line := fmt.Sprintf("%s %s", icon, run.Workflow)

		if run.Status != "" && run.Status != github.StatusCompleted {
			line = fmt.Sprintf("%s (%s)", line, run.Status)
		} else if run.Conclusion != "" {
			line = fmt.Sprintf("%s (%s)", line, run.Conclusion)
		}

		if i == m.selectedIndex {
			content.WriteString(ui.SelectedStyle.Render("> " + line))
		} else {
			content.WriteString(ui.NormalStyle.Render("  " + line))
		}

		if i < len(m.runs)-1 {
			content.WriteString("\n")
		}
	}
	return content.String()
}

// View renders the live runs pane with border.
func (m LiveRunsModel) View() string {
	style := ui.PaneStyle(m.width, m.height, m.focused)
	title := ui.TitleStyle.Render("Live Runs")
	return style.Render(title + "\n" + m.ViewContent())
}

func runStatusIcon(status, conclusion string) string {
	switch status {
	case github.StatusQueued:
		return "o"
	case github.StatusInProgress:
		return "*"
	case github.StatusCompleted:
		switch conclusion {
		case github.ConclusionSuccess:
			return "+"
		case github.ConclusionFailure:
			return "x"
		case github.ConclusionCancelled:
			return "-"
		default:
			return "?"
		}
	default:
		return "?"
	}
}

// ActiveCount returns the number of active runs.
func (m LiveRunsModel) ActiveCount() int {
	count := 0
	for _, run := range m.runs {
		if run.IsActive() {
			count++
		}
	}
	return count
}
