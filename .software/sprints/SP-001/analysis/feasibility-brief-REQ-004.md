# Feasibility Brief

**Sprint**: SP-001
**Project Root**: /Users/bruno/repo-pessoal/eve-realm/eve-realm-cli/main
**Target Entity**: REQ-004

## Entity Summary
gRPC tool client for MCP Server — internal package at `internal/mcpclient/` providing Client type with ListTools and InvokeTool methods. Wraps MCPService proto contract. In the local k3d development environment, the MCP Server's gRPC port is exposed via NodePort 30051, so the default address is `localhost:30051`. The client accepts the MCP Server address as a constructor parameter, defaulting to `localhost:30051` when not specified. Handles connection management, gRPC-to-CLI error mapping.

## Sprint Context
- Current entity count: 11
- Scope score: 11/5
- Other entities in sprint: REQ-005, SC-003, SC-004, SC-005, SC-006, SC-007, SC-008, SC-009, SC-00A, SC-00B

## Focus Questions
- Is gRPC tooling (protoc, protoc-gen-go, protoc-gen-go-grpc) available or easy to set up?
- Does the current go.mod support google.golang.org/grpc dependencies?
- Is there an existing proto definition or MCP Server contract to align with?
- Are there any blocking dependencies on the eve-realm-sdk submodule?
- Can the `make proto` target be added without conflicting with existing build targets?
- What test infrastructure is needed for gRPC testing (mock server, bufconn)?
- How should the default address `localhost:30051` be handled (const, config, env var)?
