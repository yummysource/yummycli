package cmd

import (
	"bytes"
	"context"
	"testing"

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

func TestGeminiNanoBananaRequiresOutputFlag(t *testing.T) {
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
	cmd.SetArgs([]string{"nanobanana", "--prompt", "a banana on a plate"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute returned nil error without --output")
	}

	want := "required flag(s) \"output\" not set"
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
