# eve-realm-cli

Thin client binary for the Eve Realm platform. Provides authentication, workspace
connection, marketplace management, and AI-assisted capabilities via the MCP Server.

## Build and Install

| Command | Purpose |
|---------|---------|
| `make build` | Build binary to `dist/eve-realm` |
| `make test` | Run all Go tests |
| `make install` | Copy binary to `/usr/local/bin/` |
| `make clean` | Remove `dist/` |
| `make version` | Show current version |
| `make proto` | Regenerate Go stubs from `proto/mcp/v1/mcp.proto` |
| `make release-patch` | Test, bump patch, build, install, verify |
| `make release-minor` | Test, bump minor, build, install, verify |
| `make release-major` | Test, bump major, build, install, verify |

## Development

### Regenerating Proto Stubs

The gRPC service contract is defined in `proto/mcp/v1/mcp.proto`. Generated Go stubs
live in `gen/proto/mcp/v1/`. To regenerate after modifying the proto file:

```bash
make proto
```

Requires `protoc`, `protoc-gen-go`, and `protoc-gen-go-grpc` installed in `$(go env GOPATH)/bin`.

## Tools Commands

The `tools` command group connects to the MCP Server and exposes two subcommands:
`list` (discover available tools) and `invoke` (call a tool by name).

### `eve-realm tools list`

Lists every tool registered on the MCP Server. For each tool the output shows its
name, description, and input schema.

```bash
eve-realm tools list
```

Example output:

```
Name:         ping
Description:  Check whether the MCP Server is reachable
Input Schema: {}

Name:         summarise
Description:  Summarise a block of text
Input Schema: {"type":"object","properties":{"text":{"type":"string"}},"required":["text"]}
```

Use this command to discover what tools the connected MCP Server exposes before
invoking one.

### `eve-realm tools invoke`

Invokes a named tool on the MCP Server and prints the JSON response to stdout.

```bash
eve-realm tools invoke <tool-name> [--input <json>]
```

| Argument / Flag | Description | Default |
|-----------------|-------------|---------|
| `<tool-name>` | Name of the tool to invoke (required positional argument) | — |
| `--input <json>` | JSON object passed verbatim as the tool input | `{}` |

Examples:

```bash
# Invoke without input (uses the default empty object)
eve-realm tools invoke ping

# Invoke with a JSON input object
eve-realm tools invoke summarise --input '{"text":"Hello world"}'
```

The raw JSON response from the tool is written to stdout. On error, a descriptive
message is written to stderr and the process exits non-zero. If the tool name is not
found, the error message includes the list of available tools.

## Master Skill

The `skills/master/` directory contains the master skill, published as a single
marketplace entry named `eve-realm`. It dynamically surfaces all tools available on
the MCP Server at runtime.

### Registering the skill in Claude Code

Add the following to `~/.claude/settings.json` under `extraKnownMarketplaces`:

```json
{
  "extraKnownMarketplaces": [
    {
      "name": "eve-realm",
      "skillsPath": "/path/to/eve-realm-cli/skills/master"
    }
  ]
}
```

Replace `/path/to/eve-realm-cli` with the absolute path to your checkout. Reload
Claude Code to pick up the new marketplace entry.

## MCP Server Configuration

The MCP Server address used by `eve-realm tools` is resolved in the following order
(first non-empty value wins):

1. **Environment variable** `EVE_REALM_MCP_ADDR`:

   ```bash
   export EVE_REALM_MCP_ADDR=localhost:9090
   ```

2. **YAML config** `mcp_server_addr` in `~/.eve-realm/eve-realm.yaml`:

   ```yaml
   mcp_server_addr: localhost:9090
   ```

3. **Default**: `localhost:30051`
