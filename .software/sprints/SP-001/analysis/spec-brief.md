# Spec Writer Brief

**Sprint**: SP-001
**Sprint Title**: MCP gRPC Connection Client
**Project Root**: /Users/bruno/repo-pessoal/eve-realm/eve-realm-cli/main
**Sprint Folder**: /Users/bruno/repo-pessoal/eve-realm/eve-realm-cli/main/.software/sprints/SP-001
**Date**: 2026-06-28

## Entity List

| ID | Type | Title | Partial | Scope Notes |
|----|------|-------|---------|-------------|
| REQ-004 | requirement | gRPC tool client for MCP Server | no | - |
| REQ-005 | requirement | Master skill for MCP tool discovery and invocation | no | - |
| SC-003 | scenario | Client lists tools from MCP Server | no | - |
| SC-004 | scenario | Client invokes tool and returns response | no | - |
| SC-005 | scenario | Client returns connection error when server unreachable | no | - |
| SC-006 | scenario | Client returns typed error for unknown tool | no | - |
| SC-007 | scenario | Discovery mode lists all available tools | no | - |
| SC-008 | scenario | Invocation mode returns tool response | no | - |
| SC-009 | scenario | Skill handles unreachable MCP Server gracefully | no | - |
| SC-00A | scenario | Skill suggests alternatives when tool not found | no | - |
| SC-00B | scenario | End-to-end ping invocation via master skill | no | - |

## Analysis Artifacts

- Codebase Analysis: /Users/bruno/repo-pessoal/eve-realm/eve-realm-cli/main/.software/sprints/SP-001/analysis/codebase-analysis.md
- Feasibility Reports:
  - REQ-004: /Users/bruno/repo-pessoal/eve-realm/eve-realm-cli/main/.software/sprints/SP-001/analysis/feasibility-REQ-004.md
  - REQ-005: /Users/bruno/repo-pessoal/eve-realm/eve-realm-cli/main/.software/sprints/SP-001/analysis/feasibility-REQ-005.md

## Project Context

EVE Realm CLI is the thin client binary for the Eve Realm platform. It runs natively on the user's machine providing authentication, workspace connection, marketplace management, and AI-assisted capabilities via the MCP Server.

Project statistics: 17 total entities (5 requirements, 11 scenarios, 1 change). Statuses: 2 active, 3 blocked, 2 implemented, 1 in-progress, 9 validated.

Key paths:
- Module: `github.com/beduardo/eve-realm-cli`
- Binary: `eve-realm`
- Config: `~/.eve-realm/eve-realm.yaml`
- SDK submodule: `eve-realm-sdk/` (read-only, changes in its own repo)

Build: `make build` (to `dist/`), `make test`, `make install` (to `/usr/local/bin/eve-realm`), `make proto` (generate Go stubs from proto).

## Pinned Entities

### REQ-003: Cross-cutting requirements catalog for lazy-loaded sprint policy injection

**Registry:**

| ID | Title | Trigger condition | Summary |
|----|-------|-------------------|---------|
| REQ-001 | Test-Driven Development Strategy | **Implementing or modifying Go code** in any sprint step | Defines the red-green-refactor TDD cycle, Go test framework rules (`testing` stdlib only), test patterns (table-driven, temp dirs, interface mocking, process substitution, HTTP tests, YAML round-trip, Cobra command tests), file naming conventions, and pipeline integration (spec writer generates test expectations, plan propagates them, implementer writes tests first, verifier validates coverage). |
| REQ-002 | Sprint completion and release process | **Completing a sprint and preparing a release** (typically the final steps of an implementation) | Defines the two-phase release process: Phase 1 (spec-time decisions: version increment, README update) and Phase 2 (post-implementation release sequence: commit, make release, collect metadata, append RELEASE.md, conditional README update, commit release artifacts, marketplace register). Also defines build artifact placement rules. |

**Mandatory loading rule**: If a trigger condition matches what the spec writer is about to address, call `eve_software_show <ID>` to load the full requirement before proceeding.

## Flags

- readme_update_needed: true
