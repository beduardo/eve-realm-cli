# Sprint SP-001: MCP gRPC Connection Client

**Created**: 2026-06-28
**Status**: Specified
**Entities**: 11

---

## Overview

This sprint establishes the foundational gRPC client infrastructure and the master skill
that connects the Eve Realm CLI to the MCP Server's tool registry. By delivering
`internal/mcpclient/` (REQ-004) and `skills/master/` (REQ-005), the CLI gains the ability
to dynamically discover all tools registered in the MCP Server and invoke any of them by
name — making the full MCP ecosystem accessible to Claude Code and other AI tools through
a single marketplace entry point. The nine associated scenarios (SC-003 through SC-00B)
define the exact contract for both the client library and the skill, covering happy paths,
connection failures, and tool-not-found cases, culminating in an end-to-end ping test that
validates the complete chain.

## Entity Inventory

| ID | Type | Title | Partial | Scope Notes |
|----|------|-------|---------|-------------|
| REQ-004 | requirement | gRPC tool client for MCP Server | no | - |
| REQ-005 | requirement | Master skill for MCP tool discovery and invocation | no | - |
| SC-003 | scenario | Client lists tools from MCP Server | no | - |
| SC-004 | scenario | Client invokes tool and returns response | no | - |
| SC-005 | scenario | Client returns connection error when server unreachable | no | - |
| SC-006 | scenario | Client returns typed error for unknown tool | no | - |
| SC-007 | scenario | Discovery mode lists all available tools | no | - |
| SC-008 | scenario | Invocation mode returns tool response | no | - |
| SC-009 | scenario | Skill handles unreachable MCP Server gracefully | no | - |
| SC-00A | scenario | Skill suggests alternatives when tool not found | no | - |
| SC-00B | scenario | End-to-end ping invocation via master skill | no | - |

## Technical Context

All sprint entities are net-new — no existing production code relates to gRPC, mcpclient,
or the master skill. This is the first gRPC dependency introduced into the codebase.

**Entity-to-code mapping highlights:**

- REQ-004 creates `internal/mcpclient/` (client + typed errors), `proto/mcp/v1/mcp.proto`,
  and `gen/proto/mcp/v1/` (generated stubs). Default address is `localhost:30051` (k3d
  NodePort, authoritative per AC-6).
- REQ-005 creates `skills/master/` (skill logic + SKILL.md) and `internal/config/`
  (YAML + env var config resolution). It depends on REQ-004 and must be implemented
  strictly after it.
- SC-003 through SC-006 are covered by `internal/mcpclient/client_test.go` using an
  in-process gRPC server via `bufconn`.
- SC-007 through SC-00A are covered by `skills/master/skill_test.go` using a mock
  `mcpclient` interface.
- SC-00B is covered by `skills/master/skill_test.go` as an integration test using
  `bufconn`.

**Port note:** SC-003, SC-004, and SC-006 reference `localhost:50051` in their step text.
The authoritative default is `localhost:30051` per REQ-004 AC-6. Test fixtures pass
addresses explicitly via `NewClient(addr)`, so test validity is unaffected by this
documentation drift.

**Implementation patterns to follow:**

| Pattern | Source reference | Usage |
|---------|-----------------|-------|
| Table-driven tests with `t.Run` | `eve-cli/main/internal/marketplace/skilldesc_test.go:L132-L183` | All test files |
| Interface-based dependency injection | `eve-cli/main/internal/marketplace/orchestrator.go:L29-L33` | `mockMCPClient` in skill tests; `bufconn` in client tests |
| Config loading via YAML + env var overlay | `eve-cli/main/internal/config/config.go:L144-L234` | `internal/config/config.go` |
| SKILL.md with YAML frontmatter | `eve-cli/main/internal/marketplace/skilldesc.go` | `skills/master/SKILL.md` |
| Typed error structs with `errors.As` support | `eve-cli/main/internal/marketplace/orchestrator.go` | `internal/mcpclient/errors.go` |
| YAML struct tags with `gopkg.in/yaml.v3` | `eve-cli/main/internal/config/config.go:L129-L141` | `internal/config/config.go` |
| Test function naming `TestFunctionName_Scenario` | `eve-realm-cli/main/cmd/main_test.go:L5` | All test functions |

