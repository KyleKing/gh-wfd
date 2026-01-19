package panes

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyleking/gh-lazydispatch/internal/config"
	"github.com/kyleking/gh-lazydispatch/internal/ui"
)

// ChainListModel manages the chain list display.
type ChainListModel struct {
	chains        map[string]config.Chain
	chainNames    []string
	selectedIndex int
	width         int
	height        int
	focused       bool
}

// NewChainListModel creates a new chain list model.
func NewChainListModel() ChainListModel {
	return ChainListModel{selectedIndex: 0}
}

// SetChains updates the chain definitions.
func (m *ChainListModel) SetChains(chains map[string]config.Chain) {
	m.chains = chains
	m.chainNames = make([]string, 0, len(chains))
	for name := range chains {
		m.chainNames = append(m.chainNames, name)
	}
	sort.Strings(m.chainNames)
}

// SetSize updates the pane dimensions.
func (m *ChainListModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetFocused updates the focus state.
func (m *ChainListModel) SetFocused(focused bool) {
	m.focused = focused
}

// MoveUp moves selection up.
func (m *ChainListModel) MoveUp() {
	if m.selectedIndex > 0 {
		m.selectedIndex--
	}
}

// MoveDown moves selection down.
func (m *ChainListModel) MoveDown() {
	if m.selectedIndex < len(m.chainNames)-1 {
		m.selectedIndex++
	}
}

// SelectedChain returns the currently selected chain.
func (m ChainListModel) SelectedChain() (string, config.Chain, bool) {
	if len(m.chainNames) == 0 {
		return "", config.Chain{}, false
	}
	name := m.chainNames[m.selectedIndex]
	return name, m.chains[name], true
}

// Update handles messages for the chain list.
func (m ChainListModel) Update(msg tea.Msg) (ChainListModel, tea.Cmd) {
	return m, nil
}

// ViewContent renders the chain list content without the pane border.
func (m ChainListModel) ViewContent() string {
	if len(m.chainNames) == 0 {
		var content strings.Builder
		content.WriteString(ui.SubtitleStyle.Render("No chains configured"))
		content.WriteString("\n\n")
		content.WriteString(ui.HelpStyle.Render("Add chains to"))
		content.WriteString("\n")
		content.WriteString(ui.HelpStyle.Render(".github/lazydispatch.yml"))
		return content.String()
	}

	var content strings.Builder
	for i, name := range m.chainNames {
		chain := m.chains[name]
		stepCount := len(chain.Steps)

		line := fmt.Sprintf("%s (%d steps)", name, stepCount)
		if chain.Description != "" {
			maxDescLen := m.width - len(line) - 6
			if maxDescLen > 10 {
				desc := chain.Description
				if len(desc) > maxDescLen {
					desc = desc[:maxDescLen-3] + "..."
				}
				line = fmt.Sprintf("%s (%d steps) - %s", name, stepCount, desc)
			}
		}

		if i == m.selectedIndex {
			content.WriteString(ui.SelectedStyle.Render("> " + line))
		} else {
			content.WriteString(ui.NormalStyle.Render("  " + line))
		}
		if i < len(m.chainNames)-1 {
			content.WriteString("\n")
		}
	}
	return content.String()
}

// View renders the chain list pane with border.
func (m ChainListModel) View() string {
	style := ui.PaneStyle(m.width, m.height, m.focused)
	title := ui.TitleStyle.Render("Chains")
	return style.Render(title + "\n" + m.ViewContent())
}
