"""Settings journeys J-CFG-01 to J-CFG-24 of journeys.md."""

from __future__ import annotations

import os

from . import jsonc
from .contract import CONTRACT, DEMO_PROJECT_RELATIVE, REPO
from .expect import (check_entry, configuration, entry, expect_bootstrap, expect_rejected,
                     expect_success, result_errors, settings_entries, settings_sources)
from .events import parse_jsonl
from .framework import T, journey
from .layers import PAIR_VALUES, LayerSetup, text, view_value
from .mcp import MCPError, MCPSession, is_rejected, structured
from .sandbox import Sandbox, snapshot

NOT_JSON = "not json"


def p_devcontainer(project: str) -> str:
    return f"{project}/.devcontainer/devcontainer.json"


def expect_default_entries(t: T, sb: Sandbox, doc, res, project: str) -> None:
    """The J-CFG-01 expectation for a View of P with no settings anywhere."""
    conf = configuration(doc)
    t.eq([e.get("key") for e in conf], ["devcontainer.config"] + CONTRACT.key_names,
         "data.configuration is devcontainer.config then the 15 keys in catalog order",
         result=res)
    if conf:
        t.eq(conf[0], CONTRACT.native_entry("devcontainer.config"),
             "native entry devcontainer.config", result=res)
    for spec in CONTRACT.keys:
        key = spec["key"]
        common = dict(domain=spec["domain"], redacted=False)
        if key == "output.format":
            check_entry(t, doc, key, res, value="json", layer="invocation", source="--format",
                        invocation_only=True, **common)
        elif key == "logging.format":
            check_entry(t, doc, key, res, value="jsonl", layer="invocation",
                        source="--log-format", invocation_only=True, **common)
        else:
            check_entry(t, doc, key, res, value=CONTRACT.default(key), layer="built_in",
                        source="built-in", invocation_only=False, overridden=None, **common)
    t.eq(settings_sources(doc), [
        {"layer": "system", "path": CONTRACT.system_settings_path, "state": "absent"},
        {"layer": "user", "path": "~/.config/aibox/settings.json", "state": "absent"},
        {"layer": "project", "path": sb.tilde(p_devcontainer(project)), "state": "loaded"},
    ], "data.settingsSources", result=res)


@journey("J-CFG-01", ["CFG-001", "CFG-004", "CFG-060", "CFG-061", "CFG-063", "CFG-068", "EXE-003"])
def j_cfg_01(t: T) -> None:
    sb = t.local()
    p = sb.project("P")
    res = sb.view(p)
    doc = expect_success(t, res)
    expect_default_entries(t, sb, doc, res, p)


@journey("J-CFG-02", ["CFG-001", "CFG-062"], root=True)
def j_cfg_02(t: T) -> None:
    for spec in CONTRACT.keys:
        key = spec["key"]
        layers = spec["layers"]
        pairs = [(lo, hi) for i, lo in enumerate(layers) for hi in layers[i + 1:]]
        if key == "output.format":
            pairs = [(lo, hi) for lo, hi in pairs if hi == "invocation"]
        for lo, hi in pairs:
            with t.case(f"{key}: {lo} < {hi}"):
                label = f"{key}-{lo}-{hi}"
                sb = t.container(label) if "system" in (lo, hi) else t.local(label)
                setup = LayerSetup(sb)
                log_format = True
                if key == "output.format":
                    # The higher value is the View's own --format json.
                    lo_v, hi_v = "yaml", "json"
                    setup.set(key, lo, lo_v)
                elif key == "logging.format" and hi == "invocation":
                    # The higher value is the View's own --log-format jsonl.
                    lo_v, hi_v = "plain", "jsonl"
                    setup.set(key, lo, lo_v)
                else:
                    lo_v, hi_v = PAIR_VALUES[key]
                    setup.set(key, lo, lo_v)
                    setup.set(key, hi, hi_v)
                    # Pairs of logging.format among lower layers: View without --log-format.
                    log_format = key != "logging.format"
                project, args, env = setup.apply()
                res = sb.view(project, *args, env=env, log_format=log_format)
                doc = expect_success(t, res)
                check_entry(t, doc, key, res, value=view_value(sb, key, hi_v), layer=hi,
                            source=setup.source(key, hi, project),
                            invocation_only=(hi == "invocation"),
                            overridden=[{"layer": lo, "source": setup.source(key, lo, project),
                                         "value": view_value(sb, key, lo_v), "redacted": False}])


