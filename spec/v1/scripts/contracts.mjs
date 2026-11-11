// Generate the specification's executable contracts; no product code or runtime calls.
import {readFileSync,writeFileSync} from 'node:fs';
import {fileURLToPath} from 'node:url';
import path from 'node:path';
const root=path.resolve(path.dirname(fileURLToPath(import.meta.url)),'..');
const str={type:'string',minLength:1,maxLength:4096};
const id={type:'string',pattern:'^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$'};
const digest={type:'string',pattern:'^sha256:[0-9a-f]{64}$'};
const list=(items,maxItems=1000)=>({type:'array',items,maxItems});
const object=(properties,required=[])=>({type:'object',additionalProperties:false,properties,required});
const bool={type:'boolean'};
const stamp={type:'string',format:'date-time'};
const absolute={type:'string',pattern:'^/',maxLength:4096};
const base={schemaVersion:{const:'aibox.operation-request/v1'},operation:str,requestId:id,projectRoot:absolute};
const common={configPath:absolute,expectedInputDigest:digest};
const runtime={runtimeContext:id,resourceId:id};
const specs={
 build_environment:{fields:{...common,runtimeContext:id,frozenLockfile:bool,noCache:bool},required:['expectedInputDigest','runtimeContext']},
 start_environment:{fields:{...common,runtimeContext:id,frozenLockfile:bool},required:['expectedInputDigest','runtimeContext']},
 stop_environment:{fields:{...common,...runtime},required:['expectedInputDigest','runtimeContext','resourceId']},
 remove_environment:{fields:{...common,...runtime,acknowledgeDisruption:{const:true}},required:['expectedInputDigest','runtimeContext','resourceId','acknowledgeDisruption']},
 rebuild_environment:{fields:{...common,...runtime,acknowledgeDisruption:{const:true},frozenLockfile:bool,noCache:bool},required:['expectedInputDigest','runtimeContext','resourceId','acknowledgeDisruption']},
 inspect_environment:{fields:{...common,...runtime,effectiveConfig:bool},required:['runtimeContext']},
 check_environment:{fields:{...common,runtimeContext:id,offline:bool,checks:list(id,64)},required:[]},
 inspect_workspace:{fields:{...common,effectiveConfig:bool},required:[]},
 check_workspace:{fields:{...common,offline:bool,checks:list(id,64)},required:[]},
 refresh_workspace:{fields:{...common,scope:{enum:['all','theme','tmux']},theme:str,mode:{enum:['auto','dark','light']},layout:{enum:['ai','dev','focus','cowork']}},required:['expectedInputDigest','scope']},
 read_logs:{fields:{...common,...runtime,source:{enum:['aibox','runtime','session']},tailLines:{type:'integer',minimum:1,maximum:1000},since:stamp},required:['source','tailLines']},
 inspect_operation:{fields:{operationId:id},required:['operationId']},
 migrate_preview:{fields:{sourceRoot:absolute,destinationRoot:absolute},required:['sourceRoot','destinationRoot']},
 migrate_apply:{fields:{planId:id,expectedInputDigest:digest},required:['planId','expectedInputDigest']},
 migrate_rollback:{fields:{planId:id,expectedInputDigest:digest},required:['planId','expectedInputDigest']}
};
const schema=(name,title,body)=>({$schema:'https://json-schema.org/draft/2020-12/schema',$id:`https://github.com/projectious-work/aibox/blob/v1.x-dev/spec/v1/${name}.schema.json`,title,...body});
const request=schema('operation-request','Closed per-operation aibox requests',{
 oneOf:Object.entries(specs).map(([operation,{fields,required}])=>({
  ...object({...base,operation:{const:operation},...fields},[...Object.keys(base),...required]),
  ...(operation==='read_logs'?{allOf:[{if:{properties:{source:{const:'runtime'}}},then:{required:['runtimeContext','resourceId']}}]}:{})
 }))
});
const resource=object({kind:{enum:['container','image','network','volume','file']},id:str,role:{enum:['primary','sidecar','shared','artifact']},state:{enum:['running','stopped','absent','unknown','built']},owned:bool},['kind','id','role','state','owned']);
const finding=object({code:id,severity:{enum:['error','warning','info']},status:{enum:['passed','failed','skipped','unavailable','not_authorized']},message:str,source:str,path:str,observedAt:stamp,remediation:str},['code','severity','status','message','source','observedAt','remediation']);
const effect=object({resource:str,action:str,status:{enum:['confirmed','unknown']}},['resource','action','status']);
const next={enum:['none','edit_configuration','inspect_environment','inspect_operation','request_authorization','install_dependency','retry_read','manual_recovery','rebuild_environment','restart_session','read_guide']};
const configEntry=object({key:str,value:{type:['string','number','boolean','array','object','null']},redacted:bool,layer:{enum:['built_in','system','user','project','env_file','environment','invocation','native_override']},source:str,domain:{enum:['process','native','workspace','native_override']},invocationOnly:bool},['key','value','redacted','layer','source','domain','invocationOnly']);
const payloads={
 build_environment:object({images:list(resource,32)},['images']),
 inspect_environment:object({state:{enum:['absent','running','stopped','ambiguous','unknown']},resources:list(resource,64),configuration:list(configEntry,4096),observedAt:stamp},['state','resources','configuration','observedAt']),
 check_environment:object({findings:list(finding,4096),complete:bool},['findings','complete']),
 read_logs:object({source:{enum:['aibox','runtime','session']},entries:list(object({text:{type:'string',maxLength:16384},timestamp:stamp},['text'])),truncated:bool},['source','entries','truncated']),
 refresh_workspace:object({changedFiles:list(str),rebuildRequired:bool,sessionRestartRequired:bool},['changedFiles','rebuildRequired','sessionRestartRequired']),
 inspect_operation:object({operationId:id,originalOperation:{enum:Object.keys(specs)},state:{enum:['validated','authorized','executing','inspected','succeeded','no_change','failed','cancelled','partial','timed_out','refused','rebuild_required']},lastConfirmedStep:str,completedEffects:list(effect),unknownEffects:list(effect),nextAction:next},['operationId','originalOperation','state','lastConfirmedStep','completedEffects','unknownEffects','nextAction'])
};
for(const op of ['start_environment','stop_environment','remove_environment','rebuild_environment'])payloads[op]=object({resources:list(resource,64),retainedData:list(str,128)},['resources','retainedData']);
payloads.inspect_workspace=payloads.inspect_environment;
payloads.check_workspace=payloads.check_environment;
for(const op of ['migrate_preview','migrate_apply','migrate_rollback'])payloads[op]=object({planId:id,manifest:str,conflicts:list(finding),activated:bool},['planId','manifest','conflicts','activated']);
const success=['succeeded','no_change','rebuild_required'];
const result=schema('operation-result','aibox CLI/MCP operation result',{
 ...object({
  schemaVersion:{const:'aibox.operation-result/v1'},operation:{enum:[...Object.keys(specs),'invalid_request']},requestId:id,
  actors:object({initiator:str,executor:str},['initiator','executor']),
  target:object({scope:{enum:['local','operator']},workspaceRoot:absolute,runtimeContext:id,resourceId:id},['scope','workspaceRoot']),
  inputDigest:digest,outcome:{enum:[...success,'failed','cancelled','partial','refused','timed_out']},
  changedResources:list(str),warnings:list(str),evidence:list(str),data:{type:'object'},
  error:object({code:id,category:{enum:['invalid_input','denied','missing_dependency','incompatible','timeout','interrupted','child_failure','partial_failure','internal']},message:str,retryable:bool,nextAction:next},['code','category','message','retryable','nextAction']),
  completedEffects:list(effect),unknownEffects:list(effect)
 },['schemaVersion','operation','requestId','actors','outcome','changedResources','warnings','evidence']),
 allOf:[
  {if:{properties:{outcome:{enum:success}}},then:{required:['data','target','inputDigest'],not:{required:['error']}}},
  {if:{properties:{outcome:{enum:['failed','partial','refused','timed_out','cancelled']}}},then:{required:['error']}},
  {if:{properties:{outcome:{const:'partial'}}},then:{required:['completedEffects','unknownEffects']}},
  {if:{properties:{operation:{const:'invalid_request'}}},then:{properties:{outcome:{enum:['failed','refused']}}}},
  ...Object.keys(specs).map(operation=>({if:{properties:{operation:{const:operation}}},then:{properties:{data:{$ref:`#/$defs/${operation}`}}}}))
 ],$defs:payloads
});
const policy=schema('operator-policy','External operator policy and digest-bound grants',object({
 schemaVersion:{const:'aibox.operator-policy/v1'},principal:id,
 allowedRoots:list(absolute,128),executables:object({devcontainer:absolute,runtime:absolute,compose:absolute},['devcontainer','runtime']),
 runtimeContexts:list(object({name:id,kind:{enum:['docker','podman']},endpointFingerprint:digest},['name','kind','endpointFingerprint']),32),
 operations:{...list({enum:Object.keys(specs).filter(x=>!x.endsWith('_workspace'))},32),uniqueItems:true},
 grants:list(object({id:id,requester:id,projectRoot:absolute,configPath:absolute,operation:{enum:Object.keys(specs)},inputDigest:digest,runtimeContext:id,resourceId:id,expiresAt:stamp,allowHostHooks:bool},['id','requester','projectRoot','configPath','operation','inputDigest','runtimeContext','expiresAt']),256)
},['schemaVersion','principal','allowedRoots','executables','runtimeContexts','operations','grants']));
const operationRecord=schema('operation-record','Durable local aibox operation record',object({
 schemaVersion:{const:'aibox.operation-record/v1'},operationId:id,requestFingerprint:digest,operation:{enum:Object.keys(specs)},
 projectRoot:absolute,configPath:absolute,inputDigest:digest,policyDigest:digest,grantId:id,
 createdAt:stamp,updatedAt:stamp,executor:str,state:payloads.inspect_operation.properties.state,lastConfirmedStep:str,
 resources:list(resource,64),completedEffects:list(effect),unknownEffects:list(effect),
 result:{$ref:result.$id}
},['schemaVersion','operationId','requestFingerprint','operation','projectRoot','inputDigest','createdAt','updatedAt','executor','state','lastConfirmedStep','resources','completedEffects','unknownEffects']));
const outputs={'operation-request.schema.json':request,'operation-result.schema.json':result,'operator-policy.schema.json':policy,'operation-record.schema.json':operationRecord};
for(const [file,value] of Object.entries(outputs)){
 const content=JSON.stringify(value,null,2)+'\n',destination=path.join(root,file);
 if(process.argv.includes('--check')){if(readFileSync(destination,'utf8')!==content)throw Error(`Contract drift: ${file}`);}
 else writeFileSync(destination,content);
}
console.log(`Contract generation checked: ${Object.keys(specs).length} operations, policy and operation records.`);
