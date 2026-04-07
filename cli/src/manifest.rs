//! Read/write the git-tracked files that record processkit installation
//! state in a project. Two files:
//!
//! - `context/.aibox/processkit.lock`     — what version is pinned
//! - `context/.aibox/processkit.manifest` — per-file SHA256 reference
//!
//! Both are TOML, both are committed to git. The cache (under
//! `~/.cache/aibox/processkit/`) is reproducible from the lock; the
//! manifest is the "as-installed" reference used by the three-way
//! comparison in A6 (manifest vs cache vs live working tree).
//!
//! ## Scope
//!
//! This module is intentionally policy-free. It does not know how
//! files in `<cache>/<src_path>/` map to their final locations in the
//! project (e.g. `.claude/skills/...` vs `context/processes/...`) —
//! that is A5's job. The manifest produced by [`manifest_from_cache`]
//! uses keys that are **relative to the cache `<src_path>/`**, not
//! relative to the project root.

use anyhow::{Context, Result};
use serde::{Deserialize, Serialize};
use sha2::{Digest, Sha256};
use std::collections::BTreeMap;
use std::fs;
use std::path::{Component, Path, PathBuf};

// ---------------------------------------------------------------------------
// Lock file
// ---------------------------------------------------------------------------

/// `context/.aibox/processkit.lock` contents.
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub struct ProcessKitLock {
    pub source: String,
    pub version: String,
    pub src_path: String,
    pub branch: Option<String>,
    pub resolved_commit: Option<String>,
    /// ISO 8601 UTC timestamp of the install (e.g. "2026-04-06T12:34:56Z").
    pub installed_at: String,
}

// ---------------------------------------------------------------------------
// Manifest file
// ---------------------------------------------------------------------------

/// One entry in `processkit.manifest` — describes a single shipped file
/// at install time.
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
pub struct ManifestEntry {
    /// SHA256 hex digest of the file in the cache when installed.
    pub upstream_sha: String,
    /// Logical group, e.g. `"skills/event-log"`. Used by the A6
    /// three-way comparison's "auto-update by group, not by file" rule.
    pub group: Option<String>,
}

/// `context/.aibox/processkit.manifest` contents.
///
/// File-path keys are stored as forward-slash-delimited strings so the
/// file is portable across platforms and stable across reads/writes
/// (`BTreeMap` provides deterministic ordering).
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, Default)]
pub struct ProcessKitManifest {
    pub files: BTreeMap<String, ManifestEntry>,
}

// ---------------------------------------------------------------------------
// Standard file locations
// ---------------------------------------------------------------------------

/// `<project_root>/context/.aibox/processkit.lock`
pub fn lock_path(project_root: &Path) -> PathBuf {
    project_root
        .join("context")
        .join(".aibox")
        .join("processkit.lock")
}

/// `<project_root>/context/.aibox/processkit.manifest`
pub fn manifest_path(project_root: &Path) -> PathBuf {
    project_root
        .join("context")
        .join(".aibox")
        .join("processkit.manifest")
}

// ---------------------------------------------------------------------------
// Lock read / write
// ---------------------------------------------------------------------------

/// Read the lock file. Returns `Ok(None)` if the file does not exist.
pub fn read_lock(project_root: &Path) -> Result<Option<ProcessKitLock>> {
    let path = lock_path(project_root);
    if !path.exists() {
        return Ok(None);
    }
    let body = fs::read_to_string(&path)
        .with_context(|| format!("failed to read {}", path.display()))?;
    let parsed: ProcessKitLock = toml::from_str(&body)
        .with_context(|| format!("failed to parse {} as TOML", path.display()))?;
    Ok(Some(parsed))
}

/// Write the lock file, creating parent directories as needed.
pub fn write_lock(project_root: &Path, lock: &ProcessKitLock) -> Result<()> {
    let path = lock_path(project_root);
    if let Some(parent) = path.parent() {
        fs::create_dir_all(parent)
            .with_context(|| format!("failed to create {}", parent.display()))?;
    }
    let body = toml::to_string_pretty(lock)
        .with_context(|| "failed to serialize ProcessKitLock to TOML".to_string())?;
    fs::write(&path, body).with_context(|| format!("failed to write {}", path.display()))?;
    Ok(())
}

// ---------------------------------------------------------------------------
// Manifest read / write
// ---------------------------------------------------------------------------

/// Read the manifest file. Returns `Ok(None)` if the file does not exist.
pub fn read_manifest(project_root: &Path) -> Result<Option<ProcessKitManifest>> {
    let path = manifest_path(project_root);
    if !path.exists() {
        return Ok(None);
    }
    let body = fs::read_to_string(&path)
        .with_context(|| format!("failed to read {}", path.display()))?;
    let parsed: ProcessKitManifest = toml::from_str(&body)
        .with_context(|| format!("failed to parse {} as TOML", path.display()))?;
    Ok(Some(parsed))
}