@journey("J-CFG-03", ["CFG-001", "CFG-062"], root=True)
def j_cfg_03(t: T) -> None:
    sb = t.container("color")
    setup = LayerSetup(sb)
    for layer, value in (("system", "never"), ("user", "always"), ("env_file", "never"),
                         ("environment", "always"), ("invocation", "never")):
        setup.set("output.color", layer, value)
    project, args, env = setup.apply()
    res = sb.view(project, *args, env=env)
    doc = expect_success(t, res)
    src = lambda layer: setup.source("output.color", layer, project)  # noqa: E731
    check_entry(t, doc, "output.color", res, value="never", layer="invocation", source="--color",
                invocation_only=True, overridden=[
                    {"layer": "environment", "source": src("environment"), "value": "always", "redacted": False},
                    {"layer": "env_file", "source": src("env_file"), "value": "never", "redacted": False},
                    {"layer": "user", "source": src("user"), "value": "always", "redacted": False},
                    {"layer": "system", "source": src("system"), "value": "never", "redacted": False},
                ])


J04_VALUES = {
    "logging.file": ("x.log", "y.log"),
    "logging.required": (True, False),
    "logging.rotationMiB": (5, 6),
    "logging.retentionFiles": (3, 4),
    "logging.retentionDays": (10, 20),
}


@journey("J-CFG-04", ["CFG-002", "CFG-003", "LOG-006"])
def j_cfg_04(t: T) -> None:
    for key, (env_value, user_value) in J04_VALUES.items():
        name = CONTRACT.key(key)["env"]
        section, member = key.split(".", 1)
        for layer in ("environment", "env_file"):
            with t.case(f"{key} via {layer}"):
                sb = t.local(f"{key}-{layer}")
                sb.user_settings({"schemaVersion": "1", section: {member: user_value}})
                p = sb.project("P")
                extra, env = [], {}
                if layer == "environment":
                    env[name] = text(env_value)
                else:
                    path = sb.env_file("auth.env", f"{name}={text(env_value)}\n")
                    extra = ["--env-file", path]
                res = sb.aibox_run(sb.rejection_args(p, *extra), env)
                expect_rejected(t, res, "unauthorized_source", key=key, layer=layer)


@journey("J-CFG-05", ["CFG-005", "CFG-012"])
def j_cfg_05(t: T) -> None:
    sb = t.local()
    p = sb.project("P")
    before = (snapshot(sb.home), snapshot(p))
    res = sb.view(p)
    expect_success(t, res)
    after = (snapshot(sb.home), snapshot(p))
    t.eq(sorted(after[0]), sorted(before[0]), "HOME tree paths unchanged", result=res)
    t.check(after[0] == before[0], "HOME tree identical", result=res)
    t.eq(sorted(after[1]), sorted(before[1]), "project tree paths unchanged", result=res)
    t.check(after[1] == before[1], "project tree identical", result=res)


DEBUG_SETTINGS = {"schemaVersion": "1", "logging": {"level": "debug"}}


@journey("J-CFG-06", ["CFG-010", "CFG-011"], root=True)
def j_cfg_06(t: T) -> None:
    with t.case("(a) system path"):
        sb = t.container("a")
        path = sb.system_settings(DEBUG_SETTINGS)
        p = sb.project("P")
        res = sb.view(p)
        doc = expect_success(t, res)
        check_entry(t, doc, "logging.level", res, value="debug", layer="system", source=path)
    with t.case("(b) $HOME/.config"):
        sb = t.local("b")
        sb.user_settings(DEBUG_SETTINGS)
        p = sb.project("P")
        res = sb.view(p)
        doc = expect_success(t, res)
        check_entry(t, doc, "logging.level", res, value="debug", layer="user",
                    source="~/.config/aibox/settings.json")
    with t.case("(c) absolute XDG_CONFIG_HOME outside HOME"):
        sb = t.local("c")
        xdg = f"{sb.root}/xdg"
        sb.mkdir(xdg)
        sb.user_settings(DEBUG_SETTINGS, xdg=xdg)
        p = sb.project("P")
        res = sb.view(p, env={"XDG_CONFIG_HOME": xdg})
        doc = expect_success(t, res)
        check_entry(t, doc, "logging.level", res, value="debug", layer="user",
                    source=f"{xdg}/aibox/settings.json")
    with t.case("(d) relative XDG_CONFIG_HOME"):
        sb = t.local("d")
        sb.write(f"{sb.cwd}/rel/aibox/settings.json", '{"schemaVersion":"1","logging":{"level":"debug"}}')
        p = sb.project("P")
        res = sb.view(p, env={"XDG_CONFIG_HOME": "rel"})
        doc = expect_success(t, res)
        check_entry(t, doc, "logging.level", res, value="info", layer="built_in")


