---
description: "List and invoke MCP Server tools"
argument-hint: "[tool-name] [json-input]"
disable-model-invocation: true
---

# Master Skill

The master skill connects to the MCP Server and provides two modes of operation:

- **Discovery mode** (no arguments): Lists all tools available on the MCP Server,
  including each tool's name, description, and input schema. Use this to explore
  what capabilities are available.

- **Invocation mode** (tool name provided): Invokes the named tool with the supplied
  JSON input and returns the raw JSON response from the server.

## Usage

```
/eve-realm                        # discovery mode — list all tools
/eve-realm ping {}                # invocation mode — invoke the ping tool
```

## Configuration

The MCP Server address is resolved in the following order:

1. Environment variable `EVE_REALM_MCP_ADDR`
2. YAML field `mcp_server_addr` in `~/.eve-realm/eve-realm.yaml`
3. Default: `localhost:30051`
