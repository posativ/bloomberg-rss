package storage

import (
	"database/sql"
	"errors"
	"log"
)

func (s *Storage) ClearQueue() error {
	_, err := s.db.Exec("DELETE FROM queue")
	return err
}

func (s *Storage) SubmitUrl(url string) error {
	var count int

	err := s.db.QueryRow("SELECT COUNT(*) FROM queue WHERE url = ?", url).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	_, err = s.db.Exec("INSERT INTO queue (url) VALUES (?) ON CONFLICT DO NOTHING", url)
	if err != nil {
		return err
	}

	log.Println("[queue] add", url)
	return nil
}

func (s *Storage) NextUrl() (string, error) {
	var url string
	err := s.db.QueryRow("SELECT url FROM queue ORDER BY created_at DESC LIMIT 1").Scan(&url)
	switch {
	case err == nil:
		return url, nil
	case errors.Is(err, sql.ErrNoRows):
		return "", nil
	default:
		return "", err
	}
}

func (s *Storage) RemoveUrl(url string) error {
	_, err := s.db.Exec("DELETE FROM queue WHERE url = ?", url)
	return err
}
