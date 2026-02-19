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

	fmt.Fprintln(out, theme.Section("Configuration"))
	if inputs.AppName != "" {
		printKV(out, theme, "App Name", inputs.AppName)
	}
	printKV(out, theme, "Workspace", inputs.Workspace)
	printKV(out, theme, "Scheme", inputs.Scheme)
	printKV(out, theme, "Bundle ID", inputs.BundleID)
	printKV(out, theme, "Team ID", inputs.TeamID)
	printKV(out, theme, "App ID", inputs.AppID)
	fmt.Fprintln(out)

	secretNote := ""
	if inputs.SecretsWereSet {
		secretNote = " " + theme.Success("(✓ set automatically)")
	}

	fmt.Fprintln(out, theme.Section("GitHub Secrets"+secretNote))
	fmt.Fprintf(out, "  %s\n", theme.Value("ASC_KEY_ID="+inputs.ASCKeyID))
	fmt.Fprintf(out, "  %s\n", theme.Value("ASC_ISSUER_ID="+inputs.ASCIssuerID))
	fmt.Fprintf(out, "  %s\n", theme.Value("ASC_PRIVATE_KEY_B64="+inputs.ASCPrivateKeyB64))
	fmt.Fprintln(out)

	fmt.Fprintln(out, theme.Section("GitHub Variables"))
	fmt.Fprintf(out, "  %s\n", theme.Value("ASC_APP_ID="+inputs.AppID))
	fmt.Fprintf(out, "  %s\n", theme.Value("ASC_TEAM_ID="+inputs.TeamID))
	fmt.Fprintf(out, "  %s\n", theme.Value("BUNDLE_ID="+inputs.BundleID))
	fmt.Fprintln(out)

	fmt.Fprintln(out, theme.Section("Next Steps"))
	stepNum := 1

	if !inputs.SecretsWereSet {
		fmt.Fprintf(out, theme.Muted("  %d) Add the GitHub Secrets above to your repository\n"), stepNum)
		stepNum++
		fmt.Fprintf(out, theme.Muted("  %d) Add the GitHub Variables above to your repository\n"), stepNum)
		stepNum++
	} else {
		fmt.Fprintf(out, theme.Muted("  %d) Add the GitHub Variables above to your repository\n"), stepNum)
		stepNum++
	}

	if inputs.WorkflowWasWritten && inputs.WorkflowPath != "" {
		fmt.Fprintf(out, theme.Muted("  %d) Commit and push %s\n"), stepNum, inputs.WorkflowPath)
		stepNum++
	} else if !inputs.WorkflowWasWritten {
		fmt.Fprintf(out, theme.Muted("  %d) Add a release workflow (see docs for a template)\n"), stepNum)
		stepNum++
	}

	fmt.Fprintf(out, theme.Muted("  %d) Push a v* tag to trigger your release\n"), stepNum)
}

func printKV(out io.Writer, theme term.Theme, label, value string) {
	fmt.Fprintf(out, "  %s %s\n", theme.Label(fmt.Sprintf("%-12s", label+":")), theme.Value(value))
}
