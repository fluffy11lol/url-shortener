package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"url-shortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	stmt, err := db.Prepare(`
CREATE TABLE IF NOT EXISTS urls (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	alias TEXT NOT NULL UNIQUE,
	url TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_alias ON urls(alias);`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	op := "storage.sqlite.SaveURL"
	stmt, err := s.db.Prepare("INSERT INTO urls (alias, url) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.Exec(alias, urlToSave)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) {
			return 0, storage.ErrUrlAlreadyExist
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return res.LastInsertId()
}

func (s *Storage) GetUrlByAlias(alias string) (string, error) {
	op := "storage.sqlite.GetUrlByAlias"
	stmt, err := s.db.Prepare("SELECT url FROM urls WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	var url string
	err = stmt.QueryRow(alias).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrUrlNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return url, nil
}

func (s *Storage) GetAliasByUrl(url string) (string, error) {
	op := "storage.sqlite.GetAliasByUrl"
	stmt, err := s.db.Prepare("SELECT alias FROM urls WHERE url = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	var alias string
	err = stmt.QueryRow(url).Scan(&alias)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrUrlNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return alias, nil
}

func (s *Storage) DeleteAlias(alias string) error {
	op := "storage.sqlite.DeleteAlias"
	stmt, err := s.db.Prepare("DELETE FROM urls WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrUrlNotFound)
	}
	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
