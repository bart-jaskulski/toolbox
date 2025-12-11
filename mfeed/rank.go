package main

import (
	"fmt"
	"sort"
)

type Ranking struct {
	Articles []FeedItem
}

func (r *Ranking) Len() int {
	return len(r.Articles)
}

// Rank ranks the items in the ranking
func (r *Ranking) Rank(cfg Config) error {
	if r.Len() == 0 {
		return fmt.Errorf("no items to rank")
	}

	for i := range r.Articles {
		r.Articles[i].ID = i
	}

	llm := NewLLMClient(cfg)
	scoredItems, err := llm.ScoreArticles(r.Articles)
	if err != nil {
		return fmt.Errorf("failed to score articles: %w", err)
	}

	r.Articles = scoredItems
	r.dropUninterestingItems(cfg)

	return nil
}

func (r *Ranking) sortItems() {
	sort.Slice(r.Articles, func(i, j int) bool {
		return r.Articles[i].Score > r.Articles[j].Score
	})
}

func (r *Ranking) dropUninterestingItems(cfg Config) {
	var filteredItems []FeedItem
	for _, item := range r.Articles {
		if item.Score >= cfg.MinimalScore {
			filteredItems = append(filteredItems, item)
		}
	}
	r.Articles = filteredItems
}

func (r *Ranking) pickTopItems(cfg Config) []FeedItem {
	topItems := r.Articles
	if r.Len() > cfg.MaxFeeds {
		r.sortItems()
		topItems = r.Articles[:cfg.MaxFeeds]
	}
	return topItems
}
