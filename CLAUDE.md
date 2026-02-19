# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ReleaseKit-iOS is a GitHub Actions-based CI/CD toolkit for distributing iOS apps to the App Store. It provides:
- **Two composite GitHub Actions** (`actions/archive/`, `actions/upload/`) that wrap Xcode and the `asc` CLI
- **A Go CLI wizard** (`cli/releasekit-ios-go/`) that guides developers through initial setup
- **Bash scripts** (`scripts/`) that implement the archive/upload logic
- **An installer script** (`install-cli.sh`) for end-user CLI installation

The key design: App Store Connect API handles cloud signing, so no local certificates are needed.

## Commands

### Go CLI (from `cli/releasekit-ios-go/`)

```bash
go test ./...             # Run unit tests
go build -o releasekit-ios .   # Build binary
go run . wizard           # Run wizard locally
```

### Bash validation (from repo root)

```bash
bash -n scripts/archive.sh scripts/upload.sh scripts/lib/common.sh  # Check syntax
```

### CI runs (via GitHub Actions)

- `bash -n` on all scripts in `scripts/`
- `go test ./...` in the Go CLI directory
- YAML validation via Ruby
- Action metadata key checks

## Architecture

### Actions → Scripts → Tools

The composite actions (`actions/archive/action.yml`, `actions/upload/action.yml`) pass inputs as environment variables to shell scripts, which then call `xcodebuild` and `asc`. Results flow back via `GITHUB_OUTPUT`.

### Shared Bash Library

`scripts/lib/common.sh` provides utilities used by both `archive.sh` and `upload.sh`:
- `fail()` / `require_non_empty()` — error handling with GitHub Actions annotations
- `prepare_private_key_file()` — decodes base64 (including double-encoded) ASC `.p8` keys
- `decode_base64()` — handles macOS vs Linux `base64` flag differences
- `parse_json_field()` — JSON extraction with `jq` fallback to `sed`

### Go CLI Wizard Flow

`cmd/wizard.go` → `internal/wizard/run.go` orchestrates:
1. **Detect** (`detect.go`): walks filesystem to find `.xcworkspace`, skipping build dirs
2. **Collect** (`ui.go`): huh-based interactive forms; accepts ASC key as file path or raw base64
3. **Validate** (`validate.go`): non-empty fields, file existence, valid base64; `normalizeBase64()` and `encodeFileBase64()` normalize key input
4. **Print summary** (`summary.go`): outputs values as GitHub Secrets/Variables blocks for copy-paste

Styling lives in `internal/term/style.go` (lipgloss-based `Theme` struct).

The CLI's module path is `github.com/vinceglb/releasekit-ios/cli/releasekit-ios-go`.

### Release Process

Pushing a `v*` tag triggers `.github/workflows/release-cli-beta.yml`, which:
1. Runs `go test ./...`
2. Validates the tag format (`vX.Y[.Z][-suffix]`)
3. Builds `darwin/arm64` and `darwin/amd64` binaries with `-trimpath -ldflags='-s -w'`
4. Creates tarballs and a GitHub Release via `softprops/action-gh-release@v2`

`install-cli.sh` (repo root) is the end-user installer: detects macOS architecture, resolves the latest tag from the GitHub API (or accepts `--version`), downloads the tarball, and installs to `~/.local/bin/releasekit-ios`.

### Smoke Tests

`.github/workflows/smoke.yml` is a manual `workflow_dispatch` workflow with scenario selection:
- **Local negative tests** (ubuntu): missing/invalid inputs for archive and upload
- **Split E2E tests** (macos): happy-path, bundle-mismatch, bad-app-id

### Legacy Code

The `old/` directory contains **frozen** archived files — the shell-based setup wizard (`old/releasekit-ios-setup.sh`) is replaced by the Go CLI. Do not add features there.

## Key Dependencies

- **Go**: Cobra (CLI), charmbracelet/huh (interactive forms), charmbracelet/lipgloss (styling)
- **External tools**: `xcodebuild` (macOS only), `asc` (App-Store-Connect-CLI), `jq` (optional)
- **Actions**: `rudrankriyam/setup-asc@v1` installs `asc` in workflows

## Conventions

- Scripts use `set -euo pipefail` throughout
- Secrets are masked early with `::add-mask::` in scripts
- Actions are pinned to `@v0` (the stable semver major tag)
- Go `CGO_ENABLED=0` for portable binaries
