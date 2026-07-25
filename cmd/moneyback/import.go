package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"money-backword/internal/ledger"
	"money-backword/internal/storage"
)

func cmdImport(args []string, store *storage.JSONStore) {
	fs := flag.NewFlagSet("import", flag.ExitOnError)
	filePath := fs.String("file", "", "path to import file (CSV or OFX)")
	fileFormat := fs.String("format", "csv", "file format: csv or ofx")
	accountID := fs.String("account", "", "account ID to import into")
	defaultCategory := fs.String("category", "imported", "default category for transactions")

	fs.Parse(args)

	if *filePath == "" || *accountID == "" {
		fmt.Fprintf(os.Stderr, "import: missing required flags: -file and -account\n")
		return
	}

	importer := NewImporter(store, *accountID, *defaultCategory)

	var count int
	var err error

	switch *fileFormat {
	case "csv":
		count, err = importer.ImportCSV(*filePath)
	case "ofx":
		count, err = importer.ImportOFX(*filePath)
	default:
		fmt.Fprintf(os.Stderr, "import: unsupported format: %s (csv, ofx supported)\n", *fileFormat)
		return
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "import error: %v\n", err)
		return
	}

	fmt.Printf("successfully imported %d transactions into account %s\n", count, *accountID)
}

type Importer struct {
	store           *storage.JSONStore
	accountID       string
	defaultCategory string
}

func NewImporter(store *storage.JSONStore, accountID, defaultCategory string) *Importer {
	return &Importer{
		store:           store,
		accountID:       accountID,
		defaultCategory: defaultCategory,
	}
}

func (imp *Importer) ImportCSV(filePath string) (int, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) < 2 {
		return 0, fmt.Errorf("CSV file too short or invalid")
	}

	// TODO: parse CSV header and infer column positions
	// TODO: handle various CSV formats (date format, amount sign conventions)
	// TODO: add duplicate detection

	count := 0
	for i, line := range lines[1:] { // Skip header
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) < 3 {
			fmt.Fprintf(os.Stderr, "warning: line %d skipped (invalid format)\n", i+2)
			continue
		}

		// Simple parsing: assume format is date, description, amount
		date := strings.TrimSpace(parts[0])
		description := strings.TrimSpace(parts[1])
		amountStr := strings.TrimSpace(parts[2])

		var amount float64
		_, err := fmt.Sscanf(amountStr, "%f", &amount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: line %d invalid amount: %s\n", i+2, amountStr)
			continue
		}

		txn := ledger.NewTransaction(imp.accountID, imp.defaultCategory, amount, fmt.Sprintf("[%s] %s", date, description))

		if err := imp.store.AddTransaction(txn); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to add transaction from line %d: %v\n", i+2, err)
			continue
		}

		count++
	}

	return count, nil
}

func (imp *Importer) ImportOFX(filePath string) (int, error) {
	// TODO: implement OFX parsing
	// OFX is a more complex banking format with XML-like structure
	// This stub provides the structure for future implementation

	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to read file: %w", err)
	}

	// TODO: parse OFX format
	// Look for STMTRS (statement response) sections
	// Extract STMTTRN (statement transaction) records
	// Parse date (YYMMDD format), amount, and memo fields

	if len(data) == 0 {
		return 0, fmt.Errorf("OFX file is empty")
	}

	// Stub: return no transactions
	return 0, fmt.Errorf("OFX import not yet implemented")
}
