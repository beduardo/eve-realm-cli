package tools

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// newListCmd returns the "tools list" subcommand.
// It calls client.ListTools and writes each tool's name, description,
// and input schema to cmd.OutOrStdout() in a human-readable format.
// Errors are written to cmd.ErrOrStderr() and returned so Cobra exits non-zero.
func newListCmd(client MCPClient) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List tools available on the MCP Server",
		RunE: func(cmd *cobra.Command, args []string) error {
			tools, err := client.ListTools(context.Background())
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error: %s\n", err.Error())
				return err
			}

			w := cmd.OutOrStdout()
			for i, tool := range tools {
				if i > 0 {
					fmt.Fprintln(w)
				}
				fmt.Fprintf(w, "Name:         %s\n", tool.Name)
				fmt.Fprintf(w, "Description:  %s\n", tool.Description)
				fmt.Fprintf(w, "Input Schema: %s\n", tool.InputSchema)
			}
			return nil
		},
	}
}
