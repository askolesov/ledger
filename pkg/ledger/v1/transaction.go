package v1

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Transaction struct {
	Amount     int    `json:"amount" yaml:"amount" toml:"amount"`
	IsInternal bool   `json:"is_internal" yaml:"is_internal" toml:"is_internal"`
	Comment    string `json:"comment" yaml:"comment" toml:"comment"`

	Date     string `json:"date" yaml:"date" toml:"date"`
	Category string `json:"category" yaml:"category" toml:"category"`
}

func (t Transaction) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Amount, validation.Required),
		validation.Field(&t.Comment, validation.Required),
	)
}
