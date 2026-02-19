package main

import (
	"fmt"
	"os"

	"github.com/vinceglb/releasekit-ios/cli/cmd"
	"github.com/vinceglb/releasekit-ios/cli/internal/term"
)

func main() {
	if err := cmd.Execute(); err != nil {
		theme := term.NewTheme()
		fmt.Fprintln(os.Stderr, theme.Error(err.Error()))
		os.Exit(1)
	}
}
