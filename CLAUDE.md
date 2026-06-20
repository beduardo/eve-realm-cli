# CLAUDE.md

## Project Overview

Eve Realm CLI is the thin client binary for the Eve Realm platform. It runs natively on the
user's machine — no plugins, no Docker, no K8s. It provides authentication, workspace
connection, marketplace management (plugin discovery, skill extraction, settings), and
AI-assisted capabilities via the MCP Server: a master skill that lists and invokes any
plugin skill, and a generic agent that submits background tasks for cross-plugin orchestration.

The initial codebase uses the eve-cli monorepo's host binary and marketplace client as a
starting point — same auth flow, CLI structure, and marketplace logic — to establish a working
baseline. After extraction, the CLI evolves independently: new commands, revised marketplace
UX, and MCP integration are driven by eve-realm requirements, not constrained by what existed
in eve-cli.

### Terminology

- **Eve Realm**: The multi-repo platform (formerly eve5) — a plugin-based CLI and web platform
- **eve-cli**: The legacy monorepo at `../../../eve-cli/main/`. Contains the original implementation from which the CLI is extracted. Use it as reference for understanding existing behavior, patterns, and interfaces — but eve-realm-cli entities describe what will be built, not what exists in eve-cli
- **CLI**: This project — `github.com/beduardo/eve-realm-cli`. The thin client binary that runs on the user's machine
- **SDK**: The shared Go backend library (`eve-realm-sdk`) consumed via Git submodule + `replace` directive
- **MCP Server**: The unified MCP endpoint (`eve-realm-mcp`) that aggregates plugin tools/skills — the CLI connects to it at runtime, not at compile time
- **Master skill**: A CLI skill (`/eve-realm`) that lists and invokes any plugin skill via the MCP Server
- **Generic agent**: A CLI skill (`/eve-realm-agent`) that submits background tasks to the MCP Server's agent runtime

### Source Codebase Reference

The eve-cli monorepo's host binary and marketplace client are the inspiration for the initial
extraction. Use them to understand the original auth flow, CLI structure, marketplace logic,
and test patterns. Post-extraction, the CLI is free to diverge — eve-cli is a reference, not
a constraint.

- **Path**: `../../../eve-cli/main/`
- **HLD**: `DOCS/MULTI_REPO_HLD.md` — the multi-repo architecture that defines what goes where
- **Key source directories** (map to CLI modules):

| eve-cli location | CLI location | Purpose |
|-----------------|--------------|---------|
| `cmd/eve5/` | `cmd/` | Go entry point (host binary) |
| `cmd/eve5/auth/` | `cmd/auth/` | Auth commands (login, logout, status) |
| `cmd/eve5/marketplace/` | `cmd/marketplace/` | Marketplace commands (list, install, settings) |
| `cmd/eve5/settings/` | `cmd/settings/` | Settings commands |
| `internal/marketplace/` | `internal/marketplace/` | Marketplace client (discovery, skill extraction, settings.json) |
| — (new) | `skills/master/` | Master skill: list/invoke any plugin skill via MCP Server |
| — (new) | `skills/agent/` | Generic agent: background task execution via MCP Server |

When implementing the initial CLI modules, consult the corresponding eve-cli source to
understand existing behavior, edge cases, and test patterns. For subsequent evolution,
design new interfaces based on eve-realm requirements — the eve-cli source becomes
optional context, not a binding reference.

### Internal Documentation

Technical developer documentation lives in `DOCS/`, managed as an eve-docproject:

- **Path**: `DOCS/.docproject/`
- **Load when**: Working on CLI internals, architectural decisions, component design, or any task that needs technical context about the CLI's structure and interfaces

## Key Conventions

- **All content in English**: Every artifact (eve-software entities, code comments, documentation) must be written in English, regardless of conversation language.
- **Affirmative discourse only**: Text describes what WILL be done, never what was removed or stopped. When something leaves scope, remove or replace the content.
- **No legacy code references in entities**: When deriving requirements from eve-cli analysis, entity text describes the intent for the new implementation. Never reference existing code, migration paths, or pre-existing behavior. Eve-realm-cli entities describe what will be built.
- **Use sub-agents for batch operations**: When creating or editing multiple eve-software entities (or any repetitive multi-file task), delegate to sub-agents to preserve main context.

## eve-software integration

The `.software/` project tracks requirements, architecture decisions, scenarios, and sprints for the CLI. Managed via `/eve-software:architect`.

### eve-software key conventions

