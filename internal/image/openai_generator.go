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
const openAIDefaultModelFallback = "gpt-image-1"

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

	if req.Model == "" {
		req.Model = openAIDefaultModelFallback
	}

	apiKey, err := g.credentialStore.GetAPIKey(req.Provider)
	if err != nil {
		return err
	}

	body := map[string]any{
		"model":  req.Model,
		"prompt": req.Prompt,
		"n":      1,
	}
	if req.ImageSize != "" {
		body["size"] = req.ImageSize
	}
	if req.Quality != "" {
		body["quality"] = req.Quality
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
			URL     string `json:"url"`
			B64JSON string `json:"b64_json"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	if len(result.Data) == 0 {
		return fmt.Errorf("no image data returned")
	}

	d := result.Data[0]
	if d.URL != "" {
		return downloadURL(ctx, d.URL, req.Output)
	}
	if d.B64JSON != "" {
		imgBytes, err := base64.StdEncoding.DecodeString(d.B64JSON)
		if err != nil {
			return err
		}
		return os.WriteFile(req.Output, imgBytes, 0o644)
	}
	return fmt.Errorf("no image data returned")
}

func downloadURL(ctx context.Context, url, output string) error {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download image: status %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return os.WriteFile(output, data, 0o644)
}
