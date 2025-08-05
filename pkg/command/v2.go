package command

import (
	"fmt"
	v2 "ledger/pkg/ledger/v2"
	"os"
	"sort"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func getV2ValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <file>",
		Short: "Validate OLF v2.0 file for structural correctness and data integrity",
		Long: `Validate OLF v2.0 file for structural correctness and data integrity.

Performs comprehensive validation including:
- YAML/JSON structure validation
- Data type validation
- Balance calculations and consistency checks
- Double-entry bookkeeping constraints
- Account continuity validation
- Cross-period balance verification

Examples:
  ledger validate ledger.yaml          # Validate OLF v2.0 YAML file
  ledger validate ledger.json          # Validate OLF v2.0 JSON file
  ledger validate /path/to/ledger.yaml`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			ledger, err := v2.ReadLedger(path)
			if err != nil {
				return fmt.Errorf("failed to read ledger file: %w", err)
			}

			cmd.Printf("Successfully loaded ledger with %d year(s)\n", len(ledger.Years))

			err = ledger.Validate()
			if err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			cmd.Println("✓ Ledger is valid according to OLF v2.0 specification")

			// Print summary statistics
			cmd.Printf("Total Income: %.2f\n", float64(ledger.Income())/1000)
			cmd.Printf("Total Expenses: %.2f\n", float64(ledger.Expenses())/1000)

			return nil
		},
	}
}

func getV2ReportCmd() *cobra.Command {
	var short bool

	cmd := &cobra.Command{
		Use:   "report <file>",
		Short: "Generate financial reports from OLF v2.0 file",
		Long: `Generate financial reports from OLF v2.0 file.

Reads, validates, and generates detailed financial reports showing income,
expenses, and balance information organized by year and month.

Default report shows detailed table with:
- Opening balance for each month
- Total income for the month
- Total expenses for the month
- Closing balance for the month

Use --short flag for condensed view showing only monthly expenses.

Examples:
  ledger report ledger.yaml            # Generate detailed monthly report
  ledger report ledger.json --short    # Generate condensed expense report
  ledger report ledger.yaml -s         # Short form of --short flag`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			ledger, err := v2.ReadLedger(path)
			if err != nil {
				return fmt.Errorf("failed to read ledger file: %w", err)
			}

			err = ledger.Validate()
			if err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			if short {
				v2ShortMonthlyReport(ledger)
				return nil
			}

			v2MonthlyReport(ledger)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&short, "short", "s", false, "Generate condensed report showing only monthly expenses")

	return cmd
}

func v2MonthlyReport(ledger v2.Ledger) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Year", "Month", "Opening", "Income", "Expenses", "Closing"})

	// Get sorted year numbers
	yearNums := lo.Keys(ledger.Years)
	sort.Ints(yearNums)

	for _, yearNum := range yearNums {
		year := ledger.Years[yearNum]
		t.AppendSeparator()

		// Get sorted month numbers for this year
		monthNums := year.GetMonthNumbers()

		for _, monthNum := range monthNums {
			month := year.Months[monthNum]
			t.AppendRow(table.Row{
				yearNum,
				monthNum,
				float64(month.OpeningBalance) / 1000,
				float64(month.Income()) / 1000,
				float64(month.Expenses()) / 1000,
				float64(month.ClosingBalance) / 1000,
			})
		}
	}

	t.Render()
}

func v2ShortMonthlyReport(ledger v2.Ledger) {
	// Get sorted year numbers
	yearNums := lo.Keys(ledger.Years)
	sort.Ints(yearNums)

	for _, yearNum := range yearNums {
		year := ledger.Years[yearNum]

		// Get sorted month numbers for this year
		monthNums := year.GetMonthNumbers()

		for _, monthNum := range monthNums {
			month := year.Months[monthNum]
			fmt.Printf("%.1f\n", float64(month.Expenses())/1000)
		}
	}
}
