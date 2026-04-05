package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yummysource/yummycli/internal/cmdutil"
)

// NewCmdGemini creates the Gemini provider shortcut command.
func NewCmdGemini(f *cmdutil.Factory) *cobra.Command {
	command := &cobra.Command{
		Use:   "gemini",
		Short: "Gemini provider shortcuts and presets",
	}

	command.AddCommand(
		newSimpleCommand("init", "Initialize Gemini credentials", map[string]string{
			"canonical": "auth init --provider gemini",
		}),
		newSimpleCommand("nanobanana", "Gemini image generation preset", map[string]string{
			"canonical": "image generate --provider gemini --preset nano-banana",
		}),
	)

	return command
}
