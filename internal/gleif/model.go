package gleif

// LEIRecord represents a single Legal Entity Identifier record from GLEIF golden copy CSV
type LEIRecord struct {
	LEI                 string
	LegalName           string
	LegalAddressLine1   string
	LegalAddressCity    string
	LegalAddressCountry string
	EntityStatus        string // ACTIVE, INACTIVE, etc.
}
