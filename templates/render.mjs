// Apply local Template sources using the upstream option substitution contract.
// Source distribution helper only; lifecycle operations belong to devcontainer.cli.
import { readFileSync, writeFileSync, mkdirSync, readdirSync, lstatSync, renameSync, rmSync } from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const source = path.join(path.dirname(fileURLToPath(import.meta.url)), 'src');
const args = process.argv.slice(2);
const selected = {};
for (let i = 0; i < args.length; i += 2) {
  const key = args[i];
  if (!['--template', '--output', '--harness'].includes(key) || !args[i + 1] || selected[key] !== undefined)
    throw Error('Usage: node templates/render.mjs --template minimal|curated --output /absolute/project [--harness codex|none]');
  selected[key] = args[i + 1];
}
if (!['minimal', 'curated'].includes(selected['--template']) || !path.isAbsolute(selected['--output'] ?? ''))
  throw Error('Select minimal or curated and an absolute output project directory');
const root = selected['--output'];
if (!lstatSync(root).isDirectory() || lstatSync(root).isSymbolicLink()) throw Error('Output must be an existing directory');
const template = path.join(source, selected['--template']);
const metadata = JSON.parse(readFileSync(path.join(template, 'devcontainer-template.json'), 'utf8'));
const harness = selected['--harness'] ?? metadata.options.harness.default;
if (!metadata.options.harness.enum.includes(harness)) throw Error('Unsupported harness selection');
const files = [];
function collect(directory, relative = '') {
  for (const entry of readdirSync(directory, {withFileTypes: true})) {
    const file = path.join(directory, entry.name);
    const output = path.join(relative, entry.name);
    if (entry.isDirectory()) collect(file, output);
    else if (entry.isFile()) {
      const content = readFileSync(file, 'utf8').replaceAll('${templateOption:harness}', harness);
      if (content.includes('${templateOption:')) throw Error(`Unresolved Template option in ${output}`);
      files.push({output, content, mode: lstatSync(file).mode & 0o777});
    } else throw Error(`Unsupported source file ${file}`);
  }
}
collect(path.join(template, '.devcontainer'));
// mkdir reserves the final path and refuses any existing configuration,
// including a symlink. Stage content first; failed writes leave no partial tree.
const staging = path.join(root, `.devcontainer-template-${process.pid}`);
mkdirSync(staging, {mode: 0o700});
try {
  for (const {output, content, mode} of files) {
    const file = path.join(staging, output);
    mkdirSync(path.dirname(file), {recursive: true});
    writeFileSync(file, content, {mode, flag: 'wx'});
  }
  // Reserve the destination so rename cannot replace an existing empty directory.
  const destination = path.join(root, '.devcontainer');
  mkdirSync(destination);
  try { renameSync(staging, destination); }
  catch (error) { rmSync(destination, {recursive: true}); throw error; }
} finally { rmSync(staging, {recursive: true, force: true}); }
console.log(`Applied ${metadata.id} ${metadata.version} with harness=${harness} to ${root}`);
