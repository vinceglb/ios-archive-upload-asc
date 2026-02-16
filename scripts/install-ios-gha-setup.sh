#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "[deprecated] 'install-ios-gha-setup.sh' has been renamed to 'install-releasekit-ios-setup.sh'." >&2
echo "[deprecated] Please switch to: scripts/install-releasekit-ios-setup.sh" >&2

exec "${SCRIPT_DIR}/install-releasekit-ios-setup.sh" "$@"
