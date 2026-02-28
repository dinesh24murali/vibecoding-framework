package csv

import "context"

type ProviderImportRow struct {
	Name     string
	ImageURL *string
}

type ProviderImportJob interface {
	UpsertProvider(ctx context.Context, row ProviderImportRow) error
}

type CsvRowError struct {
	Row   int    `json:"row"`
	Field string `json:"field,omitempty"`
	Issue string `json:"issue"`
}

type CsvUploadResult struct {
	TotalRows   int           `json:"total_rows"`
	SuccessRows int           `json:"success_rows"`
	FailedRows  int           `json:"failed_rows"`
	Errors      []CsvRowError `json:"errors"`
}
