// Component tests for the wave 3 tool and atomic Features, driven by the
// catalog: guards for every Feature, then the component-verified error codes
// at each Feature's first external input, and the success path where the
// installer is a replaceable port (APT, uv, go, rustup, cargo, tlmgr,
// playwright). Real installations are verified by the acceptance tests.
import assert from 'node:assert/strict';
import {existsSync, readFileSync} from 'node:fs';
import {join} from 'node:path';
import {describe, it} from 'node:test';

import {featureDirectories, loadCatalog} from '../../scripts/v1-05-features.mjs';
import {
  aptLog, errorLines, featureCopy, installedPackages, installedVersion, makeContainer, markInstalled,
  runFeature, serve, snapshot, tempDir, writeFile,
} from './harness.mjs';

const catalog = loadCatalog();
const byId = new Map(catalog.features.map((f) => [f.id, f]));
const foundation = new Set([
  'python/uv', 'python/python', 'node/node', 'go/go', 'rust/rustc', 'zig/zig', 'latex/core',
  'browser-testing/playwright',
]);
const features = featureDirectories().filter((id) => !foundation.has(id) && byId.get(id).role !== 'bundle');

// Producer architecture tokens for the catalog's {arch} placeholder.
const releaseTokens = {
  'ai/claude': ['x64', 'arm64'], 'ai/copilot': ['x64', 'arm64'], 'ai/opencode': ['x64', 'arm64'],
  'docs/mdbook': ['x86_64', 'aarch64'], typst: ['x86_64', 'aarch64'], 'cloud/cloudflared': ['amd64', 'arm64'],
  'go/golangci-lint': ['amd64', 'arm64'], 'go/gosec': ['amd64', 'arm64'], 'go/goreleaser': ['x86_64', 'arm64'],
  'node/bun': ['x64', 'aarch64'], 'typescript/biome': ['x64', 'arm64'], 'typescript/oxlint': ['x86_64', 'aarch64'],
  'zig/zls': ['x86_64', 'aarch64'], 'diagramming/d2': ['amd64', 'arm64'], 'git-ui/lazygit': ['x86_64', 'arm64'],
  'infrastructure/opentofu': ['amd64', 'arm64'], 'infrastructure/packer': ['amd64', 'arm64'],
  'kubernetes/kubectl': ['amd64', 'arm64'], 'kubernetes/helm': ['amd64', 'arm64'],
  'kubernetes/kustomize': ['amd64', 'arm64'], 'kubernetes/k9s': ['amd64', 'arm64'],
  'preview/archive/resvg': ['x86_64', 'aarch64'], 'release/shellcheck': ['x86_64', 'aarch64'],
  'release/hadolint': ['x86_64', 'arm64'], 'security/crypto/age': ['amd64', 'arm64'],
  'security/crypto/sops': ['amd64', 'arm64'], 'security/crypto/minisign': ['x86_64', 'aarch64'],
  'security/scan/gitleaks': ['x64', 'arm64'], 'security/scan/osv-scanner': ['amd64', 'arm64'],
  'security/scan/syft': ['amd64', 'arm64'], 'security/scan/grype': ['amd64', 'arm64'],
  'security/scan/cosign': ['amd64', 'arm64'],
  'preview/enhanced/glow': ['x86_64', 'arm64'], 'preview/enhanced/bat': ['x86_64', 'aarch64'],
};
// First npm package each npm Feature downloads.
const firstNpmPackage = {
  'ai/codex': ['@openai/codex', 'codex'], 'ai/continue': ['@continuedev/cli', 'cli'],
  'ai/gemini': ['@google/gemini-cli', 'gemini-cli'], 'node/pnpm': ['pnpm', 'pnpm'],
  'node/yarn': ['@yarnpkg/cli-dist', 'cli-dist'], 'typescript/tsc': ['typescript', 'typescript'],
  'browser-testing/axe-playwright': ['@axe-core/playwright', 'playwright'],
  'data/visualization': ['vega-cli', 'vega-cli', '6.4.0'],
  mermaid: ['puppeteer', 'puppeteer', '25.11.0'],
};
const archiveHelpers = (id) => [
  ...(byId.get(id).tools[0].source.split(' ')[0].endsWith('.tar.xz') ? ['xz-utils'] : []),
  ...(id === 'zig/zls' ? ['minisign'] : []),
];

