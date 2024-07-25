package config

import "time"

type Config struct {
	Addr              string
	SocksProxy        string
	Feeds             []Feed
	Cookie            string
	WaitTime          time.Duration
	KeepArticleMaxAge time.Duration
	IncludePattern    []string
}

type Feed struct {
	Url      string
	Category string
}

func NewConfig() *Config {
	return &Config{
		Addr: ":8080",
		Feeds: []Feed{
			{Url: "https://feeds.bloomberg.com/business/news.rss", Category: "business"},
			{Url: "https://feeds.bloomberg.com/markets/news.rss", Category: "markets"},
			{Url: "https://feeds.bloomberg.com/technology/news.rss", Category: "technology"},
			{Url: "https://feeds.bloomberg.com/politics/news.rss", Category: "politics"},
			{Url: "https://feeds.bloomberg.com/pursuits/news.rss", Category: "pursuits"},
			{Url: "https://feeds.bloomberg.com/economics/news.rss", Category: "economy"},
			{Url: "https://feeds.bloomberg.com/wealth/news.rss", Category: "wealth"},
			{Url: "https://feeds.bloomberg.com/crypto/news.rss", Category: "crypto"},
		},
		IncludePattern: []string{
			"/news/articles/",
		},
	}
}
