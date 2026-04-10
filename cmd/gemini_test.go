package cmd

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/yummysource/yummycli/internal/auth"
	"github.com/yummysource/yummycli/internal/cmdutil"
	internalimage "github.com/yummysource/yummycli/internal/image"
	"github.com/yummysource/yummycli/internal/providers"
)

func TestGeminiInitSavesAPIKeyWithGeminiProvider(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	secretStore := newMemorySecretStore()

	f := &cmdutil.Factory{
		IOStreams: &cmdutil.IOStreams{
			Out:    stdout,
			ErrOut: stderr,
		},
		CredentialStore: auth.NewProviderCredentialStore(secretStore),
	}

	cmd := NewCmdGemini(f)
	cmd.SetArgs([]string{"init", "--api-key", "test-api-key"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	gotSecret, err := secretStore.Get("yummycli", "provider:gemini:api_key")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if gotSecret != "test-api-key" {
		t.Fatalf("stored api key = %q, want %q", gotSecret, "test-api-key")
	}

	got := stdout.String()
	want := "{\"provider\":\"gemini\",\"configured\":true}\n"
	if got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestGeminiInitRequiresAPIKeyFlag(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	secretStore := newMemorySecretStore()

	f := &cmdutil.Factory{
		IOStreams: &cmdutil.IOStreams{
			Out:    stdout,
			ErrOut: stderr,
		},
		CredentialStore: auth.NewProviderCredentialStore(secretStore),
	}

	cmd := NewCmdGemini(f)
	cmd.SetArgs([]string{"init"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute returned nil error without --api-key")
	}

	want := "required flag(s) \"api-key\" not set"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}
}

func TestGeminiInitStoresCredentialsUnderGeminiProviderKey(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	secretStore := newMemorySecretStore()

	f := &cmdutil.Factory{
		IOStreams: &cmdutil.IOStreams{
			Out:    stdout,
			ErrOut: stderr,
		},
		CredentialStore: auth.NewProviderCredentialStore(secretStore),
	}

	cmd := NewCmdGemini(f)
	cmd.SetArgs([]string{"init", "--api-key", "test-api-key"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	got, err := secretStore.Get("yummycli", "provider:gemini:api_key")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got != "test-api-key" {
		t.Fatalf("stored api key = %q, want %q", got, "test-api-key")
	}
}

func TestGeminiNanoBananaRequiresPromptFlag(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	secretStore := newMemorySecretStore()

	f := &cmdutil.Factory{
		IOStreams: &cmdutil.IOStreams{
			Out:    stdout,
			ErrOut: stderr,
		},
		CredentialStore: auth.NewProviderCredentialStore(secretStore),
	}

	cmd := NewCmdGemini(f)
	cmd.SetArgs([]string{"nanobanana", "--output", "./result.png"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute returned nil error without --prompt")
	}

	want := "required flag(s) \"prompt\" not set"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}
}

type fakeImageGenerator struct {
	called bool
	req    internalimage.GenerateImageRequest
	err    error
}

func (g *fakeImageGenerator) GenerateImage(_ context.Context, req internalimage.GenerateImageRequest) error {
	g.called = true
	g.req = req
	return g.err
}

func TestGeminiNanoBananaUsesDefaultModelWhenGenerating(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	secretStore := newMemorySecretStore()
	generator := &fakeImageGenerator{}

	f := &cmdutil.Factory{
		IOStreams: &cmdutil.IOStreams{
			Out:    stdout,
			ErrOut: stderr,
		},
		CredentialStore: auth.NewProviderCredentialStore(secretStore),
		ImageGenerator:  generator,
	}

	cmd := NewCmdGemini(f)
	cmd.SetArgs([]string{
		"nanobanana",
		"--prompt", "a banana on a plate",
		"--output", "./result.png",
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
	if generator.req.Prompt != "a banana on a plate" {
		t.Fatalf("prompt = %q, want %q", generator.req.Prompt, "a banana on a plate")
	}
	if generator.req.Output != "./result.png" {
		t.Fatalf("output = %q, want %q", generator.req.Output, "./result.png")
	}
	if generator.req.Model != "gemini-3.1-flash-image-preview" {
		t.Fatalf("model = %q, want %q", generator.req.Model, "gemini-3.1-flash-image-preview")
	}
}

func TestGeminiNanoBananaUsesExplicitModelWhenProvided(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	secretStore := newMemorySecretStore()
	generator := &fakeImageGenerator{}

	f := &cmdutil.Factory{
		IOStreams: &cmdutil.IOStreams{
			Out:    stdout,
			ErrOut: stderr,
		},
		CredentialStore: auth.NewProviderCredentialStore(secretStore),
		ImageGenerator:  generator,
	}

	cmd := NewCmdGemini(f)
	cmd.SetArgs([]string{
		"nanobanana",
		"--prompt", "a banana on a plate",
		"--output", "./result.png",
		"--model", "gemini-3-pro-image-preview",
	})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if !generator.called {
		t.Fatal("GenerateImage was not called")
	}

	want := "gemini-3-pro-image-preview"
	if generator.req.Model != want {
		t.Fatalf("model = %q, want %q", generator.req.Model, want)
	}
}

func TestGeminiNanoBananaReturnsErrorWhenImageGeneratorIsNotConfigured(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	secretStore := newMemorySecretStore()

	f := &cmdutil.Factory{
		IOStreams: &cmdutil.IOStreams{
			Out:    stdout,
			ErrOut: stderr,
		},
		CredentialStore: auth.NewProviderCredentialStore(secretStore),
	}

	cmd := NewCmdGemini(f)
	cmd.SetArgs([]string{
		"nanobanana",
		"--prompt", "a banana on a plate",
		"--output", "./result.png",
	})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute returned nil error without an image generator")
	}

	want := "image generator is not configured"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}
}

func TestGeminiNanoBananaRejectsUnsupportedAspectRatio(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	secretStore := newMemorySecretStore()
	f := &cmdutil.Factory{
		IOStreams: &cmdutil.IOStreams{
			Out:    stdout,
			ErrOut: stderr,
		},
		CredentialStore: auth.NewProviderCredentialStore(secretStore),
	}

	cmd := NewCmdGemini(f)
	cmd.SetArgs([]string{
		"nanobanana",
		"--prompt", "a milk banana on a plate",
		"--output", "result11.png",
		"--aspect-ratio", "5:7",
	})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute returned nil error for unsupported aspect ratio")
	}

	want := "unsupported aspect ratio: 5:7"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}
}

func TestGeminiNanoBananaRejectsUnsupportedImageSize(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	secretStore := newMemorySecretStore()
	generator := &fakeImageGenerator{}

	f := &cmdutil.Factory{
		IOStreams: &cmdutil.IOStreams{
			Out:    stdout,
			ErrOut: stderr,
		},
		CredentialStore: auth.NewProviderCredentialStore(secretStore),
		ImageGenerator:  generator,
	}

	cmd := NewCmdGemini(f)
	cmd.SetArgs([]string{
		"nanobanana",
		"--prompt", "a banana on a plate",
		"--output", "./result.png",
		"--image-size", "3K",
	})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute returned nil error for unsupported image size")
	}

	want := "unsupported image size: 3K"
	if err.Error() != want {
		t.Fatalf("error = %q, want %q", err.Error(), want)
	}
}

func TestGeminiNanoBananaPassesAspectRatioAndImageSizeToGenerator(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	secretStore := newMemorySecretStore()
	generator := &fakeImageGenerator{}

	f := &cmdutil.Factory{
		IOStreams: &cmdutil.IOStreams{
			Out:    stdout,
			ErrOut: stderr,
		},
		CredentialStore: auth.NewProviderCredentialStore(secretStore),
		ImageGenerator:  generator,
	}

	cmd := NewCmdGemini(f)
	cmd.SetArgs([]string{
		"nanobanana",
		"--prompt", "a banana on a plate",
		"--output", "./result.png",
		"--aspect-ratio", "21:9",
		"--image-size", "4K",
	})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if !generator.called {
		t.Fatal("GenerateImage was not called")
	}

	if generator.req.AspectRatio != "21:9" {
		t.Fatalf("aspect ratio = %q, want %q", generator.req.AspectRatio, "21:9")
	}
	if generator.req.ImageSize != "4K" {
		t.Fatalf("image size = %q, want %q", generator.req.ImageSize, "4K")
	}
}

func TestGeminiNanoBananaUsesDefaultOutputFileNameWhenOutputIsOmitted(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	secretStore := newMemorySecretStore()
	generator := &fakeImageGenerator{}

	originalNowFunc := nowFunc
	nowFunc = func() time.Time {
		return time.Date(2026, 4, 9, 15, 30, 45, 123000000, time.Local)
	}
	defer func() {
		nowFunc = originalNowFunc
	}()

	f := &cmdutil.Factory{
		IOStreams: &cmdutil.IOStreams{
			Out:    stdout,
			ErrOut: stderr,
		},
		CredentialStore: auth.NewProviderCredentialStore(secretStore),
		ImageGenerator:  generator,
	}

	cmd := NewCmdGemini(f)
	cmd.SetArgs([]string{
		"nanobanana",
		"--prompt", "a banana on a plate",
	})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if !generator.called {
		t.Fatal("GenerateImage was not called")
	}

	want := "gemini_20260409153045_123.png"
	if generator.req.Output != want {
		t.Fatalf("output = %q, want %q", generator.req.Output, want)
	}
}

func TestGeminiNanoBananaPassesInputImageToGenerator(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	secretStore := newMemorySecretStore()
	generator := &fakeImageGenerator{}

	f := &cmdutil.Factory{
		IOStreams: &cmdutil.IOStreams{
			Out:    stdout,
			ErrOut: stderr,
		},
		CredentialStore: auth.NewProviderCredentialStore(secretStore),
		ImageGenerator:  generator,
	}

	cmd := NewCmdGemini(f)
	cmd.SetArgs([]string{
		"nanobanana",
		"--prompt", "edit this image into a watercolor style",
		"--input-image", "./source.png",
	})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if !generator.called {
		t.Fatal("GenerateImage was not called")
	}

	if generator.req.InputImages[0] != "./source.png" {
		t.Fatalf("input image = %q, want %q", generator.req.InputImages[0], "./source.png")
	}
}

func TestGeminiNanoBananaPassesMultipleInputImagesToGenerator(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	secretStore := newMemorySecretStore()
	generator := &fakeImageGenerator{}

	f := &cmdutil.Factory{
		IOStreams: &cmdutil.IOStreams{
			Out:    stdout,
			ErrOut: stderr,
		},
		CredentialStore: auth.NewProviderCredentialStore(secretStore),
		ImageGenerator:  generator,
	}

	cmd := NewCmdGemini(f)
	cmd.SetArgs([]string{
		"nanobanana",
		"--prompt", "blend these into one watercolor composition",
		"--input-image", "./a.png",
		"--input-image", "./b.png",
	})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if !generator.called {
		t.Fatal("GenerateImage was not called")
	}

	if len(generator.req.InputImages) != 2 {
		t.Fatalf("input images length = %d, want %d", len(generator.req.InputImages), 2)
	}
	if generator.req.InputImages[0] != "./a.png" {
		t.Fatalf("first input image = %q, want %q", generator.req.InputImages[0], "./a.png")
	}
	if generator.req.InputImages[1] != "./b.png" {
		t.Fatalf("second input image = %q, want %q", generator.req.InputImages[1], "./b.png")
	}
}
