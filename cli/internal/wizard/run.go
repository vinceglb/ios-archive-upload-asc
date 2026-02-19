package wizard

import (
	"errors"
	"fmt"
	"io"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/vinceglb/releasekit-ios/cli/internal/term"
)

func Run(out io.Writer) error {
	theme := term.NewTheme()

	fmt.Fprintln(out, theme.Title("ReleaseKit-iOS Wizard"))
	fmt.Fprintln(out, theme.Muted("Guided setup for distributing iOS apps to the App Store."))
	fmt.Fprintln(out)

	// Phase 0: Prerequisites.
	xcodeOK, ghOK, ghAuthed, err := checkPrerequisites(out, theme)
	if err != nil {
		return err
	}
	_ = xcodeOK
	_ = ghOK

	// Phase 1: ASC credentials.
	fmt.Fprintln(out, theme.Section("Phase 1 — App Store Connect Credentials"))
	fmt.Fprintln(out)

	issuerID, keyID, privKeyB64, err := collectPhase1ASCCredentials()
	if err != nil {
		return err
	}

	// Validate credentials with spinner (also fetches app list).
	var apps []ASCApp
	var ascErr error
	if spinErr := spinner.New().
		Title("Validating App Store Connect credentials…").
		Action(func() {
			apps, ascErr = ListASCApps(keyID, issuerID, privKeyB64)
		}).
		Run(); spinErr != nil && !errors.Is(spinErr, huh.ErrUserAborted) {
		// Spinner itself failed — continue; ascErr will be nil so we surface the spinner error.
		ascErr = spinErr
	}

	if ascErr != nil {
		fmt.Fprintf(out, "%s Failed to validate credentials: %v\n\n", theme.Error("✗"), ascErr)
		return ascErr
	}
	fmt.Fprintf(out, "%s App Store Connect credentials validated (%d apps found)\n\n",
		theme.Success("✓"), len(apps))

	// Phase 2: App selection.
	fmt.Fprintln(out, theme.Section("Phase 2 — App Selection"))
	fmt.Fprintln(out)

	appName, appID, bundleID, err := collectPhase2App(apps)
	if err != nil {
		return err
	}

	// Phase 3: Xcode project.
	fmt.Fprintln(out, theme.Section("Phase 3 — Xcode Project"))
	fmt.Fprintln(out)

	var candidates []string
	var schemes []string
	var detectedTeamID string

	if spinErr := spinner.New().
		Title("Detecting Xcode project…").
		Action(func() {
			candidates = detectAllWorkspaceCandidates(".")
			if len(candidates) == 1 {
				schemes, _ = DetectSchemes(candidates[0])
				detectedTeamID, _ = DetectTeamID(candidates[0])
			}
		}).
		Run(); spinErr != nil && !errors.Is(spinErr, huh.ErrUserAborted) {
		// Non-fatal: proceed with empty candidates.
	}

	workspace, scheme, teamID, bundleID, err := collectPhase3Xcode(candidates, schemes, detectedTeamID, bundleID)
	if err != nil {
		return err
	}

	inputs := Inputs{
		AppName:          appName,
		AppID:            appID,
		BundleID:         bundleID,
		Workspace:        workspace,
		Scheme:           scheme,
		TeamID:           teamID,
		ASCKeyID:         keyID,
		ASCIssuerID:      issuerID,
		ASCPrivateKeyB64: privKeyB64,
	}

	if err := validateInputs(inputs); err != nil {
		return err
	}

	// Phase 4: GitHub setup.
	fmt.Fprintln(out, theme.Section("Phase 4 — GitHub Setup"))
	fmt.Fprintln(out)

	var ghOwner, ghRepo string
	if ghAuthed {
		if spinErr := spinner.New().
			Title("Detecting git repository…").
			Action(func() {
				ghOwner, ghRepo, _ = DetectGitRepo()
			}).
			Run(); spinErr != nil && !errors.Is(spinErr, huh.ErrUserAborted) {
			// Non-fatal.
		}
	}

	if err := collectPhase4GitHub(out, theme, &inputs, ghAuthed, ghOwner, ghRepo); err != nil {
		return err
	}

	// Phase 5: Done.
	printSummary(out, theme, inputs)
	return nil
}
