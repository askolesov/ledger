package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLedger_Validate_YearLevelRules(t *testing.T) {
	t.Run("Y-1: Consecutive years must chain totals", func(t *testing.T) {
		// Valid case
		validLedger := Ledger{
			Years: map[int]*Year{
				2024: {
					OpeningBalance: 1000,
					ClosingBalance: 1200,
					Months: map[int]Month{
						1: {
							OpeningBalance: 1000,
							ClosingBalance: 1200,
							Accounts: map[string]Account{
								"Checking": {
									OpeningBalance: 0,
									ClosingBalance: 200,
									Entries: []Entry{
										{Amount: 200, Internal: false, Note: "Income", Date: "2024-01-15", Tag: "Salary"},
									},
								},
							},
						},
					},
				},
				2025: {
					OpeningBalance: 1200, // Must match previous year's closing
					ClosingBalance: 1300,
					Months: map[int]Month{
						1: {
							OpeningBalance: 1200,
							ClosingBalance: 1300,
							Accounts: map[string]Account{
								"Checking": {
									OpeningBalance: 200,
									ClosingBalance: 300,
									Entries: []Entry{
										{Amount: 100, Internal: false, Note: "Income", Date: "2025-01-15", Tag: "Salary"},
									},
								},
							},
						},
					},
				},
			},
		}
		err := validLedger.Validate()
		require.NoError(t, err)

		// Invalid case - broken chain
		invalidLedger := Ledger{
			Years: map[int]*Year{
				2024: {
					OpeningBalance: 1000,
					ClosingBalance: 1200,
					Months: map[int]Month{
						1: {
							OpeningBalance: 1000,
							ClosingBalance: 1200,
							Accounts: map[string]Account{
								"Checking": {
									OpeningBalance: 0,
									ClosingBalance: 200,
									Entries: []Entry{
										{Amount: 200, Internal: false, Note: "Income", Date: "2024-01-15", Tag: "Salary"},
									},
								},
							},
						},
					},
				},
				2025: {
					OpeningBalance: 1500, // Should be 1200
					ClosingBalance: 1600,
					Months: map[int]Month{
						1: {
							OpeningBalance: 1500,
							ClosingBalance: 1600,
							Accounts: map[string]Account{
								"Checking": {
									OpeningBalance: 500,
									ClosingBalance: 600,
									Entries: []Entry{
										{Amount: 100, Internal: false, Note: "Income", Date: "2025-01-15", Tag: "Salary"},
									},
								},
							},
						},
					},
				},
			},
		}
		err = invalidLedger.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "does not equal previous year closing balance")
	})

	t.Run("Y-2: Year balances must match first/last month", func(t *testing.T) {
		// Invalid opening balance
		invalidLedger := Ledger{
			Years: map[int]*Year{
				2025: {
					OpeningBalance: 1500, // Should be 1000
					ClosingBalance: 1200,
					Months: map[int]Month{
						1: {
							OpeningBalance: 1000,
							ClosingBalance: 1200,
							Accounts: map[string]Account{
								"Checking": {
									OpeningBalance: 0,
									ClosingBalance: 200,
									Entries: []Entry{
										{Amount: 200, Internal: false, Note: "Income", Date: "2025-01-15", Tag: "Salary"},
									},
								},
							},
						},
					},
				},
			},
		}
		err := invalidLedger.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "year opening balance")
	})
}

