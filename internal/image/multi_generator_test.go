package image

import (
	"context"
	"errors"
	"strings"
	"testing"
)

type stubGenerator struct {
	called  bool
	lastReq GenerateImageRequest
	err     error
}

func (s *stubGenerator) GenerateImage(_ context.Context, req GenerateImageRequest) error {
	s.called = true
	s.lastReq = req
	return s.err
}

func TestMultiGeneratorRoutesToCorrectProvider(t *testing.T) {
	gemini := &stubGenerator{}
	openai := &stubGenerator{}
	g := NewMultiGenerator(map[string]ImageGenerator{
		"gemini": gemini,
		"openai": openai,
	}, nil)

	err := g.GenerateImage(context.Background(), GenerateImageRequest{
		Provider: "gemini",
		Prompt:   "a cat",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !gemini.called {
		t.Fatal("gemini generator was not called")
	}
	if openai.called {
		t.Fatal("openai generator should not have been called")
	}
}

func TestMultiGeneratorClearsProviderSpecificFieldsOnFallback(t *testing.T) {
	primary := &stubGenerator{err: errors.New("primary failed")}
	fallback := &stubGenerator{}
	g := NewMultiGenerator(map[string]ImageGenerator{
		"openai": primary,
		"gemini": fallback,
	}, nil)

	err := g.GenerateImage(context.Background(), GenerateImageRequest{
		Provider:    "openai",
		Prompt:      "a cat",
		Model:       "dall-e-3",
		ImageSize:   "1024x1024",
		AspectRatio: "16:9",
		Quality:     "standard",
		Style:       "vivid",
		Fallback:    "gemini",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !fallback.called {
		t.Fatal("fallback generator was not called")
	}
	r := fallback.lastReq
	if r.Provider != "gemini" {
		t.Fatalf("fallback provider = %q, want gemini", r.Provider)
	}
	if r.Model != "" {
		t.Fatalf("fallback Model = %q, want empty", r.Model)
	}
	if r.ImageSize != "" {
		t.Fatalf("fallback ImageSize = %q, want empty", r.ImageSize)
	}
	if r.AspectRatio != "" {
		t.Fatalf("fallback AspectRatio = %q, want empty", r.AspectRatio)
	}
	if r.Quality != "" {
		t.Fatalf("fallback Quality = %q, want empty", r.Quality)
	}
	if r.Style != "" {
		t.Fatalf("fallback Style = %q, want empty", r.Style)
	}
}

func TestMultiGeneratorFallsBackOnProviderError(t *testing.T) {
	var logBuf strings.Builder
	gemini := &stubGenerator{err: errors.New("rate limited")}
	openai := &stubGenerator{}
	g := NewMultiGenerator(map[string]ImageGenerator{
		"gemini": gemini,
		"openai": openai,
	}, &logBuf)

	err := g.GenerateImage(context.Background(), GenerateImageRequest{
		Provider: "gemini",
		Prompt:   "a cat",
		Fallback: "openai",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !openai.called {
		t.Fatal("openai fallback generator was not called")
	}
	if !strings.Contains(logBuf.String(), "openai") {
		t.Fatalf("expected fallback log to mention openai, got %q", logBuf.String())
	}
}


func TestMultiGeneratorReturnsErrorWhenBothFail(t *testing.T) {
	gemini := &stubGenerator{err: errors.New("primary error")}
	openai := &stubGenerator{err: errors.New("fallback error")}
	g := NewMultiGenerator(map[string]ImageGenerator{
		"gemini": gemini,
		"openai": openai,
	}, nil)

	err := g.GenerateImage(context.Background(), GenerateImageRequest{
		Provider: "gemini",
		Prompt:   "a cat",
		Fallback: "openai",
	})
	if err == nil {
		t.Fatal("expected error when both providers fail")
	}
}

func TestMultiGeneratorReturnsErrorForUnknownProvider(t *testing.T) {
	g := NewMultiGenerator(map[string]ImageGenerator{}, nil)

	err := g.GenerateImage(context.Background(), GenerateImageRequest{
		Provider: "unknown",
		Prompt:   "a cat",
	})
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
}
