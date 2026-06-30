---
content_hash: 3e33f77c873d39931ee98a21cba9fbf968e752c1cceb0e8975ce30f5d93fd9d2
created: "2026-06-27"
id: SC-005
related_changes: []
related_reqs:
    - REQ-004
related_testcases: []
source: manual
status: implemented
tags:
    - grpc-client
    - error
    - connection
title: Client returns connection error when server unreachable
type: error-path
updated: "2026-06-29"
---

# SC-005: Client returns connection error when server unreachable

## Preconditions

- No MCP Server running on the target address (e.g., `localhost:59999`)

## Steps

1. Create a gRPC client with `NewClient("localhost:59999")`
2. Call `ListTools(ctx)` or `InvokeTool(ctx, "ping", "{}")`

## Expected Result

- Returns an error that is clearly a connection error, not a raw gRPC status code
- Error message mentions the unreachable address
- Error is distinguishable from tool-not-found errors via type assertion or `errors.Is`/`errors.As`
