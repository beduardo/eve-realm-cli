# Codebase Analysis

**Sprint**: SP-003
**Analyzed**: 2026-06-30
**Entities Mapped**: 7

## Entity-to-Code Mapping

| Entity ID | Type | Related Files | Notes |
|-----------|------|---------------|-------|
| REQ-006 | requirement | (new) | New implementation in `cmd/tools/`; integrates `internal/mcpclient` and `internal/config` from SP-001 |
| SC-012 | scenario | (new) | Happy path for `eve-realm tools list`; tests will use mockMCPClient + cmd.SetArgs pattern |
| SC-013 | scenario | (new) | Happy path for `eve-realm tools invoke ping` with default `{}` input |
| SC-014 | scenario | (new) | `--input` flag passthrough; tests verify the flag value reaches `InvokeTool` |
| SC-015 | scenario | (new) | Error path: unreachable server; command must exit non-zero and write to stderr |
| SC-016 | scenario | (new) | Error path: tool not found; command must exit non-zero and list alternatives on stderr |
| SC-017 | scenario | `internal/config/config.go` | Address resolution already implemented via `LoadHostConfig` + `Resolve`; SP-003 consumes it |

## Implementation Patterns

### Pattern 1: Constructor function returning `*cobra.Command`

- **Reference**: eve-cli `cmd/eve5/marketplace/marketplace.go:L13-L30`
- **Description**: Each command group is created by a `NewXxxCmd() *cobra.Command` factory. Subcommands are built by lowercase `newXxxCmd()` helpers in the same package and added via `cmd.AddCommand(...)`. The caller (main.go) calls `NewXxxCmd()` and adds it to the root.

### Pattern 2: RunE with stderr/stdout routed through Cobra

- **Reference**: eve-cli `cmd/eve5/marketplace/list.go:L34-L69`
- **Description**: Commands use `RunE` (not `Run`) so errors propagate as a non-zero exit code. Output goes to `cmd.OutOrStdout()` for success and `cmd.ErrOrStderr()` for errors — never `fmt.Println` or `os.Stderr` directly. Critical for testability via buffer injection.

### Pattern 3: Interface-based mock for gRPC client

- **Reference**: `skills/master/skill.go:L13-L16`, `skills/master/skill_test.go:L27-L38`
- **Description**: The `MCPClient` interface has `ListTools` and `InvokeTool` methods. In tests, a `mockMCPClient` struct with configurable function fields implements this interface. `cmd/tools/` should define its own local interface (Go structural typing handles satisfaction).

### Pattern 4: Address resolution via `config.LoadHostConfig` + `config.Resolve`

- **Reference**: `internal/config/config.go:L11-L40`
- **Description**: `LoadHostConfig(path)` reads YAML from `~/.eve-realm/eve-realm.yaml` and returns zero-value `HostConfig{}` if the file is missing (no error). `Resolve(envKey, yamlValue)` returns the env var when non-empty, otherwise falls back to YAML value. `mcpclient.NewClient("")` substitutes `localhost:30051` as default.

### Pattern 5: Cobra command testing via `cmd.SetArgs()` + `cmd.Execute()` + `bytes.Buffer`

- **Reference**: eve-cli `cmd/eve5/marketplace/marketplace_test.go:L56-L65`
- **Description**: Tests build the command with the factory function, inject `bytes.Buffer` via `cmd.SetOut(&stdout)` / `cmd.SetErr(&stderr)`, set args with `cmd.SetArgs(args)`, call `cmd.Execute()`, and assert on the returned error, `stdout.String()`, and `stderr.String()`.

### Pattern 6: Typed error detection with `errors.As` for non-zero exit

- **Reference**: `internal/mcpclient/errors.go`, `skills/master/skill.go:L67-L93`
- **Description**: `ConnectionError` and `ToolNotFoundError` are pointer-receiver structs supporting `errors.As`. For `ToolNotFoundError`, a secondary `ListTools` call builds the alternatives list for stderr output.

### Pattern 7: Test function naming `TestFunctionName_Scenario`

- **Reference**: `internal/mcpclient/client_test.go:L100`, `internal/config/config_test.go:L13`
- **Description**: All test functions follow `TestFunctionName_Scenario`. No testify or external assertion libraries — standard `if got != want { t.Errorf(...) }` assertions throughout.

## Files to Create

| File | Purpose | Based On |
|------|---------|----------|
| `cmd/tools/tools.go` | `NewToolsCmd()` factory registering `list` and `invoke` subcommands | eve-cli `cmd/eve5/marketplace/marketplace.go` |
| `cmd/tools/list.go` | `newListCmd()` — calls `ListTools`, formats tool descriptors to stdout | eve-cli `cmd/eve5/marketplace/list.go` |
| `cmd/tools/invoke.go` | `newInvokeCmd()` — positional arg for tool name, `--input` flag, calls `InvokeTool` | eve-cli `cmd/eve5/marketplace/register.go` |
| `cmd/tools/tools_test.go` | Shared test helpers: `runToolsCmd(t, args...)`, mock MCP client | eve-cli `cmd/eve5/marketplace/marketplace_test.go` |
| `cmd/tools/list_test.go` | Tests for `tools list` command | SC-012, SC-015, SC-017 |
| `cmd/tools/invoke_test.go` | Tests for `tools invoke` command | SC-013, SC-014, SC-015, SC-016 |

## Files to Modify

| File | Modification |
|------|--------------|
| `cmd/main.go` | Add Cobra root command wiring; import `github.com/spf13/cobra` and `cmd/tools`; register `tools` subcommand |
| `go.mod` | Add `github.com/spf13/cobra` dependency |

## Integration Points

### 1. MCP client construction in `cmd/tools/`

Both subcommands resolve the MCP Server address using `config.LoadHostConfig(path)` and `config.Resolve("EVE_REALM_MCP_ADDR", cfg.MCPServerAddr)`, then call `mcpclient.NewClient(addr)`. For tests, the MCP client must be injectable (passed into the factory function).

### 2. `MCPClient` interface reuse vs. redefinition

`skills/master` already declares `MCPClient` interface. `cmd/tools/` should define a local interface (same two methods — Go structural typing handles satisfaction) to avoid coupling `cmd/` to `skills/`.

### 3. Root command refactor from bare `main.go` to Cobra

The current `main.go` is a 21-line custom arg parser, not Cobra. SP-003 requires introducing Cobra and making `tools` a subcommand of the root. The existing `version` arg handler should be replaced with a Cobra `version` command or the built-in `--version` flag.

### 4. Config file path injection

For testability, the config file path should be injectable into the command factory (e.g., `NewToolsCmd(configPath string)`) rather than computed from `os.UserHomeDir()` at call time. This allows `t.TempDir()` to be used in SC-017 tests.

## Technical Notes

- Cobra is not yet in `go.mod`. Adding it is the first required step.
- `cmd/main.go` is a bare arg parser, not Cobra-based. SP-003 must introduce the root Cobra command.
- The `MCPClient` interface duplication approach is recommended to keep `cmd/` and `skills/` independently evolvable.
- SC-015: Use `cmd.SilenceErrors = true` at root level to avoid Cobra double-printing errors.
- SC-016: Secondary `ListTools` call mirrors the pattern in `skills/master/skill.go:L77-L89`.
- Address resolution for SC-017: three-layer precedence is fully implemented in `internal/config`. The command must apply `DefaultMCPAddr` when the resolved value is empty.
- No `cmd/auth/`, `cmd/marketplace/`, or `cmd/settings/` exist yet. SP-003 only needs `cmd/tools/` plus the Cobra root wiring.
