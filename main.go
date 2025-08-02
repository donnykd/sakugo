package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/donnykd/sakugo/model"
)

var (
	tabs = []string{"Home", "Posts", "Search", "Tags"}

	highlight = lipgloss.AdaptiveColor{
		Light: "#DD4B5F",
		Dark:  "#f6546a",
	}

	option = lipgloss.AdaptiveColor{
		Light: "#ec9da8",
		Dark:  "#ec9da8",
	}

	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(highlight)
	optionStyle = lipgloss.NewStyle().Bold(true).Foreground(option)

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
		switch msg.String() {
		case "q", "ctrl+c":
			return t, tea.Quit
		case "right":
			t.tabIndex = (t.tabIndex + 1) % len(tabs)
		case "left":
			t.tabIndex = (t.tabIndex - 1 + len(tabs)) % len(tabs)
		}
	}
	return t, nil
}

func (t tui) renderTabs() string {
	var renderedTabs []string

	for i, tab := range tabs {
		var selectedText string
		if i == t.tabIndex {
			selectedText = fmt.Sprintf("[ %s ]", tab)
			renderedTabs = append(renderedTabs, titleStyle.Render(selectedText))
		} else {
			selectedText = fmt.Sprintf("  %s  ", tab)
			renderedTabs = append(renderedTabs, optionStyle.Render(selectedText))
		}
	}

	tabLine := lipgloss.JoinHorizontal(lipgloss.Center, renderedTabs...)
	return lipgloss.NewStyle().Width(t.model.TerminalWidth).AlignHorizontal(lipgloss.Center).Render(tabLine)
}

func (t tui) renderPage(content string) string {
	tabs := t.renderTabs()
	tabHeight := lipgloss.Height(tabs)

	page := page.Width(t.model.TerminalWidth - 2).Height(t.model.TerminalHeight - tabHeight - 2).Render(content)
	layout := lipgloss.JoinVertical(lipgloss.Left, tabs, page)

	return layout
}

func (t tui) renderHome() string {
	title := titleStyle.Render("Sakugo - Sakugabooru TUI Client")
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
