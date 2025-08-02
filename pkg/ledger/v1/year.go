package v1

import (
	"fmt"

	"github.com/samber/lo"
)

type Year struct {
	Number int     `json:"number" yaml:"number" toml:"number"`
	Months []Month `json:"months" yaml:"months" toml:"months"`

	StartingBalance int `json:"starting_balance" yaml:"starting_balance" toml:"starting_balance"`
	EndingBalance   int `json:"ending_balance" yaml:"ending_balance" toml:"ending_balance"`
}

func (y Year) Validate() error {
	// validate number
	if y.Number < 1 {
		return fmt.Errorf("year number must be greater than 0")
	}

	// validate months order
	for index, month := range y.Months {
		if index == 0 {
			continue
		}
		if month.Number != y.Months[index-1].Number+1 {
			return fmt.Errorf("months must be in ascending order without gaps")
		}
		if month.StartingBalance != y.Months[index-1].EndingBalance {
			return fmt.Errorf("month %d starting amount %d doesn't equal previous month %d ending amount %d",
				month.Number, month.StartingBalance, y.Months[index-1].Number, y.Months[index-1].EndingBalance)
		}
	}

	// validate starting balance
	if len(y.Months) > 0 && y.StartingBalance != y.Months[0].StartingBalance {
		return fmt.Errorf("year starting amount %d doesn't equal first month starting amount %d",
			y.StartingBalance, y.Months[0].StartingBalance)
	}

	// validate ending balance
	if len(y.Months) > 0 && y.EndingBalance != y.Months[len(y.Months)-1].EndingBalance {
		return fmt.Errorf("year ending amount %d doesn't equal last month ending amount %d",
			y.EndingBalance, y.Months[len(y.Months)-1].EndingBalance)
	}

	// validate months
	for _, month := range y.Months {
		err := month.Validate()
		if err != nil {
			return fmt.Errorf("month %d: %w", month.Number, err)
		}
	}

	return nil
}

func (y Year) Income() int {
	return lo.SumBy(y.Months, func(month Month) int {
		return month.Income()
	})
}

func (y Year) Expenses() int {
	return lo.SumBy(y.Months, func(month Month) int {
		return month.Expenses()
	})
}
