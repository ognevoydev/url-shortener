package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"url-shortener/internal/lib/security"
	"url-shortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fn = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS url (
			id INTEGER PRIMARY KEY,
			url TEXT NOT NULL,
			alias TEXT NOT NULL UNIQUE
		);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
		`,
		`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			token TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
		`,
	}

	for _, query := range queries {
		query, err := db.Prepare(query)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		_, err = query.Exec()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const fn = "storage.sqlite.SaveURL"

	query, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	res, err := query.Exec(urlToSave, alias)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: %w", fn, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const fn = "storage.sqlite.GetURL"

	query, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	var resURL string
	err = query.QueryRow(alias).Scan(&resURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", fn, storage.ErrURLNotFound)
		}

		return "", fmt.Errorf("%s: %w", fn, err)
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const fn = "storage.sqlite.DeleteURL"

	query, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = query.Exec(alias)

	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) CreateUser(username string, password string) (int64, error) {
	const fn = "storage.sqlite.CreateUser"

	// TODO handle user already exists error

	query, err := s.db.Prepare("INSERT INTO users(username, password) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	hashedPassword, err := security.HashPassword(password)

	res, err := query.Exec(username, hashedPassword)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: %w", fn, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}

func (s *Storage) AuthenticateUser(username string, password string) (int64, error) {
	const fn = "storage.sqlite.AuthenticateUser"

	query, err := s.db.Prepare("SELECT id, password FROM users WHERE username = ?")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	var userId int64
	var hashedPassword string

	err = query.QueryRow(username).Scan(&userId, &hashedPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("%s: %w", fn, storage.ErrUserNotFound)
		}

		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	auth := security.VerifyPassword(password, hashedPassword)
	if !auth {
		return 0, nil
	}

	return userId, nil
}

func (s *Storage) CreateSession(userId int64, token string) (int64, error) {
	const fn = "storage.sqlite.CreateSession"

	query, err := s.db.Prepare("INSERT INTO sessions (user_id, token) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	res, err := query.Exec(userId, token)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}
