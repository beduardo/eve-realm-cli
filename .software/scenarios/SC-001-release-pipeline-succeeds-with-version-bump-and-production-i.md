---
content_hash: c1e260e780eaf46f8b935603276c804ab998446c77c1d8e2938a0ffe598367fa
created: "2026-06-22"
id: SC-001
related_changes: []
related_reqs:
    - REQ-002
related_testcases: []
source: manual
status: validated
tags:
    - release
    - build
    - makefile
title: Release pipeline succeeds with version bump and production install
type: happy-path
updated: "2026-06-22"
---

# SC-001: Release pipeline succeeds with version bump and production install

## Preconditions

- The project has compilable Go source code in `cmd/`
- All tests pass (`go test -count=1 ./...` exits 0)
- `VERSION` contains a valid semver string (e.g., `0.1.0`)
- `/usr/local/bin/` is writable by the current user

## Steps

1. Record the current version: `cat VERSION` → e.g., `0.1.0`
2. Run `make release-patch`
3. Observe the pipeline execution: test → bump → build → install → verify

## Expected Result

- Tests run and pass before any version change
- `VERSION` file is incremented to `0.1.1`
- Binary `dist/eve-realm` exists and is executable
- Binary `/usr/local/bin/eve-realm` exists and is executable
- Running `/usr/local/bin/eve-realm version` outputs `0.1.1`
- The embedded git hash and build date match the current commit and date
- No binary is left in the project root directory
