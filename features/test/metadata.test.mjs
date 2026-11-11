// Component tests for generated Feature metadata and its drift check
// (chapter 8, Tool Feature interface: Identity, Options, Dependencies,
// Side effects). The generator and the catalog are real collaborators.
import assert from 'node:assert/strict';
import {spawnSync} from 'node:child_process';
import {cpSync, mkdirSync, readFileSync, statSync, writeFileSync} from 'node:fs';
import {join} from 'node:path';
import {describe, it} from 'node:test';

import {featureDirectories, loadCatalog, renderMetadata} from '../../scripts/v1-05-features.mjs';
import {featuresSrc, repoRoot, tempDir} from './harness.mjs';

const catalog = loadCatalog();
const byId = new Map(catalog.features.map((f) => [f.id, f]));
const foundation = [
  'python/uv', 'python/python', 'node/node', 'go/go', 'rust/rustc', 'zig/zig', 'latex/core',
  'browser-testing/playwright',
];
const forbiddenKeys = [
  'onCreateCommand', 'updateContentCommand', 'postCreateCommand', 'postStartCommand', 'postAttachCommand',
  'mounts', 'privileged', 'capAdd', 'securityOpt', 'init', 'entrypoint', 'customizations',
];

function generator(args, featuresDir) {
  const extra = featuresDir ? ['--features-dir', featuresDir] : [];
  return spawnSync('node', [join(repoRoot, 'scripts/v1-05-features.mjs'), ...args, ...extra], {encoding: 'utf8'});
}

