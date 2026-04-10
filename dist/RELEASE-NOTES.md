# aibox v0.17.11

Bug-fix release: `[aibox].version = "latest"` no longer generates a broken Dockerfile.
No processkit version change — still compatible with v0.8.0.

## Fix: `[aibox].version = "latest"` produced invalid Docker image tag

When `[aibox].version = "latest"` was set, `aibox sync` generated:

```dockerfile
FROM ghcr.io/projectious-work/aibox:base-debian-vlatest AS aibox
```

The tag `base-debian-vlatest` does not exist in GHCR (image tags follow
`base-<flavor>-v<semver>`), causing the Docker build to fail.

The fix mirrors how processkit's `"latest"` is handled: `cmd_sync` now resolves
`[aibox].version = "latest"` to a concrete image version via the GHCR tags API
before Dockerfile generation, and logs the resolved version:

```
==> Resolved aibox image 'latest' → v0.17.11
```

The `"latest"` sentinel stays in `aibox.toml`; only the generated `Dockerfile`
uses the concrete resolved tag.
