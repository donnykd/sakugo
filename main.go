package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/donnykd/sakugo/model"
)

var (
	highlight = lipgloss.AdaptiveColor{
		Light: "#DD4B5F",
		Dark:  "#f6546a",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "|",
		Right:       "|",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "",
		Left:        "|",
		Right:       "|",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		BorderForeground(highlight).
		Padding(0, 2)

	activeTab = tab.Border(activeTabBorder, true)

	tabGap = tab.
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)
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

func (t tui) renderHomePage() string {
	homeTab := activeTab.Render("Home")
	postsTab := tab.Render("Posts")
	searchTab := tab.Render("Search")
	tagsTab := tab.Render("Tags")

	tabs := lipgloss.JoinHorizontal(lipgloss.Top, homeTab, postsTab, searchTab, tagsTab)

	gap := tabGap.Render(strings.Repeat(" ", max(0, 96-lipgloss.Width(tabs)-2)))

	title := lipgloss.NewStyle().Bold(true).Foreground(highlight).Render("Sakugo - Sakugabooru TUI Client")

	content := lipgloss.JoinVertical(lipgloss.Left, tabs, "", title, "", gap, "Press a key to navigate...")

	return content

}

func (t tui) View() string {
	switch t.m.ViewState {
	case model.Home:
		return t.renderHomePage()
	}
	return ""
}

func main() {
	m := model.NewModel()
	program := tea.NewProgram(tui{m: m})
	program.Run()
}
