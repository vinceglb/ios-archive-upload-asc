package wizard

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"
)

// DefaultWorkflowPath returns the conventional path for the release workflow.
func DefaultWorkflowPath() string {
	return ".github/workflows/release.yml"
}

// GenerateWorkflow renders the release workflow YAML from the provided inputs.
// Template uses [[ ]] delimiters to avoid collision with GitHub Actions ${{ }} syntax.
func GenerateWorkflow(inputs Inputs) (string, error) {
	// Note: delimiters are [[ ]] — NOT {{ }} — so that GitHub Actions ${{ secrets.X }}
	// syntax is passed through verbatim and not interpreted by text/template.
	tmpl, err := template.New("workflow").Delims("[[", "]]").Parse(workflowTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, inputs); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// WriteWorkflow writes content to path, creating intermediate directories.
func WriteWorkflow(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

// workflowTemplate is the GitHub Actions workflow template.
// Uses [[ ]] Go template delimiters so ${{ }} GitHub Actions expressions are untouched.
const workflowTemplate = `name: Release iOS App

on:
  workflow_dispatch:
  push:
    tags:
      - 'v*'

jobs:
  archive:
    name: Archive
    runs-on: macos-latest
    outputs:
      ipa-path: ${{ steps.archive.outputs.ipa_path }}

    steps:
      - uses: actions/checkout@v4

      - name: Setup ASC
        uses: rudrankriyam/setup-asc@v1

      - name: Archive
        id: archive
        uses: vinceglb/releasekit-ios/actions/archive@v0
        with:
          workspace: [[.Workspace]]
          scheme: [[.Scheme]]
          bundle-id: ${{ vars.BUNDLE_ID }}
          team-id: ${{ vars.ASC_TEAM_ID }}
          asc-key-id: ${{ secrets.ASC_KEY_ID }}
          asc-issuer-id: ${{ secrets.ASC_ISSUER_ID }}
          asc-private-key-b64: ${{ secrets.ASC_PRIVATE_KEY_B64 }}

  upload:
    name: Upload
    runs-on: macos-latest
    needs: archive

    steps:
      - uses: actions/checkout@v4

      - name: Setup ASC
        uses: rudrankriyam/setup-asc@v1

      - name: Upload
        uses: vinceglb/releasekit-ios/actions/upload@v0
        with:
          ipa-path: ${{ needs.archive.outputs.ipa-path }}
          app-id: ${{ vars.ASC_APP_ID }}
          asc-key-id: ${{ secrets.ASC_KEY_ID }}
          asc-issuer-id: ${{ secrets.ASC_ISSUER_ID }}
          asc-private-key-b64: ${{ secrets.ASC_PRIVATE_KEY_B64 }}
`
