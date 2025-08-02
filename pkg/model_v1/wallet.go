package model_v1

import (
	"fmt"
)

type Wallet struct {
	Name         string        `json:"name" yaml:"name" toml:"name"`
	Transactions []Transaction `json:"transactions" yaml:"transactions" toml:"transactions"`

	StartingBalance int `json:"starting_balance" yaml:"starting_balance" toml:"starting_balance"`
	EndingBalance   int `json:"ending_balance" yaml:"ending_balance" toml:"ending_balance"`
}

func (w Wallet) Validate() error {
	// validate transactions
	for index, transaction := range w.Transactions {
		err := transaction.Validate()
		if err != nil {
			return fmt.Errorf("transaction %d: %w", index, err)
		}
	}

	// validate balances
	balance := w.StartingBalance
	for _, transaction := range w.Transactions {
		balance += transaction.Amount
	}
	if balance != w.EndingBalance {
		return fmt.Errorf("starting amount + transactions %d does not equal ending amount %d", balance, w.EndingBalance)
	}

	return nil
}

func (w Wallet) InternalTransactionsSum() int {
	sum := 0
	for _, transaction := range w.Transactions {
		if transaction.IsInternal {
			sum += transaction.Amount
		}
	}
	return sum
}

func (w Wallet) Income() int {
	sum := 0

	for _, transaction := range w.Transactions {
		if transaction.IsInternal {
			continue
		}

		if transaction.Amount > 0 {
			sum += transaction.Amount
		}
	}

	return sum
}

func (w Wallet) Expenses() int {
	sum := 0

	for _, transaction := range w.Transactions {
		if transaction.IsInternal {
			continue
		}

		if transaction.Amount < 0 {
			sum += transaction.Amount
		}
	}

	return sum
}
