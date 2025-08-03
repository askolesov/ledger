package v2

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Entry represents a single transaction record in an account
type Entry struct {
	Amount   int    `json:"amount" yaml:"amount" toml:"amount"`
	Internal bool   `json:"internal" yaml:"internal" toml:"internal"`
	Note     string `json:"note" yaml:"note" toml:"note"`
	Date     string `json:"date" yaml:"date" toml:"date"`
	Tag      string `json:"tag" yaml:"tag" toml:"tag"`
}

// Validate validates an entry according to OLF v2.0 rules
func (e Entry) Validate(year, month int) error {
	// E-2: Every Entry must include both amount and non-empty note fields
	// E-3: If date is present, it must strictly follow the ISO-8601 YYYY-MM-DD format
	err := validation.ValidateStruct(&e,
		validation.Field(&e.Amount, validation.Required),
		validation.Field(&e.Note, validation.Required, validation.Length(1, 0)),
		validation.Field(&e.Date, validation.Date("2006-01-02").Error("E-3: date format must be YYYY-MM-DD")),
	)
	if err != nil {
		return err
	}

	date, ok, err := e.ParseDate()
	if err != nil {
		return fmt.Errorf("E-3: invalid date format: %w", err)
	}

	// E-1: If an entry has a date, that date must lie within the year and month of its parent Month object
	if ok {
		if year != 0 && date.Year() != year {
			return fmt.Errorf("E-1: entry date year does not match expected year (expected: %d, got: %d)", year, date.Year())
		}

		if month != 0 && int(date.Month()) != month {
			return fmt.Errorf("E-1: entry date month does not match expected month (expected: %d, got: %d)", month, int(date.Month()))
		}
	}

	return nil
}

func (e Entry) ParseDate() (time.Time, bool, error) {
	if e.Date == "" {
		return time.Time{}, false, nil
	}

	date, err := time.Parse("2006-01-02", e.Date)
	if err != nil {
		return time.Time{}, false, err
	}

	return date, true, nil
}
