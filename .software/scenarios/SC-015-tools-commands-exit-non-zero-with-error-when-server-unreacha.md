---
content_hash: 9377998943a2b7b06c6570b4f34cb8fc29d87548dd587b27899afe11abc5766e
created: "2026-06-29"
id: SC-015
related_changes: []
related_reqs:
    - REQ-006
related_testcases: []
source: manual
status: validated
tags:
    - cobra
    - tools
    - error
    - connection
title: Tools commands exit non-zero with error when server unreachable
type: error-path
updated: "2026-06-29"
---

# SC-015: Tools commands exit non-zero with error when server unreachable

## Preconditions

- No MCP Server running at the configured address

## Steps

1. Run `eve-realm tools list`
2. Run `eve-realm tools invoke ping`

## Expected Result

- Both commands exit with a non-zero status code
- stderr contains a clear error message indicating the server is unreachable
- Error message includes the configured server address
- No raw gRPC status codes or stack traces in the output
- stdout is empty
