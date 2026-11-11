---
apiVersion: processkit.projectious.work/v2
kind: WorkItem
metadata:
  id: BACK-20261024_1908-PatientTower-v0-base-image-trixie-cves
  created: '2026-10-24T19:08:59+00:00'
  labels:
    line: v0
    area: images
    release: v0.35.4
spec:
  title: 'v0 base image: track no-fix Debian trixie CVEs and apply trixie-security
    updates at build'
  state: backlog
  type: task
  priority: medium
  description: |
    Deferred from the v0.35.4 release-host vulnerability policy (WARN, non-blocking). 20 High/Critical findings in the base image, all "no fix listed": curl/libcurl (CVE-2026-10536, -11856, -18924, -19931 Critical; -12064, -13608 High), python3.13 (CVE-2026-19445 Critical; -15308, -19553), glibc (CVE-2026-19499), expat (CVE-2026-113351, -115506), ncurses (CVE-2025-69720), gcc-14 runtime (CVE-2026-102010), cyrus-sasl (CVE-2026-107161), krb5 (CVE-2026-107778), libarchive (CVE-2026-14164), setuptools wheel (CVE-2026-23949), vim (CVE-2026-26269, -28417).

    Checked against the Debian security tracker on 2026-10-24 for all five Criticals plus glibc and both expat CVEs: trixie is vulnerable, every one is marked `<no-dsa> (Minor issue)`, and fixes exist only in forky/sid. There is no trixie package to move to today.

    Follow-up actions:
    1. images/base-debian/Dockerfile never runs `apt-get upgrade`, so the image inherits whatever debian:trixie-slim shipped at its last rebuild. trixie-security already carries newer curl (8.14.1-2+deb13u7) and expat (2.9.0-1~deb13u1) builds than trixie main. Add an `apt-get upgrade` (or explicit security upgrade) step in the foundation stage, so future Debian security fixes are picked up on every release build.
    2. Re-check these CVEs at the next v0 release; promote to blocking once Debian lists a trixie fix (the release-host policy already blocks Critical findings that have a fixed version).
    3. Consider whether vim, python3.13, and libarchive need to be in the runtime image or can be slimmed.
---
