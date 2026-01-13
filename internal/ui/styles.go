package ui

import "github.com/charmbracelet/lipgloss"

// Colors used throughout the UI.
var (
	PrimaryColor   = lipgloss.Color("205")
	SecondaryColor = lipgloss.Color("240")
	AccentColor    = lipgloss.Color("86")
	MutedColor     = lipgloss.Color("245")
)

// Styles for the application.
var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	SelectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(AccentColor)

	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	HelpStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	BorderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(SecondaryColor)

	FocusedBorderStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(PrimaryColor)
)

// PaneStyle returns a style for a pane with optional focus.
func PaneStyle(width, height int, focused bool) lipgloss.Style {
	style := BorderStyle
	if focused {
		style = FocusedBorderStyle
	}
	return style.Width(width - 2).Height(height - 2)
}
