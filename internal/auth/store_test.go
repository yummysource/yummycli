package auth

import (
	"errors"
	"testing"
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
		return "", ErrSecretNotFound
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

func TestProviderCredentialStoreSaveAPIKey(t *testing.T) {
	fake := newMemorySecretStore()
	store := NewProviderCredentialStore(fake)

	err := store.SaveAPIKey("gemini", "test-api-key")
	if err != nil {
		t.Fatalf("SaveAPIKey returned error: %v", err)
	}

	got, err := fake.Get("yummycli", "provider:gemini:api_key")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got != "test-api-key" {
		t.Fatalf("stored api key = %q, want %q", got, "test-api-key")
	}
}

func TestProviderCredentialStoreRemoveAPIKeyRemovesSavedKey(t *testing.T) {
	fake := newMemorySecretStore()
	store := NewProviderCredentialStore(fake)

	err := store.SaveAPIKey("gemini", "test-api-key")
	if err != nil {
		t.Fatalf("SaveAPIKey returned error: %v", err)
	}

	err = store.RemoveAPIKey("gemini")
	if err != nil {
		t.Fatalf("RemoveAPIKey returned error: %v", err)
	}

	ok, err := store.HasAPIKey("gemini")
	if err != nil {
		t.Fatalf("HasAPIKey returned error after removal: %v", err)
	}
	if ok {
		t.Fatalf("HasAPIKey returned true after removal, want false")
	}
}

func TestProviderCredentialStoreHasAPIKeyReturnsTrueAfterSave(t *testing.T) {
	fake := newMemorySecretStore()
	store := NewProviderCredentialStore(fake)

	err := store.SaveAPIKey("gemini", "test-api-key")
	if err != nil {
		t.Fatalf("SaveAPIKey returned error: %v", err)
	}

	ok, err := store.HasAPIKey("gemini")
	if err != nil {
		t.Fatalf("HasAPIKey returned error: %v", err)
	}
	if !ok {
		t.Fatalf("HasAPIKey returned false, want true")
	}
}

func TestProviderCredentialStoreSaveAPIKeyRejectsUnsupportedProvider(t *testing.T) {
	fake := newMemorySecretStore()
	store := NewProviderCredentialStore(fake)

	err := store.SaveAPIKey("qwen", "test-api-key")
	if err == nil {
		t.Fatal("SaveAPIKey returned nil error for unsupported provider")
	}
}

func TestProviderCredentialStoreSaveAPIKeyRejectsEmptyAPIKey(t *testing.T) {
	fake := newMemorySecretStore()
	store := NewProviderCredentialStore(fake)

	err := store.SaveAPIKey("gemini", "")
	if err == nil {
		t.Fatal("SaveAPIKey returned nil error for empty api key")
	}
}

func TestProviderCredentialStoreSaveAPIKeyRejectsBlankAPIKey(t *testing.T) {
	fake := newMemorySecretStore()
	store := NewProviderCredentialStore(fake)

	err := store.SaveAPIKey("gemini", "   ")
	if err == nil {
		t.Fatal("SaveAPIKey returned nil error for blank api key")
	}
}

func TestProviderCredentialStoreAPIKeyPreviewReturnsMaskedValue(t *testing.T) {
	fake := newMemorySecretStore()
	store := NewProviderCredentialStore(fake)

	err := store.SaveAPIKey("gemini", "jkfe123456wer9")
	if err != nil {
		t.Fatalf("SaveAPIKey returned error: %v", err)
	}

	got, err := store.APIKeyPreview("gemini")
	if err != nil {
		t.Fatalf("APIKeyPreview returned error: %v", err)
	}

	want := "jkfe******wer9"
	if got != want {
		t.Fatalf("preview = %q, want %q", got, want)
	}
}

func TestProviderCredentialStoreAPIKeyPreviewReturnsNotFoundWhenMissing(t *testing.T) {
	fake := newMemorySecretStore()
	store := NewProviderCredentialStore(fake)

	_, err := store.APIKeyPreview("gemini")
	if err == nil {
		t.Fatal("APIKeyPreview returned nil error for missing api key")
	}
	if !errors.Is(err, ErrSecretNotFound) {
		t.Fatalf("error = %v, want %v", err, ErrSecretNotFound)
	}
}

func TestProviderCredentialStoreSaveAPIKeyNormalizesProviderName(t *testing.T) {
	fake := newMemorySecretStore()
	store := NewProviderCredentialStore(fake)

	err := store.SaveAPIKey(" GEMINI ", "test-api-key")
	if err != nil {
		t.Fatalf("SaveAPIKey returned error: %v", err)
	}

	got, err := fake.Get("yummycli", "provider:gemini:api_key")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got != "test-api-key" {
		t.Fatalf("stored api key = %q, want %q", got, "test-api-key")
	}
}
