package cmd

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "releasekit-ios",
		Short:         "ReleaseKit-iOS CLI",
		Long:          "ReleaseKit-iOS CLI for onboarding and setup workflows.",
		SilenceUsage:  true,
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	rootCmd.AddCommand(newWizardCmd())

	return rootCmd
}

func Execute() error {
	return NewRootCmd().Execute()
}
