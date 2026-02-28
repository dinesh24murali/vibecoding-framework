package csv

import (
	"context"
	"errors"
	"strings"
	"testing"
)

type fakePlanJob struct {
	err  error
	rows []PlanImportRow
}

func (f *fakePlanJob) UpsertPlan(_ context.Context, row PlanImportRow) error {
	f.rows = append(f.rows, row)
	return f.err
}

func TestPlansImporterImportSuccessAndValidation(t *testing.T) {
	job := &fakePlanJob{}
	importer := NewPlansImporter(job)

	input := "provider_id,name,description,price,discount,is_active\n1,Family Pack,100 channels,299,20,true\n1,,bad,100,0,true\n"
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

func TestPlansImporterRejectsMissingRequiredHeader(t *testing.T) {
	job := &fakePlanJob{}
	importer := NewPlansImporter(job)

	_, err := importer.Import(context.Background(), strings.NewReader("provider_id,name,price\n1,Plan,10\n"))
	if !errors.Is(err, ErrInvalidCSVFile) {
		t.Fatalf("Import() error = %v, want %v", err, ErrInvalidCSVFile)
	}
}
