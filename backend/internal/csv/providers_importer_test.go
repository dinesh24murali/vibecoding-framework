package csv

import (
	"context"
	"errors"
	"strings"
	"testing"
)

type fakeProviderJob struct {
	err  error
	rows []ProviderImportRow
}

func (f *fakeProviderJob) UpsertProvider(_ context.Context, row ProviderImportRow) error {
	f.rows = append(f.rows, row)
	return f.err
}

func TestProvidersImporterImportSuccessAndPartialFailure(t *testing.T) {
	job := &fakeProviderJob{}
	importer := NewProvidersImporter(job)

	input := "name,image_url\nAirtel DTH,https://img.test/a.png\n,https://img.test/b.png\n"
	result, err := importer.Import(context.Background(), strings.NewReader(input))
	if err != nil {
		t.Fatalf("Import() error = %v", err)
	}

	if result.TotalRows != 2 || result.SuccessRows != 1 || result.FailedRows != 1 {
		t.Fatalf("unexpected counts: %+v", result)
	}
	if len(result.Errors) != 1 || result.Errors[0].Field != "name" {
		t.Fatalf("errors = %+v, want name validation error", result.Errors)
	}
	if len(job.rows) != 1 {
		t.Fatalf("upsert calls = %d, want 1", len(job.rows))
	}
}

func TestProvidersImporterRejectsMissingNameHeader(t *testing.T) {
	job := &fakeProviderJob{}
	importer := NewProvidersImporter(job)

	_, err := importer.Import(context.Background(), strings.NewReader("title,image_url\na,b\n"))
	if !errors.Is(err, ErrInvalidCSVFile) {
		t.Fatalf("Import() error = %v, want %v", err, ErrInvalidCSVFile)
	}
}
