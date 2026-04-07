//! The mapping from processkit cache files to their installation paths
//! in a consuming project. **Single source of truth** used by:
//!
//! - `aibox init` (A5) — installs files via this mapping
//! - `aibox sync` (A6) — uses this mapping for the live-vs-cache 3-way comparison
//! - `aibox migrate` (A7) — references files via this mapping when working
//!   through Migration documents
//!
//! ## Install layout
//!
//! Cache file path (relative to `<cache>/<src_path>/`) → project install path
//! (relative to project root).
//!
//! | Cache path                          | Project install path                       |
//! |-------------------------------------|--------------------------------------------|
//! | `skills/<name>/...`                 | `.claude/skills/<name>/...`                |
//! | `lib/processkit/...`                | `.claude/skills/_lib/processkit/...`       |
//! | `primitives/schemas/<f>.yaml`       | `context/.aibox/schemas/<f>.yaml`          |
//! | `primitives/state-machines/<f>.yaml`| `context/.aibox/state-machines/<f>.yaml`   |
//! | `primitives/FORMAT.md`              | `context/.aibox/primitives-FORMAT.md`      |
//! | `processes/<f>.md`                  | `context/.aibox/processes/<f>.md`          |
//! | `PROVENANCE.toml` (top of src/)     | `context/.aibox/processkit-provenance.toml`|
//! | any `INDEX.md` (every level)        | skipped (processkit-internal)              |
//! | any `FORMAT.md` other than primitives/FORMAT.md | skipped (processkit-internal)  |
//! | `packages/...`                      | skipped (consumed by init logic, not installed) |
//! | anything unrecognized               | skipped, with a warning logged at install time |
//!
//! ## Why this layout
//!
//! Primitives, processes, and provenance are reference material the agent
//! reads but rarely writes. They live in `context/.aibox/` (a gitignored
//! cache mirror, set up by `aibox init`) so they can be refreshed without
//! polluting the user-facing parts of `context/`.
//!
//! Skills and the shared lib live under `.claude/skills/` because that's
//! where MCP servers' `_find_lib()` boilerplate looks for the lib at runtime
//! (it walks up from `<server>/mcp/server.py` to find `_lib/processkit/`).
//! See `processkit/src/lib/README.md` for the lookup logic.

use std::path::{Path, PathBuf};

/// What to do with a single file from the processkit cache.
#[derive(Debug, Clone, PartialEq, Eq)]
pub enum InstallAction {
    /// Install at this project-root-relative path.
    Install(PathBuf),
    /// Skip this file (processkit-internal, not user-facing).
    Skip,
}

