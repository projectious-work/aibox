# aibox v0.14.1

## Image size reduction

Moved heavy preview tools from the base image to a new **preview-enhanced** addon, significantly reducing base image size:

- **Removed from base**: `ffmpeg` (~450 MB), `imagemagick` (~50 MB), `ghostscript` (~70 MB)
- **Kept in base**: `chafa`, `librsvg2-bin`, `poppler-utils`, `timg`, `mupdf-tools`, `entr`
- **New addon**: `preview-enhanced` — install with `aibox addon add preview-enhanced` to get ffmpeg (video thumbnails), imagemagick (format conversion), and ghostscript (EPS rendering)

## Preview plugin improvements

- **eps.yazi**: Graceful fallback when ghostscript is not installed — shows helpful message directing users to the preview-enhanced addon
- **svg.yazi**: No changes (rsvg-convert fallback from v0.14.0 covers aarch64)

## Container Images

- `ghcr.io/projectious-work/aibox:base-debian-v0.14.1`
- `ghcr.io/projectious-work/aibox:base-debian-latest`

## CLI Binaries

- `aibox-v0.14.1-aarch64-unknown-linux-gnu.tar.gz`
- `aibox-v0.14.1-x86_64-unknown-linux-gnu.tar.gz`
- `aibox-v0.14.1-aarch64-apple-darwin.tar.gz` (added in Phase 2)
- `aibox-v0.14.1-x86_64-apple-darwin.tar.gz` (added in Phase 2)