const method = (id) => byId.get(id).tools[0].method;
const version = (id) => byId.get(id).options.version?.default;

function debianPackages(id) {
  return byId.get(id).tools.flatMap((t) => t.packages
    ?? t.checks.filter((c) => c.kind === 'package').map((c) => c.package)
    .concat(t.checks.some((c) => c.kind === 'package') ? [] : [t.source.split(' ').at(-1).split('=')[0]]));
}

function available(id) {
  if (method(id) === 'debian-package') return debianPackages(id);
  if (method(id) === 'texlive-installer') return ['gnupg', ...(id === 'latex/music' ? ['lilypond'] : [])];
  if (method(id) === 'producer-release') return archiveHelpers(id);
  return [];
}

function fakeTool(container, path, {exit = 0, script = ''} = {}) {
  const name = path.split('/').at(-1);
  return writeFile(join(container.root, path),
    `#!/bin/sh\nprintf '%s\\n' "$*" >> "$FAKE_STATE/${name}.log"\n${script}\nexit ${exit}\n`, 0o755);
}

function toolLog(container, name) {
  const file = join(container.state, `${name}.log`);
  return existsSync(file) ? readFileSync(file, 'utf8') : '';
}

// Markers for required runtimes and the executables of dependsOn Features.
function setUp(id, container) {
  for (const dep of byId.get(id).requires) markInstalled(container, dep, version(dep) ?? 'none');
  const coupling = byId.get(id).versionCoupling;
  if (coupling) markInstalled(container, coupling.feature, version(id));
  const deps = byId.get(id).dependsOn;
  if (deps.includes('python/uv')) fakeTool(container, 'usr/local/bin/uv');
  if (byId.get(id).requires.includes('rust/rustc')) {
    writeFile(join(container.root, 'usr/local/rustup/settings.toml'), '');
    fakeTool(container, 'usr/local/cargo/bin/rustup', {script: '[ "$1" = toolchain ] && echo "1.98.0-x86_64-unknown-linux-gnu (default)"'});
    fakeTool(container, 'usr/local/cargo/bin/cargo');
  }
  if (byId.get(id).requires.includes('go/go')) fakeTool(container, 'usr/local/bin/go');
  if (deps.includes('latex/core')) {
    for (const platform of ['x86_64-linux', 'aarch64-linux']) fakeTool(container, `usr/local/texlive/bin/${platform}/tlmgr`);
  }
}

function assertError(result, code, id, detail) {
  assert.notEqual(result.status, 0, result.output);
  const lines = errorLines(result.stderr);
  assert.equal(lines.length, 1, result.output);
  const prefix = `aibox-feature error ${code} feature=${id} `;
  assert.ok(lines[0].startsWith(prefix), `${JSON.stringify(lines[0])} does not start with ${JSON.stringify(prefix)}`);
  if (detail) assert.ok(lines[0].includes(detail), `${JSON.stringify(lines[0])} does not name ${detail}`);
}

function run(id, container, options) {
  return runFeature(featureCopy(id), container, options);
}

it('wave 3 covers every aibox-owned tool and atomic Feature', () => {
  const owned = catalog.features.filter((f) => f.owner === 'aibox' && f.role !== 'bundle').map((f) => f.id);
  const present = new Set(featureDirectories());
  assert.deepEqual(owned.filter((id) => !present.has(id)), []);
});

