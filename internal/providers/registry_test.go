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

func TestNormalizeReturnsOriginalUnsupportedProviderNameInError(t *testing.T) {
	_, err := Normalize("Qwen")
	if err == nil {
		t.Fatal("Normalize returned nil error for unsupported provider")
	}

	want := "unsupported provider: Qwen"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}
}

func TestNormalizeOpenAI(t *testing.T) {
	got, err := Normalize("openai")
	if err != nil {
		t.Fatalf("Normalize returned error: %v", err)
	}
	if got != OpenAI {
		t.Fatalf("got %q, want %q", got, OpenAI)
	}
}

func TestAllIncludesOpenAI(t *testing.T) {
	all := All()
	for _, p := range all {
		if p == OpenAI {
			return
		}
	}
	t.Fatalf("All() does not include %q", OpenAI)
}
