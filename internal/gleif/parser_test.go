package gleif

import (
	_ "embed"
	"strings"
	"testing"
)

//go:embed sample.csv
var sampleCSV string

func TestParseRecord(t *testing.T) {
	tests := []struct {
		name     string
		row      []string
		expected LEIRecord
		wantErr  bool
	}{
		{
			name: "valid active record",
			row:  []string{"7H6GLXDRUGQFU57RNE97", "Goldman Sachs Group Inc.", "200 West Street", "New York", "US", "ACTIVE"},
			expected: LEIRecord{
				LEI:                 "7H6GLXDRUGQFU57RNE97",
				LegalName:           "Goldman Sachs Group Inc.",
				LegalAddressLine1:   "200 West Street",
				LegalAddressCity:    "New York",
				LegalAddressCountry: "US",
				EntityStatus:        "ACTIVE",
			},
		},
		{
			name:    "too few columns",
			row:     []string{"BADDATA"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record, err := parseRecord(tt.row)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && record != tt.expected {
				t.Fatalf("parseRecord() = %v, expected %v", record, tt.expected)
			}
		})
	}
}

func TestParseCSV(t *testing.T) {
	csvData := `LEI,LegalName,LegalAddressLine1,LegalAddressCity,LegalAddressCountry,EntityStatus
7H6GLXDRUGQFU57RNE97,Goldman Sachs Group Inc.,200 West Street,New York,US,ACTIVE
BADDATA`

	reader := strings.NewReader(csvData)

	records, err := ParseCSV(reader)
	if err == nil {
		t.Fatalf("Expected error for invalid CSV data, got nil")
	}
	if len(records) != 0 {
		t.Fatalf("Expected 0 records due to error, got %d", len(records))
	}
}

func TestOpenSampleData(t *testing.T) {
	result, err := ParseCSV(strings.NewReader(sampleCSV))
	if err != nil {
		t.Fatalf("Failed to parse sample data: %v", err)
	}

	t.Logf("Parsed rows: %d", len(result))

}
