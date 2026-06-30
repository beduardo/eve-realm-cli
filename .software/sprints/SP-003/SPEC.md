# Sprint SP-003: Cobra command for MCP tool listing and invocation

**Created**: 2026-06-30
**Status**: Specified
**Entities**: 7

---

## Overview

This sprint delivers the `eve-realm tools` Cobra command group, exposing two operations — `list` and `invoke` — that allow users and AI agents to interact directly with the MCP Server's tool registry from the CLI. Built on the gRPC client library established in SP-001, this sprint also introduces Cobra as the root command framework for the binary, replacing the current bare argument parser in `cmd/main.go`. Together, these changes establish the foundational CLI structure that all future subcommands will follow.

## Entity Inventory

| ID | Type | Title | Partial | Scope Notes |
|----|------|-------|---------|-------------|
| REQ-006 | requirement | Cobra command for MCP tool listing and invocation | no | - |
| SC-012 | scenario | Tools list outputs tool descriptors to stdout | no | - |
| SC-013 | scenario | Tools invoke sends default empty input and returns JSON to stdout | no | - |
| SC-014 | scenario | Tools invoke passes --input flag value to the tool | no | - |
| SC-015 | scenario | Tools commands exit non-zero with error when server unreachable | no | - |
| SC-016 | scenario | Tools invoke exits non-zero and lists alternatives when tool not found | no | - |
| SC-017 | scenario | MCP Server address resolved from env var, config file, or default | no | - |

## Technical Context

The codebase analysis identified all new code as belonging to a new `cmd/tools/` package, with no modifications to existing packages other than `cmd/main.go` (Cobra root wiring) and `go.mod` (adding `github.com/spf13/cobra`).

**Entity-to-code mapping highlights:**

- REQ-006, SC-012, SC-013, SC-014, SC-015, SC-016: All new — implemented in `cmd/tools/tools.go`, `cmd/tools/list.go`, `cmd/tools/invoke.go`, and their respective test files.
- SC-017: Consumes `internal/config/config.go` (`LoadHostConfig` + `Resolve`) already in place from SP-001. No changes to the config package required.

**Implementation patterns to follow:**

1. **Constructor pattern**: `NewToolsCmd(configPath string) *cobra.Command` factory, with `newListCmd()` and `newInvokeCmd()` helpers registered via `cmd.AddCommand(...)`.
2. **RunE with routed I/O**: Commands use `RunE`; output goes to `cmd.OutOrStdout()` (success) and `cmd.ErrOrStderr()` (errors). Never `fmt.Println` or `os.Stderr` directly.
3. **Local MCPClient interface**: `cmd/tools/` defines its own local `MCPClient` interface (`ListTools` + `InvokeTool` methods). Go structural typing satisfies `internal/mcpclient.Client` without coupling `cmd/` to `skills/`.
4. **Mock via configurable function fields**: Tests use a `mockMCPClient` struct with function fields; injected via the command factory parameter.
5. **Cobra command tests**: `cmd.SetArgs(args)` + `cmd.SetOut(&stdout)` + `cmd.SetErr(&stderr)` + `cmd.Execute()`; assert on returned error, `stdout.String()`, `stderr.String()`.
6. **Typed error detection**: `errors.As(err, &mcpclient.ConnectionError{})` and `errors.As(err, &mcpclient.ToolNotFoundError{})` drive non-zero exit paths.

**Critical integration points:**

- Address resolution: `config.LoadHostConfig(configPath)` → `config.Resolve("EVE_REALM_MCP_ADDR", cfg.MCPServerAddr)` → fallback to `mcpclient.DefaultMCPAddr` when result is empty.
- `cmd/main.go` must be replaced with a Cobra root command. The existing 21-line bare arg parser is removed entirely.
- SC-016 secondary `ListTools` call mirrors the pattern in `skills/master/skill.go:L77-L89`.
- Set `SilenceErrors = true` on the root command to prevent Cobra from double-printing errors (SC-015, SC-016).

**Files to create:**

