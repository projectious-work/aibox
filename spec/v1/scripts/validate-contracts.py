# /// script
# requires-python = ">=3.10"
# dependencies = ["jsonschema[format]==4.23.0", "referencing==0.35.1"]
# ///
"""Validate the real JSON Schema vocabulary and regression fixtures offline after install."""
import copy
import json
from pathlib import Path
from jsonschema import Draft202012Validator, FormatChecker
from referencing import Registry, Resource

root = Path(__file__).resolve().parent.parent
schemas = {p.name: json.loads(p.read_text()) for p in root.glob("*.schema.json")}
registry = Registry().with_resources((s["$id"], Resource.from_contents(s)) for s in schemas.values())
validators = {}
for name, schema in schemas.items():
    Draft202012Validator.check_schema(schema)
    validators[name] = Draft202012Validator(schema, registry=registry, format_checker=FormatChecker())

count = 0
def check(name, value, valid=True):
    global count
    errors = list(validators[name].iter_errors(value))
    assert bool(errors) != valid, (name, valid, [e.message for e in errors][:3], value)
    count += 1

d = "sha256:" + "0" * 64
t = "2026-09-26T12:00:00Z"
base = dict(schemaVersion="aibox.operation-request/v1", requestId="case-1", projectRoot="/workspace")
requests = {
    "build_environment": dict(expectedInputDigest=d, runtimeContext="test", frozenLockfile=True, noCache=True),
    "start_environment": dict(expectedInputDigest=d, runtimeContext="test"),
    "stop_environment": dict(expectedInputDigest=d, runtimeContext="test", resourceId="abc"),
    "remove_environment": dict(expectedInputDigest=d, runtimeContext="test", resourceId="abc", acknowledgeDisruption=True),
    "rebuild_environment": dict(expectedInputDigest=d, runtimeContext="test", resourceId="abc", acknowledgeDisruption=True),
    "inspect_environment": dict(runtimeContext="test", effectiveConfig=True),
    "check_environment": dict(offline=True),
    "inspect_workspace": dict(effectiveConfig=True),
    "check_workspace": dict(offline=True),
    "refresh_workspace": dict(expectedInputDigest=d, scope="theme", theme="gruvbox", mode="dark"),
    "read_logs": dict(source="session", tailLines=100, since=t),
    "inspect_operation": dict(operationId="prior-1"),
    "migrate_preview": dict(sourceRoot="/old", destinationRoot="/new"),
    "migrate_apply": dict(planId="plan-1", expectedInputDigest=d),
    "migrate_rollback": dict(planId="plan-1", expectedInputDigest=d),
}
for operation, fields in requests.items():
    req = dict(base, operation=operation, **fields)
    check("operation-request.schema.json", req)
    check("operation-request.schema.json", dict(req, arbitraryShell="unsafe"), False)
    if operation.endswith("_workspace"):
        check("operation-request.schema.json", dict(req, runtimeContext="host"), False)
check("operation-request.schema.json", dict(base, operation="inspect_environment", runtimeContext="test", acknowledgeDisruption=True), False)
check("operation-request.schema.json", dict(base, operation="read_logs", source="runtime", tailLines=1001), False)
check("operation-request.schema.json", dict(base, operation="read_logs", source="session", tailLines=1, since="yesterday"), False)
check("operation-request.schema.json", dict(base, operation="rebuild_environment", **dict(requests["rebuild_environment"], acknowledgeDisruption=False)), False)
check("operation-request.schema.json", dict(base, operation="start_environment", runtimeContext="test"), False)
check("operation-request.schema.json", dict(base, operation="read_logs", source="runtime", tailLines=10), False)
check("operation-request.schema.json", dict(base, operation="read_logs", source="runtime", tailLines=10, runtimeContext="test", resourceId="abc"))

envelope = dict(schemaVersion="aibox.operation-result/v1", requestId="case-1", actors=dict(initiator="user", executor="operator"), outcome="succeeded", target=dict(scope="operator", workspaceRoot="/workspace"), inputDigest=d, changedResources=[], warnings=[], evidence=[])
inspection = dict(state="absent", resources=[], configuration=[], observedAt=t)
payloads = {op: dict(resources=[], retainedData=[]) for op in ["start_environment", "stop_environment", "remove_environment", "rebuild_environment"]}
payloads.update(build_environment=dict(images=[]), inspect_environment=inspection, inspect_workspace=inspection,
    check_environment=dict(findings=[], complete=True), check_workspace=dict(findings=[], complete=True),
    read_logs=dict(source="session", entries=[dict(text="redacted")], truncated=False),
    refresh_workspace=dict(changedFiles=[], rebuildRequired=False, sessionRestartRequired=False),
    inspect_operation=dict(operationId="prior-1", originalOperation="start_environment", state="partial", lastConfirmedStep="created", completedEffects=[], unknownEffects=[], nextAction="inspect_environment"))
for op in ["migrate_preview", "migrate_apply", "migrate_rollback"]:
    payloads[op] = dict(planId="plan-1", manifest="state:plan-1", conflicts=[], activated=False)
for operation, data in payloads.items():
    result = dict(envelope, operation=operation, data=data)
    check("operation-result.schema.json", result)
    bad = copy.deepcopy(result)
    bad["data"]["secret"] = "unexpected"
    check("operation-result.schema.json", bad, False)
