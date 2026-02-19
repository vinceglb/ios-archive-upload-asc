package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vinceglb/releasekit-ios/cli/internal/version"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print CLI version information",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "releasekit-ios %s\n", version.Version)
			fmt.Fprintf(cmd.OutOrStdout(), "commit: %s\n", version.Commit)
			fmt.Fprintf(cmd.OutOrStdout(), "build_date: %s\n", version.BuildDate)
		},
	}
}