| File | Purpose |
|------|---------|
| `cmd/tools/tools.go` | `NewToolsCmd()` factory registering `list` and `invoke` subcommands |
| `cmd/tools/list.go` | `newListCmd()` — calls `ListTools`, formats tool descriptors to stdout |
| `cmd/tools/invoke.go` | `newInvokeCmd()` — positional arg for tool name, `--input` flag, calls `InvokeTool` |
| `cmd/tools/tools_test.go` | Shared test helpers: `runToolsCmd(t, args...)`, `mockMCPClient` |
| `cmd/tools/list_test.go` | Tests for `tools list` covering SC-012, SC-015, SC-017 |
| `cmd/tools/invoke_test.go` | Tests for `tools invoke` covering SC-013, SC-014, SC-015, SC-016 |

**Files to modify:**

| File | Modification |
|------|--------------|
| `cmd/main.go` | Replace bare arg parser with Cobra root command; register `tools` subcommand |
| `go.mod` | Add `github.com/spf13/cobra` dependency |

## Implementation Sections

### REQ-006: Cobra command for MCP tool listing and invocation

**Entity**: `.software/requirements/REQ-006-cobra-command-for-mcp-tool-listing-and-invocation.md`
**Type**: requirement
**Priority**: high

**Codebase Mapping**:

New package `cmd/tools/` — all files created from scratch:
- `cmd/tools/tools.go` — `NewToolsCmd(configPath string) *cobra.Command` (based on eve-cli `cmd/eve5/marketplace/marketplace.go`)
- `cmd/tools/list.go` — `newListCmd(client MCPClient) *cobra.Command` (based on eve-cli `cmd/eve5/marketplace/list.go`)
- `cmd/tools/invoke.go` — `newInvokeCmd(client MCPClient) *cobra.Command` (based on eve-cli `cmd/eve5/marketplace/register.go`)
- `cmd/tools/tools_test.go`, `cmd/tools/list_test.go`, `cmd/tools/invoke_test.go` — test suite

Modified files:
- `cmd/main.go`: Replace with Cobra root command + `tools` subcommand wiring
- `go.mod`: Add `github.com/spf13/cobra`

**Acceptance Criteria**:

- **AC-1**: Given the MCP Server is running, when `eve-realm tools list` is executed, then it calls `ListTools` via the gRPC client and outputs each tool's name, description, and input schema to stdout.
- **AC-2**: Given the MCP Server is running and a tool named `<name>` exists, when `eve-realm tools invoke <name>` is executed, then it calls `InvokeTool` via the gRPC client with the given tool name and writes the JSON response to stdout.
- **AC-3**: Given the MCP Server is running, when `eve-realm tools invoke <name> --input '{"key":"value"}'` is executed, then the JSON string `{"key":"value"}` is passed verbatim to `InvokeTool`.
- **AC-4**: Given the MCP Server is running and no `--input` flag is provided, when `eve-realm tools invoke <name>` is executed, then `{}` is sent as the input to `InvokeTool`.
- **AC-5**: Given the MCP Server is unreachable, when either `tools list` or `tools invoke` is executed, then the command exits with a non-zero status code and writes a clear error message to stderr.
- **AC-6**: Given the MCP Server is running but the named tool does not exist, when `eve-realm tools invoke <name>` is executed, then the command exits with a non-zero status code, writes the not-found error to stderr, and lists the available tools on stderr.
- **AC-7**: Given MCP Server address configuration sources exist, when any `tools` subcommand executes, then the address is resolved from `EVE_REALM_MCP_ADDR` env var first, then from the `mcp_server_addr` field in `~/.eve-realm/eve-realm.yaml`, then defaulting to `localhost:30051`.
- **AC-8**: Given the implementation, when the package structure is inspected, then all command code lives in `cmd/tools/`.

**Implementation Notes**:

The feasibility assessment is **PROCEED-WITH-CAVEATS**. All core infrastructure from SP-001 is ready (`internal/mcpclient`, `internal/config`). Two blockers must be resolved first:

1. **Critical**: `github.com/spf13/cobra` is absent from `go.mod`. Run `go get github.com/spf13/cobra` as the first implementation step.
2. **Major**: `cmd/main.go` is a 21-line bare argument parser with no Cobra root command. Replace it entirely with a Cobra root command following the eve-cli pattern before adding the `tools` subcommand.

Complexity estimate: **M** (200-350 LOC production, 200-300 LOC test, 6 new files + 2 modified).

