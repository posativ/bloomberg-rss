package worker

import (
	"encoding/xml"
	"io"
	"log"
	"time"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Title       string `xml:"title"`
		Description string `xml:"description"`
		Link        string `xml:"link"`
		Image       struct {
			Url   string `xml:"url"`
			Title string `xml:"title"`
			Link  string `xml:"link"`
		} `xml:"image"`
		Items []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Description struct {
		XMLName xml.Name `xml:"description"`
		Text    string   `xml:",cdata"`
	} `xml:"description"`
	Link string `xml:"link"`
	Guid struct {
		IsPermaLink string `xml:"isPermaLink,attr"`
		Value       string `xml:",chardata"`
	} `xml:"guid"`
	PubDate string `xml:"pubDate"`
}

func (w *Worker) RefreshFeeds() {
	for _, feed := range w.config.Feeds {
		err := w.processFeed(feed.Url)
		if err != nil {
			log.Println(err)
		}
	}
}

func (w *Worker) processFeed(url string) error {
	resp, err := w.client.get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var rss RSS
	err = xml.Unmarshal(data, &rss)
	if err != nil {
		return err
	}

	for _, item := range rss.Channel.Items {
		pubDate, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			return err
		}
		w.Submit(item.Link, pubDate)
	}

	return nil
}
