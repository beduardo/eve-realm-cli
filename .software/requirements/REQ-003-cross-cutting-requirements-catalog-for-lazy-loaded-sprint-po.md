---
content_hash: 382100806090bf3d6ec3d218c5e96177d46bf8ae504ab03bf6cf26101df3b55a
created: "2026-06-20"
id: REQ-003
priority: high
related_adrs: []
related_changes: []
related_scenarios: []
related_testcases: []
related_userstories: []
source: manual
status: draft
tags:
    - pinned
    - cross-cutting
    - catalog
title: Cross-cutting requirements catalog for lazy-loaded sprint policy injection
updated: "2026-06-20"
---

# REQ-003: Cross-cutting requirements catalog for lazy-loaded sprint policy injection

## Description

This is the **single pinned entity** for the project. It replaces direct pinning of individual cross-cutting requirements with a lightweight catalog that agents load once, then selectively fetch only the requirements relevant to their current task.

### How it works

This entity is pinned in `software.yaml`. When sprint agents (spec writer, plan generator, step implementer, step verifier) receive the memory bundle, they get this catalog — not the full body of every cross-cutting requirement. Each agent then evaluates the trigger conditions below and calls `eve_software_show` to load the full requirement **only when the trigger matches**.

### Mandatory loading rule

**If a trigger condition below matches what you are about to do in this sprint or step, you MUST call `eve_software_show <ID>` to load the full requirement before proceeding. This is not optional. Skipping a matching requirement is a step failure.**

When in doubt about whether a trigger applies, load the requirement — the cost of loading one extra entity is far lower than the cost of violating a project-wide policy.

### Cross-cutting requirements registry

| ID | Title | Trigger condition | Summary |
|----|-------|-------------------|---------|
| REQ-001 | Test-Driven Development Strategy | **Implementing or modifying Go code** in any sprint step | Defines the red→green→refactor TDD cycle, Go test framework rules (`testing` stdlib only), test patterns (table-driven, temp dirs, interface mocking, process substitution, HTTP tests, YAML round-trip, Cobra command tests), file naming conventions, and pipeline integration (spec writer generates test expectations, plan propagates them, implementer writes tests first, verifier validates coverage). |
| REQ-002 | Sprint completion and release process | **Completing a sprint and preparing a release** (typically the final steps of an implementation) | Defines the two-phase release process: Phase 1 (spec-time decisions: version increment, README update) and Phase 2 (post-implementation release sequence: commit → `make release-*` → collect metadata → append RELEASE.md → conditional README update → commit release artifacts → `eve-realm marketplace register`). Also defines build artifact placement rules. |

### Extensibility

This catalog is the **single entry point** for all cross-cutting requirements. No other entity should be pinned directly — only REQ-003 is pinned.

#### When to add an entry

A requirement qualifies for this catalog when it defines a **project-wide policy** that sprint agents must follow conditionally — i.e., the requirement applies not to a specific feature but to a class of actions across any sprint. Signals:
- The requirement uses words like "all sprints", "every step", "cross-cutting", "policy", "strategy", "convention"
- It instructs spec/plan/implement/verify agents on how to behave
- It would previously have been pinned directly

Requirements that apply to a single sprint or feature do **not** belong here — they flow through normal sprint entity inclusion.

#### How to add an entry

1. **Create or identify the cross-cutting requirement** using standard `eve_software_create` or `eve_software_show`. The requirement must already exist before adding it to the catalog.

2. **Write the catalog row** by appending a new row to the "Cross-cutting requirements registry" table in this entity's body. Each row has four columns:

   | Column | Content | Example |
   |--------|---------|---------|
   | **ID** | The entity ID (e.g., `REQ-060`) | `REQ-060` |
   | **Title** | The entity's title, verbatim | `Database migration conventions` |
   | **Trigger condition** | Action-based phrase describing **when** agents must load this requirement. Use bold text. Phrase as a gerund: "**Doing X**" or "**Modifying Y**". Be specific enough that an agent can evaluate it in one pass. | `**Creating or modifying database migrations**` |
   | **Summary** | 2-3 sentences describing the key rules, patterns, and constraints the requirement defines. Include enough detail that an agent can decide whether to load the full entity, but not so much that the catalog becomes heavy. | `Defines migration file naming, rollback requirements, and schema change review gates.` |

3. **Update this entity** via `eve_software_update` or direct file edit. Do not create a new catalog entity.

4. **Do NOT pin the new requirement** — it is discovered through this catalog's trigger mechanism. Only REQ-003 remains pinned.

#### How to remove an entry

When a cross-cutting requirement is deprecated or no longer applies:
1. Remove its row from the registry table
2. Transition the requirement entity to `deprecated` if appropriate
3. Update this entity

#### Architect skill protocol

When the `/eve-software:architect` skill detects that a newly created or updated requirement qualifies as cross-cutting (see "When to add an entry" above), it must:
1. Ask the user: "This requirement defines a project-wide policy. Should I add it to the cross-cutting catalog (REQ-003)?"
2. On confirmation, draft the trigger condition and summary, present them for review
3. On approval, edit this entity's registry table to include the new row
4. Confirm: "Added [ID] to the cross-cutting catalog. Sprint agents will load it when [trigger condition]."

## Acceptance Criteria

- Given a sprint memory bundle is assembled and REQ-003 is pinned, when an agent receives the bundle, then it contains only this catalog — not the full body of REQ-001 or REQ-002
- Given an agent is about to implement Go code in a sprint step, when it reads this catalog, then it identifies the REQ-001 trigger as matching and calls `eve_software_show REQ-001` to load the full TDD strategy before writing any code
- Given an agent is completing a sprint and preparing a release, when it reads this catalog, then it identifies the REQ-002 trigger as matching and calls `eve_software_show REQ-002` to load the full release process
- Given a new cross-cutting requirement is created in the future, when it is added to the registry table, then agents discover and load it through the same trigger-based mechanism without any pinning changes
- Given no trigger condition matches the agent's current task, when it reads this catalog, then it proceeds without loading any cross-cutting requirement — avoiding unnecessary token consumption

## Notes

- The trigger conditions are action-based ("implementing Go code", "completing a sprint") rather than agent-type-based ("step implementer", "spec writer") — this ensures correct loading regardless of which agent is performing the action.
- Token savings: instead of multiple heavy entities always loaded (~3000+ tokens each), agents load one lightweight catalog (~800 tokens) + only the relevant heavy entities on demand.
- Future additions anticipated: Claude Code integration validation (REQ-034 equivalent) when marketplace commands are implemented.
