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
	err := validation.ValidateStruct(&e,
		validation.Field(&e.Amount, validation.Required),
		validation.Field(&e.Note, validation.Required, validation.Length(1, 0)),
		validation.Field(&e.Date, validation.Date("2006-01-02").Error("date format must be YYYY-MM-DD")),
	)
	if err != nil {
		return err
	}

	date, ok, err := e.ParseDate()
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	if ok {
		if year != 0 && date.Year() != year {
			return fmt.Errorf("entry date year %d does not match expected year %d", date.Year(), year)
		}

		if month != 0 && int(date.Month()) != month {
			return fmt.Errorf("entry date month %d does not match expected month %d", int(date.Month()), month)
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
