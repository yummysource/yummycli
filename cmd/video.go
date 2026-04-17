package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/yummysource/yummycli/internal/cmdutil"
	"github.com/yummysource/yummycli/internal/providers"
	internalvideo "github.com/yummysource/yummycli/internal/video"
)

// videoGenerateOptions is the canonical options struct for video generation.
// Shared by `video generate` (capability layer) and `gemini veo` (shortcut).
type videoGenerateOptions struct {
	Provider    string
	Prompt      string
	Output      string
	Model       string
	AspectRatio string
	Duration    int
	Resolution  string
	// InputImages holds optional paths to local images. Count determines routing:
	//   0 → text-to-video
	//   1 → starting frame (image-to-video)
	//   2-3 → ASSET reference images
	InputImages []string
}

// NewCmdVideo creates the provider-agnostic video command group.
func NewCmdVideo(f *cmdutil.Factory) *cobra.Command {
	command := &cobra.Command{
		Use:   "video",
		Short: "Generate and manage videos",
	}

	command.AddCommand(newCmdVideoGenerate(f))

	return command
}

// newCmdVideoGenerate creates the `video generate` subcommand.
func newCmdVideoGenerate(f *cmdutil.Factory) *cobra.Command {
	opts := &videoGenerateOptions{}

	command := &cobra.Command{
		Use:   "generate",
		Short: "Generate a video from a text prompt",
		Annotations: map[string]string{
			"capability": "video.generate",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVideoGenerate(f, opts)
		},
	}

	command.Flags().StringVar(&opts.Provider, "provider", "", "AI provider (e.g. gemini)")
	command.Flags().StringVar(&opts.Prompt, "prompt", "", "Video generation prompt")
	command.Flags().StringVar(&opts.Output, "output", "", "Output file path (must end in .mp4; auto-generated if omitted)")
	command.Flags().StringVar(&opts.Model, "model", "", "Veo model to use")
	command.Flags().StringVar(&opts.AspectRatio, "aspect-ratio", "", "Video aspect ratio (16:9 or 9:16)")
	command.Flags().IntVar(&opts.Duration, "duration", 0, "Duration in seconds; veo-2: 5/6/7/8, veo-3+: 4/6/8")
	command.Flags().StringVar(&opts.Resolution, "resolution", "", "Video resolution (720p, 1080p, 4k)")
	command.Flags().StringArrayVar(&opts.InputImages, "input-image", nil, "Local PNG/JPEG image path; repeat for up to 3 (1 = starting frame, 2-3 = reference images)")

	if err := command.MarkFlagRequired("provider"); err != nil {
		panic(err)
	}
	if err := command.MarkFlagRequired("prompt"); err != nil {
		panic(err)
	}

	return command
}

// runVideoGenerate is the canonical implementation for video generation.
// Shared by `video generate` and the `gemini veo` shortcut so that both entry
// points behave identically — the shortcut only pre-fills defaults.
func runVideoGenerate(f *cmdutil.Factory, opts *videoGenerateOptions) error {
	if f.VideoGenerator == nil {
		return fmt.Errorf("video generator is not configured")
	}

	if err := applyVideoDefaults(opts); err != nil {
		return err
	}
	if err := validateVideoOptions(opts); err != nil {
		return err
	}

	// Set up a cancellable context so Ctrl+C aborts the polling loop gracefully.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
		if f.Progress != nil {
			f.Progress.Clear()
		}
		fmt.Fprintln(f.IOStreams.ErrOut, "Video generation job is still running on Veo's side. No local file written.")
		os.Exit(1)
	}()

	var progressFn func(string)
	if f.Progress != nil {
		progressFn = func(msg string) { f.Progress.Set(msg) }
	}

	req := internalvideo.GenerateVideoRequest{
		Provider:    opts.Provider,
		Model:       opts.Model,
		Prompt:      opts.Prompt,
		Output:      opts.Output,
		AspectRatio: opts.AspectRatio,
		Duration:    opts.Duration,
		Resolution:  opts.Resolution,
		InputImages: opts.InputImages,
		ProgressFn:  progressFn,
	}

	result, err := f.VideoGenerator.GenerateVideo(ctx, req)
	if err != nil {
		if f.Progress != nil {
			f.Progress.Clear()
		}
		return err
	}

	if f.Progress != nil {
		f.Progress.Clear()
	}

	return f.Output.JSON(result)
}

