// Component tests for the shared install conventions in
// features/lib/aibox-feature.sh (chapter 8, Tool Feature interface: Options,
// Dependencies, Targets, Errors, Debian packages). A sample Feature uses the
// real library and real generated catalog data; only ports are faked.
import assert from 'node:assert/strict';
import {chmodSync, existsSync, readdirSync, symlinkSync, writeFileSync} from 'node:fs';
import {join} from 'node:path';
import {describe, it} from 'node:test';

import {loadCatalog, renderFeatureFiles} from '../../scripts/v1-05-features.mjs';
import {
  aptLog, errorLines, installedPackages, installedVersion, makeContainer, markInstalled,
  runFeature, serve, sha256, sha512Base64, snapshot, tempDir, writeFile,
} from './harness.mjs';

const catalog = loadCatalog();

// The sample install.sh exercises every shared helper in the order a real
// Feature uses them: guards, Debian candidates, runtime and build-only
// packages, a pinned download and its verification, then the marker.
const sampleInstall = `#!/bin/sh
set -eu
feature_dir=$(cd "\${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init
af_apt_require \${SAMPLE_RUNTIME_PACKAGES:-} \${SAMPLE_HELPER_PACKAGES:-}
[ -z "\${SAMPLE_RUNTIME_PACKAGES:-}" ] || af_apt_install $SAMPLE_RUNTIME_PACKAGES
[ -z "\${SAMPLE_HELPER_PACKAGES:-}" ] || af_helper_install $SAMPLE_HELPER_PACKAGES
[ -z "\${SAMPLE_HELPER_PACKAGES2:-}" ] || af_helper_install $SAMPLE_HELPER_PACKAGES2
if [ -n "\${SAMPLE_CLEANUP_PATH:-}" ]; then
  af_cleanup_path "$AF_ROOT/$SAMPLE_CLEANUP_PATH"
  mkdir -p "$AF_ROOT/$SAMPLE_CLEANUP_PATH"
fi
[ -z "\${SAMPLE_EXACT_PACKAGE:-}" ] || af_apt_require_version "\${SAMPLE_EXACT_PACKAGE%%=*}" "\${SAMPLE_EXACT_PACKAGE#*=}"
if [ -n "\${SAMPLE_URL:-}" ]; then
  af_download "$SAMPLE_URL" "$AF_TMP/artifact"
  case "\${SAMPLE_DIGEST:-}" in
    sha512-*) af_verify_integrity "$AF_TMP/artifact" "$SAMPLE_DIGEST" ;;
    *) af_verify_sha256 "$AF_TMP/artifact" "$SAMPLE_DIGEST" ;;
  esac
  mkdir -p "$AF_ROOT/usr/local/bin"
  cp "$AF_TMP/artifact" "$AF_ROOT/usr/local/bin/sample"
fi
printf '%s %s %s\\n' "$AF_TARGET" "$AF_ARCH" "\${AF_VERSION:-none}" > "$AF_ROOT/sample-ran"
af_finish
`;

function sampleFeature(overrides = {}) {
  const feature = {
    id: 'sample/tool',
    group: 'sample',
    role: 'tool',
    owner: 'aibox',
    featureVersion: '1.0.0',
    targets: ['debian-13-amd64', 'debian-13-arm64'],
    options: {version: {type: 'string', enum: ['1.0.0', '2.0.0'], default: '2.0.0'}},
    dependsOn: [],
    requires: [],
    tools: [{name: 'sample', docsUrl: 'https://example.invalid/', checks: [], commands: ['sample']}],
    ...overrides,
  };
  const dir = join(tempDir('sample'), 'tool');
  for (const [path, content] of Object.entries(renderFeatureFiles(feature, catalog))) {
    writeFile(join(dir, path), content);
  }
  writeFile(join(dir, 'install.sh'), sampleInstall);
  chmodSync(join(dir, 'install.sh'), 0o755);
  return dir;
}

// A PATH with the system executables except those matching `exclude`.
function pathWithout(exclude) {
  const dir = tempDir('path');
  for (const name of readdirSync('/usr/bin')) {
    if (!exclude.test(name)) symlinkSync(join('/usr/bin', name), join(dir, name));
  }
  return dir;
}