func TestLedger_Validate_MonthLevelRules(t *testing.T) {
	t.Run("M-1: Consecutive months must chain totals", func(t *testing.T) {
		invalidLedger := Ledger{
			Years: map[int]*Year{
				2025: {
					OpeningBalance: 1000,
					ClosingBalance: 1300,
					Months: map[int]Month{
						1: {
							OpeningBalance: 1000,
							ClosingBalance: 1200,
							Accounts: map[string]Account{
								"Checking": {
									OpeningBalance: 0,
									ClosingBalance: 200,
									Entries: []Entry{
										{Amount: 200, Internal: false, Note: "Income", Date: "2025-01-15", Tag: "Salary"},
									},
								},
							},
						},
						2: {
							OpeningBalance: 1500, // Should be 1200
							ClosingBalance: 1300,
							Accounts: map[string]Account{
								"Checking": {
									OpeningBalance: 500, // Should be 200
									ClosingBalance: 300,
									Entries: []Entry{
										{Amount: -200, Internal: false, Note: "Expense", Date: "2025-02-15", Tag: "Bills"},
									},
								},
							},
						},
					},
				},
			},
		}
		err := invalidLedger.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "does not equal month 2 opening balance")
	})

	t.Run("M-2: Month balances must equal sum of account balances", func(t *testing.T) {
		invalidLedger := Ledger{
			Years: map[int]*Year{
				2025: {
					OpeningBalance: 1000,
					ClosingBalance: 1200,
					Months: map[int]Month{
						1: {
							OpeningBalance: 1500, // Should be 1000 (sum of account opening balances)
							ClosingBalance: 1200,
							Accounts: map[string]Account{
								"Checking": {
									OpeningBalance: 0,
									ClosingBalance: 200,
									Entries: []Entry{
										{Amount: 200, Internal: false, Note: "Income", Date: "2025-01-15", Tag: "Salary"},
									},
								},
							},
						},
					},
				},
			},
		}
		err := invalidLedger.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "month opening balance")
	})

	t.Run("M-3: Internal entries must sum to zero", func(t *testing.T) {
		invalidLedger := Ledger{
			Years: map[int]*Year{
				2025: {
					OpeningBalance: 1000,
					ClosingBalance: 1000,
					Months: map[int]Month{
						1: {
							OpeningBalance: 1000,
							ClosingBalance: 1000,
							Accounts: map[string]Account{
								"Checking": {
									OpeningBalance: 600,
									ClosingBalance: 550,
									Entries: []Entry{
										{Amount: -50, Internal: true, Note: "Transfer to Savings", Date: "2025-01-15", Tag: "Transfer"},
									},
								},
								"Savings": {
									OpeningBalance: 400,
									ClosingBalance: 450,
									Entries: []Entry{
										{Amount: 100, Internal: true, Note: "Transfer from Checking", Date: "2025-01-15", Tag: "Transfer"}, // Should be 50
									},
								},
							},
						},
					},
				},
			},
		}
		err := invalidLedger.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "sum of internal entries must be 0")
	})
}

func TestLedger_Validate_AccountLevelRules(t *testing.T) {
	t.Run("A-1: Account balance equation", func(t *testing.T) {
		// This is tested in account_test.go, but let's test it in context
		invalidLedger := Ledger{
			Years: map[int]*Year{
				2025: {
					OpeningBalance: 1000,
					ClosingBalance: 1200,
					Months: map[int]Month{
						1: {
							OpeningBalance: 1000,
							ClosingBalance: 1200,
							Accounts: map[string]Account{
								"Checking": {
									OpeningBalance: 0,
									ClosingBalance: 300, // Should be 200
									Entries: []Entry{
										{Amount: 200, Internal: false, Note: "Income", Date: "2025-01-15", Tag: "Salary"},
									},
								},
							},
						},
					},
				},
			},
		}
		err := invalidLedger.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "does not equal closing balance")
	})

	t.Run("A-2: Account balance chaining across months", func(t *testing.T) {
		invalidLedger := Ledger{
			Years: map[int]*Year{
				2025: {
					OpeningBalance: 1000,
					ClosingBalance: 1100,
					Months: map[int]Month{
						1: {
							OpeningBalance: 1000,
							ClosingBalance: 1200,
							Accounts: map[string]Account{
								"Checking": {
									OpeningBalance: 0,
									ClosingBalance: 200,
									Entries: []Entry{
										{Amount: 200, Internal: false, Note: "Income", Date: "2025-01-15", Tag: "Salary"},
									},
								},
							},
						},
						2: {
							OpeningBalance: 1200,
							ClosingBalance: 1100,
							Accounts: map[string]Account{
								"Checking": {
									OpeningBalance: 300, // Should be 200
									ClosingBalance: 100,
									Entries: []Entry{
										{Amount: -200, Internal: false, Note: "Expense", Date: "2025-02-15", Tag: "Bills"},
									},
								},
							},
						},
					},
				},
			},
		}
		err := invalidLedger.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "previous month closing balance")
	})

	t.Run("A-3: New account must start with zero balance", func(t *testing.T) {
		invalidLedger := Ledger{
			Years: map[int]*Year{
				2025: {
					OpeningBalance: 1000,
					ClosingBalance: 1200,
					Months: map[int]Month{
						1: {
							OpeningBalance: 1000,
							ClosingBalance: 1200,
							Accounts: map[string]Account{
								"Checking": {
									OpeningBalance: 100, // Should be 0 for new account
									ClosingBalance: 300,
									Entries: []Entry{
										{Amount: 200, Internal: false, Note: "Income", Date: "2025-01-15", Tag: "Salary"},
									},
								},
							},
						},
					},
				},
			},
		}
		err := invalidLedger.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "first month account must start with opening balance 0")
	})

	t.Run("A-3: Account can only be omitted if last balance was zero", func(t *testing.T) {
		invalidLedger := Ledger{
			Years: map[int]*Year{
				2025: {
					OpeningBalance: 1000,
					ClosingBalance: 1000,
					Months: map[int]Month{
						1: {
							OpeningBalance: 1000,
							ClosingBalance: 1000,
							Accounts: map[string]Account{
								"Checking": {
									OpeningBalance: 0,
									ClosingBalance: 100, // Non-zero balance
									Entries: []Entry{
										{Amount: 100, Internal: false, Note: "Income", Date: "2025-01-15", Tag: "Salary"},
									},
								},
								"Savings": {
									OpeningBalance: 0,
									ClosingBalance: -100,
									Entries: []Entry{
										{Amount: -100, Internal: false, Note: "Expense", Date: "2025-01-15", Tag: "Bills"},
									},
								},
							},
						},
						2: {
							OpeningBalance: 1000,
							ClosingBalance: 1000,
							Accounts: map[string]Account{
								// Checking account omitted but had non-zero balance
								"Savings": {
									OpeningBalance: -100,
									ClosingBalance: -100,
									Entries:        []Entry{},
								},
							},
						},
					},
				},
			},
		}
		err := invalidLedger.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be omitted because previous month closing balance is 100")
	})
}