// applyVideoDefaults fills in provider-specific default values for empty fields.
func applyVideoDefaults(opts *videoGenerateOptions) error {
	switch opts.Provider {
	case providers.Gemini:
		if opts.Model == "" {
			opts.Model = "veo-3.1-fast-generate-preview"
		}
		if opts.AspectRatio == "" {
			opts.AspectRatio = "16:9"
		}
		if opts.Duration == 0 {
			opts.Duration = 8
		}
		if opts.Resolution == "" {
			opts.Resolution = "1080p"
		}
	default:
		return fmt.Errorf("unsupported provider: %s", opts.Provider)
	}
	return nil
}

// validDurationsForModel returns the accepted duration values (seconds) for a model.
// veo-2 accepts {5,6,7,8}; veo-3+ accepts {4,6,8}.
func validDurationsForModel(model string) []int {
	if strings.HasPrefix(model, "veo-2") {
		return []int{5, 6, 7, 8}
	}
	return []int{4, 6, 8}
}

// validResolutionsForModel returns the accepted resolution values for a model.
// veo-2 is locked to 720p; veo-3.0 adds 1080p; veo-3.1 also adds 4k.
func validResolutionsForModel(model string) []string {
	switch {
	case strings.HasPrefix(model, "veo-3.1"):
		return []string{"720p", "1080p", "4k"}
	case strings.HasPrefix(model, "veo-3"):
		return []string{"720p", "1080p"}
	default:
		return []string{"720p"}
	}
}

// validateVideoOptions runs client-side validation before submitting any API call.
func validateVideoOptions(opts *videoGenerateOptions) error {
	validDurations := validDurationsForModel(opts.Model)
	durationOK := false
	for _, v := range validDurations {
		if opts.Duration == v {
			durationOK = true
			break
		}
	}
	if !durationOK {
		allowed := make([]string, len(validDurations))
		for i, v := range validDurations {
			allowed[i] = fmt.Sprintf("%d", v)
		}
		return fmt.Errorf("duration %d is not valid for model %s (accepted: %s)", opts.Duration, opts.Model, strings.Join(allowed, ", "))
	}

	validResolutions := validResolutionsForModel(opts.Model)
	resolutionOK := false
	for _, v := range validResolutions {
		if opts.Resolution == v {
			resolutionOK = true
			break
		}
	}
	if !resolutionOK {
		return fmt.Errorf("resolution %q is not supported by model %s (accepted: %s)", opts.Resolution, opts.Model, strings.Join(validResolutions, ", "))
	}

	if (opts.Resolution == "1080p" || opts.Resolution == "4k") && opts.Duration != 8 {
		return fmt.Errorf("resolution %s requires duration=8 (got %d)", opts.Resolution, opts.Duration)
	}

	if len(opts.InputImages) > 3 {
		return fmt.Errorf("at most 3 input images are supported (got %d)", len(opts.InputImages))
	}
	for _, img := range opts.InputImages {
		if err := validateInputImage(img); err != nil {
			return err
		}
	}

	if opts.Output == "" {
		opts.Output = defaultVideoOutputPath(opts.Provider)
	} else {
		if !strings.HasSuffix(strings.ToLower(opts.Output), ".mp4") {
			return fmt.Errorf("output file must have a .mp4 extension (got %q)", opts.Output)
		}
		if _, err := os.Stat(opts.Output); err == nil {
			return fmt.Errorf("output file already exists: %s (remove it or choose a different path)", opts.Output)
		}
	}

	dir := filepath.Dir(opts.Output)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("creating output directory %s: %w", dir, err)
		}
	}

	return nil
}

// validateInputImage checks that the given file exists, is a supported image
// format (PNG or JPEG), and does not exceed the 20 MB API limit.
func validateInputImage(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("input image not found: %s", path)
	}

	const maxBytes = 20 * 1024 * 1024
	if info.Size() > maxBytes {
		return fmt.Errorf("input image exceeds the 20 MB limit (got %.1f MB): %s", float64(info.Size())/1024/1024, path)
	}

	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		return fmt.Errorf("input image must be PNG or JPEG (got %q): %s", ext, path)
	}

	return nil
}

// defaultVideoOutputPath returns an auto-generated output file name based on
// the current timestamp.
func defaultVideoOutputPath(provider string) string {
	t := time.Now()
	return fmt.Sprintf("veo_%s_%03d.mp4", t.Format("20060102_150405"), t.Nanosecond()/1e6)
}
