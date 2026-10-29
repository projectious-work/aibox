// Component tests for `scripts/v1-05-features.sh publish` (chapter 8,
// "Publishing and running tests"). The Dev Container CLI port is replaced by
// a Fake that records what it was asked to publish.
import assert from 'node:assert/strict';
import {spawnSync} from 'node:child_process';
import {cpSync, existsSync, mkdirSync, readFileSync, readdirSync, symlinkSync} from 'node:fs';
import {dirname, join} from 'node:path';
import {describe, it} from 'node:test';

import {loadCatalog, rewriteMetadata} from '../../scripts/v1-05-features.mjs';
import {featuresSrc, repoRoot, tempDir, writeFile} from './harness.mjs';

const catalog = loadCatalog();
const script = join(repoRoot, 'scripts/v1-05-features.sh');

// Fake Dev Container CLI 0.89.0: copies each published collection to
// $FAKE_STATE/published/<registry>/<namespace>/ and fails on request.
function fakeCli() {
  const state = tempDir('cli');
  const cli = writeFile(join(state, 'devcontainer'), `#!/bin/sh
if [ "$1" = --version ]; then echo 0.89.0; exit 0; fi
[ "$1 $2" = "features publish" ] || { echo "unexpected: $*" >&2; exit 2; }
shift 2
registry=""; namespace=""; target=""
while [ $# -gt 0 ]; do
  case "$1" in
    --registry|-r) registry=$2; shift 2 ;;
    --namespace|-n) namespace=$2; shift 2 ;;
    --log-level) shift 2 ;;
    *) target=$1; shift ;;
  esac
done
[ ! -f "${state}/fail" ] || { echo "(!) ERR: Failed to publish" >&2; exit 1; }
mkdir -p "${state}/published/$registry/$namespace"
cp -R "$target"/. "${state}/published/$registry/$namespace/"
`, 0o755);
  return {cli, state};
}

// The script takes devcontainer and node from PATH (chapter 8, Publishing).
function publish(args, {cli, features, path} = {}) {
  const env = {...process.env};
  if (cli) env.PATH = `${dirname(cli)}:${process.env.PATH}`;
  if (path !== undefined) env.PATH = path;
  const extra = features ? ['--features-dir', features] : [];
  const result = spawnSync(script, ['publish', ...args, ...extra], {encoding: 'utf8', env});
  return {...result, output: result.stdout + result.stderr};
}

function published(state, registry, id) {
  const file = join(state, 'published', registry, ...id.split('/').slice(0, -1), id.split('/').at(-1), 'devcontainer-feature.json');
  return existsSync(file) ? JSON.parse(readFileSync(file, 'utf8')) : undefined;
}

