package providers

import "testing"

func TestNormalizeReturnsGeminiProvider(t *testing.T) {
	got, err := Normalize(" gemini ")
	if err != nil {
		t.Fatalf("Normalize returned error: %v", err)
	}

	if got != Gemini {
		t.Fatalf("provider = %q, want %q", got, Gemini)
	}
}

func TestNormalizeLowercasesSupportedProvider(t *testing.T) {
	got, err := Normalize("GEMINI")
	if err != nil {
		t.Fatalf("Normalize returned error: %v", err)
	}

	if got != Gemini {
		t.Fatalf("provider = %q, want %q", got, Gemini)
	}
}

func TestNormalizeRejectsEmptyProvider(t *testing.T) {
	_, err := Normalize("   ")
	if err == nil {
		t.Fatal("Normalize returned nil error for empty provider")
	}

	want := "provider is required"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}
}

func TestNormalizeRejectsUnsupportedProvider(t *testing.T) {
	_, err := Normalize("qwen")
	if err == nil {
		t.Fatal("Normalize returned nil error for unsupported provider")
	}

	want := "unsupported provider: qwen"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}
}
