package main

import (
	"os"
	"runtime/pprof"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/donnykd/sakugo/model"
	"github.com/donnykd/sakugo/tui"
)

func main() {
	file, _ := os.Create("./cpu.pprof")
	pprof.StartCPUProfile(file)
	defer pprof.StopCPUProfile()

	m := model.NewModel()
	tuiInstance := tui.NewTui(m)
	program := tea.NewProgram(tuiInstance, tea.WithAltScreen())
	program.Run()
}
