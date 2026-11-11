---
apiVersion: processkit.projectious.work/v2
kind: WorkItem
metadata:
  id: BACK-20261101_1120-SmoothRiver-v1-05-parallel-runner-process-leak
  created: '2026-11-01T11:20:19+00:00'
  updated: '2026-11-01T11:40:29+00:00'
spec:
  title: Fix V1-05 AC-TOOLS parallel runner leak that exhausted the process limit
  state: backlog
  type: bug
  priority: high
  description: |
    Parallel AC-TOOLS runs (--jobs > 1) set TMPDIR per worker. With XDG_RUNTIME_DIR unset, rootless Podman keeps runtime state under TMPDIR, so containers started by the Dev Container CLI cannot be reached by plain `podman exec/stop/rm`. Reproduced with one container, no parallelism. Stuck containers keep about 27 tasks each and the cgroup limit is 1200 tasks; this most likely ended the Claude Code session on 2026-10-31 at about 21:22 CET.

    Resume after the container rebuild of 2026-11-01:
    1. /home/aibox/worktrees is not mounted and is lost on rebuild; branches survive in /workspace/.git. Recreate the worktree: `git worktree prune; git worktree add /home/aibox/worktrees/v1-05-speedup feat/v1-05-acceptance-speedup`.
    2. Apply the uncommitted fix (base df2cac38): `git apply /home/aibox/.local/share/aibox-rebuild-patches-20261101/v1-05-speedup.patch` (second copy in /var/cache/aibox-test/rebuild-patches-20261101/). It pins XDG_RUNTIME_DIR for worker engines and stops a run when a container cannot be removed or tasks reach 80% of pids.max.
    3. Check `cat /sys/fs/cgroup/pids.current` is low and remove leftover records: `podman ps -aq --filter label=aibox.v1-05.run | xargs -r podman rm -f -v`.
    4. Run unit tests, commit the fix, then repeat full and sampled runs with --jobs 4. Inspect the 3 failures of the earlier non-parallel run (full-arm64-20261031T182046Z) before closing V1-05.
---
