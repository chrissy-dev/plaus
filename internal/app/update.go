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
		// If on today view, also kick off a realtime fetch
		if m.Period == PeriodToday {
			return m, m.fetchRealtime
		}
		return m, nil
	case errMsg:
		m.Loading = false
		m.Err = msg.Err
		return m, nil
	case realtimeMsg:
		m.RealtimeVisitors = msg.Count
		return m, nil
	case realtimeErrMsg:
		// Silently ignore realtime errors — don't break the UI
		return m, nil
	case tickMsg:
		m.LiveTick = !m.LiveTick
		if m.Period == PeriodToday {
			return m, tea.Batch(m.fetchRealtime, m.tickCmd())
		}
		return m, m.tickCmd()
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			return m, m.switchPeriod(m.Period)
		case "1":
			return m, m.switchPeriod(PeriodToday)
		case "2":
			return m, m.switchPeriod(PeriodYesterday)
		case "3":
			return m, m.switchPeriod(PeriodWeek)
		case "4":
			return m, m.switchPeriod(PeriodMonth)
		}
	}
	return m, nil
}

func (m *Model) switchPeriod(p Period) tea.Cmd {
	m.Period = p
	m.Loading = true
	m.Err = nil
	m.RealtimeVisitors = 0
	return m.fetchData
}
