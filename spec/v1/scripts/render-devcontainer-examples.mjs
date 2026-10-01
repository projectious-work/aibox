// Render the example JSONC option catalog from the closed aibox schema.
// This is documentation generation, not a Dev Container configuration parser.
import { readFileSync, writeFileSync, mkdirSync } from 'node:fs';
import { fileURLToPath } from 'node:url';
import path from 'node:path';

const specRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
const schema = JSON.parse(readFileSync(path.join(specRoot, 'customization.schema.json'), 'utf8'));
const check = process.argv.includes('--check');

const native = [
  ['build', { type: 'object', properties: { dockerfile: { type: 'string' }, context: { type: 'string' } } },
    'Build from a project Dockerfile instead of image; choose exactly one source'],
  ['dockerComposeFile', { type: 'array', items: { type: 'string' } },
    'Select ordered native Compose files instead of a single-image definition'],
  ['service', { type: 'string' }, 'Name the primary service when using dockerComposeFile'],
  ['runServices', { type: 'array', items: { type: 'string' } },
    'Choose Compose services started with the primary service'],
  ['containerUser', { type: 'string' }, 'Choose the user for all container processes; align its home and mounts'],
  ['workspaceFolder', { type: 'string' }, 'Choose the project path opened inside the container'],
  ['workspaceMount', { type: 'string' }, 'Override the project mount; coordinate with workspaceFolder'],
  ['forwardPorts', { type: 'array', items: { type: 'integer' } },
    'Forward selected development ports to the client without publishing them'],
  ['containerEnv', { type: 'object', additionalProperties: { type: 'string' } },
    'Set non-secret variables for all container processes; rebuild after changes'],
  ['remoteEnv', { type: 'object', additionalProperties: { type: 'string' } },
    'Set non-secret variables for connected tools and their child processes'],
  ['features', { type: 'object', additionalProperties: { type: 'object' } },
    'Install selected, pinned Features; their own manifests define valid option keys'],
  ['initializeCommand', { type: 'string' },
    'Run a host-side initialization command; requires explicit operator review'],
  ['postCreateCommand', { type: 'string' }, 'Run a command inside the newly created container'],
  ['postStartCommand', { type: 'string' }, 'Run a command inside the container after each start'],
  ['shutdownAction', { enum: ['none', 'stopContainer', 'stopCompose'] },
    'Choose what supporting clients stop when they disconnect'],
  ['image', { type: 'string' }, 'Select a base image instead of build or dockerComposeFile; pin a release by digest'],
  ['mounts', { type: 'array', items: { type: 'string' } },
    'Persist the aibox user home with a named volume, or use a reviewed native bind'],
  ['name', { type: 'string' }, 'Give the Dev Container a human-readable display name'],
  ['remoteUser', { type: 'string' }, 'Choose the user for terminals and supporting tools'],
  ['customizations', { type: 'object', properties: { aibox: schema } },
    'Optional aibox-owned UX intent; inert until the aibox runtime Feature is installed'],
];

// Keep related native choices together while the active keys remain valid JSONC.
const nativeGroups = [
  ['Definition and identity', ['name', 'image', 'build', 'dockerComposeFile', 'service', 'runServices']],
  ['Users, workspace and retained data', ['containerUser', 'remoteUser', 'workspaceFolder', 'workspaceMount', 'mounts']],
  ['Tools, environment and connectivity', ['features', 'containerEnv', 'remoteEnv', 'forwardPorts']],
  ['Lifecycle and shutdown', ['initializeCommand', 'postCreateCommand', 'postStartCommand', 'shutdownAction']],
  ['Optional aibox workspace preferences', ['customizations']],
];
const nativeByKey = new Map(native.map(field => [field[0], field]));
if (nativeGroups.flatMap(([, keys]) => keys).length !== native.length ||
    new Set(nativeGroups.flatMap(([, keys]) => keys)).size !== native.length ||
    nativeGroups.some(([, keys]) => keys.some(key => !nativeByKey.has(key))))
  throw Error('Native option groups must contain each documented field exactly once');

