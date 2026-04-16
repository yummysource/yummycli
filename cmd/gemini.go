package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/yummysource/yummycli/internal/cmdutil"
	"github.com/yummysource/yummycli/internal/providers"
)

// geminiDefaultModel is the default image generation model for gemini nanobanana.
const geminiDefaultModel = "gemini-3.1-flash-image-preview"

// NewCmdGemini creates the Gemini provider shortcut command group.
func NewCmdGemini(f *cmdutil.Factory) *cobra.Command {
	command := &cobra.Command{
		Use:   "gemini",
		Short: "Gemini provider shortcuts and presets",
	}

	command.AddCommand(
		newCmdGeminiInit(f),
		newCmdGeminiNanoBanana(f),
		newCmdGeminiVeo(f),
	)

	return command
}

// newCmdGeminiInit creates a shortcut for initializing Gemini credentials,
// equivalent to auth init --provider gemini.
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

// newCmdGeminiNanoBanana is a user-friendly shortcut for image generate --provider gemini,
// with Gemini-specific defaults preset for model, aspect ratio, and image size.
func newCmdGeminiNanoBanana(f *cmdutil.Factory) *cobra.Command {
	opts := &imageGenerateOptions{
		Provider:    providers.Gemini,
		Model:       geminiDefaultModel,
		AspectRatio: "16:9",
		ImageSize:   "1K",
	}

	command := &cobra.Command{
		Use:   "nanobanana",
		Short: "Generate images with Gemini Nano Banana",
		Annotations: map[string]string{
			"canonical": "image generate --provider " + providers.Gemini + " --preset nano-banana",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runImageGenerate(f, opts)
		},
	}

	command.Flags().StringVar(&opts.Prompt, "prompt", "", "Image generation prompt")
	command.Flags().StringVar(&opts.Output, "output", "", "Output image path")
	command.Flags().StringVar(&opts.Model, "model", opts.Model, "Gemini image model")
	command.Flags().StringVar(&opts.AspectRatio, "aspect-ratio", opts.AspectRatio, "Gemini image aspect-ratio")
	command.Flags().StringVar(&opts.ImageSize, "image-size", opts.ImageSize, "Gemini image size")
	command.Flags().StringArrayVar(&opts.InputImages, "input-image", nil, "Gemini input image path")

	if err := command.MarkFlagRequired("prompt"); err != nil {
		panic(err)
	}
	return command
}

// newCmdGeminiVeo is a provider shortcut for video generation with Gemini Veo,
// equivalent to `video generate --provider gemini` with Gemini-specific defaults
// pre-filled. Users only need to supply --prompt; all other flags are optional.
func newCmdGeminiVeo(f *cmdutil.Factory) *cobra.Command {
	// Pre-fill Gemini Veo defaults; the user only needs --prompt.
	opts := &videoGenerateOptions{
		Provider:    providers.Gemini,
		Model:       "veo-2.0-generate-001",
		AspectRatio: "16:9",
		Duration:    5,
		Resolution:  "720p",
	}

	command := &cobra.Command{
		Use:   "veo",
		Short: "Generate videos with Gemini Veo",
		Annotations: map[string]string{
			"canonical": "video generate --provider " + providers.Gemini,
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVideoGenerate(f, opts)
		},
	}

	command.Flags().StringVar(&opts.Prompt, "prompt", "", "Video generation prompt")
	command.Flags().StringVar(&opts.Output, "output", "", "Output video path (.mp4)")
	command.Flags().StringVar(&opts.Model, "model", opts.Model, "Veo model to use")
	command.Flags().StringVar(&opts.AspectRatio, "aspect-ratio", opts.AspectRatio, "Video aspect ratio (16:9 or 9:16)")
	command.Flags().IntVar(&opts.Duration, "duration", opts.Duration, "Duration in seconds (5–8)")
	command.Flags().StringVar(&opts.Resolution, "resolution", opts.Resolution, "Video resolution (720p or 1080p)")

	if err := command.MarkFlagRequired("prompt"); err != nil {
		panic(err)
	}

	return command
}

// validateAspectRatio checks whether the given aspect ratio is supported by the specified Gemini model.
// An empty aspectRatio is allowed and skips validation.
func validateAspectRatio(model, aspectRatio string) error {
	if aspectRatio == "" {
		return nil
	}

	// Supported aspect ratios per model.
	aspectRatioGemini31FlashImagePreview := []string{"1:1", "1:4", "1:8", "2:3", "3:2", "3:4", "4:1", "4:3", "4:5", "5:4", "8:1", "9:16", "16:9", "21:9"}
	aspectRatioGemini3ProImagePreview := []string{"1:1", "2:3", "3:2", "3:4", "4:3", "4:5", "5:4", "9:16", "16:9", "21:9"}

	var allowed []string
	switch model {
	case "gemini-3.1-flash-image-preview":
		allowed = aspectRatioGemini31FlashImagePreview
	case "gemini-3-pro-image-preview":
		allowed = aspectRatioGemini3ProImagePreview
	default:
		allowed = []string{"1:1", "3:4", "4:3", "9:16", "16:9"}
	}

	for _, item := range allowed {
		if item == aspectRatio {
			return nil
		}
	}
	return fmt.Errorf("unsupported aspect ratio: %s", aspectRatio)
}

// validateImageSize validates and normalizes the image size for the specified Gemini model.
// The value is normalized to uppercase before comparison.
// An empty imageSize is allowed and returns ("", nil).
func validateImageSize(model, imageSize string) (string, error) {
	if imageSize == "" {
		return "", nil
	}

	normalized := strings.ToUpper(imageSize)

	// Supported image sizes per model.
	resolutionGemini31FlashImagePreview := []string{"512", "0.5K", "1K", "2K", "4K"}
	resolutionGemini3ProImagePreview := []string{"1K", "2K", "4K"}

	var allowed []string
	switch model {
	case "gemini-3.1-flash-image-preview":
		allowed = resolutionGemini31FlashImagePreview
	case "gemini-3-pro-image-preview":
		allowed = resolutionGemini3ProImagePreview
	default:
		allowed = []string{"1K", "2K", "4K"}
	}

	for _, item := range allowed {
		if item == normalized {
			return normalized, nil
		}
	}
	return "", fmt.Errorf("unsupported image size: %s", imageSize)
}

// nowFunc returns the current time. It is a variable so tests can substitute a fixed time.
var nowFunc = time.Now

// defaultImageOutputPath generates a default output file path based on the provider name and current time.
// An optional now argument can be passed to inject a fixed time in tests.
func defaultImageOutputPath(provider string, now ...time.Time) string {
	var t time.Time
	if len(now) > 0 {
		t = now[0]
	} else {
		t = nowFunc()
	}
	return fmt.Sprintf("%s_%s_%03d.png", provider, t.Format("20060102150405"), t.Nanosecond()/1e6)
}
