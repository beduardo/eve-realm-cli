---
content_hash: 1e2eb8b6094d0e8306ba7487e02dcbaa54d26444bf49c12b0e9e054376ab7307
created: "2026-06-20"
id: REQ-002
priority: high
related_adrs: []
related_changes: []
related_scenarios:
    - SC-001
    - SC-002
related_testcases: []
related_userstories: []
source: manual
status: blocked
tags:
    - process
    - release
    - quality
    - cross-cutting
title: Sprint completion and release process
updated: "2026-06-22"
---

# REQ-002: Sprint completion and release process

## Description

This is a **cross-cutting entity** that defines the mandatory quality and release process for every sprint. It covers decisions made at spec time, the post-implementation release sequence, and artifact management.

### Phase 1 — Spec-time decisions

At the start of every sprint specification, the spec agent must ask the user:

1. **Version increment**: Is this sprint a `major`, `minor`, or `patch` release? Record the answer in the SPEC.md frontmatter or header.
2. **README.md update**: Does this sprint introduce user-facing features or changes that warrant a README.md update? If yes, capture what the user wants to say. This decision drives step 5 of the release sequence.

### Phase 2 — Release sequence (after all tests pass)

After implementation is complete and `go test ./...` passes, execute these steps **in this exact order**:

#### Step 1: Commit implementation code

Stage and commit all implementation artifacts in a single commit:
- Source code changes (`.go` files, test files)
- Sprint artifacts (`SPEC.md`, `PLAN.md`, `IMPLEMENTATION.md`)
- Modified eve-software entities (`.software/` files changed during the sprint)

```
git add <changed files>
git commit -m "SP-XXX: <sprint title>"
```

#### Step 2: Run the release pipeline

Run the appropriate make target based on the version increment decided in Phase 1:

```
make release-patch   # or release-minor / release-major
```

This pipeline executes: `test → bump VERSION → build → install → version verify`. Production binaries are installed to `/usr/local/bin/` (standard macOS CLI location). Development builds use `dist/`. **No binaries should ever be left in the project root.**

If tests fail, the pipeline stops — no bump or build happens. Fix the issue and retry from Step 1.

#### Step 3: Collect release metadata

After the pipeline succeeds, collect:
- **Version**: `cat VERSION`
- **Git hash**: `git rev-parse --short HEAD` (points to the implementation commit from Step 1 — this is the hash embedded in the binaries)
- **Date**: current date in ISO-8601 format

#### Step 4: Append to RELEASE.md

Append a new entry to `RELEASE.md` (**append-only** — do not read or modify existing content). Use this template:

```markdown
## <version> — <date> (git: <hash>)

### Sprint: SP-XXX — <sprint title>

**Version increment**: <patch|minor|major>

**Changes**:
- <bullet list of key changes from the sprint>

**Entities affected**: <comma-separated entity IDs>

---
```

If RELEASE.md does not exist, create it with a top-level heading:

```markdown
# Changelog

---

## <version> — <date> (git: <hash>)
...
```

#### Step 5: Update README.md (conditional)

If the user decided at spec time that README.md should be updated, apply the changes now using the text captured in Phase 1. Only modify the relevant sections — do not rewrite the entire file.

#### Step 6: Commit release artifacts

Stage and commit the release artifacts:

```
git add VERSION RELEASE.md
# Include README.md only if it was updated
git commit -m "Bump version to <version>"
```

#### Step 7: Re-register marketplace

Refresh the Claude Code skill registration with the newly built binaries:

```
eve-realm marketplace register
```

### Build artifact rules

- **Production binaries** (`/usr/local/bin/`): installed by `make install` and `make release-*`
- **Development binaries** (`dist/`): written by `make build`
- **Project root**: must never contain generated binaries. If a `go build` command is run during verification, use `-o /tmp/` or `-o dist/` as the output path

## Acceptance Criteria

- Given a sprint spec is being generated, when the agent starts, then it asks the user for version increment (major/minor/patch) and whether README.md needs updating
- Given all tests pass after implementation, when the release sequence runs, then it follows Steps 1–7 in the exact order defined above
- Given the release pipeline runs, when binaries are built, then they are placed in `/usr/local/bin/` (production via `make install`) or `dist/` (development via `make build`), never in the project root
- Given the release completes, when RELEASE.md is updated, then the new entry is appended at the end without reading or modifying existing content
- Given README.md was flagged for update at spec time, when Step 5 runs, then only the relevant sections are modified with the user-provided text
- Given all release artifacts are committed, when `eve-realm marketplace register` runs, then the updated skills are available in the next Claude Code session

## Notes

- The git hash in RELEASE.md and in the binary both reference the implementation commit (Step 1), not the release commit (Step 6). This is intentional — the binary was built from the implementation commit.
- The `make release-*` pipeline is atomic: if tests fail, nothing is bumped or built.
- This requirement complements the TDD Strategy requirement. Together, they define the quality gates for every sprint.
