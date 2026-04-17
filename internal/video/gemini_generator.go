package video

import (
	"context"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/genai"

	"github.com/yummysource/yummycli/internal/auth"
	"github.com/yummysource/yummycli/internal/providers"
)

// pollInterval is how long to wait between operation status checks.
// 20 seconds matches the official SDK example cadence.
const pollInterval = 20 * time.Second

// generateTimeout is the maximum time to wait for a Veo job to complete.
// Video generation typically takes 1–3 minutes; 10 minutes gives ample headroom.
const generateTimeout = 10 * time.Minute

// GeminiVideoGenerator generates videos using the Google Gemini Veo API.
// It encapsulates the full lifecycle: job submission, LRO polling, and download.
type GeminiVideoGenerator struct {
	// credentialStore provides the API key for the Gemini provider.
	credentialStore *auth.ProviderCredentialStore
}

// NewGeminiGenerator creates a GeminiVideoGenerator backed by the given credential store.
func NewGeminiGenerator(credentialStore *auth.ProviderCredentialStore) *GeminiVideoGenerator {
	return &GeminiVideoGenerator{credentialStore: credentialStore}
}

// GenerateVideo generates a video from a text prompt and writes the result to
// req.Output. It blocks until the video is ready, polling every 20 seconds with
// a 10-minute timeout.
//
// The caller is responsible for all input validation (duration range, output
// path existence, .mp4 extension). This method trusts req fields are valid.
func (g *GeminiVideoGenerator) GenerateVideo(ctx context.Context, req GenerateVideoRequest) (*GenerateVideoResult, error) {
	if req.Provider != providers.Gemini {
		return nil, fmt.Errorf("unsupported provider: %s", req.Provider)
	}

	apiKey, err := g.credentialStore.GetAPIKey(req.Provider)
	if err != nil {
		return nil, err
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("creating Gemini client: %w", err)
	}

	// Build the video generation config from validated request fields.
	config := &genai.GenerateVideosConfig{
		AspectRatio: req.AspectRatio,
		Resolution:  req.Resolution,
	}
	if req.Duration > 0 {
		// DurationSeconds is *int32; use the genai.Ptr helper to take the address.
		d := int32(req.Duration)
		config.DurationSeconds = &d
	}

	// Step 1: Route image inputs based on count.
	//   0 images → text-to-video: pass nil starting frame, no reference images.
	//   1 image  → image-to-video: pass as the starting frame argument.
	//   2-3 images → reference-guided: populate config.ReferenceImages with ASSET type.
	var startingFrame *genai.Image
	switch len(req.InputImages) {
	case 0:
		// text-to-video, nothing to do

	case 1:
		img, err := loadImage(req.InputImages[0])
		if err != nil {
			return nil, err
		}
		startingFrame = img

	default:
		refs := make([]*genai.VideoGenerationReferenceImage, 0, len(req.InputImages))
		for _, path := range req.InputImages {
			img, err := loadImage(path)
			if err != nil {
				return nil, err
			}
			refs = append(refs, &genai.VideoGenerationReferenceImage{
				Image:         img,
				ReferenceType: genai.VideoGenerationReferenceTypeAsset,
			})
		}
		config.ReferenceImages = refs
	}

	// Step 2: Submit the job. GenerateVideos returns immediately with an operation
	// handle — the actual video generation happens asynchronously on Veo's side.
	startTime := time.Now()
	operation, err := client.Models.GenerateVideos(ctx, req.Model, req.Prompt, startingFrame, config)
	if err != nil {
		return nil, fmt.Errorf("submitting video generation job: %w", err)
	}

	// Step 3: Poll until done, honoring the parent context and the 10-minute timeout.
	// The timeout context ensures the CLI never hangs indefinitely even if the
	// parent context has no deadline.
	timeoutCtx, cancel := context.WithTimeout(ctx, generateTimeout)
	defer cancel()

	for !operation.Done {
		select {
		case <-timeoutCtx.Done():
			return nil, fmt.Errorf("video generation timed out after %v; the job may still complete on Veo's side", generateTimeout)
		default:
		}

		// Report progress to the caller if they provided a callback.
		if req.ProgressFn != nil {
			elapsed := int(time.Since(startTime).Seconds())
			req.ProgressFn(fmt.Sprintf("Generating video... %ds elapsed", elapsed))
		}

		// Wait before the next poll. Use select so Ctrl+C (context cancellation)
		// is honoured during the sleep, not just at the top of the loop.
		select {
		case <-timeoutCtx.Done():
			return nil, fmt.Errorf("video generation timed out after %v; the job may still complete on Veo's side", generateTimeout)
		case <-time.After(pollInterval):
		}

		operation, err = client.Operations.GetVideosOperation(timeoutCtx, operation, nil)
		if err != nil {
			return nil, fmt.Errorf("polling video operation status: %w", err)
		}
	}

	// Step 4: Check whether the completed operation succeeded or failed.
	if operation.Error != nil {
		return nil, fmt.Errorf("video generation failed: %v", operation.Error)
	}
	if operation.Response == nil || len(operation.Response.GeneratedVideos) == 0 {
		return nil, fmt.Errorf("video generation completed but no videos were returned")
	}

	// Step 5: Download the first generated video.
	// Veo returns a hosted URI; the bytes are not included inline in the response.
	// Files are retained for 48 hours before automatic deletion.
	gv := operation.Response.GeneratedVideos[0]
	data, err := client.Files.Download(timeoutCtx, genai.NewDownloadURIFromGeneratedVideo(gv), nil)
	if err != nil {
		return nil, fmt.Errorf("downloading generated video: %w", err)
	}

	// Step 6: Write the video bytes to the output file.
	if err := os.WriteFile(req.Output, data, 0o644); err != nil {
		return nil, fmt.Errorf("writing video to %s: %w", req.Output, err)
	}

	elapsed := int(time.Since(startTime).Seconds())
	return &GenerateVideoResult{
		Provider:        req.Provider,
		Output:          req.Output,
		Model:           req.Model,
		DurationSeconds: req.Duration,
		AspectRatio:     req.AspectRatio,
		Resolution:      req.Resolution,
		ElapsedSeconds:  elapsed,
		InputImages:     req.InputImages,
	}, nil
}

// loadImage reads a local file and returns a genai.Image with MIME type inferred
// from the file extension. Supports .png, .jpg, and .jpeg.
func loadImage(path string) (*genai.Image, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading image %s: %w", path, err)
	}
	mimeType := mime.TypeByExtension(filepath.Ext(path))
	if mimeType == "" {
		mimeType = "image/jpeg"
	}
	return &genai.Image{ImageBytes: data, MIMEType: mimeType}, nil
}
