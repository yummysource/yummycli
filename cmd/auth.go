package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yummysource/yummycli/internal/cmdutil"
)

// NewCmdAuth creates the shared auth command group.
func NewCmdAuth(f *cmdutil.Factory) *cobra.Command {
	command := newGroupCommand(
		"auth",
		"Provider and credential management",
		nil,
		newSimpleCommand("init", "Initialize provider credentials", nil),
		newSimpleCommand("list", "List configured provider credentials", nil),
		newSimpleCommand("status", "Show credential status", nil),
		newSimpleCommand("remove", "Remove provider credentials", nil),
		newSimpleCommand("clear", "Clear all provider credentials", nil),
	)

	return command
}
