// Validate native Dev Container Templates and their checked-in rendered examples.
// Runtime image build/up/exec qualification lives in scripts/verify-v1-04.sh.
import { existsSync, readFileSync, readdirSync } from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const specRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
const repoRoot = path.resolve(specRoot, '..', '..');
const readJSON = relative => JSON.parse(readFileSync(path.join(specRoot, relative), 'utf8'));
// Maintained examples use full-line JSONC comments for inactive options.
const readJSONC = relative => JSON.parse(readFileSync(path.join(specRoot, relative), 'utf8')
  .split('\n').filter(line => !line.trimStart().startsWith('//')).join('\n'));
const readExternalJSONC = file => JSON.parse(readFileSync(file, 'utf8')
  .split('\n').filter(line => !line.trimStart().startsWith('//')).join('\n'));
const assert = (condition, message) => { if (!condition) throw new Error(message); };

function checkSchema(value, schema, at = '$', rootSchema = schema) {
  if (schema.$ref) {
    assert(schema.$ref.startsWith('#/$defs/'), `${at}: external schema references are unsupported`);
    return checkSchema(value, rootSchema.$defs[schema.$ref.slice(8)], at, rootSchema);
  }
  if (schema.type) {
    const actual = Array.isArray(value) ? 'array' : value === null ? 'null'
      : Number.isInteger(value) ? 'integer' : typeof value;
    assert(actual === schema.type, `${at}: expected ${schema.type}, got ${actual}`);
  }
  if (schema.const !== undefined) assert(value === schema.const, `${at}: expected constant ${schema.const}`);
  if (schema.enum) assert(schema.enum.includes(value), `${at}: value is outside the enum`);
  if (schema.required) for (const key of schema.required) assert(Object.hasOwn(value, key), `${at}: missing ${key}`);
  if (schema.properties) {
    if (schema.additionalProperties === false) {
      for (const key of Object.keys(value)) assert(Object.hasOwn(schema.properties, key), `${at}: unknown ${key}`);
    } else if (typeof schema.additionalProperties === 'object') {
      for (const [key, child] of Object.entries(value)) {
        if (!Object.hasOwn(schema.properties, key)) checkSchema(child, schema.additionalProperties, `${at}.${key}`, rootSchema);
      }
    }
    for (const [key, child] of Object.entries(schema.properties)) {
      if (Object.hasOwn(value, key)) checkSchema(value[key], child, `${at}.${key}`, rootSchema);
    }
  }
  if (schema.items) for (const [index, child] of value.entries()) checkSchema(child, schema.items, `${at}[${index}]`, rootSchema);
  if (schema.allOf) for (const rule of schema.allOf) {
    const [key, condition] = Object.entries(rule.if.properties)[0];
    if (condition.const === value[key] || condition.enum?.includes(value[key])) checkSchema(value, rule.then, at, rootSchema);
  }
}

