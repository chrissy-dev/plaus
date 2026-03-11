package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/chrissy-dev/plaus/internal/api"
	"github.com/chrissy-dev/plaus/internal/config"
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

type GraphType int

const (
	GraphLine GraphType = iota
	GraphBar
)

type Model struct {
	Site             string
	Client           *api.Client
	Period           Period
	Graph            GraphType
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

func New(site string, client *api.Client, graphType GraphType, period Period) Model {
	return Model{
		Site:    site,
		Client:  client,
		Period:  period,
		Graph:   graphType,
		Loading: true,
		Width:   80,
	}
}

func (m *Model) savePrefs() {
	cfg, err := config.Load()
	if err != nil {
		return
	}
	cfg.GraphType = GraphTypeToString(m.Graph)
	cfg.Period = PeriodToString(m.Period)
	config.Save(cfg)
}

func (m *Model) ToggleGraph() {
	if m.Graph == GraphLine {
		m.Graph = GraphBar
	} else {
		m.Graph = GraphLine
	}
	m.savePrefs()
}

func GraphTypeFromString(s string) GraphType {
	if s == "bar" {
		return GraphBar
	}
	return GraphLine
}

func GraphTypeToString(g GraphType) string {
	if g == GraphBar {
		return "bar"
	}
	return "line"
}

func PeriodFromString(s string) Period {
	switch s {
	case "today":
		return PeriodToday
	case "yesterday":
		return PeriodYesterday
	case "7d":
		return PeriodWeek
	default:
		return PeriodMonth
	}
}

func PeriodToString(p Period) string {
	switch p {
	case PeriodToday:
		return "today"
	case PeriodYesterday:
		return "yesterday"
	case PeriodWeek:
		return "7d"
	default:
		return "30d"
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
