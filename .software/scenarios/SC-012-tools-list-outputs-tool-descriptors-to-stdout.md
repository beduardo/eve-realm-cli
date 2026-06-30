---
content_hash: 527edf77de61eade80f27d1f9380a0f6587d1641b39ee01b9dfe694347cfde62
created: "2026-06-29"
id: SC-012
related_changes: []
related_reqs:
    - REQ-006
related_testcases: []
source: manual
status: validated
tags:
    - cobra
    - tools-list
    - happy-path
title: Tools list outputs tool descriptors to stdout
type: happy-path
updated: "2026-06-29"
---

# SC-012: Tools list outputs tool descriptors to stdout

## Preconditions

- MCP Server running with multiple tools registered (e.g., `ping`, `echo`)
- MCP Server address is resolvable (via config, env var, or default)

## Steps

1. Run `eve-realm tools list`

## Expected Result

- Command exits with zero status code
- stdout contains each tool's name, description, and input schema
- All registered tools appear in the output
- Output format is human-readable and parseable by AI tools
