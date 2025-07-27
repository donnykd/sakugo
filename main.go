package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Post struct {
	ID            int    `json:"id"`
	Tags          string `json:"tags"`
	CreatedAt     int    `json:"created_at"`
	Source        string `json:"source"`
	Score         int    `json:"score"`
	FileEXT       string `json:"file_ext"`
	FileURL       string `json:"file_url"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	PreviewURL    string `json:"preview_url"`
	PreviewWidth  int    `json:"preview_width"`
	PreviewHeight int    `json:"preview_height"`
}

type PostConfig struct {
	Limit int
	Tag   string
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

func fetchPosts(cfg PostConfig) ([]Post, error) {
	limit := 8
	if cfg.Limit != 0 {
		limit = cfg.Limit
	}

	var artist string
	if cfg.Tag != "" {
		artist = cfg.Tag
	}

	url := fmt.Sprintf("https://www.sakugabooru.com/post.json?limit=%d", limit)
	if artist != "" {
		url += "&tags=" + artist
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Could not get post: %v", err)
	}
	defer resp.Body.Close()

	var posts []Post
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		return nil, fmt.Errorf("Could not decode json: %v", err)
	}

	return posts, nil
}

func main() {
	posts, err := fetchPosts(PostConfig{
		Limit: 5,
		Tag:   "shingo_yamashita", //multiple tags
		// order
	})
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	if len(posts) == 0 {
		fmt.Println("Nobody here but us chickens!")
	}

	for _, post := range posts {
		fmt.Printf("ID: %d, Tags: %s, CreatedAt: %d, Source: %s, Score: %d, FileURL: %s \n \n",
			post.ID, post.Tags, post.CreatedAt, post.Source, post.Score, post.FileURL)
	}
}
