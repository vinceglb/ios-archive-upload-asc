package wizard

import (
	"strings"
	"testing"
)

func TestGenerateWorkflowContainsWorkspaceAndScheme(t *testing.T) {
	inputs := Inputs{
		Workspace: "MyApp.xcworkspace",
		Scheme:    "MyApp",
	}

	content, err := GenerateWorkflow(inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(content, "workspace: MyApp.xcworkspace") {
		t.Errorf("expected workspace in output, got:\n%s", content)
	}
	if !strings.Contains(content, "scheme: MyApp") {
		t.Errorf("expected scheme in output, got:\n%s", content)
	}
}

func TestGenerateWorkflowPreservesGitHubSyntax(t *testing.T) {
	inputs := Inputs{
		Workspace: "App.xcworkspace",
		Scheme:    "App",
	}

	content, err := GenerateWorkflow(inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// These GitHub Actions expressions must appear verbatim in the output.
	requiredExpressions := []string{
		"${{ secrets.ASC_KEY_ID }}",
		"${{ secrets.ASC_ISSUER_ID }}",
		"${{ secrets.ASC_PRIVATE_KEY_B64 }}",
		"${{ vars.BUNDLE_ID }}",
		"${{ vars.ASC_TEAM_ID }}",
		"${{ vars.ASC_APP_ID }}",
	}
	for _, expr := range requiredExpressions {
		if !strings.Contains(content, expr) {
			t.Errorf("expected %q in generated workflow, but it was missing or transformed", expr)
		}
	}
}

func TestGenerateWorkflowContainsPinnedActions(t *testing.T) {
	inputs := Inputs{
		Workspace: "App.xcworkspace",
		Scheme:    "App",
	}

	content, err := GenerateWorkflow(inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(content, "actions/archive@v0") {
		t.Errorf("expected actions/archive@v0 in output")
	}
	if !strings.Contains(content, "actions/upload@v0") {
		t.Errorf("expected actions/upload@v0 in output")
	}
}

func TestGenerateWorkflowTriggersOnTag(t *testing.T) {
	inputs := Inputs{Workspace: "App.xcworkspace", Scheme: "App"}
	content, err := GenerateWorkflow(inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(content, "tags:") || !strings.Contains(content, "v*") {
		t.Errorf("expected tag trigger in workflow output")
	}
}

func TestDefaultWorkflowPath(t *testing.T) {
	path := DefaultWorkflowPath()
	if path != ".github/workflows/release.yml" {
		t.Errorf("unexpected default workflow path: %q", path)
	}
}
