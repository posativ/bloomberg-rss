package worker

import (
	"github.com/posativ/bloomberg-rss/src/config"
	"github.com/posativ/bloomberg-rss/src/storage"
	"log"
	"strings"
	"time"
)

type Worker struct {
	client *Client
	config *config.Config
	db     *storage.Storage
}

func NewWorker(config *config.Config, db *storage.Storage) *Worker {
	return &Worker{
		client: newClient(config.SocksProxy, config.Cookie),
		config: config,
		db:     db,
	}
}

func (w *Worker) Submit(url string, pubDate time.Time) {
	for _, include := range w.config.IncludePattern {
		if !strings.Contains(url, include) {
			log.Println("[worker] skipping excluded", url)
			return
		}
	}

	hasArticle, _ := w.db.HasArticle(url, pubDate)
	if !hasArticle {
		err := w.db.SubmitUrl(url)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (w *Worker) StartArticleCleaner(maxAge time.Duration) {
	go func() {
		for {
			err := w.db.TruncateArticles(time.Now().Add(-maxAge))
			if err != nil {
				log.Println(err)
			}
			time.Sleep(time.Duration(5) * time.Minute)
		}
	}()
}

func (w *Worker) StartFeeds(interval time.Duration) {
	go w.RefreshFeeds()

	ticker := time.NewTicker(interval)
	go func(fire <-chan time.Time) {
		for {
			select {
			case <-fire:
				log.Println("[worker] refresh feeds")
				w.RefreshFeeds()
			}
		}
	}(ticker.C)
}

func (w *Worker) StartQueue(delay time.Duration) {
	go func() {
		for {
			res := w.ProcessQueue()
			switch res {
			case NetworkSkipped:
				time.Sleep(time.Second)
			case NetworkUsed:
				time.Sleep(delay)
			}
		}
	}()
}