/// Map a cache file (path relative to `<cache>/<src_path>/`) to its install
/// action. Pure function — no I/O.
///
/// The input is a relative path using forward slashes. Backslashes are
/// normalized for Windows-friendliness (though aibox does not target Windows
/// hosts at v0.14.x; the normalization is defensive).
pub fn install_action_for(rel_path: &Path) -> InstallAction {
    let s = rel_path.to_string_lossy().replace('\\', "/");
    if s.is_empty() || s == "." {
        return InstallAction::Skip;
    }

    let parts: Vec<&str> = s.split('/').filter(|p| !p.is_empty()).collect();
    if parts.is_empty() {
        return InstallAction::Skip;
    }

    // INDEX.md files are processkit-internal docs at every level.
    if parts.last().copied() == Some("INDEX.md") {
        return InstallAction::Skip;
    }

    // Top-level files
    if parts.len() == 1 {
        return match parts[0] {
            "PROVENANCE.toml" => {
                InstallAction::Install(PathBuf::from("context/.aibox/processkit-provenance.toml"))
            }
            _ => InstallAction::Skip,
        };
    }

    // primitives/FORMAT.md is installed as a reference. Other FORMAT.md
    // files (e.g. skills/FORMAT.md) are processkit-internal.
    if s == "primitives/FORMAT.md" {
        return InstallAction::Install(PathBuf::from("context/.aibox/primitives-FORMAT.md"));
    }
    if parts.last().copied() == Some("FORMAT.md") {
        return InstallAction::Skip;
    }

    match parts[0] {
        // skills/<name>/...  →  .claude/skills/<name>/...
        "skills" if parts.len() >= 2 => {
            let mut p = PathBuf::from(".claude/skills");
            for part in &parts[1..] {
                p.push(part);
            }
            InstallAction::Install(p)
        }

        // lib/processkit/...  →  .claude/skills/_lib/processkit/...
        // (matches the _find_lib() boilerplate in MCP server scripts)
        "lib" if parts.len() >= 2 => {
            let mut p = PathBuf::from(".claude/skills/_lib");
            for part in &parts[1..] {
                p.push(part);
            }
            InstallAction::Install(p)
        }

        // primitives/schemas/X.yaml  →  context/.aibox/schemas/X.yaml
        // primitives/state-machines/X.yaml  →  context/.aibox/state-machines/X.yaml
        // primitives/<other>/...  →  context/.aibox/<rest>
        "primitives" if parts.len() >= 2 => {
            let mut p = PathBuf::from("context/.aibox");
            for part in &parts[1..] {
                p.push(part);
            }
            InstallAction::Install(p)
        }

        // processes/<f>.md  →  context/.aibox/processes/<f>.md
        // The actual user-facing process definitions live under
        // context/processes/ — installed processes are reference copies the
        // user can override by writing their own version.
        "processes" if parts.len() >= 2 => {
            let mut p = PathBuf::from("context/.aibox/processes");
            for part in &parts[1..] {
                p.push(part);
            }
            InstallAction::Install(p)
        }

        // packages/* — consumed by init logic to know which skills to install,
        // not installed into the project itself.
        "packages" => InstallAction::Skip,

        // Anything else is unknown — skip rather than guess.
        _ => InstallAction::Skip,
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    fn install(s: &str) -> InstallAction {
        install_action_for(Path::new(s))
    }

    fn assert_installs_to(input: &str, expected: &str) {
        match install(input) {
            InstallAction::Install(p) => assert_eq!(p, PathBuf::from(expected)),
            other => panic!("expected Install({}), got {:?}", expected, other),
        }
    }

    fn assert_skipped(input: &str) {
        assert_eq!(install(input), InstallAction::Skip);
    }

    #[test]
    fn skill_files_install_under_claude() {
        assert_installs_to("skills/event-log/SKILL.md", ".claude/skills/event-log/SKILL.md");
        assert_installs_to(
            "skills/workitem-management/templates/workitem.yaml",
            ".claude/skills/workitem-management/templates/workitem.yaml",
        );
        assert_installs_to(
            "skills/event-log/mcp/server.py",
            ".claude/skills/event-log/mcp/server.py",
        );
    }

    #[test]
    fn lib_installs_under_claude_lib() {
        assert_installs_to(
            "lib/processkit/__init__.py",
            ".claude/skills/_lib/processkit/__init__.py",
        );
        assert_installs_to(
            "lib/processkit/entity.py",
            ".claude/skills/_lib/processkit/entity.py",
        );
    }

    #[test]
    fn primitive_schemas_install_under_aibox_dir() {
        assert_installs_to(
            "primitives/schemas/workitem.yaml",
            "context/.aibox/schemas/workitem.yaml",
        );
        assert_installs_to(
            "primitives/state-machines/workitem.yaml",
            "context/.aibox/state-machines/workitem.yaml",
        );
    }

    #[test]
    fn primitive_format_doc_installs_with_renamed_path() {
        assert_installs_to(
            "primitives/FORMAT.md",
            "context/.aibox/primitives-FORMAT.md",
        );
    }

    #[test]
    fn provenance_installs_with_renamed_path() {
        assert_installs_to(
            "PROVENANCE.toml",
            "context/.aibox/processkit-provenance.toml",
        );
    }

    #[test]
    fn processes_install_under_aibox_processes() {
        assert_installs_to(
            "processes/release.md",
            "context/.aibox/processes/release.md",
        );
        assert_installs_to(
            "processes/code-review/template.yaml",
            "context/.aibox/processes/code-review/template.yaml",
        );
    }

    #[test]
    fn index_md_at_every_level_is_skipped() {
        assert_skipped("INDEX.md");
        assert_skipped("primitives/INDEX.md");
        assert_skipped("primitives/schemas/INDEX.md");
        assert_skipped("skills/INDEX.md");
        assert_skipped("packages/INDEX.md");
        assert_skipped("processes/INDEX.md");
    }

    #[test]
    fn skills_format_doc_is_skipped() {
        // skills/FORMAT.md is processkit-internal, only primitives/FORMAT.md is installed.
        assert_skipped("skills/FORMAT.md");
    }

    #[test]
    fn packages_are_skipped() {
        assert_skipped("packages/minimal.yaml");
        assert_skipped("packages/managed.yaml");
        assert_skipped("packages/software.yaml");
    }

    #[test]
    fn unknown_top_level_is_skipped() {
        assert_skipped("CHANGELOG.md");
        assert_skipped("LICENSE");
        assert_skipped("random/file.md");
    }

    #[test]
    fn empty_path_is_skipped() {
        assert_skipped("");
        assert_skipped(".");
    }

    #[test]
    fn windows_backslashes_are_normalized() {
        assert_installs_to(
            "skills\\event-log\\SKILL.md",
            ".claude/skills/event-log/SKILL.md",
        );
    }
}
