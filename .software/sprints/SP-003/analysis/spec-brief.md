# Spec Writer Brief

**Sprint**: SP-003
**Sprint Title**: Cobra command for MCP tool listing and invocation
**Project Root**: /Users/bruno/repo-pessoal/eve-realm/eve-realm-cli/main
**Sprint Folder**: /Users/bruno/repo-pessoal/eve-realm/eve-realm-cli/main/.software/sprints/SP-003
**Date**: 2026-06-30

## Entity List

| ID | Type | Title | Partial | Scope Notes |
|----|------|-------|---------|-------------|
| REQ-006 | requirement | Cobra command for MCP tool listing and invocation | no | - |
| SC-012 | scenario | Tools list outputs tool descriptors to stdout | no | - |
| SC-013 | scenario | Tools invoke sends default empty input and returns JSON to stdout | no | - |
| SC-014 | scenario | Tools invoke passes --input flag value to the tool | no | - |
| SC-015 | scenario | Tools commands exit non-zero with error when server unreachable | no | - |
| SC-016 | scenario | Tools invoke exits non-zero and lists alternatives when tool not found | no | - |
| SC-017 | scenario | MCP Server address resolved from env var, config file, or default | no | - |

## Analysis Artifacts

- Codebase Analysis: /Users/bruno/repo-pessoal/eve-realm/eve-realm-cli/main/.software/sprints/SP-003/analysis/codebase-analysis.md
- Feasibility Reports:
  - REQ-006: /Users/bruno/repo-pessoal/eve-realm/eve-realm-cli/main/.software/sprints/SP-003/analysis/feasibility-REQ-006.md

## Project Context

EVE Realm CLI project with 32 entities (7 requirements, 23 scenarios, 1 decision, 1 change). Sprint SP-001 (MCP gRPC Connection Client) is completed — it established the gRPC client library that SP-003 builds upon. The CLI is a thin native binary (`eve-realm`) using Cobra for command structure, with config at `~/.eve-realm/eve-realm.yaml`. Module path: `github.com/beduardo/eve-realm-cli`. Build via Makefile to `dist/`. SDK consumed via Git submodule with `replace` directive.

## Pinned Entities

The entities below are binding project policies. Every agent in every sprint phase (spec, plan, implement) MUST extract and follow all directives from these entities that are relevant to its phase. Report compliance in a Pinned Entity Compliance table.

### REQ-003: Cross-cutting requirements catalog for lazy-loaded sprint policy injection

| ID | Title | Trigger condition | Summary |
|----|-------|-------------------|---------|
| REQ-001 | Test-Driven Development Strategy | **Implementing or modifying Go code** in any sprint step | Defines the red→green→refactor TDD cycle, Go test framework rules (`testing` stdlib only), test patterns (table-driven, temp dirs, interface mocking, process substitution, HTTP tests, YAML round-trip, Cobra command tests), file naming conventions, and pipeline integration (spec writer generates test expectations, plan propagates them, implementer writes tests first, verifier validates coverage). |
| REQ-002 | Sprint completion and release process | **Completing a sprint and preparing a release** (typically the final steps of an implementation) | Defines the two-phase release process: Phase 1 (spec-time decisions: version increment, README update) and Phase 2 (post-implementation release sequence: commit → `make release-*` → collect metadata → append RELEASE.md → conditional README update → commit release artifacts → `eve-realm marketplace register`). Also defines build artifact placement rules. |

## Flags

- readme_update_needed: true
