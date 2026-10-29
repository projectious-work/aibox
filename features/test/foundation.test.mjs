// Component tests for the wave 2 foundation Features: the guards and the
// component-verified error codes (DOWNLOAD_FAILED, VERIFICATION_FAILED,
// PACKAGE_UNAVAILABLE) through each real install.sh, for every version and
// architecture cell. Successful installation is verified by the acceptance
// tests in real containers.
import assert from 'node:assert/strict';
import {join} from 'node:path';
import {describe, it} from 'node:test';

import {loadCatalog} from '../../scripts/v1-05-features.mjs';
import {
  aptLog, errorLines, featureCopy, makeContainer, markInstalled, runFeature, serve, snapshot,
  tempDir, writeFile,
} from './harness.mjs';

const catalog = loadCatalog();
const byId = new Map(catalog.features.map((f) => [f.id, f]));

// Architecture tokens per producer; `url` gives each cell's pinned artifact.
const fromSource = (id, tokens) => (version, arch) => byId.get(id).tools[0].source
  .replaceAll('{version}', version).replaceAll('{arch}', tokens[arch]);
const features = {
  'python/uv': {url: fromSource('python/uv', {amd64: 'x86_64', arm64: 'aarch64'})},
  'python/python': {preinstalled: ['python/uv'], download: false},
  'node/node': {url: fromSource('node/node', {amd64: 'x64', arm64: 'arm64'}), packages: ['xz-utils']},
  'go/go': {url: fromSource('go/go', {amd64: 'amd64', arm64: 'arm64'})},
  'rust/rustc': {
    url: (_, arch) => `https://static.rust-lang.org/rustup/archive/1.29.2/${{amd64: 'x86_64', arm64: 'aarch64'}[arch]}-unknown-linux-gnu/rustup-init`,
    packages: ['gcc', 'libc6-dev'],
  },
  'zig/zig': {url: fromSource('zig/zig', {amd64: 'x86_64', arm64: 'aarch64'}), packages: ['xz-utils', 'minisign']},
  'latex/core': {
    url: () => 'https://ftp.tu-chemnitz.de/pub/tug/historic/systems/texlive/2025/tlnet-final/install-tl-unx.tar.gz',
    packages: ['perl', 'gnupg'],
  },
  'browser-testing/playwright': {url: (version) => `https://registry.npmjs.org/@playwright/test/-/test-${version}.tgz`},
};

function setUp(id, container) {
  const spec = features[id];
  for (const dep of [...(spec.preinstalled ?? []), ...byId.get(id).requires]) {
    markInstalled(container, dep, byId.get(dep).options.version?.default ?? 'none');
  }
}

function assertError(result, code, id, detail) {
  assert.notEqual(result.status, 0, result.output);
  const expected = `aibox-feature error ${code} feature=${id} `;
  const lines = errorLines(result.stderr);
  assert.equal(lines.length, 1, result.output);
  assert.ok(lines[0].startsWith(expected), `${JSON.stringify(lines[0])} does not start with ${JSON.stringify(expected)}`);
  if (detail) assert.ok(lines[0].includes(detail), `${JSON.stringify(lines[0])} does not name ${detail}`);
}

