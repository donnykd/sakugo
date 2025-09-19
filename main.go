package main

import (
	"context"
	"fmt"
	"time"

	"github.com/donnykd/sakugo/client"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	posts, err := client.FetchPosts(ctx, client.PostConfig{
		Limit: 10,
		Tags:  []string{"order:score"},
	})
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	if len(posts) == 0 {
		fmt.Println("Nobody here but us chickens!")
	}

	for _, post := range posts {
		fmt.Println(post.Names[0].Name)
	}
}
