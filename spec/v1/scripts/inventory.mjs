// Specification evidence generator, not a product implementation.
// Reads only explicitly scoped paths from the immutable v0 release commit.
import {execFileSync} from 'node:child_process';
import {readFileSync, writeFileSync, mkdirSync} from 'node:fs';
import {createHash} from 'node:crypto';
import {fileURLToPath} from 'node:url';
import path from 'node:path';

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../../..');
const baseline = '9061bf76d2a1f6ac3c7093d8c605bcc4de7eddf6';
const git = (...args) => execFileSync('git', args, {cwd: root, encoding: 'utf8', maxBuffer: 16e6});
const show = p => git('show', `${baseline}:${p}`);
const files = (...p) => git('ls-tree', '-r', '--name-only', baseline, ...p).trim().split('\n').filter(Boolean);
const url = (p, n) => `https://github.com/projectious-work/aibox/blob/${baseline}/${p}${n ? '#L'+n : ''}`;
const sha = s => createHash('sha256').update(s).digest('hex');
const kebab = s => s.replace(/([a-z0-9])([A-Z])/g, '$1-$2').replaceAll('_','-').toLowerCase();
const types = new Map();
for (const file of ['cli/src/config.rs','cli/src/mcp_registration.rs']) {
  const lines = show(file).split('\n');
  for (let i=0;i<lines.length;i++) {
    const m = lines[i].match(/^(?:pub(?:\(crate\))? )?(struct|enum) (\w+) \{$/);
    if (!m) continue;
    let end=i+1; while (end<lines.length && lines[end]!=='}') end++;
    let start=i; while(start>0 && /^(#\[|\/\/\/|\s*$)/.test(lines[start-1])) start--;
    const attrs=lines.slice(start,i).join('\n');
    const fields=[]; let pending=[];
    for(let j=i+1;j<end;j++) {
      const f=lines[j].match(/^    (?:pub(?:\(crate\))? )?(\w+): (.+),$/);
      if(f) {
        const annotations=pending.join('\n');
        const rename=annotations.match(/rename\s*=\s*"([^"]+)"/);
        fields.push({name:f[1],key:rename?.[1] ?? (attrs.includes('rename_all = "kebab-case"')?kebab(f[1]):f[1]),type:f[2],annotations,line:j+1,aliases:[...annotations.matchAll(/alias\s*=\s*"([^"]+)"/g)].map(x=>x[1])});
        pending=[];
      } else pending.push(lines[j]);
    }
    types.set(m[2],{name:m[2],kind:m[1],file,line:i+1,fields,source:lines.slice(start,end+1).join('\n'),url:url(file,i+1)});
  }
}

function mapping(p) {
  if (/permissions|execution|security\.|allow_public/.test(p)) return ['authority-request','Native harness/operator policy; project value is a request, never a grant','AC-SEC'];
  if (/^(local\.)/.test(p)) return ['private-native','Gitignored native environment/mount/harness configuration; never copy secrets to shared JSON','AC-MIG'];
  if (/customization\.|appearance\./.test(p)) return ['aibox-ux','customizations.aibox.workspace + native user overrides; preserve exact choice','AC-UX'];
  if (/^(context|processkit|skills|process)\.|\.agents\.|^agents\./.test(p)) return ['delegate-processkit','Optional processkit Feature/integration and processkit-owned configuration; no entity engine','AC-PK'];
  if (/addons\./.test(p)) return ['feature-options','features.<selected-feature> options; explicit disabled means absent; no new resolver','AC-TOOLS'];
  if (/mcp\./.test(p)) return ['native-integration','Native harness MCP configuration / optional processkit integration; scope shared vs private','AC-HARNESS'];
  if (/audio\./.test(p)) return ['audio','Audio Feature + environment/mounts; host preparation separately authorized','AC-AUDIO'];
  if (/latex\./.test(p)) return ['latex','LaTeX Feature + customizations.aibox.latex; native Compose preview service','AC-PREVIEW'];
  if (/extra_volumes/.test(p)) return ['standard-mount','mounts or Compose volumes; preserve source/target/read-only and validate authority','AC-SEC'];
  if (/environment/.test(p)) return ['standard-env','containerEnv or Compose environment/env_file; local secrets remain private','AC-CONFIG'];
  if (/post_create_command/.test(p)) return ['standard-lifecycle','postCreateCommand; preserve ordering without host escalation','AC-CONFIG'];
  if (/keepalive/.test(p)) return ['runtime-helper','Container-local postStartCommand/helper; explicit opt-in','AC-CONFIG'];
  if (/resource_thresholds/.test(p)) return ['diagnostics','customizations.aibox.diagnostics thresholds; no host privilege implied','AC-DOCTOR'];
  if (/container\.paths/.test(p)) return ['native-path','Native devcontainer build/dockerComposeFile paths and referenced assets; explicit migration','AC-MIG'];
  if (/image\.|aibox\.(base|version)/.test(p)) return ['standard-image','image or build/Dockerfile FROM with immutable release inputs','AC-BUILD'];
  if (/container\.user/.test(p)) return ['standard-user','containerUser/remoteUser/updateRemoteUserUID; test file ownership','AC-LIFE'];
  if (/container\.hostname/.test(p)) return ['runtime-hostname','Native Compose hostname or supported runtime arguments','AC-CONFIG'];
  if (/name$/.test(p) && !p.startsWith('ai.')) return ['identity','devcontainer name plus native resource identity; display name is not deletion authority','AC-LIFE'];
  if (/^ai\./.test(p)) return ['harness','Features for install/version; customizations.aibox.harnesses for enable/order; native config','AC-HARNESS'];
  if (/^apply\./.test(p)) return ['state-retention','Preserve disabled-harness state by default; purge is explicit recovery action','AC-MIG'];
  if (/integrations\./.test(p)) return ['native-integration','Native Git/gh credential helper configuration','AC-HARNESS'];
  return ['metadata-migration','Map version/profile metadata into Dev Container artifact identity or aibox schemaVersion; retire only redundant envelope fields','AC-CONFIG'];
}
const config=[];
function walk(type,prefix,stack=[]) {
  type=type.replace(/^crate::mcp_registration::/,'');
  const t=types.get(type);
  if(!t || t.kind!=='struct') throw Error(`unhandled type ${type}`);
  for(const f of t.fields) {
    if (/serde\(skip\)/.test(f.annotations) || f.name==='legacy_theme') continue;
    let key=f.key;
    if(type==='AddonsSection' && key==='addons') key='{addon}';
    const p=prefix ? `${prefix}.${key}` : key;
    const inner=f.type.match(/(?:<|, )([A-Z]\w+)>$/)?.[1] ?? f.type.replace(/^crate::mcp_registration::/,'');
    const container=f.type.startsWith('Vec<')?'[]':/(?:HashMap|BTreeMap)</.test(f.type)?'.{key}':'';
    const base={path:p,type:f.type,aliases:f.aliases,annotations:f.annotations,source:url(t.file,f.line)};
    const child=types.get(inner);
    if(child?.kind==='struct' && !stack.includes(inner)) {
      walk(inner,p+container,[...stack,type]);
    } else {
      const [disposition,target,acceptance]=mapping(p);
      config.push({...base,id:`CFG:${p}`,disposition,target,acceptance,default_evidence:t.url,enum_evidence:child?.kind==='enum'?child.url:null});
    }
  }
}
walk('AiboxConfig',''); walk('AiboxLocalConfig','local');
// Custom deserializer inputs not represented by the effective AiSection.
for(const p of ['ai.harnesses[].harness','ai.harnesses[].enabled','ai.harnesses[].install','ai.harnesses[].version','ai.execution.{harness}.filesystem','ai.execution.{harness}.approval','ai.execution.{harness}.network']) {
  const [disposition,target,acceptance]=mapping(p);
  config.push({id:`CFG:${p}`,path:p,type:'custom deserializer input',disposition,target,acceptance,source:types.get('RawAiSection').url,default_evidence:types.get('RawAiExecutionSection').url});
}
// Flattened recursive addon group keys use another AddonToolsSection, not a literal groups table.
for(const r of config) {
  r.path=r.path.replace('addons.{addon}.{key}','addons.{addon}');
  if(r.path.includes('.groups')) { r.path='addons.{addon}.{group...}.tools.{tool}.{enabled|version}';r.target+='; recursively flattened group selection'; }
  r.id=`CFG:${r.path}`;
}
const uniqueConfig=[...new Map(config.map(x=>[x.id,x])).values()].sort((a,b)=>a.path.localeCompare(b.path));

const addonFiles=files('addons').filter(x=>x.endsWith('.yaml'));
const addons=addonFiles.map(file=>{
  const s=show(file); const header=s.split(/\n(?:builder|runtime):/)[0];
  const tools=[];
  for(const m of header.matchAll(/^  - name: (.+)\n([\s\S]*?)(?=^  - name:|^\S|$(?![\s\S]))/gm)) {
    tools.push({name:m[1].trim(),declaration:m[0].trimEnd(),acceptance:'AC-TOOLS',target:'Upstream Feature/package where suitable; otherwise bounded aibox Feature; preserve enabled/version options'});
  }
  for(const m of header.matchAll(/^  - \{ name: ([^,]+),[^\n]+/gm)) tools.push({name:m[1].trim(),declaration:m[0],acceptance:'AC-TOOLS',target:'Upstream Feature/package where suitable; otherwise bounded aibox Feature; preserve enabled/version options'});
  if(!tools.length) throw Error(`No tools parsed: ${file}`);
  return {id:`ADDON:${file.split('/').at(-1).slice(0,-5)}`,source:url(file),sha256:sha(s),metadata:header.trim(),tools};
});
const cli=show('cli/src/cli.rs');
const commandTypes=[];
for(const m of cli.matchAll(/^(?:pub )?(enum|struct) (\w+) \{\n([\s\S]*?)^\}/gm)) {
  commandTypes.push({id:`CMDTYPE:${m[2]}`,kind:m[1],source:url('cli/src/cli.rs',cli.slice(0,m.index).split('\n').length),declaration:m[0],acceptance:'AC-CLI'});
}
const runtimeFiles=files('images','cli/src/templates','cli/assets','addons').filter(x=>!x.endsWith('/AGENTS.md'));
const runtime=runtimeFiles.map(file=>({id:`ASSET:${file}`,source:url(file),sha256:sha(show(file)),disposition:'retain behavior; reuse native tool; asset-by-asset equivalence before replacement',acceptance:/theme|tmux|starship/.test(file)?'AC-UX':/yazi|preview|latex/.test(file)?'AC-PREVIEW':'AC-TOOLS'}));
const baseBuild=files('images').filter(p=>/(Dockerfile|entrypoint.*\.sh)$/.test(p)).map(file=>({source:url(file),sha256:sha(show(file)),declaration:show(file),acceptance:'AC-TOOLS'}));
const defaults=[];
for(const file of ['cli/src/config.rs','cli/src/mcp_registration.rs']) {
  const source=show(file);
  for(const m of source.matchAll(/^(?:pub )?fn (default_\w+|bool_true|bool_false)\([^\n]*\)[^\n]*\{\n[\s\S]*?^\}|^impl Default for \w+ \{\n[\s\S]*?^\}/gm)) defaults.push({source:url(file,source.slice(0,m.index).split('\n').length),declaration:m[0]});
}
const docFiles=files('docs-site/content/docs');
const docs=docFiles.filter(p=>p.endsWith('.md')).map(file=>({source:url(file),sha256:sha(show(file)),headings:show(file).split('\n').filter(l=>/^#{1,4} /.test(l))}));
const envFiles=files('cli/src','images').filter(p=>/\.(rs|sh)$/.test(p));
const env=new Map();
for(const file of envFiles) {
  for(const m of show(file).matchAll(/\b(AIBOX_[A-Z0-9_]+)\b/g)) {
    const v=env.get(m[1])??new Set();v.add(url(file));env.set(m[1],v);
  }
}
const census={baseline,config_fields:uniqueConfig.length,addon_recipes:addons.length,addon_tools:addons.reduce((n,x)=>n+x.tools.length,0),cli_declarations:commandTypes.length,runtime_assets:runtime.length,documentation_pages:docs.length,environment_identifiers:env.size};
const outputs={
  'configuration.json':uniqueConfig,
  'configuration-types.json':[...types.values()].filter(t=>t.file.endsWith('config.rs')||['McpConfig','HarnessOverride'].includes(t.name)),
  'commands.json':commandTypes,'addons.json':addons,'runtime-assets.json':runtime,
  'defaults.json':defaults,'base-build.json':baseBuild,
  'documentation.json':docs,'environment.json':[...env].sort().map(([name,sources])=>({name,sources,classification:'source occurrence; distinguish internal/test override from supported user interface before migration',acceptance:'AC-CONFIG'})),
  'census.json':census,
};
const esc=s=>String(s).replaceAll('|','\\|').replaceAll('\n',' ');
outputs['configuration.md']='# v0 configuration field ledger\n\nGenerated from the immutable baseline by `scripts/inventory.mjs`. Wildcards denote maps/arrays. Types, annotations, aliases and enum/default source evidence are preserved in the JSON companions. Targets are proposed destinations, not proven mappings. See ../02-configuration.md for custom-deserializer and authority exceptions.\n\n| v0 path | Type | Disposition / v1 destination | Acceptance |\n|---|---|---|---|\n'+uniqueConfig.map(r=>`| [\`${esc(r.path)}\`](${r.source}) | \`${esc(r.type)}\` | ${esc(r.disposition)}: ${esc(r.target)} | ${r.acceptance} |`).join('\n')+'\n';
outputs['addons.md']='# v0 tool and dependency ledger\n\nEvery recipe/tool below retains its enabled/version choices; runtime build dependencies and pins are preserved in each source-linked recipe. These are baseline pins, not selected v1 pins. All rows require AC-TOOLS; see ../04-reuse-targets.md.\n\n| Recipe | Tool | Baseline declaration |\n|---|---|---|\n'+addons.flatMap(a=>a.tools.map(t=>`| [${a.id.slice(6)}](${a.source}) | ${esc(t.name)} | ${esc(t.declaration)} |`)).join('\n')+'\n';
const out=path.join(root,'spec/v1/ledger');
const check=process.argv.includes('--check');
if(!check)mkdirSync(out,{recursive:true});
for(const [name,value] of Object.entries(outputs)) {
  const body=typeof value==='string'?value:JSON.stringify(value,null,2)+'\n';
  const dest=path.join(out,name);
  if(check) {if(readFileSync(dest,'utf8')!==body)throw Error(`Inventory drift: ${name}`);}
  else writeFileSync(dest,body);
}
console.log(JSON.stringify(census,null,2));
