// Component tests for semantic-group bundle Features (chapter 8, Tool Feature
// interface: Shapes). Bundles are generated entirely from the catalog.
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import {join} from 'node:path';
import {describe, it} from 'node:test';

import {dependencyReference, featureDirectories, loadCatalog} from '../../scripts/v1-05-features.mjs';
import {
  aptLog, errorLines, featureCopy, featuresSrc, installedVersion, makeContainer, runFeature, snapshot,
} from './harness.mjs';

const catalog = loadCatalog();
const bundles = catalog.features.filter((f) => f.role === 'bundle');
const present = new Set(featureDirectories());

describe('bundles', () => {
  it('exist for exactly the catalog bundles of semantic groups', () => {
    const semantic = catalog.groups.filter((g) => g.shape === 'semantic').map((g) => g.bundle).sort();
    assert.deepEqual(bundles.map((f) => f.id).sort(), semantic);
    for (const bundle of bundles) assert.ok(present.has(bundle.id), `features/src/${bundle.id} is missing`);
  });

  for (const bundle of bundles) {
    describe(bundle.id, () => {
      const dir = join(featuresSrc, bundle.id);
      const metadata = JSON.parse(readFileSync(join(dir, 'devcontainer-feature.json'), 'utf8'));
      const group = catalog.groups.find((g) => g.id === bundle.group);

      it('depends on every group member at its default options and has no options', () => {
        const members = group.features.filter((id) => id !== bundle.id);
        assert.deepEqual([...bundle.dependsOn].sort(), [...members].sort());
        assert.deepEqual(metadata.dependsOn, Object.fromEntries(members.map((id) => dependencyReference(catalog, id))));
        assert.deepEqual(metadata.options, {});
        assert.equal(metadata.installsAfter, undefined);
        assert.equal(metadata.containerEnv, undefined);
      });

      it('installs nothing itself on a supported target', () => {
        const container = makeContainer({arch: 'arm64'});
        const result = runFeature(featureCopy(bundle.id), container);
        assert.equal(result.status, 0, result.output);
        assert.deepEqual(aptLog(container), []);
        const files = snapshot(container).files.filter((f) => !f.startsWith('var/lib/aibox-features'));
        assert.deepEqual(files, ['etc/', `etc/os-release ${files[1].split(' ')[1]}`, 'var/', 'var/lib/', 'var/lib/apt/', 'var/lib/apt/lists/']);
        assert.equal(installedVersion(container, bundle.id), 'none');
      });

      for (const fixture of catalog.unsupportedTargets) {
        it(`fails on ${fixture.id} with UNSUPPORTED_TARGET before any change`, () => {
          const container = makeContainer({os: fixture.id, arch: 'amd64'});
          const before = snapshot(container);
          const result = runFeature(featureCopy(bundle.id), container);
          assert.notEqual(result.status, 0);
          const lines = errorLines(result.stderr);
          assert.equal(lines.length, 1, result.output);
          assert.ok(lines[0].startsWith(`aibox-feature error AIBOX_FEATURE_UNSUPPORTED_TARGET feature=${bundle.id} target=${fixture.id}-amd64`));
          assert.deepEqual(snapshot(container), before);
        });
      }
    });
  }
});
