package app

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	siteStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	sectionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245")).PaddingLeft(2)
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1)
)

func (m Model) View() string {
	s := titleStyle.Render("Plaus") + "\n"
	s += siteStyle.Render("Site: "+m.Site) + "\n\n"
	s += sectionStyle.Render("Overview (placeholder)") + "\n"
	s += sectionStyle.Render("Top Pages (placeholder)") + "\n"
	s += sectionStyle.Render("Sources (placeholder)") + "\n\n"
	s += helpStyle.Render("Press q to quit")
	return s + "\n"
}
