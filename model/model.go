package model

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/donnykd/sakugo/client"
)

type ViewState int

const (
	LoadingView ViewState = iota
	PostsView
	ErrorView
)

type Model struct {
	Posts        []client.Post
	CurrentIndex int
	SearchConfig client.PostConfig
	ViewState    ViewState
	ErrorMessage string

	TerminalWidth  int
	TerminalHeight int
}

func NewModel() *Model {
	return &Model{
		Posts:        make([]client.Post, 0),
		CurrentIndex: 0,
		SearchConfig: client.PostConfig{
			Limit: 5,
			Tags:  []string{"order:score"},
		},
		ViewState: LoadingView,
	}
}

func (m *Model) SetPosts() tea.Cmd {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	posts, err := client.FetchPosts(ctx, m.SearchConfig)
	if err != nil {
		m.SetError("error")
	}

	m.Posts = posts
	m.ViewState = PostsView
	return nil
}

func (m *Model) SetError(err string) {
	m.ErrorMessage = err
	m.ViewState = ErrorView
}

func FetchPosts(config client.PostConfig) ([]client.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return client.FetchPosts(ctx, config)
}
