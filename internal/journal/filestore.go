package journal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type FileStore struct {
	path string
	key  []byte
	mu   sync.Mutex
}

type persistedState struct {
	Records []persistedRecord `json:"records"`
}

type persistedRecord struct {
	Nonce      string `json:"nonce"`
	Ciphertext string `json:"ciphertext"`
}

func NewFileStore(path string, key []byte) (*FileStore, error) {
	if len(key) != AES256KeySize {
		return nil, fmt.Errorf("encryption key must be %d bytes", AES256KeySize)
	}
	if path == "" {
		return nil, errors.New("storage path is required")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create storage directory: %w", err)
	}

	return &FileStore{path: path, key: key}, nil
}

func (s *FileStore) Create(_ context.Context, entry Entry) (Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, err := s.readState()
	if err != nil {
		return Entry{}, err
	}

	plaintext, err := json.Marshal(entry)
	if err != nil {
		return Entry{}, fmt.Errorf("marshal entry: %w", err)
	}

	nonce, ciphertext, err := encryptJSON(s.key, plaintext)
	if err != nil {
		return Entry{}, err
	}

	state.Records = append(state.Records, persistedRecord{Nonce: nonce, Ciphertext: ciphertext})
	if err := s.writeState(state); err != nil {
		return Entry{}, err
	}

	return entry, nil
}

func (s *FileStore) ListByUser(_ context.Context, userID string) ([]Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, err := s.readState()
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(state.Records))
	for _, record := range state.Records {
		plaintext, err := decryptJSON(s.key, record.Nonce, record.Ciphertext)
		if err != nil {
			return nil, err
		}

		var entry Entry
		if err := json.Unmarshal(plaintext, &entry); err != nil {
			return nil, fmt.Errorf("unmarshal entry: %w", err)
		}
		if entry.UserID == userID {
			entries = append(entries, entry)
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].CreatedAt.Before(entries[j].CreatedAt)
	})

	return entries, nil
}

func (s *FileStore) readState() (persistedState, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return persistedState{}, nil
		}
		return persistedState{}, fmt.Errorf("read storage: %w", err)
	}
	if len(data) == 0 {
		return persistedState{}, nil
	}

	var state persistedState
	if err := json.Unmarshal(data, &state); err != nil {
		return persistedState{}, fmt.Errorf("decode storage: %w", err)
	}

	return state, nil
}

func (s *FileStore) writeState(state persistedState) error {
	raw, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("encode storage: %w", err)
	}
	if err := os.WriteFile(s.path, raw, 0o600); err != nil {
		return fmt.Errorf("write storage: %w", err)
	}

	return nil
}
