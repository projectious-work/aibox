// Dependency-free specification consistency checks. JSON is a YAML 1.2 subset.
import {readFileSync,readdirSync,existsSync} from 'node:fs';
import path from 'node:path';
import {fileURLToPath} from 'node:url';
const root=path.resolve(path.dirname(fileURLToPath(import.meta.url)),'..');
const read=p=>readFileSync(path.join(root,p),'utf8');
const json=p=>JSON.parse(read(p));
const assert=(v,m)=>{if(!v)throw Error(m);};
const docs=readdirSync(root).filter(x=>x.endsWith('.md'));
const normative=docs.filter(x=>/^0[1-6]-/.test(x)).map(read).join('\n');
const criteria=new Set([...read('05-acceptance.md').matchAll(/^### (AC-[A-Z]+) /gm)].map(x=>x[1]));
for(const m of normative.matchAll(/\bAC-[A-Z]+\b/g))assert(criteria.has(m[0]),`Missing acceptance ${m[0]}`);
const defs=new Set([...normative.matchAll(/\*\*(R-[A-Z]+):\*\*/g)].map(x=>x[1]));
for(const m of normative.matchAll(/\bR-[A-Z]+\b/g))assert(defs.has(m[0]),`Missing requirement ${m[0]}`);
for(const id of defs)assert(read('06-delivery.md').includes(`| ${id} |`),`Missing trace ${id}`);
for(const file of ['configuration.json','commands.json','runtime-assets.json']){
  const rows=json('ledger/'+file),ids=new Set();
  for(const r of rows){assert(!ids.has(r.id),`Duplicate ${r.id}`);ids.add(r.id);assert(criteria.has(r.acceptance),`No acceptance ${r.id}`);}
}
const census=json('ledger/census.json');
assert(json('ledger/configuration.json').length===census.config_fields,'Config census');
const addons=json('ledger/addons.json');
assert(addons.length===41 && addons.length===census.addon_recipes,'Addon census');
assert(addons.reduce((n,a)=>n+a.tools.length,0)===census.addon_tools,'Tool census');
assert(json('ledger/commands.json').length===census.cli_declarations,'Command census');
assert(json('ledger/runtime-assets.json').length===census.runtime_assets,'Asset census');
const roadmap=json('roadmap.yaml');
assert(roadmap.schemaVersion==='aibox.spec-roadmap/v1','Roadmap schema');
const phases=roadmap.groups.flatMap(x=>x.phases),ids=new Set(phases.map(x=>x.id));
assert(ids.size===phases.length,'Duplicate phase');
for(const p of phases){
 assert(p.title && p.summary,'Missing roadmap fields');
 assert(['idea','planned','in_progress','shipped','cancelled'].includes(p.status),'Unknown roadmap status');
 for(const dep of p.dependencies??[])assert(ids.has(dep),`Unknown dependency ${dep}`);
 if(p.status==='shipped')assert(p.release && p.devNote && existsSync(path.join(root,p.devNote)),`Missing shipped evidence ${p.id}`);
}
// Every local Markdown link, including linked source evidence ledgers.
for(const file of [...docs,'ledger/configuration.md','ledger/addons.md']){
 for(const m of read(file).matchAll(/\]\(([^)]+)\)/g)){
  const link=m[1];if(/^(https?:|#)/.test(link))continue;
  assert(existsSync(path.resolve(root,path.dirname(file),link.split('#')[0])),`Broken link ${file}: ${link}`);
 }
}
console.log(`Specification checks passed: ${defs.size} requirements, ${criteria.size} acceptance groups, ${phases.length} phases; ledger counts and local links consistent.`);
