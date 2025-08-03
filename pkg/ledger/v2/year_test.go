package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYear_Validate(t *testing.T) {
	validAccount := Account{
		OpeningBalance: 0,
		ClosingBalance: 100,
		Entries: []Entry{
			{Amount: 100, Note: "Salary"},
		},
	}
	validMonth := Month{
		OpeningBalance: 0,
		ClosingBalance: 100,
		Accounts:       map[string]Account{"main": validAccount},
	}

	tests := []struct {
		name     string
		prevYear *Year
		year     Year
		yearNum  int
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid year with one month",
			year: Year{
				OpeningBalance: 0,
				ClosingBalance: 100,
				Months:         map[int]Month{1: validMonth},
			},
			yearNum: 1,
			wantErr: false,
		},
		{
			name: "year number < 1",
			year: Year{
				OpeningBalance: 0,
				ClosingBalance: 100,
				Months:         map[int]Month{1: validMonth},
			},
			yearNum: 0,
			wantErr: true,
			errMsg:  "Y-0: year number must be greater than 0",
		},
		{
			name: "year with no months",
			year: Year{
				OpeningBalance: 0,
				ClosingBalance: 0,
				Months:         map[int]Month{},
			},
			yearNum: 1,
			wantErr: true,
			errMsg:  "Y-4: year must contain at least one month",
		},
		{
			name: "opening balance mismatch",
			year: Year{
				OpeningBalance: 50,
				ClosingBalance: 100,
				Months:         map[int]Month{1: validMonth},
			},
			yearNum: 1,
			wantErr: true,
			errMsg:  "Y-2: year opening balance does not equal first month opening balance (expected: 0, got: 50)",
		},
		{
			name: "closing balance mismatch",
			year: Year{
				OpeningBalance: 0,
				ClosingBalance: 50,
				Months:         map[int]Month{1: validMonth},
			},
			yearNum: 1,
			wantErr: true,
			errMsg:  "Y-3: year closing balance does not equal last month closing balance (expected: 100, got: 50)",
		},
		{
			name: "consecutive years",
			prevYear: &Year{
				OpeningBalance: 100,
				ClosingBalance: 100,
				Months: map[int]Month{
					1: {
						OpeningBalance: 100,
						ClosingBalance: 100,
						Accounts: map[string]Account{
							"main": {
								OpeningBalance: 100,
								ClosingBalance: 100,
							},
						},
					},
				},
			},
			year: Year{
				OpeningBalance: 100,
				ClosingBalance: 100,
				Months: map[int]Month{
					1: {
						OpeningBalance: 100,
						ClosingBalance: 100,
						Accounts: map[string]Account{
							"main": {
								OpeningBalance: 100,
								ClosingBalance: 100,
							},
						},
					},
				},
			},
			yearNum: 2,
			wantErr: false,
		},
		{
			name: "consecutive years do not match",
			prevYear: &Year{
				OpeningBalance: 100,
				ClosingBalance: 100,
				Months: map[int]Month{
					1: {
						OpeningBalance: 100,
						ClosingBalance: 100,
						Accounts: map[string]Account{
							"main": {
								OpeningBalance: 100,
								ClosingBalance: 100,
							},
						},
					},
				},
			},
			year: Year{
				OpeningBalance: 200,
				ClosingBalance: 200,
				Months: map[int]Month{
					1: {
						OpeningBalance: 200,
						ClosingBalance: 200,
						Accounts: map[string]Account{
							"main": {
								OpeningBalance: 200,
								ClosingBalance: 200,
							},
						},
					},
				},
			},
			yearNum: 2,
			wantErr: true,
			errMsg:  "Y-1: year opening balance does not equal previous year closing balance (expected: 100, got: 200)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.year.Validate(tt.yearNum, tt.prevYear)
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

func TestYear_IncomeAndExpenses(t *testing.T) {
	month1 := Month{
		OpeningBalance: 0,
		ClosingBalance: 100,
		Accounts: map[string]Account{
			"main": {
				OpeningBalance: 0,
				ClosingBalance: 100,
				Entries: []Entry{
					{Amount: 200, Note: "Salary", Date: "2024-01-01"},
					{Amount: -100, Note: "Groceries", Date: "2024-01-02"},
				},
			},
		},
	}
	month2 := Month{
		OpeningBalance: 100,
		ClosingBalance: 250,
		Accounts: map[string]Account{
			"main": {
				OpeningBalance: 100,
				ClosingBalance: 250,
				Entries: []Entry{
					{Amount: 200, Note: "Bonus", Date: "2024-02-01"},
					{Amount: -50, Note: "Utilities", Date: "2024-02-02"},
				},
			},
		},
	}

	year := Year{
		OpeningBalance: 0,
		ClosingBalance: 250,
		Months:         map[int]Month{1: month1, 2: month2},
	}

	assert.Equal(t, 400, year.Income())
	assert.Equal(t, -150, year.Expenses())
}

func TestYear_GetMonthNumbers(t *testing.T) {
	year := Year{
		Months: map[int]Month{
			5: {},
			1: {},
			3: {},
		},
	}
	assert.Equal(t, []int{1, 3, 5}, year.GetMonthNumbers())
}
