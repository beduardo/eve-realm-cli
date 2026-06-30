# Codebase Analysis

**Sprint**: SP-001
**Analyzed**: 2026-06-28 (re-run — port updated to 30051)
**Entities Mapped**: 11

---

## Entity-to-Code Mapping

| Entity ID | Type | Related Files | Lines | Notes |
|-----------|------|---------------|-------|-------|
| REQ-004 | requirement | (no existing code) | - | New: `internal/mcpclient/`, `proto/mcp/v1/mcp.proto`, `gen/proto/mcp/v1/`. Default addr is `localhost:30051` (k3d NodePort, AC-6). |
| REQ-005 | requirement | (no existing code) | - | New: `skills/master/`. Reads MCP Server address from YAML config or env var. |
| SC-003 | scenario | (no existing code) | - | Tests ListTools; scenario text uses `localhost:50051` but production default is 30051 — test fixture must use an explicit address. |
| SC-004 | scenario | (no existing code) | - | Tests InvokeTool; same note as SC-003. |
| SC-005 | scenario | (no existing code) | - | Tests connection error on unreachable address; uses `localhost:59999` explicitly. |
| SC-006 | scenario | (no existing code) | - | Tests typed ToolNotFoundError; scenario text uses `localhost:50051` — same as SC-003/SC-004 note. |
| SC-007 | scenario | (no existing code) | - | Tests discovery mode (no args); covered by `skills/master/` unit tests. |
| SC-008 | scenario | (no existing code) | - | Tests invocation mode (tool name + input); covered by `skills/master/` unit tests. |
| SC-009 | scenario | (no existing code) | - | Tests graceful error when MCP Server unreachable; covered by `skills/master/` unit tests. |
| SC-00A | scenario | (no existing code) | - | Tests alternative tool suggestions on not-found; covered by `skills/master/` unit tests. |
| SC-00B | scenario | (no existing code) | - | End-to-end ping invocation; covered by `skills/master/` integration test with bufconn. |

**Note on SC-003/SC-004/SC-006 port value**: These scenario files use `localhost:50051` in their Steps section. REQ-004 AC-6 authoritative value is `localhost:30051`. Test fixtures must pass addresses explicitly; the `NewClient` function defaults to `localhost:30051` only when an empty string is received, so test coverage is unaffected by the scenario text.

All sprint entities are net-new. No existing production code relates to gRPC, mcpclient, or master skill.

---

## Implementation Patterns

### Pattern 1: Table-driven tests with `t.Run`

- **Reference**: `eve-cli/main/internal/marketplace/skilldesc_test.go:L132-L183`
- **Description**: `[]struct{ name string; ... }` tables with `t.Run(tc.name, ...)`. All new tests must follow this style.
- **Entities Using**: REQ-004, REQ-005, SC-003 through SC-00B

### Pattern 2: Interface-based dependency injection at I/O boundaries

- **Reference**: `eve-cli/main/internal/marketplace/orchestrator.go:L29-L33` (the `Executor` interface)
- **Description**: Test doubles implement a small interface at the consumer site. For gRPC, tests use `bufconn` to run a real gRPC server in-process without OS sockets.
- **Entities Using**: REQ-004, SC-003, SC-004, SC-005, SC-006

### Pattern 3: Config loading via YAML + env var overlay

- **Reference**: `eve-cli/main/internal/config/config.go:L144-L234`
- **Description**: `LoadHostConfigFrom` reads YAML, then calls `Resolve(envKey, yamlValue)` for each field. For eve-realm-cli: YAML key `mcp_server_addr`, env key `EVE_REALM_MCP_ADDR`. When both absent, `NewClient` defaults to `localhost:30051`.
- **Entities Using**: REQ-005, SC-007, SC-008, SC-009, SC-00A, SC-00B

### Pattern 4: SKILL.md with YAML frontmatter

- **Reference**: `eve-cli/main/internal/marketplace/skilldesc.go` (ParseSkillDescription)
- **Description**: Each skill directory contains `SKILL.md` with required YAML frontmatter fields: `description`, `argument-hint`, `disable-model-invocation`, optionally `allowed_tools`.
- **Entities Using**: REQ-005

### Pattern 5: Typed error structs with `errors.As` support

- **Reference**: `eve-cli/main/internal/marketplace/orchestrator.go` (error wrapping with `%w`)
- **Description**: `ConnectionError` and `ToolNotFoundError` are concrete structs implementing `error` for `errors.As` discrimination. Replaces gRPC status codes at the package boundary.
- **Entities Using**: REQ-004, SC-005, SC-006

### Pattern 6: YAML struct tags with `gopkg.in/yaml.v3`

- **Reference**: `eve-cli/main/internal/config/config.go:L129-L141`
- **Description**: Config structs use `yaml:"..."` tags. Missing-file is a non-error condition (returns zero-value config).
- **Entities Using**: REQ-005

