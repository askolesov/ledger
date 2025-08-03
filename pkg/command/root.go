package command

import "github.com/spf13/cobra"

func GetRootCmd() *cobra.Command {
	var rootCmd = cobra.Command{
		Use:   "ledger",
		Short: "A command-line tool for managing financial ledgers using the Open Ledger Format",
		Long: `Ledger is a command-line tool for managing financial ledgers using the Open Ledger Format (OLF).

The tool supports both v1 and v2 formats of the Open Ledger Format:
- Root commands (validate, report) work with OLF v2.0 format
- v1 subcommands work with the legacy OLF v1.0 format

Supported file formats: YAML (.yaml, .yml), JSON (.json)

Examples:
  ledger validate ledger.yaml          # Validate OLF v2.0 file
  ledger report ledger.yaml            # Generate report from OLF v2.0 file
  ledger v1 validate data.yaml         # Validate OLF v1.0 file
  ledger v1 report data.yaml           # Generate report from OLF v1.0 file`,
	}

	// Add v2 commands as root commands
	rootCmd.AddCommand(getV2ValidateCmd())
	rootCmd.AddCommand(getV2ReportCmd())

	// Add v1 subcommand
	rootCmd.AddCommand(getV1Cmd())

	return &rootCmd
}