/// Write the manifest file, creating parent directories as needed.
pub fn write_manifest(project_root: &Path, manifest: &ProcessKitManifest) -> Result<()> {
    let path = manifest_path(project_root);
    if let Some(parent) = path.parent() {
        fs::create_dir_all(parent)
            .with_context(|| format!("failed to create {}", parent.display()))?;
    }
    let body = toml::to_string_pretty(manifest)
        .with_context(|| "failed to serialize ProcessKitManifest to TOML".to_string())?;
    fs::write(&path, body).with_context(|| format!("failed to write {}", path.display()))?;
    Ok(())
}

// ---------------------------------------------------------------------------
// Hashing
// ---------------------------------------------------------------------------

/// Compute the SHA256 hex digest of a file's contents. Reads the file
/// in 64 KiB chunks so very large files do not balloon memory.
pub fn sha256_of_file(path: &Path) -> Result<String> {
    use std::io::Read;
    let mut f = fs::File::open(path)
        .with_context(|| format!("failed to open {} for hashing", path.display()))?;
    let mut hasher = Sha256::new();
    let mut buf = [0u8; 64 * 1024];
    loop {
        let n = f
            .read(&mut buf)
            .with_context(|| format!("failed to read {} for hashing", path.display()))?;
        if n == 0 {
            break;
        }
        hasher.update(&buf[..n]);
    }
    Ok(format!("{:x}", hasher.finalize()))
}

// ---------------------------------------------------------------------------
// Cache walking
// ---------------------------------------------------------------------------

/// Walk a fetched cache `<src_path>` directory and produce a fresh
/// manifest mapping every shipped file to its SHA256 + group.
///
/// Keys in the resulting manifest are paths relative to `src_path`,
/// using forward slashes regardless of host OS.
///
/// Skips:
/// - `.git/` directories (in case the cache came from a `git clone`)
/// - `__pycache__/` directories and `*.pyc` files
/// - the A3 idempotency marker (`.fetch-complete`)
/// - any other dotfile or dot-directory at the root of the walk
pub fn manifest_from_cache(src_path: &Path) -> Result<ProcessKitManifest> {
    if !src_path.is_dir() {
        return Err(anyhow::anyhow!(
            "manifest_from_cache: {} is not a directory",
            src_path.display()
        ));
    }
    let mut files: BTreeMap<String, ManifestEntry> = BTreeMap::new();
    walk_into(src_path, src_path, &mut files)?;
    Ok(ProcessKitManifest { files })
}

/// Recursive helper for [`manifest_from_cache`]. `root` is the original
/// `src_path` (kept stable across recursion so relative paths are correct);
/// `dir` is the directory currently being walked.
fn walk_into(
    root: &Path,
    dir: &Path,
    files: &mut BTreeMap<String, ManifestEntry>,
) -> Result<()> {
    for entry in fs::read_dir(dir)
        .with_context(|| format!("failed to read directory {}", dir.display()))?
    {
        let entry = entry?;
        let path = entry.path();
        let file_type = entry
            .file_type()
            .with_context(|| format!("failed to stat {}", path.display()))?;
        let name = entry.file_name();
        let name_str = name.to_string_lossy().to_string();

        if should_skip(&name_str) {
            continue;
        }

        if file_type.is_dir() {
            walk_into(root, &path, files)?;
            continue;
        }
        if !file_type.is_file() {
            // Symlinks, fifos, etc. — ignore.
            continue;
        }

        let rel = path
            .strip_prefix(root)
            .with_context(|| {
                format!(
                    "failed to relativize {} against {}",
                    path.display(),
                    root.display()
                )
            })?
            .to_path_buf();

        let key = path_to_forward_slash(&rel);
        let sha = sha256_of_file(&path)?;
        let group = group_for_path(&rel);
        files.insert(
            key,
            ManifestEntry {
                upstream_sha: sha,
                group,
            },
        );
    }
    Ok(())
}

/// Returns true if a directory entry name should be skipped during the
/// cache walk.
fn should_skip(name: &str) -> bool {
    if name == ".git" || name == "__pycache__" || name == ".fetch-complete" {
        return true;
    }
    if name.starts_with('.') {
        return true;
    }
    if name.ends_with(".pyc") {
        return true;
    }
    false
}

