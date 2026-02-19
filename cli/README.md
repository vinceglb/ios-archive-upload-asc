# releasekit-ios (Go CLI)

`releasekit-ios` is the active setup CLI for ReleaseKit-iOS.

Current scope:

- `releasekit-ios wizard`

The wizard is built with Cobra + Charm (`huh`, `lipgloss`) and focuses on collecting setup values, validating local inputs, and printing manual GitHub secrets/variables.

## Local development

From repository root:

```bash
cd cli
go test ./...
go run . wizard
```

## Build locally

```bash
cd cli
go build -o releasekit-ios .
./releasekit-ios wizard
```

## Show version information

```bash
cd cli
go run . version
```

## Install from GitHub release

Latest stable:

```bash
curl -fsSL https://raw.githubusercontent.com/vinceglb/releasekit-ios/main/install-cli.sh | sh
```

Pinned version:

```bash
curl -fsSL https://raw.githubusercontent.com/vinceglb/releasekit-ios/main/install-cli.sh | sh -s -- --version v0.1.0
```

## Release flow

Push a `v*` tag to publish macOS binaries:

- `releasekit-ios-darwin-arm64.tar.gz`
- `releasekit-ios-darwin-amd64.tar.gz`

Workflow: `.github/workflows/release-cli-beta.yml`

## Deferred scope

- `check`, `apply`, `doctor`
- GitHub sync via `gh`
- ASC API validation
- Workflow file generation
- Update checks
