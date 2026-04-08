package image

import (
	"context"
	"fmt"
	"os"

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

// GenerateImage generates an image and writes it to the requested output path.
func (g *GeminiGenerator) GenerateImage(ctx context.Context, req GenerateImageRequest) error {
	if req.Provider != providers.Gemini {
		return fmt.Errorf("unsupported provider: %s", req.Provider)
	}

	apiKey, err := g.credentialStore.APIKey(req.Provider)
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

	result, err := client.Models.GenerateContent(
		ctx,
		req.Model,
		genai.Text(req.Prompt),
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