Additional implementation guidance:
- Config file path must be injectable into `NewToolsCmd(configPath string)` — not computed from `os.UserHomeDir()` at call time — so tests can use `t.TempDir()` (required for SC-017).
- Set `SilenceErrors = true` on the root command to prevent Cobra from double-printing errors for SC-015 and SC-016.
- For SC-016, after detecting `ToolNotFoundError`, make a secondary `ListTools` call to build the alternatives list on stderr, mirroring `skills/master/skill.go:L77-L89`.
- For SC-017, when `config.Resolve(...)` returns an empty string (both env var and YAML value absent), fall back to `mcpclient.DefaultMCPAddr` (`localhost:30051`).

**Test Expectations:**

- Must test: `tools list` with a mock returning multiple tool descriptors verifies each tool's name, description, and input schema appears in stdout (table-driven, one case per field).
- Must test: `tools invoke <name>` with no `--input` flag verifies `InvokeTool` is called with `{}` as the input argument.
- Must test: `tools invoke <name> --input '{"key":"value"}'` verifies `InvokeTool` is called with exactly `{"key":"value"}` (not modified or re-encoded).
- Must test: `tools list` and `tools invoke` when mock returns `ConnectionError` verify command exits non-zero and stderr contains an error message; stdout is empty.
- Must test: `tools invoke <name>` when mock returns `ToolNotFoundError` verifies exit non-zero, stderr contains the not-found message, and stderr contains the names of alternative tools returned by the secondary `ListTools` call.
- Must test: Address resolution for SC-017 — three cases in a table: (1) env var set overrides YAML value; (2) env var absent, YAML value used; (3) both absent, `localhost:30051` is used. Each case sets up a temp config file (via `t.TempDir()`) and optionally sets the env var.
- Must test: `tools list` output with zero tools registered returns exit zero with empty (or graceful) output — boundary condition.
- Must NOT rely on: Real gRPC connections or real filesystem paths for any test. All MCP client calls must go through the injected `mockMCPClient`. Config file reads must use `t.TempDir()`-scoped paths.

---

### SC-012: Tools list outputs tool descriptors to stdout

**Entity**: `.software/scenarios/SC-012-tools-list-outputs-tool-descriptors-to-stdout.md`
**Type**: scenario
**Priority**: (from scenario — happy-path)

**Codebase Mapping**:

New file `cmd/tools/list_test.go` — `TestListCmd_OutputsAllToolDescriptors` (happy-path test). Tested via `cmd.SetArgs([]string{"tools", "list"})` on the root command with injected mock MCP client returning a fixture tool list.

**Acceptance Criteria**:

- **AC-1**: Given the MCP Server is running with tools `ping` and `echo` registered, when `eve-realm tools list` is run, then the command exits with status code 0.
- **AC-2**: Given the MCP Server is running with tools registered, when `eve-realm tools list` is run, then stdout contains each tool's name, description, and input schema.
- **AC-3**: Given the MCP Server is running with multiple tools, when `eve-realm tools list` is run, then all registered tools appear in the output.
- **AC-4**: Given the output of `eve-realm tools list`, then the format is human-readable and parseable by AI tools.

**Implementation Notes**:

Feasibility: Feasible. Direct `ListTools` call; format results to stdout following the pattern from `skills/master/skill.go`. Mock the MCP client via the injected interface for tests.

**Test Expectations:**

- Must test: Mock returns two tools (`ping`, `echo`) — stdout contains both names, both descriptions, and both input schema strings; exit code is 0.
- Must test: Mock returns one tool — stdout contains that tool's name, description, input schema; no extra blank tool entries.
- Must test: Output structure is consistent — each tool section is separated and readable (no raw struct printing).
- Must NOT rely on: Real MCP Server connection. All calls go through `mockMCPClient`.

---

### SC-013: Tools invoke sends default empty input and returns JSON to stdout

**Entity**: `.software/scenarios/SC-013` (entity file not found — derived from brief and feasibility report)
**Type**: scenario
**Priority**: (happy-path)

**Codebase Mapping**:

New file `cmd/tools/invoke_test.go` — `TestInvokeCmd_DefaultEmptyInput`. Tested via `cmd.SetArgs([]string{"tools", "invoke", "ping"})` with no `--input` flag, asserting `mockMCPClient.InvokeTool` receives `{}` as input.

**Acceptance Criteria**:

