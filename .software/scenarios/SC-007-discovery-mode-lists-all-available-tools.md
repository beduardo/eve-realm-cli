---
content_hash: 266fdd34d03fad744d2e2973d9af928df234cdb7d6956ef2379a0345b6876f18
created: "2026-06-27"
id: SC-007
related_changes: []
related_reqs:
    - REQ-005
related_testcases: []
source: manual
status: implemented
tags:
    - master-skill
    - discovery
title: Discovery mode lists all available tools
type: happy-path
updated: "2026-06-29"
---

# SC-007: Discovery mode lists all available tools

## Preconditions

- MCP Server running with multiple tools registered
- Master skill configured with the correct MCP Server address

## Steps

1. Invoke the master skill without arguments (discovery mode)

## Expected Result

- Returns formatted text listing each tool with its name, description, and input schema
- Format is suitable for AI consumption (clear, structured, parseable)
- All registered tools are included in the listing
