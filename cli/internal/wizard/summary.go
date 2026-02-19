package wizard

import (
	"fmt"
	"io"

	"github.com/vinceglb/releasekit-ios/cli/internal/term"
)

func printSummary(out io.Writer, theme term.Theme, inputs Inputs) {
	fmt.Fprintln(out)
	fmt.Fprintln(out, theme.Section("Wizard Completed"))
	fmt.Fprintln(out, theme.Muted("Setup values collected successfully."))
	fmt.Fprintln(out)

	fmt.Fprintln(out, theme.Section("Summary"))
	printKV(out, theme, "Workspace", inputs.Workspace)
	printKV(out, theme, "Scheme", inputs.Scheme)
	printKV(out, theme, "Bundle ID", inputs.BundleID)
	printKV(out, theme, "Team ID", inputs.TeamID)
	printKV(out, theme, "App ID", inputs.AppID)
	fmt.Fprintln(out)

	fmt.Fprintln(out, theme.Section("GitHub Secrets"))
	fmt.Fprintf(out, "  %s\n", theme.Value("ASC_KEY_ID="+inputs.ASCKeyID))
	fmt.Fprintf(out, "  %s\n", theme.Value("ASC_ISSUER_ID="+inputs.ASCIssuerID))
	fmt.Fprintf(out, "  %s\n", theme.Value("ASC_PRIVATE_KEY_B64="+inputs.ASCPrivateKeyB64))
	fmt.Fprintln(out)

	fmt.Fprintln(out, theme.Section("GitHub Variables"))
	fmt.Fprintf(out, "  %s\n", theme.Value("ASC_APP_ID="+inputs.AppID))
	fmt.Fprintf(out, "  %s\n", theme.Value("ASC_TEAM_ID="+inputs.TeamID))
	fmt.Fprintf(out, "  %s\n", theme.Value("BUNDLE_ID="+inputs.BundleID))
	fmt.Fprintln(out)

	fmt.Fprintln(out, theme.Section("Next"))
	fmt.Fprintln(out, theme.Muted("1) Add the values above to your GitHub repository"))
	fmt.Fprintln(out, theme.Muted("2) Run your iOS release workflow"))
}

func printKV(out io.Writer, theme term.Theme, label, value string) {
	fmt.Fprintf(out, "  %s %s\n", theme.Label(fmt.Sprintf("%-12s", label+":")), theme.Value(value))
}
