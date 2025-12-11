package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"time"

	readability "github.com/go-shiori/go-readability"
)

type FeedItem struct {
	ID      int
	Title   string
	Link    string
	Content string
	Source  string
	Updated *time.Time
	Score   int
}

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Xmlns   string   `xml:"xmlns,attr"`
	Title   string   `xml:"title"`
	Link    struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Updated   string  `xml:"updated"`
	Generator string  `xml:"generator"`
	ID        string  `xml:"id"`
	Author    author  `xml:"author"`
	Entry     []entry `xml:"entry"`
}

type author struct {
	Name string `xml:"name"`
}

type entry struct {
	Title string `xml:"title"`
	Link  struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	ID      string  `xml:"id"`
	Updated string  `xml:"updated"`
	Content content `xml:"content,omitempty"`
	Author  author  `xml:"author"`
}

type content struct {
	Content string `xml:",chardata"`
	Type    string `xml:"type,attr"`
}

// GenerateFeed generates an Atom feed with the top-ranked items
func GenerateFeed(ranking Ranking) (string, error) {
	if ranking.Len() == 0 {
		return "", errors.New("no items to process")
	}

	feed := createAtomFeed(ranking)
	output, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return "", err
	}
	return xml.Header + string(output), nil
}

func createAtomFeed(ranking Ranking) Feed {
	enrichArticlesWithContent(&ranking)

	for _, item := range ranking.Articles {
		entry := entry{
			Title: fmt.Sprintf("[%d] %s", item.Score, item.Title),
			Link: struct {
				Href string `xml:"href,attr"`
			}{
				Href: item.Link,
			},
			ID:      item.Link,
			Updated: item.Updated.Format(time.RFC3339),
			Content: content{
				Content: item.Content,
				Type:    "html",
			},
			Author: author{
				Name: item.Source,
			},
		}
		feed.Entry = append(feed.Entry, entry)
	}
	return feed
}

func enrichArticlesWithContent(ranking *Ranking) {
	for i, item := range ranking.Articles {
		if ranking.Articles[i].Content != "" {
			// content already available
			continue
		}

		ranking.Articles[i].Content = fetchContent(ranking.Articles[i])

		if false && item.Content != "" && item.Score == cfg.MinimalScore {
			item.Content = createSummary(item, cfg)
		}
	}
}

func fetchContent(feedItem FeedItem) string {
	art, err := readability.FromURL(feedItem.Link, 30*time.Second)
	if err != nil {
		log.Printf("error fetching article: %v", err)
		return ""
	}

	return art.Content
}

func createSummary(feedItem FeedItem, cfg Config) string {
	llm := NewLLMClient(cfg)
	summary, err := llm.SummarizeContent(feedItem.Content)
	if err != nil {
		log.Printf("error summarizing content: %v", err)
		return ""
	}
	return summary
}
