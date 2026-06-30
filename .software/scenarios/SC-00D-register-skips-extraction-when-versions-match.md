---
content_hash: 31e0d7cfcd9c7b1784f98c478a0ad0d6f0221f69141799b4df2f4d680be6eabc
created: "2026-06-28"
id: SC-00D
related_changes: []
related_reqs:
    - REQ-007
related_testcases: []
source: manual
status: validated
tags:
    - marketplace
    - register
    - version-check
title: Register skips extraction when versions match
type: happy-path
updated: "2026-06-29"
---

# SC-00D: Register skips extraction when versions match

## Preconditions

- `claude` CLI is available in PATH
- A previous `eve-realm marketplace register` has already extracted files to `<ConfigDir>/marketplace/`
- The installed `marketplace.json` version matches the embedded version

## Steps

1. Run `eve-realm marketplace register` again

## Expected Result

- `NeedsUpdate()` returns false
- No files are overwritten or re-extracted
- The `settings.json` entry is still verified/updated (idempotent)
- Command exits with zero status code
- Output indicates the marketplace is already up to date
