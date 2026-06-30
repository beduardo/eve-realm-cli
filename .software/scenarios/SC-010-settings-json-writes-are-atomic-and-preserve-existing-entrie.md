---
content_hash: eb9b5d7c919fb51765eb4f9fa9b0ffd90da4c2fcb3dbf556dea02c56c16e4656
created: "2026-06-28"
id: SC-010
related_changes: []
related_reqs:
    - REQ-007
related_testcases: []
source: manual
status: validated
tags:
    - marketplace
    - settings
    - atomic-write
title: Settings.json writes are atomic and preserve existing entries
type: happy-path
updated: "2026-06-29"
---

# SC-010: Settings.json writes are atomic and preserve existing entries

## Preconditions

- `~/.claude/settings.json` exists with pre-existing settings (e.g., `"autoMemoryEnabled": false`, other `extraKnownMarketplaces` entries)

## Steps

1. Call `UpdateClaudeSettings()` to add the `eve-realm` marketplace entry

## Expected Result

- The `eve-realm` entry is added under `extraKnownMarketplaces`
- All pre-existing settings keys and values are preserved exactly
- Other `extraKnownMarketplaces` entries are preserved
- The write uses a temporary file (`.settings-*.tmp`) in the same directory followed by an atomic rename
- If the process crashes mid-write, the original `settings.json` is not corrupted (either old or new content, never partial)
