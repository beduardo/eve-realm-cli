---
content_hash: 62eb619d2fb0ab7bace3aaf8fb2ce12b53f3d10207448293c9bac7fd402a468a
created: "2026-06-28"
id: SC-00E
related_changes: []
related_reqs:
    - REQ-007
related_testcases: []
source: manual
status: validated
tags:
    - marketplace
    - unregister
title: Unregister removes marketplace from Claude Code settings and deletes files
type: happy-path
updated: "2026-06-29"
---

# SC-00E: Unregister removes marketplace from Claude Code settings and deletes files

## Preconditions

- Marketplace was previously registered (files extracted, `settings.json` entry present)
- `<ConfigDir>/marketplace/` exists with extracted files

## Steps

1. Run `eve-realm marketplace unregister`

## Expected Result

- The `extraKnownMarketplaces.eve-realm` entry is removed from `~/.claude/settings.json`
- All other settings in `settings.json` are preserved
- The `<ConfigDir>/marketplace/` directory and all its contents are deleted
- Command exits with zero status code
