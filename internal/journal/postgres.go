package journal

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStore struct {
	db  *sql.DB
	key []byte
}

func NewPostgresStore(ctx context.Context, databaseURL string, key []byte) (*PostgresStore, error) {
	if databaseURL == "" {
		return nil, errors.New("database url is required")
	}
	if len(key) != AES256KeySize {
		return nil, fmt.Errorf("encryption key must be %d bytes", AES256KeySize)
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(10)

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	store := &PostgresStore{db: db, key: key}
	if err := store.migrate(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

func (s *PostgresStore) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *PostgresStore) Create(ctx context.Context, entry Entry) (Entry, error) {
	plaintext, err := json.Marshal(entry)
	if err != nil {
		return Entry{}, fmt.Errorf("marshal entry: %w", err)
	}

	nonce, ciphertext, err := encryptJSON(s.key, plaintext)
	if err != nil {
		return Entry{}, err
	}

	const query = `
INSERT INTO journal_entries (
    id, user_id, nonce, ciphertext, created_at, updated_at
) VALUES ($1, $2, $3, $4, $5, $6)
`
	if _, err := s.db.ExecContext(ctx, query, entry.ID, entry.UserID, nonce, ciphertext, entry.CreatedAt, entry.UpdatedAt); err != nil {
		return Entry{}, fmt.Errorf("insert journal entry: %w", err)
	}

	return entry, nil
}

func (s *PostgresStore) ListByUser(ctx context.Context, userID string) ([]Entry, error) {
	const query = `
SELECT nonce, ciphertext
FROM journal_entries
WHERE user_id = $1
ORDER BY created_at ASC
`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("select journal entries: %w", err)
	}
	defer rows.Close()

	entries := make([]Entry, 0)
	for rows.Next() {
		var nonce string
		var ciphertext string
		if err := rows.Scan(&nonce, &ciphertext); err != nil {
			return nil, fmt.Errorf("scan journal entry: %w", err)
		}

		plaintext, err := decryptJSON(s.key, nonce, ciphertext)
		if err != nil {
			return nil, err
		}

		var entry Entry
		if err := json.Unmarshal(plaintext, &entry); err != nil {
			return nil, fmt.Errorf("unmarshal entry: %w", err)
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate journal entries: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].CreatedAt.Before(entries[j].CreatedAt)
	})

	return entries, nil
}
