package image

import (
	"context"
	"errors"
	"strings"
	"testing"
)

type stubGenerator struct {
	called   bool
	provider string
	err      error
}

func (s *stubGenerator) GenerateImage(_ context.Context, req GenerateImageRequest) error {
	s.called = true
	s.provider = req.Provider
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
