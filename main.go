package main

import (
	"fmt"

	"github.com/donnykd/sakugo/client"
)

func main() {
	posts, err := client.FetchPosts(client.PostConfig{
		Limit: 5,
		Tags:  []string{"vincent_chansard", "order:score", "yen_bm"},
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
