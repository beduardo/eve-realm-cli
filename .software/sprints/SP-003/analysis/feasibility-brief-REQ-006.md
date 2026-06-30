# Feasibility Brief

**Sprint**: SP-003
**Project Root**: /Users/bruno/repo-pessoal/eve-realm/eve-realm-cli/main
**Target Entity**: REQ-006

## Entity Summary

REQ-006 adds `eve-realm tools list` and `eve-realm tools invoke <name> [--input <json>]` Cobra commands. These commands use the gRPC client (from SP-001/REQ-004) to call `ListTools` and `InvokeTool` RPCs on the MCP Server. Address resolution uses env var > config file > default.

## Acceptance Criteria (8 total)

1. `tools list` calls ListTools via gRPC, outputs tool descriptors to stdout
2. `tools invoke <name>` calls InvokeTool, returns JSON response to stdout
3. `tools invoke <name> --input '{...}'` passes JSON input to the tool
4. No `--input` flag sends empty `{}` as default
5. Server unreachable: non-zero exit + clear stderr error
6. Tool not found: non-zero exit + error + available tools listed
7. Address resolution: env var > config file > default (`localhost:30051`)
8. Command code lives in `cmd/tools/`

## Sprint Context
- Current entity count: 7 (1 REQ + 6 SCs)
- Scope score: 7/5
- Other entities in sprint: SC-012 through SC-017 (all scenarios for REQ-006)

## Dependencies
- REQ-004 (gRPC client library) — delivered in SP-001 (completed)
- Proto definitions for ListTools and InvokeTool RPCs — must exist in eve-realm-sdk

## Focus Questions
- Does the gRPC client from SP-001 already expose ListTools and InvokeTool methods?
- Is there an existing Cobra command pattern to follow for new subcommands?
- Does the config resolution already support reading `mcp_server_addr`?
- Are proto definitions for ListTools/InvokeTool available in the SDK submodule?