const descriptions = {
  'devcontainer.build.dockerfile': 'Select the project Dockerfile path relative to this JSONC file',
  'devcontainer.build.context': 'Select the native Docker build-context directory',
  'aibox.schemaVersion': 'Select the closed aibox customization schema',
  'aibox.workspace': 'Group optional aibox workspace presentation settings',
  'aibox.workspace.emphasis': 'Choose how strongly semantic UI colors are emphasized',
  'aibox.workspace.emphasis_overrides': 'Override semantic emphasis by named UI element',
  'aibox.workspace.layout': 'Choose the initial workspace pane layout',
  'aibox.workspace.mode': 'Choose light, dark or automatic color mode',
  'aibox.workspace.prompt': 'Choose the shipped Starship prompt preset',
  'aibox.workspace.theme': 'Choose the shipped theme family',
  'aibox.workspace.variant': 'Choose a variant offered by the selected theme',
  'aibox.workspace.tmux': 'Group tmux session, status, notification and title settings',
  'aibox.workspace.tmux.layout': 'Choose the tmux pane layout',
  'aibox.workspace.tmux.prefix': 'Choose the tmux prefix key',
  'aibox.workspace.tmux.session_name': 'Name the managed tmux session',
  'aibox.workspace.tmux.layout_switch': 'Configure the interactive layout switcher',
  'aibox.workspace.tmux.layout_switch.confirm': 'Require confirmation before changing layout',
  'aibox.workspace.tmux.layout_switch.enabled': 'Enable the interactive layout switcher',
  'aibox.workspace.tmux.layout_switch.prefix-key': 'Bind the layout switcher key after the tmux prefix',
  'aibox.workspace.tmux.layout_switch.style': 'Choose the layout switcher display style',
  'aibox.workspace.tmux.notifications': 'Configure agent-attention notifications',
  'aibox.workspace.tmux.notifications.enabled': 'Enable agent-attention notifications',
  'aibox.workspace.tmux.notifications.include-message': 'Include the bounded agent message in notifications',
  'aibox.workspace.tmux.notifications.protocol': 'Choose notification delivery protocol supported by the runtime',
  'aibox.workspace.tmux.notifications.states': 'Select agent states that trigger a notification',
  'aibox.workspace.tmux.status': 'Configure the managed tmux status line',
  'aibox.workspace.tmux.status.elements': 'Select which tmux status segments appear',
  'aibox.workspace.tmux.status.elements.aibox-metrics': 'Select individual aibox diagnostic segments',
  'aibox.workspace.tmux.status.forge': 'Configure forge status sources',
  'aibox.workspace.tmux.status.forge.github-hosts': 'Select GitHub hosts queried by the forge status segment',
  'aibox.workspace.tmux.status.labels': 'Override labels shown in tmux status',
  'aibox.workspace.tmux.status.mode': 'Choose extended, plain or disabled tmux status',
  'aibox.workspace.tmux.status.model_providers': 'Configure supported model-provider status probes',
  'aibox.workspace.tmux.status.model_providers.providers': 'Select model providers shown in status',
  'aibox.workspace.tmux.status.layout': 'Place selected segments in the two status lines',
  'aibox.workspace.tmux.status.refresh': 'Configure status refresh and cache durations',
  'aibox.workspace.tmux.status.separators': 'Configure visual separators between status segments',
  'aibox.workspace.tmux.theme_switch': 'Configure interactive theme switching',
  'aibox.workspace.tmux.title': 'Configure tmux-owned terminal titles',
  'aibox.workspace.tmux.title.states': 'Configure text shown for each agent state',
  'aibox.workspace.sidebar': 'Configure the optional read-only harness sidebar',
  'aibox.workspace.sidebar.enabled': 'Show the optional harness sidebar',
  'aibox.workspace.sidebar.width': 'Set sidebar width in terminal columns',
  'aibox.workspace.sidebar.showUnknown': 'Show explicitly unavailable metrics in the sidebar',
  'aibox.workspace.review': 'Configure the optional local diff and PR-review workflow',
  'aibox.workspace.review.enabled': 'Enable local review shortcuts',
  'aibox.workspace.review.githubTui': 'Choose the optional GitHub PR viewer or web fallback',
  'aibox.harnesses': 'Group selected agent-harness launch preferences',
  'aibox.harnesses.order': 'Order selected harness panes at session launch',
  'aibox.harnesses.launch': 'Enable or disable launch for named installed harnesses',
  'aibox.latex': 'Group local LaTeX build and preview preferences',
  'aibox.latex.cache_dir': 'Choose the project-local LaTeX cache directory',
  'aibox.latex.documents': 'List named LaTeX source and output paths',
  'aibox.latex.engine': 'Choose the LaTeX engine',
  'aibox.latex.options': 'Pass explicit arguments to the selected LaTeX engine',
  'aibox.latex.preview': 'Configure local PDF preview selection',
  'aibox.latex.preview.document': 'Select the default document for preview',
  'aibox.diagnostics': 'Group thresholds for read-only workspace diagnostics',
  'aibox.diagnostics.memory_mib_warn': 'Warn above this cgroup memory usage in MiB',
  'aibox.diagnostics.oom_kill_warn': 'Warn above this observed OOM-kill count',
  'aibox.diagnostics.process_count_warn': 'Warn above this live process count',
  'aibox.diagnostics.processkit_mcp_python_warn': 'Warn above this ProcessKit MCP Python process count',
};

