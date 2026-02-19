package wizard

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/vinceglb/releasekit-ios/cli/internal/term"
)

// collectPhase1ASCCredentials collects App Store Connect API credentials.
// Returns issuerID, keyID, privKeyB64 (base64-encoded), or an error.
func collectPhase1ASCCredentials() (issuerID, keyID, privKeyB64 string, err error) {
	var keySource string = "path"
	var p8Path string

	form := huh.NewForm(
		// Group 1: Instruction note.
		huh.NewGroup(
			huh.NewNote().
				Title("App Store Connect API Key").
				Description(
					"1. Go to appstoreconnect.apple.com → Users → Integrations → API Keys\n" +
						"2. Create a key with App Manager role\n" +
						"3. Download the .p8 file — only downloadable once\n" +
						"4. Note the Key ID and Issuer ID shown on that page",
				),
		),
		// Group 2: Credentials.
		huh.NewGroup(
			huh.NewInput().
				Title("Issuer ID").
				Placeholder("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx").
				Value(&issuerID).
				Validate(requiredField("Issuer ID")),
			huh.NewInput().
				Title("Key ID").
				Placeholder("ABC123XY45").
				Value(&keyID).
				Validate(requiredField("Key ID")),
			huh.NewSelect[string]().
				Title("Private key source").
				Options(
					huh.NewOption("Local .p8 file (recommended)", "path"),
					huh.NewOption("Paste base64", "base64"),
				).
				Value(&keySource),
		),
		// Group 3: .p8 file path — hidden when keySource != "path".
		huh.NewGroup(
			huh.NewInput().
				Title("Path to .p8 file").
				Placeholder("/path/to/AuthKey_ABC123XY45.p8").
				Value(&p8Path).
				Validate(requiredField("Path to .p8 file")),
		).WithHideFunc(func() bool { return keySource != "path" }),
		// Group 4: Base64 input — hidden when keySource != "base64".
		huh.NewGroup(
			huh.NewInput().
				Title("Base64 private key").
				Placeholder("LS0tLS1CRUdJTi...").
				Value(&privKeyB64).
				Validate(requiredField("Base64 private key")),
		).WithHideFunc(func() bool { return keySource != "base64" }),
	).WithTheme(huh.ThemeCharm())

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", "", "", fmt.Errorf("wizard canceled")
		}
		return "", "", "", err
	}

	switch keySource {
	case "path":
		if _, statErr := os.Stat(p8Path); statErr != nil {
			return "", "", "", fmt.Errorf("p8 file not found: %s", p8Path)
		}
		encoded, encErr := encodeFileBase64(p8Path)
		if encErr != nil {
			return "", "", "", fmt.Errorf("could not read .p8 file: %w", encErr)
		}
		privKeyB64 = encoded
	case "base64":
		privKeyB64 = normalizeBase64(privKeyB64)
	}

	return strings.TrimSpace(issuerID), strings.TrimSpace(keyID), privKeyB64, nil
}

// collectPhase2App shows app selection from a pre-fetched list, or falls back
// to manual input when no apps are available.
func collectPhase2App(apps []ASCApp) (appName, appID, bundleID string, err error) {
	if len(apps) > 0 {
		var selectedAppID string
		appMap := make(map[string]ASCApp, len(apps))
		options := make([]huh.Option[string], 0, len(apps))
		for _, app := range apps {
			label := fmt.Sprintf("%s (%s)", app.Attributes.Name, app.Attributes.BundleID)
			options = append(options, huh.NewOption(label, app.ID))
			appMap[app.ID] = app
		}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select your app").
					Description("Fetched from App Store Connect").
					Options(options...).
					Value(&selectedAppID),
			),
		).WithTheme(huh.ThemeCharm())

		if err := form.Run(); err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				return "", "", "", fmt.Errorf("wizard canceled")
			}
			return "", "", "", err
		}

		selected := appMap[selectedAppID]
		return selected.Attributes.Name, selected.ID, selected.Attributes.BundleID, nil
	}

	// Fallback: manual input.
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("App Store Connect App ID").
				Description("Numeric app ID from App Store Connect").
				Value(&appID).
				Validate(requiredField("App ID")),
			huh.NewInput().
				Title("Bundle ID").
				Placeholder("com.example.myapp").
				Value(&bundleID).
				Validate(requiredField("Bundle ID")),
		),
	).WithTheme(huh.ThemeCharm())

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", "", "", fmt.Errorf("wizard canceled")
		}
		return "", "", "", err
	}

	return "", strings.TrimSpace(appID), strings.TrimSpace(bundleID), nil
}