for (const [id, spec] of Object.entries(features)) {
  const feature = byId.get(id);
  const versions = feature.options.version.enum;
  const available = spec.packages ?? [];

  describe(id, () => {
    for (const fixture of catalog.unsupportedTargets) {
      const os = fixture.id;
      for (const arch of ['amd64', 'arm64']) {
        it(`fails on ${os} ${arch} with UNSUPPORTED_TARGET before any change`, () => {
          const container = makeContainer({os, arch, available});
          setUp(id, container);
          const before = snapshot(container);
          const result = runFeature(featureCopy(id), container);
          assertError(result, 'AIBOX_FEATURE_UNSUPPORTED_TARGET', id, `target=${os}-${arch}`);
          assert.deepEqual(snapshot(container), before);
          assert.deepEqual(aptLog(container), []);
        });
      }
    }

    it('fails for version 0.0.0-unsupported with UNSUPPORTED_VERSION before any change', () => {
      const container = makeContainer({available});
      setUp(id, container);
      const before = snapshot(container);
      const result = runFeature(featureCopy(id), container, {options: {version: '0.0.0-unsupported'}});
      assertError(result, 'AIBOX_FEATURE_UNSUPPORTED_VERSION', id, 'version=0.0.0-unsupported');
      assert.deepEqual(snapshot(container), before);
    });

    it('fails with VERSION_CONFLICT when another version is installed', () => {
      const container = makeContainer({available});
      setUp(id, container);
      markInstalled(container, id, '0.0.1-other');
      const before = snapshot(container);
      const result = runFeature(featureCopy(id), container);
      assertError(result, 'AIBOX_FEATURE_VERSION_CONFLICT', id, 'installed=0.0.1-other');
      assert.deepEqual(snapshot(container), before);
    });

    for (const dep of feature.requires) {
      it(`fails with DEPENDENCY_MISSING naming ${dep} before any change`, () => {
        const container = makeContainer({available});
        for (const pre of spec.preinstalled ?? []) markInstalled(container, pre, '0');
        const before = snapshot(container);
        const result = runFeature(featureCopy(id), container);
        assertError(result, 'AIBOX_FEATURE_DEPENDENCY_MISSING', id, dep);
        assert.deepEqual(snapshot(container), before);
      });
    }

    if (available.length > 0) {
      for (const missing of available) {
        it(`fails with PACKAGE_UNAVAILABLE for ${missing} before any change`, () => {
          const container = makeContainer({available: available.filter((p) => p !== missing)});
          setUp(id, container);
          const before = snapshot(container);
          const result = runFeature(featureCopy(id), container);
          assertError(result, 'AIBOX_FEATURE_PACKAGE_UNAVAILABLE', id, `package=${missing}`);
          assert.deepEqual(snapshot(container), before);
          assert.ok(aptLog(container).every((line) => !/\binstall\b/.test(line)));
        });
      }
    }

    if (spec.download === false) return;

    for (const version of versions) {
      for (const arch of ['amd64', 'arm64']) {
        const url = spec.url(version, arch);

        it(`${version} ${arch}: fails with DOWNLOAD_FAILED naming ${url}`, () => {
          const container = makeContainer({arch, available});
          setUp(id, container);
          const before = snapshot(container);
          const result = runFeature(featureCopy(id), container, {options: {version}});
          assertError(result, 'AIBOX_FEATURE_DOWNLOAD_FAILED', id, `url=${url}`);
          assert.deepEqual(snapshot(container), before);
        });

        it(`${version} ${arch}: fails with VERIFICATION_FAILED for tampered content`, () => {
          const container = makeContainer({arch, available});
          setUp(id, container);
          serve(container, url, writeFile(join(tempDir('tampered'), 'artifact'), 'tampered artifact'));
          const before = snapshot(container);
          const result = runFeature(featureCopy(id), container, {options: {version}});
          assertError(result, 'AIBOX_FEATURE_VERIFICATION_FAILED', id);
          assert.deepEqual(snapshot(container), before);
        });
      }
    }
  });
}

describe('python/python and uv', () => {
  it('fails with DOWNLOAD_FAILED when uv cannot install the Python version', () => {
    const container = makeContainer();
    markInstalled(container, 'python/uv', '0.12.5');
    writeFile(join(container.root, 'usr/local/bin/uv'), '#!/bin/sh\necho "error: download failed" >&2\nexit 2\n', 0o755);
    const before = snapshot(container);
    const result = runFeature(featureCopy('python/python'), container, {options: {version: '3.13.15'}});
    assertError(result, 'AIBOX_FEATURE_DOWNLOAD_FAILED', 'python/python', 'version=3.13.15');
    assert.deepEqual(snapshot(container), before);
  });
});
