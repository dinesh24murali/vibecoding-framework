package migrate

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestListMigrationFiles(t *testing.T) {
	dir := t.TempDir()

	paths := []string{
		filepath.Join(dir, "0002_indexes.sql"),
		filepath.Join(dir, "0001_init.sql"),
		filepath.Join(dir, "README.md"),
	}

	for _, path := range paths {
		if err := os.WriteFile(path, []byte("-- test"), 0o600); err != nil {
			t.Fatalf("write file %s: %v", path, err)
		}
	}

	files, err := listMigrationFiles(dir)
	if err != nil {
		t.Fatalf("listMigrationFiles() error = %v", err)
	}

	want := []string{
		filepath.Join(dir, "0001_init.sql"),
		filepath.Join(dir, "0002_indexes.sql"),
	}

	if !reflect.DeepEqual(files, want) {
		t.Fatalf("listMigrationFiles() = %v, want %v", files, want)
	}
}

func TestVersionFromFile(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		got, err := versionFromFile("0003_payment_sms_retry.sql")
		if err != nil {
			t.Fatalf("versionFromFile() error = %v", err)
		}

		if got != 3 {
			t.Fatalf("versionFromFile() = %d, want 3", got)
		}
	})

	t.Run("invalid format", func(t *testing.T) {
		_, err := versionFromFile("migration.sql")
		if err == nil {
			t.Fatalf("versionFromFile() error = nil, want non-nil")
		}
	})
}
