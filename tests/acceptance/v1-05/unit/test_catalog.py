import unittest

from unit.fixtures import REAL, synthetic


class SelectionExpansionTest(unittest.TestCase):
    def test_requires_of_a_selected_feature_are_added_at_defaults(self):
        catalog = synthetic()
        self.assertEqual(catalog.with_requires({"kit/a": {"version": "1.0"}}),
                         {"kit/a": {"version": "1.0"}, "runtime": {}})

    def test_requires_of_a_dependson_feature_are_added(self):
        # kit/b dependsOn kit/a, which requires runtime (contract question Q1).
        self.assertEqual(synthetic().with_requires({"kit/b": {}}),
                         {"kit/b": {}, "runtime": {}})

    def test_explicit_required_feature_keeps_its_options(self):
        selection = synthetic().with_requires({"kit/a": {}, "runtime": {"version": "1.0.0"}})
        self.assertEqual(selection["runtime"], {"version": "1.0.0"})

    def test_feature_without_requires_is_unchanged(self):
        self.assertEqual(REAL.with_requires({"python/uv": {}}), {"python/uv": {}})

    def test_real_bundle_whose_member_requires_a_runtime_gets_it(self):
        self.assertIn("node/node", REAL.with_requires({"typescript/all": {}}))
        self.assertIn("node/node", REAL.with_requires({"browser-testing/chromium": {}}))

    def test_install_closure_follows_dependson_and_requires(self):
        self.assertEqual(synthetic().install_closure(["kit/b"]), {"kit/b", "kit/a", "runtime"})


class ReferenceTest(unittest.TestCase):
    def test_aibox_feature_is_referenced_by_major_tag_in_the_test_registry(self):
        self.assertEqual(synthetic().feature_reference("kit/a"),
                         "localhost:5000/aibox-features/kit/a:2")

    def test_upstream_feature_is_referenced_by_digest_with_catalog_options(self):
        catalog = synthetic()
        features = catalog.devcontainer_features({"up/tool": {}})
        ref = catalog.feature("up/tool")["upstream"]["ref"]
        self.assertEqual(features, {ref: {"version": "9.9.9"}})

    def test_upstream_feature_rejects_extra_options(self):
        with self.assertRaises(ValueError):
            synthetic().devcontainer_features({"up/tool": {"version": "1"}})


class SelectedVersionTest(unittest.TestCase):
    def test_option_value_wins_then_default_then_tool_version(self):
        catalog = synthetic()
        runtime_tool = catalog.feature("runtime")["tools"][0]
        b_tool = catalog.feature("kit/b")["tools"][0]
        self.assertEqual(catalog.selected_version("runtime", runtime_tool,
                                                  {"version": "1.0.0"}), "1.0.0")
        self.assertEqual(catalog.selected_version("runtime", runtime_tool), "2.0.0")
        self.assertEqual(catalog.selected_version("kit/b", b_tool), "3.0")

    def test_non_default_versions(self):
        self.assertEqual(synthetic().non_default_versions("runtime"), ["1.0.0"])
        self.assertEqual(synthetic().non_default_versions("kit/c"), [])

    def test_every_real_version_check_has_a_selected_version(self):
        for feature_id, feature in REAL.features.items():
            for tool in feature["tools"]:
                if any(check["kind"] == "version" for check in tool["checks"]):
                    self.assertIsNotNone(REAL.selected_version(feature_id, tool),
                                         f"{feature_id} {tool['name']}")


class TargetTest(unittest.TestCase):
    def test_architecture_maps_to_catalog_target(self):
        self.assertEqual(REAL.target_for_architecture("arm64"), "debian-13-arm64")
        self.assertEqual(REAL.target_for_architecture("amd64"), "debian-13-amd64")
        with self.assertRaises(KeyError):
            REAL.target_for_architecture("riscv64")


if __name__ == "__main__":
    unittest.main()
