package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yummysource/yummycli/internal/auth"
	"github.com/yummysource/yummycli/internal/cmdutil"
	internalimage "github.com/yummysource/yummycli/internal/image"
	"github.com/yummysource/yummycli/internal/providers"
)

// imageGenerateOptions is the canonical options struct for image generation.
// Used by both `image generate` (capability layer) and `gemini nanobanana`
// (provider shortcut) so that the shortcut is a thin wrapper, not a copy.
type imageGenerateOptions struct {
	Provider    string
	Prompt      string
	Output      string
	Model       string
	AspectRatio string
	ImageSize   string
	Quality     string
	Style       string
	InputImages []string
}

// imageGenerateResult is the JSON output written on a successful generation.
type imageGenerateResult struct {
	Provider        string `json:"provider"`
	Output          string `json:"output"`
	Model           string `json:"model"`
	InputImageCount int    `json:"inputImageCount"`
}

// NewCmdImage creates the provider-agnostic image command group.
func NewCmdImage(f *cmdutil.Factory) *cobra.Command {
	command := &cobra.Command{
		Use:   "image",
		Short: "Provider-agnostic image capabilities",
	}

	command.AddCommand(
		newCmdImageGenerate(f),
	)

	return command
}

func newCmdImageGenerate(f *cmdutil.Factory) *cobra.Command {
	opts := &imageGenerateOptions{}

	command := &cobra.Command{
		Use:   "generate",
		Short: "Generate an image",
		Annotations: map[string]string{
			"capability": "image.generate",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runImageGenerate(f, opts)
		},
	}

	command.Flags().StringVar(&opts.Provider, "provider", "", "model provider (e.g. gemini)")
	command.Flags().StringVar(&opts.Prompt, "prompt", "", "image generation prompt")
	command.Flags().StringVar(&opts.Output, "output", "", "output image path")
	command.Flags().StringVar(&opts.Model, "model", "", "model name")
	command.Flags().StringVar(&opts.AspectRatio, "aspect-ratio", "", "image aspect ratio")
	command.Flags().StringVar(&opts.ImageSize, "image-size", "", "image size")
	command.Flags().StringVar(&opts.Quality, "quality", "", "image quality (openai: standard or hd)")
	command.Flags().StringVar(&opts.Style, "style", "", "image style (openai: vivid or natural)")
	command.Flags().StringArrayVar(&opts.InputImages, "input-image", nil, "input image path (repeatable)")

	if err := command.MarkFlagRequired("prompt"); err != nil {
		panic(err)
	}

	return command
}

// runImageGenerate is the canonical implementation shared by `image generate`
// and provider shortcuts such as `gemini nanobanana`.
func runImageGenerate(f *cmdutil.Factory, opts *imageGenerateOptions) error {
	if f.ImageGenerator == nil {
		return fmt.Errorf("image generator is not configured")
	}
	if f.Output == nil {
		return fmt.Errorf("output writer is not configured")
	}

	// Resolve provider from config when not explicitly provided.
	if opts.Provider == "" {
		defaultProvider, err := f.CredentialStore.GetDefaultProvider()
		if err != nil {
			return err
		}
		if defaultProvider == "" {
			return fmt.Errorf("no provider specified and no default configured; run: yummycli init --provider <name> --api-key <key> --default")
		}
		opts.Provider = defaultProvider
	}

	// Apply provider-specific defaults and validation.
	switch opts.Provider {
	case providers.Gemini:
		if opts.Model == "" {
			opts.Model = geminiDefaultModel
		}
		if opts.AspectRatio == "" {
			opts.AspectRatio = "16:9"
		}
		if opts.ImageSize == "" {
			opts.ImageSize = "1K"
		}
		if err := validateAspectRatio(opts.Model, opts.AspectRatio); err != nil {
			return err
		}
		normalized, err := validateImageSize(opts.Model, opts.ImageSize)
		if err != nil {
			return err
		}
		opts.ImageSize = normalized
	case providers.OpenAI:
		if opts.Model == "" {
			opts.Model = openAIDefaultModel
		}
		if opts.Model != openAIDefaultModel {
			return fmt.Errorf("unsupported openai model: %s (only dall-e-3 is supported)", opts.Model)
		}
		if opts.ImageSize == "" {
			opts.ImageSize = "1024x1024"
		}
		if opts.Quality == "" {
			opts.Quality = "standard"
		}
		if opts.Style == "" {
			opts.Style = "vivid"
		}
		if err := validateOpenAISize(opts.ImageSize); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported provider: %s", opts.Provider)
	}

	if opts.Output == "" {
		opts.Output = defaultImageOutputPath(opts.Provider)
	}

	fallback, err := resolveProviderFallback(f.CredentialStore, opts.Provider)
	if err != nil {
		return err
	}

	req := internalimage.GenerateImageRequest{
		Provider:    opts.Provider,
		Prompt:      opts.Prompt,
		Output:      opts.Output,
		Model:       opts.Model,
		AspectRatio: opts.AspectRatio,
		ImageSize:   opts.ImageSize,
		Quality:     opts.Quality,
		Style:       opts.Style,
		InputImages: opts.InputImages,
		Fallback:    fallback,
	}

	if err := f.ImageGenerator.GenerateImage(context.Background(), req); err != nil {
		return err
	}

	result := imageGenerateResult{
		Provider:        opts.Provider,
		Output:          opts.Output,
		Model:           opts.Model,
		InputImageCount: len(opts.InputImages),
	}

	return f.Output.JSON(result)
}

const openAIDefaultModel = "dall-e-3"

// validateOpenAISize checks that the size is valid for dall-e-3.
func validateOpenAISize(size string) error {
	valid := []string{"1024x1024", "1024x1792", "1792x1024"}
	for _, v := range valid {
		if v == size {
			return nil
		}
	}
	return fmt.Errorf("unsupported image size for openai: %s", size)
}

// resolveProviderFallback returns the first configured provider that is not the primary.
func resolveProviderFallback(credStore *auth.ProviderCredentialStore, primary string) (string, error) {
	all := providers.All()
	for _, p := range all {
		if p == primary {
			continue
		}
		configured, err := credStore.HasAPIKey(p)
		if err != nil {
			return "", err
		}
		if configured {
			return p, nil
		}
	}
	return "", nil
}