# First build and malformed requests have no discovered container or input hash.
early = {k: v for k, v in envelope.items() if k not in ["target", "inputDigest"]}
error = dict(code="invalid_input", category="invalid_input", message="Invalid source", retryable=False, nextAction="edit_configuration")
check("operation-result.schema.json", dict(early, operation="invalid_request", outcome="failed", error=error))
check("operation-result.schema.json", dict(early, operation="build_environment", outcome="refused", error=dict(error, category="denied")))
check("operation-result.schema.json", dict(early, operation="build_environment", outcome="timed_out", error=dict(error, category="timeout")))
check("operation-result.schema.json", dict(early, operation="start_environment", outcome="partial", error=error), False)
check("operation-result.schema.json", dict(envelope, operation="read_logs"), False)
check("operation-result.schema.json", dict(envelope, operation="read_logs", data=payloads["read_logs"], error=error), False)
policy = dict(schemaVersion="aibox.operator-policy/v1", principal="operator", allowedRoots=["/workspace"], executables=dict(devcontainer="/usr/bin/devcontainer", runtime="/usr/bin/docker"), runtimeContexts=[dict(name="test", kind="docker", endpointFingerprint=d)], operations=["start_environment"], grants=[])
check("operator-policy.schema.json", policy)
check("operator-policy.schema.json", dict(policy, allowEverything=True), False)
check("operator-policy.schema.json", dict(policy, allowedRoots=["relative"]), False)
operation_record = dict(schemaVersion="aibox.operation-record/v1", operationId="case-1", requestFingerprint=d, operation="start_environment", projectRoot="/workspace", inputDigest=d, createdAt=t, updatedAt=t, executor="operator", state="executing", lastConfirmedStep="authorized", resources=[], completedEffects=[], unknownEffects=[])
check("operation-record.schema.json", operation_record)
check("operation-record.schema.json", dict(operation_record, result=dict(envelope, operation="start_environment", data=payloads["start_environment"])))
check("operation-record.schema.json", dict(operation_record, createdAt="not-a-date"), False)
roadmap = json.loads((root / "roadmap.yaml").read_text())
check("roadmap.schema.json", roadmap)
missing_demo = copy.deepcopy(roadmap)
missing_demo.pop("demos")
check("roadmap.schema.json", missing_demo, False)
invalid_demo = copy.deepcopy(roadmap)
invalid_demo["demos"]["V1-03"]["run"] = ""
check("roadmap.schema.json", invalid_demo, False)
missing_completion_commit = copy.deepcopy(roadmap)
missing_completion_commit["groups"][0]["phases"][2].pop("implementationCommit")
check("roadmap.schema.json", missing_completion_commit, False)
missing_completion_evidence = copy.deepcopy(roadmap)
missing_completion_evidence["groups"][0]["phases"][2].pop("evidence")
check("roadmap.schema.json", missing_completion_evidence, False)
unreleased_shipped = copy.deepcopy(roadmap)
unreleased_shipped["groups"][0]["phases"][2]["status"] = "shipped"
check("roadmap.schema.json", unreleased_shipped, False)
released_done = copy.deepcopy(unreleased_shipped)
released_done["groups"][0]["phases"][2]["release"] = "v1.0.0-alpha.1"
check("roadmap.schema.json", released_done)
missing_completion_note = copy.deepcopy(roadmap)
missing_completion_note["groups"][0]["phases"][2].pop("devNote")
check("roadmap.schema.json", missing_completion_note, False)
legacy_completion = copy.deepcopy(roadmap)
legacy_completion["groups"][0]["phases"][2]["implementationStatus"] = "done"
check("roadmap.schema.json", legacy_completion, False)
check("customization.schema.json", dict(schemaVersion="1", latex=dict(preview=dict(document="overview"))))
for example_name in ("minimal", "customized"):
    example_path = root / "examples" / example_name / ".devcontainer" / "devcontainer.json"
    # The maintained files are JSONC option catalogs: whole-line comments
    # describe inactive choices and must not become active configuration.
    active_lines = (line for line in example_path.read_text().splitlines()
                    if not line.lstrip().startswith("//"))
    example = json.loads("\n".join(active_lines))
    assert "features" not in example, f"{example_name}: unqualified Feature reference"
    assert "aibox" not in example, f"{example_name}: misplaced aibox configuration"
    extension = example.get("customizations", {}).get("aibox")
    if example_name == "minimal":
        assert extension is None, "minimal: customization must remain optional"
    else:
        assert extension is not None, "customized: missing aibox customization"
        check("customization.schema.json", extension)
check("customization.schema.json", dict(schemaVersion="2"), False)
check("customization.schema.json", dict(schemaVersion="1", workspace=dict(theme=None)), False)
for key, value in {"enabled": True, "port": 8765, "bind": "0.0.0.0", "allow_public": True, "engine": "native"}.items():
    check("customization.schema.json", dict(schemaVersion="1", latex=dict(preview={key: value})), False)
settings = dict(schemaVersion="1", output=dict(format="json", color="never"), logging=dict(level="info", format="jsonl", rotationMiB=10, retentionFiles=7, retentionDays=7), execution=dict(timeoutSeconds=300))
check("process-settings.schema.json", settings)
check("process-settings.schema.json", dict(settings, operatorPolicy="/unsafe"), False)
check("process-settings.schema.json", dict(settings, logging=dict(level="verbose")), False)
event = dict(schemaVersion="aibox.log-event/v1", timestamp=t, severity="info", event="operation.started", component="app", requestId="case-1", sequence=0, fields=dict(scope="operator"))
check("log-event.schema.json", event)
check("log-event.schema.json", dict(event, secret="unexpected"), False)
check("log-event.schema.json", dict(event, sequence=-1), False)
print(f"Draft 2020-12 validation passed: {len(schemas)} schemas, {count} contract fixtures.")
