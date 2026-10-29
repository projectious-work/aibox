#!/usr/bin/env python3
"""V1-05 AC-TOOLS black-box acceptance suite.

Host prerequisites: Python 3.10+ and Docker. The pinned Dev Container CLI is
bootstrapped from tools/vendor like the V1-04 host verifier. A run starts an
OCI registry on localhost:5000, publishes the Features with
`scripts/v1-05-features.sh publish --registry localhost:5000/aibox-features`,
builds workspaces from the catalog's reference base, evaluates every planned
test on the host's native target, removes the run's containers, images and
registry, and writes evidence.json to the output directory.

Examples:
  run.sh --list
  run.sh --dry-run --item present --feature 'python/*'
  run.sh --item present,version --feature go/go
  run.sh                      # every AC-TOOLS item on this host's target
"""

from __future__ import annotations

import argparse
import json
import platform
import shutil
import subprocess
import sys
import tempfile
from collections import Counter
from datetime import datetime, timezone
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))

from actools import host  # noqa: E402
from actools.catalog import REPO, SPEC, Catalog  # noqa: E402
from actools.executor import Context, Executor  # noqa: E402
from actools.plan import DOCKER_ITEMS, ITEMS, Plan, census_rows, make_plan  # noqa: E402


def now() -> str:
    return datetime.now(timezone.utc).isoformat(timespec="seconds").replace("+00:00", "Z")


def parse_args(argv: list[str]) -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__,
                                     formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument("--item", action="append", default=[],
                        help=f"AC-TOOLS items, comma separated (default all): {', '.join(ITEMS)}")
    parser.add_argument("--feature", action="append", default=[],
                        help="fnmatch pattern for catalog Feature IDs or smoke IDs (repeatable)")
    parser.add_argument("--test", action="append", default=[],
                        help="fnmatch pattern for full test IDs (repeatable)")
    parser.add_argument("--target", help="catalog target ID; a run must match the host "
                                         "(default: detected; --list/--dry-run: all targets)")
    parser.add_argument("--list", action="store_true", help="print planned test IDs and exit")
    parser.add_argument("--dry-run", action="store_true",
                        help="print planned builds with their Feature selections and tests")
    parser.add_argument("--output", help="new absolute evidence directory "
                                         "(default: a new temporary directory)")
    parser.add_argument("--work-dir", help="parent for temporary workspaces; must be shared "
                                           "with the Docker daemon (default: system temp)")
    parser.add_argument("--aibox", help="aibox binary for AC-TOOLS/inspect "
                                        "(default: build from this checkout in a Go container)")
    parser.add_argument("--registry-image", default=host.DEFAULT_REGISTRY_IMAGE,
                        help="OCI Distribution registry image for localhost:5000")
    return parser.parse_args(argv)


def build_plan(catalog: Catalog, args: argparse.Namespace, targets: list[str]) -> Plan:
    addons = json.loads((SPEC / "ledger" / "addons.json").read_text())
    return make_plan(catalog, targets, census_rows(addons), items=args.item,
                     features=args.feature, tests=args.test)


def print_dry_run(catalog: Catalog, plan: Plan) -> None:
    by_build: dict[str | None, list] = {}
    for test in plan.tests:
        by_build.setdefault(test.build, []).append(test)
    for key, tests in by_build.items():
        if key is None:
            print("offline (no container):")
        else:
            build = plan.builds[key]
            print(f"build {key}: fixture={build.fixture} "
                  f"{'expect-failure ' if build.expect_failure else ''}"
                  f"copies={list(build.copies)}")
            print("  features: " + json.dumps(catalog.devcontainer_features(
                build.selection_dict()), sort_keys=True))
        for test in tests:
            print(f"  {test.id}")
    print_summary(plan)


def print_summary(plan: Plan) -> None:
    counts = Counter(test.item for test in plan.tests)
    print(f"# {len(plan.tests)} tests, {len(plan.builds)} container builds")
    for item in ITEMS:
        if counts[item]:
            print(f"#   AC-TOOLS/{item}: {counts[item]}")


def source_state() -> dict:
    try:
        commit = subprocess.run(["git", "rev-parse", "HEAD"], cwd=REPO, text=True,
                                capture_output=True, check=True).stdout.strip()
        dirty = bool(subprocess.run(["git", "status", "--porcelain"], cwd=REPO, text=True,
                                    capture_output=True, check=True).stdout.strip())
        return {"sourceCommit": commit, "sourceDirty": dirty}
    except (OSError, subprocess.CalledProcessError):
        return {"sourceCommit": None, "sourceDirty": None}


def prepare_output(value: str | None) -> Path:
    if value is None:
        return Path(tempfile.mkdtemp(prefix="aibox-v1-05-evidence-")).resolve()
    output = Path(value)
    if not output.is_absolute() or output.exists():
        raise RuntimeError("--output must name a new absolute directory")
    output.mkdir(parents=True, mode=0o700)
    return output


