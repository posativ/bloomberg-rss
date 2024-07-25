package domain

import "encoding/xml"

type RSS struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title       string          `xml:"title"`
	Description string          `xml:"description"`
	Link        string          `xml:"link"`
	Image       RSSChannelImage `xml:"image"`
	Items       []RSSItem       `xml:"item"`
}

type RSSChannelImage struct {
	Url   string `xml:"url"`
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

type RSSItem struct {
	Title struct {
		XMLName xml.Name `xml:"title"`
		Text    string   `xml:",cdata"`
	}
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
