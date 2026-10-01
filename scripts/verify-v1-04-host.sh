#!/usr/bin/env bash
# Host-only direct Dev Container CLI qualification for V1-04.
# All containers, volumes, temporary workspaces, and explicitly tagged build
# images are identified by this run and removed individually on exit.
set -Eeuo pipefail

repo_root="$(git rev-parse --show-toplevel)"
cd "$repo_root"
cli_version='0.89.0'
output_dir=''
while (($#)); do
  case "$1" in
    --output)
      [[ $# -ge 2 && -n "$2" ]] || { printf '%s\n' '--output requires a directory' >&2; exit 2; }
      output_dir="$2"; shift 2 ;;
    -h|--help)
      cat <<'HELP'
Usage: scripts/verify-v1-04-host.sh [--output DIR]

Builds and runs the checked-in minimal/curated Template outputs plus the
custom-user and bind-home examples through pinned @devcontainers/cli 0.89.0.
Requires Node/npm and a working Docker-compatible engine. Evidence is retained
at DIR (which must not already exist), or in a unique /tmp directory.
Only containers, named home volumes, and uniquely tagged images created by this
run are removed. Temporary project directories are script-owned and deleted.
HELP
      exit 0 ;;
    *) printf 'Unknown argument: %s\n' "$1" >&2; exit 2 ;;
  esac
done

if [[ -n "$(git status --porcelain --untracked-files=all)" ]]; then
  printf 'V1-04 host candidate worktree must be clean before evidence capture.\n' >&2
  exit 2
fi
source_commit="$(git rev-parse HEAD)"
started_at="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
engine=''
if command -v docker >/dev/null 2>&1; then engine="$(command -v docker)"
elif command -v podman >/dev/null 2>&1; then engine="$(command -v podman)"
else printf '%s\n' 'A Docker-compatible CLI (docker or podman) is required.' >&2; exit 2
fi
"$engine" info >/dev/null 2>&1 || { printf 'Container engine is not ready: %s info failed.\n' "$engine" >&2; exit 2; }
command -v node >/dev/null 2>&1 && command -v npx >/dev/null 2>&1 || {
  printf '%s\n' 'Node.js and npm/npx are required.' >&2; exit 2;
}

if [[ -z "$output_dir" ]]; then
  output_dir="$(mktemp -d /tmp/aibox-v1-04-host-evidence.XXXXXX)"
else
  [[ "$output_dir" = /* ]] || { printf '%s\n' '--output must be absolute' >&2; exit 2; }
  [[ ! -e "$output_dir" ]] || { printf '%s\n' '--output must not already exist' >&2; exit 2; }
  mkdir -p -- "$output_dir"
fi
chmod 700 "$output_dir"

work_dir="$(mktemp -d /tmp/aibox-v1-04-host-run.XXXXXX)"
run_id="$(date -u +%Y%m%dT%H%M%SZ)-$$-${RANDOM}${RANDOM}"
steps_file="$output_dir/steps.tsv"
: > "$steps_file"
container_ids_file="$work_dir/container-ids"
: > "$container_ids_file"
home_volumes_file="$work_dir/home-volumes"
: > "$home_volumes_file"
observed_images_file="$work_dir/observed-images"
: > "$observed_images_file"
cleanup_failures_file="$output_dir/cleanup-failures.txt"
: > "$cleanup_failures_file"
cleanup_log="$output_dir/cleanup.log"
: > "$cleanup_log"

record_step() {
  local label="$1" log="$2"; shift 2
  local status
  printf '%s\tstarted\t%s\t%s\n' "$label" "$(basename "$log")" "$*" >> "$steps_file"
  set +e
  "$@" > "$log" 2>&1
  status=$?
  set -e
  printf '%s\t%s\t%s\t%s\n' "$label" "$status" "$(basename "$log")" "$*" >> "$steps_file"
  if ((status != 0)); then
    printf 'Host step failed: %s (exit %s; see %s)\n' "$label" "$status" "$log" >&2
    return "$status"
  fi
}

cleanup() {
  local rc=$?
  trap - EXIT
  local final_rc="$rc"
  if [[ -n "$engine" ]] && "$engine" info >/dev/null 2>&1; then
    local id inspect_output volume
    while IFS= read -r id; do
      [[ -n "$id" ]] || continue
      inspect_output="$work_dir/cleanup-inspect-$id.json"
      if "$engine" inspect "$id" > "$inspect_output" 2>/dev/null; then
        volume="$(node - "$inspect_output" <<'NODE'
const fs=require('node:fs');const x=JSON.parse(fs.readFileSync(process.argv[2],'utf8'))[0];
const m=(x.Mounts||[]).find(v=>v.Type==='volume' && v.Destination==='/home/aibox' || v.Type==='volume' && v.Destination==='/home/dev');
if(m)process.stdout.write(m.Name||'');
NODE
        )"
        if [[ "$volume" == aibox-home-* ]]; then printf '%s\n' "$volume" >> "$home_volumes_file"; fi
        if ! "$engine" rm -f "$id" >> "$cleanup_log" 2>&1; then printf 'container %s\n' "$id" >> "$cleanup_failures_file"; fi
      fi
    done < <({
      cat "$container_ids_file"
      for template in minimal curated custom-user bind-home; do
        "$engine" ps -aq --filter "label=aibox.v1-04.run=$run_id-$template" 2>/dev/null || true
      done
    } | sort -u)
    while IFS= read -r volume; do
      if [[ "$volume" == aibox-home-* ]] && ! "$engine" volume rm "$volume" >> "$cleanup_log" 2>&1; then
        printf 'volume %s\n' "$volume" >> "$cleanup_failures_file"
      fi
    done < <(sort -u "$home_volumes_file")
    for template in minimal curated custom-user bind-home; do
      image="aibox-v1-04-$run_id-$template"
      if "$engine" image inspect "$image" >/dev/null 2>&1 && ! "$engine" image rm "$image" >> "$cleanup_log" 2>&1; then
        printf 'image %s\n' "$image" >> "$cleanup_failures_file"
      fi
    done
  else
    printf 'container engine unavailable during cleanup\n' >> "$cleanup_failures_file"
  fi
  if ! cp "$observed_images_file" "$output_dir/observed-image-refs.txt"; then printf 'observed image references\n' >> "$cleanup_failures_file"; fi
  if ! rm -rf -- "$work_dir"; then printf 'temporary directory %s\n' "$work_dir" >> "$cleanup_failures_file"; fi
  if [[ -s "$cleanup_failures_file" ]]; then
    printf 'Cleanup left run-owned resources; see %s\n' "$cleanup_failures_file" >&2
    final_rc=1
  fi
  node - "$output_dir" "$source_commit" "$started_at" "$run_id" "$cli_version" "$engine" "$rc" "$final_rc" "$output_dir/observed-image-refs.txt" <<'NODE' || true
const fs=require('node:fs'),path=require('node:path'),crypto=require('node:crypto');
const [dir,sourceCommit,startedAt,runId,version,engine,exitCode,finalExitCode,retainedFile]=process.argv.slice(2);
const latest=new Map();
for(const line of fs.readFileSync(path.join(dir,'steps.tsv'),'utf8').split('\n').filter(Boolean)){
  const [name,status,log,command]=line.split('\t');
  latest.set(name,{name,status:status==='started'?'incomplete':Number(status)===0?'passed':'failed',log,command});
}
const rows=[...latest.values()];
const files=fs.readdirSync(dir).filter(name=>name!=='evidence.json').map(name=>({path:name,
  sha256:'sha256:'+crypto.createHash('sha256').update(fs.readFileSync(path.join(dir,name))).digest('hex')}));
const retained=fs.readFileSync(retainedFile,'utf8').split('\n').filter(Boolean)
  .filter(image=>!image.startsWith(`aibox-v1-04-${runId}-`));
const cleanupFailures=fs.readFileSync(path.join(dir,'cleanup-failures.txt'),'utf8').split('\n').filter(Boolean);
const evidence={phase:'V1-04',verification:'host',runId,startedAt,finishedAt:new Date().toISOString(),sourceCommit,
  devcontainerCli:{package:`@devcontainers/cli@${version}`,version},engine,exitCode:Number(exitCode),finalExitCode:Number(finalExitCode),steps:rows,artifacts:files,
  cleanup:{scope:'only run-labeled containers, their inspected aibox-home volumes, uniquely tagged build outputs, and the script temporary project root',rawContainerInspectRetained:cleanupFailures.some(x=>x.startsWith('temporary directory ')),failures:cleanupFailures,retainedUnownedImageRefs:[...new Set(retained)]},
  status:Number(finalExitCode)===0?'passed':'failed'};
fs.writeFileSync(path.join(dir,'evidence.json'),JSON.stringify(evidence,null,2)+'\n',{mode:0o600});
NODE
  exit "$final_rc"
}
trap cleanup EXIT

cd "$repo_root"
record_step engine-info "$output_dir/engine-info.log" "$engine" info
record_step cli-version "$output_dir/devcontainer-cli-version.log" npx --yes --package "@devcontainers/cli@$cli_version" devcontainer --version
[[ "$(tail -n 1 "$output_dir/devcontainer-cli-version.log" | tr -d '\r\n')" == "$cli_version" ]] || {
  printf 'Expected Dev Container CLI %s.\n' "$cli_version" >&2
  exit 1
}
record_step template-docs "$output_dir/template-docs.log" bash -c 'set -e; mkdir -p "$1/src"; cp -R "$2/templates/src/minimal" "$2/templates/src/curated" "$1/src/"; npx --yes --package "@devcontainers/cli@0.89.0" devcontainer templates generate-docs --project-folder "$1/src"; test -s "$1/src/minimal/README.md"; test -s "$1/src/curated/README.md"' _ "$work_dir/template-docs" "$repo_root"

inspect_home_mount() {
  local template="$1" project="$2" expected_user="$3" id="$4" stage="$5"
  local raw="$work_dir/inspect-$template-$stage.json" summary="$output_dir/home-mount-$template-$stage.json"
  "$engine" inspect "$id" > "$raw"
  node - "$raw" "$template" "$project" "$expected_user" "$summary" <<'NODE'
const fs=require('node:fs');const [file,name,project,user,output]=process.argv.slice(2);
const c=JSON.parse(fs.readFileSync(file,'utf8'))[0];
if(c.Config.User && c.Config.User!==user)throw Error(`${name}: container user ${c.Config.User}, expected ${user}`);
const mount=(c.Mounts||[]).find(x=>x.Destination===`/home/${user}`);
if(!mount)throw Error(`${name}: missing home mount at /home/${user}`);
if(name==='bind-home'){
  if(mount.Type!=='bind'||!mount.Source.includes(`${project}/.aibox-home`))throw Error(`${name}: home bind is not project scoped`);
}else if(mount.Type!=='volume'||!mount.Name?.startsWith('aibox-home-'))throw Error(`${name}: home is not the named project volume`);
fs.writeFileSync(output,JSON.stringify({containerId:c.Id,containerUser:c.Config.User||user,image:c.Config.Image||null,home:{type:mount.Type,name:mount.Name||null,source:mount.Source,target:mount.Destination}},null,2)+'\n',{mode:0o600});
process.stdout.write(`${mount.Type}|${mount.Name||mount.Source}|${c.Config.Image||''}`);
NODE
}

find_run_container() {
  local label="$1" id
  id="$("$engine" ps -aq --filter "label=$label" | head -n 1)"
  [[ -n "$id" ]] || return 1
  printf '%s' "$id"
}

names=(minimal curated custom-user bind-home)
for template in "${names[@]}"; do
  project="$work_dir/$template"
  mkdir -p "$project"
  if [[ "$template" == minimal || "$template" == curated ]]; then
    record_step "render-$template" "$output_dir/render-$template.log" node "$repo_root/templates/render.mjs" --template "$template" --output "$project"
  else
    cp -R "$repo_root/spec/v1/examples/$template/." "$project/"
  fi
  [[ "$template" == bind-home ]] && mkdir -p "$project/.aibox-home"
  expected_user=aibox
  [[ "$template" == custom-user ]] && expected_user=dev
  record_step "read-configuration-$template" "$output_dir/read-configuration-$template.log" npx --yes --package "@devcontainers/cli@$cli_version" devcontainer read-configuration --docker-path "$engine" --workspace-folder "$project" --include-merged-configuration
  image_tag="aibox-v1-04-$run_id-$template"
  record_step "build-$template" "$output_dir/build-$template.log" npx --yes --package "@devcontainers/cli@$cli_version" devcontainer build --docker-path "$engine" --workspace-folder "$project" --image-name "$image_tag" --label "aibox.v1-04.run=$run_id"
  run_label="aibox.v1-04.run=$run_id-$template"
  record_step "up-$template" "$output_dir/up-$template.log" npx --yes --package "@devcontainers/cli@$cli_version" devcontainer up --docker-path "$engine" --workspace-folder "$project" --id-label "$run_label" --include-configuration
  if ! container_id="$(find_run_container "$run_label")"; then
    printf 'Could not identify the exact container created for %s; stopping to avoid broad cleanup.\n' "$template" >&2
    exit 1
  fi
  printf '%s\n' "$container_id" >> "$container_ids_file"
  home_mount_before="$(inspect_home_mount "$template" "$project" "$expected_user" "$container_id" first)"
  IFS='|' read -r home_type home_identity image_ref <<< "$home_mount_before"
  if [[ "$home_type" == volume ]]; then printf '%s\n' "$home_identity" >> "$home_volumes_file"; fi
  [[ -n "$image_ref" ]] && printf '%s\n' "$image_ref" >> "$observed_images_file"
  exec_check='test "$(id -un)" = '"$expected_user"'; test -w "$HOME"; command -v tmux >/dev/null; tmux -V; command -v yazi >/dev/null; yazi --version; command -v codex >/dev/null; codex --version'
  record_step "exec-$template" "$output_dir/exec-$template.log" npx --yes --package "@devcontainers/cli@$cli_version" devcontainer exec --docker-path "$engine" --workspace-folder "$project" --id-label "$run_label" bash -lc "$exec_check"
  record_step "write-home-marker-$template" "$output_dir/write-home-marker-$template.log" npx --yes --package "@devcontainers/cli@$cli_version" devcontainer exec --docker-path "$engine" --workspace-folder "$project" --id-label "$run_label" bash -lc "printf 'v1-04-persisted\\n' > \"\$HOME/.aibox-v1-04-marker\""
  record_step "remove-for-rebuild-$template" "$output_dir/remove-for-rebuild-$template.log" "$engine" rm -f "$container_id"
  reup_args=(npx --yes --package "@devcontainers/cli@$cli_version" devcontainer up --docker-path "$engine" --workspace-folder "$project" --id-label "$run_label" --include-configuration)
  if [[ "$template" == minimal || "$template" == curated ]]; then reup_args+=(--build-no-cache); fi
  record_step "rebuild-up-$template" "$output_dir/rebuild-up-$template.log" "${reup_args[@]}"
  container_id="$(find_run_container "$run_label")"
  printf '%s\n' "$container_id" >> "$container_ids_file"
  home_mount_after="$(inspect_home_mount "$template" "$project" "$expected_user" "$container_id" rebuilt)"
  IFS='|' read -r home_type_after home_identity_after image_ref_after <<< "$home_mount_after"
  [[ "$home_type" == "$home_type_after" && "$home_identity" == "$home_identity_after" ]] || {
    printf '%s: home mount identity changed after recreate.\n' "$template" >&2
    exit 1
  }
  [[ "$home_type_after" != volume ]] || printf '%s\n' "$home_identity_after" >> "$home_volumes_file"
  [[ -z "$image_ref_after" ]] || printf '%s\n' "$image_ref_after" >> "$observed_images_file"
  record_step "verify-home-marker-$template" "$output_dir/verify-home-marker-$template.log" npx --yes --package "@devcontainers/cli@$cli_version" devcontainer exec --docker-path "$engine" --workspace-folder "$project" --id-label "$run_label" bash -lc 'test "$(cat "$HOME/.aibox-v1-04-marker")" = v1-04-persisted && test -w "$HOME"'
done

printf 'V1-04 host lifecycle verification passed. Evidence: %s\n' "$output_dir"
