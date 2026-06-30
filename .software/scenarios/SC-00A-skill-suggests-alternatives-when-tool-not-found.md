---
content_hash: 8a83089d68d755ef4bc2086056b0d83d3e99be1e9a2842c17233427bca003090
created: "2026-06-27"
id: SC-00A
related_changes: []
related_reqs:
    - REQ-005
related_testcases: []
source: manual
status: implemented
tags:
    - master-skill
    - error
    - not-found
title: Skill suggests alternatives when tool not found
type: error-path
updated: "2026-06-29"
---

# SC-00A: Skill suggests alternatives when tool not found

## Preconditions

- MCP Server running with multiple tools registered (e.g., `ping`, `echo`)
- Master skill configured with the correct MCP Server address

## Steps

1. Invoke the master skill with tool name `"nonexistent"` and empty input

## Expected Result

- Error message states the tool `"nonexistent"` was not found
- Error message lists available tools as alternatives (names of all registered tools)
- User can see what tools are available and choose the correct one
