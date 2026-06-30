# Feasibility Report: REQ-005

**Entity**: Master skill for MCP tool discovery and invocation
**Type**: requirement
**Analyzed**: 2026-06-28 (re-run — gRPC port updated to 30051)

## Recommendation

**PROCEED-WITH-CAVEATS**

REQ-005 is implementable within SP-001 in strict sequence after REQ-004. The port change (30051 vs 50051) has zero impact on the skill layer — the default lives in `mcpclient.NewClient()`, not in the skill. Three standing caveats: sequential dependency on REQ-004, `internal/config/` package must be created, and SKILL.md marketplace registration path requires a design decision.

## Prerequisite Status

| Prerequisite | Type | Required Status | Current Status | Ready? |
|---|---|---|---|---|
| REQ-004: gRPC tool client | requirement | implemented | active (not yet implemented) | No — hard sequential |
| SC-007–SC-00B | scenarios | validated | validated | Yes |
| `internal/mcpclient/` package | code | exists | does not exist | No — delivered by REQ-004 |
| `internal/config/` package | code | exists | does not exist | No — must be created |

## Dependency Graph Analysis

- **Direct dependencies**: 1 software entity (REQ-004) + 3 code artifacts
- **Blocking dependencies**: 0 at entity-status level; 1 at implementation-ordering level

## Port Change Impact

The `localhost:30051` default lives in `internal/mcpclient.NewClient()`, not in the skill. The skill's config resolution chain: env var (`EVE_REALM_MCP_ADDR`) > YAML config (`mcp_server_addr`) > empty string to `NewClient("")` which substitutes `localhost:30051`. Zero impact on REQ-005 scope or tests.

## Complexity Estimate

**Size**: M (skill alone) / L (with config package and registration)

| Factor | Assessment | Notes |
|---|---|---|
| Code changes | M | ~300-450 LOC across 4-6 files |
| Test coverage | M | Mock mcpclient interface; SC-00B needs in-process bufconn server |
| Integration risk | Medium | Skill surfacing to Claude Code is the primary unknown |
| Architectural impact | Medium | Establishes what a "skill" is in this codebase |

## Risk Factors

| Risk | Likelihood | Impact | Mitigation |
|---|---|---|---|
| Skill registration undefined | High | High | Manual `extraKnownMarketplaces` entry for SP-001 |
| REQ-004 code not ready | High | High | Serialize implementation ordering |
| `internal/config/` ownership gap | Medium | Medium | Created as part of REQ-005 implementation |
| SC-00B E2E test | Medium | Medium | In-process bufconn gRPC server |

## Blockers

| Blocker | Severity | Resolution Path |
|---------|----------|-----------------|
| Skill surfacing undefined | Major | Manual `extraKnownMarketplaces` for SP-001; automate later |
| `internal/config/` does not exist | Minor | Created during REQ-005 implementation using eve-cli reference |
| `go.mod` missing `gopkg.in/yaml.v3` | Minor | `go get` as first step |