def main(argv: list[str]) -> int:
    args = parse_args(argv)
    catalog = Catalog.load()
    if args.list or args.dry_run:
        targets = [args.target] if args.target else list(catalog.targets)
        plan = build_plan(catalog, args, targets)
        if args.dry_run:
            print_dry_run(catalog, plan)
        else:
            for test in plan.tests:
                print(test.id)
            print_summary(plan)
        return 0

    started = now()
    run_id = datetime.now(timezone.utc).strftime("%Y%m%dt%H%M%S")
    output = prepare_output(args.output)
    work = Path(tempfile.mkdtemp(prefix=f"aibox-v1-05-{run_id}-", dir=args.work_dir)).resolve()
    context = Context(catalog, run_id, output, work)
    evidence: dict = {"phase": "V1-05", "suite": "AC-TOOLS acceptance", "runId": run_id,
                      "startedAt": started, **source_state(),
                      "host": {"system": platform.system(), "machine": platform.machine(),
                               "python": platform.python_version()},
                      "selection": {"items": args.item, "features": args.feature,
                                    "tests": args.test}}
    failures: list[str] = []
    cleanup: list[str] = []
    results: list[dict] = []
    plan = None
    log = output / "host.log"
    registry = None
    try:
        probe_plan = build_plan(catalog, args, list(catalog.targets)[:1])
        needs_docker = any(t.item in DOCKER_ITEMS for t in probe_plan.tests) or (
            any(t.item == "inspect" for t in probe_plan.tests) and not args.aibox)
        target = args.target
        if needs_docker:
            context.docker = host.Docker.locate()
            architecture = context.docker.server_architecture()
            detected = catalog.target_for_architecture(architecture)
            if target and target != detected:
                raise RuntimeError(f"--target {target} does not match the host target {detected}")
            target = detected
            evidence["engine"] = {"binary": context.docker.binary,
                                  "serverVersion": context.docker.version(),
                                  "serverArch": architecture}
        evidence["target"] = target
        plan = build_plan(catalog, args, [target] if target else [])
        evidence["planned"] = len(plan.tests)
        if not plan.tests:
            raise RuntimeError("the selection plans no tests")
        if any(t.item == "inspect" for t in plan.tests):
            if args.aibox:
                context.aibox = Path(args.aibox).resolve()
                evidence["aibox"] = {"binary": str(context.aibox)}
            else:
                assert context.docker
                context.aibox, evidence["aibox"] = host.build_aibox(context.docker, work, log)
        if any(t.build for t in plan.tests):
            assert context.docker
            context.cli = host.DevContainerCli.bootstrap(work, output)
            evidence["devcontainerCli"] = context.cli.runtime
            registry = host.Registry(context.docker, args.registry_image,
                                     f"aibox-v1-05-registry-{run_id}")
            context.registry = registry
            registry.start(log)
            evidence["registry"] = {"image": args.registry_image,
                                    "digest": context.docker.image_digest(args.registry_image),
                                    "prefix": f"{host.REGISTRY_HOST}/{host.REGISTRY_PREFIX}"}
            published = host.publish(context.cli, output / "publish.log")
            evidence["publish"] = {"command": published.argv, "exitCode": published.exit_code,
                                   "log": "publish.log"}
            if published.exit_code:
                failures.append(f"publish failed (exit {published.exit_code}); see publish.log")
        executor = Executor(context, plan)
        executor.run_offline()
        if not failures:
            executor.run_builds()
        results = executor.results
        cleanup += executor.cleanup_failures
    except BaseException as exc:  # recorded; cleanup still runs
        failures.append(str(exc) or type(exc).__name__)
        print(f"V1-05 acceptance run failed: {failures[-1]}", file=sys.stderr)
    finally:
        if registry is not None:
            cleanup += registry.stop(log)
        if context.docker is not None:
            cleanup += context.docker.remove_images(f"vsc-aibox-v1-05-{run_id}-*", log)
            cleanup += context.docker.remove_images(f"aibox-v1-05-{run_id}-*", log)
        shutil.rmtree(work, ignore_errors=True)

    if plan is not None:
        recorded = {r["id"] for r in results}
        for test in plan.tests:
            if test.id not in recorded:
                results.append({"id": test.id, "item": f"AC-TOOLS/{test.item}",
                                "feature": test.feature, "tool": test.tool,
                                "target": test.target, "qualifiers": dict(test.qualifiers),
                                "build": test.build, "status": "not-run",
                                "failures": ["not run because the run failed earlier"]})
    by_item: dict[str, Counter] = {}
    for record in results:
        by_item.setdefault(record["item"], Counter())[record["status"]] += 1
    passed = bool(results) and all(r["status"] == "passed" for r in results) \
        and not failures and not cleanup
    evidence.update({"finishedAt": now(), "status": "passed" if passed else "failed",
                     "summary": {"tests": len(results),
                                 "passed": sum(r["status"] == "passed" for r in results),
                                 "failed": sum(r["status"] != "passed" for r in results),
                                 "byItem": {k: dict(v) for k, v in sorted(by_item.items())}},
                     "failures": failures, "cleanupFailures": cleanup, "results": results})
    (output / "evidence.json").write_text(json.dumps(evidence, indent=2) + "\n")
    print(f"V1-05 AC-TOOLS: {evidence['summary']['passed']} passed, "
          f"{evidence['summary']['failed']} not passed; evidence: {output / 'evidence.json'}")
    return 0 if passed else 1


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
