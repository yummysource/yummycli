package cmd

import (
	"bytes"
	"testing"

	"github.com/yummysource/yummycli/internal/auth"
	"github.com/yummysource/yummycli/internal/cmdutil"
)

type memorySecretStore struct {
	data map[string]string
}

func newMemorySecretStore() *memorySecretStore {
	return &memorySecretStore{
		data: make(map[string]string),
	}
}

func (s *memorySecretStore) Set(service, account, secret string) error {
	s.data[storeKey(service, account)] = secret
	return nil
}

func (s *memorySecretStore) Get(service, account string) (string, error) {
	value, ok := s.data[storeKey(service, account)]

	if !ok {
		return "", auth.ErrSecretNotFound
	}

	return value, nil
}

func (s *memorySecretStore) Delete(service, account string) error {
	delete(s.data, storeKey(service, account))
	return nil
}

func storeKey(service, account string) string {
	return service + "|" + account
}

func TestAuthStatusReportsGeminiAsNotConfigured(t *testing.T) {
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

	cmd := NewCmdAuth(f)

	cmd.SetArgs([]string{"status", "--provider", "gemini"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute returned error:%v", err)
	}

	got := stdout.String()
	want := "{\"provider\":\"gemini\",\"configured\":false}\n"

	if got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestAuthInitSavesGeminiAPIKey(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	fakeSecrets := newMemorySecretStore()

	f := &cmdutil.Factory{
		IOStreams: &cmdutil.IOStreams{
			Out:    stdout,
			ErrOut: stderr,
		},
		CredentialStore: auth.NewProviderCredentialStore(fakeSecrets),
	}

	cmd := NewCmdAuth(f)
	cmd.SetArgs([]string{"init", "--provider", "gemini", "--api-key", "test-api-key"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	gotSecret, err := fakeSecrets.Get("yummycli", "provider:gemini:api_key")
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

func TestAuthRemoveDeletesGeminiAPIKey(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	secretStore := newMemorySecretStore()
	err := secretStore.Set("yummycli", "provider:gemini:api_key", "test-api-key")
	if err != nil {
		t.Fatalf("api-key saved failed")
	}

	f := &cmdutil.Factory{
		IOStreams: &cmdutil.IOStreams{
			Out:    stdout,
			ErrOut: stderr,
		},
		CredentialStore: auth.NewProviderCredentialStore(secretStore),
	}

	cmd := NewCmdAuth(f)
	cmd.SetArgs([]string{"remove", "--provider", "gemini"})
	err = cmd.Execute()
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	ok, err := f.CredentialStore.HasAPIKey("gemini")
	if err != nil {
		t.Fatalf("HasAPIKey returned error: %v", err)
	}
	if ok {
		t.Fatalf("HasAPIKey returned true after remove, want false")
	}

	got := stdout.String()
	want := "{\"provider\":\"gemini\",\"removed\":true}\n"
	if got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}
