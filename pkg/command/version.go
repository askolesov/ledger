package command

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func getVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  `Print version information including version, commit hash, build date, and build environment.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("ledger version %s\n", version)
			fmt.Printf("  commit: %s\n", commit)
			fmt.Printf("  built at: %s\n", date)
			fmt.Printf("  built by: %s\n", builtBy)
			fmt.Printf("  go version: %s\n", runtime.Version())
			fmt.Printf("  platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		},
	}
}
