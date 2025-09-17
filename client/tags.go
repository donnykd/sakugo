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
	for tagName := range strings.SplitSeq(p.Tags, " ") {
		if tagName == "" {
			continue
		}

		tag, err := getTagByName(ctx, tagName)
		if err != nil {
			return fmt.Errorf("could not get tag: %v", err)
		}

		switch tag.Type {
		case GeneralTag:
			p.General = append(p.General, tag)
		case ArtistTag:
			p.Artists = append(p.Artists, tag)
		case NameTag:
			p.Names = append(p.Names, tag)
		case StyleTag:
			p.Style = append(p.Style, tag)
		case MetaTag:
			p.Meta = append(p.Meta, tag)
		default:
			p.General = append(p.General, tag)
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
	reqURL := fmt.Sprintf("https://www.sakugabooru.com/tag.json?name=%s", encodedTag)
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
	globalCache.Set(foundTag.Name, foundTag)

	return foundTag, nil
}
