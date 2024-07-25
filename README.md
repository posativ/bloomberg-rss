Bloomberg RSS
=============

Provide a text feed of Bloomberg news articles by using their public (though truncated) RSS feeds to fetch the full
article content. Best experienced with a rotating IP addresses (or SOCKS5 proxy service).

## Installation

Go + SQLite3 is required I guess. Then `make linux` or `make darwin` to build the binary in `./bin/`.

## Usage

```bash
$ ./bloomberg-rss --help
Usage of ./bloomberg-rss:
  -addr string
        address to listen on (default "localhost:8080")
  -category string
        categories to fetch (comma separated) (default "all")
  -clear-queue
        clear the queue before starting
  -cookie string
        cookie to use for bloomberg.com (default "exp_pref=EUR; country_code=DE; seen_uk=1")
  -db string
        path to database file (default "bloomberg.db")
  -keep-article-max-age duration
        time in seconds to keep articles (default 24h0m0s
  -socks string
        SOCKS5 proxy to use
  -wait int
        time in seconds to wait between requests (default 5)

$ ./bloomberg-rss -addr :7071 -clear-queue -category markets,business -socks socks5://user:pass@socks5-proxy:12345 -wait 15
2024/07/25 20:21:07 Listening on http://:7071
2024/07/25 20:21:08 [worker] skipping excluded https://www.bloomberg.com/news/audio/2024-07-25/fed-easing-cycles-and-unsustainable-fiscal-policy-macro-matters
2024/07/25 20:21:08 [queue] add https://www.bloomberg.com/news/articles/2024-07-25/lineage-shares-pop-5-after-4-4-billion-ipo-warms-tepid-market
2024/07/25 20:21:09 [worker] next https://www.bloomberg.com/news/articles/2024-07-25/concentra-shares-slump-up-to-6-5-after-529-million-us-ipo
```

Available categories: business, markets, technology, politics, pursuits, opinion, economy and crypto.

RSS feed is available at `http://localhost:8080/` (or whatever address you specified).
