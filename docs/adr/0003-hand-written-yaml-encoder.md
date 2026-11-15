# 0003. Hand-written YAML encoder for results

- **Status:** Accepted
- **Date:** 2026-11-15
- **Deciders:** implementation track of V1-06, implementing decision D7 of
  `spec/settings-and-logging/v1/overview.md`, approved by the owner on
  2026-11-11

## Context

`inspect --format yaml` must write one YAML document that a YAML 1.2 parser
reads as the same value as the JSON form (OUT-003). `output.yaml` in the
contract fixes the shape: block-style YAML in which every string and every
mapping key is a double-quoted scalar with JSON escapes, numbers, `true`,
`false` and `null` are plain, and empty collections are `[]` and `{}`.

The contract decided not to add a dependency for this (D7). The only YAML
libraries for Go are general encoders whose quoting, indentation and
key-order rules would have to be configured and pinned to meet that shape,
and the module otherwise needs no YAML. Each dependency also adds supply-chain
review work (chapter 13).

## Decision

Implement the encoder in `internal/output/yaml.go`:

- Marshal the result to JSON with `encoding/json`, so that field names,
  omission rules and values are exactly those of the JSON form.
- Read that JSON back token by token, keeping object member order.
- Write block mappings and sequences with two-space indentation; a mapping
  inside a sequence starts on the `- ` line; a nested sequence starts on the
  next line.
- Quote every string with the escapes that JSON and YAML 1.2 share, and write
  characters that YAML does not allow unescaped (C0 and C1 controls, DEL,
  U+2028, U+2029, BOM) as `\uXXXX`.

## Consequences

- No new module dependency; the binary stays small and static.
- The encoder handles only what JSON can express, which is all a result can
  hold. It is not a general YAML writer and is not exported for other use.
- Correctness rests on golden tests in `internal/output` and on the
  acceptance suite, which parses the output with a vendored YAML 1.2 parser
  (`ruamel.yaml` 0.19.1) and compares it with the JSON form.
- Member order in YAML follows the JSON form, so the two are easy to diff.
