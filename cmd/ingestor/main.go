package main

import (
	"log"
	"os"

	"github.com/ratulotron/lei_explorer/internal/gleif"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) < 1 {
		log.Fatalf("Usage: ingestor <path_to_csv_file>")
	}
	csvFilePath := argsWithoutProg[0]

	// Open the CSV file
	file, err := os.Open(csvFilePath)
	if err != nil {
		log.Fatalf("Error opening file: %s", err.Error())
	}
	defer file.Close()

	// Parse the CSV file
	records, err := gleif.ParseCSV(file)
	if err != nil {
		log.Fatalf("Error parsing CSV: %s", err.Error())
	}

	log.Printf("Parsed %d records from CSV", len(records))
}