- **AC-1**: Given the MCP Server is running and `ping` tool is registered, when `eve-realm tools invoke ping` is run (no `--input` flag), then the command exits with status code 0.
- **AC-2**: Given no `--input` flag is provided, when `tools invoke` executes, then `InvokeTool` is called with `{}` as the input argument.
- **AC-3**: Given the tool returns a JSON response, when `tools invoke` completes, then the JSON response is written to stdout.

**Implementation Notes**:

Feasibility: Feasible. The default value of the `--input` Cobra string flag is `"{}"`. The flag value is passed verbatim to `InvokeTool`.

**Test Expectations:**

- Must test: No `--input` flag provided — mock verifies `InvokeTool` called with `{}` as input string; stdout contains the mock's JSON response; exit code 0.
- Must test: Mock `InvokeTool` returns a multi-field JSON object — stdout contains the exact JSON string returned by the mock.
- Must NOT rely on: Real MCP Server. The `mockMCPClient.InvokeTool` field captures and returns controlled values.

---

### SC-014: Tools invoke passes --input flag value to the tool

**Entity**: `.software/scenarios/SC-014-tools-invoke-passes-input-flag-value-to-the-tool.md`
**Type**: scenario
**Priority**: (happy-path)

**Codebase Mapping**:

New file `cmd/tools/invoke_test.go` — `TestInvokeCmd_InputFlagPassthrough`. Tested via `cmd.SetArgs([]string{"tools", "invoke", "echo", "--input", `{"key":"value"}`})` asserting `mockMCPClient.InvokeTool` receives `{"key":"value"}` verbatim.

**Acceptance Criteria**:

- **AC-1**: Given the MCP Server is running and an `echo` tool is registered, when `eve-realm tools invoke echo --input '{"key":"value"}'` is run, then the command exits with status code 0.
- **AC-2**: Given `--input '{"key":"value"}'` is provided, when `tools invoke` executes, then `InvokeTool` receives `{"key":"value"}` as the input (not re-encoded or modified).
- **AC-3**: Given the tool returns a JSON response, when `tools invoke` completes, then the JSON response is written to stdout.
- **AC-4**: Given `--input` flag is set, then the `--input` flag value is passed verbatim to the gRPC `InvokeTool` RPC.

**Implementation Notes**:

Feasibility: Feasible. The `--input` flag is a Cobra `StringVar`. Its value is forwarded directly to `client.InvokeTool(name, inputFlag)` without transformation.

**Test Expectations:**

- Must test: `--input '{"key":"value"}'` — mock captures the input argument, assert it equals `{"key":"value"}` byte-for-byte; exit code 0.
- Must test: `--input '{}'` explicit empty object — mock receives `{}` (same as default, but explicitly set); behavior is identical.
- Must test: `--input` with nested JSON — input containing nested structures is passed verbatim without re-serialization.
- Must NOT rely on: JSON parsing or re-encoding of the `--input` value in production code. The value must travel as a raw string.

---

### SC-015: Tools commands exit non-zero with error when server unreachable

**Entity**: `.software/scenarios/SC-015` (entity file not found — derived from brief and feasibility report)
**Type**: scenario
**Priority**: (error-path)

**Codebase Mapping**:

- `cmd/tools/list_test.go` — `TestListCmd_ConnectionError`
- `cmd/tools/invoke_test.go` — `TestInvokeCmd_ConnectionError`

Both tests inject a `mockMCPClient` that returns a `mcpclient.ConnectionError`. Assert exit non-zero, stderr non-empty, stdout empty.

**Acceptance Criteria**:

- **AC-1**: Given the MCP Server is unreachable, when `eve-realm tools list` is executed, then the command exits with a non-zero status code.
- **AC-2**: Given the MCP Server is unreachable, when `eve-realm tools invoke <name>` is executed, then the command exits with a non-zero status code.
- **AC-3**: Given the server is unreachable, when either command fails, then a clear error message is written to stderr.
- **AC-4**: Given the server is unreachable, when either command fails, then stdout contains no output.

**Implementation Notes**:

Feasibility: Feasible. Use `errors.As(err, &connErr)` where `connErr` is of type `*mcpclient.ConnectionError`. Write the error message to `cmd.ErrOrStderr()` and return the error from `RunE`. Set `SilenceErrors = true` on root to prevent double-printing.

**Test Expectations:**

