package video

import "context"

// VideoGenerator generates a video from a text prompt and writes it to disk.
// Implementations are responsible for the full lifecycle: API submission,
// polling until the operation completes, downloading the result, and writing
// the output file.
//
// The interface is intentionally single-step (no separate Submit/Poll/Download
// methods) so callers remain simple. Internal polling is an implementation detail.
// If Phase 2 introduces detached operations, a new interface can be added without
// changing this one.
//
// This mirrors internal/image.ImageGenerator exactly so the Factory pattern
// and test injection approach are consistent across capability packages.
type VideoGenerator interface {
	GenerateVideo(ctx context.Context, req GenerateVideoRequest) (*GenerateVideoResult, error)
}

// GenerateVideoRequest holds all parameters for a single video generation call.
// All validation (duration range, resolution/duration compatibility, output path
// constraints) must be performed by the caller before passing this struct to
// GenerateVideo.
type GenerateVideoRequest struct {
	// Provider identifies the AI backend, e.g. "gemini".
	Provider string

	// Model is the provider-specific model ID, e.g. "veo-2.0-generate-001".
	Model string

	// Prompt is the natural-language description of the video to generate.
	Prompt string

	// Output is the destination file path. Must end in ".mp4" and must not
	// already exist. If empty, the generator should return an error.
	Output string

	// AspectRatio is the desired video aspect ratio, e.g. "16:9" or "9:16".
	AspectRatio string

	// Duration is the video length in seconds. Valid values: 4, 6, 8.
	Duration int

	// Resolution is the video resolution, e.g. "720p" or "1080p".
	// 1080p requires Duration == 8.
	Resolution string

	// ProgressFn is an optional callback invoked during polling with a status
	// message such as "Generating video... 45s elapsed". A nil value is safe.
	// Using a callback rather than an interface avoids circular package imports
	// between video and cmdutil.
	ProgressFn func(message string)
}

// GenerateVideoResult holds metadata about a successfully generated video.
// It is returned on success and serialised to JSON on stdout by the command layer.
type GenerateVideoResult struct {
	// Provider is the AI provider used, e.g. "gemini".
	Provider string `json:"provider"`

	// Output is the path of the written video file.
	Output string `json:"output"`

	// Model is the model that generated the video.
	Model string `json:"model"`

	// DurationSeconds is the video length as requested.
	DurationSeconds int `json:"duration_seconds"`

	// AspectRatio is the video aspect ratio, e.g. "16:9".
	AspectRatio string `json:"aspect_ratio"`

	// Resolution is the video resolution, e.g. "720p".
	Resolution string `json:"resolution"`

	// ElapsedSeconds is the wall-clock time from submission to file written.
	ElapsedSeconds int `json:"elapsed_seconds"`
}
