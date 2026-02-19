package wizard

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/vinceglb/releasekit-ios/cli/internal/term"
)

// checkPrerequisites checks for required and optional tools.
// Returns an error (and prints a hard-exit message) only when `asc` is missing.
// xcodeOK, ghOK, ghAuthed are non-blocking signals for graceful fallback.
func checkPrerequisites(out io.Writer, theme term.Theme) (xcodeOK, ghOK, ghAuthed bool, err error) {
	fmt.Fprintln(out, theme.Section("Prerequisites"))

	ascOK := commandExists("asc")
	xcodeOK = commandExists("xcodebuild")
	ghOK = commandExists("gh")

	printCheck(out, theme, "asc", ascOK)
	printCheck(out, theme, "xcodebuild", xcodeOK)
	printCheck(out, theme, "gh", ghOK)

	if ghOK {
		ghAuthed = ghIsAuthenticated()
		if ghAuthed {
			fmt.Fprintf(out, "  %s gh authenticated\n", theme.Success("✓"))
		} else {
			fmt.Fprintf(out, "  %s gh not authenticated (run: gh auth login)\n", theme.Muted("○"))
		}
	}

	fmt.Fprintln(out)

	if !ascOK {
		fmt.Fprintln(out, theme.Error("asc CLI not found. Install it with:"))
		fmt.Fprintln(out, theme.Value("  brew install rudrankriyam/tap/asc"))
		fmt.Fprintln(out)
		return false, false, false, fmt.Errorf("asc CLI is required")
	}

	return xcodeOK, ghOK, ghAuthed, nil
}

func printCheck(out io.Writer, theme term.Theme, name string, ok bool) {
	if ok {
		fmt.Fprintf(out, "  %s %s\n", theme.Success("✓"), name)
	} else {
		fmt.Fprintf(out, "  %s %s\n", theme.Muted("✗"), name)
	}
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func ghIsAuthenticated() bool {
	cmd := exec.Command("gh", "auth", "status")
	return cmd.Run() == nil
}
