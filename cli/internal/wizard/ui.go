package wizard

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
)

type wizardFormState struct {
	Workspace     string
	Scheme        string
	BundleID      string
	TeamID        string
	AppID         string
	ASCKeyID      string
	ASCIssuerID   string
	KeySource     string
	P8Path        string
	PrivateKeyB64 string
}

func collectInputs() (Inputs, error) {
	workspaceDefault, _ := detectWorkspaceCandidate(".")

	state := wizardFormState{
		Workspace: workspaceDefault,
		KeySource: "path",
	}

	required := func(label string) func(string) error {
		return func(value string) error {
			if strings.TrimSpace(value) == "" {
				return fmt.Errorf("%s is required", label)
			}
			return nil
		}
	}

	mainForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Xcode workspace path").
				Description("Path to your .xcworkspace (auto-detected when possible)").
				Value(&state.Workspace).
				Validate(required("Xcode workspace path")),
			huh.NewInput().
				Title("Xcode scheme").
				Value(&state.Scheme).
				Validate(required("Xcode scheme")),
			huh.NewInput().
				Title("Bundle ID").
				Value(&state.BundleID).
				Validate(required("Bundle ID")),
			huh.NewInput().
				Title("Apple Team ID").
				Value(&state.TeamID).
				Validate(required("Apple Team ID")),
			huh.NewInput().
				Title("App Store Connect app ID").
				Value(&state.AppID).
				Validate(required("App Store Connect app ID")),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("ASC Key ID").
				Value(&state.ASCKeyID).
				Validate(required("ASC Key ID")),
			huh.NewInput().
				Title("ASC Issuer ID").
				Value(&state.ASCIssuerID).
				Validate(required("ASC Issuer ID")),
			huh.NewSelect[string]().
				Title("How do you want to provide your private key?").
				Options(
					huh.NewOption("Local .p8 file path (recommended)", "path"),
					huh.NewOption("Paste base64 private key", "base64"),
				).
				Value(&state.KeySource),
		),
	)

	if err := mainForm.Run(); err != nil {
		return Inputs{}, fmt.Errorf("wizard canceled: %w", err)
	}

	switch state.KeySource {
	case "path":
		pathForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Path to .p8 private key file").
					Value(&state.P8Path).
					Validate(required("Path to .p8 private key file")),
			),
		)
		if err := pathForm.Run(); err != nil {
			return Inputs{}, fmt.Errorf("wizard canceled: %w", err)
		}
		if _, err := os.Stat(state.P8Path); err != nil {
			return Inputs{}, fmt.Errorf("p8 file not found: %s", state.P8Path)
		}
		encoded, err := encodeFileBase64(state.P8Path)
		if err != nil {
			return Inputs{}, fmt.Errorf("could not read .p8 file: %w", err)
		}
		state.PrivateKeyB64 = encoded
	case "base64":
		b64Form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("ASC private key base64").
					Value(&state.PrivateKeyB64).
					Validate(required("ASC private key base64")),
			),
		)
		if err := b64Form.Run(); err != nil {
			return Inputs{}, fmt.Errorf("wizard canceled: %w", err)
		}
		state.PrivateKeyB64 = normalizeBase64(state.PrivateKeyB64)
	default:
		return Inputs{}, fmt.Errorf("unknown key source: %s", state.KeySource)
	}

	return Inputs{
		Workspace:        strings.TrimSpace(state.Workspace),
		Scheme:           strings.TrimSpace(state.Scheme),
		BundleID:         strings.TrimSpace(state.BundleID),
		TeamID:           strings.TrimSpace(state.TeamID),
		AppID:            strings.TrimSpace(state.AppID),
		ASCKeyID:         strings.TrimSpace(state.ASCKeyID),
		ASCIssuerID:      strings.TrimSpace(state.ASCIssuerID),
		ASCPrivateKeyB64: state.PrivateKeyB64,
	}, nil
}
