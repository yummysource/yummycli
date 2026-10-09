package image

import "context"

// GenerateImageRequest describes a text-to-image generation request.
type GenerateImageRequest struct {
	Provider    string
	Prompt      string
	Output      string
	Model       string
	AspectRatio string
	ImageSize   string
	Quality      string   // OpenAI only — low, medium, high, auto
	OutputFormat string   // OpenAI only — png, jpeg, webp
	Fallback     string   // optional fallback provider name
	InputImages  []string
}

// GenerateImageResult reports which provider and model actually produced the image,
// which differs from the request when a fallback provider was used.
type GenerateImageResult struct {
	Provider string
	Model    string
}

// ImageGenerator generates images from text prompts.
type ImageGenerator interface {
	GenerateImage(ctx context.Context, req GenerateImageRequest) (GenerateImageResult, error)
}
