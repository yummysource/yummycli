package image

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/genai"

	"github.com/yummysource/yummycli/internal/auth"
	"github.com/yummysource/yummycli/internal/providers"
)

// GeminiGenerator generates images with the Gemini API.
type GeminiGenerator struct {
	credentialStore      *auth.ProviderCredentialStore
	interactionsEndpoint string
	httpClient           *http.Client
}

// NewGeminiGenerator creates a GeminiGenerator.
func NewGeminiGenerator(credentialStore *auth.ProviderCredentialStore) *GeminiGenerator {
	return &GeminiGenerator{
		credentialStore:      credentialStore,
		interactionsEndpoint: "https://generativelanguage.googleapis.com/v1beta/interactions",
		httpClient:           http.DefaultClient,
	}
}

const geminiDefaultModelFallback = "gemini-nano-banana-2.1"

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
	if req.Model == geminiDefaultModelFallback {
		return g.generateWithInteractions(ctx, apiKey, req)
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

type interactionInput struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	Data     string `json:"data,omitempty"`
	MIMEType string `json:"mime_type,omitempty"`
}

type interactionImageFormat struct {
	Type        string `json:"type"`
	AspectRatio string `json:"aspect_ratio,omitempty"`
	ImageSize   string `json:"image_size,omitempty"`
}

type createInteractionRequest struct {
	Model          string                 `json:"model"`
	Input          []interactionInput     `json:"input"`
	ResponseFormat interactionImageFormat `json:"response_format"`
}

type interactionContent struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type createInteractionResponse struct {
	Steps []struct {
		Content []interactionContent `json:"content"`
	} `json:"steps"`
}

func (g *GeminiGenerator) generateWithInteractions(ctx context.Context, apiKey string, req GenerateImageRequest) error {
	inputs := []interactionInput{{Type: "text", Text: req.Prompt}}
	for _, imagePath := range req.InputImages {
		data, err := os.ReadFile(imagePath)
		if err != nil {
			return err
		}
		mimeType, err := mimeTypeFromPath(imagePath)
		if err != nil {
			return err
		}
		inputs = append(inputs, interactionInput{
			Type:     "image",
			Data:     base64.StdEncoding.EncodeToString(data),
			MIMEType: mimeType,
		})
	}

	payload, err := json.Marshal(createInteractionRequest{
		Model: req.Model,
		Input: inputs,
		ResponseFormat: interactionImageFormat{
			Type:        "image",
			AspectRatio: req.AspectRatio,
			ImageSize:   req.ImageSize,
		},
	})
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, g.interactionsEndpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-goog-api-key", apiKey)

	resp, err := g.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		responseBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		if err != nil {
			return err
		}
		return fmt.Errorf("gemini interactions API returned %s: %s", resp.Status, strings.TrimSpace(string(responseBody)))
	}

	var result createInteractionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	for _, step := range result.Steps {
		for _, content := range step.Content {
			if content.Type != "image" || content.Data == "" {
				continue
			}
			imageData, err := base64.StdEncoding.DecodeString(content.Data)
			if err != nil {
				return err
			}
			return os.WriteFile(req.Output, imageData, 0o644)
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