describe('Feature sources', () => {
  it('include every wave 2 foundation Feature', () => {
    const present = featureDirectories();
    for (const id of foundation) assert.ok(present.includes(id), `features/src/${id} is missing`);
  });

  for (const id of featureDirectories()) {
    describe(id, () => {
      const feature = byId.get(id);
      const dir = join(featuresSrc, id);
      const metadata = JSON.parse(readFileSync(join(dir, 'devcontainer-feature.json'), 'utf8'));

      it('is an aibox-owned catalog Feature', () => {
        assert.ok(feature, `${id} is not in the catalog`);
        assert.equal(feature.owner, 'aibox');
      });

      it('has the catalog identity', () => {
        assert.equal(metadata.id, id.split('/').at(-1));
        assert.equal(metadata.version, feature.featureVersion);
      });

      it('declares exactly the catalog options', () => {
        assert.deepEqual(Object.keys(metadata.options).sort(), Object.keys(feature.options).sort());
        if (feature.options.version) {
          assert.equal(metadata.options.version.type, 'string');
          assert.deepEqual(metadata.options.version.enum, feature.options.version.enum);
          assert.equal(metadata.options.version.default, feature.options.version.default);
        }
        assert.equal('enabled' in metadata.options, false);
      });

      it('maps dependsOn to <registry>/<id>:<major> and requires to installsAfter', () => {
        const expected = Object.fromEntries(feature.dependsOn.map((dep) => {
          const target = byId.get(dep);
          return target.owner === 'upstream'
            ? [target.upstream.ref, target.upstream.options ?? {}]
            : [`${catalog.registry}/${dep}:${target.featureVersion.split('.')[0]}`, {}];
        }));
        assert.deepEqual(metadata.dependsOn ?? {}, expected);
        assert.deepEqual(metadata.installsAfter ?? [], feature.requires.map((dep) => `${catalog.registry}/${dep}`));
      });

      it('has exactly the catalog containerEnv and no lifecycle, mount or privilege keys', () => {
        assert.deepEqual(metadata.containerEnv, feature.containerEnv);
        for (const key of forbiddenKeys) assert.equal(key in metadata, false, `${key} must not be declared`);
      });

      it('has an executable install.sh that uses the shared library copy', () => {
        const install = join(dir, 'install.sh');
        assert.ok(statSync(install).mode & 0o111, 'install.sh is not executable');
        const source = readFileSync(install, 'utf8');
        assert.match(source, /^#!\/bin\/sh\n/);
        assert.match(source, /\. "\$feature_dir\/lib\/aibox-feature\.sh"/);
        assert.match(source, /\naf_init\n/);
        assert.match(source, /\naf_finish\n/);
        assert.equal(
          readFileSync(join(dir, 'lib/aibox-feature.sh'), 'utf8'),
          readFileSync(join(repoRoot, 'features/lib/aibox-feature.sh'), 'utf8'),
        );
      });
    });
  }
});

describe('drift check', () => {
  function copyOfSources() {
    const dir = join(tempDir('sources'), 'src');
    cpSync(featuresSrc, dir, {recursive: true});
    return dir;
  }

  it('passes for the committed sources', () => {
    const result = generator(['check']);
    assert.equal(result.status, 0, result.stderr);
  });

  it('fails when a shared library copy drifts', () => {
    const dir = copyOfSources();
    writeFileSync(join(dir, 'python/uv/lib/aibox-feature.sh'), '# edited\n', {flag: 'a'});
    const result = generator(['check'], dir);
    assert.notEqual(result.status, 0);
    assert.match(result.stderr, /python\/uv\/lib\/aibox-feature\.sh differs/);
  });

  it('fails when metadata drifts from the catalog', () => {
    const dir = copyOfSources();
    const file = join(dir, 'node/node/devcontainer-feature.json');
    const metadata = JSON.parse(readFileSync(file, 'utf8'));
    metadata.options.version.enum.push('27.0.0');
    writeFileSync(file, JSON.stringify(metadata, null, 2));
    const result = generator(['check'], dir);
    assert.notEqual(result.status, 0);
    assert.match(result.stderr, /node\/node\/devcontainer-feature\.json differs/);
  });

  it('fails for a Feature directory that is not an aibox-owned catalog Feature', () => {
    const dir = copyOfSources();
    mkdirSync(join(dir, 'docs/hugo'), {recursive: true});
    writeFileSync(join(dir, 'docs/hugo/install.sh'), '#!/bin/sh\n');
    const result = generator(['check'], dir);
    assert.notEqual(result.status, 0);
    assert.match(result.stderr, /docs\/hugo is not an aibox-owned catalog Feature/);
  });

  it('generate repairs drift', () => {
    const dir = copyOfSources();
    writeFileSync(join(dir, 'go/go/lib/catalog.sh'), '');
    assert.equal(generator(['generate'], dir).status, 0);
    assert.equal(generator(['check'], dir).status, 0);
  });

  it('--complete names every aibox-owned Feature without sources', () => {
    const dir = copyOfSources();
    const result = generator(['check', '--complete'], dir);
    const present = new Set(featureDirectories(dir));
    const missing = catalog.features.filter((f) => f.owner === 'aibox' && !present.has(f.id));
    assert.equal(result.status === 0, missing.length === 0);
    for (const f of missing) assert.match(result.stderr, new RegExp(`Feature ${f.id} has no sources`));
  });
});

describe('dependsOn rendering', () => {
  it('uses upstream.ref and upstream.options for an upstream Feature and <registry>/<id>:<major> otherwise', () => {
    const bundle = catalog.features.find((f) => f.dependsOn.some((d) => byId.get(d).owner === 'upstream'));
    assert.ok(bundle, 'the catalog has a Feature that depends on an upstream Feature');
    const metadata = JSON.parse(renderMetadata(bundle, catalog));
    for (const dep of bundle.dependsOn) {
      const target = byId.get(dep);
      if (target.owner === 'upstream') {
        assert.deepEqual(metadata.dependsOn[target.upstream.ref], target.upstream.options ?? {});
      } else {
        assert.deepEqual(metadata.dependsOn[`${catalog.registry}/${dep}:${target.featureVersion.split('.')[0]}`], {});
      }
    }
    assert.equal(Object.keys(metadata.dependsOn).length, bundle.dependsOn.length);
  });
});