J07_CASES = [
    ("mode 000", '{"schemaVersion":"1"}', 0o000, "unreadable"),
    ("text not json", NOT_JSON, 0o600, "malformed"),
    ("unknown member operatorPolicy", '{"schemaVersion":"1","operatorPolicy":{}}', 0o600, "unknown_key"),
    ("schemaVersion 2", '{"schemaVersion":"2"}', 0o600, "unsupported_version"),
    ("no schemaVersion", '{"logging":{"level":"debug"}}', 0o600, "unsupported_version"),
    ("logging.level verbose", '{"schemaVersion":"1","logging":{"level":"verbose"}}', 0o600, "invalid_value"),
    ("logging.rotationMiB 0", '{"schemaVersion":"1","logging":{"rotationMiB":0}}', 0o600, "invalid_value"),
    ("member logging twice",
     '{"schemaVersion":"1","logging":{"level":"debug"},"logging":{"level":"warn"}}', 0o600, "duplicate_key"),
    ("output null", '{"schemaVersion":"1","output":null}', 0o600, "null_value"),
]


@journey("J-CFG-07", ["CFG-006", "CFG-013", "CFG-014", "CFG-015", "CFG-016"])
def j_cfg_07(t: T) -> None:
    for label, content, mode, reason in J07_CASES:
        with t.case(label):
            sb = t.local(label)
            sb.user_settings(content, mode=mode)
            p = sb.project("P")
            res = sb.aibox_run(sb.rejection_args(p))
            expect_rejected(t, res, reason)


@journey("J-CFG-08", ["CFG-017", "CFG-018"], root=True)
def j_cfg_08(t: T) -> None:
    sb = t.container()
    sb.system_settings({"schemaVersion": "1", "logging": {"level": "warn"}})
    sb.user_settings({"schemaVersion": "1", "output": {"color": "always"}})
    other = sb.write(f"{sb.cwd}/other.json", '{"schemaVersion":"1","output":{"color":"never"}}')
    p = sb.project("P")
    with t.case("--settings-file other.json"):
        res = sb.view(p, "--settings-file", "other.json")
        doc = expect_success(t, res)
        check_entry(t, doc, "output.color", res, value="never", source=other)
        check_entry(t, doc, "logging.level", res, value="warn", layer="system",
                    source=CONTRACT.system_settings_path)
        sources = settings_sources(doc) or []
        users = [s for s in sources if isinstance(s, dict) and s.get("layer") == "user"]
        t.eq(users, [{"layer": "user", "path": other, "state": "loaded"}],
             "the user element of settingsSources is other.json", result=res)
        t.check(not any(isinstance(s, dict) and s.get("path") == "~/.config/aibox/settings.json"
                        for s in sources), "the default user file is not listed", result=res)
    with t.case("--settings-file missing.json"):
        res = sb.aibox_run(sb.rejection_args(p, "--settings-file", "missing.json"))
        expect_rejected(t, res, "missing_selected_file")


TRACE_SETTINGS = '{"schemaVersion":"1","logging":{"level":"trace"}}'


@journey("J-CFG-09", ["CFG-019", "CFG-020", "CFG-022"])
def j_cfg_09(t: T) -> None:
    sb = t.local()
    p = sb.project("P", edit=lambda s: jsonc.add_members(
        s, {"customizations": {"aibox": {"logging": {"level": "trace"}}}}))
    for d in (p, f"{p}/.devcontainer", sb.cwd):
        sb.write(f"{d}/settings.json", TRACE_SETTINGS)
    res = sb.view(p)
    doc = expect_success(t, res)
    check_entry(t, doc, "logging.level", res, value="info", layer="built_in")


@journey("J-CFG-10", ["CFG-020", "CFG-021"])
def j_cfg_10(t: T) -> None:
    with t.case("(a) View of W"):
        sb = t.local("a")
        w = sb.project("W")
        res = sb.view(w)
        doc = expect_success(t, res)
        for key, value in (("workspace.theme", "nord"), ("workspace.mode", "dark"),
                           ("workspace.layout", "dev")):
            check_entry(t, doc, key, res, value=value, layer="project", domain="workspace",
                        source=sb.tilde(p_devcontainer(w)))
    with t.case("(b) W with layout grid"):
        sb = t.local("b")
        w = sb.project("W", edit=w_with_grid)
        expect_rejected(t, sb.aibox_run(sb.rejection_args(w)), "invalid_value")
    with t.case("(c) user file workspace.theme no-such-theme"):
        sb = t.local("c")
        sb.user_settings({"schemaVersion": "1", "workspace": {"theme": "no-such-theme"}})
        p = sb.project("P")
        expect_rejected(t, sb.aibox_run(sb.rejection_args(p)), "invalid_value")
    with t.case("(d) AIBOX_WORKSPACE__MODE=sepia"):
        sb = t.local("d")
        p = sb.project("P")
        expect_rejected(t, sb.aibox_run(sb.rejection_args(p), {"AIBOX_WORKSPACE__MODE": "sepia"}),
                        "invalid_value")


