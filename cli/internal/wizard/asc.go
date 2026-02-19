package wizard

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ASCApp represents a single app from the App Store Connect API.
type ASCApp struct {
	ID         string `json:"id"`
	Attributes struct {
		Name     string `json:"name"`
		BundleID string `json:"bundleId"`
	} `json:"attributes"`
}

// ListASCApps authenticates with the asc CLI and returns all apps from App Store Connect.
func ListASCApps(keyID, issuerID, privKeyB64 string) ([]ASCApp, error) {
	tmpDir, err := os.MkdirTemp("", "releasekit-wizard-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	// Decode base64 private key to temp .p8 file.
	keyBytes, err := base64.StdEncoding.DecodeString(normalizeBase64(privKeyB64))
	if err != nil {
		return nil, fmt.Errorf("invalid private key base64: %w", err)
	}
	privKeyPath := filepath.Join(tmpDir, "AuthKey.p8")
	if err := os.WriteFile(privKeyPath, keyBytes, 0600); err != nil {
		return nil, err
	}

	ascHome := filepath.Join(tmpDir, "asc-home")
	if err := os.MkdirAll(ascHome, 0755); err != nil {
		return nil, err
	}

	// Login to asc.
	loginCmd := exec.Command("asc", "auth", "login",
		"--bypass-keychain",
		"--skip-validation",
		"--name", "releasekit-ios-wizard",
		"--key-id", keyID,
		"--issuer-id", issuerID,
		"--private-key", privKeyPath,
	)
	loginCmd.Env = buildASCEnv(ascHome, keyID, issuerID, privKeyPath)
	if out, loginErr := loginCmd.CombinedOutput(); loginErr != nil {
		return nil, fmt.Errorf("asc auth login failed: %w\n%s", loginErr, strings.TrimSpace(string(out)))
	}

	// List apps.
	listCmd := exec.Command("asc", "apps", "list", "--output", "json")
	listCmd.Env = buildASCEnv(ascHome, keyID, issuerID, privKeyPath)
	out, err := listCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("asc apps list failed: %w", err)
	}

	return parseASCAppsOutput(out)
}

// parseASCAppsOutput parses both bare-array and {"data":[]} JSON envelope formats.
func parseASCAppsOutput(raw []byte) ([]ASCApp, error) {
	raw = []byte(strings.TrimSpace(string(raw)))
	if len(raw) == 0 {
		return nil, nil
	}

	// Try bare array first.
	if raw[0] == '[' {
		var apps []ASCApp
		if err := json.Unmarshal(raw, &apps); err != nil {
			return nil, fmt.Errorf("failed to parse asc output: %w", err)
		}
		return apps, nil
	}

	// Try {"data":[...]} envelope.
	var envelope struct {
		Data []ASCApp `json:"data"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, fmt.Errorf("failed to parse asc output: %w", err)
	}
	return envelope.Data, nil
}

// buildASCEnv builds the environment for running asc commands, stripping any
// existing conflicting variables and injecting the provided credentials.
func buildASCEnv(ascHome, keyID, issuerID, privKeyPath string) []string {
	skipPrefixes := []string{
		"HOME=",
		"ASC_KEY_ID=",
		"ASC_ISSUER_ID=",
		"ASC_PRIVATE_KEY_PATH=",
		"ASC_BYPASS_KEYCHAIN=",
	}

	var filtered []string
	for _, env := range os.Environ() {
		skip := false
		for _, prefix := range skipPrefixes {
			if strings.HasPrefix(env, prefix) {
				skip = true
				break
			}
		}
		if !skip {
			filtered = append(filtered, env)
		}
	}

	filtered = append(filtered,
		"HOME="+ascHome,
		"ASC_BYPASS_KEYCHAIN=1",
		"ASC_KEY_ID="+keyID,
		"ASC_ISSUER_ID="+issuerID,
		"ASC_PRIVATE_KEY_PATH="+privKeyPath,
	)
	return filtered
}
