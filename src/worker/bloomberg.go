package worker

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
)

type NetworkUsage int

const (
	NetworkUsed NetworkUsage = iota
	NetworkSkipped
)

func (w *Worker) ProcessQueue() NetworkUsage {
	url, err := w.db.NextUrl()
	if err != nil {
		log.Println("[worker] processFeed queue url:", url, "error:", err)
		return NetworkSkipped
	}

	if url == "" {
		return NetworkSkipped
	}

	log.Println("[worker] next", url)

	title, data, err := w.processArticle(url)
	if err != nil {
		log.Println("[worker] error processing url:", url, "error:", err)
		return NetworkUsed
	}

	err = w.db.RemoveFromQueueAndWriteArticle(url, title, data)
	if err != nil {
		log.Println("[worker] error writing article:", url, "error:", err)
		return NetworkUsed
	}

	return NetworkUsed
}

func (w *Worker) processArticle(url string) (string, string, error) {
	url = strings.Replace(url, "https://www.bloomberg.com", "https://www.bloomberg.com.", 1)

	resp, err := w.client.get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	str := string(content)
	titleExpr := regexp.MustCompile(`<title>([^<]*)</title>`)
	title := titleExpr.FindStringSubmatch(str)
	if len(title) < 2 {
		return "", "", fmt.Errorf("unable to find title")
	}

	nextDataExpr := regexp.MustCompile(`<script id="__NEXT_DATA__" type="application/json">([^<]*)</script>`)
	nextData := nextDataExpr.FindStringSubmatch(str)
	if len(nextData) < 2 {
		if strings.Contains(str, "Bloomberg - Are you a robot?") {
			return "", "", fmt.Errorf("bloomberg thinks we are a robot")
		}

		if strings.Contains(url, "/news/videos/") {
			return title[1], "{}", nil
		}

		return title[1], "{}", fmt.Errorf("unable to find __NEXT_DATA__")
	}

	return title[1], nextData[1], nil
}