def _replace_once(text_: str, old: str, new: str) -> str:
    if text_.count(old) != 1:
        raise ValueError(f"{old!r} occurs {text_.count(old)} times")
    return text_.replace(old, new)


def w_with_grid(text_: str) -> str:
    """Project W with layout grid; the edit is checked by parsing the result."""
    out = _replace_once(text_, '"mode": "dark", "layout": "dev"}', '"mode": "dark", "layout": "grid"}')
    if jsonc.loads(out)["customizations"]["aibox"]["workspace"]["layout"] != "grid":
        raise ValueError("layout grid not set")
    return out


@journey("J-CFG-11", ["CFG-030", "CFG-031"])
def j_cfg_11(t: T) -> None:
    with t.case("(a) .env in working directory and project"):
        sb = t.local("a")
        p = sb.project("P")
        sb.env_file(".env", "AIBOX_LOGGING__LEVEL=trace\n")
        sb.env_file(".env", "AIBOX_LOGGING__LEVEL=trace\n", where=p)
        res = sb.view(p)
        doc = expect_success(t, res)
        check_entry(t, doc, "logging.level", res, value="info", layer="built_in")

    def setup(label):
        sb = t.local(label)
        p = sb.project("P")
        a = sb.env_file("a.env", "AIBOX_LOGGING__LEVEL=debug\n")
        b = sb.env_file("b.env", "AIBOX_LOGGING__LEVEL=warn\n")
        return sb, p, a, b

    with t.case("(b) --env-file a.env --env-file b.env"):
        sb, p, a, b = setup("b")
        res = sb.view(p, "--env-file", "a.env", "--env-file", "b.env")
        doc = expect_success(t, res)
        check_entry(t, doc, "logging.level", res, value="warn", layer="env_file", source=b,
                    overridden=[{"layer": "env_file", "source": a, "value": "debug", "redacted": False}])
        envs = [s for s in (settings_sources(doc) or []) if isinstance(s, dict) and s.get("layer") == "env_file"]
        t.eq([s.get("path") for s in envs], [a, b], "settingsSources lists a.env then b.env", result=res)
    with t.case("(c) AIBOX_ENV_FILE=a.env"):
        sb, p, a, b = setup("c")
        res = sb.view(p, env={"AIBOX_ENV_FILE": a})
        doc = expect_success(t, res)
        check_entry(t, doc, "logging.level", res, value="debug", layer="env_file", source=a)
    with t.case("(d) AIBOX_ENV_FILE=a.env and --env-file b.env"):
        sb, p, a, b = setup("d")
        res = sb.view(p, "--env-file", "b.env", env={"AIBOX_ENV_FILE": a})
        doc = expect_success(t, res)
        check_entry(t, doc, "logging.level", res, value="warn", layer="env_file", source=b)
        t.check(not any(isinstance(s, dict) and s.get("path") == a for s in (settings_sources(doc) or [])),
                "a.env is not in settingsSources", result=res)


@journey("J-CFG-12", ["CFG-032", "CFG-033", "CFG-034", "CFG-035", "CFG-036", "CFG-038", "CFG-044"])
def j_cfg_12(t: T) -> None:
    def env_case(label, content, reason, mode=0o600):
        with t.case(label):
            sb = t.local(label)
            p = sb.project("P")
            path = sb.env_file("case.env", content, mode=mode)
            expect_rejected(t, sb.aibox_run(sb.rejection_args(p, "--env-file", path)), reason)

    with t.case("AIBOX_ENV_FILE=rel.env"):
        sb = t.local("rel")
        p = sb.project("P")
        sb.env_file("rel.env", "AIBOX_LOGGING__LEVEL=debug\n")
        expect_rejected(t, sb.aibox_run(sb.rejection_args(p), {"AIBOX_ENV_FILE": "rel.env"}),
                        "relative_env_file")
    env_case("export prefix", "export AIBOX_LOGGING__LEVEL=debug\n", "malformed")
    env_case("leading space", " AIBOX_LOGGING__LEVEL=debug\n", "malformed")
    env_case("lowercase key", "aibox_logging__level=debug\n", "malformed")
    env_case("same key twice", "AIBOX_LOGGING__LEVEL=debug\nAIBOX_LOGGING__LEVEL=warn\n", "duplicate_key")
    with t.case("--env-file path missing"):
        sb = t.local("missing")
        p = sb.project("P")
        expect_rejected(t, sb.aibox_run(sb.rejection_args(p, "--env-file", f"{sb.cwd}/missing.env")),
                        "missing_selected_file")
    with t.case("path is a directory"):
        sb = t.local("dir")
        p = sb.project("P")
        sb.mkdir(f"{sb.cwd}/dir.env")
        expect_rejected(t, sb.aibox_run(sb.rejection_args(p, "--env-file", f"{sb.cwd}/dir.env")),
                        "not_regular_file")
    env_case("file mode 000", "AIBOX_LOGGING__LEVEL=debug\n", "unreadable", mode=0o000)
    env_case("key AIBOX_NOPE", "AIBOX_NOPE=1\n", "unknown_key")
    env_case("key PATH", "PATH=/x\n", "unknown_key")
    with t.case("17 --env-file options"):
        sb = t.local("seventeen")
        p = sb.project("P")
        extra = []
        for i in range(CONTRACT.max_env_files + 1):
            extra += ["--env-file", sb.env_file(f"f{i:02d}.env", "AIBOX_LOGGING__LEVEL=debug\n")]
        expect_rejected(t, sb.aibox_run(sb.rejection_args(p, *extra)), "usage_error")


