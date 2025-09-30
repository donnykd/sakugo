package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	_ "github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/donnykd/sakugo/client"
	"github.com/donnykd/sakugo/model"
)

var (
	highlight = lipgloss.AdaptiveColor{
		Light: "#FF4757",
		Dark:  "#FF4757",
	}
	title = lipgloss.AdaptiveColor{
		Light: "#FF4757",
		Dark:  "#FF4757",
	}
	option = lipgloss.AdaptiveColor{
		Light: "#A4B0BE",
		Dark:  "#A4B0BE",
	}
	bg = lipgloss.AdaptiveColor{
		Light: "#222222",
		Dark:  "#222222",
	}
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(title)
	optionStyle = lipgloss.NewStyle().Bold(true).Foreground(option)
	pageBorder  = lipgloss.Border{
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
		Border(pageBorder, true).BorderForeground(highlight).Padding(0, 2).BorderBackground(bg)
)

func cleanTab(s string) string {
	// Remove ANSI reset code and then resetForeground
	return strings.ReplaceAll(s, "\x1b[0m", "") + "\x1b[39m"
}

type Tui struct {
	model    *model.Model
	tabIndex int
	spinner  spinner.Model
}

func NewTui(m *model.Model) *Tui {
	return &Tui{
		model:    m,
		tabIndex: 0,
	}
}

func (t *Tui) Init() tea.Cmd {
	t.model.SetPosts()
	return nil
}

func (t *Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		}
	}
	return t, nil
}

func (t *Tui) View() string {
	switch t.model.ViewState {
	case model.PostsView:
		return t.renderPosts()
	}
	return ""
}

func (t *Tui) renderSearchBar() string {
	searchBar := page.Width(t.model.TerminalWidth - 2).Height(1).Background(bg).Render("")
	return searchBar
}

func (t *Tui) renderPage(content string) string {
	searchBar := t.renderSearchBar()
	page := page.Width(t.model.TerminalWidth - 2).
		Height(t.model.TerminalHeight - 5).Background(bg).Render(content)

	pageCombined := lipgloss.JoinVertical(lipgloss.Left, searchBar, page)
	layout := lipgloss.Place(
		t.model.TerminalWidth, t.model.TerminalHeight, lipgloss.Center, lipgloss.Bottom, pageCombined)
	return layout
}

func (t *Tui) cleanPostName(name string) string {
	return strings.TrimSuffix(strings.ReplaceAll(name, "_", " "), " series")
}

func (t *Tui) postTab(p client.Post) string {
	seen := make(map[string]bool)
	var postNames []string
	for _, name := range p.Names {
		if !seen[name.Name] {
			cleanedName := t.cleanPostName(name.Name)
			postNames = append(postNames, cleanedName)
			seen[name.Name] = true
		}
	}
	tabName := strings.Join(postNames, " • ")
	title := titleStyle.Render(tabName)
	metadata := lipgloss.NewStyle().Foreground(option).
		Render(fmt.Sprintf("ID: %d | Score: %d", p.ID, p.Score))
	tab := lipgloss.JoinVertical(lipgloss.Left, title, metadata)
	return cleanTab(tab)
}

func (t *Tui) renderPosts() string {
	postsList := t.model.Posts
	var createdTabs []string
	for _, post := range postsList {
		postStyle := lipgloss.NewStyle().
			Padding(1).
			Width(t.model.TerminalWidth - 6)
		tab := t.postTab(post)
		styledTab := postStyle.Render(tab)
		createdTabs = append(createdTabs, styledTab)
	}
	allTabs := strings.Join(createdTabs, "\n")
	centeredContent := lipgloss.NewStyle().Width(t.model.TerminalWidth).
		AlignHorizontal(lipgloss.Top).Render(allTabs)
	posts := t.renderPage(centeredContent)
	return posts
}
