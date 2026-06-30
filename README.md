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

## Master Skill

The `skills/master/` directory contains the master skill, published as a single
marketplace entry named `eve-realm`. It dynamically surfaces all tools available on
the MCP Server at runtime.

**Discovery mode** (no arguments):

```
/eve-realm
```

Lists all registered tools with their names, descriptions, and input schemas.

**Invocation mode** (tool name + optional JSON input):

```
/eve-realm ping {}
```

Invokes the named tool and returns the raw JSON response.

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

The MCP Server address is resolved in the following order (first non-empty value wins):

1. **Environment variable** `EVE_REALM_MCP_ADDR`:

   ```bash
   export EVE_REALM_MCP_ADDR=localhost:9090
   ```

2. **YAML config** `mcp_server_addr` in `~/.eve-realm/eve-realm.yaml`:

   ```yaml
   mcp_server_addr: localhost:9090
   ```

3. **Default**: `localhost:30051`
