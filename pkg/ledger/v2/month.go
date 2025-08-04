package v2

import (
	"fmt"
	"sort"

	"github.com/samber/lo"
)

// Month represents a month with accounts
type Month struct {
	OpeningBalance int                `json:"opening_balance" yaml:"opening_balance" toml:"opening_balance"`
	ClosingBalance int                `json:"closing_balance" yaml:"closing_balance" toml:"closing_balance"`
	Accounts       map[string]Account `json:"accounts" yaml:"accounts" toml:"accounts"`
}

// Validate validates a month according to OLF v2.0 rules
func (m Month) Validate(year, monthNum int, prevMonth *Month) error {
	// M-0: Month key (monthNum) must be between 1 and 12 (inclusive)
	if monthNum < 1 || monthNum > 12 {
		return fmt.Errorf("M-0: month number must be between 1 and 12 (got: %d)", monthNum)
	}

	// M-5: A Month must contain at least one Account entry
	if len(m.Accounts) == 0 {
		return fmt.Errorf("M-5: month must contain at least one account")
	}

	// Validate accounts
	for accountName, account := range m.Accounts {
		var prevAccount *Account
		if prevMonth != nil {
			if val, ok := prevMonth.Accounts[accountName]; ok {
				prevAccount = &val
			}
		}

		if err := account.Validate(year, monthNum, prevAccount, prevMonth != nil); err != nil {
			return fmt.Errorf("account %s: %w", accountName, err)
		}
	}

	// M-2: Month opening_balance equals sum of all account opening_balance values
	accountsOpeningSum := lo.SumBy(lo.Values(m.Accounts), func(account Account) int {
		return account.OpeningBalance
	})
	if m.OpeningBalance != accountsOpeningSum {
		return fmt.Errorf("M-2: month opening balance does not equal sum of account opening balances (expected: %d, got: %d)",
			accountsOpeningSum, m.OpeningBalance)
	}

	// M-3: Month closing_balance equals sum of all account closing_balance values
	accountsClosingSum := lo.SumBy(lo.Values(m.Accounts), func(account Account) int {
		return account.ClosingBalance
	})
	if m.ClosingBalance != accountsClosingSum {
		return fmt.Errorf("M-3: month closing balance does not equal sum of account closing balances (expected: %d, got: %d)",
			accountsClosingSum, m.ClosingBalance)
	}

	// M-4: Within each month, Σ(entry.amount where internal = true) must equal 0 (double-entry constraint)
	internalSum := lo.SumBy(lo.Values(m.Accounts), func(account Account) int {
		return account.InternalEntriesSum()
	})
	if internalSum != 0 {
		return fmt.Errorf("M-4: sum of internal entries must equal 0 (got: %d)", internalSum)
	}

	if prevMonth != nil {
		// M-1: Consecutive months must chain totals
		if prevMonth.ClosingBalance != m.OpeningBalance {
			return fmt.Errorf("M-1: month opening balance does not equal previous month closing balance (expected: %d, got: %d)",
				prevMonth.ClosingBalance, m.OpeningBalance)
		}

		// A-4: An account may be omitted in later months only if its last closing_balance = 0
		for accountName, prevAccount := range prevMonth.Accounts {
			if _, exists := m.Accounts[accountName]; !exists {
				if prevAccount.ClosingBalance != 0 {
					return fmt.Errorf("A-4: account '%s' cannot be omitted with non-zero closing balance (got: %d)",
						accountName, prevAccount.ClosingBalance)
				}
			}
		}
	}

	return nil
}

// Income returns the sum of income from all accounts
func (m Month) Income() int {
	return lo.SumBy(lo.Values(m.Accounts), func(account Account) int {
		return account.Income()
	})
}

// Expenses returns the sum of expenses from all accounts
func (m Month) Expenses() int {
	return lo.SumBy(lo.Values(m.Accounts), func(account Account) int {
		return account.Expenses()
	})
}

// GetAccountNames returns sorted list of account names
func (m Month) GetAccountNames() []string {
	names := lo.Keys(m.Accounts)
	sort.Strings(names)
	return names
}
