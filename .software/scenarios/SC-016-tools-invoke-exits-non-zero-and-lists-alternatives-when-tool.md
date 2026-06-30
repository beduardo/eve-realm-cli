---
content_hash: a519a355cadb295c9fc0deb866c64d07a97a116c39f35891c9218d554fadf75f
created: "2026-06-29"
id: SC-016
related_changes: []
related_reqs:
    - REQ-006
related_testcases: []
source: manual
status: validated
tags:
    - cobra
    - tools-invoke
    - error
    - not-found
title: Tools invoke exits non-zero and lists alternatives when tool not found
type: error-path
updated: "2026-06-29"
---

# SC-016: Tools invoke exits non-zero and lists alternatives when tool not found

## Preconditions

- MCP Server running with multiple tools registered (e.g., `ping`, `echo`)
- MCP Server address is resolvable

## Steps

1. Run `eve-realm tools invoke nonexistent`

## Expected Result

- Command exits with a non-zero status code
- stderr states that the tool `"nonexistent"` was not found
- stderr lists available tools as alternatives (e.g., `ping`, `echo`)
- stdout is empty
