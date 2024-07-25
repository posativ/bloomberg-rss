package server

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/posativ/bloomberg-rss/src/domain"
	"html"
	"log"
	"net/http"
	"strings"
)

func unwrapLink(href string, webUrl string, inner string) string {
	if href == "" {
		return inner
	}

	if strings.HasPrefix(href, "bbg://people/") {
		return inner
	}

	switch {
	case strings.HasPrefix(href, "http:"):
	case strings.HasPrefix(href, "https:"):
	case strings.HasPrefix(href, "mailto:"):
		// noop
	case strings.HasPrefix(href, "bbg://securities/"):
		asset := strings.Split(href, "/")[1] // TODO: -1 is the correct index
		href = fmt.Sprintf("https://www.bloomberg.com/quote/%s", asset)
	case strings.HasPrefix(href, "bbg://news/"):
		if webUrl == "" {
			// requires Bloomberg Terminal
			return inner
		}

		href = strings.Replace(webUrl, "https://www.bloomberg.com", "https://www.bloomberg.com.", 1)
	case strings.HasPrefix(href, "bbg://screens/"):
		// e.g. bbg://screens/FED
		return inner
	case strings.HasPrefix(href, "bbg://msg/"):
		// e.g. bbg://msg/foo@bar.com
		return "mailto:" + href[9:]
	default:
		log.Println("Unable to parse href:", href)
		return inner
	}

	return fmt.Sprintf("<a href=\"%s\">%s</a>", href, inner)
}

func transform(content domain.Content) (string, error) {

	inner := ""
	if content.Content != nil {
		for _, c := range content.Content {
			res, err := transform(c)
			if err != nil {
				return "", err
			}
			inner += res
		}
	}

	switch content.Type {
	case "":
		// unknown
		return inner, nil
	case "document":
		// root element
		return inner, nil
	case "heading":
		level := content.Data.Level
		if level != 0 {
			return fmt.Sprintf("<h%d>%s</h%d>", level, inner, level), nil
		} else {
			return "", fmt.Errorf("heading without level")
		}
	case "paragraph":
		return fmt.Sprintf("<p>%s</p>", inner), nil
	case "quote", "aside":
		return fmt.Sprintf("<blockquote>%s</blockquote>", inner), nil
	case "text":
		return html.UnescapeString(content.Value) + inner, nil
	case "br":
		return "<br />" + inner, nil
	case "link":
		return unwrapLink(content.Data.Href, content.Data.WebUrl, inner), nil
	case "list":
		switch content.SubType {
		case "ordered":
			return fmt.Sprintf("<ol>%s</ol>", inner), nil
		case "unordered":
			return fmt.Sprintf("<ul>%s</ul>", inner), nil
		default:
			return "", fmt.Errorf("unknown list subtype: %s", content.SubType)
		}
	case "listItem":
		return fmt.Sprintf("<li>%s</li>", inner), nil
	case "div":
		return fmt.Sprintf("<div>%s</div>", inner), nil
	case "embed":
		return content.IFrameData.Html, nil
	case "media":
		switch content.SubType {
		case "photo":
			return fmt.Sprintf("<img src=\"%s\" alt=\"%s\" /><br/><em>%s</em>", content.Data.Photo.Src, content.Data.Photo.Alt, content.Data.Photo.Caption), nil
		case "chart":
			return fmt.Sprintf("<img src=\"%s\" alt=\"%s\" />", content.Data.Chart.Fallback, content.Data.Chart.Caption), nil
		case "audio":
			return fmt.Sprintf("<audio controls src=\"%s\">%s</audio>", content.Data.Attachment.Url, content.Data.Attachment.Title), nil
		case "video":
			// difficult to support, we only have the video playlist as x-mpegURL or m3u8
			return "", nil
		default:
			return "", fmt.Errorf("unknown media subtype: %s", content.SubType)
		}
	case "entity":
		switch content.SubType {
		case "security":
			return content.Meta.Security, nil
		case "story":
			return unwrapLink(content.Data.Href, content.Data.WebUrl, inner), nil
		case "person":
			return unwrapLink(content.Data.Href, content.Data.WebUrl, inner), nil
		default:
			return "", fmt.Errorf("unknown entity subtype: %s", content.SubType)

		}
	case "inline-newsletter", "inline-recirc", "ad", "columns", "row", "cell", "tabularData", "callout", "byTheNumbers", "footnoteRef":
		return "", nil
	default:
		return "", fmt.Errorf("unsupported content type '%s'", content.Type)
	}
}

func (s *Server) RssHandler(w http.ResponseWriter, r *http.Request) {
	items, err := s.db.GetRssItems()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var transformed []domain.RSSItem
	for _, item := range items {
		var description domain.Html
		err := json.Unmarshal([]byte(item.Description.Text), &description)
		if err != nil {
			log.Printf("[server] unable to parse HTML url=%s error=%s\n", item.Link, err)
			continue
		}

		content, err := transform(description.Props.PageProps.Story.Body)
		if err != nil {
			log.Printf("[server] unable to transform JSON url=%s error=%s\n", item.Link, err)
			continue
		}

		if content == "" {
			continue
		}

		item.Description.Text = content

		transformed = append(transformed, item)
	}

	rss := domain.RSS{
		XMLName: xml.Name{Local: "rss"},
		Version: "2.0",
		Channel: domain.RSSChannel{
			Title:       "Bloomberg",
			Description: "",
			Image: domain.RSSChannelImage{
				Url:   "https://www.bloomberg.com/feeds/static/images/bloomberg_logo_blue.png",
				Title: "Bloomberg",
				Link:  "https://www.bloomberg.com",
			},
			Link:  "https://www.bloomberg.com",
			Items: transformed,
		},
	}
	result, err := xml.Marshal(rss)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/rss+xml")
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8" ?>`))
	w.Write(result)

}
