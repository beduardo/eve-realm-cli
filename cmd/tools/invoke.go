package tools

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/beduardo/eve-realm-cli/internal/mcpclient"
	"github.com/spf13/cobra"
)

// newInvokeCmd returns the "tools invoke" subcommand.
// It accepts a positional tool-name argument and an optional --input flag (default "{}").
// The flag value is passed verbatim to client.InvokeTool without re-serialization.
// On success the JSON response is written to cmd.OutOrStdout().
// On ToolNotFoundError a secondary ListTools call fetches alternatives for the error message.
// All errors are written to cmd.ErrOrStderr() and returned so Cobra exits non-zero.
func newInvokeCmd(client MCPClient) *cobra.Command {
	var input string

	cmd := &cobra.Command{
		Use:   "invoke <tool-name>",
		Short: "Invoke a tool on the MCP Server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			toolName := args[0]

			result, err := client.InvokeTool(context.Background(), toolName, input)
			if err != nil {
				var notFoundErr *mcpclient.ToolNotFoundError
				if errors.As(err, &notFoundErr) {
					return handleToolNotFound(cmd, client, err)
				}
				fmt.Fprintf(cmd.ErrOrStderr(), "Error: %s\n", err.Error())
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), result)
			return nil
		},
	}

	cmd.Flags().StringVar(&input, "input", "{}", "JSON input for the tool (default: empty object)")

	return cmd
}

// handleToolNotFound writes the not-found error to stderr and optionally appends
// available alternatives gathered from a secondary ListTools call.
func handleToolNotFound(cmd *cobra.Command, client MCPClient, notFoundErr error) error {
	w := cmd.ErrOrStderr()
	fmt.Fprintf(w, "Error: %s\n", notFoundErr.Error())

	tools, listErr := client.ListTools(context.Background())
	if listErr != nil || len(tools) == 0 {
		return notFoundErr
	}

	names := make([]string, 0, len(tools))
	for _, t := range tools {
		names = append(names, t.Name)
	}
	fmt.Fprintf(w, "Available tools: %s\n", strings.Join(names, ", "))

	return notFoundErr
}
