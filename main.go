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

func main() {
	resp, err := http.Get("https://www.sakugabooru.com/post.json?tags=id:217484")
	if err != nil {
		fmt.Printf("Could not load post: %v", err)
	}
	defer resp.Body.Close()

	var posts []Post
	json.NewDecoder(resp.Body).Decode(&posts)
	

	if len(posts) == 1{
		post := posts[0]
		fmt.Printf("ID: %d, Tags: %s, CreatedAt: %d, Source: %s, Score: %d, FileURL: %s",
			post.ID, post.Tags, post.CreatedAt, post.Source, post.Score, post.FileURL)
	}
}
