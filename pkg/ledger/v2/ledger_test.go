package v2

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestReadLedger(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	t.Run("read YAML file", func(t *testing.T) {
		yamlContent := `years:
  2025:
    opening_balance: 1000
    closing_balance: 1050
    months:
      1:
        opening_balance: 1000
        closing_balance: 1050
        accounts:
          Checking:
            opening_balance: 600
            closing_balance: 620
            entries:
              - amount: 50
                internal: false
                note: "Salary"
                date: "2025-01-28"
                tag: "Income"
          Savings:
            opening_balance: 400
            closing_balance: 430
            entries:
              - amount: 30
                internal: false
                note: "Interest"
                date: "2025-01-31"
                tag: "Income"
`

		yamlFile := filepath.Join(tempDir, "test.yaml")
		err := os.WriteFile(yamlFile, []byte(yamlContent), 0666)
		require.NoError(t, err)

		ledger, err := ReadLedger(yamlFile)
		require.NoError(t, err)

		// Validate structure
		assert.Equal(t, 1000, ledger.Years[2025].OpeningBalance)
		assert.Equal(t, 1050, ledger.Years[2025].ClosingBalance)
		assert.Equal(t, 600, ledger.Years[2025].Months[1].Accounts["Checking"].OpeningBalance)
		assert.Equal(t, "Salary", ledger.Years[2025].Months[1].Accounts["Checking"].Entries[0].Note)
	})

	t.Run("read JSON file", func(t *testing.T) {
		jsonContent := `{
  "years": {
    "2025": {
      "opening_balance": 1000,
      "closing_balance": 1050,
      "months": {
        "1": {
          "opening_balance": 1000,
          "closing_balance": 1050,
          "accounts": {
            "Checking": {
              "opening_balance": 600,
              "closing_balance": 620,
              "entries": [
                {
                  "amount": 50,
                  "internal": false,
                  "note": "Salary",
                  "date": "2025-01-28",
                  "tag": "Income"
                }
              ]
            },
            "Savings": {
              "opening_balance": 400,
              "closing_balance": 430,
              "entries": [
                {
                  "amount": 30,
                  "internal": false,
                  "note": "Interest",
                  "date": "2025-01-31",
                  "tag": "Income"
                }
              ]
            }
          }
        }
      }
    }
  }
}`

		jsonFile := filepath.Join(tempDir, "test.json")
		err := os.WriteFile(jsonFile, []byte(jsonContent), 0644)
		require.NoError(t, err)

		ledger, err := ReadLedger(jsonFile)
		require.NoError(t, err)

		// Validate structure
		assert.Equal(t, 1000, ledger.Years[2025].OpeningBalance)
		assert.Equal(t, 1050, ledger.Years[2025].ClosingBalance)
		assert.Equal(t, 600, ledger.Years[2025].Months[1].Accounts["Checking"].OpeningBalance)
		assert.Equal(t, "Salary", ledger.Years[2025].Months[1].Accounts["Checking"].Entries[0].Note)
	})

	t.Run("unsupported file format", func(t *testing.T) {
		txtFile := filepath.Join(tempDir, "test.txt")
		err := os.WriteFile(txtFile, []byte("some content"), 0644)
		require.NoError(t, err)

		_, err = ReadLedger(txtFile)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported file format")
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := ReadLedger("nonexistent.yaml")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read file")
	})

	t.Run("invalid YAML", func(t *testing.T) {
		invalidYamlFile := filepath.Join(tempDir, "invalid.yaml")
		err := os.WriteFile(invalidYamlFile, []byte("invalid: yaml: content: ["), 0644)
		require.NoError(t, err)

		_, err = ReadLedger(invalidYamlFile)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse file")
	})
}

func TestWriteLedger(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "ledger_write_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test data
	testLedger := Ledger{
		Years: map[int]Year{
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
								ClosingBalance: 650,
								Entries: []Entry{
									{
										Amount:   50,
										Internal: false,
										Note:     "Salary",
										Date:     "2025-01-28",
										Tag:      "Income",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	t.Run("write YAML file", func(t *testing.T) {
		yamlFile := filepath.Join(tempDir, "output.yaml")
		err := WriteLedger(testLedger, yamlFile)
		require.NoError(t, err)

		// Read it back and verify
		ledger, err := ReadLedger(yamlFile)
		require.NoError(t, err)
		assert.Equal(t, testLedger, ledger)
	})

	t.Run("write JSON file", func(t *testing.T) {
		jsonFile := filepath.Join(tempDir, "output.json")
		err := WriteLedger(testLedger, jsonFile)
		require.NoError(t, err)

		// Read it back and verify
		ledger, err := ReadLedger(jsonFile)
		require.NoError(t, err)
		assert.Equal(t, testLedger, ledger)
	})

	t.Run("unsupported file format", func(t *testing.T) {
		txtFile := filepath.Join(tempDir, "output.txt")
		err := WriteLedger(testLedger, txtFile)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported file format")
	})
}

func TestLedger_String(t *testing.T) {
	ledger := Ledger{
		Years: map[int]Year{
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
								ClosingBalance: 650,
								Entries: []Entry{
									{
										Amount:   50,
										Internal: false,
										Note:     "Salary",
										Date:     "2025-01-28",
										Tag:      "Income",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	yamlBytes, err := yaml.Marshal(ledger)
	yamlStr := string(yamlBytes)
	require.NoError(t, err)
	assert.Contains(t, yamlStr, "years:")
	assert.Contains(t, yamlStr, "2025:")
	assert.Contains(t, yamlStr, "opening_balance: 1000")
	assert.Contains(t, yamlStr, "Salary")
}

func TestLedger_IncomeAndExpenses(t *testing.T) {
	ledger := Ledger{
		Years: map[int]Year{
			2025: {
				OpeningBalance: 1000,
				ClosingBalance: 1150,
				Months: map[int]Month{
					1: {
						OpeningBalance: 1000,
						ClosingBalance: 1150,
						Accounts: map[string]Account{
							"Checking": {
								OpeningBalance: 600,
								ClosingBalance: 750,
								Entries: []Entry{
									{Amount: 200, Internal: false, Note: "Salary", Date: "2025-01-15", Tag: "Income"},
									{Amount: -50, Internal: false, Note: "Groceries", Date: "2025-01-20", Tag: "Expense"},
									{Amount: -30, Internal: true, Note: "Transfer", Date: "2025-01-25", Tag: "Transfer"},
								},
							},
							"Savings": {
								OpeningBalance: 400,
								ClosingBalance: 400,
								Entries: []Entry{
									{Amount: 30, Internal: true, Note: "Transfer", Date: "2025-01-25", Tag: "Transfer"},
								},
							},
						},
					},
				},
			},
		},
	}

	income := ledger.Income()
	expenses := ledger.Expenses()

	assert.Equal(t, 200, income)   // Only non-internal positive amounts
	assert.Equal(t, -50, expenses) // Only non-internal negative amounts
}

func TestLedger_GetYearNumbers(t *testing.T) {
	ledger := Ledger{
		Years: map[int]Year{
			2025: {OpeningBalance: 1000, ClosingBalance: 1100, Months: map[int]Month{}},
			2023: {OpeningBalance: 800, ClosingBalance: 900, Months: map[int]Month{}},
			2024: {OpeningBalance: 900, ClosingBalance: 1000, Months: map[int]Month{}},
		},
	}

	yearNums := lo.Keys(ledger.Years)
	sort.Ints(yearNums)
	assert.Equal(t, []int{2023, 2024, 2025}, yearNums) // Should be sorted
}
