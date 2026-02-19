package wizard

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

func TestValidateInputsSuccess(t *testing.T) {
	tmpDir := t.TempDir()
	workspace := filepath.Join(tmpDir, "App.xcworkspace")
	if err := os.Mkdir(workspace, 0o755); err != nil {
		t.Fatalf("mkdir workspace: %v", err)
	}

	inputs := Inputs{
		Workspace:        workspace,
		Scheme:           "App",
		BundleID:         "com.example.app",
		TeamID:           "TEAMID123",
		AppID:            "123456789",
		ASCKeyID:         "KEYID123",
		ASCIssuerID:      "00000000-0000-0000-0000-000000000000",
		ASCPrivateKeyB64: base64.StdEncoding.EncodeToString([]byte("private-key")),
	}

	if err := validateInputs(inputs); err != nil {
		t.Fatalf("expected no validation error, got: %v", err)
	}
}

func TestValidateInputsMissingWorkspace(t *testing.T) {
	inputs := Inputs{
		Workspace:        "/tmp/does-not-exist.xcworkspace",
		Scheme:           "App",
		BundleID:         "com.example.app",
		TeamID:           "TEAMID123",
		AppID:            "123456789",
		ASCKeyID:         "KEYID123",
		ASCIssuerID:      "00000000-0000-0000-0000-000000000000",
		ASCPrivateKeyB64: "cHJpdmF0ZS1rZXk=",
	}

	err := validateInputs(inputs)
	if err == nil {
		t.Fatal("expected workspace error")
	}
}

func TestValidateInputsInvalidBase64(t *testing.T) {
	tmpDir := t.TempDir()
	workspace := filepath.Join(tmpDir, "App.xcworkspace")
	if err := os.Mkdir(workspace, 0o755); err != nil {
		t.Fatalf("mkdir workspace: %v", err)
	}

	inputs := Inputs{
		Workspace:        workspace,
		Scheme:           "App",
		BundleID:         "com.example.app",
		TeamID:           "TEAMID123",
		AppID:            "123456789",
		ASCKeyID:         "KEYID123",
		ASCIssuerID:      "00000000-0000-0000-0000-000000000000",
		ASCPrivateKeyB64: "not-base64@@",
	}

	err := validateInputs(inputs)
	if err == nil {
		t.Fatal("expected invalid base64 error")
	}
}

func TestNormalizeBase64RemovesWhitespace(t *testing.T) {
	input := "YWJj\n ZGVm\t"
	got := normalizeBase64(input)
	want := "YWJjZGVm"
	if got != want {
		t.Fatalf("normalizeBase64() = %q, want %q", got, want)
	}
}

func TestEncodeFileBase64(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "AuthKey_TEST.p8")
	if err := os.WriteFile(path, []byte("private-key"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	got, err := encodeFileBase64(path)
	if err != nil {
		t.Fatalf("encodeFileBase64 returned error: %v", err)
	}
	want := base64.StdEncoding.EncodeToString([]byte("private-key"))
	if got != want {
		t.Fatalf("encodeFileBase64() = %q, want %q", got, want)
	}
}
