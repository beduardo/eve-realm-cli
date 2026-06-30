---
content_hash: 812dd4c325492d0adf236a082002bd829fa3dd5920c31f3396acbab59fb04e87
created: "2026-06-27"
id: REQ-004
priority: high
related_adrs: []
related_changes: []
related_scenarios:
    - SC-003
    - SC-004
    - SC-005
    - SC-006
related_testcases: []
related_userstories: []
source: manual
status: implemented
tags:
    - grpc
    - mcp-client
    - tool-registry
title: gRPC tool client for MCP Server
updated: "2026-06-29"
---

# REQ-004: gRPC tool client for MCP Server

## Description

The CLI implements a gRPC client library that connects to the MCP Server's tool registry
service. In the local k3d development environment, the MCP Server's gRPC port is exposed
via NodePort 30051, so the default address is `localhost:30051`. This client wraps the
`MCPService` proto contract, providing Go functions to list available tools and invoke
them by name.

The client is a reusable internal package (`internal/mcpclient/`) consumed by the master
skill (REQ-005) and any future CLI component that needs to interact with MCP Server tools.
It handles connection management, error mapping (gRPC status codes to CLI-friendly errors),
and response deserialization.

The proto definition is shared with the MCP Server project. For the initial implementation,
the proto file is copied locally (`proto/mcp/v1/mcp.proto`). Once the SDK submodule workflow
is established, both projects will consume the generated code from the SDK.

## Acceptance Criteria

1. An internal package at `internal/mcpclient/` provides a `Client` type with `ListTools(ctx) ([]Tool, error)` and `InvokeTool(ctx, name, input string) (string, error)` methods.
2. `ListTools` calls the MCP Server's `ListTools` gRPC RPC and returns a slice of tool descriptors (name, description, input schema).
3. `InvokeTool` calls the MCP Server's `InvokeTool` gRPC RPC with the given tool name and JSON input, returning the JSON output string.
4. When the MCP Server is unreachable, methods return a clear error indicating connection failure (not a raw gRPC error).
5. When `InvokeTool` receives a gRPC `NOT_FOUND` status, the client returns a typed error distinguishable from connection errors.
6. The client accepts the MCP Server address as a constructor parameter (e.g., `NewClient(addr string)`), defaulting to `localhost:30051` when not specified.
7. The proto definition at `proto/mcp/v1/mcp.proto` matches the MCP Server's definition. Generated Go code lives in `gen/proto/mcp/v1/`.
8. `make proto` generates the Go client stubs from the proto definition.
