// Package client provides a client for the sakugabooru API
//
// This package allows fetching posts and processing different tags
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Post struct {
	ID            int    `json:"id"`
	Tags          string `json:"tags"`
	Names         []Tag  `json:"-"`
	Artists       []Tag  `json:"-"`
	Style         []Tag  `json:"-"`
	Meta          []Tag  `json:"-"`
	General       []Tag  `json:"-"`
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

type Tag struct {
	Name string `json:"name"`
	Type int    `json:"type"`
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

func makeRequest(ctx context.Context, reqURL string) ([]Post, error) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %v", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not get post: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var posts []Post
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		return nil, fmt.Errorf("could not decode JSON: %v", err)
	}

	var filteredPosts []Post
	for _, post := range posts {
		if validPost, err := validatePost(post); err == nil {
			if err = validPost.setTags(ctx); err != nil {
				return nil, fmt.Errorf("could not set tags for post '%d' error: %v", validPost.ID, err)
			}
			filteredPosts = append(filteredPosts, validPost)
		}
	}

	return filteredPosts, nil
}

func FetchPosts(ctx context.Context, cfg PostConfig) ([]Post, error) {
	limit := 8
	if cfg.Limit > 0 {
		limit = cfg.Limit
	}

	reqURL := fmt.Sprintf("https://www.sakugabooru.com/post.json?limit=%d", limit)

	if len(cfg.Tags) > 0 {
		tagString := strings.Join(cfg.Tags, "+")
		fmt.Println(tagString)
		reqURL += "&tags=" + tagString
	}

	return makeRequest(ctx, reqURL)
}
