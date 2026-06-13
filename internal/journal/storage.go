package journal

import (
	"context"
	"fmt"
	"os"
)

func NewRepository(ctx context.Context, databaseURL, dataFile string, key []byte) (Repository, func() error, error) {
	if databaseURL != "" {
		store, err := NewPostgresStore(ctx, databaseURL, key)
		if err != nil {
			return nil, nil, err
		}
		return store, store.Close, nil
	}

	if dataFile == "" {
		return nil, nil, fmt.Errorf("either DATABASE_URL or JOURNAL_DATA_FILE must be set")
	}

	store, err := NewFileStore(dataFile, key)
	if err != nil {
		return nil, nil, err
	}
	return store, func() error { return nil }, nil
}

func EnsureDataDir(path string) error {
	if path == "" {
		return nil
	}
	return os.MkdirAll(path, 0o755)
}
