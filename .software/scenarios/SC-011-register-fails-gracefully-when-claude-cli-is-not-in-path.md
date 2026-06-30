---
content_hash: 3dd0af746344533da057e7bf0c91778dc1874d3f569b1916e25cec4b40c352ad
created: "2026-06-28"
id: SC-011
related_changes: []
related_reqs:
    - REQ-007
related_testcases: []
source: manual
status: validated
tags:
    - marketplace
    - register
    - error
title: Register fails gracefully when claude CLI is not in PATH
type: error-path
updated: "2026-06-29"
---

# SC-011: Register fails gracefully when claude CLI is not in PATH

## Preconditions

- `claude` CLI is NOT available in PATH

## Steps

1. Run `eve-realm marketplace register`

## Expected Result

- Command exits with a non-zero status code
- Error message clearly states that the `claude` CLI was not found in PATH
- No files are extracted
- `~/.claude/settings.json` is not modified
