package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/beduardo/eve-realm-cli/cmd/tools"
	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	GitHash   = "unknown"
	BuildDate = "unknown"
)

func defaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".eve-realm", "eve-realm.yaml")
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "eve-realm",
		Short:         "eve-realm — thin client for the Eve Realm platform",
		SilenceErrors: true,
	}

	configPath := defaultConfigPath()

	root.AddCommand(newVersionCmd())
	root.AddCommand(tools.NewToolsCmd(configPath))

	return root
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "eve-realm %s (git: %s, built: %s)\n", Version, GitHash, BuildDate)
		},
	}
}

func main() {
	root := newRootCmd()
	root.Execute() //nolint:errcheck
}
