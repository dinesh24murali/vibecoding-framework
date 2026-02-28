package csv

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type PlansImporter struct {
	job PlanImportJob
}

func NewPlansImporter(job PlanImportJob) *PlansImporter {
	return &PlansImporter{job: job}
}

func (i *PlansImporter) Import(ctx context.Context, reader io.Reader) (CsvUploadResult, error) {
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

	indexByName, err := mapPlanColumns(header)
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

		row, rowErr := parsePlanRow(record, indexByName)
		if rowErr != nil {
			result.FailedRows++
			result.Errors = append(result.Errors, CsvRowError{Row: rowNumber, Field: rowErr.field, Issue: rowErr.issue})
			continue
		}

		if err := i.job.UpsertPlan(ctx, row); err != nil {
			result.FailedRows++
			field, issue := classifyPlanUpsertError(err)
			result.Errors = append(result.Errors, CsvRowError{Row: rowNumber, Field: field, Issue: issue})
			continue
		}

		result.SuccessRows++
	}

	return result, nil
}

func mapPlanColumns(header []string) (map[string]int, error) {
	required := []string{"provider_id", "name", "description", "price", "discount", "is_active"}
	indexes := make(map[string]int, len(required))

	for idx, col := range header {
		indexes[strings.ToLower(strings.TrimSpace(col))] = idx
	}

	for _, key := range required {
		if _, ok := indexes[key]; !ok {
			return nil, ErrInvalidCSVFile
		}
	}

	return indexes, nil
}

func parsePlanRow(record []string, indexByName map[string]int) (PlanImportRow, *csvRowIssue) {
	providerID, err := parseInt64Field(record, indexByName["provider_id"])
	if err != nil || providerID <= 0 {
		return PlanImportRow{}, &csvRowIssue{field: "provider_id", issue: "invalid_integer"}
	}

	name, err := parseStringField(record, indexByName["name"])
	if err != nil || len(name) > 120 {
		if err != nil {
			return PlanImportRow{}, &csvRowIssue{field: "name", issue: "required"}
		}
		return PlanImportRow{}, &csvRowIssue{field: "name", issue: "max_length_exceeded"}
	}

	description, err := parseStringField(record, indexByName["description"])
	if err != nil || len(description) > 1000 {
		if err != nil {
			return PlanImportRow{}, &csvRowIssue{field: "description", issue: "required"}
		}
		return PlanImportRow{}, &csvRowIssue{field: "description", issue: "max_length_exceeded"}
	}

	price, err := parseFloatField(record, indexByName["price"])
	if err != nil || price < 0.01 {
		return PlanImportRow{}, &csvRowIssue{field: "price", issue: "invalid_amount"}
	}

	discount, err := parseFloatField(record, indexByName["discount"])
	if err != nil || discount < 0 || discount > price {
		return PlanImportRow{}, &csvRowIssue{field: "discount", issue: "invalid_amount"}
	}

	isActive, err := parseBoolField(record, indexByName["is_active"])
	if err != nil {
		return PlanImportRow{}, &csvRowIssue{field: "is_active", issue: "invalid_boolean"}
	}

	return PlanImportRow{
		ProviderID:  providerID,
		Name:        name,
		Description: description,
		Price:       price,
		Discount:    discount,
		IsActive:    isActive,
	}, nil
}

func parseStringField(record []string, idx int) (string, error) {
	if idx >= len(record) {
		return "", errors.New("missing")
	}
	v := strings.TrimSpace(record[idx])
	if v == "" {
		return "", errors.New("empty")
	}
	return v, nil
}

func parseInt64Field(record []string, idx int) (int64, error) {
	v, err := parseStringField(record, idx)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(v, 10, 64)
}

func parseFloatField(record []string, idx int) (float64, error) {
	v, err := parseStringField(record, idx)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(v, 64)
}

func parseBoolField(record []string, idx int) (bool, error) {
	v, err := parseStringField(record, idx)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(v)
}

func classifyPlanUpsertError(_ error) (string, string) {
	return "", "upsert_failed"
}
