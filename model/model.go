package model

import (
	"github.com/donnykd/sakugo/client"
)

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

func NewModel() *Model {
	return &Model{
		Posts:        make([]client.Post, 0),
		CurrentIndex: 0,
		ViewState:    Home,
		CurrentPage:  1,
	}
}

func (m *Model) GoHome() {
	m.ViewState = Home
}

func (m *Model) LoadPosts() {
	m.ViewState = Loading
	
	posts, err := client.FetchPosts(m.SearchConfig)
	
	if err != nil{
		m.ViewState = Error
		m.ErrorMessage = "Nobody here but us chickens!"
		return
	}
	
	m.Posts = posts
	m.CurrentIndex = 0
	m.ViewState = PostsView
}
