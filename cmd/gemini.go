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
const geminiDefaultModel = "gemini-3.1-flash-image"

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
		newCmdGeminiSpeak(f),
		newCmdGeminiVoices(f),
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
		Short: "Generate or edit images with Gemini",
		Long: `Shortcut for "yummycli image generate --provider gemini" with Gemini defaults pre-filled.

DEFAULTS
  --model        gemini-3.1-flash-image
  --aspect-ratio 16:9
  --image-size   1K

ASPECT RATIO
  Supported: 1:1 1:4 1:8 2:3 3:2 3:4 4:1 4:3 4:5 5:4 8:1 9:16 16:9 21:9

IMAGE SIZE
  Supported: 512 0.5K 1K 2K 4K
  Note: 512 and 0.5K are gemini-3.1-flash-image only.

MODELS
  gemini-3.1-flash-image   (default) fast, supports all aspect ratios and 512/0.5K sizes
  gemini-3-pro-image-preview  higher quality; aspect ratios: 1:1 2:3 3:2 3:4 4:3 4:5 5:4 9:16 16:9 21:9

IMAGE EDITING
  Pass --input-image to edit an existing image instead of generating from scratch.`,
		Example: `  # Text-to-image with defaults (16:9, 1K)
  yummycli gemini nanobanana --prompt "a panda eating bamboo"

  # Widescreen 4K
  yummycli gemini nanobanana --prompt "city skyline at night" --aspect-ratio 21:9 --image-size 4K

  # Edit an image
  yummycli gemini nanobanana --prompt "make it watercolor style" --input-image ./photo.png

  # Multi-image compositing
  yummycli gemini nanobanana --prompt "blend into one scene" \
    --input-image ./subject.png --input-image ./background.jpg`,
		Annotations: map[string]string{
			"canonical": "image generate --provider " + providers.Gemini + " --preset nano-banana",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runImageGenerate(f, opts)
		},
	}

	command.Flags().StringVar(&opts.Prompt, "prompt", "", "image generation or editing prompt")
	command.Flags().StringVar(&opts.Output, "output", "", "output file path (auto-generated if omitted)")
	command.Flags().StringVar(&opts.Model, "model", opts.Model, "Gemini image model (default: gemini-3.1-flash-image)")
	command.Flags().StringVar(&opts.AspectRatio, "aspect-ratio", opts.AspectRatio, "output aspect ratio (default: 16:9)")
	command.Flags().StringVar(&opts.ImageSize, "image-size", opts.ImageSize, "output resolution: 512 | 0.5K | 1K | 2K | 4K (default: 1K)")
	command.Flags().StringArrayVar(&opts.InputImages, "input-image", nil, "input image path for editing; repeat for multiple")

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
		Model:       "veo-3.1-fast-generate-preview",
		AspectRatio: "16:9",
		Duration:    8,
		Resolution:  "1080p",
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
	command.Flags().IntVar(&opts.Duration, "duration", opts.Duration, "Duration in seconds; veo-2: 5/6/7/8, veo-3+: 4/6/8")
	command.Flags().StringVar(&opts.Resolution, "resolution", opts.Resolution, "Video resolution (720p, 1080p, 4k; veo-3.1 only for 4k)")
	command.Flags().StringArrayVar(&opts.InputImages, "input-image", nil, "Local PNG/JPEG image path; repeat for up to 3 (1 = starting frame, 2-3 = reference images)")

	if err := command.MarkFlagRequired("prompt"); err != nil {
		panic(err)
	}

	return command
}

// newCmdGeminiSpeak is a provider shortcut for speech synthesis with Gemini TTS,
// equivalent to `audio speak --provider gemini` with Gemini-specific defaults pre-filled.
func newCmdGeminiSpeak(f *cmdutil.Factory) *cobra.Command {
	opts := &audioSpeakOptions{
		Provider: providers.Gemini,
	}

	command := &cobra.Command{
		Use:   "speak",
		Short: "Synthesise text to speech with Gemini TTS",
		Annotations: map[string]string{
			"canonical": "audio speak --provider " + providers.Gemini,
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAudioSpeak(f, opts)
		},
	}

	command.Flags().StringVar(&opts.Text, "text", "", "Text to synthesise")
	command.Flags().StringVar(&opts.Output, "output", "", "Output WAV file path (auto-generated if omitted)")
	command.Flags().StringVar(&opts.Model, "model", "", "TTS model to use")
	command.Flags().StringVar(&opts.Voice, "voice", "", "Prebuilt voice name (e.g. Aoede); mutually exclusive with --speaker")
	command.Flags().StringVar(&opts.LanguageCode, "language", "", "BCP-47 language code (auto-detected from text if omitted)")
	command.Flags().StringArrayVar(&opts.Speakers, "speaker", nil, "Multi-speaker mapping Name:Voice; repeat up to 2 (e.g. --speaker Alice:Aoede)")

	if err := command.MarkFlagRequired("text"); err != nil {
		panic(err)
	}

	return command
}

// newCmdGeminiVoices is a provider shortcut for listing Gemini TTS voices,
// equivalent to `audio voices --provider gemini`.
func newCmdGeminiVoices(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "voices",
		Short: "List available Gemini TTS voices",
		Annotations: map[string]string{
			"canonical": "audio voices --provider " + providers.Gemini,
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAudioVoices(f, providers.Gemini)
		},
	}
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
	case "gemini-3.1-flash-image":
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
	case "gemini-3.1-flash-image":
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

// defaultImageOutputPath generates a default output file path.
// format sets the file extension (e.g. "jpeg"); empty defaults to "png".
// An optional now argument can be passed to inject a fixed time in tests.
func defaultImageOutputPath(provider, format string, now ...time.Time) string {
	var t time.Time
	if len(now) > 0 {
		t = now[0]
	} else {
		t = nowFunc()
	}
	ext := format
	if ext == "" {
		ext = "png"
	}
	return fmt.Sprintf("%s_%s_%03d.%s", provider, t.Format("20060102150405"), t.Nanosecond()/1e6, ext)
}
