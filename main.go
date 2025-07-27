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

type Artist struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Tags struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func main() {
	postResp, _ := http.Get("https://www.sakugabooru.com/post.json?tags=id:217484")
	defer postResp.Body.Close()

	var posts []Post
	json.NewDecoder(postResp.Body).Decode(&posts)

	post := posts[0]
	fmt.Printf("ID: %d, Tags: %s, CreatedAt: %d, Source: %s, Score: %d, FileURL: %s \n \n",
		post.ID, post.Tags, post.CreatedAt, post.Source, post.Score, post.FileURL)

	artistResp, _ := http.Get("https://www.sakugabooru.com/artist.json?name=shingo_yamashita")
	defer artistResp.Body.Close()

	var artists []Artist
	json.NewDecoder(artistResp.Body).Decode(&artists)

	artist := artists[0]
	fmt.Printf("ID: %d, Name: %s \n \n", artist.ID, artist.Name)

	tagResp, _ := http.Get("https://www.sakugabooru.com/tag.json?name=kanada_light_flare")
	defer tagResp.Body.Close()

	var tags []Tags
	json.NewDecoder(tagResp.Body).Decode(&tags)

	tag := tags[0]
	fmt.Printf("ID: %d, Name: %s, Count: %d \n", tag.ID, tag.Name, tag.Count)
}