### Pattern 7: Test function naming convention

- **Reference**: `eve-realm-cli/main/cmd/main_test.go:L5`, `eve-cli/main/internal/marketplace/orchestrator_test.go:L162`
- **Description**: `TestFunctionName_Scenario`. No testify — standard library `testing` only.
- **Entities Using**: REQ-004, REQ-005, all SCs

---

## Files to Create

| File | Purpose | Entities |
|------|---------|----------|
| `proto/mcp/v1/mcp.proto` | Proto definition for MCPService (ListTools, InvokeTool RPCs) | REQ-004 (AC-7) |
| `gen/proto/mcp/v1/mcp.pb.go` | Generated protobuf message types | REQ-004 (AC-7) |
| `gen/proto/mcp/v1/mcp_grpc.pb.go` | Generated gRPC client stub | REQ-004 (AC-7) |
| `internal/mcpclient/client.go` | `Client` type with `NewClient(addr)`, `ListTools`, `InvokeTool`; defaults to `localhost:30051` | REQ-004 (AC-1–AC-6) |
| `internal/mcpclient/errors.go` | `ConnectionError` and `ToolNotFoundError` typed errors | REQ-004 (AC-4, AC-5) |
| `internal/mcpclient/client_test.go` | Unit tests using bufconn in-process server | SC-003, SC-004, SC-005, SC-006 |
| `internal/config/config.go` | `HostConfig` with `MCPServerAddr`; `LoadHostConfig` + `Resolve` | REQ-005 (AC-7) |
| `internal/config/config_test.go` | Config loading, env var overlay, missing-file tests | REQ-005 (AC-7) |
| `skills/master/skill.go` | Master skill — discovery and invocation modes | REQ-005 (AC-1–AC-8) |
| `skills/master/skill_test.go` | Unit tests with mock mcpclient | SC-007–SC-00B |
| `skills/master/SKILL.md` | Skill metadata with YAML frontmatter | REQ-005 (AC-1) |
| `buf.yaml` | Buf configuration for proto generation | REQ-004 (AC-8) |

## Files to Modify

| File | Modification | Entities |
|------|--------------|----------|
| `go.mod` | Add `google.golang.org/grpc`, `google.golang.org/protobuf`, `gopkg.in/yaml.v3` | REQ-004, REQ-005 |
| `Makefile` | Add `proto` target; extend `.PHONY` | REQ-004 (AC-8) |

---

## Integration Points

1. **go.mod**: Empty of external deps — must add gRPC, protobuf, yaml.v3. SDK submodule not needed for SP-001.
2. **Makefile `proto` target**: `buf generate` or `protoc` to regenerate `gen/proto/mcp/v1/`. Idempotent, no conflict with existing targets.
3. **`internal/mcpclient` -> `skills/master`**: Primary dependency chain: `skills/master/skill.go` -> `internal/mcpclient.Client` -> `gen/proto/mcp/v1.MCPServiceClient` -> gRPC.
4. **`internal/config` -> `skills/master`**: Config resolution: env var > YAML > `localhost:30051` default via `NewClient("")`.
5. **Default port `localhost:30051`**: Per AC-6, `NewClient(addr)` defaults to `localhost:30051` when addr is empty. k3d NodePort for MCP Server gRPC service.

---

## Test Patterns

| Pattern | Reference | Usage |
|---------|-----------|-------|
| Table-driven `t.Run` | eve-cli `skilldesc_test.go` | All test files |
| `t.Helper()` + `t.TempDir()` | eve-cli `orchestrator_test.go` | Config file tests, bufconn teardown |
| Interface mock (no mock libs) | eve-cli `orchestrator_test.go` | `mockMCPClient` for skill tests |
| In-process gRPC with bufconn | Standard pattern | `internal/mcpclient/client_test.go` |
| External test package `package foo_test` | eve-cli convention | Black-box testing |
| Standard library `testing` only | eve-cli convention | No testify/gomock |

---

## Technical Notes

- **Port change**: REQ-004 AC-6 specifies `localhost:30051` (k3d NodePort). SC-003, SC-004, SC-006 still reference `50051` in step text — documentation drift, not a code concern.
- **No gRPC anywhere in the workspace**: First gRPC dependency introduced by this sprint.
- **No SDK submodule on disk**: Not needed for SP-001; all proto and client code is local.
- **No existing `internal/` packages**: Both `internal/mcpclient` and `internal/config` created from scratch.
- **`cmd/main.go` is a 20-line scaffold**: No Cobra, no config, no gRPC. SP-001 does not modify it.
- **`buf` vs `protoc`**: Either works. `buf` is modern and avoids managing protoc plugins.
- **Skill registration unresolved**: Simplest path is manual `extraKnownMarketplaces` entry in `~/.claude/settings.json`.
- **Eve-cli has no gRPC source files**: No working proto/gRPC code to copy from reference codebase.
