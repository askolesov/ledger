package v2

import (
	"fmt"

	"github.com/samber/lo"
)

// Account represents a financial account with entries
type Account struct {
	OpeningBalance int     `json:"opening_balance" yaml:"opening_balance" toml:"opening_balance"`
	ClosingBalance int     `json:"closing_balance" yaml:"closing_balance" toml:"closing_balance"`
	Entries        []Entry `json:"entries" yaml:"entries" toml:"entries"`
}

// Validate validates an account according to OLF v2.0 rules
func (a Account) Validate(year, month int, prevAccount *Account) error {
	// Validate all entries
	for i, entry := range a.Entries {
		if err := entry.Validate(year, month); err != nil {
			return fmt.Errorf("entry %d: %w", i, err)
		}
	}

	// A-1: For every account: opening_balance + Σ(entry.amount) = closing_balance
	calculatedBalance := a.OpeningBalance + a.EntriesSum()
	if calculatedBalance != a.ClosingBalance {
		return fmt.Errorf("A-1: account balance calculation incorrect (opening: %d + entries: %d = %d, expected closing: %d)",
			a.OpeningBalance, a.EntriesSum(), calculatedBalance, a.ClosingBalance)
	}

	if prevAccount != nil {
		// A-2: If an account exists in consecutive months, prev.closing_balance = next.opening_balance
		if a.OpeningBalance != prevAccount.ClosingBalance {
			return fmt.Errorf("A-2: account opening balance does not equal previous month closing balance (expected: %d, got: %d)",
				prevAccount.ClosingBalance, a.OpeningBalance)
		}
	} else {
		// A-3: A new account must start with opening_balance = 0
		if a.OpeningBalance != 0 {
			return fmt.Errorf("A-3: new account must start with opening balance 0 (got: %d)", a.OpeningBalance)
		}
	}

	return nil
}

// EntriesSum returns the sum of all entry amounts
func (a Account) EntriesSum() int {
	return lo.SumBy(a.Entries, func(entry Entry) int {
		return entry.Amount
	})
}

// InternalEntriesSum returns the sum of all internal entry amounts
func (a Account) InternalEntriesSum() int {
	return lo.SumBy(a.Entries, func(entry Entry) int {
		if entry.Internal {
			return entry.Amount
		}
		return 0
	})
}

// Income returns the sum of positive non-internal entries
func (a Account) Income() int {
	return lo.SumBy(a.Entries, func(entry Entry) int {
		if !entry.Internal && entry.Amount > 0 {
			return entry.Amount
		}
		return 0
	})
}

// Expenses returns the sum of negative non-internal entries
func (a Account) Expenses() int {
	return lo.SumBy(a.Entries, func(entry Entry) int {
		if !entry.Internal && entry.Amount < 0 {
			return entry.Amount
		}
		return 0
	})
}
