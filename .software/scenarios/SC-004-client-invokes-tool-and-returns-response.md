---
content_hash: d9538037443d7baa591fc766576a11bae35c30ed535175415e4cc97d943d8990
created: "2026-06-27"
id: SC-004
related_changes: []
related_reqs:
    - REQ-004
related_testcases: []
source: manual
status: implemented
tags:
    - grpc-client
    - invoke-tool
title: Client invokes tool and returns response
type: happy-path
updated: "2026-06-29"
---

# SC-004: Client invokes tool and returns response

## Preconditions

- MCP Server running on a known address with a registered `ping` tool

## Steps

1. Create a gRPC client with `NewClient("localhost:50051")`
2. Call `InvokeTool(ctx, "ping", "{}")` on the client

## Expected Result

- Returns a JSON string containing the tool output (e.g., `{"message": "pong", ...}`)
- No error returned
