package cmd

import (
	"bytes"
	"testing"

	"github.com/yummysource/yummycli/internal/auth"
	"github.com/yummysource/yummycli/internal/cmdutil"
	"github.com/yummysource/yummycli/internal/output"
	"github.com/yummysource/yummycli/internal/providers"
)

func newInitFactory(stdout, stderr *bytes.Buffer) *cmdutil.Factory {
	return &cmdutil.Factory{
		IOStreams:        &cmdutil.IOStreams{Out: stdout, ErrOut: stderr},
		CredentialStore:  auth.NewProviderCredentialStore(newMemorySecretStore()),
		Output:           output.New(stdout),
	}
}

func TestInitSavesAPIKey(t *testing.T) {
	stdout := &bytes.Buffer{}
	f := newInitFactory(stdout, &bytes.Buffer{})

	cmd := NewCmdInit(f)
	cmd.SetArgs([]string{"--provider", providers.Gemini, "--api-key", "test-key"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	configured, err := f.CredentialStore.HasAPIKey(providers.Gemini)
	if err != nil {
		t.Fatalf("HasAPIKey returned error: %v", err)
	}
	if !configured {
		t.Fatal("expected gemini to be configured after init")
	}
}

func TestInitWithDefaultFlagSetsDefaultProvider(t *testing.T) {
	stdout := &bytes.Buffer{}
	f := newInitFactory(stdout, &bytes.Buffer{})

	cmd := NewCmdInit(f)
	cmd.SetArgs([]string{"--provider", providers.Gemini, "--api-key", "test-key", "--default"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	got, err := f.CredentialStore.GetDefaultProvider()
	if err != nil {
		t.Fatalf("GetDefaultProvider returned error: %v", err)
	}
	if got != providers.Gemini {
		t.Fatalf("default provider = %q, want %q", got, providers.Gemini)
	}
}

func TestInitWithoutDefaultFlagDoesNotChangeDefaultProvider(t *testing.T) {
	stdout := &bytes.Buffer{}
	f := newInitFactory(stdout, &bytes.Buffer{})

	// Pre-set gemini as default.
	_ = f.CredentialStore.SetDefaultProvider(providers.Gemini)

	cmd := NewCmdInit(f)
	cmd.SetArgs([]string{"--provider", providers.OpenAI, "--api-key", "test-key"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	got, err := f.CredentialStore.GetDefaultProvider()
	if err != nil {
		t.Fatalf("GetDefaultProvider returned error: %v", err)
	}
	if got != providers.Gemini {
		t.Fatalf("default provider = %q, want %q (should be unchanged)", got, providers.Gemini)
	}
}

func TestInitWritesJSONResult(t *testing.T) {
	stdout := &bytes.Buffer{}
	f := newInitFactory(stdout, &bytes.Buffer{})

	cmd := NewCmdInit(f)
	cmd.SetArgs([]string{"--provider", providers.Gemini, "--api-key", "test-key", "--default"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	want := "{\"provider\":\"gemini\",\"configured\":true,\"default\":true}\n"
	if got := stdout.String(); got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestInitWritesJSONResultWithDefaultFalseWhenNoDefaultFlag(t *testing.T) {
	stdout := &bytes.Buffer{}
	f := newInitFactory(stdout, &bytes.Buffer{})

	cmd := NewCmdInit(f)
	cmd.SetArgs([]string{"--provider", providers.Gemini, "--api-key", "test-key"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	want := "{\"provider\":\"gemini\",\"configured\":true,\"default\":false}\n"
	if got := stdout.String(); got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestInitRequiresProviderFlag(t *testing.T) {
	stdout := &bytes.Buffer{}
	f := newInitFactory(stdout, &bytes.Buffer{})
	cmd := NewCmdInit(f)
	cmd.SetArgs([]string{"--api-key", "test-key"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error without --provider")
	}
}

func TestInitRequiresAPIKeyFlag(t *testing.T) {
	stdout := &bytes.Buffer{}
	f := newInitFactory(stdout, &bytes.Buffer{})
	cmd := NewCmdInit(f)
	cmd.SetArgs([]string{"--provider", providers.Gemini})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error without --api-key")
	}
}