**Critical integration points:**

1. `go.mod` — currently bare; must add `google.golang.org/grpc`, `google.golang.org/protobuf`,
   and `gopkg.in/yaml.v3`.
2. `Makefile` — must add a `proto` target (via `buf generate` or `protoc`) and extend
   `.PHONY`.
3. Dependency chain: `skills/master/skill.go` -> `internal/mcpclient.Client` ->
   `gen/proto/mcp/v1.MCPServiceClient` -> gRPC.
4. Config resolution: env var `EVE_REALM_MCP_ADDR` > YAML `mcp_server_addr` > empty
   string to `NewClient("")` which substitutes `localhost:30051`.

**Files to create:**

| File | Purpose |
|------|---------|
| `proto/mcp/v1/mcp.proto` | Proto definition for MCPService (ListTools, InvokeTool RPCs) |
| `gen/proto/mcp/v1/mcp.pb.go` | Generated protobuf message types |
| `gen/proto/mcp/v1/mcp_grpc.pb.go` | Generated gRPC client stub |
| `internal/mcpclient/client.go` | `Client` type: `NewClient`, `ListTools`, `InvokeTool`; defaults to `localhost:30051` |
| `internal/mcpclient/errors.go` | `ConnectionError` and `ToolNotFoundError` typed errors |
| `internal/mcpclient/client_test.go` | Unit tests using bufconn in-process gRPC server |
| `internal/config/config.go` | `HostConfig` with `MCPServerAddr`; `LoadHostConfig` + `Resolve` |
| `internal/config/config_test.go` | Config loading, env var overlay, missing-file tests |
| `skills/master/skill.go` | Master skill — discovery and invocation modes |
| `skills/master/skill_test.go` | Unit tests with mock mcpclient + SC-00B bufconn integration test |
| `skills/master/SKILL.md` | Skill metadata with YAML frontmatter |
| `buf.yaml` | Buf configuration for proto generation |

**Files to modify:**

| File | Modification |
|------|--------------|
| `go.mod` | Add `google.golang.org/grpc`, `google.golang.org/protobuf`, `gopkg.in/yaml.v3` |
| `Makefile` | Add `proto` target; extend `.PHONY` |

---

## Implementation Sections

### REQ-004: gRPC tool client for MCP Server

**Entity**: `.software/entities/requirements/REQ-004.md`
**Type**: requirement
**Priority**: high

**Codebase Mapping**:

Files to create:
- `proto/mcp/v1/mcp.proto` — MCPService proto definition (ListTools + InvokeTool RPCs); maps to AC-7
- `gen/proto/mcp/v1/mcp.pb.go` — generated protobuf message types; maps to AC-7
- `gen/proto/mcp/v1/mcp_grpc.pb.go` — generated gRPC client stub; maps to AC-7
- `internal/mcpclient/client.go` — `Client` type with `NewClient(addr string)`, `ListTools`, `InvokeTool`; `DefaultMCPAddr = "localhost:30051"`; maps to AC-1 through AC-6
- `internal/mcpclient/errors.go` — `ConnectionError` and `ToolNotFoundError` concrete structs implementing `error`; maps to AC-4, AC-5
- `internal/mcpclient/client_test.go` — bufconn in-process gRPC server tests; maps to SC-003, SC-004, SC-005, SC-006

Files to modify:
- `go.mod` — add `google.golang.org/grpc`, `google.golang.org/protobuf`
- `Makefile` — add `proto` target, extend `.PHONY`

**Acceptance Criteria**:

- **AC-1**: Given the `internal/mcpclient` package, when consumed by any caller, then it exposes a `Client` type with `ListTools(ctx context.Context) ([]Tool, error)` and `InvokeTool(ctx context.Context, name, input string) (string, error)` methods.
- **AC-2**: Given a running MCP Server with registered tools, when `ListTools(ctx)` is called, then it invokes the server's `ListTools` gRPC RPC and returns a `[]Tool` slice where each entry has populated `name`, `description`, and `input_schema` fields.
- **AC-3**: Given a running MCP Server with a registered `ping` tool, when `InvokeTool(ctx, "ping", "{}")` is called, then it invokes the server's `InvokeTool` gRPC RPC with the given name and JSON input, and returns the JSON output string.
- **AC-4**: Given an MCP Server that is unreachable at the configured address, when any client method is called, then it returns a clear error indicating connection failure — not a raw gRPC status code or error message.
- **AC-5**: Given an MCP Server running and receiving an `InvokeTool` call for a nonexistent tool name, when the server returns gRPC `NOT_FOUND` status, then the client returns a typed error (e.g., `ToolNotFoundError`) that is distinguishable from connection errors via `errors.As`.
- **AC-6**: Given `NewClient` is called with an empty string as address, when the client connects, then it defaults to `localhost:30051`; when a non-empty address is provided, it uses that address.
- **AC-7**: Given the `proto/mcp/v1/mcp.proto` file, when it matches the MCP Server's definition, then generated Go code in `gen/proto/mcp/v1/` provides the correct client stubs.
- **AC-8**: Given the `Makefile` `proto` target, when `make proto` is executed, then it generates Go client stubs in `gen/proto/mcp/v1/` from the proto definition without error.

**Implementation Notes**:

Complexity estimate: **M** (~250–400 LOC production, ~200–300 LOC tests).

Key risks and mitigations:
- **gRPC tooling absent (confirmed)**: Install `buf` + `protoc-gen-go` + `protoc-gen-go-grpc` plugins before writing any proto. Alternatively, commit the generated files to avoid a toolchain prerequisite.
- **No proto definition exists (confirmed)**: Author `proto/mcp/v1/mcp.proto` as the first task of REQ-004. Define MCPService with two RPCs: `ListTools` (empty request, returns list of tool descriptors) and `InvokeTool` (tool name + JSON input, returns JSON output string).
- **go.mod is bare (confirmed)**: Run `go get google.golang.org/grpc google.golang.org/protobuf` as the first implementation step.
- **make proto target missing (confirmed)**: Add `.PHONY: proto` to Makefile alongside the `buf generate` or `protoc` invocation.
- **Proto contract divergence from MCP Server (low)**: Treat the local proto as canonical for SP-001. SDK submodule proto sharing is deferred.
- **Default address**: Implement as `const DefaultMCPAddr = "localhost:30051"` in `client.go`; `NewClient` uses it when `addr == ""`.

**Test Expectations**:

- Must test: `ListTools` returns a non-empty `[]Tool` slice when server has tools registered — each entry has non-empty `name`, `description`, and `input_schema` (SC-003, AC-2). Uses bufconn in-process server.
- Must test: `InvokeTool` with a valid tool name and `"{}"` input returns a non-empty JSON string and no error (SC-004, AC-3). Uses bufconn in-process server.
- Must test: `ListTools` and `InvokeTool` return a `ConnectionError` (detectable via `errors.As`) when the server address is unreachable (SC-005, AC-4). Uses `localhost:59999` or equivalent unused port.
- Must test: `InvokeTool` with an unknown tool name returns a `ToolNotFoundError` (detectable via `errors.As`) and not a `ConnectionError` (SC-006, AC-5). Uses bufconn server that responds with gRPC `NOT_FOUND`.
- Must test: `NewClient("")` connects to `localhost:30051`; `NewClient("localhost:9999")` connects to the given address (AC-6). Verified via the default constant, not by actually connecting.
- Must test: `ToolNotFoundError` message includes the tool name that was not found (SC-006 expected result).
- Must NOT rely on: external network — all tests use bufconn or a known-unused port. No testify or external assertion libraries. No global test state.

---

### REQ-005: Master skill for MCP tool discovery and invocation

**Entity**: `.software/entities/requirements/REQ-005.md`
**Type**: requirement
**Priority**: high

**Codebase Mapping**:

Files to create:
- `skills/master/skill.go` — discovery and invocation modes; reads config via `internal/config`; delegates to `internal/mcpclient.Client`; maps to AC-1 through AC-8
- `skills/master/skill_test.go` — unit tests with mock mcpclient; SC-00B integration test with bufconn; maps to SC-007, SC-008, SC-009, SC-00A, SC-00B
- `skills/master/SKILL.md` — YAML frontmatter with `description`, `argument-hint`, `disable-model-invocation`; maps to AC-1
- `internal/config/config.go` — `HostConfig` struct with `MCPServerAddr string \`yaml:"mcp_server_addr"\``; `LoadHostConfig(path string)` returning zero-value on missing file; `Resolve(envKey, yamlValue string) string`; maps to AC-7
- `internal/config/config_test.go` — YAML round-trip, env var overlay, missing-file (non-error) tests; maps to AC-7

Files to modify: none beyond what REQ-004 already modifies.

**Acceptance Criteria**:

- **AC-1**: Given the CLI marketplace, when the master skill is registered, then it appears under the name `eve-realm` (or configured name) and its `SKILL.md` contains valid YAML frontmatter with required fields.
- **AC-2**: Given the master skill invoked without arguments (discovery mode), when `ListTools` is called on the gRPC client, then the skill returns formatted text listing each available tool with its name, description, and input schema in a format suitable for AI consumption.
- **AC-3**: Given the master skill invoked with a specific tool name and input (invocation mode), when `InvokeTool` is called on the gRPC client, then the skill returns the tool's JSON response directly with no wrapper formatting.
- **AC-4**: Given the MCP Server is unreachable when the skill is invoked (either mode), when the client returns a `ConnectionError`, then the skill returns a user-friendly message (e.g., "MCP Server is not available at `<address>`. Check that the server is running.") with no stack traces or raw gRPC error codes.
- **AC-5**: Given the master skill is invoked with a tool name that does not exist, when the client returns a `ToolNotFoundError`, then the skill returns a clear message stating the tool was not found and lists available tools as alternatives.
- **AC-6**: Given the skill implementation, when code is examined, then all skill logic resides in `skills/master/`.
- **AC-7**: Given a `~/.eve-realm/eve-realm.yaml` file with `mcp_server_addr` set, or the `EVE_REALM_MCP_ADDR` environment variable set, when the skill initializes the gRPC client, then it uses that address; when neither is set, `NewClient("")` defaults to `localhost:30051`.
- **AC-8**: Given the master skill invoked with tool name `"ping"`, when the MCP Server has the ping tool registered, then the skill returns `{"message": "pong", "timestamp": "<RFC 3339>"}` with a recent timestamp and valid JSON.

**Implementation Notes**:

Complexity estimate: **M (skill alone) / L (with config package and registration)** (~300–450 LOC across 4–6 files).

Key risks and mitigations:
- **Skill registration undefined (high likelihood, high impact)**: For SP-001, use a manual `extraKnownMarketplaces` entry in `~/.claude/settings.json`. Automated registration is deferred to a later sprint.
- **Sequential dependency on REQ-004 (hard)**: REQ-004 must be fully implemented (code compiles, tests pass) before REQ-005 implementation begins. The `internal/mcpclient` package is the direct import dependency.
- **`internal/config/` does not exist**: Create during REQ-005 implementation, referencing `eve-cli/main/internal/config/config.go` for YAML loading and env-var overlay patterns.
- **SC-00B E2E test**: Implement as an integration test within `skills/master/skill_test.go` using a bufconn in-process gRPC server that handles the ping tool. This avoids any external network dependency.

**Test Expectations**:

- Must test: Discovery mode (no args) calls `ListTools` and returns formatted text containing tool names, descriptions, and input schemas — all registered tools present (SC-007, AC-2). Uses `mockMCPClient` interface.
- Must test: Invocation mode (tool name + input) calls `InvokeTool` and returns the raw JSON response string with no additional formatting (SC-008, AC-3). Uses `mockMCPClient`.
- Must test: When `mockMCPClient.ListTools` or `mockMCPClient.InvokeTool` returns a `ConnectionError`, the skill returns a user-friendly message containing the server address, with no gRPC error codes or stack traces (SC-009, AC-4). Uses `mockMCPClient`.
- Must test: When `mockMCPClient.InvokeTool` returns a `ToolNotFoundError` for tool `"nonexistent"`, the skill response includes the tool name and lists available alternatives (SC-00A, AC-5). Uses `mockMCPClient`.
- Must test: End-to-end — invoking the skill with tool `"ping"` against a bufconn in-process server returns `{"message": "pong", "timestamp": "..."}` with valid JSON and a recent timestamp (SC-00B, AC-8).
- Must test: Config loading reads `mcp_server_addr` from YAML, env var `EVE_REALM_MCP_ADDR` overrides YAML, and missing config file is a non-error condition returning zero-value (AC-7). Uses `t.TempDir()` for YAML round-trip tests.
- Must NOT rely on: external network — all tests use `mockMCPClient` or bufconn. No testify or external assertion libraries. No global test state.

