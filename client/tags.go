package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	GeneralTag = 0
	ArtistTag  = 1
	NameTag    = 3
	StyleTag   = 4
	MetaTag    = 5
)

var globalCache = newTagCache()

type TagCache struct {
	cache map[string]Tag
	mu    sync.RWMutex
}

type tagResult struct {
	tag Tag
	err error
}

func newTagCache() *TagCache {
	return &TagCache{
		cache: make(map[string]Tag),
	}
}

func (tc *TagCache) Set(key string, tag Tag) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.cache[key] = tag
}

func (tc *TagCache) Get(key string) (Tag, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	tag, exists := tc.cache[key]
	return tag, exists
}

func (p *Post) setTags(ctx context.Context) error {
	tagNames := strings.Fields(p.Tags)
	if len(tagNames) == 0 {
		return nil
	}

	ch := make(chan tagResult, len(tagNames))
	var wg sync.WaitGroup
	wg.Add(len(tagNames))

	for _, tagName := range tagNames {
		go func(name string) {
			defer wg.Done()
			tag, err := getTagByName(ctx, name)
			ch <- tagResult{tag, err}
		}(tagName)
	}

	wg.Wait()
	close(ch)

	for result := range ch {
		if result.err != nil {
			return fmt.Errorf("could not get tag: %v", result.err)
		}

		switch result.tag.Type {
		case GeneralTag:
			p.General = append(p.General, result.tag)
		case ArtistTag:
			p.Artists = append(p.Artists, result.tag)
		case NameTag:
			p.Names = append(p.Names, result.tag)
		case StyleTag:
			p.Style = append(p.Style, result.tag)
		case MetaTag:
			p.Meta = append(p.Meta, result.tag)
		default:
			p.General = append(p.General, result.tag)
		}
	}
	return nil
}

func getTagByName(ctx context.Context, tagName string) (Tag, error) {
	if cachedTag, found := globalCache.Get(tagName); found {
		return cachedTag, nil
	}

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	encodedTag := url.QueryEscape(tagName)
	reqURL := fmt.Sprintf("https://www.sakugabooru.com/tag.json?name=%s&order=count&limit=1", encodedTag)
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return Tag{}, fmt.Errorf("could not create request: %v", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return Tag{}, fmt.Errorf("could not get tag: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return Tag{}, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var tags []Tag
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return Tag{}, fmt.Errorf("could not decode JSON: %v", err)
	}

	if len(tags) == 0 {
		return Tag{}, fmt.Errorf("tag '%s' not found", tagName)
	}

	foundTag := tags[0]
	fmt.Println(foundTag)
	globalCache.Set(foundTag.Name, foundTag)

	return foundTag, nil
}
