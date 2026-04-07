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
		newCmdGeminiInit(f),
		newSimpleCommand("nanobanana", "Gemini image generation preset", map[string]string{
			"canonical": "image generate --provider gemini --preset nano-banana",
		}),
	)

	return command
}

func newCmdGeminiInit(f *cmdutil.Factory) *cobra.Command {
	opts := &authInitOptions{
		Provider: "gemini",
	}

	command := &cobra.Command{
		Use:   "init",
		Short: "Initialize Gemini credentials",
		Annotations: map[string]string{
			"canonical": "auth init --provider gemini",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAuthInit(f, opts)
		},
	}

	command.Flags().StringVar(&opts.APIKey, "api-key", "", "api key")
	if err := command.MarkFlagRequired("api-key"); err != nil {
		panic(err)
	}

	return command
}
