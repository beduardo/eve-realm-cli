# Codebase Analysis Brief

**Sprint**: SP-003
**Project Root**: /Users/bruno/repo-pessoal/eve-realm/eve-realm-cli/main
**Entity IDs**: REQ-006, SC-012, SC-013, SC-014, SC-015, SC-016, SC-017

## Entity Details

### REQ-006: Cobra command for MCP tool listing and invocation
- Type: requirement
- Status: active
- Tags: cobra, cli, tools, grpc
- The CLI exposes `eve-realm tools list` and `eve-realm tools invoke <name> [--input <json>]`. Both use the gRPC client library (REQ-004) to communicate with the MCP Server. Address resolution follows: env var > config file > default (`localhost:30051`). Command code lives in `cmd/tools/`.

### SC-012: Tools list outputs tool descriptors to stdout
- Type: scenario
- Status: validated
- Tags: cobra, tools-list, happy-path
- Run `eve-realm tools list` against a running MCP Server; expect zero exit, stdout with tool name/description/schema for all registered tools.

### SC-013: Tools invoke sends default empty input and returns JSON to stdout
- Type: scenario
- Status: validated
- Tags: cobra, tools-invoke, happy-path
- Run `eve-realm tools invoke ping` without `--input`; expect zero exit, JSON response on stdout, empty `{}` sent as input.

### SC-014: Tools invoke passes --input flag value to the tool
- Type: scenario
- Status: validated
- Tags: cobra, tools-invoke, input-flag
- Run `eve-realm tools invoke <name> --input '{"key":"value"}'`; expect the JSON passed verbatim to `InvokeTool`, response on stdout.

### SC-015: Tools commands exit non-zero with error when server unreachable
- Type: scenario
- Status: validated
- Tags: cobra, tools, error, connection
- Run `tools list` and `tools invoke` with no server; expect non-zero exit, clear stderr error with server address, empty stdout.

### SC-016: Tools invoke exits non-zero and lists alternatives when tool not found
- Type: scenario
- Status: validated
- Tags: cobra, tools-invoke, error, not-found
- Run `eve-realm tools invoke nonexistent`; expect non-zero exit, stderr with "not found" message plus available tool names.

### SC-017: MCP Server address resolved from env var, config file, or default
- Type: scenario
- Status: validated
- Tags: cobra, tools, config, address-resolution
- Test all three address resolution paths: EVE_REALM_MCP_ADDR env var, `mcp_server_addr` in config YAML, default `localhost:30051`. Precedence: env > config > default.

## Focus Areas
- How existing Cobra commands are structured (cmd/ directory patterns)
- How the gRPC client from SP-001 is wired and used
- Config/address resolution patterns already in place
- Test patterns for CLI commands