---

### SC-003: Client lists tools from MCP Server

**Entity**: `.software/entities/scenarios/SC-003.md`
**Type**: scenario
**Priority**: (inherited from REQ-004)

**Codebase Mapping**:

Covered by `internal/mcpclient/client_test.go` using a bufconn in-process gRPC server that registers at least one tool (e.g., `ping`). Test passes the server address explicitly — not the `localhost:30051` default.

**Acceptance Criteria**:

- **Given** an MCP Server running at a known address with at least one tool registered (e.g., `ping`), **when** `NewClient(addr)` creates a client and `ListTools(ctx)` is called, **then** the method returns a `[]Tool` slice with at least one entry, each having non-empty `name`, `description`, and `input_schema` fields, and no error is returned.

**Implementation Notes**:

Feasibility not assessed separately for scenarios. This scenario is covered directly by REQ-004 AC-2. Use bufconn for the in-process server; do not rely on external network. Note: the scenario step text uses `localhost:50051` — tests must pass the address explicitly to avoid confusion with the `localhost:30051` default.

**Test Expectations**:

- Must test: `ListTools` returns `[]Tool` with at least one entry when the in-process server has a tool registered.
- Must test: Each returned `Tool` has non-empty `Name`, `Description`, and `InputSchema` fields.
- Must test: No error is returned on a successful `ListTools` call.
- Must NOT rely on: external network; external assertion libraries.

---

### SC-004: Client invokes tool and returns response

**Entity**: `.software/entities/scenarios/SC-004.md`
**Type**: scenario
**Priority**: (inherited from REQ-004)

**Codebase Mapping**:

Covered by `internal/mcpclient/client_test.go` using bufconn in-process server with a `ping` tool handler. Test passes the address explicitly.

**Acceptance Criteria**:

- **Given** an MCP Server running at a known address with a `ping` tool registered, **when** `NewClient(addr)` creates a client and `InvokeTool(ctx, "ping", "{}")` is called, **then** the method returns a non-empty JSON string (e.g., `{"message": "pong", ...}`) and no error.

**Implementation Notes**:

Feasibility not assessed separately for scenarios. Covered by REQ-004 AC-3. Use the same bufconn server test fixture as SC-003, extended with an `InvokeTool` RPC handler for the ping tool.

**Test Expectations**:

- Must test: `InvokeTool("ping", "{}")` returns a non-empty string with no error against an in-process server.
- Must test: The returned string is valid JSON.
- Must NOT rely on: external network; external assertion libraries.

---

### SC-005: Client returns connection error when server unreachable

**Entity**: `.software/entities/scenarios/SC-005.md`
**Type**: scenario
**Priority**: (inherited from REQ-004)

**Codebase Mapping**:

Covered by `internal/mcpclient/client_test.go`. Test uses `localhost:59999` (a port with no listener). Both `ListTools` and `InvokeTool` branches must be tested.

**Acceptance Criteria**:

- **Given** no MCP Server running at the target address (e.g., `localhost:59999`), **when** `NewClient("localhost:59999")` creates a client and `ListTools(ctx)` or `InvokeTool(ctx, "ping", "{}")` is called, **then** an error is returned that: (a) is clearly a connection error (not a raw gRPC status), (b) mentions the unreachable address, and (c) is distinguishable from `ToolNotFoundError` via `errors.As`.

**Implementation Notes**:

Feasibility not assessed separately for scenarios. Covered by REQ-004 AC-4. The `ConnectionError` struct must include the address field for message composition. gRPC connection failures surface as `codes.Unavailable` or dial errors — the client wraps these into `ConnectionError` before returning.

**Test Expectations**:

- Must test: `ListTools` on an unreachable address returns an error where `errors.As(err, &ConnectionError{})` is true.
- Must test: `InvokeTool` on an unreachable address returns the same typed error.
- Must test: The `ConnectionError` message includes the address string `"localhost:59999"`.
- Must test: `errors.As(err, &ToolNotFoundError{})` returns false for a `ConnectionError`.
- Must NOT rely on: external network (use a known-unused local port, not a remote host).

---

### SC-006: Client returns typed error for unknown tool

**Entity**: `.software/entities/scenarios/SC-006.md`
**Type**: scenario
**Priority**: (inherited from REQ-004)

**Codebase Mapping**:

Covered by `internal/mcpclient/client_test.go` using a bufconn in-process server that returns gRPC `NOT_FOUND` status for any tool invocation with an unknown name.

**Acceptance Criteria**:

- **Given** an MCP Server running at a known address with at least one tool registered, **when** `NewClient(addr)` creates a client and `InvokeTool(ctx, "nonexistent", "{}")` is called, **then** a typed `ToolNotFoundError` is returned that: (a) is detectable via `errors.As`, (b) is distinguishable from `ConnectionError`, and (c) includes the tool name `"nonexistent"` in its message.

**Implementation Notes**:

Feasibility not assessed separately for scenarios. Covered by REQ-004 AC-5. The bufconn server handler should respond with `status.Error(codes.NotFound, "tool not found: nonexistent")` when the tool name is unknown. The client maps this gRPC status code to `ToolNotFoundError`.

**Test Expectations**:

- Must test: `InvokeTool("nonexistent", "{}")` returns an error where `errors.As(err, &ToolNotFoundError{})` is true.
- Must test: `errors.As(err, &ConnectionError{})` returns false for a `ToolNotFoundError`.
- Must test: The `ToolNotFoundError` message includes the tool name `"nonexistent"`.
- Must NOT rely on: external network; external assertion libraries.

---

### SC-007: Discovery mode lists all available tools

**Entity**: `.software/entities/scenarios/SC-007.md`
**Type**: scenario
**Priority**: (inherited from REQ-005)

**Codebase Mapping**:

Covered by `skills/master/skill_test.go` using a `mockMCPClient` interface that returns a predefined list of tools (e.g., `ping`, `echo`).

**Acceptance Criteria**:

- **Given** the master skill and a mock MCP client returning multiple tools (`ping`, `echo`), **when** the skill is invoked without arguments (discovery mode), **then** the returned text contains each tool's name, description, and input schema, and all registered tools appear in the output in a format suitable for AI consumption.

**Implementation Notes**:

Feasibility not assessed separately for scenarios. Covered by REQ-005 AC-2. The mock interface returns a fixed `[]Tool` slice. Assert that the output string contains each tool name and description.

**Test Expectations**:

- Must test: Discovery mode output contains the name and description of each tool returned by `mockMCPClient.ListTools`.
- Must test: Discovery mode output contains the input schema of each tool.
- Must test: All tools from the mock are present in the output (none silently omitted).
- Must NOT rely on: real gRPC connection; external assertion libraries.

---

### SC-008: Invocation mode returns tool response

**Entity**: `.software/entities/scenarios/SC-008.md`
**Type**: scenario
**Priority**: (inherited from REQ-005)

**Codebase Mapping**:

Covered by `skills/master/skill_test.go` using a `mockMCPClient` that returns a predefined JSON string for the `ping` tool.

**Acceptance Criteria**:

- **Given** the master skill and a mock MCP client configured for the `ping` tool, **when** the skill is invoked with tool name `"ping"` and empty input, **then** it returns the tool's JSON response string directly (e.g., `{"message": "pong", ...}`) with no wrapper or additional formatting.

**Implementation Notes**:

Feasibility not assessed separately for scenarios. Covered by REQ-005 AC-3. The mock returns a fixed JSON string; the test asserts the skill output equals that string exactly, with no extra wrapping.

