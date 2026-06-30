---
content_hash: 943f43bb924192c36bc263ca851422f092248ec8b3534cf834e8a85f5c053f29
created: "2026-06-27"
id: REQ-005
priority: high
related_adrs: []
related_changes: []
related_scenarios:
    - SC-007
    - SC-008
    - SC-009
    - SC-00A
    - SC-00B
related_testcases: []
related_userstories: []
source: manual
status: implemented
tags:
    - skill
    - marketplace
    - master-skill
    - mcp
title: Master skill for MCP tool discovery and invocation
updated: "2026-06-29"
---

# REQ-005: Master skill for MCP tool discovery and invocation

## Description

A marketplace skill registered in the CLI that acts as the gateway between Claude Code
(or any AI tool) and the MCP Server's tool registry. The skill is registered as a single
entry in the marketplace, but dynamically surfaces all tools available in the MCP Server.

The skill operates in two modes:

- **Discovery mode** — When invoked without a specific tool request, the skill calls
  `eve-realm tools list` (REQ-006) and constructs a descriptive prompt listing all
  available tools with their names, descriptions, and input schemas. This prompt helps
  the AI choose which tool to invoke. The tool list is live — it reflects the current state
  of the MCP Server's registry (as plugins connect/disconnect, the list changes).

- **Invocation mode** — When invoked with a specific tool name and input (or when the AI
  selects a tool from the discovery prompt), the skill calls
  `eve-realm tools invoke <name> --input <json>` (REQ-006) and returns the tool's
  response to the AI.

The skill is also flexible enough to receive a complete prompt for direct execution,
acting as a pass-through agent that determines the appropriate tool and invokes it.

## Acceptance Criteria

1. A skill is registered in the CLI marketplace under the name `eve-realm` (or configured name).
2. When the skill is invoked without arguments, it calls `eve-realm tools list` and returns a formatted prompt listing all available tools with name, description, and input schema.
3. When the skill is invoked with a tool name and input, it calls `eve-realm tools invoke <name>` and returns the tool's JSON response.
4. When the MCP Server is unreachable, the skill surfaces the error from the Cobra command as a user-friendly message.
5. When a requested tool is not found, the skill surfaces the Cobra command's error which lists available tools as alternatives.
6. The skill code lives in `skills/master/`.
7. The skill delegates MCP Server address resolution to the Cobra command (REQ-006 handles config/env/default).
8. End-to-end: invoking the skill with tool name `ping` returns the MCP Server's ping response (`{"message": "pong", "timestamp": "..."}`).
