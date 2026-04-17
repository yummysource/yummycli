package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/yummysource/yummycli/internal/cmdspec"
	"github.com/yummysource/yummycli/internal/cmdutil"
	"github.com/yummysource/yummycli/internal/providers"
	internalvideo "github.com/yummysource/yummycli/internal/video"
)

// NewCmdVideo creates the provider-agnostic video command group, loading its
// subcommands dynamically from the embedded JSON spec files.
//
// This mirrors how lark-cli's RegisterServiceCommands works: the JSON spec
// drives the command tree, and Go handles the actual execution logic.
func NewCmdVideo(f *cmdutil.Factory) *cobra.Command {
	specs, err := cmdspec.LoadAll()
	if err != nil {
		// Non-fatal: return an empty group so the binary still starts.
		fmt.Fprintf(os.Stderr, "warning: could not load capability specs: %v\n", err)
		return &cobra.Command{Use: "video", Short: "Generate and manage videos"}
	}

	// Find the "video" capability spec and build its command group.
	for _, spec := range specs {
		if spec.Capability == "video" {
			return buildCapabilityCommand(f, spec)
		}
	}

	return &cobra.Command{Use: "video", Short: "Generate and manage videos"}
}

// buildCapabilityCommand creates a Cobra command group for a CapabilitySpec and
// registers one subcommand per operation. This is the yummycli equivalent of
// lark-cli's registerService + registerResource + registerMethod chain.
func buildCapabilityCommand(f *cmdutil.Factory, spec cmdspec.CapabilitySpec) *cobra.Command {
	group := &cobra.Command{
		Use:   spec.Capability,
		Short: spec.Short,
	}

	for _, op := range spec.Operations {
		op := op // capture loop variable for the closure
		group.AddCommand(buildOperationCommand(f, spec, op))
	}

	return group
}

// buildOperationCommand creates a single Cobra leaf command for an OperationSpec.
// Flags are registered dynamically from the spec; the run function is resolved
// by capability+operation name via dispatchVideoOperation.
//
// flagValues stores per-flag pointers so Cobra can write flag values into them.
// Map values are not addressable in Go, so we use a pointer-per-flag pattern
// rather than a plain map[string]string.
func buildOperationCommand(f *cmdutil.Factory, cap cmdspec.CapabilitySpec, op cmdspec.OperationSpec) *cobra.Command {
	// fvStr and fvInt hold pointers to string and int flag values respectively.
	// Each pointer is allocated once and passed to Cobra via Flags().StringVar /
	// Flags().IntVar, so the CLI runtime can write the parsed value into it.
	fvStr := make(map[string]*string)
	fvInt := make(map[string]*int)

	cmd := &cobra.Command{
		Use:   op.Use,
		Short: op.Short,
		RunE: func(cmd *cobra.Command, args []string) error {
			return dispatchVideoOperation(f, cap.Capability+"."+op.Use, fvStr, fvInt)
		},
	}

	// Register each flag defined in the spec onto the Cobra command.
	// Enum values are embedded in the usage string so --help shows them.
	for _, flag := range op.Flags {
		usage := buildFlagUsage(flag)

		switch flag.Type {
		case "string":
			p := new(string)
			*p = flag.Default
			fvStr[flag.Name] = p
			cmd.Flags().StringVar(p, flag.Name, flag.Default, usage)

		case "int":
			p := new(int)
			if v, err := strconv.Atoi(flag.Default); err == nil {
				*p = v
			}
			fvInt[flag.Name] = p
			defaultInt, _ := strconv.Atoi(flag.Default)
			cmd.Flags().IntVar(p, flag.Name, defaultInt, usage)
		}

		if flag.Required {
			if err := cmd.MarkFlagRequired(flag.Name); err != nil {
				panic(fmt.Sprintf("cmdspec: marking %s as required: %v", flag.Name, err))
			}
		}
	}

	return cmd
}

// buildFlagUsage constructs a usage string for a FlagSpec, appending the enum
// list when present so users can see valid values in --help output.
func buildFlagUsage(flag cmdspec.FlagSpec) string {
	if len(flag.Enum) > 0 {
		return fmt.Sprintf("%s (one of: %s)", flag.Usage, strings.Join(flag.Enum, ", "))
	}
	return flag.Usage
}

