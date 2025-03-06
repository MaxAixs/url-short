package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
)

type URLRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{db: db}
}

func (r *URLRepository) SaveURL(url string, shortUrl string) (int, error) {
	q := `INSERT INTO urls (original_url, short_url) VALUES ($1, $2) RETURNING id`

	var id int

	err := r.db.QueryRow(q, url, shortUrl).Scan(&id)
	if err != nil {
		if uniqueViolation(err) {
			return 0, ErrUrlExists
		}

		return 0, fmt.Errorf("failed to insert URL %w", err)
	}

	return id, nil
}

func uniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code == "23505" {
			return true
		}
	}

	return false
}

func (r *URLRepository) DeleteURL(shortUrl string) error {
	q := `DELETE FROM urls WHERE short_url = $1`

	_, err := r.db.Exec(q, shortUrl)
	if err != nil {
		return fmt.Errorf("failed to delete URL %w", err)
	}

	return nil
}

func (r *URLRepository) GetURL(shortUrl string) (string, error) {
	q := `SELECT original_url FROM urls WHERE short_url = $1`

	var url string

	err := r.db.QueryRow(q, shortUrl).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		}

		return "", fmt.Errorf("failed to get URL %w", err)
	}

	return url, nil
}
