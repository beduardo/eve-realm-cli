---
content_hash: 23c5ed6518429ce12d7e777f8d014f43ff2b14cbbca80a41d33b7dca3be5f4e6
created: "2026-06-20"
id: REQ-001
priority: high
related_adrs: []
related_changes: []
related_scenarios: []
related_testcases: []
related_userstories: []
source: manual
status: blocked
tags:
    - tdd
    - testing
    - cross-cutting
testing_strategy: tdd
title: Test-Driven Development Strategy
updated: "2026-06-22"
---

# REQ-001: Test-Driven Development Strategy

## Strategy

This project uses **TDD** (Test-Driven Development). All implementation steps must follow
the red → green → refactor cycle.

**TDD cycle per acceptance criterion:**

1. **Red** — Write a failing test that encodes the acceptance criterion's expected
   behavior. Run it to confirm failure.
2. **Green** — Write the minimum production code to make the test pass. No speculative
   code.
3. **Refactor** — Remove duplication, improve naming, simplify structure. All tests
   must remain green.

Each step in a sprint plan addresses one or more acceptance criteria. For each
criterion, the implementer executes one complete red → green → refactor micro-cycle
before moving to the next.

## Scope

- **Sprint entities**: All code implementing sprint requirements, scenarios, bugfixes,
  and user stories follows TDD.
- **ADRs are excluded**: Architecture decisions are not implementable units and do not
  receive test expectations.

## Framework

- **Language**: Go
- **Test framework**: `testing` (standard library only — no testify or external
  assertion libraries)
- **Test runner**: `go test ./...`
- **Coverage**: `go test -coverprofile=coverage.out ./...`
- **Assertion style**: Standard `if got != want { t.Errorf(...) }` pattern

## Patterns

- **Table-driven tests**: Use `[]struct{ name string; input T; want U }` with
  `t.Run(tc.name, ...)` subtests for functions with multiple scenarios.
- **Temp directories**: Always use `t.TempDir()` for file and directory operations —
  never write to the real filesystem.
- **Interface-based mocking**: Define small interfaces at the consumer site and pass
  test doubles. No external mock generation libraries.
- **Process substitution for exec**: When testing code that calls `os/exec`, use
  `exec.CommandContext` injection or the `TestHelperProcess` pattern
  (`-test.run=TestHelperProcess`).
- **HTTP tests**: Use `httptest.NewServer()` or `httptest.NewRecorder()` for HTTP
  client and handler tests.
- **YAML round-trip tests**: For config parsing, write YAML to `t.TempDir()`, load
  it, assert fields, serialize back, and compare.
- **Cobra command tests**: Use `cmd.SetArgs()` + `cmd.Execute()` with captured
  stdout/stderr via `bytes.Buffer`.

## Constraints

- Tests must live in `*_test.go` files alongside their source files in the same
  package (white-box testing) or in a `_test` package suffix (black-box testing)
  depending on what is being tested.
- No global test state — each test must be self-contained.
- Test function naming follows `TestFunctionName_Scenario` convention
  (e.g., `TestDiscoverPlugins_EmptyDirectory`).
- No test should depend on external services or network access — all I/O is mocked
  or uses temp directories.
- No external dependencies for testing — standard library only.

## Test Expectations (pipeline integration)

This requirement flows to all sprint agents via the cross-cutting catalog (pinned).
When agents load it:

- **Spec writer**: Generates a "Test Expectations" subsection per entity, mapping each
  acceptance criterion to the tests that must verify it, the mocking strategy, and the
  test type (unit/integration).
- **Plan generator**: Propagates test expectations per step. Each step's acceptance
  criteria include test-to-AC traceability (e.g., "Test verifies AC-1 from REQ-003").
  Step ordering co-locates test files with production files (TDD = test-first).
- **Step implementer**: Follows the red → green → refactor cycle as described in this
  entity's Strategy section. Reads test expectations from the step brief and writes
  tests covering all listed expectations before production code.
- **Step verifier**: Compares written tests against the step's test expectations. For
  each expectation, verifies at least one corresponding test exists with meaningful
  assertions. Missing test coverage is a **step failure** — triggers implementer retry,
  then abort. Hollow tests (trivial assertions, implementation-detail testing) are
  flagged.

## Acceptance Criteria

- Given a sprint that includes this entity (or has it in the cross-cutting catalog), when the spec writer generates SPEC.md, then each entity's Implementation Section includes a "Test Expectations" subsection.
- Given a SPEC.md with test expectations, when the plan generator produces PLAN.md, then each implementation step includes test expectations and acceptance criteria that reference test-to-AC traceability.
- Given a plan step with test expectations and TDD strategy, when the step implementer executes, then tests are written before production code following the red → green → refactor cycle.
- Given a completed step with test expectations, when the step verifier runs, then it validates that written tests cover all listed expectations — and fails the step if any expectation lacks a corresponding test.