- Must test: `tools list` with mock returning `ConnectionError` — `cmd.Execute()` returns non-nil error; stderr contains error text; stdout is empty.
- Must test: `tools invoke <name>` with mock returning `ConnectionError` — same assertions.
- Must test: Error message in stderr is user-readable (not a raw Go error struct string).
- Must NOT rely on: Real network connectivity. `mockMCPClient.ListTools` and `mockMCPClient.InvokeTool` return a constructed `mcpclient.ConnectionError` directly.

---

### SC-016: Tools invoke exits non-zero and lists alternatives when tool not found

**Entity**: `.software/scenarios/SC-016` (entity file not found — derived from brief and feasibility report)
**Type**: scenario
**Priority**: (error-path)

**Codebase Mapping**:

`cmd/tools/invoke_test.go` — `TestInvokeCmd_ToolNotFound`. Mock `InvokeTool` returns `ToolNotFoundError`; a second mock `ListTools` call returns the available tools list. Assert exit non-zero, stderr contains not-found message and alternative tool names.

**Acceptance Criteria**:

- **AC-1**: Given the MCP Server is running but the named tool does not exist, when `eve-realm tools invoke <name>` is executed, then the command exits with a non-zero status code.
- **AC-2**: Given a tool-not-found error, when `tools invoke` handles it, then the not-found error message is written to stderr.
- **AC-3**: Given a tool-not-found error, when `tools invoke` handles it, then the names of all available tools are listed on stderr.

**Implementation Notes**:

Feasibility: Feasible. Pattern mirrors `skills/master/skill.go:L77-L89`. After detecting `ToolNotFoundError` via `errors.As`, make a secondary `client.ListTools()` call and write the alternatives to `cmd.ErrOrStderr()`. If the secondary `ListTools` call also fails, write only the original not-found error.

**Test Expectations:**

- Must test: `tools invoke unknown-tool` with `InvokeTool` returning `ToolNotFoundError` and `ListTools` returning `["ping", "echo"]` — stderr contains "unknown-tool" (or not-found text), "ping", and "echo"; exit non-zero; stdout empty.
- Must test: Secondary `ListTools` call fails (returns `ConnectionError`) — stderr still contains the original not-found error; command still exits non-zero.
- Must test: Tool not found with zero alternatives available (empty list from `ListTools`) — stderr contains not-found message; no crash or panic.
- Must NOT rely on: Real MCP Server. Both `InvokeTool` and the secondary `ListTools` call in the same test use configurable mock function fields.

---

### SC-017: MCP Server address resolved from env var, config file, or default

**Entity**: `.software/scenarios/SC-017` (entity file not found — derived from brief and feasibility report)
**Type**: scenario
**Priority**: (configuration)

**Codebase Mapping**:

`cmd/tools/list_test.go` — `TestListCmd_AddressResolution` (table-driven, three cases). Uses `t.TempDir()` for config file isolation. Sets/unsets `EVE_REALM_MCP_ADDR` env var per case. The mock MCP client captures the address passed to `mcpclient.NewClient(addr)` (requires address to be injectable or the test inspects via side-effect capture).

Existing file consumed (read-only): `internal/config/config.go` — `LoadHostConfig` and `Resolve`.

**Acceptance Criteria**:

- **AC-1**: Given `EVE_REALM_MCP_ADDR` env var is set to `custom-host:9999`, when a `tools` command executes, then `custom-host:9999` is used as the MCP Server address, regardless of the config file value.
- **AC-2**: Given `EVE_REALM_MCP_ADDR` is not set and `~/.eve-realm/eve-realm.yaml` contains `mcp_server_addr: yaml-host:8888`, when a `tools` command executes, then `yaml-host:8888` is used.
- **AC-3**: Given `EVE_REALM_MCP_ADDR` is not set and the config file is absent or has no `mcp_server_addr`, when a `tools` command executes, then `localhost:30051` is used.

**Implementation Notes**:

Feasibility: Feasible. `internal/config` already implements the three-layer resolution. The command factory must: (1) call `config.LoadHostConfig(configPath)` with the injectable `configPath`; (2) call `config.Resolve("EVE_REALM_MCP_ADDR", cfg.MCPServerAddr)`; (3) if the result is empty, use `mcpclient.DefaultMCPAddr`. For testability, `configPath` must come from the `NewToolsCmd(configPath string)` parameter — not hardcoded.

**Test Expectations:**

