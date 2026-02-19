package wizard

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

func normalizeBase64(value string) string {
	if value == "" {
		return ""
	}
	parts := strings.Fields(value)
	return strings.Join(parts, "")
}

func validateInputs(inputs Inputs) error {
	if strings.TrimSpace(inputs.Workspace) == "" {
		return fmt.Errorf("Xcode workspace path is required")
	}
	if strings.TrimSpace(inputs.Scheme) == "" {
		return fmt.Errorf("Xcode scheme is required")
	}
	if strings.TrimSpace(inputs.BundleID) == "" {
		return fmt.Errorf("Bundle ID is required")
	}
	if strings.TrimSpace(inputs.TeamID) == "" {
		return fmt.Errorf("Apple Team ID is required")
	}
	if strings.TrimSpace(inputs.AppID) == "" {
		return fmt.Errorf("App Store Connect app ID is required")
	}
	if strings.TrimSpace(inputs.ASCKeyID) == "" {
		return fmt.Errorf("ASC Key ID is required")
	}
	if strings.TrimSpace(inputs.ASCIssuerID) == "" {
		return fmt.Errorf("ASC Issuer ID is required")
	}
	if strings.TrimSpace(inputs.ASCPrivateKeyB64) == "" {
		return fmt.Errorf("ASC private key base64 is required")
	}

	if _, err := os.Stat(inputs.Workspace); err != nil {
		return fmt.Errorf("workspace path does not exist: %s", inputs.Workspace)
	}

	normalized := normalizeBase64(inputs.ASCPrivateKeyB64)
	if _, err := base64.StdEncoding.DecodeString(normalized); err != nil {
		return fmt.Errorf("ASC private key is not valid base64")
	}

	return nil
}

func encodeFileBase64(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(content)
	return encoded, nil
}
