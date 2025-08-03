package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntry_Validate(t *testing.T) {
	tests := []struct {
		name    string
		entry   Entry
		year    int
		month   int
		wantErr bool
		errMsg  string
	}{
		{
			name: "minimalvalid entry",
			entry: Entry{
				Amount: 100,
				Note:   "Salary payment",
			},
			wantErr: false,
		},
		{
			name: "full valid entry",
			entry: Entry{
				Amount:   100,
				Internal: true,
				Note:     "Salary payment",
				Date:     "2025-01-15",
				Tag:      "Income",
			},
			year:    2025,
			month:   1,
			wantErr: false,
		},
		{
			name: "missing amount",
			entry: Entry{
				Internal: false,
				Note:     "Salary payment",
				Date:     "2025-01-15",
				Tag:      "Income",
			},
			year:    2025,
			month:   1,
			wantErr: true,
			errMsg:  "amount: cannot be blank",
		},
		{
			name: "missing note",
			entry: Entry{
				Amount:   100,
				Internal: false,
				Date:     "2025-01-15",
				Tag:      "Income",
			},
			year:    2025,
			month:   1,
			wantErr: true,
			errMsg:  "note: cannot be blank",
		},
		{
			name: "invalid date format",
			entry: Entry{
				Amount:   100,
				Internal: false,
				Note:     "Salary payment",
				Date:     "2025/01/15",
				Tag:      "Income",
			},
			year:    2025,
			month:   1,
			wantErr: true,
			errMsg:  "date: E-3: date format must be YYYY-MM-DD",
		},
		{
			name: "invalid date value",
			entry: Entry{
				Amount:   100,
				Internal: false,
				Note:     "Salary payment",
				Date:     "2025-13-45",
				Tag:      "Income",
			},
			year:    2025,
			month:   1,
			wantErr: true,
			errMsg:  "date format",
		},
		{
			name: "year mismatch",
			entry: Entry{
				Amount:   100,
				Internal: false,
				Note:     "Salary payment",
				Date:     "2024-01-15",
				Tag:      "Income",
			},
			year:    2025,
			month:   1,
			wantErr: true,
			errMsg:  "E-1: entry date year does not match expected year (expected: 2025, got: 2024)",
		},
		{
			name: "month mismatch",
			entry: Entry{
				Amount:   100,
				Internal: false,
				Note:     "Salary payment",
				Date:     "2025-02-15",
				Tag:      "Income",
			},
			year:    2025,
			month:   1,
			wantErr: true,
			errMsg:  "E-1: entry date month does not match expected month (expected: 1, got: 2)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.entry.Validate(tt.year, tt.month)
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
