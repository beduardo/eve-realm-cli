---
content_hash: a10275732e7ed1f88e98d3443c2e5e66c7c08d8dc46a227e3f8864d66dd1d70f
created: "2026-06-28"
id: SC-00C
related_changes: []
related_reqs:
    - REQ-007
related_testcases: []
source: manual
status: validated
tags:
    - marketplace
    - register
    - happy-path
title: Register extracts embedded files and updates Claude Code settings
type: happy-path
updated: "2026-06-29"
---

# SC-00C: Register extracts embedded files and updates Claude Code settings

## Preconditions

- `claude` CLI is available in PATH
- No previous marketplace installation exists at `<ConfigDir>/marketplace/`
- `~/.claude/settings.json` exists with other settings but no `extraKnownMarketplaces` entry for `eve-realm`

## Steps

1. Run `eve-realm marketplace register`

## Expected Result

- All embedded files are extracted to `<ConfigDir>/marketplace/` preserving the directory structure
- `<ConfigDir>/marketplace/.claude-plugin/marketplace.json` exists and contains `"name": "eve-realm"`
- `<ConfigDir>/marketplace/plugins/eve-realm/.claude-plugin/plugin.json` exists
- `<ConfigDir>/marketplace/plugins/eve-realm/skills/master/SKILL.md` exists with the master skill content
- `~/.claude/settings.json` contains `extraKnownMarketplaces.eve-realm` with `source.source: "directory"` and `source.path` pointing to the extracted directory
- All pre-existing settings in `settings.json` are preserved
- Command exits with zero status code
- Output lists available skills
