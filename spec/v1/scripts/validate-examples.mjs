// Static native Dev Container examples only. V1-04 owns runtime build/up/exec
// qualification and publication of Dev Container Templates.
import { readFileSync } from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
const readJson = relative => JSON.parse(readFileSync(path.join(root, relative), 'utf8'));
// Maintained examples use whole-line JSONC comments for inactive options.
// Read only active fields; comments are documentation, not configuration.
const readExample = relative => JSON.parse(readFileSync(path.join(root, relative), 'utf8')
  .split('\n').filter(line => !line.trimStart().startsWith('//')).join('\n'));
const assert = (condition, message) => {
  if (!condition) throw new Error(message);
};

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
  if (schema.required) {
    for (const key of schema.required) assert(Object.hasOwn(value, key), `${at}: missing ${key}`);
  }
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
  if (schema.items) {
    for (const [index, child] of value.entries()) checkSchema(child, schema.items, `${at}[${index}]`, rootSchema);
  }
  if (schema.allOf) {
    for (const rule of schema.allOf) {
      const [key, condition] = Object.entries(rule.if.properties)[0];
      if (value[key] === condition.const) checkSchema(value, rule.then, at, rootSchema);
    }
  }
}

const schema = readJson('customization.schema.json');
const minimal = readExample('examples/minimal/.devcontainer/devcontainer.json');
const customized = readExample('examples/customized/.devcontainer/devcontainer.json');

for (const [name, example] of [['minimal', minimal], ['customized', customized]]) {
  assert(example.image, `${name}: expected a native image declaration`);
  assert(example.remoteUser, `${name}: expected an explicit remoteUser`);
  assert(Array.isArray(example.mounts) && example.mounts.length > 0, `${name}: expected native persistent-home mount`);
  assert(!Object.hasOwn(example, 'features'), `${name}: example must not invent unverified Feature references`);
  assert(!Object.hasOwn(example, 'build'), `${name}: example must not claim an unverified image build`);
  assert(!Object.hasOwn(example, 'aibox'), `${name}: aibox preferences belong under customizations.aibox`);
  assert(!JSON.stringify(example).includes('.aibox-home'),
    `${name}: starter example must not require .aibox-home`);
  assert(!JSON.stringify(example).includes('.aibox/'),
    `${name}: starter example must not require .aibox`);
}

assert(!Object.hasOwn(minimal, 'customizations'), 'minimal: aibox customization must remain optional');
assert(customized.customizations?.aibox, 'customized: expected optional customizations.aibox example');
checkSchema(customized.customizations.aibox, schema, '$.customizations.aibox');

const validOptional = { schemaVersion: '1', workspace: { emphasis: 'standard', prompt: 'plain' } };
checkSchema(validOptional, schema, '$.validOptional');
for (const invalid of [
  { schemaVersion: '2' },
  { schemaVersion: '1', unknownPreference: true },
  { schemaVersion: '1', workspace: { theme: null } },
  { schemaVersion: '1', workspace: { theme: 'made-up-theme' } },
  { schemaVersion: '1', workspace: { tmux: { layout_switch: { style: 'dialog' } } } },
  { schemaVersion: '1', workspace: { tmux: { notifications: { protocol: 'host-shell' } } } },
]) {
  let rejected = false;
  try { checkSchema(invalid, schema, '$.negativeFixture'); } catch { rejected = true; }
  assert(rejected, `closed customization schema accepted invalid fixture ${JSON.stringify(invalid)}`);
}

console.log('Static native Dev Container examples: passed.');
console.log('V1-04 runtime build/up/exec qualification and template publication: deferred.');
