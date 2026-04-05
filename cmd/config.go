package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yummysource/yummycli/internal/cmdutil"
)

// NewCmdConfig creates the shared config command group.
func NewCmdConfig(f *cmdutil.Factory) *cobra.Command {
	command := newGroupCommand(
		"config",
		"Local configuration management",
		nil,
		newSimpleCommand("list", "List configuration values", nil),
		newSimpleCommand("get", "Get a configuration value", nil),
		newSimpleCommand("set", "Set a configuration value", nil),
	)

	return command
}
