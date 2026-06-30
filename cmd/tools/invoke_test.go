package tools

import (
	"context"
	"strings"
	"testing"

	"github.com/beduardo/eve-realm-cli/internal/mcpclient"
)

// TestNewInvokeCmd_NoInputFlag verifies that when no --input flag is provided,
// InvokeTool is called with "{}" as the default input and the response is written to stdout.
func TestNewInvokeCmd_NoInputFlag(t *testing.T) {
	var capturedInput string
	mock := &mockMCPClient{
		InvokeToolFn: func(ctx context.Context, name, input string) (string, error) {
			capturedInput = input
			return `{"result":"pong"}`, nil
		},
	}

	stdout, stderr, err := runToolsCmd(t, mock, "tools", "invoke", "ping")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if capturedInput != "{}" {
		t.Errorf("expected input %q, got %q", "{}", capturedInput)
	}
	if !strings.Contains(stdout.String(), `{"result":"pong"}`) {
		t.Errorf("expected stdout to contain %q, got %q", `{"result":"pong"}`, stdout.String())
	}
	if stderr.String() != "" {
		t.Errorf("expected empty stderr, got %q", stderr.String())
	}
}

// TestNewInvokeCmd_WithInputFlag verifies that the --input flag value is passed
// verbatim to InvokeTool byte-for-byte.
func TestNewInvokeCmd_WithInputFlag(t *testing.T) {
	var capturedInput string
	mock := &mockMCPClient{
		InvokeToolFn: func(ctx context.Context, name, input string) (string, error) {
			capturedInput = input
			return `{"echo":"value"}`, nil
		},
	}

	stdout, stderr, err := runToolsCmd(t, mock, "tools", "invoke", "echo", "--input", `{"key":"value"}`)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if capturedInput != `{"key":"value"}` {
		t.Errorf("expected input %q, got %q", `{"key":"value"}`, capturedInput)
	}
	if !strings.Contains(stdout.String(), `{"echo":"value"}`) {
		t.Errorf("expected stdout to contain %q, got %q", `{"echo":"value"}`, stdout.String())
	}
	if stderr.String() != "" {
		t.Errorf("expected empty stderr, got %q", stderr.String())
	}
}

// TestNewInvokeCmd_ExplicitEmptyInput verifies that --input '{}' behaves identically
// to the default (no --input), passing "{}" to InvokeTool.
func TestNewInvokeCmd_ExplicitEmptyInput(t *testing.T) {
	var capturedInput string
	mock := &mockMCPClient{
		InvokeToolFn: func(ctx context.Context, name, input string) (string, error) {
			capturedInput = input
			return `{"ok":true}`, nil
		},
	}

	stdout, stderr, err := runToolsCmd(t, mock, "tools", "invoke", "echo", "--input", "{}")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if capturedInput != "{}" {
		t.Errorf("expected input %q, got %q", "{}", capturedInput)
	}
	if !strings.Contains(stdout.String(), `{"ok":true}`) {
		t.Errorf("expected stdout to contain %q, got %q", `{"ok":true}`, stdout.String())
	}
	if stderr.String() != "" {
		t.Errorf("expected empty stderr, got %q", stderr.String())
	}
}

// TestNewInvokeCmd_NestedJSONPassedVerbatim verifies that nested JSON in --input
// is passed to InvokeTool without re-serialization.
func TestNewInvokeCmd_NestedJSONPassedVerbatim(t *testing.T) {
	nestedInput := `{"outer":{"inner":"val","arr":[1,2,3]}}`
	var capturedInput string
	mock := &mockMCPClient{
		InvokeToolFn: func(ctx context.Context, name, input string) (string, error) {
			capturedInput = input
			return `{"done":true}`, nil
		},
	}

	_, _, err := runToolsCmd(t, mock, "tools", "invoke", "echo", "--input", nestedInput)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if capturedInput != nestedInput {
		t.Errorf("expected input %q, got %q", nestedInput, capturedInput)
	}
}