CRLF_ENV = b"\r\n# comment\r\nAIBOX_WORKSPACE__THEME=nord\r\n"


@journey("J-CFG-13", ["CFG-033"])
def j_cfg_13(t: T) -> None:
    with t.case("CRLF, blank line, comment"):
        sb = t.local("crlf")
        p = sb.project("P")
        path = sb.env_file("crlf.env", CRLF_ENV)
        res = sb.view(p, "--env-file", path)
        doc = expect_success(t, res)
        check_entry(t, doc, "workspace.theme", res, value="nord", layer="env_file")
    with t.case('quoted "nord"'):
        sb = t.local("quoted")
        p = sb.project("P")
        path = sb.env_file("quoted.env", 'AIBOX_WORKSPACE__THEME="nord"\n')
        expect_rejected(t, sb.aibox_run(sb.rejection_args(p, "--env-file", path)), "invalid_value")
    with t.case("mode 0644"):
        sb = t.local("0644")
        p = sb.project("P")
        path = sb.env_file("crlf.env", CRLF_ENV, mode=0o644)
        res = sb.view(p, "--env-file", path)
        doc = expect_success(t, res)
        check_entry(t, doc, "workspace.theme", res, value="nord", layer="env_file")


@journey("J-CFG-14", ["CFG-039", "CFG-040", "CFG-041", "CFG-042"])
def j_cfg_14(t: T) -> None:
    with t.case("(a) AIBOX_LOGGING__LEVEL=debug"):
        sb = t.local("a")
        p = sb.project("P")
        res = sb.view(p, env={"AIBOX_LOGGING__LEVEL": "debug"})
        doc = expect_success(t, res)
        check_entry(t, doc, "logging.level", res, value="debug", layer="environment",
                    source="AIBOX_LOGGING__LEVEL")
    with t.case("(b) AIBOX_LOGING__LEVEL=debug"):
        sb = t.local("b")
        p = sb.project("P")
        expect_rejected(t, sb.aibox_run(sb.rejection_args(p), {"AIBOX_LOGING__LEVEL": "debug"}),
                        "unknown_environment_name", suggestion="AIBOX_LOGGING__LEVEL")
    with t.case("(c) registered runtime identifiers"):
        sb = t.local("c")
        p = sb.project("P")
        res = sb.view(p, env={"AIBOX_TMUX_SESSION": "x", "AIBOX_TMUX_ANYTHING": "y",
                              "AIBOX_VERSION": "0.35.4", "AIBOX_LOG_LEVEL": "trace"})
        doc = expect_success(t, res)
        expect_default_entries(t, sb, doc, res, p)
    for label, name, value in (("(d)", "AIBOX_DIAGNOSTICS__OFFLINE", "TRUE"),
                               ("(e)", "AIBOX_EXECUTION__TIMEOUT_SECONDS", "+5"),
                               ("(f)", "AIBOX_EXECUTION__TIMEOUT_SECONDS", "0"),
                               ("(g)", "AIBOX_OUTPUT__COLOR", "")):
        with t.case(f"{label} {name}={value}"):
            sb = t.local(label)
            p = sb.project("P")
            expect_rejected(t, sb.aibox_run(sb.rejection_args(p), {name: value}), "invalid_value")


