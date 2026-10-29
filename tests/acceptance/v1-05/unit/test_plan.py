import unittest

from actools.plan import ITEMS, census_rows, known_risk_stacks, make_plan, normalize_items
from unit.fixtures import CENSUS, REAL, synthetic

AMD, ARM = "debian-13-amd64", "debian-13-arm64"
SYNTHETIC_ITEMS = [item for item in ITEMS if item not in ("stack",)]


def plan_for(catalog, item, target=AMD, **kwargs):
    return make_plan(catalog, [target], [], items=[item], **kwargs)


def ids(plan):
    return [test.id for test in plan.tests]


def selection(plan, test):
    return plan.builds[test.build].selection_dict()


def subjects(test):
    return {e.subject for e in test.expectations}


def scopes(test):
    return {e.scope.split()[0] for e in test.expectations}


class TestIdTest(unittest.TestCase):
    """Chapter 5 test-ID slots per item."""

    def test_slot_examples_for_every_item(self):
        plan = make_plan(REAL, [ARM], CENSUS)
        by_item = {}
        for test in plan.tests:
            by_item.setdefault(test.item, set()).add(test.id)
        examples = {
            "present": "AC-TOOLS/present python/uv uv debian-13-arm64",
            "version": "AC-TOOLS/version go/go go debian-13-arm64 version=1.25.12",
            "absent": "AC-TOOLS/absent debian-13-arm64",
            "isolated": "AC-TOOLS/isolated go/goimports goimports debian-13-arm64",
            "depends": "AC-TOOLS/depends ai/aider uv debian-13-arm64",
            "bundle": "AC-TOOLS/bundle go/all debian-13-arm64",
            "side-effects": "AC-TOOLS/side-effects python/uv debian-13-arm64",
            "requires": "AC-TOOLS/requires ai/codex debian-13-arm64 missing=node/node",
            "unsupported-version": "AC-TOOLS/unsupported-version python/uv debian-13-arm64",
            "unsupported-target": "AC-TOOLS/unsupported-target python/uv fixture=debian-12",
            "conflict": "AC-TOOLS/conflict go/go debian-13-arm64 version=1.25.12",
            "smoke": "AC-TOOLS/smoke SMOKE-LATEX debian-13-arm64",
            "stack": "AC-TOOLS/stack STACK-PYTHON debian-13-arm64",
            "inspect": "AC-TOOLS/inspect example=customized",
            "trace": "AC-TOOLS/trace ADDON:python/pip",
        }
        self.assertEqual(set(examples), set(ITEMS))
        for item, example in examples.items():
            self.assertIn(example, by_item[item], item)

    def test_ids_are_unique_when_planning_both_targets(self):
        plan = make_plan(REAL, [AMD, ARM], CENSUS)
        self.assertEqual(len(ids(plan)), len(set(ids(plan))))

    def test_target_slot_is_used_except_where_chapter_five_omits_it(self):
        for test in make_plan(REAL, [ARM], CENSUS).tests:
            if test.item in ("unsupported-target", "inspect", "trace"):
                self.assertNotIn(ARM, test.id)
            else:
                self.assertIn(f" {ARM}", test.id)

    def test_unknown_item_is_rejected(self):
        with self.assertRaises(ValueError):
            normalize_items(["presence"])
        self.assertEqual(normalize_items(["AC-TOOLS/present,version"]), ["present", "version"])


