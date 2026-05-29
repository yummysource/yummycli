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

	"github.com/yummysource/yummycli/internal/auth"
	"github.com/yummysource/yummycli/internal/providers"
)

const defaultOpenAIImagesURL = "https://api.openai.com/v1/images/generations"

// OpenAIGenerator generates images with the OpenAI DALL-E API.
type OpenAIGenerator struct {
	credentialStore *auth.ProviderCredentialStore
	baseURL         string
}

// NewOpenAIGenerator creates an OpenAIGenerator.
// Pass an empty baseURL to use the default OpenAI endpoint.
func NewOpenAIGenerator(credentialStore *auth.ProviderCredentialStore, baseURL string) *OpenAIGenerator {
	if baseURL == "" {
		baseURL = defaultOpenAIImagesURL
	}
	return &OpenAIGenerator{
		credentialStore: credentialStore,
		baseURL:         baseURL,
	}
}

// GenerateImage generates an image and writes it to the requested output path.
func (g *OpenAIGenerator) GenerateImage(ctx context.Context, req GenerateImageRequest) error {
	if req.Provider != providers.OpenAI {
		return fmt.Errorf("unsupported provider: %s", req.Provider)
	}

	apiKey, err := g.credentialStore.GetAPIKey(req.Provider)
	if err != nil {
		return err
	}

	body := map[string]any{
		"model":           req.Model,
		"prompt":          req.Prompt,
		"n":               1,
		"response_format": "b64_json",
	}
	if req.ImageSize != "" {
		body["size"] = req.ImageSize
	}
	if req.Quality != "" {
		body["quality"] = req.Quality
	}
	if req.Style != "" {
		body["style"] = req.Style
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, g.baseURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("openai api error %d: %s", resp.StatusCode, string(raw))
	}

	var result struct {
		Data []struct {
			B64JSON string `json:"b64_json"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	if len(result.Data) == 0 || result.Data[0].B64JSON == "" {
		return fmt.Errorf("no image data returned")
	}

	imgBytes, err := base64.StdEncoding.DecodeString(result.Data[0].B64JSON)
	if err != nil {
		return err
	}

	return os.WriteFile(req.Output, imgBytes, 0o644)
}
