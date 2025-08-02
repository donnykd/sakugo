package main

import (
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
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		BorderForeground(highlight).
		Padding(0, 3)

	activeTab = tab.Border(activeTabBorder, true)

	pageBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}

	page = lipgloss.NewStyle().
		Border(pageBorder, true).BorderForeground(highlight).Padding(0, 2)
)

type tui struct {
	model *model.Model
}

func (t tui) Init() tea.Cmd {
	t.model.LoadHome()
	return nil
}

func (t tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.model.TerminalHeight = msg.Height
		t.model.TerminalWidth = msg.Width

		if t.model.TerminalWidth < 70 {
			t.model.TerminalWidth = 70
		}
		if t.model.TerminalHeight < 20 {
			t.model.TerminalHeight = 20
		}
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
	centeredTabs := lipgloss.NewStyle().Width(t.model.TerminalWidth).AlignHorizontal(lipgloss.Center).Render(tabs)

	title := lipgloss.NewStyle().Bold(true).Foreground(highlight).Render("Sakugo - Sakugabooru TUI Client")
	centeredTitle := lipgloss.NewStyle().Width(t.model.TerminalWidth).AlignHorizontal(lipgloss.Center).Render(title)

	content := lipgloss.JoinVertical(lipgloss.Left, centeredTabs, "", centeredTitle, "", "Press a key to navigate...")

	page := page.Width(t.model.TerminalWidth - 2).Height(t.model.TerminalHeight - 2).Render(content)

	return page
}

func (t tui) View() string {
	switch t.model.ViewState {
	case model.Home:
		return t.renderHomePage()
	}
	return ""
}

func main() {
	m := model.NewModel()
	program := tea.NewProgram(tui{model: m}, tea.WithAltScreen())
	program.Run()
}