const dynamicOptions = {
  'aibox.workspace.emphasis_overrides': [
    ['aibox.workspace.emphasis_overrides.<element>', 'Override one named semantic element', 'string emphasis choice supported by the renderer'],
  ],
  'aibox.workspace.tmux.status.labels': [
    ['aibox.workspace.tmux.status.labels.<segment>', 'Set the label of one documented status segment', 'JSON string'],
  ],
  'aibox.harnesses.launch': [
    ['aibox.harnesses.launch.<installed-harness>.enabled', 'Start an installed harness in its pane', 'true | false'],
  ],
  'aibox.latex.documents': [
    ['aibox.latex.documents[].name', 'Name one document for build and preview selection', 'JSON string'],
    ['aibox.latex.documents[].source', 'Select its TeX source relative to the project', 'JSON string path'],
    ['aibox.latex.documents[].output_dir', 'Select its separate output directory', 'JSON string path'],
  ],
};

const nativeValid = {
  'devcontainer.dockerComposeFile': 'relative path or ordered array of relative paths',
  'devcontainer.forwardPorts': 'array of port numbers or "host:port" strings',
  'devcontainer.mounts': 'array of native mount strings or mount objects',
  'devcontainer.features': 'map of selected Feature IDs to their declared option objects',
  'devcontainer.initializeCommand': 'string, argument array or named-command object',
  'devcontainer.postCreateCommand': 'string, argument array or named-command object',
  'devcontainer.postStartCommand': 'string, argument array or named-command object',
};

function purposeFor(pointer) {
  if (descriptions[pointer]) return descriptions[pointer];
  if (/\.status\.elements\.aibox-metrics\.[^.]+$/.test(pointer))
    return 'Show this aibox diagnostic metric in tmux status';
  if (/\.status\.elements\.[^.]+$/.test(pointer))
    return 'Show this tmux status segment';
  if (/\.status\.labels\.[^.]+$/.test(pointer))
    return 'Set the visible label for this status segment';
  if (/\.status\.layout\.line[12]-(left|right)$/.test(pointer))
    return 'List segment IDs at this position in tmux status';
  if (/\.status\.refresh\.[^.]+$/.test(pointer))
    return 'Set this refresh or cache duration';
  if (/\.status\.model_providers\.[^.]+$/.test(pointer))
    return 'Control this model-provider status probe';
  if (/\.status\.separators\.[^.]+$/.test(pointer))
    return 'Choose this visual separator setting';
  if (/\.theme_switch\.[^.]+$/.test(pointer))
    return 'Control this interactive theme-switcher setting';
  if (/\.title\.states\.[^.]+$/.test(pointer))
    return 'Set terminal title text for this agent state';
  if (/\.title\.[^.]+$/.test(pointer))
    return 'Control this tmux-owned title setting';
  throw Error(`Missing purpose for ${pointer}`);
}

function valid(schema) {
  if ('const' in schema) return JSON.stringify(schema.const);
  if (schema.enum) return schema.enum.map(JSON.stringify).join(' | ');
  if (schema.type === 'boolean') return 'true | false';
  if (schema.type === 'integer') return `integer >= ${schema.minimum ?? 0}`;
  if (schema.type === 'array') return `array of ${schema.items?.type ?? 'values'}`;
  if (schema.type === 'object') return schema.additionalProperties ? 'object with named entries' : 'object';
  return 'JSON string; use a value supported by the owning runtime';
}

