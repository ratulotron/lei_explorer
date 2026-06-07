package gleif

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
)

var ErrInvalidRecord = errors.New("invalid record: expected at least 6 columns")

func parseRecord(row []string) (LEIRecord, error) {
	if len(row) < 6 {
		return LEIRecord{}, ErrInvalidRecord
	}
	return LEIRecord{
		LEI:                 row[0],
		LegalName:           row[1],
		LegalAddressLine1:   row[2],
		LegalAddressCity:    row[3],
		LegalAddressCountry: row[4],
		EntityStatus:        row[5],
	}, nil
}

// ParseCSV reads all LEI records containing CSV data.
// Accepts an io.Reader - file, HTTP response body, etc.
// and returns a slice of LEIRecord structs.
func ParseCSV(r io.Reader) ([]LEIRecord, error) {
	reader := csv.NewReader(r)
	_, e := reader.Read() // Skip header row
	if e != nil {
		return nil, fmt.Errorf("reading csv header: %w", e)
	}

	var records []LEIRecord
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading csv: %w", err)
		}

		record, err := parseRecord(row)
		if err != nil {
			return nil, fmt.Errorf("parsing record: %w", err)
		}
		records = append(records, record)
	}
	return records, nil
}