**Test Expectations**:

- Must test: Invocation mode output equals the raw JSON string returned by `mockMCPClient.InvokeTool`, with no modification.
- Must test: No additional text, headers, or formatting is added around the tool response.
- Must NOT rely on: real gRPC connection; external assertion libraries.

---

### SC-009: Skill handles unreachable MCP Server gracefully

**Entity**: `.software/entities/scenarios/SC-009.md`
**Type**: scenario
**Priority**: (inherited from REQ-005)

**Codebase Mapping**:

Covered by `skills/master/skill_test.go` using a `mockMCPClient` that returns a `ConnectionError` from `ListTools` and `InvokeTool`.

**Acceptance Criteria**:

- **Given** the master skill and a mock MCP client that returns a `ConnectionError` for any call, **when** the skill is invoked in either mode (discovery or invocation), **then** it returns a user-friendly error message (e.g., "MCP Server is not available at `<address>`. Check that the server is running.") that includes the configured server address, with no stack traces or raw gRPC error codes.

**Implementation Notes**:

Feasibility not assessed separately for scenarios. Covered by REQ-005 AC-4. The skill must detect `ConnectionError` via `errors.As` and format a user-readable message from the address stored in the error.

**Test Expectations**:

- Must test: When the mock returns a `ConnectionError`, the skill output is a user-friendly string (not a Go error string or gRPC status message).
- Must test: The output includes the server address from the `ConnectionError`.
- Must test: No gRPC status codes, error type names, or stack traces appear in the output.
- Must NOT rely on: real gRPC connection; external assertion libraries.

---

### SC-00A: Skill suggests alternatives when tool not found

**Entity**: `.software/entities/scenarios/SC-00A.md`
**Type**: scenario
**Priority**: (inherited from REQ-005)

**Codebase Mapping**:

Covered by `skills/master/skill_test.go`. The mock is configured to: return `ToolNotFoundError` for `InvokeTool("nonexistent", ...)`, and return a known list of tools (`ping`, `echo`) from `ListTools` — which the skill calls as a follow-up to populate the alternatives list.

**Acceptance Criteria**:

- **Given** the master skill with a mock MCP client that has `ping` and `echo` registered, **when** the skill is invoked with tool name `"nonexistent"`, **then** the skill returns a message that: (a) states the tool `"nonexistent"` was not found, and (b) lists the available tools (`ping`, `echo`) as alternatives.

**Implementation Notes**:

Feasibility not assessed separately for scenarios. Covered by REQ-005 AC-5. The skill detects `ToolNotFoundError` via `errors.As`, then calls `ListTools` to retrieve available tools and includes them in the error message.

**Test Expectations**:

- Must test: The output message includes the string `"nonexistent"` (the tool name that was not found).
- Must test: The output message lists available tool names (`ping`, `echo`) as alternatives.
- Must test: The behavior is consistent in both discovery mode and invocation mode when the tool is not found.
- Must NOT rely on: real gRPC connection; external assertion libraries.

---

### SC-00B: End-to-end ping invocation via master skill

**Entity**: `.software/entities/scenarios/SC-00B.md`
**Type**: scenario
**Priority**: (inherited from REQ-005)

**Codebase Mapping**:

Covered by an integration test in `skills/master/skill_test.go` using a bufconn in-process gRPC server that implements the `ping` tool handler. The CLI config is set to the bufconn address for this test.

**Acceptance Criteria**:

- **Given** the master skill and a bufconn in-process MCP Server with the `ping` tool registered, and CLI configured with the bufconn address, **when** the skill is invoked with tool name `"ping"`, **then** it returns `{"message": "pong", "timestamp": "<RFC 3339>"}` where the timestamp is recent (within seconds of invocation) and the response is valid JSON.

**Implementation Notes**:

Feasibility not assessed separately for scenarios. Covered by REQ-005 AC-8. This is the one integration test in the sprint — it exercises the complete chain from skill invocation through config resolution, gRPC client, and the in-process server, without mocks. Use `bufconn` to avoid OS sockets. Validate the timestamp is parseable as RFC 3339 and not a zero value.

