package tools

import (
	"context"

	"github.com/beduardo/eve-realm-cli/internal/config"
	"github.com/beduardo/eve-realm-cli/internal/mcpclient"
	"github.com/spf13/cobra"
)

// MCPClient is the local interface that cmd/tools depends on.
// It matches the subset of mcpclient.Client methods used by the tools commands.
type MCPClient interface {
	ListTools(ctx context.Context) ([]mcpclient.Tool, error)
	InvokeTool(ctx context.Context, name, input string) (string, error)
}

// NewToolsCmd builds the "tools" cobra.Command for production use.
// It resolves the MCP Server address from the config file at configPath,
// the EVE_REALM_MCP_ADDR environment variable, or the default address, in
// that order of precedence.
func NewToolsCmd(configPath string) *cobra.Command {
	cfg, _ := config.LoadHostConfig(configPath)
	addr := config.Resolve("EVE_REALM_MCP_ADDR", cfg.MCPServerAddr)
	if addr == "" {
		addr = mcpclient.DefaultMCPAddr
	}
	client := mcpclient.NewClient(addr)
	return newToolsCmdWithClient(client)
}

// newToolsCmdWithClient constructs the "tools" command tree using the supplied
// client. This internal constructor is used by tests to inject a mock.
func newToolsCmdWithClient(client MCPClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tools",
		Short: "Interact with MCP Server tools",
	}

	cmd.AddCommand(newListCmd(client))
	cmd.AddCommand(newInvokeCmd(client))

	return cmd
}

