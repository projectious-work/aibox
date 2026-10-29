import subprocess
import unittest

from actools import checks
from actools.checks import Outcome, Probe
from unit.fixtures import REAL

VERSION_REGEX = "(^|[^0-9.]){version}([^0-9]|$)"


def version_check(version, regex=VERSION_REGEX):
    return checks.check_expectation({"kind": "version", "argv": ["t", "--version"],
                                     "regex": regex}, version, "t")


class VersionCheckTest(unittest.TestCase):
    def test_selected_version_matches_as_a_token(self):
        expectation = version_check("1.2.3")
        self.assertTrue(checks.passes(expectation, Outcome(0, "go version go1.2.3 linux/arm64")))
        self.assertTrue(checks.passes(expectation, Outcome(0, "tool\n1.2.3\n")))

    def test_dots_are_literal_and_neighbouring_digits_reject(self):
        expectation = version_check("1.2.3")
        for output in ("1x2y3", "11.2.3", "1.2.34", "0.1.2.3"):
            self.assertFalse(checks.passes(expectation, Outcome(0, output)), output)

    def test_nonzero_exit_fails_even_when_output_matches(self):
        self.assertFalse(checks.passes(version_check("1.2.3"), Outcome(1, "1.2.3")))

    def test_caret_matches_at_line_start(self):
        expectation = checks.check_expectation(
            {"kind": "command", "argv": ["python", "-m", "pip", "--version"],
             "regex": "^pip [0-9]"}, None, "python")
        self.assertTrue(checks.passes(expectation, Outcome(0, "warning\npip 25.1 from x")))
        self.assertFalse(checks.passes(expectation, Outcome(0, "the pip 25.1")))

    def test_every_catalog_regex_compiles_after_substitution(self):
        for feature_id, feature in REAL.features.items():
            for tool in feature["tools"]:
                version = REAL.selected_version(feature_id, tool)
                for expectation in checks.presence_expectations(tool, version):
                    if expectation.regex is not None:
                        checks.compile_ecmascript(expectation.regex)

    def test_missing_version_for_placeholder_is_a_planning_error(self):
        with self.assertRaises(ValueError):
            version_check(None)


class OtherKindsTest(unittest.TestCase):
    def test_package_requires_exact_version(self):
        expectation = checks.check_expectation(
            {"kind": "package", "package": "gcc", "version": "4:14.2.0-1"}, None, "gcc")
        self.assertEqual(expectation.probe.argv,
                         ("dpkg-query", "-W", "-f", "${Version}", "gcc"))
        self.assertTrue(checks.passes(expectation, Outcome(0, "4:14.2.0-1")))
        self.assertFalse(checks.passes(expectation, Outcome(0, "4:14.2.0-1+b1")))
        self.assertFalse(checks.passes(expectation, Outcome(1, "4:14.2.0-1")))

    def test_texlive_requires_the_package_name_on_a_line(self):
        expectation = checks.check_expectation({"kind": "texlive", "package": "amsmath"},
                                               None, "tex")
        self.assertTrue(checks.passes(expectation, Outcome(0, "amsmath\n")))
        self.assertFalse(checks.passes(expectation, Outcome(0, "")))
        self.assertFalse(checks.passes(expectation, Outcome(0, "amsmath-doc")))

    def test_path_and_command_presence(self):
        path = checks.check_expectation({"kind": "path", "path": "/opt/x"}, None, "x")
        self.assertEqual(path.probe.argv, ("test", "-d", "/opt/x"))
        self.assertTrue(checks.passes(path, Outcome(0, "")))
        (command,) = [e for e in checks.presence_expectations(
            {"name": "x", "checks": [], "commands": ["x"]}, None)]
        self.assertTrue(checks.passes(command, Outcome(0, "/usr/local/bin/x")))
        self.assertFalse(checks.passes(command, Outcome(1, "")))

    def test_absence_rules(self):
        tool = {"name": "x", "commands": ["x", "shared"],
                "checks": [{"kind": "package", "package": "libx", "version": "1"},
                           {"kind": "path", "path": "/opt/x"},
                           {"kind": "texlive", "package": "x"}]}
        expectations = checks.absence_expectations(tool, skip_commands={"shared"})
        self.assertEqual([(e.rule, e.subject) for e in expectations],
                         [("command-absent", "x"), ("package-absent", "libx"),
                          ("path-absent", "/opt/x")])
        command, package, path = expectations
        self.assertTrue(checks.passes(command, Outcome(1, "")))
        self.assertFalse(checks.passes(command, Outcome(0, "/usr/bin/x")))
        self.assertTrue(checks.passes(package, Outcome(1, "dpkg-query: no packages found")))
        self.assertTrue(checks.passes(package, Outcome(0, "config-files")))
        self.assertFalse(checks.passes(package, Outcome(0, "installed")))
        self.assertTrue(checks.passes(path, Outcome(1, "")))

    def test_command_check_is_absent_unless_it_would_pass(self):
        tool = {"name": "clippy", "commands": [],
                "checks": [{"kind": "command", "argv": ["cargo", "clippy", "--version"],
                            "regex": "^clippy "},
                           {"kind": "version", "argv": ["x"], "regex": "{version}"}]}
        (expectation,) = checks.absence_expectations(tool)
        self.assertEqual(expectation.rule, "command-check-absent")
        self.assertTrue(checks.passes(expectation, Outcome(101, "error: no such command")))
        self.assertTrue(checks.passes(expectation, Outcome(0, "cargo-clippy missing")))
        self.assertFalse(checks.passes(expectation, Outcome(0, "clippy 0.1.98")))


