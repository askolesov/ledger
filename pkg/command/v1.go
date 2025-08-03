package command

import (
	"fmt"
	v1 "go-finances/pkg/ledger/v1"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func getV1Cmd() *cobra.Command {
	var v1Cmd = &cobra.Command{
		Use:   "v1",
		Short: "Commands for working with Open Ledger Format v1.0",
		Long: `Commands for working with Open Ledger Format v1.0 (legacy format).

The v1 subcommand provides access to the original OLF v1.0 functionality
for backward compatibility with existing ledger files.`,
	}

	v1Cmd.AddCommand(getV1ValidateCmd())
	v1Cmd.AddCommand(getV1ReportCmd())

	return v1Cmd
}

func getV1ValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <file>",
		Short: "Validate an Open Ledger Format v1.0 file",
		Long: `Validate an Open Ledger Format v1.0 file for structural correctness and data integrity.

This command reads the specified ledger file and performs comprehensive validation
according to OLF v1.0 rules, including:
- File format validation (YAML structure)
- Data type validation
- Balance calculations and consistency checks
- Transaction integrity verification

Supported file formats: YAML (.yaml, .yml)

Examples:
  ledger v1 validate data.yaml         # Validate OLF v1.0 file
  ledger v1 validate /path/to/ledger.yaml`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			data, err := v1.ReadData(path)
			if err != nil {
				return fmt.Errorf("failed to read ledger file: %w", err)
			}

			cmd.Println("Successfully loaded ledger data")
			cmd.Println(data)

			err = data.Validate()
			if err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			cmd.Println("✓ Data is valid")
			return nil
		},
	}
}

func getV1ReportCmd() *cobra.Command {
	var short bool

	cmd := &cobra.Command{
		Use:   "report <file>",
		Short: "Generate financial reports from an Open Ledger Format v1.0 file",
		Long: `Generate comprehensive financial reports from an Open Ledger Format v1.0 file.

This command reads the specified ledger file, validates it, and generates
detailed financial reports showing income, expenses, and balance information
organized by year and month.

The default report shows a detailed table with:
- Starting balance for each month
- Total income for the month
- Total expenses for the month  
- Ending balance for the month

Use the --short flag for a condensed view showing only monthly expenses.

Supported file formats: YAML (.yaml, .yml)

Examples:
  ledger v1 report data.yaml           # Generate detailed monthly report
  ledger v1 report data.yaml --short   # Generate condensed expense report
  ledger v1 report data.yaml -s        # Short form of --short flag`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			data, err := v1.ReadData(path)
			if err != nil {
				return fmt.Errorf("failed to read ledger file: %w", err)
			}

			err = data.Validate()
			if err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			if short {
				v1ShortMonthlyReport(data)
				return nil
			}

			v1MonthlyReport(data)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&short, "short", "s", false, "Generate a condensed report showing only monthly expenses")

	return cmd
}

func v1MonthlyReport(d v1.Data) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Year", "Month", "Starting", "Income", "Expenses", "Ending"})

	for _, year := range d.Years {
		t.AppendSeparator()

		for _, month := range year.Months {
			t.AppendRow(table.Row{
				year.Number,
				month.Number,
				float64(month.StartingBalance) / 1000,
				float64(month.Income()) / 1000,
				float64(month.Expenses()) / 1000,
				float64(month.EndingBalance) / 1000,
			})
		}
	}

	t.Render()
}

func v1ShortMonthlyReport(d v1.Data) {
	for _, year := range d.Years {
		for _, month := range year.Months {
			fmt.Printf("%.1f\n", float64(month.Expenses())/1000)
		}
	}
}
