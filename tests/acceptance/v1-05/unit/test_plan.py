import re
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


class TestIdTest(unittest.TestCase):
    def test_ids_name_item_feature_tool_and_target(self):
        plan = plan_for(REAL, "present", ARM, features=["python/uv"])
        self.assertEqual(ids(plan), ["AC-TOOLS/present python/uv uv debian-13-arm64"])

    def test_ids_are_unique_for_both_targets(self):
        plan = make_plan(REAL, [AMD, ARM], CENSUS)
        self.assertEqual(len(ids(plan)), len(set(ids(plan))))

    def test_container_tests_name_the_target_and_offline_tests_do_not(self):
        for test in make_plan(REAL, [ARM], CENSUS).tests:
            pattern = rf"^AC-TOOLS/{re.escape(test.item)} \S+"
            self.assertRegex(test.id, pattern)
            if test.build:
                self.assertIn(f" {ARM}", test.id)
            else:
                self.assertNotIn(ARM, test.id)

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
    def test_absent_uses_the_bare_base_and_skips_baseline_provides(self):
        plan = plan_for(synthetic(), "absent")
        self.assertEqual({tuple(selection(plan, t)) for t in plan.tests}, {()})
        subjects = {e.subject for t in plan.tests for e in t.expectations}
        self.assertNotIn("shared", subjects)
        self.assertNotIn("libshared", subjects)
        self.assertIn("libb", subjects)

    def test_tool_without_commands_packages_or_paths_has_no_absence_test(self):
        plan = plan_for(synthetic(), "absent")
        self.assertNotIn("kit/d", {t.feature for t in plan.tests})

    def test_real_baseline_openssl_is_never_checked_absent(self):
        plan = make_plan(REAL, [AMD], [], items=["absent", "isolated"])
        self.assertNotIn("openssl", {e.subject for t in plan.tests for e in t.expectations})


class IsolationTest(unittest.TestCase):
    def test_dependency_siblings_are_not_checked(self):
        plan = plan_for(synthetic(), "isolated", ARM, features=["kit/b"])
        self.assertEqual({dict(t.qualifiers)["sibling"] for t in plan.tests}, {"kit/c"})

    def test_also_provides_is_excluded_only_on_its_target(self):
        amd = plan_for(synthetic(), "isolated", AMD, features=["kit/a"])
        arm = plan_for(synthetic(), "isolated", ARM, features=["kit/a"])
        self.assertNotIn("kit/c", {dict(t.qualifiers)["sibling"] for t in amd.tests})
        self.assertIn("kit/c", {dict(t.qualifiers)["sibling"] for t in arm.tests})

    def test_real_rustc_also_provides_cross_gcc_only_on_amd64(self):
        def subjects(target):
            plan = plan_for(REAL, "isolated", target, features=["rust/rustc"])
            return {e.subject for t in plan.tests for e in t.expectations}
        also = {entry["name"] for tool in REAL.feature("rust/rustc")["tools"]
                for entry in tool.get("alsoProvides", [])}
        self.assertTrue(also)
        self.assertFalse(also & subjects(AMD))
        self.assertEqual(also & subjects(ARM), also)

    def test_member_selection_never_contains_a_sibling(self):
        plan = plan_for(REAL, "isolated")
        for test in plan.tests:
            self.assertNotIn(dict(test.qualifiers)["sibling"],
                             REAL.install_closure(selection(plan, test)))


class BundleAndDependsTest(unittest.TestCase):
    def test_bundle_checks_every_member_tool(self):
        plan = plan_for(synthetic(), "bundle")
        self.assertEqual({t.tool for t in plan.tests}, {"a", "b", "c", "d"})
        self.assertEqual(selection(plan, plan.tests[0]), {"kit/all": {}, "runtime": {}})

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
        self.assertEqual(selection(plan, test),
                         {"kit/a": {"version": "0.0.0-unsupported"}, "runtime": {}})

    def test_unsupported_target_runs_each_fixture_and_accepts_the_closure(self):
        plan = plan_for(synthetic(), "unsupported-target", features=["kit/b"])
        self.assertEqual({plan.builds[t.build].fixture for t in plan.tests},
                         {u["id"] for u in REAL.data["unsupportedTargets"]})
        self.assertEqual(plan.tests[0].error["features"], ["kit/a", "kit/b", "runtime"])

    def test_upstream_features_have_no_aibox_error_tests(self):
        for item in ("requires", "unsupported-version", "unsupported-target", "side-effects"):
            plan = make_plan(REAL, [AMD], [], items=[item])
            self.assertTrue(all(REAL.is_aibox(t.feature) for t in plan.tests), item)

    def test_conflict_selects_bundle_and_member_with_non_default_version(self):
        plan = plan_for(synthetic(), "conflict")
        (test,) = plan.tests
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
        uv_tools = {f for f, feature in REAL.features.items()
                    if any(t["method"] == "uv-tool" for t in feature["tools"])}
        self.assertTrue(uv_tools <= set(stacks["python/all"]))
        plan = plan_for(REAL, "stack")
        selections = {dict(t.qualifiers)["stack"]: selection(plan, t) for t in plan.tests}
        self.assertEqual(set(selections), set(stacks))
        for lead, selected in selections.items():
            installed = REAL.depends_closure(selected)
            required = {r for f in installed for r in REAL.feature(f)["requires"]}
            self.assertLessEqual(required, installed, lead)
        self.assertIn("node/node", REAL.depends_closure(selections["node/all"]))

    def test_trace_has_one_test_per_census_row(self):
        plan = make_plan(REAL, [], CENSUS, items=["trace"])
        self.assertEqual(len(plan.tests), 98)

    def test_inspect_covers_customized_and_minimal(self):
        plan = make_plan(REAL, [], [], items=["inspect"])
        self.assertEqual(ids(plan), ["AC-TOOLS/inspect customized", "AC-TOOLS/inspect minimal"])

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