function assertError(result, code, featureId, detailPattern) {
  assert.notEqual(result.status, 0, `expected failure, output:\n${result.output}`);
  const lines = errorLines(result.stderr);
  assert.equal(lines.length, 1, `expected one error line on stderr, got:\n${result.output}`);
  const pattern = new RegExp(`^aibox-feature error ${code} feature=${featureId.replace('/', '\\/')} (.+)$`);
  const match = lines[0].match(pattern);
  assert.ok(match, `error line ${JSON.stringify(lines[0])} does not match ${pattern}`);
  if (detailPattern) assert.match(match[1], detailPattern);
}

function assertUnchanged(container, before) {
  assert.deepEqual(snapshot(container), before, 'a before-install error changed the container');
}

describe('target guard', () => {
  for (const arch of ['amd64', 'arm64']) {
    it(`accepts debian-13-${arch}`, () => {
      const container = makeContainer({arch});
      const result = runFeature(sampleFeature(), container);
      assert.equal(result.status, 0, result.output);
      assert.equal(errorLines(result.output).length, 0);
    });
  }

  const unsupported = [
    {os: 'debian-12', arch: 'amd64', target: 'debian-12-amd64'},
    {os: 'debian-12', arch: 'arm64', target: 'debian-12-arm64'},
    {os: 'ubuntu-24.04', arch: 'amd64', target: 'ubuntu-24.04-amd64'},
    {os: 'debian-13', arch: 'riscv64', target: 'debian-13-riscv64'},
    {os: 'debian-13', arch: 'i386', target: 'debian-13-i386'},
  ];
  for (const {os, arch, target} of unsupported) {
    it(`rejects ${target} before any change`, () => {
      const container = makeContainer({os, arch, available: ['xz-utils']});
      const before = snapshot(container);
      const result = runFeature(sampleFeature(), container, {
        env: {SAMPLE_HELPER_PACKAGES: 'xz-utils', SAMPLE_URL: 'https://example.invalid/a'},
      });
      assertError(result, 'AIBOX_FEATURE_UNSUPPORTED_TARGET', 'sample/tool', new RegExp(`target=${target}( |$)`));
      assertUnchanged(container, before);
      assert.deepEqual(aptLog(container), []);
    });
  }

  it('rejects a Feature whose targets omit the running architecture', () => {
    const container = makeContainer({arch: 'arm64'});
    const before = snapshot(container);
    const result = runFeature(sampleFeature({targets: ['debian-13-amd64']}), container);
    assertError(result, 'AIBOX_FEATURE_UNSUPPORTED_TARGET', 'sample/tool', /target=debian-13-arm64/);
    assertUnchanged(container, before);
  });

  it('rejects a base without dpkg', () => {
    const container = makeContainer();
    const before = snapshot(container);
    const result = runFeature(sampleFeature(), container, {env: {PATH: pathWithout(/^dpkg/)}});
    assert.notEqual(result.status, 0);
    assert.match(result.stderr, /aibox-feature error AIBOX_FEATURE_UNSUPPORTED_TARGET feature=sample\/tool /);
    assertUnchanged(container, before);
  });
});

describe('version guard', () => {
  it('uses the option default when no value is given', () => {
    const container = makeContainer();
    const result = runFeature(sampleFeature(), container);
    assert.equal(result.status, 0, result.output);
    assert.equal(installedVersion(container, 'sample/tool'), '2.0.0');
  });

  it('accepts every enum value', () => {
    for (const version of ['1.0.0', '2.0.0']) {
      const container = makeContainer();
      const result = runFeature(sampleFeature(), container, {options: {version}});
      assert.equal(result.status, 0, result.output);
      assert.equal(installedVersion(container, 'sample/tool'), version);
    }
  });

  for (const version of ['0.0.0-unsupported', 'latest', '>=1.0.0', 'https://example.invalid/x', '1.0.0 2.0.0', '*']) {
    it(`rejects ${JSON.stringify(version)} before any change`, () => {
      const container = makeContainer();
      const before = snapshot(container);
      const result = runFeature(sampleFeature(), container, {options: {version}});
      assertError(result, 'AIBOX_FEATURE_UNSUPPORTED_VERSION', 'sample/tool', /version=/);
      assertUnchanged(container, before);
    });
  }

  it('ignores a VERSION variable for a Feature without a version option', () => {
    const container = makeContainer();
    const result = runFeature(sampleFeature({options: {}}), container, {options: {version: '9.9.9'}});
    assert.equal(result.status, 0, result.output);
  });
});

