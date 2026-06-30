# Feasibility Brief

**Sprint**: SP-001
**Project Root**: /Users/bruno/repo-pessoal/eve-realm/eve-realm-cli/main
**Target Entity**: REQ-005

## Entity Summary
Master skill for MCP tool discovery and invocation — marketplace skill at `skills/master/` that acts as gateway between AI tools and MCP Server. Two modes: discovery (ListTools) and invocation (InvokeTool). Depends on REQ-004's gRPC client. The gRPC client defaults to `localhost:30051` (k3d NodePort), but the skill reads the address from CLI config or env var.

## Sprint Context
- Current entity count: 11
- Scope score: 11/5
- Other entities in sprint: REQ-004, SC-003, SC-004, SC-005, SC-006, SC-007, SC-008, SC-009, SC-00A, SC-00B

## Focus Questions
- Is there an existing marketplace/skill registration system in the codebase?
- How do skills integrate with the CLI command structure?
- What is the skill interface contract (input/output format)?
- Does the skill need to be registered as a Claude Code MCP tool, a CLI command, or both?
- How does the skill read config (MCP Server address) — viper, direct YAML, env var?
- What is the dependency chain: skill -> mcpclient -> gRPC -> MCP Server?
- Are there patterns in eve-cli for skill implementation to reference?
- How does the default address (localhost:30051) interact with the config/env var resolution?
