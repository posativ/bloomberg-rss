package storage

import (
	"github.com/posativ/bloomberg-rss/src/domain"
)

func (s *Storage) GetRssItems() ([]domain.RSSItem, error) {
	rows, err := s.db.Query(`
		SELECT articles.title, articles.content, articles.url, articles.created_at FROM articles
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.RSSItem
	for rows.Next() {
		var item domain.RSSItem
		if err := rows.Scan(&item.Title.Text, &item.Description.Text, &item.Link, &item.PubDate); err != nil {
			return nil, err
		}

		item.Guid.Value = item.Link
		item.Guid.IsPermaLink = "true"

		items = append(items, item)
	}

	return items, nil

}
