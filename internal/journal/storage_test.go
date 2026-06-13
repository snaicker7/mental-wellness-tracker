package journal

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestNewRepositoryRequiresStorageConfiguration(t *testing.T) {
	repo, cleanup, err := NewRepository(context.Background(), "", "", mustKey(t))
	if err == nil {
		t.Fatal("expected missing storage configuration error")
	}
	if repo != nil {
		t.Fatalf("expected nil repo, got %#v", repo)
	}
	if cleanup != nil {
		t.Fatal("expected nil cleanup")
	}
}

func TestNewFileStoreRejectsInvalidConfiguration(t *testing.T) {
	if _, err := NewFileStore(filepath.Join(t.TempDir(), "journals.json"), []byte("short-key")); err == nil {
		t.Fatal("expected invalid encryption key error")
	}
	if _, err := NewFileStore("", mustKey(t)); err == nil {
		t.Fatal("expected empty storage path error")
	}
}

func TestFileStoreReturnsDecodeErrorForCorruptState(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "journals.json")
	if err := os.WriteFile(storePath, []byte("{not-json"), 0o600); err != nil {
		t.Fatalf("write corrupt store: %v", err)
	}

	store, err := NewFileStore(storePath, mustKey(t))
	if err != nil {
		t.Fatalf("NewFileStore() error = %v", err)
	}

	if _, err := store.ListByUser(context.Background(), "student-1"); err == nil {
		t.Fatal("expected corrupt storage decode error")
	}
}

func TestEnsureDataDir(t *testing.T) {
	if err := EnsureDataDir(""); err != nil {
		t.Fatalf("EnsureDataDir empty path error = %v", err)
	}

	path := filepath.Join(t.TempDir(), "nested", "data")
	if err := EnsureDataDir(path); err != nil {
		t.Fatalf("EnsureDataDir() error = %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("expected directory to exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("expected %s to be a directory", path)
	}
}
