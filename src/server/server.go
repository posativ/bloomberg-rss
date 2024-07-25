package server

import (
	"flag"
	"github.com/posativ/bloomberg-rss/src/config"
	"github.com/posativ/bloomberg-rss/src/storage"
	"github.com/posativ/bloomberg-rss/src/worker"
	"log"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	config *config.Config
	db     *storage.Storage
	worker *worker.Worker
}

func NewServer() (*Server, error) {
	addrFlag := flag.String("addr", "localhost:8080", "address to listen on")
	cookieFlag := flag.String("cookie", "exp_pref=EUR; country_code=DE; seen_uk=1", "cookie to use for bloomberg.com")
	dbFlag := flag.String("db", "bloomberg.db", "path to database file")
	waitTimeFlag := flag.Int("wait", 5, "time in seconds to wait between requests")
	keepArticleMaxAgeFlag := flag.Duration("keep-article-max-age", time.Duration(24*60*60)*time.Second, "time in seconds to keep articles")
	clearQueueFlag := flag.Bool("clear-queue", false, "clear the queue before starting")
	socksFlag := flag.String("socks", "", "SOCKS5 proxy to use")
	categoryFlag := flag.String("category", "all", "categories to fetch (comma separated)")

	flag.Parse()

	cfg := config.NewConfig()
	cfg.Addr = *addrFlag
	cfg.SocksProxy = *socksFlag
	cfg.Cookie = *cookieFlag
	cfg.WaitTime = time.Duration(*waitTimeFlag) * time.Second
	cfg.KeepArticleMaxAge = *keepArticleMaxAgeFlag

	if *categoryFlag != "all" {
		categories := strings.Split(*categoryFlag, ",")
		selectedFeeds := make([]config.Feed, 0)
		for _, feed := range cfg.Feeds {
			for _, category := range categories {
				if feed.Category == category {
					selectedFeeds = append(selectedFeeds, feed)
				}
			}
		}
		cfg.Feeds = selectedFeeds
	}

	db, err := storage.New(*dbFlag)
	if err != nil {
		return nil, err
	}

	if *clearQueueFlag {
		err = db.ClearQueue()
		if err != nil {
			log.Println(err)
		}
	}

	return &Server{
		config: cfg,
		db:     db,
		worker: worker.NewWorker(cfg, db),
	}, nil
}

func (s *Server) Start() {
	s.worker.StartArticleCleaner(s.config.KeepArticleMaxAge)
	s.worker.StartFeeds(time.Minute * time.Duration(10))
	s.worker.StartQueue(s.config.WaitTime)

	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: s.handler(),
	}
	log.Printf("Listening on http://%s\n", s.config.Addr)
	log.Fatal(server.ListenAndServe())
}

func (s *Server) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.RssHandler)
	return mux
}