/// Convert a relative `Path` into its forward-slash string form. Used so
/// manifest keys are stable across Windows and Unix hosts.
fn path_to_forward_slash(rel: &Path) -> String {
    let mut parts: Vec<String> = Vec::new();
    for c in rel.components() {
        if let Component::Normal(os) = c {
            parts.push(os.to_string_lossy().to_string());
        }
    }
    parts.join("/")
}

// ---------------------------------------------------------------------------
// Group heuristic
// ---------------------------------------------------------------------------

/// Compute a logical group name for a `src_path`-relative file path.
///
/// The grouping drives the A6 "auto-update by group, never by individual
/// file" rule, so the goal is "files that should always move together
/// share a group". Heuristic, in priority order:
///
/// 1. `PROVENANCE.toml` at the top → `"PROVENANCE"`
/// 2. `skills/<name>/...`           → `"skills/<name>"` (everything below
///    a skill directory belongs to that skill)
/// 3. `primitives/schemas/<X>.yaml` → `"primitives/schemas/<X>"`
/// 4. `primitives/state-machines/<X>.yaml` → `"primitives/state-machines/<X>"`
/// 5. `primitives/<other>`          → `"primitives"`
/// 6. `lib/...`                     → `"lib"`
/// 7. `processes/<name>/...` (or `processes/<name>.md`) → `"processes/<name>"`
/// 8. anything else                 → the immediate parent directory or `None`
pub fn group_for_path(rel_path: &Path) -> Option<String> {
    let parts: Vec<String> = rel_path
        .components()
        .filter_map(|c| match c {
            Component::Normal(os) => Some(os.to_string_lossy().to_string()),
            _ => None,
        })
        .collect();

    if parts.is_empty() {
        return None;
    }

    // 1. PROVENANCE.toml at top
    if parts.len() == 1 && parts[0] == "PROVENANCE.toml" {
        return Some("PROVENANCE".to_string());
    }

    // 2. skills/<name>/...
    if parts[0] == "skills" && parts.len() >= 2 {
        return Some(format!("skills/{}", parts[1]));
    }

    // 3, 4, 5. primitives/...
    if parts[0] == "primitives" {
        if parts.len() >= 3
            && (parts[1] == "schemas" || parts[1] == "state-machines")
        {
            // Strip a trailing extension from the leaf so all files
            // describing the same primitive share a group.
            let leaf = strip_known_ext(&parts[2]);
            return Some(format!("primitives/{}/{}", parts[1], leaf));
        }
        return Some("primitives".to_string());
    }

    // 6. lib/...
    if parts[0] == "lib" {
        return Some("lib".to_string());
    }

    // 7. processes/<name>(/...)
    if parts[0] == "processes" && parts.len() >= 2 {
        let leaf = strip_known_ext(&parts[1]);
        return Some(format!("processes/{}", leaf));
    }

    // 8. fallback — immediate parent dir, or None for top-level loose files.
    if parts.len() >= 2 {
        return Some(parts[..parts.len() - 1].join("/"));
    }
    None
}

