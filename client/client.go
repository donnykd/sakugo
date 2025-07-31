package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Post struct {
	ID            int    `json:"id"`
	Tags          string `json:"tags"`
	CreatedAt     int    `json:"created_at"`
	Source        string `json:"source"`
	Score         int    `json:"score"`
	FileExt       string `json:"file_ext"`
	FileURL       string `json:"file_url"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	PreviewURL    string `json:"preview_url"`
	PreviewWidth  int    `json:"preview_width"`
	PreviewHeight int    `json:"preview_height"`
}

type Artist struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Tag struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type PostConfig struct {
	Limit int
	Tags  []string
}

func validatePost(post Post) (Post, error) {
	if post.ID > 0 && post.FileURL != "" {
		return post, nil
	}
	return Post{}, fmt.Errorf("invalid post: ID=%d, FileURL='%s'", post.ID, post.FileURL)
}

func makeRequest(url string) ([]Post, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Could not get post: %v", err)
	}
	defer resp.Body.Close()

	var posts []Post
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		return nil, fmt.Errorf("Could not decode json: %v", err)
	}

	for i, post := range posts {
		if _, error := validatePost(post); error != nil {
			posts = append(posts[:i], posts[:i+1]...)
		}
	}
	return posts, nil
}

func FetchPosts(cfg PostConfig) ([]Post, error) {
	limit := 8
	if cfg.Limit != 0 {
		limit = cfg.Limit
	}

	url := fmt.Sprintf("https://www.sakugabooru.com/post.json?limit=%d", limit)

	if len(cfg.Tags) > 0 {
		tagString := strings.Join(cfg.Tags, "+")
		url += "&tags=" + tagString
	}

	return makeRequest(url)
}
