---
content_hash: 5aa6b5423f3f675b1241332f5c798a962cc520c748b250081ad9a28d7ce9f509
created: "2026-06-27"
id: SC-00B
related_changes: []
related_reqs:
    - REQ-005
related_testcases: []
source: manual
status: implemented
tags:
    - master-skill
    - ping
    - e2e
title: End-to-end ping invocation via master skill
type: happy-path
updated: "2026-06-29"
---

# SC-00B: End-to-end ping invocation via master skill

## Preconditions

- MCP Server running with the `ping` tool registered
- CLI configured with the correct MCP Server address (via config or env var)

## Steps

1. Invoke the master skill with tool `"ping"`

## Expected Result

- Returns `{"message": "pong", "timestamp": "<RFC 3339>"}` 
- Timestamp is recent (within seconds of invocation)
- Response is valid JSON
