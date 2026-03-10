package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	siteStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	headerStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).MarginTop(1)
	metricLabel  = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	metricValue  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("255"))
	rowDim       = lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	rowVal       = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	errStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1)
)

func (m Model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Plaus") + "\n")
	b.WriteString(siteStyle.Render("Site: "+m.Site+" · Last 30 days") + "\n")

	if m.Loading {
		b.WriteString("\n  Loading...\n")
		return b.String()
	}

	if m.Err != nil {
		b.WriteString("\n" + errStyle.Render("Error: "+m.Err.Error()) + "\n")
		b.WriteString(helpStyle.Render("Press r to retry · q to quit") + "\n")
		return b.String()
	}

	// Overview
	b.WriteString(headerStyle.Render("  Overview") + "\n")
	a := m.Aggregate
	stats := []struct{ label, value string }{
		{"Visitors", fmt.Sprintf("%d", a.Visitors)},
		{"Visits", fmt.Sprintf("%d", a.Visits)},
		{"Pageviews", fmt.Sprintf("%d", a.Pageviews)},
		{"Views/Visit", fmt.Sprintf("%.1f", a.ViewsPerVisit)},
		{"Bounce Rate", fmt.Sprintf("%.0f%%", a.BounceRate)},
		{"Visit Duration", formatDuration(a.VisitDuration)},
	}
	for _, s := range stats {
		b.WriteString(fmt.Sprintf("  %s %s\n",
			metricLabel.Render(fmt.Sprintf("%-14s", s.label)),
			metricValue.Render(s.value),
		))
	}

	// Top Pages
	b.WriteString(headerStyle.Render("  Top Pages") + "\n")
	for _, p := range m.Pages {
		b.WriteString(fmt.Sprintf("  %s %s\n",
			rowDim.Render(fmt.Sprintf("%-40s", truncate(p.Page, 40))),
			rowVal.Render(fmt.Sprintf("%d visitors", p.Visitors)),
		))
	}

	// Sources
	b.WriteString(headerStyle.Render("  Top Sources") + "\n")
	for _, s := range m.Sources {
		source := s.Source
		if source == "" {
			source = "(direct)"
		}
		b.WriteString(fmt.Sprintf("  %s %s\n",
			rowDim.Render(fmt.Sprintf("%-40s", truncate(source, 40))),
			rowVal.Render(fmt.Sprintf("%d visitors", s.Visitors)),
		))
	}

	b.WriteString(helpStyle.Render("  r refresh · q quit") + "\n")
	return b.String()
}

func formatDuration(seconds int) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	return fmt.Sprintf("%dm %ds", seconds/60, seconds%60)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
