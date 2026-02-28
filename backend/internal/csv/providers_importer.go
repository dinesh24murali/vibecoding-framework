package csv

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
)

var ErrInvalidCSVFile = errors.New("invalid csv file")

type ProvidersImporter struct {
	job ProviderImportJob
}

func NewProvidersImporter(job ProviderImportJob) *ProvidersImporter {
	return &ProvidersImporter{job: job}
}

func (i *ProvidersImporter) Import(ctx context.Context, reader io.Reader) (CsvUploadResult, error) {
	if reader == nil {
		return CsvUploadResult{}, ErrInvalidCSVFile
	}

	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true
	csvReader.FieldsPerRecord = -1

	header, err := csvReader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return CsvUploadResult{}, ErrInvalidCSVFile
		}
		return CsvUploadResult{}, fmt.Errorf("read csv header: %w", err)
	}

	nameIdx, imageIdx, err := mapProviderColumns(header)
	if err != nil {
		return CsvUploadResult{}, err
	}

	result := CsvUploadResult{Errors: make([]CsvRowError, 0)}
	rowNumber := 1

	for {
		rowNumber++
		record, readErr := csvReader.Read()
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				break
			}
			return CsvUploadResult{}, fmt.Errorf("read csv row %d: %w", rowNumber, readErr)
		}

		result.TotalRows++

		row, rowErr := parseProviderRow(record, nameIdx, imageIdx)
		if rowErr != nil {
			result.FailedRows++
			result.Errors = append(result.Errors, CsvRowError{Row: rowNumber, Field: rowErr.field, Issue: rowErr.issue})
			continue
		}

		if err := i.job.UpsertProvider(ctx, row); err != nil {
			result.FailedRows++
			result.Errors = append(result.Errors, CsvRowError{Row: rowNumber, Issue: "upsert_failed"})
			continue
		}

		result.SuccessRows++
	}

	return result, nil
}

type csvRowIssue struct {
	field string
	issue string
}

func mapProviderColumns(header []string) (nameIdx int, imageIdx int, err error) {
	nameIdx = -1
	imageIdx = -1

	for idx, col := range header {
		normalized := strings.ToLower(strings.TrimSpace(col))
		switch normalized {
		case "name":
			nameIdx = idx
		case "image_url":
			imageIdx = idx
		}
	}

	if nameIdx == -1 {
		return -1, -1, ErrInvalidCSVFile
	}

	return nameIdx, imageIdx, nil
}

func parseProviderRow(record []string, nameIdx int, imageIdx int) (ProviderImportRow, *csvRowIssue) {
	if nameIdx >= len(record) {
		return ProviderImportRow{}, &csvRowIssue{field: "name", issue: "required"}
	}

	name := strings.TrimSpace(record[nameIdx])
	if name == "" {
		return ProviderImportRow{}, &csvRowIssue{field: "name", issue: "required"}
	}
	if len(name) > 120 {
		return ProviderImportRow{}, &csvRowIssue{field: "name", issue: "max_length_exceeded"}
	}

	row := ProviderImportRow{Name: name}
	if imageIdx >= 0 && imageIdx < len(record) {
		imageRaw := strings.TrimSpace(record[imageIdx])
		if imageRaw != "" {
			parsed, err := url.ParseRequestURI(imageRaw)
			if err != nil || parsed.Scheme == "" || parsed.Host == "" {
				return ProviderImportRow{}, &csvRowIssue{field: "image_url", issue: "invalid_uri"}
			}
			row.ImageURL = &imageRaw
		}
	}

	return row, nil
}
