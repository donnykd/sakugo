package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/donnykd/sakugo/model"
)

type tui struct {
	m *model.Model
}

func (t tui) Init() tea.Cmd {
	t.m.LoadHome()
	return nil
}

func (t tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return t, tea.Quit
		}
	}
	return t, nil
}

func (t tui) View() string {
	switch t.m.ViewState {
	case model.Home:
		return "Sakugo, The Sakugabooru TUI Client running locally!"
	}
	return ""
}

func main() {
	m := model.NewModel()
	program := tea.NewProgram(tui{m: m})
	program.Run()
}
