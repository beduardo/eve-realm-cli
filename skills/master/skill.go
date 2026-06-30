package master

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/beduardo/eve-realm-cli/internal/mcpclient"
)

// MCPClient abstracts the gRPC client for dependency injection and testing.
type MCPClient interface {
	ListTools(ctx context.Context) ([]mcpclient.Tool, error)
	InvokeTool(ctx context.Context, name, input string) (string, error)
}

// Run executes the master skill.
//   - No args: discovery mode — lists all tools from the MCP Server.
//   - Args: invocation mode — invokes the named tool with optional JSON input.
func Run(client MCPClient, args []string) string {
	if len(args) == 0 {
		return runDiscovery(client)
	}
	return runInvocation(client, args)
}

// runDiscovery calls ListTools and formats the result for AI consumption.
func runDiscovery(client MCPClient) string {
	tools, err := client.ListTools(context.Background())
	if err != nil {
		return formatError(client, err)
	}

	if len(tools) == 0 {
		return "No tools are available on the MCP Server."
	}

	var sb strings.Builder
	sb.WriteString("Available tools:\n\n")
	for _, t := range tools {
		fmt.Fprintf(&sb, "Tool: %s\n", t.Name)
		fmt.Fprintf(&sb, "Description: %s\n", t.Description)
		fmt.Fprintf(&sb, "Input schema: %s\n\n", t.InputSchema)
	}
	return strings.TrimRight(sb.String(), "\n")
}

// runInvocation calls InvokeTool with the supplied name and input.
func runInvocation(client MCPClient, args []string) string {
	toolName := args[0]
	input := "{}"
	if len(args) >= 2 {
		input = args[1]
	}

	output, err := client.InvokeTool(context.Background(), toolName, input)
	if err != nil {
		return formatError(client, err)
	}
	return output
}

// formatError maps typed errors to user-friendly messages.
// It uses client to fetch alternatives when a ToolNotFoundError is detected.
func formatError(client MCPClient, err error) string {
	var connErr *mcpclient.ConnectionError
	if errors.As(err, &connErr) {
		return fmt.Sprintf(
			"MCP Server is not available at %s. Check that the server is running.",
			connErr.Addr,
		)
	}

	var notFoundErr *mcpclient.ToolNotFoundError
	if errors.As(err, &notFoundErr) {
		tools, listErr := client.ListTools(context.Background())
		if listErr != nil {
			return fmt.Sprintf("Tool %q was not found on the MCP Server.", notFoundErr.Name)
		}
		names := make([]string, 0, len(tools))
		for _, t := range tools {
			names = append(names, t.Name)
		}
		return fmt.Sprintf(
			"Tool %q was not found on the MCP Server. Available tools: %s",
			notFoundErr.Name,
			strings.Join(names, ", "),
		)
	}

	return fmt.Sprintf("Unexpected error: %v", err)
}
