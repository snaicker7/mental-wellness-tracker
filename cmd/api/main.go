package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"mental-wellness-tracker/internal/journal"
)

func main() {
	if err := loadDotEnv(".env"); err != nil {
		log.Printf("dotenv: %v", err)
	}

	ctx := context.Background()
	addr := envOrDefault("ADDR", ":8080")
	dataFile := envOrDefault("JOURNAL_DATA_FILE", "data/journals.json")
	databaseURL := os.Getenv("DATABASE_URL")
	key, err := loadEncryptionKey()
	if err != nil {
		log.Fatal(err)
	}

	repo, cleanup, err := journal.NewRepository(ctx, databaseURL, dataFile, key)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if cleanup == nil {
			return
		}
		if err := cleanup(); err != nil {
			log.Printf("cleanup error: %v", err)
		}
	}()

	server := journal.NewServer(repo)
	log.Printf("mental wellness tracker API listening on %s", addr)
	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatal(err)
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func loadDotEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" || os.Getenv(key) != "" {
			continue
		}

		value = strings.Trim(value, `"`)
		value = strings.Trim(value, "'")
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func loadEncryptionKey() ([]byte, error) {
	value := os.Getenv("JOURNAL_ENCRYPTION_KEY")
	if value == "" {
		return nil, fmt.Errorf("JOURNAL_ENCRYPTION_KEY is required and must be base64-encoded 32 bytes")
	}

	key, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, fmt.Errorf("decode JOURNAL_ENCRYPTION_KEY: %w", err)
	}
	if len(key) != journal.AES256KeySize {
		return nil, fmt.Errorf("JOURNAL_ENCRYPTION_KEY must decode to %d bytes", journal.AES256KeySize)
	}

	return key, nil
}
