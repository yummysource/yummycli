package cmd

import (
	"bytes"
	"testing"

	"github.com/yummysource/yummycli/internal/auth"
	"github.com/yummysource/yummycli/internal/cmdutil"
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
