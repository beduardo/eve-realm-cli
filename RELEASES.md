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