@journey("J-CFG-15", ["CFG-050", "CFG-051", "CFG-056", "CFG-008"])
def j_cfg_15(t: T) -> None:
    sb = t.local("flags")
    p = sb.project("P")
    for label, extra in (("--log-level debug", ["--log-level", "debug"]),
                         ("--log-level=debug", ["--log-level=debug"])):
        with t.case(f"(a) {label}"):
            res = sb.view(p, *extra)
            doc = expect_success(t, res)
            check_entry(t, doc, "logging.level", res, value="debug", layer="invocation",
                        source="--log-level")
    with t.case("(a) --timeout 5"):
        res = sb.view(p, "--timeout", "5")
        doc = expect_success(t, res)
        check_entry(t, doc, "execution.timeoutSeconds", res, value=5, layer="invocation",
                    source="--timeout")
    for label, extra, reason in (("(b) --log-level loud", ["--log-level", "loud"], "invalid_value"),
                                 ("(c) --timeout 0", ["--timeout", "0"], "invalid_value"),
                                 ("(d) --format json twice", ["--format", "json"], "usage_error"),
                                 ("(e) --set a=b", ["--set", "a=b"], "usage_error"),
                                 ("(f) positional extra", ["extra"], "usage_error")):
        with t.case(label):
            expect_rejected(t, sb.aibox_run(sb.rejection_args(p, *extra)), reason)
    with t.case("(g) aibox version --format json"):
        res = sb.aibox_run(["version", "--format", "json"])
        expect_bootstrap(t, res, "usage_error")
        t.eq(len(res.err.rstrip("\n").split("\n")) if res.err else 0, 1,
             "stderr is exactly one line", result=res)
    with t.case("(h) --log-level loud with a user file that is not JSON"):
        sb2 = t.local("h")
        sb2.user_settings(NOT_JSON)
        p2 = sb2.project("P")
        expect_rejected(t, sb2.aibox_run(sb2.rejection_args(p2, "--log-level", "loud")), "malformed")


@journey("J-CFG-16", ["CFG-052"])
def j_cfg_16(t: T) -> None:
    sb = t.local()
    d = f"{sb.root}/D"
    sb.mkdir(f"{d}/sub")
    target = sb.project("P", name="p", parent=f"{d}/sub")
    res = sb.aibox_run(["inspect", "--project", "sub/p", "--format", "json"], cwd=d)
    doc = expect_success(t, res)
    if doc is not None:
        t.eq((doc.get("target") or {}).get("workspaceRoot"), os.path.realpath(target),
             "target.workspaceRoot is the canonical path of D/sub/p", result=res)


@journey("J-CFG-17", ["CFG-007", "CFG-065"])
def j_cfg_17(t: T) -> None:
    sb = t.local()
    sb.user_settings({"schemaVersion": "1", "logging": {"level": "warn"}, "output": {"color": "never"}})
    p = sb.project("P")
    cli_res = sb.view(p)
    cli = expect_success(t, cli_res, "CLI view")
    session = MCPSession([sb.aibox, "mcp", "serve", "--context", "local"], sb.env(), sb.cwd)
    mcp_doc = None
    try:
        session.initialize()
        resp = session.call_tool("inspect_workspace", {"requestId": "j-cfg-17", "projectRoot": p,
                                                       "effectiveConfig": True})
        mcp_doc = structured(resp)
    except MCPError as exc:
        t.check(False, f"MCP session: {exc}")
    finally:
        mres = session.close()
    if not t.check(mcp_doc is not None, "inspect_workspace returned an operation result", result=mres):
        return
    t.check(not result_errors(mcp_doc), "MCP result validates against the schema", result=mres,
            actual=result_errors(mcp_doc)[:20])
    t.eq(mcp_doc.get("outcome"), "succeeded", "MCP outcome", result=mres)
    t.eq(settings_sources(mcp_doc), settings_sources(cli), "data.settingsSources equal", result=mres)
    cli_entries = {e.get("key"): e for e in settings_entries(cli)}
    mcp_entries = {e.get("key"): e for e in settings_entries(mcp_doc)}
    t.eq(list(mcp_entries), CONTRACT.key_names, "MCP settings entries in catalog order", result=mres)
    for key in CONTRACT.key_names:
        if key == "output.format":
            check_entry(t, mcp_doc, key, mres, value="json", layer="built_in", source="built-in")
        elif key == "logging.format":
            check_entry(t, mcp_doc, key, mres, value="jsonl", layer="built_in", source="built-in")
        else:
            t.eq(mcp_entries.get(key), cli_entries.get(key), f"entry({key}) equal for CLI and MCP",
                 result=mres)
    t.check(not any(e.get("layer") == "invocation" for e in configuration(mcp_doc)),
            "no MCP entry has layer invocation", result=mres)


