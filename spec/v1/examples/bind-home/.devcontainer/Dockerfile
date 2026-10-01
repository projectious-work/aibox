# Native starter: this file is project-owned after Template application.
# The index digest contains linux/amd64 and linux/arm64 manifests.
FROM debian:trixie-slim@sha256:a99cfc517144bc59b1978475ec53b46ecabec7e43635402ee5b77cc54cd1b20a
ARG USERNAME=aibox
ARG HARNESS=codex
# Immutable Debian snapshot pins package versions and their transitive repository.
RUN printf '%s\n' 'deb [check-valid-until=no] http://snapshot.debian.org/archive/debian/20260930T000000Z/ trixie main' > /etc/apt/sources.list \
    && rm -f /etc/apt/sources.list.d/debian.sources \
    && apt-get -o Acquire::https::Verify-Peer=true update \
    && apt-get install -y --no-install-recommends \
       bash=5.2.37-2+b10 ca-certificates=20250419 curl=8.14.1-2+deb13u5 \
       git=1:2.47.3-0+deb13u1 procps=2:4.0.4-9 tmux=3.5a-3 unzip=6.0-29+deb13u1 \
    && rm -rf /var/lib/apt/lists/*
COPY install-tools.sh /tmp/install-tools.sh
RUN sh /tmp/install-tools.sh "$HARNESS" && rm /tmp/install-tools.sh \
    && useradd --create-home --uid 1000 --shell /bin/bash "$USERNAME"
ENV HOME=/home/${USERNAME}
USER ${USERNAME}
CMD ["sleep", "infinity"]
