package journal

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func (s *PostgresStore) migrate(ctx context.Context) error {
	if err := s.ensureMigrationTable(ctx); err != nil {
		return err
	}

	migrations, err := fs.Glob(migrationFiles, "migrations/*.up.sql")
	if err != nil {
		return fmt.Errorf("list migrations: %w", err)
	}
	sort.Strings(migrations)

	for _, migrationPath := range migrations {
		version := migrationVersion(migrationPath)
		applied, err := s.migrationApplied(ctx, version)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		raw, err := migrationFiles.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", migrationPath, err)
		}
		if _, err := s.db.ExecContext(ctx, string(raw)); err != nil {
			return fmt.Errorf("apply migration %s: %w", migrationPath, err)
		}
		if err := s.markMigrationApplied(ctx, version); err != nil {
			return err
		}
	}

	return nil
}

func migrationVersion(path string) string {
	base := path[strings.LastIndex(path, "/")+1:]
	return strings.TrimSuffix(base, ".up.sql")
}

func (s *PostgresStore) ensureMigrationTable(ctx context.Context) error {
	const query = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version TEXT PRIMARY KEY,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
`
	if _, err := s.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("ensure migration table: %w", err)
	}

	return nil
}

func (s *PostgresStore) migrationApplied(ctx context.Context, version string) (bool, error) {
	const query = `SELECT 1 FROM schema_migrations WHERE version = $1`
	var found int
	scanErr := s.db.QueryRowContext(ctx, query, version).Scan(&found)
	if scanErr == nil {
		return true, nil
	}
	if scanErr == sqlErrNoRows {
		return false, nil
	}

	return false, fmt.Errorf("check migration %s: %w", version, scanErr)
}

func (s *PostgresStore) markMigrationApplied(ctx context.Context, version string) error {
	const query = `INSERT INTO schema_migrations (version) VALUES ($1)`
	if _, err := s.db.ExecContext(ctx, query, version); err != nil {
		return fmt.Errorf("record migration %s: %w", version, err)
	}

	return nil
}