// TestNewInvokeCmd_ConnectionError verifies that a ConnectionError from InvokeTool
// causes a non-zero exit, routes the error message to stderr, and leaves stdout empty.
func TestNewInvokeCmd_ConnectionError(t *testing.T) {
	mock := &mockMCPClient{
		InvokeToolFn: func(ctx context.Context, name, input string) (string, error) {
			return "", &mcpclient.ConnectionError{Addr: "localhost:30051"}
		},
	}

	stdout, stderr, err := runToolsCmd(t, mock, "tools", "invoke", "ping")

	if err == nil {
		t.Fatal("expected non-nil error, got nil")
	}
	if stdout.String() != "" {
		t.Errorf("expected empty stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "localhost:30051") {
		t.Errorf("expected stderr to contain address, got %q", stderr.String())
	}
}

// TestNewInvokeCmd_ToolNotFoundWithAlternatives verifies that when InvokeTool
// returns a ToolNotFoundError, the command makes a secondary ListTools call and
// includes the available alternatives in stderr output.
func TestNewInvokeCmd_ToolNotFoundWithAlternatives(t *testing.T) {
	mock := &mockMCPClient{
		InvokeToolFn: func(ctx context.Context, name, input string) (string, error) {
			return "", &mcpclient.ToolNotFoundError{Name: "unknown-tool"}
		},
		ListToolsFn: func(ctx context.Context) ([]mcpclient.Tool, error) {
			return []mcpclient.Tool{
				{Name: "ping", Description: "ping tool"},
				{Name: "echo", Description: "echo tool"},
			}, nil
		},
	}

	stdout, stderr, err := runToolsCmd(t, mock, "tools", "invoke", "unknown-tool")

	if err == nil {
		t.Fatal("expected non-nil error, got nil")
	}
	if stdout.String() != "" {
		t.Errorf("expected empty stdout, got %q", stdout.String())
	}
	stderrStr := stderr.String()
	if !strings.Contains(stderrStr, "unknown-tool") {
		t.Errorf("expected stderr to contain %q, got %q", "unknown-tool", stderrStr)
	}
	if !strings.Contains(stderrStr, "ping") {
		t.Errorf("expected stderr to contain %q, got %q", "ping", stderrStr)
	}
	if !strings.Contains(stderrStr, "echo") {
		t.Errorf("expected stderr to contain %q, got %q", "echo", stderrStr)
	}
}

// TestNewInvokeCmd_ToolNotFoundSecondaryListFails verifies that when the secondary
// ListTools call also fails, the command still exits non-zero with the original
// not-found error in stderr, and does not panic.
func TestNewInvokeCmd_ToolNotFoundSecondaryListFails(t *testing.T) {
	mock := &mockMCPClient{
		InvokeToolFn: func(ctx context.Context, name, input string) (string, error) {
			return "", &mcpclient.ToolNotFoundError{Name: "unknown-tool"}
		},
		ListToolsFn: func(ctx context.Context) ([]mcpclient.Tool, error) {
			return nil, &mcpclient.ConnectionError{Addr: "localhost:30051"}
		},
	}

	stdout, stderr, err := runToolsCmd(t, mock, "tools", "invoke", "unknown-tool")

	if err == nil {
		t.Fatal("expected non-nil error, got nil")
	}
	if stdout.String() != "" {
		t.Errorf("expected empty stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "unknown-tool") {
		t.Errorf("expected stderr to contain %q, got %q", "unknown-tool", stderr.String())
	}
}

// TestNewInvokeCmd_ToolNotFoundZeroAlternatives verifies that when secondary ListTools
// returns an empty slice, the command exits non-zero with the not-found message and
// does not crash.
func TestNewInvokeCmd_ToolNotFoundZeroAlternatives(t *testing.T) {
	mock := &mockMCPClient{
		InvokeToolFn: func(ctx context.Context, name, input string) (string, error) {
			return "", &mcpclient.ToolNotFoundError{Name: "unknown-tool"}
		},
		ListToolsFn: func(ctx context.Context) ([]mcpclient.Tool, error) {
			return []mcpclient.Tool{}, nil
		},
	}

	stdout, stderr, err := runToolsCmd(t, mock, "tools", "invoke", "unknown-tool")

	if err == nil {
		t.Fatal("expected non-nil error, got nil")
	}
	if stdout.String() != "" {
		t.Errorf("expected empty stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "unknown-tool") {
		t.Errorf("expected stderr to contain %q, got %q", "unknown-tool", stderr.String())
	}
}

// TestNewInvokeCmd_ResponseWrittenToStdout verifies that a multi-field JSON response
// from InvokeTool is written exactly to stdout.
func TestNewInvokeCmd_ResponseWrittenToStdout(t *testing.T) {
	response := `{"status":"ok","count":42,"data":{"msg":"hello"}}`
	mock := &mockMCPClient{
		InvokeToolFn: func(ctx context.Context, name, input string) (string, error) {
			return response, nil
		},
	}

	stdout, _, err := runToolsCmd(t, mock, "tools", "invoke", "ping")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(stdout.String(), response) {
		t.Errorf("expected stdout to contain %q, got %q", response, stdout.String())
	}
}
