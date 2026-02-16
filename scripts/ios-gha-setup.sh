#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "[deprecated] 'ios-gha-setup' has been renamed to 'releasekit-ios-setup'." >&2
echo "[deprecated] Please switch to: scripts/releasekit-ios-setup.sh" >&2

exec "${SCRIPT_DIR}/releasekit-ios-setup.sh" "$@"
