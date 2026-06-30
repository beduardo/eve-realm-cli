---
content_hash: 44176c3c88ad62acfaef3c8c372091959c6175dccfbfa792dc1b3708bf1fc33f
created: "2026-06-28"
id: SC-00F
related_changes: []
related_reqs:
    - REQ-007
related_testcases: []
source: manual
status: validated
tags:
    - marketplace
    - list
title: List displays all embedded plugins and skills with descriptions
type: happy-path
updated: "2026-06-29"
---

# SC-00F: List displays all embedded plugins and skills with descriptions

## Preconditions

- The binary is compiled with the embedded marketplace files
- At least one plugin (`eve-realm`) with at least one skill (`master`) is embedded

## Steps

1. Run `eve-realm marketplace list`

## Expected Result

- Output lists the `eve-realm` plugin with its version and description
- Under the plugin, the `master` skill is listed with the description extracted from `SKILL.md` YAML frontmatter
- Skills named `workflow` are excluded from the listing
- Output is sorted alphabetically by plugin name, then by skill name
- Command exits with zero status code
