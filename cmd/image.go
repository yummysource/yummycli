package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yummysource/yummycli/internal/cmdutil"
)

// NewCmdImage creates the provider-agnostic image command group.
func NewCmdImage(f *cmdutil.Factory) *cobra.Command {
	command := &cobra.Command{
		Use:   "image",
		Short: "Provider-agnostic image capabilities",
	}

	command.AddCommand(
		newSimpleCommand("generate", "Generate an image", map[string]string{
			"capability": "image.generate",
		}),
		newSimpleCommand("edit", "Edit an existing image", map[string]string{
			"capability": "image.edit",
		}),
		newSimpleCommand("models", "List available image models", map[string]string{
			"capability": "image.models",
		}),
	)

	return command
}
