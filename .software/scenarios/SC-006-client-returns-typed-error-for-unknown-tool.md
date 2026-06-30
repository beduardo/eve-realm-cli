---
content_hash: 273e719f3e02845751372748c6c0f4ccb95d3625d8ee4a80c48a9ad591086d73
created: "2026-06-27"
id: SC-006
related_changes: []
related_reqs:
    - REQ-004
related_testcases: []
source: manual
status: implemented
tags:
    - grpc-client
    - error
    - not-found
title: Client returns typed error for unknown tool
type: error-path
updated: "2026-06-29"
---

# SC-006: Client returns typed error for unknown tool

## Preconditions

- MCP Server running on a known address with at least one tool registered

## Steps

1. Create a gRPC client with `NewClient("localhost:50051")`
2. Call `InvokeTool(ctx, "nonexistent", "{}")`

## Expected Result

- Returns a typed `ToolNotFoundError` (or similar sentinel/typed error)
- Error is distinguishable from connection errors via type assertion or `errors.Is`/`errors.As`
- Error message includes the tool name that was not found
