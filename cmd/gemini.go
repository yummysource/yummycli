package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/yummysource/yummycli/internal/cmdutil"
	internalimage "github.com/yummysource/yummycli/internal/image"
	"github.com/yummysource/yummycli/internal/providers"
)

type geminiNanoBananaOptions struct {
	Prompt      string
	Output      string
	Model       string
	AspectRatio string
	ImageSize   string
	InputImages []string
}

type geminiNanoBananaResult struct {
	Provider        string `json:"provider"`
	Output          string `json:"output"`
	Model           string `json:"model"`
	InputImageCount int    `json:"inputImageCount"`
}

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
	command.Flags().StringVar(&opts.AspectRatio, "aspect-ratio", "16:9", "Gemini image aspect-ratio")
	command.Flags().StringVar(&opts.ImageSize, "image-size", "1K", "Gemini image size")
	command.Flags().StringArrayVar(&opts.InputImages, "input-image", nil, "Gemini input image path")

	if err := command.MarkFlagRequired("prompt"); err != nil {
		panic(err)
	}
	return command
}

func runGeminiNanoBanana(f *cmdutil.Factory, opts *geminiNanoBananaOptions) error {
	errAspectRatio := validateAspectRatio(opts)
	errImageSize := validateImageSize(opts)

	if errAspectRatio != nil {
		return errAspectRatio
	}

	if errImageSize != nil {
		return errImageSize
	}

	if f.ImageGenerator == nil {
		return fmt.Errorf("image generator is not configured")
	}
	if f.Output == nil {
		return fmt.Errorf("output writer is not configured")
	}

	if opts.Output == "" {
		opts.Output = defaultImageOutputPath(providers.Gemini)
	}

	req := internalimage.GenerateImageRequest{
		Provider:    providers.Gemini,
		Prompt:      opts.Prompt,
		Output:      opts.Output,
		Model:       opts.Model,
		AspectRatio: opts.AspectRatio,
		ImageSize:   opts.ImageSize,
		InputImages: opts.InputImages,
	}

	if err := f.ImageGenerator.GenerateImage(context.Background(), req); err != nil {
		return err
	}

	result := geminiNanoBananaResult{
		Provider:        providers.Gemini,
		Output:          opts.Output,
		Model:           opts.Model,
		InputImageCount: len(opts.InputImages),
	}

	return f.Output.JSON(result)
}

func validateAspectRatio(opts *geminiNanoBananaOptions) error {
	aspectRatioGemini31FlashImagePreview := []string{"1:1", "1:4", "1:8", "2:3", "3:2", "3:4", "4:1", "4:3", "4:5", "5:4", "8:1", "9:16", "16:9", "21:9"}

	aspectRatioGemini3ProImagePreview := []string{"1:1", "2:3", "3:2", "3:4", "4:3", "4:5", "5:4", "9:16", "16:9", "21:9"}

	var aspectRatio []string
	switch opts.Model {
	case "gemini-3.1-flash-image-preview":
		aspectRatio = aspectRatioGemini31FlashImagePreview

	case "gemini-3-pro-image-preview":
		aspectRatio = aspectRatioGemini3ProImagePreview
	default:
		aspectRatio = []string{"1:1", "3:4", "4:3", "9:16", "16:9"}
	}

	isInAspectRatio := false
	for _, item := range aspectRatio {
		if item == opts.AspectRatio {
			isInAspectRatio = true
		}
	}

	if !isInAspectRatio {
		// return errors.New("unsupported aspect-ratio")
		return fmt.Errorf("unsupported aspect ratio: %s", opts.AspectRatio)
	}
	return nil
}

func validateImageSize(opts *geminiNanoBananaOptions) error {
	resolutionGemini31FlashImagePreview := []string{"512", "0.5K", "1K", "2K", "4K"}
	resolutionGemini3ProImagePreview := []string{"1K", "2K", "4K"}

	var resolution []string
	switch opts.Model {
	case "gemini-3.1-flash-image-preview":
		resolution = resolutionGemini31FlashImagePreview

	case "gemini-3-pro-image-preview":
		resolution = resolutionGemini3ProImagePreview
	default:
		resolution = []string{"1K", "2K", "4K"}
	}

	opts.ImageSize = strings.ToUpper(opts.ImageSize)
	isInResolution := false
	for _, item := range resolution {
		if item == opts.ImageSize {
			isInResolution = true
		}
	}

	if !isInResolution {
		// return errors.New("unsupported aspect-ratio")
		return fmt.Errorf("unsupported image size: %s", opts.ImageSize)
	}
	return nil
}

var nowFunc = time.Now

func defaultImageOutputPath(provider string, now ...time.Time) string {
	var t time.Time
	if len(now) > 0 {
		t = now[0]
	} else {
		t = nowFunc()
	}
	return fmt.Sprintf("%s_%s_%03d.png", provider, t.Format("20060102150405"), t.Nanosecond()/1e6)
}
