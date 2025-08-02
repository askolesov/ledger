package v1

import (
	"fmt"

	"github.com/samber/lo"
)

type Month struct {
	Number  int      `json:"number" yaml:"number" toml:"number"`
	Wallets []Wallet `json:"wallets" yaml:"wallets" toml:"wallets"`

	StartingBalance int `json:"starting_balance" yaml:"starting_balance" toml:"starting_balance"`
	EndingBalance   int `json:"ending_balance" yaml:"ending_balance" toml:"ending_balance"`
}

func (m Month) Validate() error {
	// validate number
	if m.Number < 1 || m.Number > 12 {
		return fmt.Errorf("month number must be between 1 and 12")
	}

	// validate wallets
	for _, wallet := range m.Wallets {
		err := wallet.Validate()
		if err != nil {
			return fmt.Errorf("wallet %s: %w", wallet.Name, err)
		}
	}

	// starting balance
	walletsStarting := lo.SumBy(m.Wallets, func(wallet Wallet) int {
		return wallet.StartingBalance
	})

	if m.StartingBalance != walletsStarting {
		return fmt.Errorf("month starting amount %d doesn't equal sum of wallet starting amounts %d",
			m.StartingBalance, walletsStarting)
	}

	// ending balance
	walletsEnding := lo.SumBy(m.Wallets, func(wallet Wallet) int {
		return wallet.EndingBalance
	})

	if m.EndingBalance != walletsEnding {
		return fmt.Errorf("month ending amount %d doesn't equal sum of wallet ending amounts %d",
			m.EndingBalance, walletsEnding)
	}

	// internal transactions
	internalTransactions := lo.SumBy(m.Wallets, func(wallet Wallet) int {
		return wallet.InternalTransactionsSum()
	})

	if internalTransactions != 0 {
		return fmt.Errorf("month internal transactions sum must be 0, got %d", internalTransactions)
	}

	return nil
}

func (m Month) Income() int {
	return lo.SumBy(m.Wallets, func(wallet Wallet) int {
		return wallet.Income()
	})
}

func (m Month) Expenses() int {
	return lo.SumBy(m.Wallets, func(wallet Wallet) int {
		return wallet.Expenses()
	})
}
