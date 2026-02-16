# `ios-gha-setup` (Deprecated Alias)

`ios-gha-setup` has been renamed to **`releasekit-ios-setup`**.

Use the new guide:
- [`docs/releasekit-ios-setup.md`](releasekit-ios-setup.md)

Use the new installer:

```bash
curl -fsSL https://raw.githubusercontent.com/vinceglb/releasekit-ios/main/scripts/install-releasekit-ios-setup.sh | sh
```

Then run:

```bash
releasekit-ios-setup wizard
```

Compatibility note:
- `scripts/ios-gha-setup.sh` still forwards to `scripts/releasekit-ios-setup.sh`.
- `scripts/install-ios-gha-setup.sh` still forwards to `scripts/install-releasekit-ios-setup.sh`.