@journey("J-CFG-18", ["CFG-054", "CFG-057"])
def j_cfg_18(t: T) -> None:
    cases = [
        ("(a) AIBOX_LOGGING__FORMAT=plain", ["mcp", "serve", "--context", "local"],
         {"AIBOX_LOGGING__FORMAT": "plain"}, None, "invalid_value"),
        ("(b) --log-format plain", ["mcp", "serve", "--log-format", "plain"], {}, None, "invalid_value"),
        ("(c) user file logging.format plain", ["mcp", "serve"], {},
         '{"schemaVersion":"1","logging":{"format":"plain"}}', "invalid_value"),
        ("(d) user file not JSON, AIBOX_OUTPUT__FORMAT=json", ["mcp", "serve"],
         {"AIBOX_OUTPUT__FORMAT": "json"}, NOT_JSON, "malformed"),
    ]
    for label, args, env, user, reason in cases:
        with t.case(label):
            sb = t.local(label)
            if user is not None:
                sb.user_settings(user)
            res = sb.aibox_run(args, env, stdin=b"")
            expect_rejected(t, res, reason, envelope=False)
            expect_bootstrap(t, res)


@journey("J-CFG-19", ["CFG-055"])
def j_cfg_19(t: T) -> None:
    sb = t.local()
    p = sb.project("P")
    session = MCPSession([sb.aibox, "mcp", "serve", "--log-level", "warn"], sb.env(), sb.cwd)
    try:
        session.initialize()
        for i, (member, value) in enumerate((("logLevel", "debug"),
                                             ("logging", {"level": "debug"}),
                                             ("outputFormat", "yaml"))):
            with t.case(f"extra input member {member}"):
                resp = session.call_tool("inspect_workspace", {"requestId": f"j-cfg-19-{i}",
                                                               "projectRoot": p, member: value})
                t.check(is_rejected(resp), "rejected by the tool's input schema", actual=resp)
        with t.case("correct call"):
            resp = session.call_tool("inspect_workspace", {"requestId": "j-cfg-19-ok", "projectRoot": p})
            doc = structured(resp)
            t.check(not is_rejected(resp) and doc is not None and doc.get("outcome") == "succeeded",
                    "the correct call succeeds", actual=resp)
    except MCPError as exc:
        t.check(False, f"MCP session: {exc}")
    finally:
        res = session.close()
    parsed = parse_jsonl(res.stderr)
    t.check(not parsed.other, "stderr is JSON Lines", result=res, actual=parsed.other[:5])
    low = [e.get("event") for e in parsed.events
           if e.get("severity") not in ("warn", "error")]
    t.eq(low, [], "stderr has no event below warn", result=res)


@journey("J-CFG-20", ["CFG-066"])
def j_cfg_20(t: T) -> None:
    sb = t.local()
    markers = f"{sb.root}/markers"
    fakebin = f"{sb.root}/fakebin"
    sb.mkdir(markers)
    sb.mkdir(fakebin)
    for tool in ("docker", "podman", "devcontainer"):
        sb.write(f"{fakebin}/{tool}", f"#!/bin/sh\n: > '{markers}/{tool}'\n", mode=0o700)
    p = sb.project("P", edit=lambda s: jsonc.add_members(s, {
        "initializeCommand": f"touch {markers}/initializeCommand",
        "postCreateCommand": f"touch {markers}/postCreateCommand"}))
    res = sb.aibox_run(sb.view_args(p), path=f"{fakebin}:/usr/bin:/bin")
    expect_success(t, res)
    t.eq(sorted(os.listdir(markers)), [], "no marker exists", result=res)


@journey("J-CFG-21", ["CFG-067"])
def j_cfg_21(t: T) -> None:
    for kind in ("P", "W"):
        with t.case(kind):
            sb = t.local(kind)
            proj = sb.project(kind)
            res = sb.aibox_run(["inspect", "--project", proj, "--format", "json"])
            doc = expect_success(t, res)
            if doc is None:
                continue
            t.eq(configuration(doc), [CONTRACT.native_entry("devcontainer.config")],
                 "data.configuration is exactly the native entries", result=res)
            t.check("settingsSources" not in (doc.get("data") or {}), "data has no settingsSources",
                    result=res)


CANARY = "canary-0123456789"


def _strings(obj, skip_path=(), path=()):
    if isinstance(obj, dict):
        for k, v in obj.items():
            yield from _strings(v, skip_path, path + (k,))
    elif isinstance(obj, list):
        for v in obj:
            yield from _strings(v, skip_path, path)
    elif isinstance(obj, str) and path != skip_path:
        yield path, obj
    if isinstance(obj, dict):
        for k in obj:
            if path + (k,) != skip_path:
                yield path + (k,), k


