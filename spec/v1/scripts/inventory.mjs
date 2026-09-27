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
  if (/^latex\.preview\.(allow_public|bind|enabled|engine|port)$/.test(p)) return ['native-compose','compose.yaml services.latex-preview: service presence, image and ports; allow_public is an operator-policy request, never a local UX grant','AC-PREVIEW'];
  if (/permissions|execution|security\.|allow_public/.test(p)) return ['authority-request','Native harness/operator policy; project value is a request, never a grant','AC-SEC'];
  if (/^(local\.)/.test(p)) return ['private-native','Gitignored native environment/mount/harness configuration; never copy secrets to shared JSON','AC-MIG'];
  if (/^(customization|appearance)\./.test(p)) return ['aibox-ux',`customizations.aibox.workspace.${p.replace(/^(customization|appearance)\./,'')} (managed native output; user-local native override is separately owned)`,'AC-UX'];
  if (/^(context|processkit|skills|process)\.|\.agents\.|^agents\./.test(p)) return ['delegate-processkit','Optional processkit Feature/integration and processkit-owned configuration; no entity engine','AC-PK'];
  if (/addons\./.test(p)) return ['feature-options','features.<selected-feature> options; explicit disabled means absent; no new resolver','AC-TOOLS'];
  if (/mcp\./.test(p)) return ['native-integration','Native harness MCP configuration / optional processkit integration; scope shared vs private','AC-HARNESS'];
  if (/audio\./.test(p)) return ['audio','Audio Feature + environment/mounts; host preparation separately authorized','AC-AUDIO'];
  if (/^latex\./.test(p)) return ['latex',`customizations.aibox.${p}; LaTeX Feature/native Compose preview service`,'AC-PREVIEW'];
  if (/extra_volumes/.test(p)) return ['standard-mount','mounts or Compose volumes; preserve source/target/read-only and validate authority','AC-SEC'];
  if (/environment/.test(p)) return ['standard-env','containerEnv or Compose environment/env_file; local secrets remain private','AC-CONFIG'];
  if (/post_create_command/.test(p)) return ['standard-lifecycle','postCreateCommand; preserve ordering without host escalation','AC-CONFIG'];
  if (/keepalive/.test(p)) return ['runtime-helper','Container-local postStartCommand/helper; explicit opt-in','AC-CONFIG'];
  if (/resource_thresholds/.test(p)) return ['diagnostics',`customizations.aibox.diagnostics.${p.split('resource_thresholds.')[1]}; no host privilege implied`,'AC-DOCTOR'];
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
const objectSchema=()=>({type:'object',additionalProperties:false,properties:{}});
function valueSchema(type){
  const t=type.replace(/^Option<(.+)>$/,'$1');
  if(t==='bool')return {type:'boolean'};
  if(/^(u\d+|usize)$/.test(t))return {type:'integer',minimum:0};
  if(t==='Vec<String>')return {type:'array',items:{type:'string'}};
  if(/(?:HashMap|BTreeMap)<String, String>/.test(t))return {type:'object',additionalProperties:{type:'string'}};
  const named=types.get(t);
  if(named?.kind==='enum'){
    const values=[];let annotations='';
    for(const line of named.source.split('\n')){
      const variant=line.match(/^    ([A-Z]\w+),/);
      if(variant){
        const explicit=annotations.match(/serde\(rename\s*=\s*"([^"]+)"/);
        values.push(explicit?.[1]??kebab(variant[1]));annotations='';
      }else if(line.trim().startsWith('#['))annotations+=line;
    }
    if(values.length)return {type:'string',enum:values};
  }
  return {type:'string'};
}
function addSchemaPath(root,path,type){
  const parts=path.split('.');let current=root;
  for(let i=0;i<parts.length;i++){
    const isArray=parts[i].endsWith('[]'),key=parts[i].replace(/\[\]$/,'');
    if(i===parts.length-1)current.properties[key]=valueSchema(type);
    else {
      if(!current.properties[key])current.properties[key]=isArray?{type:'array',items:objectSchema()}:objectSchema();
      current=isArray?current.properties[key].items:current.properties[key];
    }
  }
}
const workspaceSchema=objectSchema(),latexSchema=objectSchema(),diagnosticsSchema=objectSchema();
for(const row of uniqueConfig){
  if(row.path.startsWith('customization.'))addSchemaPath(workspaceSchema,row.path.slice(14),row.type);
  if(row.path.startsWith('latex.')&&!/^latex\.preview\.(allow_public|bind|enabled|engine|port)$/.test(row.path))addSchemaPath(latexSchema,row.path.slice(6),row.type);
  if(row.path.startsWith('container.resource_thresholds.'))addSchemaPath(diagnosticsSchema,row.path.slice(30),row.type);
}
workspaceSchema.properties.sidebar={type:'object',additionalProperties:false,properties:{enabled:{type:'boolean'},width:{type:'integer',minimum:20},showUnknown:{type:'boolean'}}};
workspaceSchema.properties.review={type:'object',additionalProperties:false,properties:{enabled:{type:'boolean'},githubTui:{enum:['gh-dash','web']}}};
const harnessSchema={type:'object',additionalProperties:false,properties:{
  order:{type:'array',uniqueItems:true,items:{type:'string'}},
  launch:{type:'object',additionalProperties:{type:'object',additionalProperties:false,properties:{enabled:{type:'boolean'}}}}
}};
const customizationSchema={
  $schema:'https://json-schema.org/draft/2020-12/schema',
  $id:'https://github.com/projectious-work/aibox/blob/v1.x-dev/spec/v1/customization.schema.json',
  title:'aibox devcontainer.json customizations.aibox v1',type:'object',additionalProperties:false,
  required:['schemaVersion'],properties:{schemaVersion:{const:'1'},workspace:workspaceSchema,harnesses:harnessSchema,
    latex:latexSchema,diagnostics:diagnosticsSchema}
};

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
// A type declaration is not a command ledger: expand every public verb,
// resource discriminator and argument so removals cannot hide inside a blob.
const commandRoutes={
  Init:'native Template/copyable starter; no recurring config wizard',
  Apply:'split into build, rebuild, local refresh, or one-time conversion by effect',
  Up:'start_environment; interactive attachment is a separate CLI step',
  Emergency:'attach recovery path, bypassing tmux/Yazi/status',
  Prune:'scoped native cleanup with preview and fresh consent; no global prune',
  Down:'stop_environment', Get:'inspect_environment or versioned knowledge resource',
  Describe:'inspect_environment or versioned knowledge resource',
  Set:'ordinary file edit plus validate/refresh; no config mutation API',
  Edit:'ordinary editor plus how-to resource',
  Reset:'scoped recovery or processkit-owned workflow',
  Delete:'remove_environment or processkit-owned workflow, by resource',
  Doctor:'read-only check_environment',
  Create:'native scaffold or separately reviewed snapshot/backup utility',
  SelfCmd:'distribution package manager; help/version/completion remain CLI-only'
};
const resourceRoutes={
 OutputFormat:{Table:'human renderer',Json:'versioned JSON result',Yaml:'versioned YAML projection'},
 Layout:{Dev:'workspace.layout=dev',Focus:'workspace.layout=focus',Cowork:'workspace.layout=cowork',Ai:'workspace.layout=ai'},
 ApplyResource:{Audio:'external authorized host-audio setup',GeneratedRuntime:'local refresh for managed UX only',Migration:'processkit-owned migration workflow',Env:'snapshot switch utility with rollback'},
 GetResource:{Runtime:'inspect_environment',Addon:'versioned tool catalog',Env:'snapshot utility list',Kit:'processkit catalog',Skill:'processkit catalog',SkillCategory:'processkit catalog',Process:'processkit catalog',Migration:'processkit migration catalog'},
 DescribeResource:{Runtime:'inspect_environment',Addon:'versioned tool catalog',AddonCatalog:'versioned tool catalog',ImageProvenancePolicy:'signed artifact metadata',ProviderBackends:'versioned capability catalog',WorkspaceManifest:'effective config and provenance inspection',Env:'snapshot utility inspect',Kit:'processkit catalog',Skill:'processkit catalog',Process:'processkit catalog'},
 EditResource:{Config:'ordinary editor on devcontainer.json/native referenced files'},
 ResetResource:{Project:'scoped backup and recovery utility',Context:'processkit-owned recovery guidance'},
 DeleteResource:{Runtime:'remove_environment',Addon:'edit Feature references then validate',Skill:'processkit-owned configuration',Env:'snapshot utility delete with consent',Migration:'processkit-owned migration operation'},
 DoctorTarget:{Project:'read-only local/operator project checks',Audio:'read-only context-scoped audio checks',Security:'delegate scanners; report findings and scope'},
 PruneScope:{Safe:'scoped native cleanup preview',BuildCache:'native builder cache cleanup with exact scope',RuntimeHome:'scoped local cache cleanup',AgentWorktrees:'provider worktree cleanup with dirty protection',Containers:'owned-resource cleanup with exact identity',All:'retire broad aggregate; documented individually scoped actions'},
 CreateAction:{Env:'snapshot utility save',Backup:'scoped backup utility'},
 SelfAction:{Update:'distribution package manager update/check',Completion:'generated CLI shell completion',Uninstall:'distribution package manager uninstall; purge separately consented'}
};
const routeFor=(parent,name)=>parent==='Commands'?commandRoutes[name]:resourceRoutes[parent]?.[name];
const commandRows=[];
const argumentRows=[];
for(const t of commandTypes.filter(x=>x.kind==='enum')) {
  const declaration=t.declaration.split('\n');
  for(let i=1;i<declaration.length-1;i++) {
    const m=declaration[i].match(/^    (\w+)(?:\s*\{|,)$/);
    if(!m) continue;
    const start=i, name=m[1];
    let end=i+1;
    while(end<declaration.length-1 && !/^    \w+(?:\s*\{|,)$/.test(declaration[end])) end++;
    const body=declaration.slice(start,end).join('\n');
    const flags=[...body.matchAll(/#\[arg\((?:[^\]]|\n)*?\)\]/g)].map(x=>x[0]);
    const typeLine=Number(t.source.match(/#L(\d+)$/)?.[1]);
    commandRows.push({id:`CMD:${t.id.slice(8)}.${name}`,parent:t.id.slice(8),name,source:url('cli/src/cli.rs',typeLine+start),
      declaration:body,argument_annotations:flags,
      target:routeFor(t.id.slice(8),name),
      acceptance:'AC-CLI'});
    let annotation=[];
    for(let j=start+1;j<end;j++){
      const field=declaration[j].match(/^        (?:pub )?(\w+): ([^,]+),$/);
      if(field){
        argumentRows.push({id:`ARG:${t.id.slice(8)}.${name}.${field[1]}`,command:`CMD:${t.id.slice(8)}.${name}`,
          name:field[1],type:field[2],source:url('cli/src/cli.rs',typeLine+j),
          declaration:annotation.join('\n'),target:routeFor(t.id.slice(8),name),acceptance:'AC-CLI'});
        annotation=[];
      } else if(declaration[j].trim())annotation.push(declaration[j]);
    }
    i=end-1;
  }
}
const cliStruct=commandTypes.find(x=>x.id==='CMDTYPE:Cli');
const globalRoutes={config:'replace old aibox.toml locator with explicit devcontainer.json/workspace-folder locator',
 log_level:'stderr/structured log level only; never changes machine result',
 yes:'retire global consent bypass; destructive operations require fresh scoped authorization'};
for(const [index,line] of cliStruct.declaration.split('\n').entries()){
 const field=line.match(/^    pub (\w+): ([^,]+),$/);
 if(!field || field[1]==='command')continue;
 argumentRows.push({id:`ARG:Cli.${field[1]}`,command:'CMDTYPE:Cli',name:field[1],type:field[2],
  source:url('cli/src/cli.rs',Number(cliStruct.source.match(/#L(\d+)$/)?.[1])+index),
  declaration:cliStruct.declaration.split('\n').slice(Math.max(0,index-3),index+1).join('\n'),
  target:globalRoutes[field[1]],acceptance:'AC-CLI'});
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
const census={baseline,config_fields:uniqueConfig.length,addon_recipes:addons.length,addon_tools:addons.reduce((n,x)=>n+x.tools.length,0),cli_declarations:commandTypes.length,cli_actions:commandRows.length,cli_arguments:argumentRows.length,runtime_assets:runtime.length,documentation_pages:docs.length,environment_identifiers:env.size};
const outputs={
  'configuration.json':uniqueConfig,
  'configuration-types.json':[...types.values()].filter(t=>t.file.endsWith('config.rs')||['McpConfig','HarnessOverride'].includes(t.name)),
  'commands.json':commandTypes,'command-actions.json':commandRows,'command-arguments.json':argumentRows,'addons.json':addons,'runtime-assets.json':runtime,
  'defaults.json':defaults,'base-build.json':baseBuild,
  'documentation.json':docs,'environment.json':[...env].sort().map(([name,sources])=>({name,sources,classification:'source occurrence; distinguish internal/test override from supported user interface before migration',acceptance:'AC-CONFIG'})),
  'census.json':census,
  '../customization.schema.json':customizationSchema,
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
