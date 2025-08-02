package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccount_Validate(t *testing.T) {
	tests := []struct {
		name    string
		account Account
		year    int
		month   int
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid account with balanced entries",
			account: Account{
				OpeningBalance: 1000,
				ClosingBalance: 1050,
				Entries: []Entry{
					{
						Amount: 100,
						Note:   "Salary",
					},
					{
						Amount: -50,
						Note:   "Groceries",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "account with no entries",
			account: Account{
				OpeningBalance: 1000,
				ClosingBalance: 1000,
				Entries:        []Entry{},
			},
			wantErr: false,
		},
		{
			name: "unbalanced account",
			account: Account{
				OpeningBalance: 1000,
				ClosingBalance: 1100, // Should be 1050
				Entries: []Entry{
					{
						Amount: 100,
						Note:   "Salary",
					},
					{
						Amount: -50,
						Note:   "Groceries",
					},
				},
			},
			wantErr: true,
			errMsg:  "does not equal closing balance",
		},
		{
			name: "entry with invalid date",
			account: Account{
				OpeningBalance: 1000,
				ClosingBalance: 1050,
				Entries: []Entry{
					{
						Amount:   50,
						Internal: false,
						Note:     "Invalid date entry",
						Date:     "invalid-date",
						Tag:      "Income",
					},
				},
			},
			year:    2025,
			month:   1,
			wantErr: true,
			errMsg:  "entry 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.account.Validate(tt.year, tt.month)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAccount_EntriesSum(t *testing.T) {
	account := Account{
		Entries: []Entry{
			{Amount: 100},
			{Amount: -50},
			{Amount: 25},
		},
	}

	sum := account.EntriesSum()
	assert.Equal(t, 75, sum)
}

func TestAccount_InternalEntriesSum(t *testing.T) {
	account := Account{
		Entries: []Entry{
			{Amount: 100, Internal: false},
			{Amount: -50, Internal: true},
			{Amount: 25, Internal: true},
			{Amount: 10, Internal: false},
		},
	}

	sum := account.InternalEntriesSum()
	assert.Equal(t, -25, sum) // -50 + 25
}

func TestAccount_Income(t *testing.T) {
	account := Account{
		Entries: []Entry{
			{Amount: 100, Internal: false}, // Income
			{Amount: -50, Internal: false}, // Expense (ignored)
			{Amount: 25, Internal: true},   // Internal (ignored)
			{Amount: 75, Internal: false},  // Income
		},
	}

	income := account.Income()
	assert.Equal(t, 175, income) // 100 + 75
}

func TestAccount_Expenses(t *testing.T) {
	account := Account{
		Entries: []Entry{
			{Amount: 100, Internal: false}, // Income (ignored)
			{Amount: -50, Internal: false}, // Expense
			{Amount: -25, Internal: true},  // Internal (ignored)
			{Amount: -75, Internal: false}, // Expense
		},
	}

	expenses := account.Expenses()
	assert.Equal(t, -125, expenses) // -50 + -75
}