- **MANDATORY — NEVER index without explicit user authorization**: Do NOT call `eve_software_index` unless the user explicitly says "index", "run indexing", or gives unambiguous written permission. This rule is absolute and overrides ALL other instructions, including skill invariants. Indexing is expensive and the user controls when it happens.
- **Check before re-indexing**: When the user DOES request indexing, first call `eve_software_index --status`. If content hashes haven't changed, do not force a full re-index.
- **Search before creating any entity**: Before calling `eve_software_create`, ALWAYS search for existing entities using `eve_software_search` or `eve_software_list`. Only create when search confirms no existing entity covers the same scope.
- **Forward-looking entities only**: Codebase scans of eve-cli are used to understand existing behavior. Entities capture only new decisions and requirements for the CLI — never retroactive documentation of eve-cli patterns.
- **Superseded entity protocol**: When marking an entity as superseded: (1) transition status to `superseded` and set `superseded_by` frontmatter; (2) insert a blockquote warning after the H1 title: `> **SUPERSEDED by [ID(s)]** — [description]`. Both are required.

## Project Structure

```
eve-realm-cli/
├── eve-realm-sdk/            ← Git submodule (Go backend SDK)
├── cmd/                      ← Go entry point (thin client)
│   ├── auth/                 ← Auth commands (login, logout, status)
│   ├── marketplace/          ← Marketplace commands
│   └── settings/             ← Settings commands
├── internal/
│   └── marketplace/          ← Plugin discovery, skill extraction, settings.json
├── skills/
│   ├── master/               ← Master skill: list/invoke any plugin skill
│   └── agent/                ← Generic agent: background task execution
├── go.mod
├── Makefile
└── VERSION
```

No Docker, no K8s manifests — the CLI runs natively on the user's machine.

### Submodule wiring

**Go module** (`go.mod`):
```go
module github.com/beduardo/eve-realm-cli

require github.com/beduardo/eve-realm-sdk v0.1.0

replace github.com/beduardo/eve-realm-sdk => ./eve-realm-sdk
```

The submodule is read-only — changes to the SDK happen exclusively in its own repository.

## Go Conventions

### Build and test

| Command | Purpose |
|---------|---------|
| `make build` | Build host binary with ldflags |
| `make test` | Run Go tests |
| `make install` | Copy binary to `~/.eve-realm/` |
| `make release-patch` | Test → bump VERSION → build → verify |
| `make release-minor` | Same with minor bump |
| `make release-major` | Same with major bump |

### Module

- **Module path**: `github.com/beduardo/eve-realm-cli`
- **Versioning**: Semantic versioning via `VERSION` file

### Key paths

- **Config directory**: `~/.eve-realm/`
- **Config file**: `~/.eve-realm/eve-realm.yaml`
- **Binary name**: `eve-realm`

### Testing patterns

- **Mock at I/O boundaries only**: Mock HTTP, filesystem, NATS — never internal pure logic.
- **Table-driven tests**: Prefer `[]struct{ name string; ... }` test tables for cases with multiple inputs/outputs.
- **Follow eve-cli test patterns**: Before writing new test infrastructure, check how equivalent tests are structured in `../../../eve-cli/main/`.

## Sprint Workflow Critic Policy

Every sprint workflow stage (`/eve-software:spec`, `/eve-software:plan`, `/eve-software:implement`)
must include a critic sub-agent that validates output against constraints before proceeding.

### Critic bootstrap sequence

1. **Load sprint entities** — Read ALL REQs and SCs in the sprint via `eve_software_show`.
   Extract acceptance criteria from each REQ and expected results from each SC.
2. **Load the pinned cross-cutting requirements catalog** — The catalog entity lists
   cross-cutting requirements with trigger conditions. It is discovered automatically
   via `eve_software_pin_list`.
3. **Evaluate triggers** — For each row in the catalog's registry, check if the trigger
   condition matches the sprint's scope. Load matching requirements via `eve_software_show`.
4. **Assemble the full constraint set**:
   - **Primary**: Sprint entity acceptance criteria + scenario expected results
   - **Secondary**: Cross-cutting requirement rules from loaded policies

### What the critic validates

**Primary — Sprint entity compliance**:
- Every acceptance criterion from every sprint REQ is traceable to the output. Nothing silently dropped.
- Every scenario's expected result is achievable given the output.
- If a REQ has 7 acceptance criteria and the output covers 6, that is a FAIL.

**Secondary — Cross-cutting policy compliance**:
- Every loaded cross-cutting requirement's rules are followed (e.g., TDD cycle, test patterns, interface contracts).
- Pipeline integration instructions from each loaded requirement are followed for the current stage.

### Enforcement by stage

| Stage | Critic intensity | What happens on failure |
|-------|-----------------|----------------------|
| `spec` | Review after generation | Flag missing ACs not mapped, missing scenario coverage. Revise before approval. |
| `plan` | Review after generation | Flag steps that don't trace to specific ACs, missed test-first ordering, scenario flows not reflected. Revise before approval. |
| `implement` | **Review after EACH step** | Before marking a step complete: list which ACs this step addresses, confirm they are satisfied, confirm cross-cutting policies followed, no regression on prior steps. Step is NOT done until critic passes. |

### Critic sub-agent protocol