describe('publish', () => {
  const selection = ['python/uv', 'python/python', 'node/node', 'browser-testing/playwright'];
  const featureArgs = selection.flatMap((id) => ['--feature', id]);

  it('publishes each Feature under <host>/<prefix>/<catalog-id> with rewritten references', () => {
    const {cli, state} = fakeCli();
    const result = publish(['--registry', 'localhost:5000/aibox-features', ...featureArgs], {cli});
    assert.equal(result.status, 0, result.output);
    for (const id of selection) {
      const metadata = published(state, 'localhost:5000', `aibox-features/${id}`);
      assert.ok(metadata, `${id} was not published under localhost:5000/aibox-features/${id}`);
      assert.equal(metadata.id, id.split('/').at(-1));
      assert.ok(existsSync(join(state, 'published/localhost:5000/aibox-features', id, 'install.sh')));
      assert.ok(existsSync(join(state, 'published/localhost:5000/aibox-features', id, 'lib/aibox-feature.sh')));
    }
    assert.deepEqual(published(state, 'localhost:5000', 'aibox-features/python/python').dependsOn,
      {'localhost:5000/aibox-features/python/uv:1': {}});
    assert.deepEqual(published(state, 'localhost:5000', 'aibox-features/browser-testing/playwright').installsAfter,
      ['localhost:5000/aibox-features/node/node']);
  });

  it('publishes every aibox-owned Feature, bundles included, when no --feature is given', () => {
    const {cli, state} = fakeCli();
    const result = publish(['--registry', 'localhost:5000/aibox-features'], {cli});
    assert.equal(result.status, 0, result.output);
    for (const feature of catalog.features) {
      const metadata = published(state, 'localhost:5000', `aibox-features/${feature.id}`);
      assert.equal(metadata !== undefined, feature.owner === 'aibox', feature.id);
    }
    const bundle = published(state, 'localhost:5000', 'aibox-features/git-ui/all');
    const upstream = catalog.features.find((f) => f.id === 'git-ui/gh').upstream;
    assert.deepEqual(bundle.dependsOn, {
      [upstream.ref]: upstream.options,
      'localhost:5000/aibox-features/git-ui/lazygit:1': {},
    });
  });

  it('leaves the source metadata unchanged', () => {
    const {cli} = fakeCli();
    const before = readFileSync(join(featuresSrc, 'python/python/devcontainer-feature.json'), 'utf8');
    assert.equal(publish(['--registry', 'localhost:5000/aibox-features', ...featureArgs], {cli}).status, 0);
    assert.equal(readFileSync(join(featuresSrc, 'python/python/devcontainer-feature.json'), 'utf8'), before);
  });

  it('leaves upstream references unchanged', () => {
    const upstream = catalog.features.find((f) => f.owner === 'upstream');
    const metadata = {
      dependsOn: {[`${catalog.registry}/go/go:1`]: {}, [upstream.upstream.ref]: upstream.upstream.options},
      installsAfter: [`${catalog.registry}/node/node`],
    };
    assert.deepEqual(rewriteMetadata(metadata, catalog, 'localhost:5000/aibox-features'), {
      dependsOn: {'localhost:5000/aibox-features/go/go:1': {}, [upstream.upstream.ref]: upstream.upstream.options},
      installsAfter: ['localhost:5000/aibox-features/node/node'],
    });
  });

  it('fails with a message when an aibox-owned Feature has no sources', () => {
    const {cli, state} = fakeCli();
    const features = join(tempDir('partial'), 'src');
    mkdirSync(join(features, 'python'), {recursive: true});
    cpSync(join(featuresSrc, 'python/uv'), join(features, 'python/uv'), {recursive: true});
    const result = publish(['--registry', 'localhost:5000/aibox-features'], {cli, features});
    assert.notEqual(result.status, 0);
    assert.match(result.stderr, /aibox-owned Feature node\/node has no sources/);
    assert.equal(existsSync(join(state, 'published')), false, 'nothing may be published');
  });

  it('fails with a message when the CLI fails to publish', () => {
    const {cli, state} = fakeCli();
    writeFile(join(state, 'fail'), '');
    const result = publish(['--registry', 'localhost:5000/aibox-features', '--feature', 'go/go'], {cli});
    assert.notEqual(result.status, 0);
    assert.match(result.stderr, /v1-05-features: publishing localhost:5000\/aibox-features\/go failed/);
  });

  for (const registry of ['localhost:5000', 'localhost:5000/', '/aibox-features', 'localhost:5000/Aibox']) {
    it(`rejects --registry ${registry}`, () => {
      const {cli, state} = fakeCli();
      const result = publish(['--registry', registry, '--feature', 'go/go'], {cli});
      assert.notEqual(result.status, 0);
      assert.match(result.stderr, /<host>\/<prefix>/);
      assert.equal(existsSync(join(state, 'published')), false);
    });
  }

  it('rejects a Dev Container CLI other than the pinned version', () => {
    const cli = writeFile(join(tempDir('old-cli'), 'devcontainer'), '#!/bin/sh\necho 0.80.0\n', 0o755);
    const result = publish(['--registry', 'localhost:5000/aibox-features', '--feature', 'go/go'], {cli});
    assert.notEqual(result.status, 0);
    assert.match(result.stderr, /Dev Container CLI 0\.89\.0 is required/);
  });

  it('fails with a message when devcontainer is not on PATH', () => {
    const bin = tempDir('bin');
    symlinkSync(process.execPath, join(bin, 'node'));
    for (const tool of ['bash', 'env', 'sed', 'mktemp', 'rm', 'ls', 'tr', 'cat', 'dirname']) {
      const target = ['/usr/bin', '/bin'].map((d) => join(d, tool)).find((p) => existsSync(p));
      if (target) symlinkSync(target, join(bin, tool));
    }
    const result = publish(['--registry', 'localhost:5000/aibox-features', '--feature', 'go/go'], {path: bin});
    assert.notEqual(result.status, 0);
    assert.match(result.stderr, /devcontainer/);
  });

  it('publishes under a multi-segment prefix', () => {
    const {cli, state} = fakeCli();
    assert.equal(publish(['--registry', 'registry.example:5000/a/b', '--feature', 'zig/zig'], {cli}).status, 0);
    assert.deepEqual(readdirSync(join(state, 'published/registry.example:5000/a/b/zig')), ['zig']);
  });
});
