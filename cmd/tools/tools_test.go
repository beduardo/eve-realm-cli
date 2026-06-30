package tools

import (
	"bytes"
	"context"
	"testing"

	"github.com/beduardo/eve-realm-cli/internal/mcpclient"
	"github.com/spf13/cobra"
)

// mockMCPClient is a test double for the MCPClient interface.
// Callers assign function fields to control behavior per test case.
type mockMCPClient struct {
	ListToolsFn  func(ctx context.Context) ([]mcpclient.Tool, error)
	InvokeToolFn func(ctx context.Context, name, input string) (string, error)
}

func (m *mockMCPClient) ListTools(ctx context.Context) ([]mcpclient.Tool, error) {
	return m.ListToolsFn(ctx)
}

func (m *mockMCPClient) InvokeTool(ctx context.Context, name, input string) (string, error) {
	return m.InvokeToolFn(ctx, name, input)
}

// runToolsCmd creates a root cobra.Command, attaches the tools subcommand tree
// with the injected mock client, captures stdout and stderr, executes the
// command with the given args, and returns the captured buffers and any error.
func runToolsCmd(t *testing.T, mock MCPClient, args ...string) (stdout, stderr *bytes.Buffer, err error) {
	t.Helper()

	root := &cobra.Command{
		Use:           "eve-realm",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	stdout = new(bytes.Buffer)
	stderr = new(bytes.Buffer)
	root.SetOut(stdout)
	root.SetErr(stderr)

	toolsCmd := newToolsCmdWithClient(mock)
	root.AddCommand(toolsCmd)

	root.SetArgs(args)
	err = root.Execute()
	return stdout, stderr, err
}

// TestMockMCPClient_SatisfiesInterface verifies at compile time that
// *mockMCPClient satisfies the local MCPClient interface.
func TestMockMCPClient_SatisfiesInterface(t *testing.T) {
	var _ MCPClient = (*mockMCPClient)(nil)
}

// TestNewToolsCmd_RegistersSubcommands verifies that NewToolsCmd returns a
// cobra.Command that has "list" and "invoke" registered as subcommands.
func TestNewToolsCmd_RegistersSubcommands(t *testing.T) {
	configPath := t.TempDir() + "/eve-realm.yaml"
	cmd := NewToolsCmd(configPath)

	subNames := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		subNames[sub.Name()] = true
	}

	for _, want := range []string{"list", "invoke"} {
		if !subNames[want] {
			t.Errorf("expected subcommand %q to be registered, got %v", want, subNames)
		}
	}
}

// TestNewToolsCmd_Use verifies the Use field of the tools command.
func TestNewToolsCmd_Use(t *testing.T) {
	configPath := t.TempDir() + "/eve-realm.yaml"
	cmd := NewToolsCmd(configPath)

	if cmd.Use != "tools" {
		t.Errorf("expected Use = %q, got %q", "tools", cmd.Use)
	}
}
