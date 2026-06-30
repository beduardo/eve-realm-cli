# Releases

---

## 0.2.0 — 2026-06-30 (git: c34d663)

### Sprint: SP-001 — MCP gRPC Connection Client

**Version increment**: minor

**Changes**:
- gRPC client library (`internal/mcpclient/`) with typed errors (`ConnectionError`, `ToolNotFoundError`) and bufconn-tested transport
- Master skill (`skills/master/`) for tool discovery and invocation via the MCP Server
- Config resolver (`internal/config/`) with env var / YAML / default address precedence
- Proto definition (`proto/mcp/v1/mcp.proto`) and generated Go stubs (`gen/proto/mcp/v1/`)
- Makefile `proto` target for stub regeneration

**Entities affected**: REQ-004, REQ-005, SC-003, SC-004, SC-005, SC-006, SC-007, SC-008, SC-009, SC-00A, SC-00B

---

## 0.3.0 — 2026-06-30 (git: 52ed513)

### Sprint: SP-003 — Cobra command for MCP tool listing and invocation

**Version increment**: minor

**Changes**:
- Adopted Cobra as root command framework, replacing the bare argument parser in `cmd/main.go`
- New `eve-realm tools list` command for discovering MCP Server tools with human-readable output (name, description, input schema)
- New `eve-realm tools invoke <name> [--input <json>]` command for calling MCP Server tools with verbatim JSON pass-through and not-found alternative listing
- Added `github.com/spf13/cobra` dependency to `go.mod`
- MCP Server address resolution via three-tier precedence (env var → YAML config → default)

**Entities affected**: REQ-006, SC-012, SC-013, SC-014, SC-015, SC-016, SC-017

---
