package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Runner struct {
	db            *sql.DB
	migrationsDir string
}

func NewRunner(db *sql.DB, migrationsDir string) *Runner {
	return &Runner{db: db, migrationsDir: migrationsDir}
}

func (r *Runner) Up(ctx context.Context) error {
	if err := r.ensureVersionTable(ctx); err != nil {
		return err
	}

	dirtyVersion, err := r.dirtyVersion(ctx)
	if err != nil {
		return err
	}

	if dirtyVersion > 0 {
		return fmt.Errorf("cannot continue: migration version %d is marked dirty", dirtyVersion)
	}

	files, err := listMigrationFiles(r.migrationsDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		version, err := versionFromFile(file)
		if err != nil {
			return err
		}

		applied, err := r.isApplied(ctx, version)
		if err != nil {
			return err
		}

		if applied {
			continue
		}

		if err := r.markDirty(ctx, version); err != nil {
			return fmt.Errorf("mark migration %d dirty: %w", version, err)
		}

		if err := r.executeMigration(ctx, file); err != nil {
			return fmt.Errorf("apply migration %d (%s): %w", version, filepath.Base(file), err)
		}

		if err := r.markClean(ctx, version); err != nil {
			return fmt.Errorf("mark migration %d clean: %w", version, err)
		}
	}

	return nil
}

func (r *Runner) ensureVersionTable(ctx context.Context) error {
	const query = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version BIGINT PRIMARY KEY,
    dirty BOOLEAN NOT NULL DEFAULT FALSE,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("ensure schema_migrations: %w", err)
	}

	return nil
}

func (r *Runner) dirtyVersion(ctx context.Context) (int64, error) {
	const query = `SELECT version FROM schema_migrations WHERE dirty = TRUE ORDER BY version DESC LIMIT 1`

	var version int64
	err := r.db.QueryRowContext(ctx, query).Scan(&version)
	if err == sql.ErrNoRows {
		return 0, nil
	}

	if err != nil {
		return 0, fmt.Errorf("query dirty version: %w", err)
	}

	return version, nil
}

func (r *Runner) isApplied(ctx context.Context, version int64) (bool, error) {
	const query = `SELECT COUNT(1) FROM schema_migrations WHERE version = $1 AND dirty = FALSE`

	var count int
	err := r.db.QueryRowContext(ctx, query, version).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("query migration %d state: %w", version, err)
	}

	return count > 0, nil
}

func (r *Runner) markDirty(ctx context.Context, version int64) error {
	const query = `
INSERT INTO schema_migrations (version, dirty, applied_at)
VALUES ($1, TRUE, NOW())
ON CONFLICT (version)
DO UPDATE SET dirty = EXCLUDED.dirty, applied_at = EXCLUDED.applied_at`

	_, err := r.db.ExecContext(ctx, query, version)
	return err
}

func (r *Runner) markClean(ctx context.Context, version int64) error {
	const query = `UPDATE schema_migrations SET dirty = FALSE, applied_at = NOW() WHERE version = $1`

	_, err := r.db.ExecContext(ctx, query, version)
	return err
}

func (r *Runner) executeMigration(ctx context.Context, path string) error {
	payload, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read migration file: %w", err)
	}

	_, err = r.db.ExecContext(ctx, string(payload))
	if err != nil {
		return fmt.Errorf("execute migration SQL: %w", err)
	}

	return nil
}

func listMigrationFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read migrations dir %q: %w", dir, err)
	}

	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		files = append(files, filepath.Join(dir, entry.Name()))
	}

	sort.Strings(files)
	return files, nil
}

func versionFromFile(path string) (int64, error) {
	name := filepath.Base(path)
	parts := strings.SplitN(name, "_", 2)
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid migration filename %q: expected <version>_<name>.sql", name)
	}

	version, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid migration version in %q: %w", name, err)
	}

	return version, nil
}