describe('repeat installation', () => {
  it('succeeds without change for an identical repeat', () => {
    const container = makeContainer();
    const feature = sampleFeature();
    assert.equal(runFeature(feature, container, {options: {version: '1.0.0'}}).status, 0);
    const before = snapshot(container);
    const result = runFeature(feature, container, {options: {version: '1.0.0'}});
    assert.equal(result.status, 0, result.output);
    assert.equal(errorLines(result.output).length, 0);
    assertUnchanged(container, before);
  });

  it('fails with VERSION_CONFLICT for a different version before any change', () => {
    const container = makeContainer();
    const feature = sampleFeature();
    assert.equal(runFeature(feature, container, {options: {version: '2.0.0'}}).status, 0);
    const before = snapshot(container);
    const result = runFeature(feature, container, {options: {version: '1.0.0'}});
    assertError(result, 'AIBOX_FEATURE_VERSION_CONFLICT', 'sample/tool', /installed=2\.0\.0 requested=1\.0\.0/);
    assertUnchanged(container, before);
  });
});

describe('requires', () => {
  const withRuntime = () => sampleFeature({requires: ['sample/runtime']});

  it('fails with DEPENDENCY_MISSING naming the missing Feature before any change', () => {
    const container = makeContainer();
    const before = snapshot(container);
    const result = runFeature(withRuntime(), container);
    assertError(result, 'AIBOX_FEATURE_DEPENDENCY_MISSING', 'sample/tool', /sample\/runtime/);
    assertUnchanged(container, before);
  });

  it('accepts any installed version of the required runtime', () => {
    const container = makeContainer();
    markInstalled(container, 'sample/runtime', '7.1.0');
    const result = runFeature(withRuntime(), container);
    assert.equal(result.status, 0, result.output);
  });
});

describe('version coupling', () => {
  const equal = () => sampleFeature({versionCoupling: {feature: 'sample/base', rule: 'equal'}});
  const table = () => sampleFeature({
    versionCoupling: {feature: 'sample/base', rule: 'table', table: {'1.0.0': ['0.9.0', '0.9.1'], '2.0.0': '1.0.0'}},
  });

  it('accepts an equal coupled version', () => {
    const container = makeContainer();
    markInstalled(container, 'sample/base', '1.0.0');
    assert.equal(runFeature(equal(), container, {options: {version: '1.0.0'}}).status, 0);
  });

  it('fails with COUPLING_MISMATCH for an unequal coupled version before any change', () => {
    const container = makeContainer();
    markInstalled(container, 'sample/base', '2.0.0');
    const before = snapshot(container);
    const result = runFeature(equal(), container, {options: {version: '1.0.0'}});
    assertError(result, 'AIBOX_FEATURE_COUPLING_MISMATCH', 'sample/tool', /sample\/base.*2\.0\.0/);
    assertUnchanged(container, before);
  });

  it('fails with COUPLING_MISMATCH when the coupled Feature is not installed', () => {
    const container = makeContainer();
    const before = snapshot(container);
    const result = runFeature(equal(), container);
    assertError(result, 'AIBOX_FEATURE_COUPLING_MISMATCH', 'sample/tool', /sample\/base/);
    assertUnchanged(container, before);
  });

  it('accepts versions the table lists for the selected version', () => {
    for (const [version, coupled] of [['1.0.0', '0.9.0'], ['1.0.0', '0.9.1'], ['2.0.0', '1.0.0']]) {
      const container = makeContainer();
      markInstalled(container, 'sample/base', coupled);
      const result = runFeature(table(), container, {options: {version}});
      assert.equal(result.status, 0, `${version}/${coupled}: ${result.output}`);
    }
  });

  it('fails with COUPLING_MISMATCH when the table rule finds no coupled Feature', () => {
    const container = makeContainer();
    const before = snapshot(container);
    const result = runFeature(table(), container, {options: {version: '2.0.0'}});
    assertError(result, 'AIBOX_FEATURE_COUPLING_MISMATCH', 'sample/tool', /installed=none/);
    assertUnchanged(container, before);
  });

  it('fails with COUPLING_MISMATCH for a version the table does not list', () => {
    const container = makeContainer();
    markInstalled(container, 'sample/base', '1.0.0');
    const result = runFeature(table(), container, {options: {version: '1.0.0'}});
    assertError(result, 'AIBOX_FEATURE_COUPLING_MISMATCH', 'sample/tool', /sample\/base/);
  });
});

