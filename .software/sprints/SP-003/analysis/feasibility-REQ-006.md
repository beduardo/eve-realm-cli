# Feasibility Report: REQ-006

**Entity**: Cobra command for MCP tool listing and invocation
**Type**: requirement
**Analyzed**: 2026-06-30

## Recommendation

**PROCEED-WITH-CAVEATS**

All core prerequisites are implemented and ready. Two caveats: (1) Cobra is not yet declared as a Go module dependency; (2) `cmd/main.go` is a bare `fmt.Println` stub with no Cobra root command — wiring the `tools` subcommand requires bootstrapping Cobra first.

## Prerequisite Status

| Prerequisite | Required Status | Current Status | Ready? |
|---|---|---|---|
| REQ-004: gRPC tool client | implemented | implemented | Yes |
| REQ-005: Master skill | implemented | implemented | Yes |
| `internal/mcpclient.Client` — `ListTools` + `InvokeTool` | exists | present | Yes |
| `internal/mcpclient.ConnectionError` typed error | exists | present | Yes |
| `internal/mcpclient.ToolNotFoundError` typed error | exists | present | Yes |
| `internal/config.LoadHostConfig` + `Resolve` | exists | present | Yes |
| `github.com/spf13/cobra` Go module dependency | declared in `go.mod` | absent | **No** |
| Cobra root command in `cmd/main.go` | must exist | absent — stub only | **No** |

## Dependency Graph Analysis

- REQ-006 -> REQ-004: consumes `internal/mcpclient.Client` — implemented, all required signatures present
- REQ-006 -> REQ-005: master skill calls gRPC client directly (not via Cobra commands). REQ-006 does not block REQ-005's current implementation
- Blocking dependencies: 0

## Complexity Estimate

**Size**: M (200-350 LOC production, 200-300 LOC test, 4-6 files)

| Factor | Assessment | Notes |
|---|---|---|
| Code changes | M | New `cmd/tools/` package + `cmd/main.go` refactor |
| Test coverage | M | Table-driven unit tests using mock MCPClient |
| Integration risk | Low | mcpclient and config packages are fully tested |
| Architectural impact | Low | Follows existing eve-cli patterns exactly |

## Risk Factors

| Risk | Likelihood | Impact | Mitigation |
|---|---|---|---|
| Cobra not in `go.mod` | Certain | High | `go get github.com/spf13/cobra` as first step |
| SC-016 secondary ListTools call adds latency | Medium | Low | Mirror proven pattern from `skills/master/skill.go` |
| SC-012 output format underspecified | Medium | Low | Use format from `skills/master/skill.go` |
| `cmd/main.go` refactor from bare parser to Cobra | Low | Low | File is 21 lines; straightforward replacement |

## Blockers

| Blocker | Severity | Resolution |
|---|---|---|
| `github.com/spf13/cobra` absent from `go.mod` | Critical | `go get github.com/spf13/cobra` as first implementation step |
| `cmd/main.go` has no Cobra root command | Major | Replace 21-line stub with Cobra root (following eve-cli pattern) |

## Scenario Feasibility

| Scenario | Verdict | Notes |
|---|---|---|
| SC-012: tools list outputs descriptors | Feasible | Direct `client.ListTools` call, format to stdout |
| SC-013: tools invoke default empty input | Feasible | Default `--input` flag value is `"{}"` |
| SC-014: tools invoke --input passthrough | Feasible | Cobra string flag, pass verbatim to InvokeTool |
| SC-015: non-zero exit on unreachable | Feasible | `errors.As(err, &ConnectionError{})`, write to stderr |
| SC-016: not-found + list alternatives | Feasible | Same pattern as `skills/master/skill.go` `formatError()` |
| SC-017: address resolution precedence | Feasible | Fully implemented in `internal/config`; consume with fallback to `DefaultMCPAddr` |

## Notes

- The `config.Resolve` function returns empty string when both env var and yaml value are absent. The command must apply `mcpclient.DefaultMCPAddr` when resolved value is empty.
- For testability, config file path should be injectable into command factory.
- Estimated new files: `cmd/tools/tools.go`, `cmd/tools/list.go`, `cmd/tools/invoke.go`, `cmd/tools/tools_test.go`, `cmd/tools/list_test.go`, `cmd/tools/invoke_test.go`. `cmd/main.go` is replaced. Total: 7 files.
