package app

import (
	"fmt"
	"strings"

	"github.com/NimbleMarkets/ntcharts/sparkline"
	"github.com/charmbracelet/lipgloss"

	"github.com/chriswk/plaus/internal/api"
)

var (
	purple    = lipgloss.Color("63")
	pink      = lipgloss.Color("212")
	grey      = lipgloss.Color("241")
	lightGrey = lipgloss.Color("245")
	white     = lipgloss.Color("255")
	red       = lipgloss.Color("196")
	dimWhite  = lipgloss.Color("250")
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(purple).
			PaddingLeft(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(grey).
			PaddingLeft(1)

	metricValueStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(white)

	metricLabelStyle = lipgloss.NewStyle().
				Foreground(lightGrey)

	panelHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(pink).
				PaddingBottom(1)

	rowNameStyle = lipgloss.NewStyle().
			Foreground(dimWhite)

	rowValueStyle = lipgloss.NewStyle().
			Foreground(lightGrey)

	errStyle = lipgloss.NewStyle().
			Foreground(red).
			PaddingLeft(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(grey).
			PaddingLeft(1)

	chartBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238")).
			Padding(0, 1)
)

func (m Model) View() string {
	w := m.Width
	if w < 40 {
		w = 40
	}

	var b strings.Builder

	b.WriteString(titleStyle.Render("Plaus") + "\n")
	b.WriteString(subtitleStyle.Render(m.Site+" · Last 30 days") + "\n\n")

	if m.Loading {
		b.WriteString(subtitleStyle.Render("Loading...") + "\n")
		return b.String()
	}

	if m.Err != nil {
		b.WriteString(errStyle.Render("Error: "+m.Err.Error()) + "\n\n")
		b.WriteString(helpStyle.Render("r retry · q quit") + "\n")
		return b.String()
	}

	b.WriteString(renderMetricCards(m.Aggregate, w))
	b.WriteString("\n\n")
	b.WriteString(renderChart(m.TimeSeries, w))
	b.WriteString("\n\n")
	b.WriteString(renderTwoPanels(m.Pages, m.Sources, w))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("r refresh · q quit") + "\n")

	return b.String()
}

func renderMetricCards(a api.Aggregate, width int) string {
	cards := []struct {
		label string
		value string
	}{
		{"VISITORS", formatNumber(a.Visitors)},
		{"VISITS", formatNumber(a.Visits)},
		{"PAGEVIEWS", formatNumber(a.Pageviews)},
		{"VIEWS/VISIT", fmt.Sprintf("%.1f", a.ViewsPerVisit)},
		{"BOUNCE RATE", fmt.Sprintf("%.0f%%", a.BounceRate)},
		{"VISIT DURATION", formatDuration(a.VisitDuration)},
	}

	cardWidth := (width - 2) / len(cards)
	if cardWidth < 12 {
		cardWidth = 12
	}

	cardStyle := lipgloss.NewStyle().
		Width(cardWidth).
		PaddingLeft(1).
		PaddingRight(1)

	rendered := make([]string, len(cards))
	for i, c := range cards {
		rendered[i] = cardStyle.Render(
			metricValueStyle.Render(c.value) + "\n" + metricLabelStyle.Render(c.label),
		)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
}

func renderChart(ts []api.TimeSeriesPoint, width int) string {
	if len(ts) == 0 {
		return ""
	}

	chartWidth := width - 6
	if chartWidth < 20 {
		chartWidth = 20
	}

	values := make([]float64, len(ts))
	for i, p := range ts {
		values[i] = float64(p.Visitors)
	}

	// Interpolate to fill chart width so sparkline spans full width
	interpolated := interpolate(values, chartWidth)

	sl := sparkline.New(chartWidth, 5)
	sl.Style = lipgloss.NewStyle().Foreground(purple)
	sl.PushAll(interpolated)
	sl.DrawBraille()

	// Date range labels
	firstDate := ts[0].Date
	lastDate := ts[len(ts)-1].Date
	dateLabel := lipgloss.NewStyle().Foreground(grey)
	padding := max(0, chartWidth-len(firstDate)-len(lastDate))
	labelLine := dateLabel.Render(firstDate) +
		strings.Repeat(" ", padding) +
		dateLabel.Render(lastDate)

	header := panelHeaderStyle.Render("Visitors")
	content := header + "\n" + sl.View() + "\n" + labelLine

	return chartBorder.Width(chartWidth).Render(content)
}

func interpolate(data []float64, targetLen int) []float64 {
	if len(data) == 0 {
		return nil
	}
	if len(data) == 1 || targetLen <= len(data) {
		return data
	}
	result := make([]float64, targetLen)
	for i := range result {
		pos := float64(i) * float64(len(data)-1) / float64(targetLen-1)
		low := int(pos)
		high := low + 1
		if high >= len(data) {
			high = len(data) - 1
		}
		frac := pos - float64(low)
		result[i] = data[low]*(1-frac) + data[high]*frac
	}
	return result
}

func renderTwoPanels(pages []api.PageStats, sources []api.SourceStats, width int) string {
	panelWidth := (width - 3) / 2
	if panelWidth < 30 {
		panelWidth = 30
	}

	pageRows := make([]listRow, len(pages))
	for i, p := range pages {
		pageRows[i] = listRow{Name: p.Page, Value: formatNumber(p.Visitors)}
	}

	sourceRows := make([]listRow, len(sources))
	for i, s := range sources {
		name := s.Source
		if name == "" {
			name = "(direct)"
		}
		sourceRows[i] = listRow{Name: name, Value: formatNumber(s.Visitors)}
	}

	left := renderPanel("Top Pages", panelWidth, pageRows)
	right := renderPanel("Top Sources", panelWidth, sourceRows)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, " ", right)
}

type listRow struct {
	Name  string
	Value string
}

func renderPanel(title string, width int, rows []listRow) string {
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("238")).
		Width(width - 2).
		Padding(0, 1)

	var b strings.Builder
	b.WriteString(panelHeaderStyle.Render(title) + "\n")

	if len(rows) == 0 {
		b.WriteString(lipgloss.NewStyle().Foreground(grey).Render("No data"))
		return border.Render(b.String())
	}

	nameWidth := width - 14
	if nameWidth < 10 {
		nameWidth = 10
	}

	for i, r := range rows {
		name := rowNameStyle.Render(fmt.Sprintf("%-*s", nameWidth, truncate(r.Name, nameWidth)))
		value := rowValueStyle.Render(fmt.Sprintf("%6s", r.Value))
		b.WriteString(name + " " + value)
		if i < len(rows)-1 {
			b.WriteString("\n")
		}
	}

	return border.Render(b.String())
}

func formatDuration(seconds int) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	return fmt.Sprintf("%dm %ds", seconds/60, seconds%60)
}

func formatNumber(n int) string {
	if n >= 1_000_000 {
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	}
	if n >= 1_000 {
		return fmt.Sprintf("%.1fk", float64(n)/1_000)
	}
	return fmt.Sprintf("%d", n)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
