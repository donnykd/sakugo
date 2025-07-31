package client

import (
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
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Could not get post: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var posts []Post
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		return nil, fmt.Errorf("Could not decode JSON: %v", err)
	}

	var filteredPosts []Post
	for _, post := range posts {
		if validPost, error := validatePost(post); error == nil {
			filteredPosts = append(filteredPosts, validPost)
		}
	}
	return filteredPosts, nil
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