for (const id of features) {
  const feature = byId.get(id);

  describe(id, () => {
    for (const fixture of catalog.unsupportedTargets) {
      it(`fails on ${fixture.id} with UNSUPPORTED_TARGET before any change`, () => {
        const container = makeContainer({os: fixture.id, arch: 'arm64', available: available(id)});
        setUp(id, container);
        const before = snapshot(container);
        assertError(run(id, container), 'AIBOX_FEATURE_UNSUPPORTED_TARGET', id, `target=${fixture.id}-arm64`);
        assert.deepEqual(snapshot(container), before);
        assert.deepEqual(aptLog(container), []);
      });
    }

    if (feature.options.version) {
      it('fails for version 0.0.0-unsupported with UNSUPPORTED_VERSION before any change', () => {
        const container = makeContainer({available: available(id)});
        setUp(id, container);
        const before = snapshot(container);
        assertError(run(id, container, {options: {version: '0.0.0-unsupported'}}),
          'AIBOX_FEATURE_UNSUPPORTED_VERSION', id, 'version=0.0.0-unsupported');
        assert.deepEqual(snapshot(container), before);
      });
    }

    it('fails with VERSION_CONFLICT when another version is installed', () => {
      const container = makeContainer({available: available(id)});
      setUp(id, container);
      markInstalled(container, id, '0.0.1-other');
      const before = snapshot(container);
      assertError(run(id, container), 'AIBOX_FEATURE_VERSION_CONFLICT', id, 'installed=0.0.1-other');
      assert.deepEqual(snapshot(container), before);
    });

    for (const dep of feature.requires) {
      it(`fails with DEPENDENCY_MISSING naming ${dep} before any change`, () => {
        const container = makeContainer({available: available(id)});
        const before = snapshot(container);
        assertError(run(id, container), 'AIBOX_FEATURE_DEPENDENCY_MISSING', id, `missing=${dep}`);
        assert.deepEqual(snapshot(container), before);
      });
    }

    for (const missing of available(id)) {
      it(`fails with PACKAGE_UNAVAILABLE for ${missing} before any change`, () => {
        const container = makeContainer({available: available(id).filter((p) => p !== missing)});
        setUp(id, container);
        const before = snapshot(container);
        assertError(run(id, container), 'AIBOX_FEATURE_PACKAGE_UNAVAILABLE', id, `package=${missing}`);
        assert.deepEqual(snapshot(container), before);
      });
    }

    switch (method(id)) {
      case 'producer-release':
        for (const [index, arch] of ['amd64', 'arm64'].entries()) {
          const url = feature.tools[0].source.split(' ')[0]
            .replaceAll('{version}', version(id)).replaceAll('{arch}', releaseTokens[id][index]);
          it(`${arch}: fails with DOWNLOAD_FAILED naming ${url}`, () => {
            const container = makeContainer({arch, available: available(id)});
            setUp(id, container);
            const before = snapshot(container);
            assertError(run(id, container), 'AIBOX_FEATURE_DOWNLOAD_FAILED', id, `url=${url}`);
            assert.deepEqual(snapshot(container), before);
          });
          it(`${arch}: fails with VERIFICATION_FAILED for tampered content`, () => {
            const container = makeContainer({arch, available: available(id)});
            setUp(id, container);
            serve(container, url, writeFile(join(tempDir('tampered'), 'artifact'), 'tampered'));
            const before = snapshot(container);
            assertError(run(id, container), 'AIBOX_FEATURE_VERIFICATION_FAILED', id, 'algorithm=sha256');
            assert.deepEqual(snapshot(container), before);
          });
        }
        break;

      case 'debian-package':
        it('installs every package at its candidate without recommends', () => {
          const container = makeContainer({available: debianPackages(id)});
          setUp(id, container);
          const result = run(id, container);
          assert.equal(result.status, 0, result.output);
          assert.deepEqual(installedPackages(container), [...new Set(debianPackages(id))].sort());
          for (const line of aptLog(container).filter((l) => /\binstall\b/.test(l))) assert.match(line, /--no-install-recommends/);
          assert.equal(installedVersion(container, id), 'none');
        });
        break;

      case 'uv-tool': {
        const requirement = feature.tools[0].source.split(' ').find((w) => w.includes('=='))
          .replace('{version}', version(id));
        it(`installs ${requirement} with uv`, () => {
          const container = makeContainer();
          setUp(id, container);
          const result = run(id, container);
          assert.equal(result.status, 0, result.output);
          assert.match(toolLog(container, 'uv'), new RegExp(`^tool install .*${requirement.replaceAll('.', '\\.')}`, 'm'));
          assert.equal(installedVersion(container, id), version(id));
        });
        it('fails with DOWNLOAD_FAILED when uv cannot install it', () => {
          const container = makeContainer();
          setUp(id, container);
          fakeTool(container, 'usr/local/bin/uv', {exit: 2});
          const before = snapshot(container);
          assertError(run(id, container), 'AIBOX_FEATURE_DOWNLOAD_FAILED', id, requirement);
          assert.deepEqual(snapshot(container), before);
        });
        break;
      }

      case 'npm-global': {
        const [name, base, fixed] = firstNpmPackage[id];
        const url = `https://registry.npmjs.org/${name}/-/${base}-${fixed ?? version(id)}.tgz`;
        it(`fails with DOWNLOAD_FAILED naming ${url}`, () => {
          const container = makeContainer();
          setUp(id, container);
          const before = snapshot(container);
          assertError(run(id, container), 'AIBOX_FEATURE_DOWNLOAD_FAILED', id, `url=${url}`);
          assert.deepEqual(snapshot(container), before);
        });
        it('fails with VERIFICATION_FAILED for a tampered package', () => {
          const container = makeContainer();
          setUp(id, container);
          serve(container, url, writeFile(join(tempDir('tampered'), 'package.tgz'), 'tampered'));
          const before = snapshot(container);
          assertError(run(id, container), 'AIBOX_FEATURE_VERIFICATION_FAILED', id, 'algorithm=sha512');
          assert.deepEqual(snapshot(container), before);
        });
        break;
      }

      case 'go-install': {
        const module = feature.tools[0].source.split(' ').at(-1);
        it(`installs ${module} with the installed go`, () => {
          const container = makeContainer();
          setUp(id, container);
          const result = run(id, container);
          assert.equal(result.status, 0, result.output);
          assert.match(toolLog(container, 'go'), new RegExp(`^install ${module.replaceAll('.', '\\.')}$`, 'm'));
        });
        it('fails with DOWNLOAD_FAILED naming the module when go install fails', () => {
          const container = makeContainer();
          setUp(id, container);
          fakeTool(container, 'usr/local/bin/go', {exit: 1});
          const before = snapshot(container);
          assertError(run(id, container), 'AIBOX_FEATURE_DOWNLOAD_FAILED', id, `module=${module}`);
          assert.deepEqual(snapshot(container), before);
        });
        break;
      }

      case 'texlive-installer': {
        const packages = feature.tools[0].checks.filter((c) => c.kind === 'texlive').map((c) => c.package);
        it(`installs ${packages.join(' ')} with tlmgr and removes the build-only GnuPG`, () => {
          for (const arch of ['amd64', 'arm64']) {
            const container = makeContainer({arch, available: available(id)});
            setUp(id, container);
            const result = run(id, container);
            assert.equal(result.status, 0, result.output);
            assert.match(toolLog(container, 'tlmgr'), new RegExp(`^install ${packages.join(' ')}$`, 'm'));
            assert.ok(!installedPackages(container).includes('gnupg'));
          }
        });
        it('fails with DOWNLOAD_FAILED when tlmgr cannot install the packages', () => {
          const container = makeContainer({available: available(id)});
          setUp(id, container);
          for (const platform of ['x86_64-linux', 'aarch64-linux']) fakeTool(container, `usr/local/texlive/bin/${platform}/tlmgr`, {exit: 1});
          const before = snapshot(container);
          assertError(run(id, container), 'AIBOX_FEATURE_DOWNLOAD_FAILED', id, 'source=texlive');
          assert.deepEqual(snapshot(container), before);
        });
        break;
      }

      case 'playwright-install': {
        const path = feature.tools[0].checks.find((c) => c.kind === 'path').path;
        const name = id.split('/').at(-1);
        it(`installs ${path} with the installed playwright`, () => {
          const container = makeContainer();
          setUp(id, container);
          fakeTool(container, 'usr/local/bin/playwright', {script: `mkdir -p "$PLAYWRIGHT_BROWSERS_PATH/${path.split('/').at(-1)}"`});
          const result = run(id, container);
          assert.equal(result.status, 0, result.output);
          assert.match(toolLog(container, 'playwright'), new RegExp(`^install --with-deps --no-shell ${name}$`, 'm'));
          assert.ok(existsSync(join(container.root, path)));
        });
        it('fails with DOWNLOAD_FAILED when playwright cannot install the browser', () => {
          const container = makeContainer();
          setUp(id, container);
          fakeTool(container, 'usr/local/bin/playwright', {exit: 1});
          const before = snapshot(container);
          assertError(run(id, container), 'AIBOX_FEATURE_DOWNLOAD_FAILED', id, `browser=${name}`);
          assert.deepEqual(snapshot(container), before);
        });
        it('fails when playwright installs another revision', () => {
          const container = makeContainer();
          setUp(id, container);
          fakeTool(container, 'usr/local/bin/playwright');
          const result = run(id, container);
          assert.notEqual(result.status, 0);
          assert.equal(installedVersion(container, id), undefined);
        });
        break;
      }

      default:
        break;
    }
  });
}

