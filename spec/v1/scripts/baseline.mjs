// Read-only V1-01 evidence check. This validates the checked-in v0 inventory;
// it does not assert that a v1 runtime implements any inventoried behavior.
import {execFileSync} from 'node:child_process';
import {createHash} from 'node:crypto';
import {readFileSync} from 'node:fs';
import {fileURLToPath} from 'node:url';
import path from 'node:path';

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../../..');
const spec = path.join(root, 'spec/v1');
const fixture = JSON.parse(readFileSync(path.join(spec, 'fixtures/baseline/v0.35.0.json')));
const ledger = name => JSON.parse(readFileSync(path.join(spec, `ledger/${name}.json`)));
const git = (...args) => execFileSync('git', args, {cwd: root, encoding: 'utf8', maxBuffer: 16e6});
const fail = message => {throw Error(message);};
const assert = (condition, message) => {if (!condition) fail(message);};
const hash = body => createHash('sha256').update(body).digest('hex');
const commit = fixture.commit;
assert(/^[0-9a-f]{40}$/.test(commit), 'fixture commit must be a full SHA');
assert(git('rev-parse', `${commit}^{commit}`).trim() === commit, 'baseline commit is unavailable');
const census = ledger('census');
assert(census.baseline === commit, 'census baseline differs from fixture');

const tree = git('ls-tree', '-r', '--name-only', commit).trim().split('\n');
const paths = new Set(tree);
const contents = new Map();
function source(file) {
  assert(paths.has(file), `source absent at baseline: ${file}`);
  if (!contents.has(file)) contents.set(file, git('show', `${commit}:${file}`));
  return contents.get(file);
}
const linkPattern = new RegExp(`^https://github\\.com/projectious-work/aibox/blob/${commit}/([^#]+)(?:#L([1-9]\\d*))?$`);
function verifyLink(link, label) {
  assert(typeof link === 'string', `${label}: missing source link`);
  const match = link.match(linkPattern);
  assert(match, `${label}: source is not pinned to ${commit}: ${link}`);
  const body = source(match[1]);
  if (match[2]) assert(Number(match[2]) <= body.split('\n').length, `${label}: source line out of range: ${link}`);
  return {file: match[1], body};
}
function unique(rows, key, label) {
  const seen = new Set();
  for (const row of rows) {
    const value = row[key];
    assert(typeof value === 'string' && value.length, `${label}: unresolved ${key}`);
    assert(!seen.has(value), `${label}: duplicate ${key} ${value}`);
    seen.add(value);
  }
  return seen;
}
function sameSet(actual, expected, label) {
  assert(actual.size === expected.size, `${label}: source/ledger cardinality differs (${actual.size} vs ${expected.size})`);
  for (const value of expected) assert(actual.has(value), `${label}: missing row for ${value}`);
}
function count(key, n) {
  assert(n === fixture.counts[key], `${key}: expected ${fixture.counts[key]}, found ${n}`);
  assert(census[key] === n, `${key}: census says ${census[key]}, found ${n}`);
}
function resolved(row, label, fields) {
  for (const field of fields) assert(typeof row[field] === 'string' && row[field].trim(), `${label}: unresolved ${field}`);
}

const config = ledger('configuration');
count('config_fields', config.length);
unique(config, 'id', 'configuration');
for (const row of config) {
  assert(row.id === `CFG:${row.path}`, `${row.id}: path/id mismatch`);
  resolved(row, row.id, ['type', 'disposition', 'target', 'acceptance']);
  for (const field of ['source', 'default_evidence', 'enum_evidence']) {
    if (row[field]) verifyLink(row[field], `${row.id}.${field}`);
  }
  assert(row.source, `${row.id}: unresolved source`);
}

const addons = ledger('addons');
count('addon_recipes', addons.length);
unique(addons, 'id', 'addons');
const addonPaths = new Set();
let toolCount = 0;
for (const row of addons) {
  const {file, body} = verifyLink(row.source, row.id);
  addonPaths.add(file);
  assert(row.sha256 === hash(body), `${row.id}: source hash mismatch`);
  assert(Array.isArray(row.tools) && row.tools.length, `${row.id}: no tools`);
  unique(row.tools, 'name', row.id);
  for (const tool of row.tools) resolved(tool, `${row.id}.${tool.name}`, ['declaration', 'target', 'acceptance']);
  toolCount += row.tools.length;
}
count('addon_tools', toolCount);
sameSet(addonPaths, new Set(tree.filter(file => file.startsWith('addons/') && file.endsWith('.yaml'))), 'addon recipes');

