package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/chriswk/plaus/internal/api"
)

type Model struct {
	Site      string
	Client    *api.Client
	Aggregate api.Aggregate
	Pages     []api.PageStats
	Sources   []api.SourceStats
	Loading   bool
	Err       error
	Width     int
	Height    int
}

type dataMsg struct {
	Aggregate api.Aggregate
	Pages     []api.PageStats
	Sources   []api.SourceStats
}

type errMsg struct{ Err error }

func New(site string, client *api.Client) Model {
	return Model{
		Site:    site,
		Client:  client,
		Loading: true,
		Width:   80,
	}
}

func (m Model) Init() tea.Cmd {
	return m.fetchData
}

func (m Model) fetchData() tea.Msg {
	agg, err := m.Client.GetAggregate("30d")
	if err != nil {
		return errMsg{Err: err}
	}
	pages, err := m.Client.GetTopPages("30d", 10)
	if err != nil {
		return errMsg{Err: err}
	}
	sources, err := m.Client.GetTopSources("30d", 10)
	if err != nil {
		return errMsg{Err: err}
	}
	return dataMsg{
		Aggregate: agg,
		Pages:     pages,
		Sources:   sources,
	}
}
