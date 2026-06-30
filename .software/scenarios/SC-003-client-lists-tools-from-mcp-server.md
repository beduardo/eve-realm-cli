---
content_hash: 206709d81a362d6d9b20f5be6892c012a1a852e4513c67f0f04804afde2af6a5
created: "2026-06-27"
id: SC-003
related_changes: []
related_reqs:
    - REQ-004
related_testcases: []
source: manual
status: implemented
tags:
    - grpc-client
    - list-tools
title: Client lists tools from MCP Server
type: happy-path
updated: "2026-06-29"
---

# SC-003: Client lists tools from MCP Server

## Preconditions

- MCP Server running on a known address (e.g., `localhost:50051`)
- At least one tool registered in the MCP Server (e.g., `ping`)

## Steps

1. Create a gRPC client with `NewClient("localhost:50051")`
2. Call `ListTools(ctx)` on the client

## Expected Result

- Returns a `[]Tool` slice with at least one entry
- Each Tool struct has populated `name`, `description`, and `input_schema` fields
- No error returned
