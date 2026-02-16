# ios-archive-upload-asc

Composite GitHub Action to archive an iOS app, export an `.ipa`, and upload it to App Store Connect using [`asc`](https://github.com/rudrankriyam/App-Store-Connect-CLI).

## Requirements

- Runner: `macos-14` (or newer macOS runner with Xcode).
- App Store Connect API key with permissions for build upload and provisioning updates.
- Project configured for automatic signing.

## Inputs

### Required

| Name | Description |
| --- | --- |
| `workspace` | Path to the `.xcworkspace` file or directory. |
| `scheme` | Xcode scheme to archive. |
| `app_id` | App Store Connect app ID. |
| `bundle_id` | Expected bundle identifier for validation against the archive. |
| `asc_key_id` | App Store Connect API Key ID. |
| `asc_issuer_id` | App Store Connect API Issuer ID. |
| `asc_private_key_b64` | Base64-encoded `.p8` private key content. |
| `asc_team_id` | Apple Developer Team ID used for signing/export. |

### Optional

| Name | Default | Description |
| --- | --- | --- |
| `configuration` | `Release` | Xcode build configuration. |
| `archive_path` | `${{ runner.temp }}/archive/App.xcarchive` | Output archive path. |
| `export_path` | `${{ runner.temp }}/export` | Output export directory (contains IPA). |
| `asc_version` | `0.28.8` | `asc` CLI version to install. |
| `wait_for_processing` | `false` | If `true`, waits for ASC build processing. |
| `poll_interval` | `30s` | Poll interval when waiting. |
| `xcodebuild_extra_args` | `""` | Extra args appended to `xcodebuild archive`. |

## Outputs

| Name | Description |
| --- | --- |
| `archive_path` | Absolute path to generated `.xcarchive`. |
| `ipa_path` | Absolute path to generated `.ipa`. |
| `upload_id` | ASC upload ID returned by `asc builds upload`. |
| `file_id` | ASC upload file ID returned by `asc builds upload`. |
| `asc_result_json` | Compact JSON output from `asc builds upload`. |

## Usage

```yaml
name: Build and Upload iOS

on:
  workflow_dispatch:

jobs:
  upload:
    runs-on: macos-14
    steps:
      - uses: actions/checkout@v4

      - name: Archive, export and upload
        id: ios_upload
        uses: ga-vince/ios-archive-upload-asc@v1
        with:
          workspace: MyApp.xcworkspace
          scheme: MyApp
          app_id: ${{ vars.ASC_APP_ID }}
          bundle_id: com.example.myapp
          asc_key_id: ${{ secrets.ASC_KEY_ID }}
          asc_issuer_id: ${{ secrets.ASC_ISSUER_ID }}
          asc_private_key_b64: ${{ secrets.ASC_PRIVATE_KEY_B64 }}
          asc_team_id: ${{ secrets.ASC_TEAM_ID }}
          wait_for_processing: "true"

      - name: Print outputs
        run: |
          echo "Archive: ${{ steps.ios_upload.outputs.archive_path }}"
          echo "IPA: ${{ steps.ios_upload.outputs.ipa_path }}"
          echo "Upload ID: ${{ steps.ios_upload.outputs.upload_id }}"
```

## Troubleshooting

- `asc upload failed`
  - Verify `asc_key_id`, `asc_issuer_id`, `asc_private_key_b64`.
  - Verify API key permissions in App Store Connect.
- Signing/provisioning failures during archive/export
  - Confirm automatic signing is enabled for the target and `asc_team_id` is correct.
  - Confirm bundle ID is registered under the team.
- `Bundle ID mismatch`
  - The archive’s bundle identifier does not match `bundle_id`.
  - Check the scheme/target and input value.
- Upload succeeds locally but not in CI
  - Ensure the runner is macOS with compatible Xcode and network access to App Store Connect.
