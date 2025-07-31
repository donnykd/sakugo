package model

import "github.com/donnykd/sakugo/client"

type ViewState int

const (
	Home ViewState = iota
	PostsView
	PostDetail
	Loading
	Error
)

type Model struct {
	Posts        []client.Post
	CurrentIndex int
	SearchConfig client.PostConfig
	isLoading    bool
	ViewState    ViewState
	ErrorMessage string

	// pagination
	HasMorePages bool
	CurrentPage  int
}

func (m *Model) NewModel() (*Model, error) {
	model := Model{
		Posts:        make([]client.Post, 0),
		CurrentIndex: 0,
		SearchConfig: client.PostConfig{},
		isLoading:    false,
		ViewState:    Home,
		ErrorMessage: "",
		HasMorePages: false,
		CurrentPage:  1,
	}

	return &model, nil
}
