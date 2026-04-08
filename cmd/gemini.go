package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yummysource/yummycli/internal/cmdutil"
	internalimage "github.com/yummysource/yummycli/internal/image"
	"github.com/yummysource/yummycli/internal/providers"
)

// NewCmdGemini creates the Gemini provider shortcut command.
func NewCmdGemini(f *cmdutil.Factory) *cobra.Command {
	command := &cobra.Command{
		Use:   "gemini",
		Short: "Gemini provider shortcuts and presets",
	}

	command.AddCommand(
		newCmdGeminiInit(f),
		newCmdGeminiNanoBanana(f),
	)

	return command
}

func newCmdGeminiInit(f *cmdutil.Factory) *cobra.Command {
	opts := &authInitOptions{
		Provider: providers.Gemini,
	}

	command := &cobra.Command{
		Use:   "init",
		Short: "Initialize Gemini credentials",
		Annotations: map[string]string{
			"canonical": "auth init --provider " + providers.Gemini,
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

type geminiNanoBananaOptions struct {
	Prompt string
	Output string
	Model  string
}

func newCmdGeminiNanoBanana(f *cmdutil.Factory) *cobra.Command {
	opts := &geminiNanoBananaOptions{
		Model: "gemini-3.1-flash-image-preview",
	}

	command := &cobra.Command{
		Use:   "nanobanana",
		Short: "Generate images with Gemini Nano Banana",
		Annotations: map[string]string{
			"canonical": "image generate --provider " + providers.Gemini + " --preset nano-banana",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGeminiNanoBanana(f, opts)
		},
	}

	command.Flags().StringVar(&opts.Prompt, "prompt", "", "Image generation prompt")
	command.Flags().StringVar(&opts.Output, "output", "", "Output image path")
	command.Flags().StringVar(&opts.Model, "model", opts.Model, "Gemini image model")

	if err := command.MarkFlagRequired("prompt"); err != nil {
		panic(err)
	}
	if err := command.MarkFlagRequired("output"); err != nil {
		panic(err)
	}

	return command
}

func runGeminiNanoBanana(f *cmdutil.Factory, opts *geminiNanoBananaOptions) error {
	if f.ImageGenerator == nil {
		return fmt.Errorf("image generator is not configured")
	}

	req := internalimage.GenerateImageRequest{
		Provider: providers.Gemini,
		Prompt:   opts.Prompt,
		Output:   opts.Output,
		Model:    opts.Model,
	}

	return f.ImageGenerator.GenerateImage(context.Background(), req)
}
