package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/donnykd/sakugo/model"
	"github.com/donnykd/sakugo/tui"
)

func main() {
	m := model.NewModel()
	tuiInstance := tui.NewTui(m)
	program := tea.NewProgram(tuiInstance, tea.WithAltScreen())
	program.Run()
}