func TestLedger_Validate_ValidComplexExample(t *testing.T) {
	// Test the example from the OLF specification
	validLedger := Ledger{
		Years: map[int]*Year{
			2025: {
				OpeningBalance: 1000,
				ClosingBalance: 1050,
				Months: map[int]Month{
					1: {
						OpeningBalance: 1000,
						ClosingBalance: 1050,
						Accounts: map[string]Account{
							"Checking": {
								OpeningBalance: 600,
								ClosingBalance: 620,
								Entries: []Entry{
									{
										Amount:   50,
										Internal: false,
										Note:     "Salary (January)",
										Date:     "2025-01-28",
										Tag:      "Income",
									},
									{
										Amount:   -30,
										Internal: true,
										Note:     "Transfer to Savings",
										Date:     "2025-01-30",
										Tag:      "Transfer",
									},
								},
							},
							"Savings": {
								OpeningBalance: 400,
								ClosingBalance: 430,
								Entries: []Entry{
									{
										Amount:   30,
										Internal: true,
										Note:     "Transfer from Checking",
										Date:     "2025-01-30",
										Tag:      "Transfer",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	err := validLedger.Validate()
	require.NoError(t, err)
}

func TestLedger_Validate_ConsecutiveYears(t *testing.T) {
	t.Run("valid consecutive years", func(t *testing.T) {
		ledger := Ledger{
			Years: map[int]*Year{
				2023: {OpeningBalance: 1000, ClosingBalance: 1100, Months: map[int]Month{}},
				2024: {OpeningBalance: 1100, ClosingBalance: 1200, Months: map[int]Month{}},
				2025: {OpeningBalance: 1200, ClosingBalance: 1300, Months: map[int]Month{}},
			},
		}
		err := ledger.Validate()
		require.NoError(t, err)
	})

	t.Run("gap in years", func(t *testing.T) {
		ledger := Ledger{
			Years: map[int]*Year{
				2023: {OpeningBalance: 1000, ClosingBalance: 1100, Months: map[int]Month{}},
				2025: {OpeningBalance: 1100, ClosingBalance: 1200, Months: map[int]Month{}}, // Missing 2024
			},
		}
		err := ledger.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "years must be consecutive")
	})
}