function example(pointer, schema) {
  const samples = {
    'devcontainer.build.dockerfile': 'Dockerfile',
    'devcontainer.build.context': '..',
    'devcontainer.dockerComposeFile': ['compose.yaml'],
    'devcontainer.service': 'workspace',
    'devcontainer.runServices': ['workspace'],
    'devcontainer.containerUser': 'vscode',
    'devcontainer.workspaceFolder': '/workspaces/project',
    'devcontainer.workspaceMount': 'source=${localWorkspaceFolder},target=/workspaces/project,type=bind',
    'devcontainer.forwardPorts': [3000],
    'devcontainer.containerEnv': { EXAMPLE: 'value' },
    'devcontainer.remoteEnv': { EXAMPLE: 'value' },
    'devcontainer.features': {},
    'devcontainer.initializeCommand': 'echo reviewed-host-command',
    'devcontainer.postCreateCommand': 'echo created',
    'devcontainer.postStartCommand': 'echo started',
    'aibox.workspace.emphasis_overrides': { selected_element: 'standard' },
    'aibox.workspace.tmux.notifications.states': ['question', 'error'],
    'aibox.workspace.tmux.status.forge.github-hosts': ['github.com'],
    'aibox.workspace.tmux.status.model_providers.providers': ['openai'],
    'aibox.harnesses.order': ['codex'],
    'aibox.harnesses.launch': { codex: { enabled: true } },
    'aibox.latex.documents': [{ name: 'main', source: 'docs/main.tex', output_dir: '.latex-cache/main' }],
    'aibox.latex.options': [],
  };
  if (pointer in samples) return samples[pointer];
  if ('const' in schema) return schema.const;
  if (schema.enum) return schema.enum[0];
  if (schema.type === 'boolean') return false;
  if (schema.type === 'integer') return schema.minimum ?? 0;
  if (schema.type === 'array') return [];
  if (schema.type === 'object') return {};
  return 'example';
}

function fieldsFor(schema) { return Object.entries(schema.properties ?? {}); }

function emitField(lines, key, child, pointer, indent, selected, activeValues, isLastActive, description) {
  const active = selected.has(pointer);
  const prefix = active ? '' : '// ';
  const pad = ' '.repeat(indent);
  const doc = description ?? purposeFor(pointer);
  lines.push(`${pad}// ${doc}. Valid: ${nativeValid[pointer] ?? valid(child)}.`);
  for (const [, purpose, values] of dynamicOptions[pointer] ?? [])
    lines.push(`${pad}// ${purpose}. Valid: ${values}.`);
  if (child.properties) {
    lines.push(`${pad}${prefix}${JSON.stringify(key)}: {`);
    emitFields(lines, fieldsFor(child), pointer, indent + 2, selected, activeValues);
    lines.push(`${pad}${prefix}}${active && isLastActive ? '' : ','}`);
  } else {
    const value = active ? activeValues[pointer] : example(pointer, child);
    lines.push(`${pad}${prefix}${JSON.stringify(key)}: ${JSON.stringify(value)}${active && isLastActive ? '' : ','}`);
  }
}

function emitFields(lines, fields, parent, indent, selected, activeValues) {
  const ordered = [...fields].sort(([a], [b]) => Number(selected.has(`${parent}.${a}`)) - Number(selected.has(`${parent}.${b}`)));
  const activeKeys = ordered.filter(([key]) => selected.has(`${parent}.${key}`)).map(([key]) => key);
  for (const [index, [key, child, description]] of ordered.entries()) {
    if (index > 0 && (child.properties || ordered[index - 1][1].properties)) lines.push('');
    const pointer = `${parent}.${key}`;
    emitField(lines, key, child, pointer, indent, selected, activeValues,
      key === activeKeys.at(-1), description);
  }
}