describe('rust components', () => {
  for (const component of ['clippy', 'rustfmt']) {
    it(`rust/${component} adds the ${component} component`, () => {
      const container = makeContainer();
      setUp(`rust/${component}`, container);
      assert.equal(run(`rust/${component}`, container).status, 0);
      assert.match(toolLog(container, 'rustup'), new RegExp(`^component add ${component}$`, 'm'));
    });
    it(`rust/${component} fails with DOWNLOAD_FAILED when rustup fails`, () => {
      const container = makeContainer();
      setUp(`rust/${component}`, container);
      fakeTool(container, 'usr/local/cargo/bin/rustup', {exit: 1});
      assertError(run(`rust/${component}`, container), 'AIBOX_FEATURE_DOWNLOAD_FAILED', `rust/${component}`, `component=${component}`);
    });
  }

  it('rust/cargo-audit builds the pinned crate with --locked', () => {
    const container = makeContainer();
    setUp('rust/cargo-audit', container);
    assert.equal(run('rust/cargo-audit', container).status, 0);
    assert.match(toolLog(container, 'cargo'), /^install --locked cargo-audit --version 0\.22\.2$/m);
  });

  it('rust/cargo-audit fails with DOWNLOAD_FAILED when cargo install fails', () => {
    const container = makeContainer();
    setUp('rust/cargo-audit', container);
    fakeTool(container, 'usr/local/cargo/bin/cargo', {exit: 101});
    assertError(run('rust/cargo-audit', container), 'AIBOX_FEATURE_DOWNLOAD_FAILED', 'rust/cargo-audit', 'cargo-audit@0.22.2');
  });

  it('rust/x86-64-cross changes nothing but its marker on amd64', () => {
    const container = makeContainer({arch: 'amd64'});
    setUp('rust/x86-64-cross', container);
    const result = run('rust/x86-64-cross', container);
    assert.equal(result.status, 0, result.output);
    assert.deepEqual(aptLog(container), []);
    assert.equal(toolLog(container, 'rustup'), '');
  });

  it('rust/x86-64-cross installs the cross toolchain and adds the target on arm64', () => {
    const container = makeContainer({arch: 'arm64', available: ['gcc-x86-64-linux-gnu', 'libc6-dev-amd64-cross']});
    setUp('rust/x86-64-cross', container);
    const result = run('rust/x86-64-cross', container);
    assert.equal(result.status, 0, result.output);
    assert.deepEqual(installedPackages(container), ['gcc-x86-64-linux-gnu', 'libc6-dev-amd64-cross']);
    assert.match(toolLog(container, 'rustup'), /^target add --toolchain 1\.98\.0-x86_64-unknown-linux-gnu x86_64-unknown-linux-gnu$/m);
  });

  it('rust/x86-64-cross fails with PACKAGE_UNAVAILABLE on arm64 without the cross packages', () => {
    const container = makeContainer({arch: 'arm64', available: ['libc6-dev-amd64-cross']});
    setUp('rust/x86-64-cross', container);
    const before = snapshot(container);
    assertError(run('rust/x86-64-cross', container), 'AIBOX_FEATURE_PACKAGE_UNAVAILABLE', 'rust/x86-64-cross', 'package=gcc-x86-64-linux-gnu');
    assert.deepEqual(snapshot(container), before);
  });
});

describe('cloud/gcp', () => {
  const keyUrl = 'https://packages.cloud.google.com/apt/doc/apt-key.gpg';
  const keyFile = () => writeFile(join(tempDir('key'), 'apt-key.gpg'), 'not the pinned key');

  it('fails with DOWNLOAD_FAILED naming the key URL and leaves the APT sources unchanged', () => {
    const container = makeContainer();
    const before = snapshot(container);
    assertError(run('cloud/gcp', container), 'AIBOX_FEATURE_DOWNLOAD_FAILED', 'cloud/gcp', `url=${keyUrl}`);
    assert.deepEqual(snapshot(container), before);
  });

  it('fails with VERIFICATION_FAILED for another key', () => {
    const container = makeContainer();
    serve(container, keyUrl, keyFile());
    const before = snapshot(container);
    assertError(run('cloud/gcp', container), 'AIBOX_FEATURE_VERIFICATION_FAILED', 'cloud/gcp', 'algorithm=sha256');
    assert.deepEqual(snapshot(container), before);
  });
});