describe('Debian packages', () => {
  it('fails with PACKAGE_UNAVAILABLE naming the package before any change', () => {
    const container = makeContainer({available: ['gcc']});
    const before = snapshot(container);
    const result = runFeature(sampleFeature(), container, {env: {SAMPLE_RUNTIME_PACKAGES: 'gcc libc6-dev'}});
    assertError(result, 'AIBOX_FEATURE_PACKAGE_UNAVAILABLE', 'sample/tool', /package=libc6-dev/);
    assertUnchanged(container, before);
    assert.ok(aptLog(container).every((line) => !/\binstall\b/.test(line)), 'no package was installed');
  });

  it('fails with PACKAGE_UNAVAILABLE for a missing build-only helper before any change', () => {
    const container = makeContainer({available: ['gcc']});
    const before = snapshot(container);
    const result = runFeature(sampleFeature(), container, {
      env: {SAMPLE_RUNTIME_PACKAGES: 'gcc', SAMPLE_HELPER_PACKAGES: 'minisign'},
    });
    assertError(result, 'AIBOX_FEATURE_PACKAGE_UNAVAILABLE', 'sample/tool', /package=minisign/);
    assertUnchanged(container, before);
  });

  it('fails with PACKAGE_UNAVAILABLE when the package lists cannot be updated', () => {
    const container = makeContainer({available: ['gcc']});
    writeFileSync(join(container.state, 'update-fails'), '');
    const before = snapshot(container);
    const result = runFeature(sampleFeature(), container, {env: {SAMPLE_RUNTIME_PACKAGES: 'gcc'}});
    assertError(result, 'AIBOX_FEATURE_PACKAGE_UNAVAILABLE', 'sample/tool', /gcc/);
    assertUnchanged(container, before);
  });

  it('installs candidates without recommended packages and keeps runtime packages', () => {
    const container = makeContainer({available: ['gcc', 'libc6-dev']});
    const result = runFeature(sampleFeature(), container, {env: {SAMPLE_RUNTIME_PACKAGES: 'gcc libc6-dev'}});
    assert.equal(result.status, 0, result.output);
    assert.deepEqual(installedPackages(container), ['gcc', 'libc6-dev']);
    const installs = aptLog(container).filter((line) => /\binstall\b/.test(line));
    assert.ok(installs.length > 0);
    for (const line of installs) assert.match(line, /--no-install-recommends/);
  });

  it('removes build-only helpers and their new dependencies after success', () => {
    const container = makeContainer({available: ['gcc', 'xz-utils', 'minisign', 'libsodium23'], installed: ['bash']});
    writeFileSync(join(container.state, 'requires-minisign'), 'libsodium23\n');
    const result = runFeature(sampleFeature(), container, {
      env: {SAMPLE_RUNTIME_PACKAGES: 'gcc', SAMPLE_HELPER_PACKAGES: 'xz-utils minisign'},
    });
    assert.equal(result.status, 0, result.output);
    assert.deepEqual(installedPackages(container), ['bash', 'gcc']);
  });

  it('fails when a build-only helper cannot be removed', () => {
    const container = makeContainer({available: ['xz-utils']});
    writeFileSync(join(container.state, 'purge-fails'), '');
    const result = runFeature(sampleFeature(), container, {env: {SAMPLE_HELPER_PACKAGES: 'xz-utils'}});
    assert.notEqual(result.status, 0, result.output);
    assert.match(result.stderr, /cannot remove build-only packages: xz-utils/);
  });

  it('removes the helpers of every af_helper_install call', () => {
    const container = makeContainer({available: ['xz-utils', 'minisign']});
    const result = runFeature(sampleFeature(), container, {
      env: {SAMPLE_HELPER_PACKAGES: 'xz-utils', SAMPLE_HELPER_PACKAGES2: 'minisign'},
    });
    assert.equal(result.status, 0, result.output);
    assert.deepEqual(installedPackages(container), []);
  });

  it('accepts an exact package version the repositories offer', () => {
    const container = makeContainer({available: ['google-cloud-cli']});
    writeFileSync(join(container.state, 'versions-google-cloud-cli'), '592.0.0-0\n591.0.0-0\n');
    const result = runFeature(sampleFeature(), container, {env: {SAMPLE_EXACT_PACKAGE: 'google-cloud-cli=591.0.0-0'}});
    assert.equal(result.status, 0, result.output);
  });

  it('fails with PACKAGE_UNAVAILABLE for an exact version the repositories lack', () => {
    const container = makeContainer({available: ['google-cloud-cli']});
    writeFileSync(join(container.state, 'versions-google-cloud-cli'), '592.0.0-0\n591.0.0-01\n');
    const before = snapshot(container);
    const result = runFeature(sampleFeature(), container, {env: {SAMPLE_EXACT_PACKAGE: 'google-cloud-cli=591.0.0-0'}});
    assertError(result, 'AIBOX_FEATURE_PACKAGE_UNAVAILABLE', 'sample/tool', /package=google-cloud-cli=591\.0\.0-0/);
    assertUnchanged(container, before);
  });

  it('removes a registered cleanup path on success and on failure', () => {
    for (const url of [undefined, 'https://example.invalid/missing']) {
      const container = makeContainer();
      const result = runFeature(sampleFeature(), container, {
        env: {SAMPLE_CLEANUP_PATH: 'etc/apt/sources.list.d', ...(url ? {SAMPLE_URL: url, SAMPLE_DIGEST: '0'.repeat(64)} : {})},
      });
      assert.equal(result.status === 0, url === undefined, result.output);
      assert.equal(existsSync(join(container.root, 'etc/apt/sources.list.d')), false);
    }
  });

  it('keeps a helper that was installed before the Feature ran', () => {
    const container = makeContainer({available: ['xz-utils', 'minisign'], installed: ['xz-utils']});
    const result = runFeature(sampleFeature(), container, {env: {SAMPLE_HELPER_PACKAGES: 'xz-utils minisign'}});
    assert.equal(result.status, 0, result.output);
    assert.deepEqual(installedPackages(container), ['xz-utils']);
  });

  it('leaves no package lists behind when the base had none', () => {
    const container = makeContainer({available: ['gcc']});
    const result = runFeature(sampleFeature(), container, {env: {SAMPLE_RUNTIME_PACKAGES: 'gcc'}});
    assert.equal(result.status, 0, result.output);
    assert.deepEqual(readdirSync(join(container.root, 'var/lib/apt/lists')), []);
  });
});

