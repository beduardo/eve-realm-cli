# Codebase Analysis Brief

**Sprint**: SP-001
**Project Root**: /Users/bruno/repo-pessoal/eve-realm/eve-realm-cli/main
**Entity IDs**: REQ-004, REQ-005, SC-003, SC-004, SC-005, SC-006, SC-007, SC-008, SC-009, SC-00A, SC-00B

## Entity Details

### REQ-004: gRPC tool client for MCP Server
- Type: requirement
- Status: active
- Tags: grpc, mcp-client, tool-registry
- Internal package at `internal/mcpclient/` providing Client type with ListTools and InvokeTool methods. Wraps MCPService proto contract. In the local k3d development environment, the MCP Server's gRPC port is exposed via NodePort 30051, so the default address is `localhost:30051`. Handles connection management, gRPC-to-CLI error mapping, response deserialization. Proto at `proto/mcp/v1/mcp.proto`, generated code at `gen/proto/mcp/v1/`.
- AC-6: The client accepts the MCP Server address as a constructor parameter (e.g., `NewClient(addr string)`), defaulting to `localhost:30051` when not specified.

### REQ-005: Master skill for MCP tool discovery and invocation
- Type: requirement
- Status: active
- Tags: skill, marketplace, master-skill, mcp
- Marketplace skill at `skills/master/` acting as gateway between AI tools and MCP Server. Two modes: discovery (ListTools formatted prompt) and invocation (InvokeTool pass-through). Reads MCP Server address from config or env var.

### SC-003: Client lists tools from MCP Server
- Type: scenario
- Status: validated
- Tests ListTools returning tool descriptors from MCP Server.

### SC-004: Client invokes tool and returns response
- Type: scenario
- Status: validated
- Tests InvokeTool with tool name and JSON input returning JSON output.

### SC-005: Client returns connection error when server unreachable
- Type: scenario
- Status: validated
- Tests clear connection error when MCP Server is unreachable.

### SC-006: Client returns typed error for unknown tool
- Type: scenario
- Status: validated
- Tests typed error distinguishable from connection errors for NOT_FOUND.

### SC-007: Discovery mode lists all available tools
- Type: scenario
- Status: validated
- Tests skill invocation without arguments returning formatted tool listing.

### SC-008: Invocation mode returns tool response
- Type: scenario
- Status: validated
- Tests skill invocation with tool name returning JSON response.

### SC-009: Skill handles unreachable MCP Server gracefully
- Type: scenario
- Status: validated
- Tests user-friendly error when MCP Server is down.

### SC-00A: Skill suggests alternatives when tool not found
- Type: scenario
- Status: validated
- Tests error message listing available tools when requested tool not found.

### SC-00B: End-to-end ping invocation via master skill
- Type: scenario
- Status: validated
- Tests end-to-end: skill invokes ping tool, gets pong response.

## Focus Areas
- Existing Go module structure (go.mod, go.sum), dependencies, and eve-realm-sdk submodule wiring
- Current cmd/ and internal/ package layout and patterns
- Makefile targets and build patterns
- Existing test patterns and infrastructure
- Config file handling (if any)
- Reference codebase at ../../../eve-cli/main/ for equivalent patterns
- Proto/gRPC tooling setup (any existing .proto files, buf.yaml, protoc config)
- Default address handling patterns (localhost:30051 default for k3d NodePort)
