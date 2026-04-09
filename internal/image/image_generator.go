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
}

// ImageGenerator generates images from text prompts.
type ImageGenerator interface {
	GenerateImage(ctx context.Context, req GenerateImageRequest) error
}
