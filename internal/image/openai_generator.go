package image

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"

	"github.com/yummysource/yummycli/internal/auth"
	"github.com/yummysource/yummycli/internal/providers"
)

const defaultOpenAIAPIBase = "https://api.openai.com/v1"
const openAIDefaultModelFallback = "gpt-image-2"

// OpenAIGenerator generates and edits images with the OpenAI API.
type OpenAIGenerator struct {
	credentialStore *auth.ProviderCredentialStore
	apiBase         string
}

// NewOpenAIGenerator creates an OpenAIGenerator.
// Pass an empty apiBase to use the default OpenAI API base URL.
func NewOpenAIGenerator(credentialStore *auth.ProviderCredentialStore, apiBase string) *OpenAIGenerator {
	if apiBase == "" {
		apiBase = defaultOpenAIAPIBase
	}
	return &OpenAIGenerator{
		credentialStore: credentialStore,
		apiBase:         apiBase,
	}
}

// GenerateImage generates or edits an image and writes it to the requested output path.
// When req.InputImages is non-empty, the images/edits endpoint is used.
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

	if len(req.InputImages) > 0 {
		return g.editImage(ctx, apiKey, req)
	}
	return g.generateImage(ctx, apiKey, req)
}

func (g *OpenAIGenerator) generateImage(ctx context.Context, apiKey string, req GenerateImageRequest) error {
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
	if req.OutputFormat != "" {
		body["output_format"] = req.OutputFormat
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	url := g.apiBase + "/images/generations"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	return g.doImageRequest(ctx, httpReq, req.Output)
}

func (g *OpenAIGenerator) editImage(ctx context.Context, apiKey string, req GenerateImageRequest) error {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)

	for _, imgPath := range req.InputImages {
		data, err := os.ReadFile(imgPath)
		if err != nil {
			return fmt.Errorf("reading input image %s: %w", imgPath, err)
		}
		mimeType, err := mimeTypeFromPath(imgPath)
		if err != nil {
			return err
		}
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="image[]"; filename="%s"`, filepath.Base(imgPath)))
		h.Set("Content-Type", mimeType)
		part, err := w.CreatePart(h)
		if err != nil {
			return err
		}
		if _, err := part.Write(data); err != nil {
			return err
		}
	}

	_ = w.WriteField("model", req.Model)
	_ = w.WriteField("prompt", req.Prompt)
	_ = w.WriteField("n", "1")
	if req.ImageSize != "" {
		_ = w.WriteField("size", req.ImageSize)
	}
	if req.Quality != "" {
		_ = w.WriteField("quality", req.Quality)
	}
	if req.OutputFormat != "" {
		_ = w.WriteField("output_format", req.OutputFormat)
	}
	w.Close()

	url := g.apiBase + "/images/edits"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &body)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", w.FormDataContentType())

	return g.doImageRequest(ctx, httpReq, req.Output)
}

func (g *OpenAIGenerator) doImageRequest(ctx context.Context, httpReq *http.Request, output string) error {
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
		return downloadURL(ctx, d.URL, output)
	}
	if d.B64JSON != "" {
		imgBytes, err := base64.StdEncoding.DecodeString(d.B64JSON)
		if err != nil {
			return err
		}
		return os.WriteFile(output, imgBytes, 0o644)
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
