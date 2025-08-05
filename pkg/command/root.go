package command

import "github.com/spf13/cobra"

func GetRootCmd() *cobra.Command {
	var rootCmd = cobra.Command{
		Use:   "ledger",
		Short: "Command-line tool for managing financial ledgers using Open Ledger Format",
		Long: `Command-line tool for managing financial ledgers using Open Ledger Format (OLF).

Supports both OLF v1.0 (legacy) and v2.0 formats:
- Root commands work with OLF v2.0 (YAML/JSON)
- v1 subcommands work with OLF v1.0 (YAML only)

Examples:
  ledger validate ledger.yaml          # Validate OLF v2.0 file
  ledger report ledger.yaml            # Generate OLF v2.0 report
  ledger v1 validate data.yaml         # Validate OLF v1.0 file
  ledger v1 report data.yaml           # Generate OLF v1.0 report`,
		Version: version,
	}

	// Add v2 commands as root commands
	rootCmd.AddCommand(getV2ValidateCmd())
	rootCmd.AddCommand(getV2ReportCmd())

	// Add version command
	rootCmd.AddCommand(getVersionCmd())

	// Add v1 subcommand
	rootCmd.AddCommand(getV1Cmd())

	return &rootCmd
}
