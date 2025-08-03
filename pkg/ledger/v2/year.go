package v2

import (
	"fmt"
	"sort"

	"github.com/samber/lo"
)

// Year represents a year with months
type Year struct {
	OpeningBalance int           `json:"opening_balance" yaml:"opening_balance" toml:"opening_balance"`
	ClosingBalance int           `json:"closing_balance" yaml:"closing_balance" toml:"closing_balance"`
	Months         map[int]Month `json:"months" yaml:"months" toml:"months"`
}

// Validate validates a year according to OLF v2.0 rules
func (y Year) Validate(yearNum int, prevYear *Year) error {
	if yearNum < 1 {
		return fmt.Errorf("year number must be greater than 0, got %d", yearNum)
	}

	if len(y.Months) == 0 {
		return fmt.Errorf("year %d has no months", yearNum)
	}

	// Validate months
	var prevMonth *Month

	if prevYear != nil {
		// Y-1: Consecutive years must chain totals: prev.closing_balance = next.opening_balance
		if y.OpeningBalance != prevYear.ClosingBalance {
			return fmt.Errorf("year opening balance %d does not equal previous year closing balance %d",
				y.OpeningBalance, prevYear.ClosingBalance)
		}

		prevYearMonthNums := lo.Keys(prevYear.Months)
		sort.Ints(prevYearMonthNums)

		lastMonthNum := prevYearMonthNums[len(prevYearMonthNums)-1]
		prevMonthValue := prevYear.Months[lastMonthNum]
		prevMonth = &prevMonthValue
	}

	monthNums := lo.Keys(y.Months)
	sort.Ints(monthNums)

	for _, monthNum := range monthNums {
		month := y.Months[monthNum]

		if err := month.Validate(yearNum, monthNum, prevMonth); err != nil {
			return fmt.Errorf("month %d: %w", monthNum, err)
		}

		prevMonth = &month
	}

	// Y-2: Year opening_balance equals first month's opening_balance
	firstMonth := y.Months[monthNums[0]]
	if y.OpeningBalance != firstMonth.OpeningBalance {
		return fmt.Errorf("year opening balance %d does not equal first month opening balance %d",
			y.OpeningBalance, firstMonth.OpeningBalance)
	}

	// Y-2: Year closing_balance equals last month's closing_balance
	lastMonth := y.Months[monthNums[len(monthNums)-1]]
	if y.ClosingBalance != lastMonth.ClosingBalance {
		return fmt.Errorf("year closing balance %d does not equal last month closing balance %d",
			y.ClosingBalance, lastMonth.ClosingBalance)
	}

	return nil
}

// Income returns the sum of income from all months
func (y Year) Income() int {
	return lo.SumBy(lo.Values(y.Months), func(month Month) int {
		return month.Income()
	})
}

// Expenses returns the sum of expenses from all months
func (y Year) Expenses() int {
	return lo.SumBy(lo.Values(y.Months), func(month Month) int {
		return month.Expenses()
	})
}

// GetMonthNumbers returns sorted list of month numbers
func (y Year) GetMonthNumbers() []int {
	monthNums := lo.Keys(y.Months)
	sort.Ints(monthNums)
	return monthNums
}