class ProbeScriptTest(unittest.TestCase):
    def test_script_round_trip_merges_streams_and_keeps_exit_codes(self):
        probes = [Probe(("sh", "-c", "echo out; echo err >&2; exit 3")),
                  checks.command_probe("sh"),
                  checks.command_probe("aibox-definitely-missing"),
                  Probe(("printf", "%s", "it's ${HOME} $(x) \"q\"")),
                  checks.path_probe("/")]
        completed = subprocess.run(["sh", "-c", checks.probe_script(probes)], text=True,
                                   capture_output=True, check=True)
        outcomes = checks.parse_probe_output("noise\n" + completed.stdout, probes)
        self.assertEqual(outcomes[probes[0]], Outcome(3, "out\nerr"))
        self.assertEqual(outcomes[probes[1]].exit_code, 0)
        self.assertNotEqual(outcomes[probes[2]].exit_code, 0)
        self.assertEqual(outcomes[probes[3]].output, "it's ${HOME} $(x) \"q\"")
        self.assertEqual(outcomes[probes[4]].exit_code, 0)

    def test_incomplete_output_is_an_error(self):
        with self.assertRaises(RuntimeError):
            checks.parse_probe_output("", [Probe(("true",))])


class ErrorLineTest(unittest.TestCase):
    CODE = "AIBOX_FEATURE_DEPENDENCY_MISSING"

    def test_line_with_build_prefix_matches(self):
        output = ("#8 0.31 step\n#8 0.32 aibox-feature error AIBOX_FEATURE_DEPENDENCY_MISSING "
                  "feature=ai/codex requires node/node\nERROR: failed")
        self.assertTrue(checks.error_line_matches(output, self.CODE, ["ai/codex"], ["node/node"]))

    def test_other_feature_or_code_does_not_match(self):
        output = "aibox-feature error AIBOX_FEATURE_DEPENDENCY_MISSING feature=go/goimports x"
        self.assertFalse(checks.error_line_matches(output, self.CODE, ["go/go"]))
        self.assertFalse(checks.error_line_matches(
            output, "AIBOX_FEATURE_UNSUPPORTED_TARGET", ["go/goimports"]))

    def test_detail_must_name_the_missing_feature(self):
        output = "aibox-feature error AIBOX_FEATURE_DEPENDENCY_MISSING feature=ai/codex missing"
        self.assertFalse(checks.error_line_matches(output, self.CODE, ["ai/codex"], ["node/node"]))
        registry_form = output + " localhost:5000/aibox-features/node/node"
        self.assertTrue(checks.error_line_matches(registry_form, self.CODE, ["ai/codex"],
                                                  ["node/node"]))
        self.assertFalse(checks.error_line_matches(output + " node/nodex", self.CODE,
                                                   ["ai/codex"], ["node/node"]))

    def test_any_feature_of_the_closure_is_accepted(self):
        output = "aibox-feature error AIBOX_FEATURE_UNSUPPORTED_TARGET feature=node/node debian 12"
        self.assertTrue(checks.error_line_matches(output, "AIBOX_FEATURE_UNSUPPORTED_TARGET",
                                                  ["ai/codex", "node/node"]))


if __name__ == "__main__":
    unittest.main()
