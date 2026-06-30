---
content_hash: 64a19f68f14949b3fbc4f6565d6f59fb69e40bfe8817b94fa299d615506f51da
created: "2026-06-27"
id: SC-008
related_changes: []
related_reqs:
    - REQ-005
related_testcases: []
source: manual
status: implemented
tags:
    - master-skill
    - invocation
title: Invocation mode returns tool response
type: happy-path
updated: "2026-06-29"
---

# SC-008: Invocation mode returns tool response

## Preconditions

- MCP Server running with the `ping` tool registered
- Master skill configured with the correct MCP Server address

## Steps

1. Invoke the master skill with tool name `"ping"` and empty input

## Expected Result

- Returns the tool's JSON response directly (e.g., `{"message": "pong", ...}`)
- No wrapper or additional formatting around the tool response