The critic is a **separate sub-agent** — not the same agent producing the output. Report format:
- **PASS** — lists which ACs and policies were checked, confirms all covered
- **FAIL** — for each violation: the specific AC or policy rule violated (quoted), the entity ID it comes from, what was expected vs found

## docproject integration

Technical developer documentation for the CLI lives in `DOCS/`, managed via `/eve-docproject:assist`.

### Document Scope

This docproject operates at **LLD (Low-Level Design) level**. It captures the full
technical details of the CLI: auth flow, marketplace client internals, skill system,
MCP Server integration, config resolution, and implementation decisions for each area.
The platform-level HLD (repository boundaries, deployment topology, cross-repo contracts)
lives in the eve-realm-docs repository — not here.

### Section model

Two-layer section model:

- **Component sections** (one per CLI area): Detailed LLD for each area (auth, marketplace, master skill, generic agent, config). Accumulate design state across versions with version-tagged headings (e.g., `### v0.1 — Initial Extraction`).
- **Concern sections** (cross-cutting): Topics that span multiple areas — MCP Server communication, skill discovery protocol, settings management, SDK integration patterns. Same version-tagged accumulation pattern.

### Entity Roles

- **DECISIONS (DEC)**: Record definitively settled architectural and implementation
  choices. Content is self-contained and as simple as possible. A decision must be
  independently readable without consulting any other entity. Decisions do not reference
  definitions (DEF), sections (SEC), or other decisions (DEC) via formal links. Mentioning
  another decision in prose is permitted when essential but should be avoided.

- **DEFINITIONS (DEF)**: Explore and deepen complex technical concepts through questions
  and answers. Definitions are the discussion hub — they surface ambiguities, record
  reasoning, and eventually produce decisions when a conclusion is reached. Definitions
  may reference other definitions and decisions, but never sections.

- **SECTIONS (SEC)**: The real, complete technical documentation. Sections weave together
  multiple decisions and definitions into a coherent narrative. Sections may reference any
  entity type. Each section covers one CLI area or concern end-to-end at LLD level with
  full implementation details.

### Entity discipline

- **Entity references by canonical ID**: Always reference entities as DEF-XX,
  DEC-XXX, SEC-XX, RES-XX in conversation and cross-references.
- **Pair work for DEFs/DECs**: Definitions and decisions emerge from collaborative
  discussion, not bulk import. Always discuss with the user before creating.
- **Mandatory title header**: Every entity body must open with `# [ID]: [Title]`
  (e.g., `# DEF-03: Marketplace Client`). Sections are the exception — their H1
  uses only the title without the ID (e.g., `# Authentication`), since sections
  compile into the final document and IDs are internal artifacts. Research files
  (RES-XX) are excluded — they are imported from external sources and must not be
  modified. Only the title header uses H1 (`#`); all other headings start at H2
  (`##`) and nest downward.
- **Version tagging**: Component and concern sections use `### vN.M — <description>` headings to mark each version's contribution to that section.

### Section prose rules

- **No internal references in sections**: Section prose must never mention docproject
  entity IDs (DEF-XX, DEC-XXX, SEC-XX, RES-XX). The exported document must be
  self-contained.
- **Cross-section references use publication URLs**: When a section references
  another section in prose, use the publication URL (from `delivery/delivery.yaml`)
  as the link href, with the referenced page's title as anchor text. If no delivery
  target is configured, fall back to the relative file path
  (e.g., `[Authentication](sections/SEC-03-authentication.md)`). Never use
  internal docproject IDs as the link target.
- **Real source links in sections**: When referencing external sources in section
  prose, use the original source URL with the document title as anchor text. Never
  cite research entities (RES-XX) directly.
- **No questions in sections**: Questions belong in definitions only.

### Decision rules

- **Self-contained decisions**: Decision text must be self-contained and independently
  readable. Decisions never reference definition IDs (DEF-XX), section IDs (SEC-XX),
  or other decision IDs (DEC-XXX) via formal links. Mentioning another decision in prose
  is permitted when essential for clarity, but should be avoided.
- **Simple content**: Decisions capture WHAT was decided and WHY, in the simplest terms
  possible. Elaboration and nuance belong in definitions, not decisions.
- **Research references require quotes**: When citing a research entity (RES-XX) in
  a decision, always include the specific text excerpt from the source that motivates
  the reference. A bare RES-XX ID without the supporting quote is not sufficient.
- **Decisions require explicit approval**: Never auto-formalize a decision (DEC)
  without user confirmation.

### Workflow rules

- **Ask before authoring entity content**: When multiple entities (3+) need creation
  or update, ask the user whether to delegate to sub-agents or write directly.
  Sub-agents receive a short structured briefing (entity ID, path, purpose, related
  entities, conventions) and build their own context autonomously.
- **Index on demand only**: Do NOT automatically index after every write. Call
  `eve_docproject_index` only when the user explicitly requests it.