const customizationSchema = readJSON('customization.schema.json');
const nativeFeatures = readJSON('fixtures/v1-05-native-features.json');
const exampleNames = ['minimal', 'customized', 'custom-user', 'bind-home'];
const examples = new Map();
for (const name of exampleNames) {
  const relative = `examples/${name}/.devcontainer/devcontainer.json`;
  const file = path.join(specRoot, relative);
  assert(existsSync(file), `${name}: missing rendered native example`);
  const example = readJSONC(relative);
  examples.set(name, example);

  assert(example.build && !example.image && !example.dockerComposeFile,
    `${name}: select exactly one project-owned native build definition`);
  assert(example.build.dockerfile === 'Dockerfile' && example.build.context === '.',
    `${name}: build must use its checked-in Dockerfile and project context`);
  for (const required of ['Dockerfile', 'install-tools.sh']) {
    assert(existsSync(path.join(specRoot, 'examples', name, '.devcontainer', required)),
      `${name}: missing .devcontainer/${required}`);
  }
  assert(example.containerUser && example.remoteUser, `${name}: both containerUser and remoteUser must be explicit`);
  assert(example.containerUser === example.remoteUser, `${name}: current starter must use one aligned user`);
  assert(Array.isArray(example.mounts) && example.mounts.length === 1,
    `${name}: declare exactly one persistent-home mount`);
  const mount = example.mounts[0];
  const userHome = `/home/${example.remoteUser}`;
  if (typeof mount === 'string') {
    const fields = Object.fromEntries(mount.split(',').map(field => {
      const index = field.indexOf('=');
      return [field.slice(0, index), field.slice(index + 1)];
    }));
    assert(fields.source === 'aibox-home-${devcontainerId}' && fields.target === userHome && fields.type === 'volume',
      `${name}: named home volume must be project-scoped and target ${userHome}`);
  } else {
    assert(name === 'bind-home' && mount.type === 'bind' &&
      mount.source === '${localWorkspaceFolder}/.aibox-home' && mount.target === userHome,
    `${name}: only bind-home may opt into the explicit project-local bind`);
  }
  assert(!Object.hasOwn(example, 'aibox'), `${name}: aibox data belongs under customizations.aibox`);
  const extension = example.customizations?.aibox;
  if (name === 'minimal') {
    assert(extension === undefined, 'minimal: aibox UX configuration must be optional');
    assert(!Object.hasOwn(example, 'features'), 'minimal: optional tools must remain absent');
  }
  if (name === 'customized') {
    assert(extension !== undefined, 'customized: expected optional customizations.aibox example');
    checkSchema(extension, customizationSchema, '$.customizations.aibox');
    assert(JSON.stringify(example.features) === JSON.stringify(nativeFeatures.selected),
      'customized: native Feature selection must match the reviewed fixture');
  }
}

const customUserDockerfile = readFileSync(path.join(specRoot, 'examples/custom-user/.devcontainer/Dockerfile'), 'utf8');
assert(/ARG USERNAME=dev\b/.test(customUserDockerfile), 'custom-user: Dockerfile user must match the example');
assert(/useradd[^\n]*\$USERNAME/.test(customUserDockerfile), 'custom-user: Dockerfile must create the declared user');
assert(examples.get('custom-user').remoteUser === 'dev', 'custom-user: expected dev remote user');

