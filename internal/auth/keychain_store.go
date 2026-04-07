package auth

import (
	"errors"

	"github.com/zalando/go-keyring"
)

// KeychainSecretStore stores secrets in the OS keychain.
type KeychainSecretStore struct{}

// NewKeychainSecretStore creates a KeychainSecretStore.
func NewKeychainSecretStore() *KeychainSecretStore {
	return &KeychainSecretStore{}
}

// Set stores a secret in the OS keychain.
func (k *KeychainSecretStore) Set(service, account, secret string) error {
	return keyring.Set(service, account, secret)
}

// Get reads a secret from the OS keychain.
func (k *KeychainSecretStore) Get(service, account string) (string, error) {
	value, err := keyring.Get(service, account)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return "", ErrSecretNotFound
		}
		return "", err
	}

	return value, nil
}

// Delete removes a secret from the OS keychain.
func (k *KeychainSecretStore) Delete(service, account string) error {
	err := keyring.Delete(service, account)
	if err != nil && !errors.Is(err, keyring.ErrNotFound) {
		return err
	}

	return nil
}
