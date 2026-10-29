// Component-test harness for aibox tool Feature install scripts.
//
// Ports replaced with Test Doubles (Meszaros): the filesystem root
// (AIBOX_FEATURE_ROOT), and the child executables dpkg, dpkg-query, apt-get,
// apt-cache and curl (Fakes in ./fakes). Everything else — the shared
// library, the generated catalog data and each install.sh — is the real code.
import {spawnSync} from 'node:child_process';
import {createHash} from 'node:crypto';
import {
  cpSync, existsSync, mkdirSync, mkdtempSync, readdirSync, readFileSync, statSync, writeFileSync,
} from 'node:fs';
import {tmpdir} from 'node:os';
import {dirname, join, relative} from 'node:path';
import {fileURLToPath} from 'node:url';

export const repoRoot = join(dirname(fileURLToPath(import.meta.url)), '..', '..');
export const featuresSrc = join(repoRoot, 'features', 'src');
const fakesDir = join(repoRoot, 'features', 'test', 'fakes');

export const osReleases = {
  'debian-13': 'PRETTY_NAME="Debian GNU/Linux 13 (trixie)"\nNAME="Debian GNU/Linux"\nVERSION_ID="13"\nVERSION="13 (trixie)"\nVERSION_CODENAME=trixie\nID=debian\n',
  'debian-12': 'PRETTY_NAME="Debian GNU/Linux 12 (bookworm)"\nNAME="Debian GNU/Linux"\nVERSION_ID="12"\nVERSION="12 (bookworm)"\nVERSION_CODENAME=bookworm\nID=debian\n',
  'ubuntu-24.04': 'PRETTY_NAME="Ubuntu 24.04.3 LTS"\nNAME="Ubuntu"\nVERSION_ID="24.04"\nVERSION="24.04.3 LTS (Noble Numbat)"\nVERSION_CODENAME=noble\nID=ubuntu\nID_LIKE=debian\n',
};

export function tempDir(prefix) {
  return mkdtempSync(join(tmpdir(), `aibox-feature-${prefix}-`));
}

// A fake container: a root directory plus the state of the fake executables.
export function makeContainer({os = 'debian-13', arch = 'amd64', available = [], installed = []} = {}) {
  const root = tempDir('root');
  const state = tempDir('state');
  mkdirSync(join(root, 'etc'), {recursive: true});
  mkdirSync(join(root, 'var/lib/apt/lists'), {recursive: true});
  writeFileSync(join(root, 'etc/os-release'), osReleases[os]);
  writeFileSync(join(state, 'arch'), `${arch}\n`);
  writeFileSync(join(state, 'available'), available.map((p) => `${p}\n`).join(''));
  writeFileSync(join(state, 'installed'), installed.map((p) => `${p}\n`).join(''));
  writeFileSync(join(state, 'urls'), '');
  return {root, state};
}

export function serve(container, url, file) {
  writeFileSync(join(container.state, 'urls'), `${url} ${file}\n`, {flag: 'a'});
}

export function installedPackages(container) {
  return readFileSync(join(container.state, 'installed'), 'utf8').split('\n').filter(Boolean).sort();
}

export function aptLog(container) {
  const file = join(container.state, 'apt.log');
  return existsSync(file) ? readFileSync(file, 'utf8').split('\n').filter(Boolean) : [];
}

// Records the observable state of the fake root: every path with its digest.
export function snapshot(container) {
  const entries = [];
  const walk = (dir) => {
    for (const name of readdirSync(dir).sort()) {
      const path = join(dir, name);
      const info = statSync(path);
      const rel = relative(container.root, path);
      if (info.isDirectory()) {
        entries.push(`${rel}/`);
        walk(path);
      } else {
        entries.push(`${rel} ${createHash('sha256').update(readFileSync(path)).digest('hex')}`);
      }
    }
  };
  walk(container.root);
  return {files: entries, packages: installedPackages(container)};
}

export function markInstalled(container, featureId, version) {
  const dir = join(container.root, 'var/lib/aibox-features', featureId);
  mkdirSync(dir, {recursive: true});
  writeFileSync(join(dir, 'version'), `${version}\n`);
}

export function installedVersion(container, featureId) {
  const file = join(container.root, 'var/lib/aibox-features', featureId, 'version');
  return existsSync(file) ? readFileSync(file, 'utf8').trim() : undefined;
}

// Runs a Feature's install.sh the way the Dev Container CLI does: as an
// executable in the Feature directory with options as upper-case variables.
export function runFeature(featureDir, container, {options = {}, env = {}} = {}) {
  const optionEnv = Object.fromEntries(Object.entries(options).map(([k, v]) => [k.toUpperCase(), String(v)]));
  const result = spawnSync(join(featureDir, 'install.sh'), [], {
    cwd: featureDir,
    encoding: 'utf8',
    env: {
      PATH: `${fakesDir}:/usr/bin:/bin`,
      HOME: container.state,
      TMPDIR: container.state,
      AIBOX_FEATURE_ROOT: container.root,
      FAKE_STATE: container.state,
      ...optionEnv,
      ...env,
    },
  });
  return {status: result.status, stdout: result.stdout, stderr: result.stderr, output: result.stdout + result.stderr};
}

export function errorLines(output) {
  return output.split('\n').filter((line) => line.includes('aibox-feature error '));
}

// Copies a Feature source directory so a test can run it without touching
// the repository (install scripts never write into their own directory, but
// the copy keeps tests independent of each other).
export function featureCopy(id) {
  const target = join(tempDir('feature'), id.replaceAll('/', '-'));
  cpSync(join(featuresSrc, id), target, {recursive: true});
  return target;
}

export function sha256(file) {
  return createHash('sha256').update(readFileSync(file)).digest('hex');
}

export function sha512Base64(file) {
  return createHash('sha512').update(readFileSync(file)).digest('base64');
}

export function writeFile(path, content, mode) {
  mkdirSync(dirname(path), {recursive: true});
  writeFileSync(path, content, mode === undefined ? undefined : {mode});
  return path;
}
