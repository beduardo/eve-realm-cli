# Releases

---

## SP-001 — MCP gRPC Connection Client (2026-06-28)

This sprint delivers the foundational gRPC transport layer between the CLI and the MCP
Server. It introduces the gRPC client library (`internal/mcpclient/`) with a typed error
model, the master skill (`skills/master/`) that exposes tool discovery and invocation via
the CLI, and the config resolver (`internal/config/`) that supplies the MCP Server address
at runtime. Proto definitions and generated stubs (`proto/`, `gen/`) are included alongside
a Makefile `proto` target that regenerates the stubs from source.

**Entities**: REQ-004, REQ-005, SC-003, SC-004, SC-005, SC-006, SC-007, SC-008, SC-009, SC-00A, SC-00B

---