describe('downloads and verification', () => {
  const url = 'https://example.invalid/sample-1.0.0.tar.gz';

  function artifact(content) {
    return writeFile(join(tempDir('artifact'), 'sample'), content);
  }

  it('installs a verified artifact and records the selected version', () => {
    const container = makeContainer();
    const file = artifact('sample-binary');
    serve(container, url, file);
    const result = runFeature(sampleFeature(), container, {env: {SAMPLE_URL: url, SAMPLE_DIGEST: sha256(file)}});
    assert.equal(result.status, 0, result.output);
    assert.ok(existsSync(join(container.root, 'usr/local/bin/sample')));
    assert.equal(installedVersion(container, 'sample/tool'), '2.0.0');
  });

  it('fails with DOWNLOAD_FAILED naming the URL and removes build-only helpers', () => {
    const container = makeContainer({available: ['xz-utils']});
    const result = runFeature(sampleFeature(), container, {
      env: {SAMPLE_URL: url, SAMPLE_DIGEST: '0'.repeat(64), SAMPLE_HELPER_PACKAGES: 'xz-utils'},
    });
    assertError(result, 'AIBOX_FEATURE_DOWNLOAD_FAILED', 'sample/tool', new RegExp(`url=${url.replaceAll('.', '\\.')}`));
    assert.equal(installedVersion(container, 'sample/tool'), undefined);
    assert.deepEqual(installedPackages(container), []);
    assert.deepEqual(readdirSync(join(container.root, 'var/lib/apt/lists')), []);
  });

  it('fails with VERIFICATION_FAILED for a sha256 mismatch and installs nothing', () => {
    const container = makeContainer();
    serve(container, url, artifact('tampered'));
    const result = runFeature(sampleFeature(), container, {env: {SAMPLE_URL: url, SAMPLE_DIGEST: 'a'.repeat(64)}});
    assertError(result, 'AIBOX_FEATURE_VERIFICATION_FAILED', 'sample/tool', /sha256/);
    assert.equal(existsSync(join(container.root, 'usr/local/bin/sample')), false);
    assert.equal(installedVersion(container, 'sample/tool'), undefined);
  });

  it('verifies an npm sha512 integrity value', () => {
    const container = makeContainer();
    const file = artifact('package');
    serve(container, url, file);
    const integrity = `sha512-${sha512Base64(file)}`;
    assert.equal(runFeature(sampleFeature(), container, {env: {SAMPLE_URL: url, SAMPLE_DIGEST: integrity}}).status, 0);
    const other = makeContainer();
    serve(other, url, artifact('other'));
    const result = runFeature(sampleFeature(), other, {env: {SAMPLE_URL: url, SAMPLE_DIGEST: integrity}});
    assertError(result, 'AIBOX_FEATURE_VERIFICATION_FAILED', 'sample/tool', /sha512/);
  });
});

