---
content_hash: 13b4bac5dc72f551210db8ef2d647bbbde99a51418cd6155bd50566e45451929
created: "2026-06-29"
id: SC-017
related_changes: []
related_reqs:
    - REQ-006
related_testcases: []
source: manual
status: validated
tags:
    - cobra
    - tools
    - config
    - address-resolution
title: MCP Server address resolved from env var, config file, or default
type: happy-path
updated: "2026-06-29"
---

# SC-017: MCP Server address resolved from env var, config file, or default

## Preconditions

- MCP Server running on a non-default address (e.g., `localhost:9999`)

## Steps

1. Set `EVE_REALM_MCP_ADDR=localhost:9999` and run `eve-realm tools list`
2. Unset the env var, write `mcp_server_addr: localhost:9999` to `~/.eve-realm/eve-realm.yaml`, and run `eve-realm tools list`
3. Unset the env var and remove the config file, then run `eve-realm tools list`

## Expected Result

- Step 1: Command connects to `localhost:9999` (env var takes precedence)
- Step 2: Command connects to `localhost:9999` (config file used when env var is absent)
- Step 3: Command connects to `localhost:30051` (default when neither is set)
- Precedence order: env var > config file > default
