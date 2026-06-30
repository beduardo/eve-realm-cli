# Feasibility Report: REQ-004

**Entity**: gRPC tool client for MCP Server
**Type**: requirement
**Analyzed**: 2026-06-28 (re-run — port change from 50051 to 30051, AC-6 default address update)

## Recommendation

**PROCEED-WITH-CAVEATS**

REQ-004 is internally coherent and all four dependent scenarios are in `validated` status. The port change from 50051 to 30051 (k3d NodePort) is a single-site constant change with no structural impact. Three setup gaps remain: `go.mod` has zero dependencies, no proto definition exists, and gRPC tooling is absent. One documentation inconsistency: SC-003, SC-004, SC-006 still reference `localhost:50051` in step text.

## Prerequisite Status

| Prerequisite | Type | Required Status | Current Status | Ready? |
|-------------|------|-----------------|----------------|--------|
| SC-003 | scenario | validated | validated | Yes |
| SC-004 | scenario | validated | validated | Yes |
| SC-005 | scenario | validated | validated | Yes |
| SC-006 | scenario | validated | validated | Yes |

## Dependency Graph Analysis

- **Direct dependencies**: 4 (SC-003, SC-004, SC-005, SC-006)
- **Blocking dependencies**: 0
- REQ-005 depends on REQ-004 (consumer direction), not the reverse

## Port Change Impact

The change from 50051 to 30051 is a single-site change: the default constant in `NewClient`.

Recommended implementation: exported `const DefaultMCPAddr = "localhost:30051"` in `internal/mcpclient/client.go`, with `NewClient` using it when `addr` is empty.

**Scenario text drift**: SC-003, SC-004, SC-006 reference `localhost:50051`. Documentation drift only — tests pass explicit addresses, so the port in scenario prose does not affect test validity.

## Complexity Estimate

**Size**: M

| Factor | Assessment | Notes |
|--------|-----------|-------|
| Code changes | M | ~250-400 LOC production code |
| Test coverage | M | Four scenarios map to bufconn tests; ~200-300 LOC tests |
| Integration risk | Low | Pure internal package, no coupling to `cmd/` |
| Architectural impact | Low | New directories only; go.mod and Makefile modified |

## Risk Factors

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| gRPC tooling absent | Confirmed | Medium | Install buf + protoc-gen-go plugins; OR commit generated code |
| go.mod is bare | Confirmed | Low | `go get` as first step |
| No proto definition exists | Confirmed | Medium | Author minimal MCPService (2 RPCs) in spec phase |
| Proto contract divergence from MCP Server | Low | Medium | Define proto as canonical contract for SP-001 |
| SDK submodule not initialized | Confirmed | None | Explicitly deferred |
| Scenario text references 50051 | Confirmed | None | Documentation drift; should be updated for accuracy |

## Blockers

| Blocker | Severity | Resolution Path |
|---------|----------|-----------------|
| gRPC tooling not installed | Major | Install `buf` + Go protoc plugins; or commit generated files |
| No proto definition exists | Major | Author `proto/mcp/v1/mcp.proto` in spec phase |
| go.mod has no dependencies | Minor | `go get` as first implementation step |
| `make proto` target missing | Minor | Add `.PHONY: proto` target to Makefile |
