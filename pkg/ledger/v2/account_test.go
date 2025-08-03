package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccount_Validate(t *testing.T) {
	tests := []struct {
		name        string
		prevAccount *Account
		account     Account
		year        int
		month       int
		wantErr     bool
		errMsg      string
	}{
		{
			name: "valid account with no entries",
			prevAccount: &Account{
				ClosingBalance: 1000,
			},
			account: Account{
				OpeningBalance: 1000,
				ClosingBalance: 1000,
				Entries:        []Entry{},
			},
			year:    2025,
			month:   1,
			wantErr: false,
		},
		{
			name: "valid account with entries",
			prevAccount: &Account{
				ClosingBalance: 1000,
			},
			account: Account{
				OpeningBalance: 1000,
				ClosingBalance: 1200,
				Entries: []Entry{
					{Amount: 200, Note: "Income", Date: "2025-01-15"},
				},
			},
			year:    2025,
			month:   1,
			wantErr: false,
		},
		{
			name: "valid account with multiple entries",
			prevAccount: &Account{
				ClosingBalance: 1000,
			},
			account: Account{
				OpeningBalance: 1000,
				ClosingBalance: 1100,
				Entries: []Entry{
					{Amount: 200, Note: "Income", Date: "2025-01-15"},
					{Amount: -100, Note: "Expense", Date: "2025-01-16"},
				},
			},
			year:    2025,
			month:   1,
			wantErr: false,
		},
		{
			name: "invalid balance calculation",
			prevAccount: &Account{
				ClosingBalance: 1000,
			},
			account: Account{
				OpeningBalance: 1000,
				ClosingBalance: 1200,
				Entries: []Entry{
					{Amount: 100, Note: "Income", Date: "2025-01-15"},
				},
			},
			year:    2025,
			month:   1,
			wantErr: true,
			errMsg:  "A-1: account balance calculation incorrect (opening: 1000 + entries: 100 = 1100, expected closing: 1200)",
		},
		{
			name: "invalid entry in account",
			prevAccount: &Account{
				ClosingBalance: 1000,
			},
			account: Account{
				OpeningBalance: 1000,
				ClosingBalance: 1000,
				Entries: []Entry{
					{Amount: 0, Note: "", Date: "2025-01-15"}, // Invalid entry
				},
			},
			year:    2025,
			month:   1,
			wantErr: true,
			errMsg:  "entry 0: amount: cannot be blank; note: cannot be blank",
		},
		{
			name: "entry with year mismatch",
			prevAccount: &Account{
				ClosingBalance: 1000,
			},
			account: Account{
				OpeningBalance: 1000,
				ClosingBalance: 1000,
				Entries: []Entry{
					{Amount: 100, Note: "Test", Date: "2024-01-15"}, // Wrong year
				},
			},
			year:    2025,
			month:   1,
			wantErr: true,
			errMsg:  "entry 0: E-1: entry date year does not match expected year (expected: 2025, got: 2024)",
		},
		{
			name: "previous account closing balance mismatch",
			prevAccount: &Account{
				ClosingBalance: 1000,
			},
			account: Account{
				OpeningBalance: 1100,
				ClosingBalance: 1100,
			},
			year:    2025,
			month:   1,
			wantErr: true,
			errMsg:  "A-2: account opening balance does not equal previous month closing balance (expected: 1000, got: 1100)",
		},
		{
			name: "new account with non-zero opening balance",
			account: Account{
				OpeningBalance: 1000,
				ClosingBalance: 1000,
			},
			year:    2025,
			month:   1,
			wantErr: true,
			errMsg:  "A-3: new account must start with opening balance 0 (got: 1000)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.account.Validate(tt.year, tt.month, tt.prevAccount)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAccount_EntriesSum(t *testing.T) {
	tests := []struct {
		name    string
		account Account
		want    int
	}{
		{
			name: "empty entries",
			account: Account{
				Entries: []Entry{},
			},
			want: 0,
		},
		{
			name: "single positive entry",
			account: Account{
				Entries: []Entry{
					{Amount: 100, Note: "Income"},
				},
			},
			want: 100,
		},
		{
			name: "single negative entry",
			account: Account{
				Entries: []Entry{
					{Amount: -50, Note: "Expense"},
				},
			},
			want: -50,
		},
		{
			name: "multiple entries",
			account: Account{
				Entries: []Entry{
					{Amount: 100, Note: "Income"},
					{Amount: -30, Note: "Expense"},
					{Amount: 25, Note: "Refund"},
				},
			},
			want: 95,
		},
		{
			name: "mixed internal and external entries",
			account: Account{
				Entries: []Entry{
					{Amount: 100, Internal: false, Note: "External"},
					{Amount: 50, Internal: true, Note: "Internal"},
					{Amount: -20, Internal: false, Note: "External"},
				},
			},
			want: 130,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.account.EntriesSum()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAccount_InternalEntriesSum(t *testing.T) {
	tests := []struct {
		name    string
		account Account
		want    int
	}{
		{
			name: "empty entries",
			account: Account{
				Entries: []Entry{},
			},
			want: 0,
		},
		{
			name: "no internal entries",
			account: Account{
				Entries: []Entry{
					{Amount: 100, Internal: false, Note: "External"},
					{Amount: -50, Internal: false, Note: "External"},
				},
			},
			want: 0,
		},
		{
			name: "only internal entries",
			account: Account{
				Entries: []Entry{
					{Amount: 100, Internal: true, Note: "Internal"},
					{Amount: -30, Internal: true, Note: "Internal"},
				},
			},
			want: 70,
		},
		{
			name: "mixed internal and external entries",
			account: Account{
				Entries: []Entry{
					{Amount: 100, Internal: false, Note: "External"},
					{Amount: 50, Internal: true, Note: "Internal"},
					{Amount: -20, Internal: false, Note: "External"},
					{Amount: -10, Internal: true, Note: "Internal"},
				},
			},
			want: 40,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.account.InternalEntriesSum()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAccount_Income(t *testing.T) {
	tests := []struct {
		name    string
		account Account
		want    int
	}{
		{
			name: "empty entries",
			account: Account{
				Entries: []Entry{},
			},
			want: 0,
		},
		{
			name: "no positive external entries",
			account: Account{
				Entries: []Entry{
					{Amount: -100, Internal: false, Note: "Expense"},
					{Amount: 50, Internal: true, Note: "Internal"},
				},
			},
			want: 0,
		},
		{
			name: "only positive external entries",
			account: Account{
				Entries: []Entry{
					{Amount: 100, Internal: false, Note: "Salary"},
					{Amount: 25, Internal: false, Note: "Bonus"},
				},
			},
			want: 125,
		},
		{
			name: "mixed entries with positive external",
			account: Account{
				Entries: []Entry{
					{Amount: 100, Internal: false, Note: "Salary"},
					{Amount: 50, Internal: true, Note: "Internal"},
					{Amount: -30, Internal: false, Note: "Expense"},
					{Amount: 25, Internal: false, Note: "Refund"},
				},
			},
			want: 125,
		},
		{
			name: "zero amount external entry",
			account: Account{
				Entries: []Entry{
					{Amount: 0, Internal: false, Note: "Zero"},
				},
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.account.Income()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAccount_Expenses(t *testing.T) {
	tests := []struct {
		name    string
		account Account
		want    int
	}{
		{
			name: "empty entries",
			account: Account{
				Entries: []Entry{},
			},
			want: 0,
		},
		{
			name: "no negative external entries",
			account: Account{
				Entries: []Entry{
					{Amount: 100, Internal: false, Note: "Income"},
					{Amount: -50, Internal: true, Note: "Internal"},
				},
			},
			want: 0,
		},
		{
			name: "only negative external entries",
			account: Account{
				Entries: []Entry{
					{Amount: -100, Internal: false, Note: "Rent"},
					{Amount: -25, Internal: false, Note: "Food"},
				},
			},
			want: -125,
		},
		{
			name: "mixed entries with negative external",
			account: Account{
				Entries: []Entry{
					{Amount: 100, Internal: false, Note: "Salary"},
					{Amount: -50, Internal: true, Note: "Internal"},
					{Amount: -30, Internal: false, Note: "Rent"},
					{Amount: -25, Internal: false, Note: "Food"},
				},
			},
			want: -55,
		},
		{
			name: "zero amount external entry",
			account: Account{
				Entries: []Entry{
					{Amount: 0, Internal: false, Note: "Zero"},
				},
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.account.Expenses()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAccount_Integration(t *testing.T) {
	// Test a realistic account scenario
	prevAccount := &Account{
		ClosingBalance: 1000,
	}
	account := Account{
		OpeningBalance: 1000,
		ClosingBalance: 1250,
		Entries: []Entry{
			{Amount: 500, Internal: false, Note: "Salary", Date: "2025-01-15"},
			{Amount: -200, Internal: false, Note: "Rent", Date: "2025-01-16"},
			{Amount: -50, Internal: false, Note: "Food", Date: "2025-01-17"},
			{Amount: 100, Internal: true, Note: "Transfer"},
			{Amount: -100, Internal: true, Note: "Transfer", Date: "2025-01-18"},
		},
	}

	t.Run("validate account", func(t *testing.T) {
		err := account.Validate(2025, 1, prevAccount)
		require.NoError(t, err)
	})

	t.Run("check entries sum", func(t *testing.T) {
		sum := account.EntriesSum()
		assert.Equal(t, 250, sum) // 500 - 200 - 50 + 100 - 100 = 250
	})

	t.Run("check internal entries sum", func(t *testing.T) {
		sum := account.InternalEntriesSum()
		assert.Equal(t, 0, sum) // 100 - 100 = 0
	})

	t.Run("check income", func(t *testing.T) {
		income := account.Income()
		assert.Equal(t, 500, income) // Only the 500 salary
	})

	t.Run("check expenses", func(t *testing.T) {
		expenses := account.Expenses()
		assert.Equal(t, -250, expenses) // -200 rent - 50 food = -250
	})

	t.Run("verify balance calculation", func(t *testing.T) {
		calculatedBalance := account.OpeningBalance + account.EntriesSum()
		assert.Equal(t, account.ClosingBalance, calculatedBalance)
	})
}
