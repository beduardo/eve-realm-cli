---
content_hash: c196730e4ddd5ca6a3d0ad9c17a144dca013993cde139d3847ba8b6aed9441a1
created: "2026-06-22"
id: SC-002
related_changes: []
related_reqs:
    - REQ-002
related_testcases: []
source: manual
status: automated
tags:
    - release
    - build
    - makefile
title: Release pipeline aborts when tests fail
type: error-path
updated: "2026-06-29"
---

# SC-002: Release pipeline aborts when tests fail

## Preconditions

- The project has compilable Go source code in `cmd/`
- At least one test fails (`go test -count=1 ./...` exits non-zero)
- `VERSION` contains a valid semver string (e.g., `0.1.0`)

## Steps

1. Record the current version: `cat VERSION` → e.g., `0.1.0`
2. Run `make release-patch`
3. Observe the pipeline execution stops at the test step

## Expected Result

- Tests run and fail
- `VERSION` file remains unchanged at `0.1.0` (no bump)
- No binary is produced in `dist/`
- No binary is copied to `/usr/local/bin/`
- The make command exits with a non-zero status code
