package wizard

import (
	"fmt"
	"io"

	"github.com/vinceglb/releasekit-ios/cli/releasekit-ios-go/internal/term"
)

func Run(out io.Writer) error {
	theme := term.NewTheme()

	fmt.Fprintln(out, theme.Title("ReleaseKit-iOS Wizard"))
	fmt.Fprintln(out, theme.Muted("Guided setup for App Store Connect values and manual GitHub configuration."))
	fmt.Fprintln(out)

	inputs, err := collectInputs()
	if err != nil {
		return err
	}

	if err := validateInputs(inputs); err != nil {
		return err
	}

	printSummary(out, theme, inputs)
	return nil
}