const types = ledger('commands');
const actions = ledger('command-actions');
const args = ledger('command-arguments');
count('cli_declarations', types.length);
count('cli_actions', actions.length);
count('cli_arguments', args.length);
const typeIds = unique(types, 'id', 'command declarations');
const actionIds = unique(actions, 'id', 'command actions');
unique(args, 'id', 'command arguments');
for (const row of types) {verifyLink(row.source, row.id); resolved(row, row.id, ['declaration', 'acceptance']);}
for (const row of actions) {
  assert(typeIds.has(`CMDTYPE:${row.parent}`), `${row.id}: missing parent declaration`);
  verifyLink(row.source, row.id);
  resolved(row, row.id, ['declaration', 'target', 'acceptance']);
}
for (const row of args) {
  assert(actionIds.has(row.command) || typeIds.has(row.command), `${row.id}: missing command`);
  verifyLink(row.source, row.id);
  resolved(row, row.id, ['target', 'acceptance']);
}

const runtime = ledger('runtime-assets');
count('runtime_assets', runtime.length);
unique(runtime, 'id', 'runtime assets');
const runtimePaths = new Set();
for (const row of runtime) {
  const {file, body} = verifyLink(row.source, row.id);
  assert(row.id === `ASSET:${file}`, `${row.id}: path/id mismatch`);
  assert(row.sha256 === hash(body), `${row.id}: source hash mismatch`);
  resolved(row, row.id, ['disposition', 'acceptance']);
  runtimePaths.add(file);
}
sameSet(runtimePaths, new Set(tree.filter(file => fixture.source_scopes.runtime.some(prefix => file.startsWith(`${prefix}/`)) && !file.endsWith('/AGENTS.md'))), 'runtime assets');

const docs = ledger('documentation');
count('documentation_pages', docs.length);
const docPaths = new Set();
for (const row of docs) {
  const {file, body} = verifyLink(row.source, 'documentation');
  assert(!docPaths.has(file), `documentation: duplicate ${file}`);
  assert(row.sha256 === hash(body), `documentation: source hash mismatch ${file}`);
  docPaths.add(file);
}
sameSet(docPaths, new Set(tree.filter(file => file.startsWith('docs-site/content/docs/') && file.endsWith('.md'))), 'documentation pages');

const env = ledger('environment');
count('environment_identifiers', env.length);
const envNames = unique(env, 'name', 'environment identifiers');
const sourceEnv = new Set();
for (const file of tree.filter(file => (file.startsWith('cli/src/') || file.startsWith('images/')) && /\.(rs|sh)$/.test(file))) {
  for (const match of source(file).matchAll(/\b(AIBOX_[A-Z0-9_]+)\b/g)) sourceEnv.add(match[1]);
}
sameSet(envNames, sourceEnv, 'environment identifiers');
for (const row of env) resolved(row, row.name, ['classification', 'acceptance']);

for (const name of ['configuration-types', 'defaults', 'base-build']) {
  const rows = ledger(name);
  assert(rows.length, `${name}: empty ledger`);
  for (const [index, row] of rows.entries()) {
    const label = `${name}[${index}]`;
    if (row.source && row.source.startsWith('https://')) verifyLink(row.source, label);
    if (row.url) verifyLink(row.url, label);
    if (row.file) {
      source(row.file);
      assert(Number.isInteger(row.line) && row.line > 0, `${label}: missing source line`);
    }
    if (row.sha256) assert(row.sha256 === hash(verifyLink(row.source, label).body), `${label}: source hash mismatch`);
  }
}

// The generator's exact-output comparison catches omissions in parsed fields,
// declarations, and metadata that cardinality and link checks alone cannot.
execFileSync(process.execPath, [path.join(spec, 'scripts/inventory.mjs'), '--check'], {cwd: root, encoding: 'utf8', stdio: 'pipe', maxBuffer: 16e6});
console.log(`V1-01 baseline ${commit}: verified ${config.length} config fields, ${addons.length} recipes/${toolCount} tools, ${actions.length} actions/${args.length} arguments, ${runtime.length} assets, ${docs.length} docs, ${env.length} env IDs; source links, hashes and inventory match.`);
