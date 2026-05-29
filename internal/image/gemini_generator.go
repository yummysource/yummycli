package image

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/genai"

	"github.com/yummysource/yummycli/internal/auth"
	"github.com/yummysource/yummycli/internal/providers"
)

// GeminiGenerator generates images with the Gemini API.
type GeminiGenerator struct {
	credentialStore *auth.ProviderCredentialStore
}

// NewGeminiGenerator creates a GeminiGenerator.
func NewGeminiGenerator(credentialStore *auth.ProviderCredentialStore) *GeminiGenerator {
	return &GeminiGenerator{
		credentialStore: credentialStore,
	}
}

const geminiDefaultModelFallback = "gemini-3.1-flash-image-preview"

// GenerateImage generates an image and writes it to the requested output path.
func (g *GeminiGenerator) GenerateImage(ctx context.Context, req GenerateImageRequest) error {
	if req.Provider != providers.Gemini {
		return fmt.Errorf("unsupported provider: %s", req.Provider)
	}

	if req.Model == "" {
		req.Model = geminiDefaultModelFallback
	}

	apiKey, err := g.credentialStore.GetAPIKey(req.Provider)
	if err != nil {
		return err
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return err
	}

	config := &genai.GenerateContentConfig{
		ResponseModalities: []string{"IMAGE", "TEXT"},
	}

	if req.AspectRatio != "" || req.ImageSize != "" {
		config.ImageConfig = &genai.ImageConfig{}
	}

	if req.AspectRatio != "" {
		config.ImageConfig.AspectRatio = req.AspectRatio
	}

	if req.ImageSize != "" {
		config.ImageConfig.ImageSize = req.ImageSize
	}

	parts, err := buildGenerateContentParts(req)
	if err != nil {
		return err
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, err := client.Models.GenerateContent(
		ctx,
		req.Model,
		contents,
		config,
	)
	if err != nil {
		return err
	}

	for _, candidate := range result.Candidates {
		if candidate.Content == nil {
			continue
		}
		for _, part := range candidate.Content.Parts {
			if part.InlineData != nil && len(part.InlineData.Data) > 0 {
				return os.WriteFile(req.Output, part.InlineData.Data, 0o644)
			}
		}
	}

	return fmt.Errorf("no image data returned")
}

func mimeTypeFromPath(path string) (string, error) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".png":
		return "image/png", nil
	case ".jpg", ".jpeg":
		return "image/jpeg", nil
	case ".webp":
		return "image/webp", nil
	default:
		return "", fmt.Errorf("unsupported input image type: %s", path)
	}
}

func buildGenerateContentParts(req GenerateImageRequest) ([]*genai.Part, error) {
	parts := make([]*genai.Part, 0, 2)

	for _, imagePath := range req.InputImages {
		data, err := os.ReadFile(imagePath)
		if err != nil {
			return nil, err
		}

		mimeType, err := mimeTypeFromPath(imagePath)
		if err != nil {
			return nil, err
		}

		parts = append(parts, genai.NewPartFromBytes(data, mimeType))
	}

	parts = append(parts, genai.NewPartFromText(req.Prompt))

	return parts, nil
}
