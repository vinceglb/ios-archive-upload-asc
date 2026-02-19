package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vinceglb/releasekit-ios/cli/internal/wizard"
)

func newWizardCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "wizard",
		Short: "Run guided setup wizard",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return wizard.Run(cmd.OutOrStdout())
		},
	}
}
