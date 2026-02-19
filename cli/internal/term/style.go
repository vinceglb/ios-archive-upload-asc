package term

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	titleStyle   lipgloss.Style
	sectionStyle lipgloss.Style
	labelStyle   lipgloss.Style
	valueStyle   lipgloss.Style
	mutedStyle   lipgloss.Style
	errorStyle   lipgloss.Style
	successStyle lipgloss.Style
}

func NewTheme() Theme {
	return Theme{
		titleStyle:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")),
		sectionStyle: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63")),
		labelStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("246")),
		valueStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("252")),
		mutedStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
		errorStyle:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196")),
		successStyle: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("82")),
	}
}

func (t Theme) Title(value string) string {
	return t.titleStyle.Render(value)
}

func (t Theme) Section(value string) string {
	return t.sectionStyle.Render(value)
}

func (t Theme) Label(value string) string {
	return t.labelStyle.Render(value)
}

func (t Theme) Value(value string) string {
	return t.valueStyle.Render(value)
}

func (t Theme) Muted(value string) string {
	return t.mutedStyle.Render(value)
}

func (t Theme) Error(value string) string {
	return t.errorStyle.Render("[error] " + value)
}

func (t Theme) Success(value string) string {
	return t.successStyle.Render(value)
}
