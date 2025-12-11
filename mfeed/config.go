package main

import (
	"time"
)

const DefaultLookBackDays = -7 // Constant look-back days

type Config struct {
	LLMModel       string
	OpenAIEndpoint string
	MaxFeeds       int
	HistoricalDate time.Time
	MinimalScore   int
	FeedTitle      string
	FeedHref       string
	FeedAuthor     string
	Feeds          []string
}

func defaultConfig() Config {
	return Config{
		LLMModel:       "llama3-groq-70b-8192-tool-use-preview",
		OpenAIEndpoint: "https://api.groq.com/openai/v1",
		MaxFeeds:       120,
		HistoricalDate: time.Now().AddDate(0, 0, DefaultLookBackDays),
		MinimalScore:   3,
		FeedTitle:      "Most interesting tech news",
		FeedHref:       "https://bart-jaskulski.github.io/mfeed/",
		FeedAuthor:     "Bart Jaskulski",
		Feeds: []string{
			"https://www.garfieldtech.com/blog/feed",
			"https://doeken.org/blog/feed.atom",
			"https://martinfowler.com/feed.atom",
			"https://whynothugo.nl/posts.xml",
			"https://tomasvotruba.com/rss",
			"https://www.dantleech.com/index.xml",
			"https://developer.woocommerce.com/feed/",
			"https://developer.wordpress.org/news/feed/",
			"https://make.wordpress.org/core/feed/",
			"https://make.wordpress.org/plugins/feed/",
			"https://make.wordpress.org/performance/feed/",
			"https://feeds.feedburner.com/symfony/blog",
			"https://phpstan.org/rss.xml",
			"https://blog.jetbrains.com/phpstorm/feed/",
			"http://blog.cleancoder.com/atom.xml",
			"https://jvns.ca/atom.xml",
			"https://stackoverflow.blog/feed",
			"https://hnrss.org/best?description=0",
			"https://www.freecodecamp.org/news/rss/",
			"https://tonsky.me/atom.xml",
			"https://blog.jim-nielsen.com/feed.xml",
			"https://www.smashingmagazine.com/feed/",
			"https://css-tricks.com/feed/",
			"https://buttondown.email/hillelwayne/rss",
			"https://lucasfcosta.com/feed.xml",
		},
	}
}
