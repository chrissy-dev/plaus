package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/chrissy-dev/plaus/internal/api"
)

type Period int

const (
	PeriodToday Period = iota
	PeriodYesterday
	PeriodWeek
	PeriodMonth
)

func (p Period) Label() string {
	switch p {
	case PeriodToday:
		return "Today"
	case PeriodYesterday:
		return "Yesterday"
	case PeriodWeek:
		return "Last 7 days"
	case PeriodMonth:
		return "Last 30 days"
	}
	return ""
}

func (p Period) DateRange() string {
	switch p {
	case PeriodToday:
		return "day"
	case PeriodYesterday:
		y := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		return y // will be wrapped as [y, y] in query
	case PeriodWeek:
		return "7d"
	case PeriodMonth:
		return "30d"
	}
	return "30d"
}

func (p Period) TimeDimension() string {
	switch p {
	case PeriodToday, PeriodYesterday:
		return "time:hour"
	default:
		return "time:day"
	}
}

type Model struct {
	Site             string
	Client           *api.Client
	Period           Period
	Aggregate        api.Aggregate
	Pages            []api.PageStats
	Sources          []api.SourceStats
	TimeSeries       []api.TimeSeriesPoint
	RealtimeVisitors int
	Loading          bool
	Err              error
	Width            int
	Height           int
	LiveTick         bool // toggles for blinking indicator
}

type dataMsg struct {
	Aggregate  api.Aggregate
	Pages      []api.PageStats
	Sources    []api.SourceStats
	TimeSeries []api.TimeSeriesPoint
}

type errMsg struct{ Err error }
type realtimeMsg struct{ Count int }
type realtimeErrMsg struct{ Err error }
type tickMsg time.Time

func New(site string, client *api.Client) Model {
	return Model{
		Site:    site,
		Client:  client,
		Period:  PeriodMonth,
		Loading: true,
		Width:   80,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.fetchData, m.tickCmd())
}

func (m Model) fetchData() tea.Msg {
	dateRange := m.Period.DateRange()

	// For yesterday, use a custom date range array
	var dr any = dateRange
	if m.Period == PeriodYesterday {
		dr = []string{dateRange, dateRange}
	}

	agg, err := m.Client.GetAggregate(dr)
	if err != nil {
		return errMsg{Err: err}
	}
	pages, err := m.Client.GetTopPages(dr, 10)
	if err != nil {
		return errMsg{Err: err}
	}
	sources, err := m.Client.GetTopSources(dr, 10)
	if err != nil {
		return errMsg{Err: err}
	}
	ts, err := m.Client.GetTimeSeries(dr, m.Period.TimeDimension())
	if err != nil {
		return errMsg{Err: err}
	}
	return dataMsg{
		Aggregate:  agg,
		Pages:      pages,
		Sources:    sources,
		TimeSeries: ts,
	}
}

func (m Model) fetchRealtime() tea.Msg {
	count, err := m.Client.GetRealtimeVisitors()
	if err != nil {
		return realtimeErrMsg{Err: err}
	}
	return realtimeMsg{Count: count}
}

func (m Model) tickCmd() tea.Cmd {
	return tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
