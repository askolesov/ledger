package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMonth_Validate(t *testing.T) {
	tests := []struct {
		name      string
		prevMonth *Month
		month     Month
		year      int
		monthNum  int
		wantErr   bool
		errMsg    string
	}{
		{
			name: "valid month with single account",
			prevMonth: &Month{
				OpeningBalance: 1000,
				ClosingBalance: 1000,
				Accounts: map[string]Account{
					"checking": {
						OpeningBalance: 1000,
						ClosingBalance: 1000,
					},
				},
			},
			month: Month{
				OpeningBalance: 1000,
				ClosingBalance: 1200,
				Accounts: map[string]Account{
					"checking": {
						OpeningBalance: 1000,
						ClosingBalance: 1200,
						Entries: []Entry{
							{Amount: 200, Note: "Salary", Date: "2024-01-15", Internal: false},
						},
					},
				},
			},
			year:     2024,
			monthNum: 1,
			wantErr:  false,
		},
		{
			name: "invalid month number - too low",
			month: Month{
				OpeningBalance: 1000,
				ClosingBalance: 1200,
				Accounts: map[string]Account{
					"checking": {
						OpeningBalance: 1000,
						ClosingBalance: 1200,
					},
				},
			},
			year:     2024,
			monthNum: 0,
			wantErr:  true,
			errMsg:   "month number must be between 1 and 12, got 0",
		},
		{
			name: "invalid month number - too high",
			month: Month{
				OpeningBalance: 1000,
				ClosingBalance: 1200,
				Accounts: map[string]Account{
					"checking": {
						OpeningBalance: 1000,
						ClosingBalance: 1200,
					},
				},
			},
			year:     2024,
			monthNum: 13,
			wantErr:  true,
			errMsg:   "month number must be between 1 and 12, got 13",
		},
		{
			name: "month with no accounts",
			month: Month{
				OpeningBalance: 1000,
				ClosingBalance: 1200,
				Accounts:       map[string]Account{},
			},
			year:     2024,
			monthNum: 1,
			wantErr:  true,
			errMsg:   "month 1 has no accounts",
		},
		{
			name: "opening balance mismatch with account sums",
			prevMonth: &Month{
				OpeningBalance: 1000,
				ClosingBalance: 1000,
				Accounts: map[string]Account{
					"checking": {
						OpeningBalance: 500,
						ClosingBalance: 500,
					},
				},
			},
			month: Month{
				OpeningBalance: 1000,
				ClosingBalance: 1200,
				Accounts: map[string]Account{
					"checking": {
						OpeningBalance: 500,
						ClosingBalance: 1200,
						Entries: []Entry{
							{Amount: 700, Note: "Salary", Date: "2024-01-15", Internal: false},
						},
					},
				},
			},
			year:     2024,
			monthNum: 1,
			wantErr:  true,
			errMsg:   "month opening balance 1000 does not equal sum of account opening balances 500",
		},
		{
			name: "closing balance mismatch with account sums",
			prevMonth: &Month{
				OpeningBalance: 1000,
				ClosingBalance: 1000,
				Accounts: map[string]Account{
					"checking": {
						OpeningBalance: 1000,
						ClosingBalance: 1000,
					},
				},
			},
			month: Month{
				OpeningBalance: 1000,
				ClosingBalance: 1200,
				Accounts: map[string]Account{
					"checking": {
						OpeningBalance: 1000,
						ClosingBalance: 1500,
						Entries: []Entry{
							{Amount: 500, Note: "Salary", Date: "2024-01-15", Internal: false},
						},
					},
				},
			},
			year:     2024,
			monthNum: 1,
			wantErr:  true,
			errMsg:   "month closing balance 1200 does not equal sum of account closing balances 1500",
		},
		{
			name: "internal entries sum not zero",
			prevMonth: &Month{
				OpeningBalance: 1000,
				ClosingBalance: 1000,
				Accounts: map[string]Account{
					"checking": {
						OpeningBalance: 1000,
						ClosingBalance: 1000,
					},
				},
			},
			month: Month{
				OpeningBalance: 1000,
				ClosingBalance: 1200,
				Accounts: map[string]Account{
					"checking": {
						OpeningBalance: 1000,
						ClosingBalance: 1200,
						Entries: []Entry{
							{Amount: 200, Note: "Transfer", Date: "2024-01-15", Internal: true},
						},
					},
				},
			},
			year:     2024,
			monthNum: 1,
			wantErr:  true,
			errMsg:   "sum of internal entries must be 0, got 200",
		},
		{
			name: "valid month with balanced internal entries",
			prevMonth: &Month{
				OpeningBalance: 1000,
				ClosingBalance: 1000,
				Accounts: map[string]Account{
					"checking": {
						OpeningBalance: 1000,
						ClosingBalance: 1000,
					},
				},
			},
			month: Month{
				OpeningBalance: 1000,
				ClosingBalance: 1000,
				Accounts: map[string]Account{
					"checking": {
						OpeningBalance: 1000,
						ClosingBalance: 900,
						Entries: []Entry{
							{Amount: -100, Note: "Transfer to savings", Date: "2024-01-15", Internal: true},
						},
					},
					"savings": {
						OpeningBalance: 0,
						ClosingBalance: 100,
						Entries: []Entry{
							{Amount: 100, Note: "Transfer from checking", Date: "2024-01-15", Internal: true},
						},
					},
				},
			},
			year:     2024,
			monthNum: 1,
			wantErr:  false,
		},
		{
			name: "previous month closing balance mismatch",
			prevMonth: &Month{
				OpeningBalance: 1000,
				ClosingBalance: 900,
				Accounts: map[string]Account{
					"checking": {
						OpeningBalance: 1000,
						ClosingBalance: 900,
					},
				},
			},
			month: Month{
				OpeningBalance: 1000,
				ClosingBalance: 1200,
				Accounts: map[string]Account{
					"checking": {
						OpeningBalance: 1000,
						ClosingBalance: 1200,
						Entries: []Entry{
							{Amount: 200, Note: "Salary"},
						},
					},
				},
			},
			year:     2024,
			monthNum: 2,
			wantErr:  true,
			errMsg:   "opening balance 1000 does not equal previous month's closing balance 900",
		},
		{
			name: "account omitted with non-zero previous closing balance",
			prevMonth: &Month{
				OpeningBalance: 800,
				ClosingBalance: 1000,
				Accounts: map[string]Account{
					"checking": {
						OpeningBalance: 800,
						ClosingBalance: 1000,
						Entries: []Entry{
							{Amount: 200, Note: "Salary"},
						},
					},
					"savings": {
						OpeningBalance: 0,
						ClosingBalance: 100,
						Entries: []Entry{
							{Amount: 100, Note: "Deposit"},
						},
					},
				},
			},
			month: Month{
				OpeningBalance: 1000,
				ClosingBalance: 1200,
				Accounts: map[string]Account{
					"checking": {
						OpeningBalance: 1000,
						ClosingBalance: 1200,
						Entries: []Entry{
							{Amount: 200, Note: "Salary"},
						},
					},
				},
			},
			year:     2024,
			monthNum: 2,
			wantErr:  true,
			errMsg:   "account savings: cannot be omitted because previous month closing balance is 100 (must be 0)",
		},
		{
			name: "account omitted with zero previous closing balance - valid",
			prevMonth: &Month{
				OpeningBalance: 200,
				ClosingBalance: 200,
				Accounts: map[string]Account{
					"checking": {
						OpeningBalance: 100,
						ClosingBalance: 200,
						Entries: []Entry{
							{Amount: 100, Note: "Salary", Date: "2024-01-15", Internal: false},
						},
					},
					"temp": {
						OpeningBalance: 100,
						ClosingBalance: 0,
						Entries: []Entry{
							{Amount: -100, Note: "Test spending"},
						},
					},
				},
			},
			month: Month{
				OpeningBalance: 200,
				ClosingBalance: 200,
				Accounts: map[string]Account{
					"checking": {
						OpeningBalance: 200,
						ClosingBalance: 200,
					},
				},
			},
			year:     2024,
			monthNum: 2,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.month.Validate(tt.year, tt.monthNum, tt.prevMonth)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMonth_Income(t *testing.T) {
	tests := []struct {
		name  string
		month Month
		want  int
	}{
		{
			name: "single income entry",
			month: Month{
				Accounts: map[string]Account{
					"checking": {
						Entries: []Entry{
							{Amount: 1000, Note: "Salary", Date: "2024-01-15", Internal: false},
						},
					},
				},
			},
			want: 1000,
		},
		{
			name: "multiple income entries",
			month: Month{
				Accounts: map[string]Account{
					"checking": {
						Entries: []Entry{
							{Amount: 1000, Note: "Salary", Date: "2024-01-15", Internal: false},
							{Amount: 500, Note: "Bonus", Date: "2024-01-20", Internal: false},
						},
					},
				},
			},
			want: 1500,
		},
		{
			name: "income and expenses mixed",
			month: Month{
				Accounts: map[string]Account{
					"checking": {
						Entries: []Entry{
							{Amount: 1000, Note: "Salary", Date: "2024-01-15", Internal: false},
							{Amount: -200, Note: "Rent", Date: "2024-01-01", Internal: false},
						},
					},
				},
			},
			want: 1000,
		},
		{
			name: "internal entries ignored",
			month: Month{
				Accounts: map[string]Account{
					"checking": {
						Entries: []Entry{
							{Amount: 1000, Note: "Transfer", Date: "2024-01-15", Internal: true},
						},
					},
				},
			},
			want: 0,
		},
		{
			name: "multiple accounts",
			month: Month{
				Accounts: map[string]Account{
					"checking": {
						Entries: []Entry{
							{Amount: 1000, Note: "Salary", Date: "2024-01-15", Internal: false},
						},
					},
					"savings": {
						Entries: []Entry{
							{Amount: 500, Note: "Interest", Date: "2024-01-31", Internal: false},
						},
					},
				},
			},
			want: 1500,
		},
		{
			name: "no income entries",
			month: Month{
				Accounts: map[string]Account{
					"checking": {
						Entries: []Entry{
							{Amount: -200, Note: "Rent", Date: "2024-01-01", Internal: false},
						},
					},
				},
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.month.Income()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMonth_Expenses(t *testing.T) {
	tests := []struct {
		name  string
		month Month
		want  int
	}{
		{
			name: "single expense entry",
			month: Month{
				Accounts: map[string]Account{
					"checking": {
						Entries: []Entry{
							{Amount: -200, Note: "Rent", Date: "2024-01-01", Internal: false},
						},
					},
				},
			},
			want: -200,
		},
		{
			name: "multiple expense entries",
			month: Month{
				Accounts: map[string]Account{
					"checking": {
						Entries: []Entry{
							{Amount: -200, Note: "Rent", Date: "2024-01-01", Internal: false},
							{Amount: -100, Note: "Utilities", Date: "2024-01-15", Internal: false},
						},
					},
				},
			},
			want: -300,
		},
		{
			name: "income and expenses mixed",
			month: Month{
				Accounts: map[string]Account{
					"checking": {
						Entries: []Entry{
							{Amount: 1000, Note: "Salary", Date: "2024-01-15", Internal: false},
							{Amount: -200, Note: "Rent", Date: "2024-01-01", Internal: false},
						},
					},
				},
			},
			want: -200,
		},
		{
			name: "internal entries ignored",
			month: Month{
				Accounts: map[string]Account{
					"checking": {
						Entries: []Entry{
							{Amount: -100, Note: "Transfer", Date: "2024-01-15", Internal: true},
						},
					},
				},
			},
			want: 0,
		},
		{
			name: "multiple accounts",
			month: Month{
				Accounts: map[string]Account{
					"checking": {
						Entries: []Entry{
							{Amount: -200, Note: "Rent", Date: "2024-01-01", Internal: false},
						},
					},
					"savings": {
						Entries: []Entry{
							{Amount: -50, Note: "Service fee", Date: "2024-01-31", Internal: false},
						},
					},
				},
			},
			want: -250,
		},
		{
			name: "no expense entries",
			month: Month{
				Accounts: map[string]Account{
					"checking": {
						Entries: []Entry{
							{Amount: 1000, Note: "Salary", Date: "2024-01-15", Internal: false},
						},
					},
				},
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.month.Expenses()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMonth_GetAccountNames(t *testing.T) {
	tests := []struct {
		name  string
		month Month
		want  []string
	}{
		{
			name: "single account",
			month: Month{
				Accounts: map[string]Account{
					"checking": {},
				},
			},
			want: []string{"checking"},
		},
		{
			name: "multiple accounts - already sorted",
			month: Month{
				Accounts: map[string]Account{
					"checking": {},
					"savings":  {},
					"credit":   {},
				},
			},
			want: []string{"checking", "credit", "savings"},
		},
		{
			name: "multiple accounts - unsorted",
			month: Month{
				Accounts: map[string]Account{
					"zebra":    {},
					"alpha":    {},
					"checking": {},
				},
			},
			want: []string{"alpha", "checking", "zebra"},
		},
		{
			name:  "no accounts",
			month: Month{Accounts: map[string]Account{}},
			want:  []string{},
		},
		{
			name: "case sensitive sorting",
			month: Month{
				Accounts: map[string]Account{
					"Checking": {},
					"checking": {},
					"CHECKING": {},
				},
			},
			want: []string{"CHECKING", "Checking", "checking"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.month.GetAccountNames()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMonth_Validate_AccountValidationErrors(t *testing.T) {
	// Test that account validation errors are properly wrapped
	month := Month{
		OpeningBalance: 1000,
		ClosingBalance: 1200,
		Accounts: map[string]Account{
			"checking": {
				OpeningBalance: 1000,
				ClosingBalance: 1200,
				Entries: []Entry{
					{Amount: 200, Note: "", Date: "2024-01-15", Internal: false}, // Invalid: empty note
				},
			},
		},
	}

	err := month.Validate(2024, 1, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "account checking:")
	assert.Contains(t, err.Error(), "note: cannot be blank")
}

func TestMonth_Validate_ComplexScenario(t *testing.T) {
	prevMonth := &Month{
		ClosingBalance: 2000,
		Accounts: map[string]Account{
			"checking": {
				ClosingBalance: 1500,
			},
			"savings": {
				ClosingBalance: 500,
			},
		},
	}

	// Test a complex scenario with multiple accounts and internal transfers
	month := Month{
		OpeningBalance: 2000,
		ClosingBalance: 3000,
		Accounts: map[string]Account{
			"checking": {
				OpeningBalance: 1500,
				ClosingBalance: 1800,
				Entries: []Entry{
					{Amount: 1000, Note: "Salary", Date: "2024-01-15", Internal: false},
					{Amount: -200, Note: "Rent", Date: "2024-01-01", Internal: false},
					{Amount: -500, Note: "Transfer to savings", Date: "2024-01-20", Internal: true},
				},
			},
			"savings": {
				OpeningBalance: 500,
				ClosingBalance: 1200,
				Entries: []Entry{
					{Amount: 500, Note: "Transfer from checking", Date: "2024-01-20", Internal: true},
					{Amount: 200, Note: "Interest", Date: "2024-01-31", Internal: false},
				},
			},
		},
	}

	err := month.Validate(2024, 1, prevMonth)
	assert.NoError(t, err)

	// Verify calculations
	assert.Equal(t, 1200, month.Income())   // 1000 + 200
	assert.Equal(t, -200, month.Expenses()) // -200 only
	assert.Equal(t, []string{"checking", "savings"}, month.GetAccountNames())
}
