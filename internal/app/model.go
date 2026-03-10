package app

import tea "github.com/charmbracelet/bubbletea"

type Model struct {
	Site string
}

func New(site string) Model {
	return Model{Site: site}
}

func (m Model) Init() tea.Cmd {
	return nil
}
