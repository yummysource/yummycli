package image

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestGeminiDefaultModelFallbackUsesNanoBanana21(t *testing.T) {
	if geminiDefaultModelFallback != "gemini-nano-banana-2.1" {
		t.Fatalf("default model = %q, want %q", geminiDefaultModelFallback, "gemini-nano-banana-2.1")
	}
}

func TestGenerateWithInteractionsSendsNanoBanana21RequestAndWritesImage(t *testing.T) {
	var gotRequest createInteractionRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want %q", r.Method, http.MethodPost)
		}
		if got := r.Header.Get("x-goog-api-key"); got != "test-api-key" {
			t.Errorf("api key header = %q, want %q", got, "test-api-key")
		}
		if err := json.NewDecoder(r.Body).Decode(&gotRequest); err != nil {
			t.Errorf("decode request: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"test-interaction","steps":[{"type":"model_output","content":[{"type":"image","data":"` + base64.StdEncoding.EncodeToString([]byte("generated-image")) + `","mime_type":"image/png"}]}]}`))
	}))
	defer server.Close()

	inputPath := filepath.Join(t.TempDir(), "reference.png")
	if err := os.WriteFile(inputPath, []byte("reference-image"), 0o644); err != nil {
		t.Fatalf("write input image: %v", err)
	}
	outputPath := filepath.Join(t.TempDir(), "output.png")
	generator := &GeminiGenerator{
		interactionsEndpoint: server.URL,
		httpClient:           server.Client(),
	}
	req := GenerateImageRequest{
		Model:       "gemini-nano-banana-2.1",
		Prompt:      "make a poster",
		Output:      outputPath,
		AspectRatio: "16:9",
		ImageSize:   "2K",
		InputImages: []string{inputPath},
	}

	if err := generator.generateWithInteractions(t.Context(), "test-api-key", req); err != nil {
		t.Fatalf("generateWithInteractions returned error: %v", err)
	}

	if gotRequest.Model != req.Model {
		t.Fatalf("model = %q, want %q", gotRequest.Model, req.Model)
	}
	if gotRequest.ResponseFormat.Type != "image" {
		t.Fatalf("response format type = %q, want %q", gotRequest.ResponseFormat.Type, "image")
	}
	if gotRequest.ResponseFormat.AspectRatio != "16:9" {
		t.Fatalf("aspect ratio = %q, want %q", gotRequest.ResponseFormat.AspectRatio, "16:9")
	}
	if gotRequest.ResponseFormat.ImageSize != "2K" {
		t.Fatalf("image size = %q, want %q", gotRequest.ResponseFormat.ImageSize, "2K")
	}
	if len(gotRequest.Input) != 2 {
		t.Fatalf("input length = %d, want %d", len(gotRequest.Input), 2)
	}
	if gotRequest.Input[0].Type != "text" || gotRequest.Input[0].Text != req.Prompt {
		t.Fatalf("text input = %#v", gotRequest.Input[0])
	}
	if gotRequest.Input[1].Type != "image" || gotRequest.Input[1].MIMEType != "image/png" {
		t.Fatalf("image input = %#v", gotRequest.Input[1])
	}

	gotImage, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read output image: %v", err)
	}
	if string(gotImage) != "generated-image" {
		t.Fatalf("output image = %q, want %q", gotImage, "generated-image")
	}
}

func TestMIMETypeFromPath(t *testing.T) {
	for _, testCase := range []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{name: "png", path: "a.png", want: "image/png"},
		{name: "jpg", path: "a.jpg", want: "image/jpeg"},
		{name: "jpeg", path: "a.jpeg", want: "image/jpeg"},
		{name: "webp", path: "a.webp", want: "image/webp"},
		{name: "unsupported", path: "a.gif", wantErr: true},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := mimeTypeFromPath(testCase.path)
			if testCase.wantErr {
				if err == nil {
					t.Fatal("mimeTypeFromPath returned nil error")
				}
				return
			}

			if err != nil {
				t.Fatalf("mimeTypeFromPath returned error: %v", err)
			}
			if got != testCase.want {
				t.Fatalf("mime type = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestBuildGenerateContentPartsWithoutInputImage(t *testing.T) {
	req := GenerateImageRequest{
		Prompt: "turn this into watercolor",
	}

	parts, err := buildGenerateContentParts(req)
	if err != nil {
		t.Fatalf("buildGenerateContentParts returned error: %v", err)
	}

	if len(parts) != 1 {
		t.Fatalf("parts length = %d, want %d", len(parts), 1)
	}
	if parts[0].Text != "turn this into watercolor" {
		t.Fatalf("text part = %q, want %q", parts[0].Text, "turn this into watercolor")
	}
}

func TestBuildGenerateContentPartsWithInputImage(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "source.png")

	err := os.WriteFile(inputPath, []byte("fake-image-bytes"), 0o644)
	if err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	req := GenerateImageRequest{
		Prompt:      "turn this into watercolor",
		InputImages: []string{inputPath},
	}

	parts, err := buildGenerateContentParts(req)
	if err != nil {
		t.Fatalf("buildGenerateContentParts returned error: %v", err)
	}

	if len(parts) != 2 {
		t.Fatalf("parts length = %d, want %d", len(parts), 2)
	}
	if parts[0].InlineData == nil {
		t.Fatal("first part inline data is nil")
	}
	if parts[0].InlineData.MIMEType != "image/png" {
		t.Fatalf("mime type = %q, want %q", parts[0].InlineData.MIMEType, "image/png")
	}
	if string(parts[0].InlineData.Data) != "fake-image-bytes" {
		t.Fatalf("image bytes = %q, want %q", string(parts[0].InlineData.Data), "fake-image-bytes")
	}
	if parts[1].Text != "turn this into watercolor" {
		t.Fatalf("text part = %q, want %q", parts[1].Text, "turn this into watercolor")
	}
}