for (const name of ['minimal', 'curated']) {
  const templateRoot = path.join(repoRoot, 'templates', 'src', name);
  assert(existsSync(templateRoot), `${name}: missing source Template`);
  const metadata = JSON.parse(readFileSync(path.join(templateRoot, 'devcontainer-template.json'), 'utf8'));
  assert(metadata.id === name && metadata.version === '1.0.0', `${name}: Template id/version mismatch`);
  for (const field of ['name', 'description', 'documentationURL', 'licenseURL', 'publisher'])
    assert(typeof metadata[field] === 'string' && metadata[field].length > 0, `${name}: missing Template ${field}`);
  assert(metadata.publisher === 'projectious-work', `${name}: unexpected Template publisher`);
  assert(Array.isArray(metadata.platforms) && metadata.platforms.includes('Any'), `${name}: Template must declare supported platforms`);
  const option = metadata.options?.harness;
  assert(option?.type === 'string' && option.default === 'codex' &&
    JSON.stringify(option.enum) === JSON.stringify(['codex', 'none']), `${name}: harness option contract changed`);
  for (const relative of ['.devcontainer/devcontainer.json', '.devcontainer/Dockerfile', '.devcontainer/install-tools.sh'])
    assert(existsSync(path.join(templateRoot, relative)), `${name}: Template is missing ${relative}`);
  const sourceConfig = readExternalJSONC(path.join(templateRoot, '.devcontainer/devcontainer.json'));
  assert(sourceConfig.build?.args?.HARNESS === '${templateOption:harness}',
    `${name}: harness choice must flow through the standard Template option`);
  assert(sourceConfig.containerUser === 'aibox' && sourceConfig.remoteUser === 'aibox',
    `${name}: source Template user contract changed`);
  assert(sourceConfig.build?.dockerfile === 'Dockerfile' && sourceConfig.build?.context === '.',
    `${name}: source Template must use its project-owned Dockerfile and context`);
  assert(sourceConfig.mounts?.length === 1 && sourceConfig.mounts[0] ===
    'source=aibox-home-${devcontainerId},target=/home/aibox,type=volume', `${name}: source Template home-volume contract changed`);
  const templateIgnore = readFileSync(path.join(templateRoot, '.devcontainer/.gitignore'), 'utf8').split('\n');
  for (const ignored of ['*.env', 'compose.local.yaml'])
    assert(templateIgnore.includes(ignored), `${name}: Template must ignore ${ignored}`);
  const dockerfile = readFileSync(path.join(templateRoot, '.devcontainer/Dockerfile'), 'utf8');
  assert(/USER \$\{USERNAME\}/.test(dockerfile) && /ARG USERNAME=aibox/.test(dockerfile),
    `${name}: Dockerfile must create and run as the declared non-root user`);
  assert(/FROM [^\n]+@sha256:[0-9a-f]{64}/.test(dockerfile) && /snapshot\.debian\.org/.test(dockerfile),
    `${name}: base image and package source must be pinned`);
  const installer = readFileSync(path.join(templateRoot, '.devcontainer/install-tools.sh'), 'utf8');
  assert(/yazi_sha=[0-9a-f]{64}/.test(installer) && /codex_sha=[0-9a-f]{64}/.test(installer),
    `${name}: optional binaries must have architecture-specific SHA256 checks`);
  if (name === 'curated') assert(existsSync(path.join(templateRoot, '.devcontainer/tmux.conf')),
    'curated: missing the documented native tmux defaults');
  assert(!/(?:COPY|ADD)\s+.*\.(?:aibox|aibox-home)(?:\/|\s)/.test(dockerfile), `${name}: must not bake private aibox state into the image`);
  const allFiles = readdirSync(templateRoot, { recursive: true }).map(String);
  assert(!allFiles.some(file => file === '.aibox' || file.startsWith('.aibox/')),
    `${name}: Template must not create a default .aibox tree`);
  assert(!allFiles.some(file => file === '.aibox-home' || file.startsWith('.aibox-home/')),
    `${name}: Template must not create a default .aibox-home tree`);
}

for (const name of exampleNames) {
  const ignore = readFileSync(path.join(specRoot, 'examples', name, '.devcontainer/.gitignore'), 'utf8').split('\n');
  for (const ignored of ['*.env', 'compose.local.yaml'])
    assert(ignore.includes(ignored), `${name}: native starter must ignore ${ignored}`);
}

assert(existsSync(path.join(specRoot, 'examples/bind-home/.gitignore')), 'bind-home: missing private-path ignores');
const bindIgnore = readFileSync(path.join(specRoot, 'examples/bind-home/.gitignore'), 'utf8');
for (const ignored of ['.aibox-home/', '.env', '.devcontainer/compose.local.yaml'])
  assert(bindIgnore.split('\n').includes(ignored), `bind-home: missing ${ignored} ignore`);

assert(!Object.hasOwn({ ...examples.get('minimal') }, 'customizations'),
  'minimal: aibox customization must remain absent');
const validOptional = { schemaVersion: '1', workspace: { emphasis: 'standard', prompt: 'plain' } };
checkSchema(validOptional, customizationSchema, '$.validOptional');
for (const invalid of [
  { schemaVersion: '2' },
  { schemaVersion: '1', unknownPreference: true },
  { schemaVersion: '1', workspace: { theme: null } },
  { schemaVersion: '1', workspace: { theme: 'made-up-theme' } },
  { schemaVersion: '1', workspace: { tmux: { layout_switch: { style: 'dialog' } } } },
  { schemaVersion: '1', workspace: { tmux: { notifications: { protocol: 'host-shell' } } } },
]) {
  let rejected = false;
  try { checkSchema(invalid, customizationSchema, '$.negativeFixture'); } catch { rejected = true; }
  assert(rejected, `closed customization schema accepted invalid fixture ${JSON.stringify(invalid)}`);
}

console.log(`Native Template metadata/content and ${examples.size} rendered starter examples: passed.`);
