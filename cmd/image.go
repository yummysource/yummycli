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
	Provider     string
	Prompt       string
	Output       string
	Model        string
	AspectRatio  string
	ImageSize    string
	Quality      string
	Style        string
	OutputFormat string
	InputImages  []string
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
		Short: "Generate or edit an image",
		Long: `Generate or edit an image using Gemini or OpenAI.

Provider is resolved from config when --provider is omitted.
Run "yummycli init --provider <name> --api-key <key> --default" to set a default provider.

PROVIDERS
  gemini   Google Gemini image generation (gemini-3.1-flash-image by default)
           Docs: https://ai.google.dev/gemini-api/docs/image-generation
  openai   OpenAI image generation (gpt-image-2 by default)
           Docs: https://developers.openai.com/api/docs/guides/image-generation

GEMINI FLAGS
  --aspect-ratio   Output aspect ratio.
                   Supported: 1:1 1:4 1:8 2:3 3:2 3:4 4:1 4:3 4:5 5:4 8:1 9:16 16:9 21:9
                   Default: 16:9
  --image-size     Output resolution.
                   Supported: 512 0.5K 1K 2K 4K
                   Default: 1K
                   Note: 512 and 0.5K are gemini-3.1-flash-image only.

OPENAI FLAGS
  --image-size     Output dimensions (W x H, both must be multiples of 16).
                   Presets: 1536x864 (16:9, default) 1024x576 2048x1152
                            1024x1024 (square) 1536x1024 1024x1536
  --quality        Rendering quality: low | medium | high | auto
                   Default: API auto
  --output-format  Output file format: png | jpeg | webp
                   Default: png
                   Note: jpeg significantly reduces file size vs png.

MODELS
  gemini   gemini-3.1-flash-image (default), gemini-3-pro-image-preview
  openai   gpt-image-2 (default), gpt-5.5
           Unknown models trigger a warning and fall back to gpt-image-2.

IMAGE EDITING
  Pass one or more --input-image flags to edit existing images instead of generating from scratch.
  Gemini: prompt + input images → multimodal edit (up to many images)
  OpenAI: prompt + input images → /images/edits endpoint (multipart upload)`,
		Example: `  # Text-to-image (uses configured default provider)
  yummycli image generate --prompt "a sunset over the ocean"

  # Gemini with explicit aspect ratio
  yummycli image generate --provider gemini --prompt "city skyline" --aspect-ratio 21:9 --image-size 4K

  # OpenAI 16:9 JPEG
  yummycli image generate --provider openai --prompt "a rabbit" --output-format jpeg

  # Edit an existing image (Gemini)
  yummycli image generate --provider gemini --prompt "make it watercolor" --input-image ./photo.png

  # Multi-image compositing (Gemini)
  yummycli image generate --provider gemini --prompt "blend these" \
    --input-image ./a.png --input-image ./b.png`,
		Annotations: map[string]string{
			"capability": "image.generate",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runImageGenerate(f, opts)
		},
	}

	command.Flags().StringVar(&opts.Provider, "provider", "", "provider: gemini or openai (uses configured default if omitted)")
	command.Flags().StringVar(&opts.Prompt, "prompt", "", "image generation or editing prompt")
	command.Flags().StringVar(&opts.Output, "output", "", "output file path (auto-generated if omitted)")
	command.Flags().StringVar(&opts.Model, "model", "", "model name (see MODELS in --help)")
	command.Flags().StringVar(&opts.AspectRatio, "aspect-ratio", "", "gemini: output aspect ratio, e.g. 16:9 (default), 9:16, 1:1")
	command.Flags().StringVar(&opts.ImageSize, "image-size", "", "gemini: 1K/2K/4K (default 1K) | openai: WxH e.g. 1536x864 (default)")
	command.Flags().StringVar(&opts.Quality, "quality", "", "openai: rendering quality — low | medium | high | auto")
	command.Flags().StringVar(&opts.OutputFormat, "output-format", "", "openai: output format — png (default) | jpeg | webp")
	command.Flags().StringArrayVar(&opts.InputImages, "input-image", nil, "input image path for editing; repeat for multiple (e.g. --input-image a.png --input-image b.png)")

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
		if !isKnownOpenAIModel(opts.Model) {
			fmt.Fprintf(f.IOStreams.ErrOut, "warning: unknown openai model %q, using default %q\n", opts.Model, openAIDefaultModel)
			opts.Model = openAIDefaultModel
		}
		if opts.ImageSize == "" {
			opts.ImageSize = "1536x864"
		}
		if err := validateOpenAISize(opts.ImageSize); err != nil {
			return err
		}
		if opts.OutputFormat != "" {
			if err := validateOpenAIOutputFormat(opts.OutputFormat); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unsupported provider: %s", opts.Provider)
	}

	if opts.Output == "" {
		opts.Output = defaultImageOutputPath(opts.Provider, opts.OutputFormat)
	}

	fallback, err := resolveProviderFallback(f.CredentialStore, opts.Provider)
	if err != nil {
		return err
	}

	req := internalimage.GenerateImageRequest{
		Provider:     opts.Provider,
		Prompt:       opts.Prompt,
		Output:       opts.Output,
		Model:        opts.Model,
		AspectRatio:  opts.AspectRatio,
		ImageSize:    opts.ImageSize,
		Quality:      opts.Quality,
		Style:        opts.Style,
		OutputFormat: opts.OutputFormat,
		InputImages:  opts.InputImages,
		Fallback:     fallback,
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

const openAIDefaultModel = "gpt-image-2"

// knownOpenAIModels lists the OpenAI image models supported by yummycli.
var knownOpenAIModels = []string{"gpt-image-2", "gpt-5.5"}

// isKnownOpenAIModel reports whether model is a recognised OpenAI image model.
func isKnownOpenAIModel(model string) bool {
	for _, m := range knownOpenAIModels {
		if m == model {
			return true
		}
	}
	return false
}

// validateOpenAIOutputFormat checks that the output format is valid for OpenAI.
func validateOpenAIOutputFormat(format string) error {
	switch format {
	case "png", "jpeg", "webp":
		return nil
	}
	return fmt.Errorf("unsupported output format for openai: %s (supported: png, jpeg, webp)", format)
}

// validateOpenAISize checks that the size is a known OpenAI image size.
func validateOpenAISize(size string) error {
	valid := []string{
		"1536x864",            // 16:9 default
		"1024x576",            // 16:9 smaller
		"2048x1152",           // 16:9 larger
		"1024x1024",           // square
		"1536x1024",           // 3:2 landscape
		"1024x1536",           // 2:3 portrait
		"1024x1792", "1792x1024", // dall-e-3
	}
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
