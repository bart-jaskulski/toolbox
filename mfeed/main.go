package main

import (
	"log"
	"os"
	"sync"
	"time"
)

var (
	cfg  = defaultConfig()
	feed = Feed{
		Xmlns:     "http://www.w3.org/2005/Atom",
		Title:     cfg.FeedTitle,
		Generator: "mfeed",
		Link: struct {
			Href string `xml:"href,attr"`
		}{
			Href: cfg.FeedHref,
		},
		Updated: time.Now().Format(time.RFC3339),
		ID:      cfg.FeedHref,
		Author: author{
			Name: cfg.FeedAuthor,
		},
	}
)

func main() {
	if len(cfg.Feeds) == 0 {
		log.Fatalf("no feed URLs found")
	}

	var wg sync.WaitGroup
	rankingsChan := make(chan Ranking, len(cfg.Feeds))

	for _, url := range cfg.Feeds {
		wg.Add(1)
		go func(feedURL string) {
			defer wg.Done()
			rank, err := processFeed(feedURL)
			if err != nil {
				log.Printf("error generating feed: %v", err)
				return
			}
			if rank.Len() > 0 {
				rankingsChan <- rank
			}
		}(url)
	}

	wg.Wait()
	close(rankingsChan)

	var metaRanking Ranking
	for rank := range rankingsChan {
		metaRanking.Articles = append(metaRanking.Articles, rank.Articles...)
	}

	if mErr := metaRanking.Rank(cfg); mErr != nil {
		log.Fatalf("error ranking meta feed: %v", mErr)
	}

	if metaRanking.Len() == 0 {
		log.Fatalf("no new items after ranking")
	}

	feedXML, genErr := GenerateFeed(metaRanking)
	if genErr != nil {
		log.Fatalf("error generating feed: %v", genErr)
	}

	os.Stdout.Write([]byte(feedXML))
}

func processFeed(feedURL string) (Ranking, error) {
	items, fetchErr := FetchFeed(feedURL, cfg.HistoricalDate)
	if fetchErr != nil {
		return Ranking{}, fetchErr
	}

	if len(items) == 0 {
		return Ranking{}, nil
	}

	feedRanking := Ranking{Articles: items}
	if rankErr := feedRanking.Rank(cfg); rankErr != nil {
		return Ranking{}, rankErr
	}

	return feedRanking, nil
}