**Test Expectations**:

- Must test: Full chain — skill invocation -> config resolution -> `NewClient` -> `InvokeTool` -> in-process bufconn server -> response returned.
- Must test: Returned string is valid JSON parseable into a struct with `message` and `timestamp` fields.
- Must test: `message` field equals `"pong"`.
- Must test: `timestamp` field parses as RFC 3339 and is within a reasonable window of `time.Now()`.
- Must NOT rely on: external network; mocked gRPC client (this is the integration test that exercises the real client code).

---

## Documentation Tasks

### RELEASES.md Entry

**Required**: Always

Add an entry to RELEASES.md documenting:
- Sprint ID and title: SP-001 — MCP gRPC Connection Client
- Summary of changes delivered: gRPC client library (`internal/mcpclient/`) connecting to the MCP Server's tool registry, master skill (`skills/master/`) enabling tool discovery and invocation from Claude Code, config resolution via YAML and env var overlay (`internal/config/`), proto definition and generated stubs (`proto/mcp/v1/`, `gen/proto/mcp/v1/`), and Makefile `proto` target.
- Entity IDs included: REQ-004, REQ-005, SC-003, SC-004, SC-005, SC-006, SC-007, SC-008, SC-009, SC-00A, SC-00B
- Date of completion: to be filled at completion time

This entry should be appended to the existing RELEASES.md file. Do not read or modify existing entries.

### README.md Update

**Required**: User-facing changes detected

Update README.md to reflect:
- New `skills/master/` skill registered as `eve-realm` in the marketplace, enabling Claude Code to list and invoke any MCP Server tool via a single skill entry.
- MCP Server address configuration: set `mcp_server_addr` in `~/.eve-realm/eve-realm.yaml` or export `EVE_REALM_MCP_ADDR` to point the skill at a non-default server. Default address is `localhost:30051` (k3d NodePort).
- New `make proto` command that regenerates Go client stubs from `proto/mcp/v1/mcp.proto`.
- Skill registration for SP-001 uses a manual `extraKnownMarketplaces` entry in `~/.claude/settings.json` — instructions for this manual step.

---

## Pinned Entity Compliance

| Entity | Directive | How Addressed |
|--------|-----------|---------------|
| REQ-003: Cross-cutting requirements catalog for lazy-loaded sprint policy injection | No spec-phase directives found in entity body; functions as a registry loader with mandatory loading rule for matched triggers. Trigger "Implementing or modifying Go code" matches this sprint — REQ-001 (TDD) was loaded and applied. | REQ-001 (`testing_strategy: tdd`) trigger matched. Test Expectations subsections generated for all applicable entities (REQ-004, REQ-005, SC-003 through SC-00B). REQ-002 trigger (sprint completion/release) does not match the spec phase and is deferred to implementation completion. |

---

## Out of Scope

- SDK submodule initialization and proto sharing between CLI and MCP Server via SDK.
- Automated skill registration (marketplace auto-discovery). SP-001 uses manual `extraKnownMarketplaces` entry.
- Generic agent skill (`skills/agent/`) — deferred to a future sprint.
- Auth, marketplace, and settings commands — no changes in this sprint.
- `cmd/main.go` modifications — the 20-line scaffold is unchanged.
- Updating SC-003, SC-004, and SC-006 scenario step text to reflect the correct port (`30051` instead of `50051`) — documentation drift only, deferred.
- Proto contract validation against a live MCP Server — the local proto is canonical for SP-001.

## Prerequisites

- `buf` CLI installed (or `protoc` + `protoc-gen-go` + `protoc-gen-go-grpc` plugins) to run `make proto`. Alternatively, generated files can be committed to the repository to skip the toolchain requirement.
- `go get` available to add `google.golang.org/grpc`, `google.golang.org/protobuf`, and `gopkg.in/yaml.v3` to `go.mod`.
- REQ-004 must be fully implemented (all tests passing, package compiles) before REQ-005 implementation begins. This is a hard sequential constraint.
- Manual `extraKnownMarketplaces` entry in `~/.claude/settings.json` configured before end-to-end skill verification (SC-00B manual test).