- Must test: Table-driven with three rows — env-var-wins, yaml-wins, default-used. Each row uses a temp config file (via `t.TempDir()`) and controls env var state via `t.Setenv()`.
- Must test: Env var set + YAML file present → env var value is used (YAML value ignored).
- Must test: Env var unset + YAML file with `mcp_server_addr` → YAML value is used.
- Must test: Env var unset + YAML file absent (or field missing) → `localhost:30051` is used.
- Must NOT rely on: The real `~/.eve-realm/eve-realm.yaml` file or the real `EVE_REALM_MCP_ADDR` env var from the test runner environment. Use `t.TempDir()` and `t.Setenv()` for full isolation.

---

## Documentation Tasks

### RELEASES.md Entry

**Required**: Always

Add an entry to RELEASES.md documenting:
- Sprint ID and title: SP-003 — Cobra command for MCP tool listing and invocation
- Summary of changes delivered: New `eve-realm tools list` and `eve-realm tools invoke` Cobra commands; introduction of Cobra as the root command framework; Cobra dependency added to `go.mod`; `cmd/main.go` refactored from bare arg parser to Cobra root.
- Entity IDs included: REQ-006, SC-012, SC-013, SC-014, SC-015, SC-016, SC-017
- Date of completion

This entry should be appended to the existing RELEASES.md file. Do not read or modify existing entries.

### README.md Update

**Required**: User-facing changes detected

Update README.md to reflect:
- New `eve-realm tools list` command — describe usage, output format (tool name, description, input schema), and when to use it.
- New `eve-realm tools invoke <name> [--input <json>]` command — describe positional argument, `--input` flag (default: `{}`), and JSON response output.
- MCP Server address configuration — document the three-layer precedence: `EVE_REALM_MCP_ADDR` env var overrides `mcp_server_addr` in `~/.eve-realm/eve-realm.yaml`, which overrides the default `localhost:30051`.
- Root command change — the binary now uses Cobra; update any usage examples that reference the old bare argument parser behavior.

## Pinned Entity Compliance

| Entity | Directive | How Addressed |
|--------|-----------|---------------|
| REQ-003: Cross-cutting requirements catalog for lazy-loaded sprint policy injection | Acknowledged — no spec-phase action required (REQ-003 is a catalog/loader entity; its directives are fulfilled by loading REQ-001 and REQ-002 and applying their rules below) | REQ-001 (TDD strategy) and REQ-002 (release process) are loaded and applied throughout this spec. |
| REQ-001: Test-Driven Development Strategy | Spec writer must generate a "Test Expectations" subsection per entity, mapping each acceptance criterion to the tests that must verify it, the mocking strategy, and the test type. | Every implementation section in this spec includes a "Test Expectations" subsection with specific, testable behaviors, mocking approach, and anti-patterns. |
| REQ-002: Sprint completion and release process | Phase 1 (spec-time): spec writer must decide version increment and README update requirement. | `readme_update_needed: true` is set in the brief (user-facing commands added). README.md Update section is included. Version increment decision is deferred to the release phase (post-implementation). |

## Out of Scope

- `cmd/auth/`, `cmd/marketplace/`, and `cmd/settings/` command packages — not part of this sprint.
- The `eve-realm marketplace register` command and embedded marketplace functionality (REQ-007) — separate sprint.
- Generic agent skill (`skills/agent/`) — not addressed in this sprint.
- gRPC client implementation changes — `internal/mcpclient` is consumed read-only; no modifications.
- Changes to `internal/config` — consumed read-only; address resolution is already implemented.
- The `protect` command for guarding Claude Code settings — explicitly excluded per REQ-007 notes.
- Production error retry logic for the secondary `ListTools` call in SC-016 — a single best-effort attempt is sufficient.

## Prerequisites

- `github.com/spf13/cobra` must be added to `go.mod` before any `cmd/tools/` code is written (`go get github.com/spf13/cobra`).
- `cmd/main.go` must be replaced with a Cobra root command before registering the `tools` subcommand.
- `internal/mcpclient.Client` (`ListTools` + `InvokeTool`) must be present and passing its own tests — confirmed implemented from SP-001.
- `internal/mcpclient.ConnectionError` and `internal/mcpclient.ToolNotFoundError` typed errors must be present — confirmed from SP-001.
- `internal/config.LoadHostConfig` and `internal/config.Resolve` must be present — confirmed from SP-001.