function render(flavor) {
  const customized = flavor === 'customized';
  const customUser = flavor === 'custom-user';
  const bindHome = flavor === 'bind-home';
  const username = customUser ? 'dev' : 'aibox';
  const activeValues = {
    'devcontainer.name': customized ? 'Native Dev Container with Optional aibox Preferences' : 'Native Dev Container',
    'devcontainer.build.dockerfile': 'Dockerfile',
    'devcontainer.build.context': '.',
    'devcontainer.containerUser': username,
    'devcontainer.remoteUser': username,
    'devcontainer.mounts': bindHome
      ? [{source: '${localWorkspaceFolder}/.aibox-home', target: '/home/aibox', type: 'bind'}]
      : [`source=aibox-home-\${devcontainerId},target=/home/${username},type=volume`],
    'aibox.schemaVersion': '1',
    'aibox.workspace.theme': 'gruvbox',
    'aibox.workspace.mode': 'dark',
    'aibox.workspace.layout': 'dev',
  };
  const selected = new Set(['devcontainer.name', 'devcontainer.build', 'devcontainer.build.dockerfile', 'devcontainer.build.context', 'devcontainer.containerUser', 'devcontainer.remoteUser', 'devcontainer.mounts']);
  if (customized) for (const pointer of ['devcontainer.customizations', 'devcontainer.customizations.aibox',
    'aibox.schemaVersion', 'aibox.workspace', 'aibox.workspace.theme', 'aibox.workspace.mode', 'aibox.workspace.layout'])
    selected.add(pointer);
  const lines = [
    '// aibox v1 example: standard Dev Container JSON with Comments (JSONC).',
    '// Uncomment only settings needed by this project; add a comma before the next active key.',
    '// image, build and dockerComposeFile are alternative native definitions.',
    '// aibox does not parse or transform native build/up fields; devcontainer.cli owns them.',
    '// Native options below are limited to aibox journeys; upstream documents all others.',
    '// Optional aibox keys come from spec/v1/customization.schema.json, not a parallel build schema.',
    '// For unconstrained strings, comments state schema-valid types; check runtime support before activation.',
    '// Feature-specific options come from each selected Feature manifest, not this file.',
    '{',
  ];
  const activeNativeKeys = nativeGroups.flatMap(([, keys]) => keys)
    .filter(key => selected.has(`devcontainer.${key}`));
  for (const [heading, keys] of nativeGroups) {
    lines.push('', `  // --- ${heading} ---`);
    for (const key of keys) {
      const [, child, description] = nativeByKey.get(key);
      const pointer = `devcontainer.${key}`;
      if (key === 'customizations') {
        lines.push('  // Product-specific metadata; optional for aibox UX. Valid: object.');
        const active = selected.has(pointer);
        lines.push(`  ${active ? '' : '// '}\"customizations\": {`);
        lines.push('    // aibox-owned UX preferences; inert without its runtime Feature. Valid: object.');
        lines.push(`    ${active ? '' : '// '}\"aibox\": {`);
        emitFields(lines, fieldsFor(schema), 'aibox', 6, selected, activeValues);
        lines.push(`    ${active ? '' : '// '}}`);
        lines.push(`  ${active ? '' : '// '}}`);
      } else {
        emitField(lines, key, child, pointer, 2, selected, activeValues,
          key === activeNativeKeys.at(-1), description);
      }
    }
  }
  lines.push('}', '');
  return lines.join('\n');
}

for (const flavor of ['minimal', 'customized', 'custom-user', 'bind-home']) {
  const file = path.join(specRoot, 'examples', flavor, '.devcontainer', 'devcontainer.json');
  const expected = render(flavor);
  // Every emitted example is also strict JSON once its documentation comments
  // are removed. This checks our active keys without interpreting native fields.
  const activeJSON = expected.split('\n').filter(line => !line.trimStart().startsWith('//')).join('\n');
  const active = JSON.parse(activeJSON);
  if (flavor === 'minimal' && 'customizations' in active)
    throw Error('The minimal example must remain native-only');
  if (flavor === 'customized' && active.customizations?.aibox?.schemaVersion !== '1')
    throw Error('The customized example must select the current aibox schema');
  if (check) {
    if (readFileSync(file, 'utf8') !== expected) throw Error(`Example option catalog drift: ${file}`);
  } else { mkdirSync(path.dirname(file), {recursive: true}); writeFileSync(file, expected); }
  const template = flavor === 'customized' ? 'curated' : 'minimal';
  const sourceRoot = path.resolve(specRoot, '../../templates/src', template, '.devcontainer');
  for (const name of ['Dockerfile', 'install-tools.sh', '.gitignore', ...(template === 'curated' ? ['tmux.conf'] : [])]) {
    let content = readFileSync(path.join(sourceRoot, name), 'utf8');
    if (flavor === 'custom-user' && name === 'Dockerfile') content = content.replace('ARG USERNAME=aibox', 'ARG USERNAME=dev');
    const destination = path.join(path.dirname(file), name);
    if (check) {
      if (readFileSync(destination, 'utf8') !== content) throw Error(`Example native file drift: ${destination}`);
    } else writeFileSync(destination, content);
  }
  if (flavor === 'bind-home') {
    const ignore = path.join(specRoot, 'examples', flavor, '.gitignore');
    const contents = '.aibox-home/\n.env\n.devcontainer/compose.local.yaml\n';
    if (check) { if (readFileSync(ignore, 'utf8') !== contents) throw Error(`Example private-path ignore drift: ${ignore}`); }
    else writeFileSync(ignore, contents);
  }
}
console.log('Dev Container example option catalogs checked.');
