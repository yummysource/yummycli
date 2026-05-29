package image

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/yummysource/yummycli/internal/auth"
	"github.com/yummysource/yummycli/internal/providers"
)

// memorySecretStore is an in-memory SecretStore used in tests.
type memorySecretStore struct {
	data map[string]string
}

func newMemorySecretStore() *memorySecretStore {
	return &memorySecretStore{data: make(map[string]string)}
}

func (s *memorySecretStore) Set(service, account, secret string) error {
	s.data[service+"|"+account] = secret
	return nil
}

func (s *memorySecretStore) Get(service, account string) (string, error) {
	v, ok := s.data[service+"|"+account]
	if !ok {
		return "", auth.ErrSecretNotFound
	}
	return v, nil
}

func (s *memorySecretStore) Delete(service, account string) error {
	delete(s.data, service+"|"+account)
	return nil
}

// newOpenAITestServer returns a test server that serves:
//   - POST /  → JSON with a url pointing back to the server's /image path
//   - GET /image → the provided imageData bytes
func newOpenAITestServer(t *testing.T, imageData []byte) *httptest.Server {
	t.Helper()
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "image/png")
			w.Write(imageData)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]any{
			"data": []map[string]string{{"url": srv.URL + "/image"}},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	return srv
}

func TestOpenAIGeneratorRejectsNonOpenAIProvider(t *testing.T) {
	store := auth.NewProviderCredentialStore(newMemorySecretStore())
	g := NewOpenAIGenerator(store, "")

	err := g.GenerateImage(context.Background(), GenerateImageRequest{
		Provider: providers.Gemini,
		Prompt:   "a cat",
		Output:   "out.png",
		Model:    "dall-e-3",
	})
	if err == nil {
		t.Fatal("expected error for non-openai provider")
	}
}

func TestOpenAIGeneratorWritesImageFile(t *testing.T) {
	fakeBytes := []byte("fake-image-data")
	srv := newOpenAITestServer(t, fakeBytes)
	defer srv.Close()

	secretStore := newMemorySecretStore()
	_ = secretStore.Set("yummycli", "provider:openai:api_key", "test-key")
	store := auth.NewProviderCredentialStore(secretStore)

	g := NewOpenAIGenerator(store, srv.URL)

	dir := t.TempDir()
	outPath := filepath.Join(dir, "out.png")

	err := g.GenerateImage(context.Background(), GenerateImageRequest{
		Provider:  providers.OpenAI,
		Prompt:    "a cat",
		Output:    outPath,
		Model:     "dall-e-3",
		ImageSize: "1024x1024",
		Quality:   "standard",
		Style:     "vivid",
	})
	if err != nil {
		t.Fatalf("GenerateImage returned error: %v", err)
	}

	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if string(got) != string(fakeBytes) {
		t.Fatalf("file contents = %q, want %q", got, fakeBytes)
	}
}

func TestOpenAIGeneratorReturnsErrorOnAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]string{"message": "invalid api key"},
		})
	}))
	defer srv.Close()

	secretStore := newMemorySecretStore()
	_ = secretStore.Set("yummycli", "provider:openai:api_key", "bad-key")
	store := auth.NewProviderCredentialStore(secretStore)

	g := NewOpenAIGenerator(store, srv.URL)

	err := g.GenerateImage(context.Background(), GenerateImageRequest{
		Provider: providers.OpenAI,
		Prompt:   "a cat",
		Output:   "out.png",
		Model:    "dall-e-3",
	})
	if err == nil {
		t.Fatal("expected error for API error response")
	}
}
