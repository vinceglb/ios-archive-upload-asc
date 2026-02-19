package wizard

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectTeamIDFromPbxproj(t *testing.T) {
	// Create a minimal project.pbxproj fixture.
	tmpDir := t.TempDir()
	xcodeprojDir := filepath.Join(tmpDir, "MyApp.xcodeproj")
	if err := os.MkdirAll(xcodeprojDir, 0755); err != nil {
		t.Fatal(err)
	}

	pbxprojContent := `// !$*UTF8*$!
{
	archiveVersion = 1;
	classes = {};
	objectVersion = 56;
	objects = {
		ABC = {
			DEVELOPMENT_TEAM = ABCDE12345;
			PRODUCT_BUNDLE_IDENTIFIER = "com.example.app";
		};
		DEF = {
			DEVELOPMENT_TEAM = ABCDE12345;
			PRODUCT_BUNDLE_IDENTIFIER = "com.example.app.tests";
		};
	};
}
`
	pbxprojPath := filepath.Join(xcodeprojDir, "project.pbxproj")
	if err := os.WriteFile(pbxprojPath, []byte(pbxprojContent), 0644); err != nil {
		t.Fatal(err)
	}

	// DetectTeamID walks from the workspace directory, which is tmpDir.
	// We pass a fake workspace path inside tmpDir.
	fakeWorkspace := filepath.Join(tmpDir, "MyApp.xcworkspace")
	teamID, err := DetectTeamID(fakeWorkspace)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if teamID != "ABCDE12345" {
		t.Errorf("expected team ID 'ABCDE12345', got %q", teamID)
	}
}

func TestDetectTeamIDNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	fakeWorkspace := filepath.Join(tmpDir, "MyApp.xcworkspace")

	teamID, err := DetectTeamID(fakeWorkspace)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if teamID != "" {
		t.Errorf("expected empty team ID, got %q", teamID)
	}
}

func TestDetectTeamIDMostFrequent(t *testing.T) {
	tmpDir := t.TempDir()
	xcodeprojDir := filepath.Join(tmpDir, "MyApp.xcodeproj")
	if err := os.MkdirAll(xcodeprojDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Two entries for AAAAA11111 and one for BBBBB22222 — should return AAAAA11111.
	pbxprojContent := `
		DEVELOPMENT_TEAM = AAAAA11111;
		DEVELOPMENT_TEAM = AAAAA11111;
		DEVELOPMENT_TEAM = BBBBB22222;
`
	if err := os.WriteFile(filepath.Join(xcodeprojDir, "project.pbxproj"), []byte(pbxprojContent), 0644); err != nil {
		t.Fatal(err)
	}

	fakeWorkspace := filepath.Join(tmpDir, "MyApp.xcworkspace")
	teamID, err := DetectTeamID(fakeWorkspace)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if teamID != "AAAAA11111" {
		t.Errorf("expected 'AAAAA11111', got %q", teamID)
	}
}

func TestDetectAllWorkspaceCandidates(t *testing.T) {
	tmpDir := t.TempDir()

	// Create two workspace directories.
	ws1 := filepath.Join(tmpDir, "Alpha.xcworkspace")
	ws2 := filepath.Join(tmpDir, "sub", "Beta.xcworkspace")
	for _, d := range []string{ws1, ws2} {
		if err := os.MkdirAll(d, 0755); err != nil {
			t.Fatal(err)
		}
	}
	// Create a workspace inside Pods (should be skipped).
	podsWS := filepath.Join(tmpDir, "Pods", "Ignored.xcworkspace")
	if err := os.MkdirAll(podsWS, 0755); err != nil {
		t.Fatal(err)
	}

	candidates := detectAllWorkspaceCandidates(tmpDir)
	if len(candidates) != 2 {
		t.Fatalf("expected 2 candidates, got %d: %v", len(candidates), candidates)
	}
	// Results are sorted; Alpha comes before sub/Beta.
	if candidates[0] != "Alpha.xcworkspace" {
		t.Errorf("expected first candidate 'Alpha.xcworkspace', got %q", candidates[0])
	}
}
