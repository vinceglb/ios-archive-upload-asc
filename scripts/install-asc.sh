#!/usr/bin/env bash
set -euo pipefail

fail() {
  echo "::error::$*" >&2
  exit 1
}

require_non_empty() {
  local name="$1"
  local value="${2:-}"
  if [[ -z "${value}" ]]; then
    fail "${name} is required"
  fi
}

version="${INPUT_ASC_VERSION:-}"
require_non_empty "INPUT_ASC_VERSION" "${version}"

runner_temp="${RUNNER_TEMP:-}"
if [[ -z "${runner_temp}" ]]; then
  fail "RUNNER_TEMP is not set"
fi

arch="$(uname -m)"
case "${arch}" in
  arm64)
    asset_arch="arm64"
    ;;
  x86_64)
    asset_arch="amd64"
    ;;
  *)
    fail "Unsupported macOS architecture: ${arch}"
    ;;
esac

asset_name="asc_${version}_macOS_${asset_arch}"
checksums_name="asc_${version}_checksums.txt"
release_url="https://github.com/rudrankriyam/App-Store-Connect-CLI/releases/download/${version}"

tmp_dir="$(mktemp -d "${runner_temp}/asc-install.XXXXXX")"
cleanup() {
  rm -rf "${tmp_dir}"
}
trap cleanup EXIT

binary_path="${tmp_dir}/${asset_name}"
checksums_path="${tmp_dir}/${checksums_name}"

curl -fsSL --retry 3 --retry-all-errors -o "${binary_path}" "${release_url}/${asset_name}" \
  || fail "Failed to download asc binary for version ${version}"
curl -fsSL --retry 3 --retry-all-errors -o "${checksums_path}" "${release_url}/${checksums_name}" \
  || fail "Failed to download checksum file for version ${version}"

expected_sha="$(awk -v target="${asset_name}" '$2 == target {print $1}' "${checksums_path}")"
if [[ -z "${expected_sha}" ]]; then
  fail "Could not find checksum entry for ${asset_name}"
fi

actual_sha="$(shasum -a 256 "${binary_path}" | awk '{print $1}')"
if [[ "${actual_sha}" != "${expected_sha}" ]]; then
  fail "Checksum mismatch for ${asset_name}. Expected ${expected_sha}, got ${actual_sha}"
fi

bin_dir="${runner_temp}/bin"
mkdir -p "${bin_dir}"
install -m 0755 "${binary_path}" "${bin_dir}/asc"

if [[ -n "${GITHUB_PATH:-}" ]]; then
  echo "${bin_dir}" >> "${GITHUB_PATH}"
else
  export PATH="${bin_dir}:${PATH}"
fi

"${bin_dir}/asc" --version
