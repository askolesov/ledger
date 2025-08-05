package command

import (
	"fmt"
	v1 "ledger/pkg/ledger/v1"
	v2 "ledger/pkg/ledger/v2"
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
	v1Cmd.AddCommand(getV1MigrateCmd())

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

func getV1MigrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate <input-file> <output-file>",
		Short: "Migrate OLF v1.0 file to OLF v2.0 format",
		Long: `Migrate OLF v1.0 file to OLF v2.0 format.

Reads an existing OLF v1.0 ledger file, validates it, converts it to the new
OLF v2.0 format, and writes the result to the specified output file.

The migration process:
- Validates the input v1.0 file for correctness
- Converts the data structure from v1.0 to v2.0 format
- Maps wallets to accounts and transactions to entries
- Preserves all financial data and balances
- Validates the converted v2.0 data before writing

Examples:
  ledger v1 migrate data-v1.yaml data-v2.yaml    # Migrate YAML file
  ledger v1 migrate data-v1.json data-v2.json    # Migrate JSON file`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputPath := args[0]
			outputPath := args[1]

			return migrateV1ToV2(cmd, inputPath, outputPath)
		},
	}
}

func migrateV1ToV2(cmd *cobra.Command, inputPath, outputPath string) error {
	// Read and validate v1 data
	cmd.Printf("Reading v1.0 ledger file: %s\n", inputPath)
	v1Data, err := v1.ReadData(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read v1.0 ledger file: %w", err)
	}

	cmd.Println("Validating v1.0 data...")
	err = v1Data.Validate()
	if err != nil {
		return fmt.Errorf("v1.0 data validation failed: %w", err)
	}

	// Convert v1 to v2
	cmd.Println("Converting v1.0 to v2.0 format...")
	v2Ledger := convertV1ToV2(v1Data)

	// Validate v2 data
	cmd.Println("Validating converted v2.0 data...")
	err = v2Ledger.Validate()
	if err != nil {
		return fmt.Errorf("converted v2.0 data validation failed: %w", err)
	}

	// Write v2 data
	cmd.Printf("Writing v2.0 ledger file: %s\n", outputPath)
	err = v2.WriteLedger(v2Ledger, outputPath)
	if err != nil {
		return fmt.Errorf("failed to write v2.0 ledger file: %w", err)
	}

	cmd.Println("✓ Migration completed successfully")
	return nil
}

func convertV1ToV2(v1Data v1.Data) v2.Ledger {
	v2Years := make(map[int]v2.Year)

	for _, v1Year := range v1Data.Years {
		v2Months := make(map[int]v2.Month)

		for _, v1Month := range v1Year.Months {
			v2Accounts := make(map[string]v2.Account)

			for _, v1Wallet := range v1Month.Wallets {
				var v2Entries []v2.Entry

				for _, v1Transaction := range v1Wallet.Transactions {
					v2Entry := v2.Entry{
						Amount:   v1Transaction.Amount,
						Internal: v1Transaction.IsInternal,
						Note:     v1Transaction.Comment,
						Date:     v1Transaction.Date,
						Tag:      v1Transaction.Category,
					}
					v2Entries = append(v2Entries, v2Entry)
				}

				v2Account := v2.Account{
					OpeningBalance: v1Wallet.StartingBalance,
					ClosingBalance: v1Wallet.EndingBalance,
					Entries:        v2Entries,
				}

				v2Accounts[v1Wallet.Name] = v2Account
			}

			v2Month := v2.Month{
				OpeningBalance: v1Month.StartingBalance,
				ClosingBalance: v1Month.EndingBalance,
				Accounts:       v2Accounts,
			}

			v2Months[v1Month.Number] = v2Month
		}

		v2Year := v2.Year{
			OpeningBalance: v1Year.StartingBalance,
			ClosingBalance: v1Year.EndingBalance,
			Months:         v2Months,
		}

		v2Years[v1Year.Number] = v2Year
	}

	return v2.Ledger{
		Years: v2Years,
	}
}
