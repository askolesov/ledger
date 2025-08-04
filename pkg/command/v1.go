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
		Short: "Commands for Open Ledger Format v1.0 (legacy)",
		Long: `Commands for Open Ledger Format v1.0 (legacy format).

Provides backward compatibility with existing OLF v1.0 ledger files.
Supports YAML format only.`,
	}

	v1Cmd.AddCommand(getV1ValidateCmd())
	v1Cmd.AddCommand(getV1ReportCmd())

	return v1Cmd
}

func getV1ValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <file>",
		Short: "Validate OLF v1.0 file for structural correctness and data integrity",
		Long: `Validate OLF v1.0 file for structural correctness and data integrity.

Performs comprehensive validation including:
- YAML structure validation
- Data type validation  
- Balance calculations and consistency checks
- Transaction integrity verification

Examples:
  ledger v1 validate data.yaml         # Validate OLF v1.0 file
  ledger v1 validate /path/to/data.yaml`,
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
		Short: "Generate financial reports from OLF v1.0 file",
		Long: `Generate financial reports from OLF v1.0 file.

Reads, validates, and generates detailed financial reports showing income,
expenses, and balance information organized by year and month.

Default report shows detailed table with:
- Opening balance for each month
- Total income for the month
- Total expenses for the month  
- Closing balance for the month

Use --short flag for condensed view showing only monthly expenses.

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

	cmd.Flags().BoolVarP(&short, "short", "s", false, "Generate condensed report showing only monthly expenses")

	return cmd
}

func v1MonthlyReport(d v1.Data) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Year", "Month", "Opening", "Income", "Expenses", "Closing"})

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
