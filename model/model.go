package model

import (
	"github.com/donnykd/sakugo/client"
)

type ViewState int

const (
	HomeView ViewState = iota
	PostsView
	SearchView
	TagsView
	PostDetail
	Loading
	Error
)

type Model struct {
	Posts        []client.Post
	CurrentIndex int
	SearchConfig client.PostConfig
	ViewState    ViewState
	ErrorMessage string

	// pagination
	HasMorePages bool
	CurrentPage  int

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
		ViewState:   HomeView,
		CurrentPage: 1,
	}
}

func (m *Model) Loading() {
	m.ViewState = Loading
}

func (m *Model) LoadHome() {
	m.ViewState = HomeView
}

func (m *Model) LoadPosts(err error) {
	posts, err := client.FetchPosts(m.SearchConfig)

	if err != nil {
		m.ViewState = Error
		m.ErrorMessage = "Nobody here but us chickens!"
		return
	}

	m.Posts = posts
	m.CurrentIndex = 0
	m.ViewState = PostsView
}
