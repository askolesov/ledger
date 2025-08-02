package command

import "github.com/spf13/cobra"

func GetRootCmd() *cobra.Command {
	var rootCmd = cobra.Command{
		Use:   "finances",
		Short: "Finances is a tool to manage your finances",
	}

	rootCmd.AddCommand(getValidateCmd())
	rootCmd.AddCommand(getReportCmd())

	return &rootCmd
}
