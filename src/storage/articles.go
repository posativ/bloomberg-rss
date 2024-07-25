package storage

import "time"

func (s *Storage) RemoveFromQueueAndWriteArticle(url string, title string, content string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO articles (url, title, content) VALUES (?, ?, ?)
		ON CONFLICT (url) DO UPDATE SET title = ?, content = ?, updated_at = CURRENT_TIMESTAMP`,
		url, title, content, title, content,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM queue WHERE url = ?", url)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *Storage) HasArticle(url string, pubDate time.Time) (bool, error) {
	var count int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM articles WHERE url = ? AND updated_at >= ?`, url, pubDate,
	).Scan(&count)
	return count > 0, err
}

func (s *Storage) ReadArticle(url string) (string, string, error) {
	var title, content string
	err := s.db.QueryRow(`
		SELECT title, content FROM articles WHERE url = ? `, url,
	).Scan(&title, &content)
	return title, content, err
}

func (s *Storage) TruncateArticles(before time.Time) error {
	_, err := s.db.Exec("DELETE FROM articles WHERE updated_at < ?", before)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("VACUUM ")
	return err
}