// collectPhase3Xcode collects Xcode project settings using detected candidates,
// schemes, and team ID. prefillBundleID is pre-filled from Phase 2.
func collectPhase3Xcode(
	candidates []string,
	schemes []string,
	detectedTeamID string,
	prefillBundleID string,
) (workspace, scheme, teamID, bundleID string, err error) {
	var useDetectedTeam bool = true
	var teamInput string
	bundleID = prefillBundleID

	var groups []*huh.Group

	// Workspace group.
	switch len(candidates) {
	case 0:
		groups = append(groups, huh.NewGroup(
			huh.NewInput().
				Title("Xcode workspace path").
				Description("Path to your .xcworkspace").
				Value(&workspace).
				Validate(requiredField("Xcode workspace path")),
		))
	case 1:
		workspace = candidates[0]
		groups = append(groups, huh.NewGroup(
			huh.NewNote().
				Title("Xcode Workspace").
				Description("Detected: "+candidates[0]),
		))
	default:
		opts := make([]huh.Option[string], len(candidates))
		for i, c := range candidates {
			opts[i] = huh.NewOption(c, c)
		}
		groups = append(groups, huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Xcode workspace").
				Options(opts...).
				Value(&workspace),
		))
	}

	// Scheme group.
	if len(schemes) == 0 {
		groups = append(groups, huh.NewGroup(
			huh.NewInput().
				Title("Xcode scheme").
				Value(&scheme).
				Validate(requiredField("Xcode scheme")),
		))
	} else {
		opts := make([]huh.Option[string], len(schemes))
		for i, s := range schemes {
			opts[i] = huh.NewOption(s, s)
		}
		groups = append(groups, huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Xcode scheme").
				Options(opts...).
				Value(&scheme),
		))
	}

	// Team ID group.
	if detectedTeamID != "" {
		groups = append(groups,
			huh.NewGroup(
				huh.NewConfirm().
					Title(fmt.Sprintf("Use detected Apple Team ID: %s?", detectedTeamID)).
					Affirmative("Yes").
					Negative("Enter manually").
					Inline(true).
					Value(&useDetectedTeam),
			),
			huh.NewGroup(
				huh.NewInput().
					Title("Apple Team ID").
					Placeholder("XXXXXXXXXX").
					Value(&teamInput).
					Validate(requiredField("Apple Team ID")),
			).WithHideFunc(func() bool { return useDetectedTeam }),
		)
	} else {
		groups = append(groups, huh.NewGroup(
			huh.NewInput().
				Title("Apple Team ID").
				Placeholder("XXXXXXXXXX").
				Value(&teamInput).
				Validate(requiredField("Apple Team ID")),
		))
	}

	// Bundle ID group (pre-filled, editable).
	groups = append(groups, huh.NewGroup(
		huh.NewInput().
			Title("Bundle ID").
			Placeholder("com.example.myapp").
			Value(&bundleID).
			Validate(requiredField("Bundle ID")),
	))

	form := huh.NewForm(groups...).WithTheme(huh.ThemeCharm())
	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", "", "", "", fmt.Errorf("wizard canceled")
		}
		return "", "", "", "", err
	}

	// Resolve team ID.
	if detectedTeamID != "" && useDetectedTeam {
		teamID = detectedTeamID
	} else {
		teamID = strings.TrimSpace(teamInput)
	}

	return strings.TrimSpace(workspace), strings.TrimSpace(scheme), teamID, strings.TrimSpace(bundleID), nil
}

