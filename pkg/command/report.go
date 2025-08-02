package command

import (
	"fmt"
	"os"

	v1 "go-finances/pkg/ledger/v1"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func getReportCmd() *cobra.Command {
	var short bool

	cmd := &cobra.Command{
		Use:   "report",
		Short: "Generate a report of your finances",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			data, err := v1.ReadData(path)
			if err != nil {
				return err
			}

			err = data.Validate()
			if err != nil {
				return err
			}

			// check for short flag and call the appropriate function

			if short {
				ShortMonthlyReport(data)
				return nil
			}

			MonthlyReport(data)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&short, "short", "s", false, "Print a shorter version of the report")

	return cmd
}

func MonthlyReport(d v1.Data) {
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

func ShortMonthlyReport(d v1.Data) {
	for _, year := range d.Years {
		for _, month := range year.Months {
			fmt.Printf("%.1f\n", float64(month.Expenses())/1000)
		}
	}
}