/// Strip a single trailing known extension from a file name.
fn strip_known_ext(name: &str) -> String {
    for ext in [".yaml", ".yml", ".toml", ".md", ".py"] {
        if let Some(stripped) = name.strip_suffix(ext) {
            return stripped.to_string();
        }
    }
    name.to_string()
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

#[cfg(test)]
mod tests {
    use super::*;
    use std::fs;
    use tempfile::TempDir;

    fn sample_lock() -> ProcessKitLock {
        ProcessKitLock {
            source: "https://github.com/projectious-work/processkit.git".to_string(),
            version: "v0.4.0".to_string(),
            src_path: "src".to_string(),
            branch: None,
            resolved_commit: Some("deadbeefcafebabe".to_string()),
            installed_at: "2026-04-06T12:00:00Z".to_string(),
        }
    }

    fn sample_manifest() -> ProcessKitManifest {
        let mut files = BTreeMap::new();
        for (k, sha, group) in [
            (
                "PROVENANCE.toml",
                "1111111111111111111111111111111111111111111111111111111111111111",
                "PROVENANCE",
            ),
            (
                "skills/event-log/SKILL.md",
                "2222222222222222222222222222222222222222222222222222222222222222",
                "skills/event-log",
            ),
            (
                "skills/event-log/templates/entry.yaml",
                "3333333333333333333333333333333333333333333333333333333333333333",
                "skills/event-log",
            ),
            (
                "primitives/schemas/workitem.yaml",
                "4444444444444444444444444444444444444444444444444444444444444444",
                "primitives/schemas/workitem",
            ),
            (
                "lib/processkit/entity.py",
                "5555555555555555555555555555555555555555555555555555555555555555",
                "lib",
            ),
        ] {
            files.insert(
                k.to_string(),
                ManifestEntry {
                    upstream_sha: sha.to_string(),
                    group: Some(group.to_string()),
                },
            );
        }
        ProcessKitManifest { files }
    }

    // -- Round trips --------------------------------------------------------

    #[test]
    fn lock_round_trip() {
        let tmp = TempDir::new().unwrap();
        let lock = sample_lock();
        write_lock(tmp.path(), &lock).unwrap();
        let back = read_lock(tmp.path()).unwrap().unwrap();
        assert_eq!(back, lock);
    }

    #[test]
    fn manifest_round_trip() {
        let tmp = TempDir::new().unwrap();
        let mf = sample_manifest();
        write_manifest(tmp.path(), &mf).unwrap();
        let back = read_manifest(tmp.path()).unwrap().unwrap();
        assert_eq!(back, mf);
    }

    #[test]
    fn lock_returns_none_if_missing() {
        let tmp = TempDir::new().unwrap();
        assert!(read_lock(tmp.path()).unwrap().is_none());
    }

    #[test]
    fn manifest_returns_none_if_missing() {
        let tmp = TempDir::new().unwrap();
        assert!(read_manifest(tmp.path()).unwrap().is_none());
    }

    #[test]
    fn write_lock_creates_parent_directories() {
        let tmp = TempDir::new().unwrap();
        write_lock(tmp.path(), &sample_lock()).unwrap();
        assert!(tmp.path().join("context/.aibox/processkit.lock").exists());
    }

    #[test]
    fn write_manifest_creates_parent_directories() {
        let tmp = TempDir::new().unwrap();
        write_manifest(tmp.path(), &sample_manifest()).unwrap();
        assert!(
            tmp.path()
                .join("context/.aibox/processkit.manifest")
                .exists()
        );
    }

    // -- Hashing ------------------------------------------------------------

    #[test]
    fn sha256_of_file_known_value() {
        let tmp = TempDir::new().unwrap();
        let p = tmp.path().join("hello.txt");
        fs::write(&p, b"hello\n").unwrap();
        let sha = sha256_of_file(&p).unwrap();
        // sha256("hello\n")
        assert_eq!(
            sha,
            "5891b5b522d5df086d0ff0b110fbd9d21bb4fc7163af34d08286a2e846f6be03"
        );
    }

    #[test]
    fn sha256_of_file_handles_empty_file() {
        let tmp = TempDir::new().unwrap();
        let p = tmp.path().join("empty");
        fs::write(&p, b"").unwrap();
        let sha = sha256_of_file(&p).unwrap();
        // sha256("") — known value
        assert_eq!(
            sha,
            "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
        );
    }

    // -- Cache walk ---------------------------------------------------------

    /// Build a synthetic processkit-shaped src tree under `root`.
    fn synth_src_tree(root: &Path) {
        fs::create_dir_all(root.join("skills/event-log/templates")).unwrap();
        fs::create_dir_all(root.join("skills/event-log/__pycache__")).unwrap();
        fs::create_dir_all(root.join("primitives/schemas")).unwrap();
        fs::create_dir_all(root.join("primitives/state-machines")).unwrap();
        fs::create_dir_all(root.join("lib/processkit")).unwrap();
        fs::create_dir_all(root.join(".git/objects")).unwrap();

        fs::write(root.join("PROVENANCE.toml"), "version = \"v0.4.0\"\n").unwrap();
        fs::write(root.join("skills/event-log/SKILL.md"), "# event log\n").unwrap();
        fs::write(
            root.join("skills/event-log/templates/entry.yaml"),
            "name: x\n",
        )
        .unwrap();
        fs::write(root.join("skills/event-log/__pycache__/x.pyc"), b"junk").unwrap();
        fs::write(
            root.join("primitives/schemas/workitem.yaml"),
            "name: workitem\n",
        )
        .unwrap();
        fs::write(
            root.join("primitives/state-machines/workflow.yaml"),
            "name: workflow\n",
        )
        .unwrap();
        fs::write(root.join("lib/processkit/entity.py"), "print('x')\n").unwrap();
        fs::write(root.join(".git/objects/foo"), b"git internals").unwrap();
        fs::write(root.join(".fetch-complete"), b"deadbeef").unwrap();
        fs::write(root.join(".DS_Store"), b"junk").unwrap();
    }

    #[test]
    fn manifest_from_cache_walks_synthetic_tree() {
        let tmp = TempDir::new().unwrap();
        let src = tmp.path().join("src");
        synth_src_tree(&src);

        let mf = manifest_from_cache(&src).unwrap();

        // Expected entries
        let expected_keys = [
            "PROVENANCE.toml",
            "skills/event-log/SKILL.md",
            "skills/event-log/templates/entry.yaml",
            "primitives/schemas/workitem.yaml",
            "primitives/state-machines/workflow.yaml",
            "lib/processkit/entity.py",
        ];
        for k in expected_keys {
            assert!(
                mf.files.contains_key(k),
                "expected manifest entry {}, got keys {:?}",
                k,
                mf.files.keys().collect::<Vec<_>>()
            );
        }

        // Spot-check groups
        assert_eq!(
            mf.files["PROVENANCE.toml"].group.as_deref(),
            Some("PROVENANCE")
        );
        assert_eq!(
            mf.files["skills/event-log/SKILL.md"].group.as_deref(),
            Some("skills/event-log")
        );
        assert_eq!(
            mf.files["skills/event-log/templates/entry.yaml"]
                .group
                .as_deref(),
            Some("skills/event-log")
        );
        assert_eq!(
            mf.files["primitives/schemas/workitem.yaml"]
                .group
                .as_deref(),
            Some("primitives/schemas/workitem")
        );
        assert_eq!(
            mf.files["primitives/state-machines/workflow.yaml"]
                .group
                .as_deref(),
            Some("primitives/state-machines/workflow")
        );
        assert_eq!(
            mf.files["lib/processkit/entity.py"].group.as_deref(),
            Some("lib")
        );
    }

    #[test]
    fn manifest_from_cache_skips_dotfiles_and_pycache() {
        let tmp = TempDir::new().unwrap();
        let src = tmp.path().join("src");
        synth_src_tree(&src);
        let mf = manifest_from_cache(&src).unwrap();

        // .git/, .DS_Store, __pycache__/x.pyc must not appear
        for k in mf.files.keys() {
            assert!(!k.starts_with(".git"), "leaked .git entry: {}", k);
            assert!(!k.contains(".DS_Store"), "leaked dotfile: {}", k);
            assert!(!k.contains("__pycache__"), "leaked pycache: {}", k);
            assert!(!k.ends_with(".pyc"), "leaked pyc: {}", k);
        }
    }

    #[test]
    fn manifest_from_cache_skips_fetch_complete_marker() {
        let tmp = TempDir::new().unwrap();
        let src = tmp.path().join("src");
        synth_src_tree(&src);
        let mf = manifest_from_cache(&src).unwrap();
        assert!(!mf.files.contains_key(".fetch-complete"));
    }

    // -- group_for_path -----------------------------------------------------

    #[test]
    fn group_for_path_skill_subfile() {
        assert_eq!(
            group_for_path(Path::new("skills/event-log/SKILL.md")),
            Some("skills/event-log".to_string())
        );
    }

    #[test]
    fn group_for_path_skill_template() {
        assert_eq!(
            group_for_path(Path::new("skills/event-log/templates/foo.yaml")),
            Some("skills/event-log".to_string())
        );
    }

    #[test]
    fn group_for_path_primitive_schema() {
        assert_eq!(
            group_for_path(Path::new("primitives/schemas/workitem.yaml")),
            Some("primitives/schemas/workitem".to_string())
        );
    }

    #[test]
    fn group_for_path_primitive_state_machine() {
        assert_eq!(
            group_for_path(Path::new("primitives/state-machines/workflow.yaml")),
            Some("primitives/state-machines/workflow".to_string())
        );
    }

    #[test]
    fn group_for_path_primitive_format_doc() {
        // Single-file primitive (e.g. primitives/FORMAT.md) collapses
        // to the catch-all "primitives" group.
        assert_eq!(
            group_for_path(Path::new("primitives/FORMAT.md")),
            Some("primitives".to_string())
        );
    }

    #[test]
    fn group_for_path_lib() {
        assert_eq!(
            group_for_path(Path::new("lib/processkit/entity.py")),
            Some("lib".to_string())
        );
    }

    #[test]
    fn group_for_path_provenance() {
        assert_eq!(
            group_for_path(Path::new("PROVENANCE.toml")),
            Some("PROVENANCE".to_string())
        );
    }

    #[test]
    fn group_for_path_process_subfile() {
        assert_eq!(
            group_for_path(Path::new("processes/foo.md")),
            Some("processes/foo".to_string())
        );
        assert_eq!(
            group_for_path(Path::new("processes/foo/template.yaml")),
            Some("processes/foo".to_string())
        );
    }

    #[test]
    fn group_for_path_unknown_top_level_file_returns_none() {
        // A loose top-level file with no recognized prefix has no group.
        assert_eq!(group_for_path(Path::new("README.md")), None);
    }
}
