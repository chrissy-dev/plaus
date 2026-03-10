package app

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil
	case dataMsg:
		m.Loading = false
		m.Aggregate = msg.Aggregate
		m.Pages = msg.Pages
		m.Sources = msg.Sources
		m.TimeSeries = msg.TimeSeries
		return m, nil
	case errMsg:
		m.Loading = false
		m.Err = msg.Err
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			m.Loading = true
			m.Err = nil
			return m, m.fetchData
		}
	}
	return m, nil
}