// dispatchVideoOperation routes a capability+operation key to its run function.
// This mirrors lark-cli's serviceMethodRun: the dispatch layer receives resolved
// flag values and forwards them to the appropriate implementation.
//
// Adding a new operation only requires:
//  1. A new entry in video.json.
//  2. A new case here pointing to a new runXxx function.
func dispatchVideoOperation(f *cmdutil.Factory, key string, fvStr map[string]*string, fvInt map[string]*int) error {
	getString := func(name string) string {
		if p, ok := fvStr[name]; ok {
			return *p
		}
		return ""
	}
	getInt := func(name string) int {
		if p, ok := fvInt[name]; ok {
			return *p
		}
		return 0
	}

	switch key {
	case "video.generate":
		opts := &videoGenerateOptions{
			Provider:    getString("provider"),
			Prompt:      getString("prompt"),
			Output:      getString("output"),
			Model:       getString("model"),
			AspectRatio: getString("aspect-ratio"),
			Duration:    getInt("duration"),
			Resolution:  getString("resolution"),
			InputImage:  getString("input-image"),
		}
		return runVideoGenerate(f, opts)

	default:
		return fmt.Errorf("unknown operation: %s", key)
	}
}

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
	// InputImage is an optional path to a local image file used as the starting
	// frame. When set the request is treated as image-to-video; otherwise text-to-video.
	InputImage string
}

// runVideoGenerate is the canonical implementation for video generation.
// It validates inputs, sets up signal handling and progress reporting, calls
// the VideoGenerator, and writes the result as JSON to stdout.
//
// Shared by `video generate` and the `gemini veo` shortcut so that both entry
// points behave identically — the shortcut only pre-fills defaults.
func runVideoGenerate(f *cmdutil.Factory, opts *videoGenerateOptions) error {
	if f.VideoGenerator == nil {
		return fmt.Errorf("video generator is not configured")
	}

	// Apply provider-specific defaults and validate inputs before any API call.
	if err := applyVideoDefaults(opts); err != nil {
		return err
	}
	if err := validateVideoOptions(opts); err != nil {
		return err
	}

	// Set up a cancellable context so Ctrl+C aborts the polling loop gracefully.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Catch SIGINT (Ctrl+C) and cancel the context with a user-visible message.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
		// Clear the progress line before printing so the message is not garbled.
		if f.Progress != nil {
			f.Progress.Clear()
		}
		fmt.Fprintln(f.IOStreams.ErrOut, "Video generation job is still running on Veo's side. No local file written.")
		os.Exit(1)
	}()

	// Wire progress reporting: update ErrOut on each poll cycle.
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
		InputImage:  opts.InputImage,
		ProgressFn:  progressFn,
	}

	result, err := f.VideoGenerator.GenerateVideo(ctx, req)
	if err != nil {
		if f.Progress != nil {
			f.Progress.Clear()
		}
		return err
	}

	// Clear the progress line before writing JSON to stdout.
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

// validDurationsForModel returns the set of accepted duration values (seconds) for a model.
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
	default: // veo-2
		return []string{"720p"}
	}
}

// validateVideoOptions runs client-side validation before submitting any API call.
// Errors returned here are fast-fail: they prevent wasting quota on a bad request.
func validateVideoOptions(opts *videoGenerateOptions) error {
	// Duration must be one of the discrete values the model accepts.
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

	// Resolution must be one the model supports.
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

	// High-resolution clips require a full 8-second duration.
	if (opts.Resolution == "1080p" || opts.Resolution == "4k") && opts.Duration != 8 {
		return fmt.Errorf("resolution %s requires duration=8 (got %d)", opts.Resolution, opts.Duration)
	}

	// Validate the input image when provided.
	if opts.InputImage != "" {
		if err := validateInputImage(opts.InputImage); err != nil {
			return err
		}
	}

	// Validate or generate the output path.
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

	// Ensure the output directory is writable.
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

	const maxBytes = 20 * 1024 * 1024 // 20 MB — Veo API hard limit
	if info.Size() > maxBytes {
		return fmt.Errorf("input image exceeds the 20 MB limit (got %.1f MB): %s", float64(info.Size())/1024/1024, path)
	}

	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		return fmt.Errorf("input image must be PNG or JPEG (got %q): %s", ext, path)
	}

	return nil
}

// defaultVideoOutputPath returns an auto-generated output file name based on the
// provider name and current timestamp, following the same pattern as defaultImageOutputPath.
func defaultVideoOutputPath(provider string) string {
	t := time.Now()
	return fmt.Sprintf("veo_%s_%03d.mp4", t.Format("20060102_150405"), t.Nanosecond()/1e6)
}
