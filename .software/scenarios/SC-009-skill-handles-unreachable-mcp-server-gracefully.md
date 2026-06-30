---
content_hash: 1450ede0fc0ab23f9527f75431debab48d820cc0fd05f3121ada9de409e6be14
created: "2026-06-27"
id: SC-009
related_changes: []
related_reqs:
    - REQ-005
related_testcases: []
source: manual
status: implemented
tags:
    - master-skill
    - error
    - connection
title: Skill handles unreachable MCP Server gracefully
type: error-path
updated: "2026-06-29"
---

# SC-009: Skill handles unreachable MCP Server gracefully

## Preconditions

- MCP Server not running or unreachable at the configured address

## Steps

1. Invoke the master skill in either mode (discovery or invocation)

## Expected Result

- Returns a user-friendly error message (e.g., "MCP Server is not available at <address>. Check that the server is running.")
- No stack traces or raw gRPC error codes exposed to the user
- Error message includes the configured server address for debugging
