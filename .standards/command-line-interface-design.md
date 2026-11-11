# Command-line interface design standard

## Purpose and scope

This standard defines the command grammar of projectious.work CLIs: how
commands, arguments, and options are named and shaped. It offers three styles
as options and a baseline every style follows. Streams, machine output, exit
codes, and configuration precedence are defined in the
[output](application-output-logging-and-evidence.md#streams-and-formats),
[compatibility](compatibility-and-machine-interfaces.md#cli-profile), and
[configuration](application-configuration.md#precedence) standards.

Terms in *italics* are [Semantic Anchors](https://github.com/llm-coding/Semantic-Anchors/tree/6554641a11fdf99fdbdf09cc0602268e3db88966); see
[README](README.md#conventions).

## Choosing a style

Each CLI selects exactly one style, records it with the reason in an *ADR
according to Nygard*, and applies it to every command. A new command follows
the recorded style even where another style would read better locally.

| Style | Grammar | Choose when |
|---|---|---|
| [A. Unix utility](#a-unix-utility-style) | `tool [options] [operands]` | The tool does one job on files or streams and is used in pipelines. |
| [B. Verb–resource](#b-verbresource-style-kubectl) | `tool <verb> <type> [<name>] [options]` | The tool manages many resource types of one API with a shared verb set and declarative state. |
| [C. Object–action](#c-objectaction-style-openstackclient) | `tool <object> <action> [<object>] [options]` | The tool fronts several services with many object kinds and a small fixed action set. |

Docker-style management commands (`tool container ls`) are a variant of C,
and Git-style subcommands (`tool commit`) are a variant of A with
subcommands; a CLI that uses either records it as that style.

## Baseline for every style

Every style follows the *Unix Philosophy*: each command does one job, results
go to stdout and compose through pipes, and success is quiet unless output is
requested.

1. `-h` and `--help` show help and ignore other arguments; help opens with
   examples. A command that requires arguments and gets none shows short help
   and fails with a usage exit code.
2. `--version` prints the version. `-v` never means version.
3. Every option has a long name; one-letter forms are reserved for frequent
   options. Reused names keep their common meaning: `--all`, `--force`,
   `--dry-run`, `--quiet`, `--verbose`, `--no-input`, `--no-color`.
4. Options are preferred to positional arguments. Positional arguments carry
   only the primary inputs or names of the command.
5. `--` ends option parsing; `-` as an operand means stdin or stdout.
6. Human output is the default; its wording and layout are not a contract,
   but its stream assignment and exit codes are. Scripts use the machine
   format, selected by the style's format option. It emits the versioned
   envelope of the
   [compatibility standard](compatibility-and-machine-interfaces.md#machine-readable-results);
   a raw resource dump is a payload inside that envelope.
7. Secrets are never accepted as option values; their sources follow the
   [security standard](security-and-software-supply-chain.md#secrets-and-sensitive-data).
8. The CLI prompts only when stdin is a terminal. `--no-input` disables every
   prompt. Destructive commands offer `--dry-run`, confirm in proportion to
   the damage (typing the resource name for irreversible actions), and accept
   `--force` or an explicit confirmation option for scripts.
9. Errors state what happened and how to fix it; the most important line is
   last.
10. Verbs, option names, and formats mean the same thing in every command.
    Near-synonyms (`update` and `upgrade`) are not both used.

## A. Unix utility style

Source: POSIX Utility Syntax Guidelines and the GNU Coding Standards.

- Short options are one character after `-`; options without an argument MAY
  be grouped (`-abc`).
- An option argument is never optional. `-o value`, `-ovalue`, `--opt value`,
  and `--opt=value` are all accepted.
- Inputs are positional operands; outputs are named with `-o`/`--output`, not
  positionally.
- Options precede operands; GNU-style interleaving MAY be accepted.
  Order between options does not matter, except for repeated and mutually
  exclusive options.
- `-n` MAY mean `--dry-run`, as in common utilities.

## B. Verb–resource style (kubectl)

Source: Kubernetes kubectl reference and usage conventions.

- Grammar: `tool <verb> <type> [<name>...]`, or `<type>/<name>` and
  `-f <file|dir|url>` for several targets. One verb set serves all types.
- Type names are case-insensitive and accept singular, plural, and one
  documented short name; resource names are case-sensitive. Omitting the name
  addresses all resources of that type.
- Uniform scope and filter options: `-n`/`--namespace` for scope, `-A` for all
  scopes, `-l`/`--selector` for label selectors, `-f` for input files.
  `-n` therefore never means dry run.
- `-o`/`--output` selects `json`, `yaml`, `name`, or `wide`; scripts never
  parse the default table.
- Mutating verbs offer `--dry-run=none|client|server`: `client` renders
  locally, `server` validates without persisting.
- Every resource created imperatively can also be applied declaratively
  (`apply -f`); production changes use the declarative path.
- Plugins are executables named `tool-<name>` on `PATH` and never shadow
  built-in commands.

## C. Object–action style (OpenStackClient)

Source: OpenStackClient Human Interface Guide.

- Grammar: `tool <object> <action> [<object>]`. Objects are singular, may be
  several words (`volume type`), and are never abbreviated. With two objects,
  the primary object precedes the action (`group add user`).
- Actions come from a fixed set; a new action needs a recorded reason:
  - `create`: at most a name; everything else is an option;
  - `show`: one name or ID;
  - `list`: no positional arguments; filters are options;
  - `delete`: one or more names or IDs;
  - `set` and `unset`: one name or ID, changing only that object's own
    attributes;
  - `add` and `remove`: the parent and the child.
- Command options are long-only, lowercase, and hyphen-separated. Repeated
  values repeat the option (`--property a=1 --property b=2`); booleans are
  pairs (`--enable`/`--disable` or `--x`/`--no-x`); timestamps use ISO 8601.
- Global options map to environment variables under the
  [configuration standard's](application-configuration.md#precedence)
  precedence and authority rules; secret-bearing or authority-restricted
  options have no environment equivalent.
- `list` and `show` select output with `-f`/`--format`; the default is a
  table.

## Agent briefing

> Each CLI picks one recorded style (ADR): A Unix utility `tool [opts]
> [operands]`, B verb–resource `tool <verb> <type> [<name>]` (kubectl), or C
> object–action `tool <object> <action>` (OpenStackClient). All follow the
> *Unix Philosophy* and the baseline: `--help` with examples, `--version`,
> long names for all options, `--` and `-`, machine envelope as contract, no
> secrets in options, prompts only on a TTY, `--dry-run` for destructive use.
