// Dependency-free specification consistency checks. JSON is a YAML 1.2 subset.
import {readFileSync,readdirSync,existsSync} from 'node:fs';
import path from 'node:path';
import {fileURLToPath} from 'node:url';
const root=path.resolve(path.dirname(fileURLToPath(import.meta.url)),'..');
const read=p=>readFileSync(path.join(root,p),'utf8');
const json=p=>JSON.parse(read(p));
const assert=(v,m)=>{if(!v)throw Error(m);};
const docs=readdirSync(root).filter(x=>x.endsWith('.md'));
const normative=docs.filter(x=>/^(0[1-9]|1[0-4])-/.test(x)).map(read).join('\n');
const criteria=new Set([...read('05-acceptance.md').matchAll(/^### (AC-[A-Z]+) /gm)].map(x=>x[1]));
const featureIds=new Set([...read('05-acceptance.md').matchAll(/^\| (F\d{2}) \|/gm)].map(x=>x[1]));
const tracedIds=[...read('feature-trace.md').matchAll(/^\| (F\d{2}) \|/gm)].map(x=>x[1]);
assert(new Set(tracedIds).size===tracedIds.length,'Duplicate feature trace');
assert(featureIds.size===tracedIds.length && tracedIds.every(x=>featureIds.has(x)),'Feature trace coverage');
for(const m of normative.matchAll(/\bAC-[A-Z]+\b/g))assert(criteria.has(m[0]),`Missing acceptance ${m[0]}`);
const defs=new Set([...normative.matchAll(/\*\*(R-[A-Z]+):\*\*/g)].map(x=>x[1]));
for(const m of normative.matchAll(/\bR-[A-Z]+\b/g))assert(defs.has(m[0]),`Missing requirement ${m[0]}`);
for(const id of defs)assert(read('06-delivery.md').includes(`| ${id} |`),`Missing trace ${id}`);
for(const file of ['configuration.json','commands.json','command-actions.json','command-arguments.json','runtime-assets.json']){
  const rows=json('ledger/'+file),ids=new Set();
  for(const r of rows){assert(!ids.has(r.id),`Duplicate ${r.id}`);ids.add(r.id);assert(criteria.has(r.acceptance),`No acceptance ${r.id}`);}
}
const census=json('ledger/census.json');
assert(json('ledger/configuration.json').length===census.config_fields,'Config census');
const addons=json('ledger/addons.json');
assert(addons.length===41 && addons.length===census.addon_recipes,'Addon census');
assert(addons.reduce((n,a)=>n+a.tools.length,0)===census.addon_tools,'Tool census');
assert(json('ledger/commands.json').length===census.cli_declarations,'Command census');
assert(json('ledger/command-actions.json').length===census.cli_actions,'Action census');
assert(json('ledger/command-arguments.json').length===census.cli_arguments,'Argument census');
for(const row of [...json('ledger/command-actions.json'),...json('ledger/command-arguments.json')])assert(row.target,'Unresolved command route '+row.id);
assert(json('ledger/runtime-assets.json').length===census.runtime_assets,'Asset census');
const roadmap=json('roadmap.yaml');
const roadmapSchema=json('roadmap.schema.json');
// Validate the declared, intentionally small schema vocabulary without an
// implicit online npm dependency. Unknown schema keywords fail closed.
function schemaCheck(value,schema,at='$',rootSchema=schema){
 const known=new Set(['$schema','$id','title','type','additionalProperties','required','properties','items','minItems','minLength','minimum','pattern','enum','const','uniqueItems','$defs','$ref','allOf','if','then']);
 for(const key of Object.keys(schema))assert(known.has(key),`Unsupported schema keyword ${key}`);
 if(schema.$ref){assert(schema.$ref.startsWith('#/$defs/'),'External schema reference');return schemaCheck(value,rootSchema.$defs[schema.$ref.slice(8)],at,rootSchema);}
 if(schema.type){const actual=Array.isArray(value)?'array':value===null?'null':Number.isInteger(value)?'integer':typeof value;assert(actual===schema.type,`${at}: expected ${schema.type}`);}
 if(schema.const!==undefined)assert(value===schema.const,`${at}: const`);
 if(schema.enum)assert(schema.enum.includes(value),`${at}: enum`);
 if(schema.minLength!==undefined)assert(value.length>=schema.minLength,`${at}: minLength`);
 if(schema.minimum!==undefined)assert(value>=schema.minimum,`${at}: minimum`);
 if(schema.pattern)assert(new RegExp(schema.pattern).test(value),`${at}: pattern`);
 if(schema.minItems!==undefined)assert(value.length>=schema.minItems,`${at}: minItems`);
 if(schema.uniqueItems)assert(new Set(value.map(JSON.stringify)).size===value.length,`${at}: uniqueItems`);
 if(schema.required)for(const key of schema.required)assert(Object.hasOwn(value,key),`${at}: missing ${key}`);
 if(schema.properties){
  if(schema.additionalProperties===false)for(const key of Object.keys(value))assert(Object.hasOwn(schema.properties,key),`${at}: unknown ${key}`);
  else if(typeof schema.additionalProperties==='object')for(const key of Object.keys(value))if(!Object.hasOwn(schema.properties,key))schemaCheck(value[key],schema.additionalProperties,`${at}.${key}`,rootSchema);
  for(const [key,child] of Object.entries(schema.properties))if(Object.hasOwn(value,key))schemaCheck(value[key],child,`${at}.${key}`,rootSchema);
 }
 else if(typeof schema.additionalProperties==='object')for(const [key,child] of Object.entries(value))schemaCheck(child,schema.additionalProperties,`${at}.${key}`,rootSchema);
 if(schema.items)for(let i=0;i<value.length;i++)schemaCheck(value[i],schema.items,`${at}[${i}]`,rootSchema);
 if(schema.allOf)for(const rule of schema.allOf){
  const [key,condition]=Object.entries(rule.if.properties)[0];
  if(value[key]===condition.const)schemaCheck(value,rule.then,at,rootSchema);
 }
}
schemaCheck(roadmap,roadmapSchema);
const customizationSchema=json('customization.schema.json');
const customizationExample={schemaVersion:'1',workspace:{theme:'gruvbox',mode:'dark',layout:'dev',tmux:{status:{mode:'extended',refresh:{'interval-seconds':5}}}}};
schemaCheck(customizationExample,customizationSchema);
for(const invalid of [
 {schemaVersion:'1',workspace:{theme:'unknown-palette'}},
 {schemaVersion:'1',workspace:{tmux:{unexpected:true}}},
 {schemaVersion:'1',workspace:{theme:null}},
 {schemaVersion:'1',workspace:{tmux:{status:{refresh:{'interval-seconds':-1}}}}}
]){
 let rejected=false;try{schemaCheck(invalid,customizationSchema);}catch{rejected=true;}
 assert(rejected,'Customization negative fixture accepted');
}
for(const row of json('ledger/configuration.json').filter(x=>x.path.startsWith('customization.'))){
 let current=customizationSchema.properties.workspace;
 for(const segment of row.path.slice(14).split('.')){
  const key=segment.replace(/\[\]$/,'');
  assert(current.properties?.[key],`Missing customization schema key ${row.path}`);
  current=current.properties[key];if(segment.endsWith('[]'))current=current.items;
 }
}
const resultSchema=json('operation-result.schema.json');
const requestSchema=json('operation-request.schema.json');
const requestExample={schemaVersion:'aibox.operation-request/v1',operation:'remove_environment',requestId:'example-1',projectRoot:'/workspace',expectedInputDigest:'sha256:'+'0'.repeat(64),runtimeContext:'test',resourceId:'abc'};
schemaCheck(requestExample,requestSchema);
for(const invalid of [
 {...requestExample,operation:'shell'},
 {...requestExample,resourceId:''},
 {...requestExample,expectedInputDigest:'not-a-digest'},
 {...requestExample,resourceId:undefined},
 {...requestExample,unbounded:true}
]){
 let rejected=false;try{schemaCheck(invalid,requestSchema);}catch{rejected=true;}
 assert(rejected,'Request negative fixture accepted');
}
const resultExample={schemaVersion:'aibox.operation-result/v1',operation:'check_environment',requestId:'example-1',
 actors:{initiator:'local-user',executor:'aibox-local'},target:{scope:'local',workspaceRoot:'/workspace'},
 inputDigest:'sha256:'+'0'.repeat(64),outcome:'no_change',changedResources:[],warnings:[],evidence:[]};
schemaCheck(resultExample,resultSchema);
for(const invalid of [
 {...resultExample,outcome:'healthy'},
 {...resultExample,actors:{initiator:'local-user'}},
 {...resultExample,inputDigest:'not-a-digest'},
 {...resultExample,outcome:'failed'},
 {...resultExample,target:{scope:'operator',workspaceRoot:'/workspace'}},
 {...resultExample,secret:'leak'}
]){
 let rejected=false;try{schemaCheck(invalid,resultSchema);}catch{rejected=true;}
 assert(rejected,'Result negative fixture accepted');
}
const phases=roadmap.groups.flatMap(x=>x.phases),ids=new Set(phases.map(x=>x.id));
assert(new Set(roadmap.groups.map(x=>x.id)).size===roadmap.groups.length,'Duplicate roadmap group');
assert(ids.size===phases.length,'Duplicate phase');
for(const p of phases){
 assert(p.title && p.summary,'Missing roadmap fields');
 assert(['idea','planned','in_progress','shipped','cancelled'].includes(p.status),'Unknown roadmap status');
 for(const dep of p.dependencies??[])assert(ids.has(dep),`Unknown dependency ${dep}`);
 assert(p.spec?.length>0,`Missing spec references ${p.id}`);
 for(const ref of p.spec)assert(existsSync(path.join(root,ref)),`Missing phase spec ${p.id}: ${ref}`);
 if(p.status==='shipped')assert(p.release && p.devNote && existsSync(path.join(root,p.devNote)),`Missing shipped evidence ${p.id}`);
}
const visited=new Set(),active=new Set();
const byId=new Map(phases.map(x=>[x.id,x]));
function visit(id){assert(!active.has(id),`Roadmap cycle at ${id}`);if(visited.has(id))return;active.add(id);for(const dep of byId.get(id).dependencies??[])visit(dep);active.delete(id);visited.add(id);}
for(const id of ids)visit(id);
// Every local Markdown link, including linked source evidence ledgers.
for(const file of [...docs,'ledger/configuration.md','ledger/addons.md']){
 for(const m of read(file).matchAll(/\]\(([^)]+)\)/g)){
  const link=m[1];if(/^(https?:|#)/.test(link))continue;
  assert(existsSync(path.resolve(root,path.dirname(file),link.split('#')[0])),`Broken link ${file}: ${link}`);
 }
}
console.log(`Specification checks passed: ${defs.size} requirements, ${criteria.size} acceptance groups, ${phases.length} phases; ledger counts and local links consistent.`);
