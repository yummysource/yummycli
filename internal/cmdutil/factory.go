package cmdutil

import "github.com/yummysource/yummycli/internal/auth"

// Factory holds shared dependencies injected into commands.
type Factory struct {
	IOStreams       *IOStreams
	CredentialStore *auth.ProviderCredentialStore
}
