package cmd

import (
	"bytes"
	"testing"
	"time"

	"github.com/yummysource/yummycli/internal/auth"
	"github.com/yummysource/yummycli/internal/cmdutil"
	"github.com/yummysource/yummycli/internal/output"
	"github.com/yummysource/yummycli/internal/providers"
)

// newImageGenerateFactory builds a Factory for image generate tests.
// Pass nil for generator to leave ImageGenerator unset (used to test nil guards).
func newImageGenerateFactory(stdout, stderr *bytes.Buffer, generator *fakeImageGenerator) *cmdutil.Factory {
	f := &cmdutil.Factory{
		IOStreams: &cmdutil.IOStreams{
			Out:    stdout,
			ErrOut: stderr,
		},
		CredentialStore: auth.NewProviderCredentialStore(newMemorySecretStore()),
		Output:          output.New(stdout),
	}
	if generator != nil {
		f.ImageGenerator = generator
	}
	return f
}

func TestImageGenerateRequiresProviderFlag(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	f := newImageGenerateFactory(stdout, stderr, &fakeImageGenerator{})

	cmd := NewCmdImage(f)
	cmd.SetArgs([]string{"generate", "--prompt", "a cat"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute returned nil error without --provider")
	}

	want := "required flag(s) \"provider\" not set"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}
}

func TestImageGenerateRequiresPromptFlag(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	f := newImageGenerateFactory(stdout, stderr, &fakeImageGenerator{})

	cmd := NewCmdImage(f)
	cmd.SetArgs([]string{"generate", "--provider", providers.Gemini})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute returned nil error without --prompt")
	}

	want := "required flag(s) \"prompt\" not set"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}
}

func TestImageGenerateRejectsUnsupportedProvider(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	f := newImageGenerateFactory(stdout, stderr, &fakeImageGenerator{})

	cmd := NewCmdImage(f)
	cmd.SetArgs([]string{"generate", "--provider", "openai", "--prompt", "a cat"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute returned nil error for unsupported provider")
	}

	want := "unsupported provider: openai"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}
}

func TestImageGenerateReturnsErrorWhenImageGeneratorIsNotConfigured(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	f := newImageGenerateFactory(stdout, stderr, nil)

	cmd := NewCmdImage(f)
	cmd.SetArgs([]string{"generate", "--provider", providers.Gemini, "--prompt", "a cat"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute returned nil error without an image generator")
	}

	want := "image generator is not configured"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}
}

func TestImageGenerateUsesDefaultModelForGemini(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	generator := &fakeImageGenerator{}
	f := newImageGenerateFactory(stdout, stderr, generator)

	cmd := NewCmdImage(f)
	cmd.SetArgs([]string{
		"generate",
		"--provider", providers.Gemini,
		"--prompt", "a sunset",
		"--output", "out.png",
	})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if !generator.called {
		t.Fatal("GenerateImage was not called")
	}
	if generator.req.Provider != providers.Gemini {
		t.Fatalf("provider = %q, want %q", generator.req.Provider, providers.Gemini)
	}
	if generator.req.Model != geminiDefaultModel {
		t.Fatalf("model = %q, want %q", generator.req.Model, geminiDefaultModel)
	}
}

func TestImageGenerateRejectsUnsupportedAspectRatioForGemini(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	generator := &fakeImageGenerator{}
	f := newImageGenerateFactory(stdout, stderr, generator)

	cmd := NewCmdImage(f)
	cmd.SetArgs([]string{
		"generate",
		"--provider", providers.Gemini,
		"--prompt", "a cat",
		"--output", "out.png",
		"--aspect-ratio", "7:3",
	})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute returned nil error for unsupported aspect ratio")
	}

	want := "unsupported aspect ratio: 7:3"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}
}

func TestImageGenerateRejectsUnsupportedImageSizeForGemini(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	generator := &fakeImageGenerator{}
	f := newImageGenerateFactory(stdout, stderr, generator)

	cmd := NewCmdImage(f)
	cmd.SetArgs([]string{
		"generate",
		"--provider", providers.Gemini,
		"--prompt", "a cat",
		"--output", "out.png",
		"--image-size", "8K",
	})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute returned nil error for unsupported image size")
	}

	want := "unsupported image size: 8K"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}
}

func TestImageGenerateNormalizesImageSizeToUpperCase(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	generator := &fakeImageGenerator{}
	f := newImageGenerateFactory(stdout, stderr, generator)

	cmd := NewCmdImage(f)
	cmd.SetArgs([]string{
		"generate",
		"--provider", providers.Gemini,
		"--prompt", "a cat",
		"--output", "out.png",
		"--image-size", "4k",
	})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if generator.req.ImageSize != "4K" {
		t.Fatalf("image size = %q, want %q", generator.req.ImageSize, "4K")
	}
}

func TestImageGenerateUsesDefaultOutputPathWhenOutputIsOmitted(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	generator := &fakeImageGenerator{}
	f := newImageGenerateFactory(stdout, stderr, generator)

	originalNowFunc := nowFunc
	nowFunc = func() time.Time {
		return time.Date(2026, 4, 9, 10, 20, 30, 456000000, time.Local)
	}
	defer func() { nowFunc = originalNowFunc }()

	cmd := NewCmdImage(f)
	cmd.SetArgs([]string{
		"generate",
		"--provider", providers.Gemini,
		"--prompt", "a cat",
	})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	want := "gemini_20260409102030_456.png"
	if generator.req.Output != want {
		t.Fatalf("output = %q, want %q", generator.req.Output, want)
	}
}

func TestImageGenerateWritesJSONResultOnSuccess(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	generator := &fakeImageGenerator{}
	f := newImageGenerateFactory(stdout, stderr, generator)

	originalNowFunc := nowFunc
	nowFunc = func() time.Time {
		return time.Date(2026, 4, 9, 10, 20, 30, 456000000, time.Local)
	}
	defer func() { nowFunc = originalNowFunc }()

	cmd := NewCmdImage(f)
	cmd.SetArgs([]string{
		"generate",
		"--provider", providers.Gemini,
		"--prompt", "a sunset",
		"--input-image", "./a.png",
	})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	got := stdout.String()
	want := "{\"provider\":\"gemini\",\"output\":\"gemini_20260409102030_456.png\",\"model\":\"gemini-3.1-flash-image-preview\",\"inputImageCount\":1}\n"
	if got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}
