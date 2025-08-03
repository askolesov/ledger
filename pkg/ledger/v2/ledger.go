package v2

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

// Ledger represents the root structure of the Open Ledger Format v2.0
type Ledger struct {
	Years map[int]Year `json:"years" yaml:"years" toml:"years"`
}

// Validate validates the entire ledger according to OLF v2.0 rules
func (l Ledger) Validate() error {
	// Get sorted year numbers for consistent validation order
	yearNums := lo.Keys(l.Years)
	sort.Ints(yearNums)

	// Validate all years
	var prevYear *Year
	for _, yearNum := range yearNums {
		year := l.Years[yearNum]

		if err := year.Validate(yearNum, prevYear); err != nil {
			return fmt.Errorf("year %d: %w", yearNum, err)
		}

		prevYear = &year
	}

	return nil
}

// Income returns the total income across all years
func (l Ledger) Income() int {
	return lo.SumBy(lo.Values(l.Years), func(year Year) int {
		return year.Income()
	})
}

// Expenses returns the total expenses across all years
func (l Ledger) Expenses() int {
	return lo.SumBy(lo.Values(l.Years), func(year Year) int {
		return year.Expenses()
	})
}

// ReadLedger reads and parses a ledger file in YAML, JSON, or TOML format
func ReadLedger(path string) (Ledger, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return Ledger{}, fmt.Errorf("failed to read file: %w", err)
	}

	ledger := Ledger{}
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		err = json.Unmarshal(bytes, &ledger)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(bytes, &ledger)
	default:
		return Ledger{}, fmt.Errorf("unsupported file format: %s", filepath.Ext(path))
	}

	if err != nil {
		return Ledger{}, fmt.Errorf("failed to parse file: %w", err)
	}

	return ledger, nil
}

// WriteLedger writes a ledger to a file in the specified format
func WriteLedger(ledger Ledger, path string) error {
	var data []byte
	var err error

	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		data, err = json.MarshalIndent(ledger, "", "  ")
	case ".yaml", ".yml":
		data, err = yaml.Marshal(ledger)
	default:
		return fmt.Errorf("unsupported file format: %s", filepath.Ext(path))
	}

	if err != nil {
		return fmt.Errorf("failed to marshal ledger: %w", err)
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