@journey("J-CFG-22", ["CFG-064", "RED-001", "RED-003"])
def j_cfg_22(t: T) -> None:
    sb = t.local()
    sb.user_settings({"schemaVersion": "1", "logging": {"file": "run.log"}})
    p = sb.project("P", parent=f"{sb.home}/work")
    env_path = sb.env_file(f"{CANARY}.env", "AIBOX_WORKSPACE__THEME=nord\n", where=sb.home)
    res = sb.view(p, "--env-file", env_path, env={"EXAMPLE_API_TOKEN": CANARY})
    doc = expect_success(t, res)
    check_entry(t, doc, "logging.file", res, value="~/.local/state/aibox/logs/run.log", redacted=False)
    theme = entry(doc, "workspace.theme") or {}
    src = theme.get("source") or ""
    t.check(CONTRACT.replacement in src and CANARY not in src and src.startswith("~"),
            "entry(workspace.theme).source contains [redacted], not the canary, and starts with ~",
            result=res, actual=src)
    sources = settings_sources(doc) or []
    envs = [s for s in sources if isinstance(s, dict) and s.get("layer") == "env_file"]
    t.eq(len(envs), 1, "one env_file element in settingsSources", result=res)
    for s in envs:
        path = s.get("path") or ""
        t.check(CONTRACT.replacement in path and CANARY not in path and path.startswith("~"),
                "the env_file settingsSources path contains [redacted], not the canary, and starts with ~",
                result=res, actual=path)
    for s in sources:
        path = s.get("path") if isinstance(s, dict) else None
        if isinstance(path, str):
            t.check(not path.startswith(sb.home), f"settingsSources path {path!r} uses ~ for HOME",
                    result=res)
    t.eq(next((s.get("path") for s in sources if isinstance(s, dict) and s.get("layer") == "user"), None),
         "~/.config/aibox/settings.json", "user settingsSources path starts with ~", result=res)
    t.eq(next((s.get("path") for s in sources if isinstance(s, dict) and s.get("layer") == "project"), None),
         sb.tilde(p_devcontainer(p)), "project settingsSources path starts with ~", result=res)
    if doc is not None:
        leaks = [(".".join(map(str, path)), s) for path, s in _strings(doc, ("target", "workspaceRoot"))
                 if sb.home in s]
        t.eq(leaks, [], "apart from target.workspaceRoot, stdout does not contain the HOME path",
             result=res)
    t.check(CANARY not in res.out, "stdout does not contain the canary", result=res)
    t.check(CANARY not in res.err, "stderr does not contain the canary", result=res)


@journey("J-CFG-23", ["CFG-009", "CFG-024"])
def j_cfg_23(t: T) -> None:
    for cmd in ("version", "help"):
        with t.case(f"(a) aibox {cmd}"):
            sb = t.local(cmd)
            sb.user_settings(NOT_JSON)
            res = sb.aibox_run([cmd], {"AIBOX_NOPE": "1"})
            t.eq(res.exit, 0, "exit code", result=res)
            t.check(res.out.strip() != "", f"stdout holds the {cmd} text", result=res)
            t.eq(res.err, "", "stderr is empty", result=res)
    with t.case("(b) HOME unset"):
        sb = t.local("nohome")
        p = sb.project("P")
        res = sb.aibox_run(sb.rejection_args(p), home=False)
        expect_rejected(t, res, "no_home")


@journey("J-CFG-24", ["CFG-060", "CFG-061"])
def j_cfg_24(t: T) -> None:
    sb = t.local()
    sb.user_settings({"schemaVersion": "1", "logging": {"level": "debug"}})
    demo = REPO / DEMO_PROJECT_RELATIVE
    res = sb.aibox_run(["inspect", "--context", "local", "--project", DEMO_PROJECT_RELATIVE,
                        "--effective-config", "--format", "json"], cwd=str(REPO))
    doc = expect_success(t, res)
    text_ = (demo / ".devcontainer" / "devcontainer.json").read_text()
    native = [CONTRACT.native_entry("devcontainer.config")]
    features = jsonc.feature_references(text_)
    if features:
        native.append(CONTRACT.native_entry("devcontainer.features", features))
    conf = configuration(doc)
    t.eq(conf[:len(native)], native, "the native entries for the project", result=res)
    t.eq([e.get("key") for e in conf[len(native):]], CONTRACT.key_names,
         "then 15 settings entries in catalog order", result=res)
    check_entry(t, doc, "logging.level", res, value="debug", layer="user",
                source="~/.config/aibox/settings.json")