// collectPhase4GitHub handles GitHub secret setting and optional workflow generation.
func collectPhase4GitHub(out io.Writer, theme term.Theme, inputs *Inputs, ghAuthed bool, owner, repo string) error {
	repoSlug := ""
	if owner != "" && repo != "" {
		repoSlug = owner + "/" + repo
		inputs.GitHubRepo = repoSlug
	}

	// Ask to auto-set secrets.
	var wantAutoSecrets bool
	if ghAuthed && repoSlug != "" {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title(fmt.Sprintf("Set GitHub secrets automatically for %s?", repoSlug)).
					Affirmative("Yes, use gh CLI").
					Negative("I'll do it manually").
					Inline(true).
					Value(&wantAutoSecrets),
			),
		).WithTheme(huh.ThemeCharm())
		if err := form.Run(); err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				return fmt.Errorf("wizard canceled")
			}
			return err
		}
	}

	if wantAutoSecrets {
		secrets := map[string]string{
			"ASC_KEY_ID":          inputs.ASCKeyID,
			"ASC_ISSUER_ID":       inputs.ASCIssuerID,
			"ASC_PRIVATE_KEY_B64": inputs.ASCPrivateKeyB64,
		}
		errs := SetGitHubSecrets(repoSlug, secrets)

		allOK := true
		for name, setErr := range errs {
			if setErr != nil {
				fmt.Fprintf(out, "  %s Failed to set %s: %v\n", theme.Error("✗"), name, setErr)
				allOK = false
			} else {
				fmt.Fprintf(out, "  %s %s set\n", theme.Success("✓"), name)
			}
		}
		if allOK {
			inputs.SecretsWereSet = true
		}
		fmt.Fprintln(out)
	}

	// Ask to generate workflow file.
	var wantWorkflow bool
	workflowForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Generate .github/workflows/release.yml?").
				Affirmative("Yes").
				Negative("No, I'll write it myself").
				Inline(true).
				Value(&wantWorkflow),
		),
	).WithTheme(huh.ThemeCharm())
	if err := workflowForm.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return fmt.Errorf("wizard canceled")
		}
		return err
	}

	if wantWorkflow {
		workflowPath := DefaultWorkflowPath()
		inputs.WorkflowPath = workflowPath

		// Check if file already exists.
		if _, statErr := os.Stat(workflowPath); statErr == nil {
			var overwrite bool
			overwriteForm := huh.NewForm(
				huh.NewGroup(
					huh.NewConfirm().
						Title(workflowPath+" already exists. Overwrite?").
						Affirmative("Yes, overwrite").
						Negative("No, keep existing").
						Inline(true).
						Value(&overwrite),
				),
			).WithTheme(huh.ThemeCharm())
			if err := overwriteForm.Run(); err != nil {
				if errors.Is(err, huh.ErrUserAborted) {
					return fmt.Errorf("wizard canceled")
				}
				return err
			}
			if !overwrite {
				return nil
			}
		}

		content, err := GenerateWorkflow(*inputs)
		if err != nil {
			return fmt.Errorf("failed to generate workflow: %w", err)
		}
		if err := WriteWorkflow(workflowPath, content); err != nil {
			return fmt.Errorf("failed to write workflow: %w", err)
		}
		inputs.WorkflowWasWritten = true
		fmt.Fprintf(out, "  %s %s written\n", theme.Success("✓"), workflowPath)
		fmt.Fprintln(out)
	}

	return nil
}

// requiredField returns a validation function that rejects blank values.
func requiredField(label string) func(string) error {
	return func(value string) error {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s is required", label)
		}
		return nil
	}
}