class PresentAndVersionTest(unittest.TestCase):
    def test_present_covers_every_tool_of_every_non_bundle_feature(self):
        plan = plan_for(REAL, "present")
        expected = {(f, t["name"]) for f, feature in REAL.features.items()
                    if feature["role"] != "bundle" for t in feature["tools"]}
        self.assertEqual({(t.feature, t.tool) for t in plan.tests}, expected)

    def test_present_checks_default_version_and_commands(self):
        plan = plan_for(synthetic(), "present", features=["runtime"])
        rules = [(e.rule, e.regex) for e in plan.tests[0].expectations]
        self.assertIn(("version", r"(^|[^0-9.])2\.0\.0([^0-9]|$)"), rules)
        self.assertIn(("command-present", None), rules)
        self.assertEqual(selection(plan, plan.tests[0]), {"runtime": {}})

    def test_present_selection_includes_required_features(self):
        plan = plan_for(synthetic(), "present", features=["kit/a"])
        self.assertEqual(selection(plan, plan.tests[0]), {"kit/a": {}, "runtime": {}})

    def test_version_covers_each_non_default_value_with_checks_only(self):
        plan = plan_for(REAL, "version")
        expected = {(f, v) for f in REAL.features for v in REAL.non_default_versions(f)}
        self.assertEqual({(t.feature, dict(t.qualifiers)["version"]) for t in plan.tests},
                         expected)
        for test in plan.tests:
            self.assertEqual(selection(plan, test)[test.feature],
                             {"version": dict(test.qualifiers)["version"]})
            self.assertNotIn("command-present", {e.rule for e in test.expectations})


class AbsenceTest(unittest.TestCase):
    def test_absent_is_one_test_on_the_bare_base(self):
        plan = plan_for(synthetic(), "absent")
        (test,) = plan.tests
        self.assertEqual(test.id, f"AC-TOOLS/absent {AMD}")
        self.assertEqual(selection(plan, test), {})

    def test_absence_observables_follow_chapter_eight(self):
        (test,) = plan_for(synthetic(), "absent").tests
        rules = {(e.rule, e.subject) for e in test.expectations}
        self.assertIn(("command-absent", "c-cmd"), rules)
        self.assertIn(("command-check-absent", "c-cmd --help"), rules)
        self.assertIn(("package-absent", "libb"), rules)
        self.assertNotIn("version", {r for r, _ in rules})
        self.assertNotIn("texlive", {r for r, _ in rules})
        self.assertFalse({"shared", "libshared"} & {s for _, s in rules})
        self.assertNotIn("kit/d", scopes(test))

    def test_real_baseline_openssl_is_never_checked_absent(self):
        plan = make_plan(REAL, [AMD], [], items=["absent", "isolated"])
        self.assertNotIn("openssl", {e.subject for t in plan.tests for e in t.expectations})

    def test_real_rust_components_are_observable_through_command_checks(self):
        (test,) = plan_for(REAL, "absent").tests
        self.assertIn("cargo clippy --version", subjects(test))
        self.assertIn("rustfmt --version", subjects(test))


class IsolationTest(unittest.TestCase):
    def test_one_test_per_member_tool_covering_non_dependency_siblings(self):
        plan = plan_for(synthetic(), "isolated", ARM, features=["kit/b"])
        (test,) = plan.tests
        self.assertEqual(test.id, f"AC-TOOLS/isolated kit/b b {ARM}")
        self.assertEqual(scopes(test), {"kit/c", "kit/e"})  # kit/a is a dependency
        self.assertEqual(selection(plan, test), {"kit/b": {}, "runtime": {}})

    def test_also_provides_exempts_only_what_it_names_on_its_target(self):
        amd = plan_for(synthetic(), "isolated", AMD, features=["kit/a"]).tests[0]
        arm = plan_for(synthetic(), "isolated", ARM, features=["kit/a"]).tests[0]
        # command c-cmd is exempt on amd64, but kit/c's command check is not
        self.assertNotIn(("command-absent", "c-cmd"), {(e.rule, e.subject)
                                                       for e in amd.expectations})
        self.assertIn(("command-check-absent", "c-cmd --help"),
                      {(e.rule, e.subject) for e in amd.expectations})
        self.assertIn(("command-absent", "c-cmd"), {(e.rule, e.subject)
                                                    for e in arm.expectations})
        # feature kit/e is exempt entirely on amd64 only
        self.assertNotIn("kit/e", scopes(amd))
        self.assertIn("kit/e", scopes(arm))

    def test_real_rustc_exempts_cross_feature_only_on_amd64(self):
        def isolated_scopes(target):
            (test,) = plan_for(REAL, "isolated", target, features=["rust/rustc"]).tests
            return scopes(test)
        self.assertNotIn("rust/x86-64-cross", isolated_scopes(AMD))
        self.assertIn("rust/x86-64-cross", isolated_scopes(ARM))
        self.assertIn("rust/clippy", isolated_scopes(AMD))

    def test_member_selection_never_installs_a_checked_sibling(self):
        plan = plan_for(REAL, "isolated")
        for test in plan.tests:
            self.assertFalse(scopes(test) & REAL.install_closure(selection(plan, test)), test.id)


