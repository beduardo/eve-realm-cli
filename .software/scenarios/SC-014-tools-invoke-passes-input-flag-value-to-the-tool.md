---
content_hash: 03f757759bea91171c52212ae6163171bdcc7476ec6466db370410d09f70d44d
created: "2026-06-29"
id: SC-014
related_changes: []
related_reqs:
    - REQ-006
related_testcases: []
source: manual
status: validated
tags:
    - cobra
    - tools-invoke
    - input-flag
title: Tools invoke passes --input flag value to the tool
type: happy-path
updated: "2026-06-29"
---

# SC-014: Tools invoke passes --input flag value to the tool

## Preconditions

- MCP Server running with a tool that echoes or processes its input
- MCP Server address is resolvable

## Steps

1. Run `eve-realm tools invoke <name> --input '{"key":"value"}'`

## Expected Result

- Command exits with zero status code
- The tool receives `{"key":"value"}` as its input (not the default `{}`)
- stdout contains the tool's JSON response
- The `--input` flag value is passed verbatim to the gRPC `InvokeTool` RPC
