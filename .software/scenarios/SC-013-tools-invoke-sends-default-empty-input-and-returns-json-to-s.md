---
content_hash: 45f08eb37a555bac131abdaaa991d7efa2222290550f6066a1804db08248ccba
created: "2026-06-29"
id: SC-013
related_changes: []
related_reqs:
    - REQ-006
related_testcases: []
source: manual
status: validated
tags:
    - cobra
    - tools-invoke
    - happy-path
title: Tools invoke sends default empty input and returns JSON to stdout
type: happy-path
updated: "2026-06-29"
---

# SC-013: Tools invoke sends default empty input and returns JSON to stdout

## Preconditions

- MCP Server running with the `ping` tool registered
- MCP Server address is resolvable

## Steps

1. Run `eve-realm tools invoke ping` (no `--input` flag)

## Expected Result

- Command exits with zero status code
- stdout contains the tool's JSON response (e.g., `{"message": "pong", ...}`)
- The tool receives an empty JSON object `{}` as input (default when `--input` is omitted)
- No additional formatting or wrapper around the JSON output