class BundleAndDependsTest(unittest.TestCase):
    def test_bundle_is_one_test_with_every_member_presence(self):
        plan = plan_for(synthetic(), "bundle")
        (test,) = plan.tests
        self.assertEqual(test.id, f"AC-TOOLS/bundle kit/all {AMD}")
        self.assertEqual(scopes(test), {"kit/a", "kit/b", "kit/c", "kit/d", "kit/e"})
        self.assertEqual(selection(plan, test), {"kit/all": {}, "runtime": {}})

    def test_depends_checks_direct_dependencies_only(self):
        plan = plan_for(synthetic(), "depends", features=["kit/b"])
        self.assertEqual({(t.feature, t.tool) for t in plan.tests if t.feature == "kit/b"},
                         {("kit/b", "a")})

    def test_real_depends_covers_every_feature_with_dependson(self):
        plan = plan_for(REAL, "depends")
        self.assertEqual({t.feature for t in plan.tests},
                         {f for f, feature in REAL.features.items() if feature["dependsOn"]})


class ErrorItemsTest(unittest.TestCase):
    def test_requires_omits_the_required_feature_and_expects_it_named(self):
        plan = plan_for(synthetic(), "requires")
        (test,) = plan.tests
        self.assertEqual(test.id, f"AC-TOOLS/requires kit/a {AMD} missing=runtime")
        self.assertEqual(selection(plan, test), {"kit/a": {}})
        self.assertEqual(test.error, {"code": "AIBOX_FEATURE_DEPENDENCY_MISSING",
                                      "features": ["kit/a"], "mustName": ["runtime"]})
        self.assertTrue(plan.builds[test.build].expect_failure)

    def test_unsupported_version_selects_the_contract_value(self):
        plan = plan_for(synthetic(), "unsupported-version")
        self.assertEqual({t.feature for t in plan.tests}, {"runtime", "kit/a"})
        test = next(t for t in plan.tests if t.feature == "kit/a")
        self.assertEqual(test.id, f"AC-TOOLS/unsupported-version kit/a {AMD}")
        self.assertEqual(selection(plan, test),
                         {"kit/a": {"version": "0.0.0-unsupported"}, "runtime": {}})

    def test_unsupported_target_runs_each_fixture_once_and_accepts_the_closure(self):
        plan = make_plan(synthetic(), [AMD, ARM], [], items=["unsupported-target"],
                         features=["kit/b"])
        self.assertEqual(sorted(ids(plan)),
                         [f"AC-TOOLS/unsupported-target kit/b fixture={u['id']}"
                          for u in sorted(REAL.data["unsupportedTargets"],
                                          key=lambda u: u["id"])])
        self.assertEqual({plan.builds[t.build].fixture for t in plan.tests},
                         {u["id"] for u in REAL.data["unsupportedTargets"]})
        self.assertEqual(plan.tests[0].error["features"], ["kit/a", "kit/b", "runtime"])

    def test_upstream_features_have_no_aibox_error_tests(self):
        for item in ("requires", "unsupported-version", "unsupported-target", "side-effects"):
            plan = make_plan(REAL, [AMD], [], items=[item])
            self.assertTrue(all(REAL.is_aibox(t.feature) for t in plan.tests), item)

    def test_conflict_selects_bundle_and_member_with_first_non_default_version(self):
        plan = plan_for(synthetic(), "conflict")
        (test,) = plan.tests
        self.assertEqual(test.id, f"AC-TOOLS/conflict kit/a {AMD} version=1.0")
        self.assertEqual(selection(plan, test),
                         {"kit/all": {}, "kit/a": {"version": "1.0"}, "runtime": {}})
        self.assertEqual(test.error["code"], "AIBOX_FEATURE_VERSION_CONFLICT")
        self.assertEqual(test.error["features"], ["kit/a"])

    def test_real_conflict_covers_members_with_several_versions(self):
        plan = plan_for(REAL, "conflict")
        self.assertEqual({t.feature for t in plan.tests},
                         {"go/go", "node/node", "python/python", "rust/rustc"})


