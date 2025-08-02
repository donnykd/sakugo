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
		Top:      "─",
		Left:     "│",
		Right:    "│",
		TopLeft:  "╭",
		TopRight: "╮",
	}

	tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		BorderForeground(highlight).
		Padding(0, 3)

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
	model    *model.Model
	tabIndex int
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

func (t tui) renderTabs() string {
	homeTab := tab.Render("Home")
	postsTab := tab.Render("Posts")
	searchTab := tab.Render("Search")
	tagsTab := tab.Render("Tags")

	tabs := lipgloss.JoinHorizontal(lipgloss.Top, homeTab, postsTab, searchTab, tagsTab)
	centeredTabs := lipgloss.NewStyle().Width(t.model.TerminalWidth).AlignHorizontal(lipgloss.Center).Render(tabs)
	return centeredTabs
}

func (t tui) renderPage(content string) string {
	tabs := t.renderTabs()
	tabHeight := lipgloss.Height(tabs)

	page := page.Width(t.model.TerminalWidth - 2).Height(t.model.TerminalHeight - tabHeight - 2).Render(content)
	layout := lipgloss.JoinVertical(lipgloss.Left, tabs, page)

	return layout
}

func (t tui) renderHome() string {
	title := lipgloss.NewStyle().Bold(true).Foreground(highlight).Render("Sakugo - Sakugabooru TUI Client")
	centeredTitle := lipgloss.NewStyle().Width(t.model.TerminalWidth).AlignHorizontal(lipgloss.Center).Render(title)

	content := lipgloss.JoinVertical(lipgloss.Left, "", centeredTitle, "", "Press a key to navigate...")

	home := t.renderPage(content)

	return home
}

func (t tui) View() string {
	switch t.model.ViewState {
	case model.Home:
		return t.renderHome()
	}
	return ""
}

func main() {
	m := model.NewModel()
	program := tea.NewProgram(tui{model: m}, tea.WithAltScreen())
	program.Run()
}
