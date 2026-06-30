---
content_hash: 5255d526fb3994a1805354e22233f7e55b115467f6161afe1e2c31d05b4ba808
created: "2026-06-28"
id: REQ-006
priority: high
related_adrs: []
related_changes: []
related_scenarios:
    - SC-012
    - SC-013
    - SC-014
    - SC-015
    - SC-016
    - SC-017
related_testcases: []
related_userstories: []
source: manual
status: active
tags:
    - cobra
    - cli
    - tools
    - grpc
title: Cobra command for MCP tool listing and invocation
updated: "2026-06-29"
---

# REQ-006: Cobra command for MCP tool listing and invocation

## Description

The `eve-realm` CLI exposes a `tools` subcommand with two operations: `list` and `invoke`.
These Cobra commands are the CLI surface that the master skill (REQ-005) calls to interact
with the MCP Server's tool registry. They use the gRPC client library (REQ-004) internally.

- `eve-realm tools list` — Connects to the MCP Server via gRPC, calls `ListTools`, and
  outputs each tool's name, description, and input schema in a human-readable format
  (also parseable by AI tools).

- `eve-realm tools invoke <name> [--input <json>]` — Connects to the MCP Server via gRPC,
  calls `InvokeTool` with the given tool name and JSON input, and outputs the tool's
  JSON response to stdout.

Both commands read the MCP Server address from the CLI configuration
(`~/.eve-realm/eve-realm.yaml`) or the `EVE_REALM_MCP_ADDR` environment variable,
defaulting to `localhost:30051`.

## Acceptance Criteria

1. `eve-realm tools list` calls `ListTools` via the gRPC client and outputs tool descriptors (name, description, input schema) to stdout.
2. `eve-realm tools invoke <name>` calls `InvokeTool` via the gRPC client with the given tool name and returns the JSON response to stdout.
3. `eve-realm tools invoke <name> --input '{"key":"value"}'` passes the JSON input to the tool.
4. When no `--input` flag is provided, an empty JSON object `{}` is sent as input.
5. When the MCP Server is unreachable, both commands exit with a non-zero code and a clear error message to stderr.
6. When the tool name is not found, `invoke` exits with a non-zero code and prints the error message plus available tools to stderr.
7. The MCP Server address is read from config file, `EVE_REALM_MCP_ADDR` env var, or defaults to `localhost:30051` (in that precedence order).
8. Command code lives in `cmd/tools/`.