class OtherItemsTest(unittest.TestCase):
    def test_side_effects_carry_catalog_container_env(self):
        plan = plan_for(REAL, "side-effects", features=["rust/rustc"])
        self.assertEqual(plan.tests[0].payload["containerEnv"],
                         REAL.feature("rust/rustc")["containerEnv"])

    def test_every_smoke_entry_copies_its_fixture(self):
        plan = plan_for(REAL, "smoke")
        self.assertEqual({t.feature for t in plan.tests}, {s["id"] for s in REAL.data["smoke"]})
        for test in plan.tests:
            smoke = next(s for s in REAL.data["smoke"] if s["id"] == test.feature)
            self.assertEqual(plan.builds[test.build].copies, (smoke["fixture"],))

    def test_stacks_follow_chapter_five(self):
        stacks = dict(known_risk_stacks(REAL))
        self.assertEqual(set(stacks), {"STACK-PYTHON", "STACK-NODE", "STACK-DOCS",
                                       "STACK-NATIVE", "STACK-SECURITY"})
        uv_tools = {f for f, feature in REAL.features.items()
                    if any(t["method"] == "uv-tool" for t in feature["tools"])}
        self.assertTrue(uv_tools <= set(stacks["STACK-PYTHON"]))
        plan = plan_for(REAL, "stack")
        self.assertEqual({t.feature for t in plan.tests}, set(stacks))
        for test in plan.tests:
            installed = REAL.depends_closure(selection(plan, test))
            required = {r for f in installed for r in REAL.feature(f)["requires"]}
            self.assertLessEqual(required, installed, test.id)
            self.assertEqual(scopes(test), {f for f in installed if not REAL.is_bundle(f)})

    def test_trace_has_one_test_per_census_row(self):
        plan = make_plan(REAL, [], CENSUS, items=["trace"])
        self.assertEqual(len(plan.tests), 98)

    def test_inspect_covers_customized_and_minimal(self):
        plan = make_plan(REAL, [], [], items=["inspect"])
        self.assertEqual(ids(plan), ["AC-TOOLS/inspect example=customized",
                                     "AC-TOOLS/inspect example=minimal"])

    def test_census_rows_join_recipe_and_tool(self):
        self.assertEqual(census_rows([{"id": "ADDON:x", "tools": [{"name": "y"}]}]),
                         ["ADDON:x/y"])


class FilterTest(unittest.TestCase):
    def test_feature_pattern_selects_tests_and_only_their_builds(self):
        plan = make_plan(synthetic(), [AMD], [], items=SYNTHETIC_ITEMS, features=["kit/c"])
        self.assertTrue(all("kit/c" in t.subjects for t in plan.tests))
        self.assertEqual(set(plan.builds), {t.build for t in plan.tests if t.build})

    def test_test_pattern_selects_by_full_id(self):
        plan = make_plan(REAL, [ARM], [], items=["present"], tests=["* go/go *"])
        self.assertEqual(ids(plan), ["AC-TOOLS/present go/go go debian-13-arm64"])

    def test_a_build_is_either_expected_to_fail_or_brought_up(self):
        plan = make_plan(REAL, [AMD], [])
        for test in plan.tests:
            if test.build:
                self.assertEqual(plan.builds[test.build].expect_failure, test.error is not None)


if __name__ == "__main__":
    unittest.main()
