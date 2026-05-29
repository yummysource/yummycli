package auth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/yummysource/yummycli/internal/providers"
)

const secretServiceName = "yummycli"

// ErrSecretNotFound reports that the requested secret does not exist.
var ErrSecretNotFound = errors.New("secret not found")

// SecretStore abstracts secret persistence for provider credentials.
type SecretStore interface {
	Set(service, account, secret string) error
	Get(service, account string) (string, error)
	Delete(service, account string) error
}

// ProviderCredentialStore manages provider API keys.
type ProviderCredentialStore struct {
	store SecretStore
}

// NewProviderCredentialStore creates a ProviderCredentialStore.
func NewProviderCredentialStore(store SecretStore) *ProviderCredentialStore {
	return &ProviderCredentialStore{
		store: store,
	}
}

// SaveAPIKey stores an API key for the given provider.
func (s *ProviderCredentialStore) SaveAPIKey(provider, apiKey string) error {
	normalizedProvider, err := providers.Normalize(provider)
	if err != nil {
		return err
	}
	if strings.TrimSpace(apiKey) == "" {
		return fmt.Errorf("api key is required")
	}

	account := accountName(normalizedProvider)
	return s.store.Set(secretServiceName, account, apiKey)
}

// HasAPIKey reports whether an API key exists for the given provider.
func (s *ProviderCredentialStore) HasAPIKey(provider string) (bool, error) {
	normalizedProvider, err := providers.Normalize(provider)
	if err != nil {
		return false, err
	}

	value, err := s.store.Get(secretServiceName, accountName(normalizedProvider))
	if err != nil {
		if errors.Is(err, ErrSecretNotFound) {
			return false, nil
		}
		return false, err
	}

	return strings.TrimSpace(value) != "", nil
}

// RemoveAPIKey deletes the API key for the given provider.
func (s *ProviderCredentialStore) RemoveAPIKey(provider string) error {
	normalizedProvider, err := providers.Normalize(provider)
	if err != nil {
		return err
	}

	return s.store.Delete(secretServiceName, accountName(normalizedProvider))
}

// APIKeyPreview returns a masked preview of the API key for the given provider.
func (s *ProviderCredentialStore) APIKeyPreview(provider string) (string, error) {
	normalizedProvider, err := providers.Normalize(provider)
	if err != nil {
		return "", err
	}

	value, err := s.store.Get(secretServiceName, accountName(normalizedProvider))
	if err != nil {
		return "", err
	}

	return maskAPIKey(value), nil
}

func accountName(provider string) string {
	return "provider:" + provider + ":api_key"
}

func maskAPIKey(value string) string {
	if len(value) < 8 {
		return "******"
	}

	return value[:4] + "******" + value[len(value)-4:]
}

// ProviderStatus holds the credential status for a single provider.
type ProviderStatus struct {
	Provider      string
	Configured    bool
	APIKeyPreview string // non-empty only when Configured is true
}

// ListConfigured returns the credential status for every registered provider.
func (s *ProviderCredentialStore) ListConfigured() ([]ProviderStatus, error) {
	all := providers.All()
	result := make([]ProviderStatus, 0, len(all))

	for _, provider := range all {
		configured, err := s.HasAPIKey(provider)
		if err != nil {
			return nil, err
		}

		status := ProviderStatus{
			Provider:   provider,
			Configured: configured,
		}

		if configured {
			preview, err := s.APIKeyPreview(provider)
			if err != nil {
				return nil, err
			}
			status.APIKeyPreview = preview
		}

		result = append(result, status)
	}

	return result, nil
}

const configDefaultProviderAccount = "config:default_provider"

// SetDefaultProvider persists the default provider name.
func (s *ProviderCredentialStore) SetDefaultProvider(provider string) error {
	return s.store.Set(secretServiceName, configDefaultProviderAccount, provider)
}

// GetDefaultProvider returns the configured default provider, or "" if not set.
func (s *ProviderCredentialStore) GetDefaultProvider() (string, error) {
	value, err := s.store.Get(secretServiceName, configDefaultProviderAccount)
	if err != nil {
		if errors.Is(err, ErrSecretNotFound) {
			return "", nil
		}
		return "", err
	}
	return value, nil
}

// GetAPIKey returns the stored API key for the given provider.
func (s *ProviderCredentialStore) GetAPIKey(provider string) (string, error) {
	normalizedProvider, err := providers.Normalize(provider)
	if err != nil {
		return "", err
	}

	return s.store.Get(secretServiceName, accountName(normalizedProvider))
}
