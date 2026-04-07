package cmdutil

import (
	"os"

	"github.com/yummysource/yummycli/internal/auth"
	"golang.org/x/term"
)

// NewDefault creates a production Factory with the process IO streams.
func NewDefault() *Factory {
	secretStore := auth.NewKeychainSecretStore()
	return &Factory{
		IOStreams: &IOStreams{
			In:         os.Stdin,
			Out:        os.Stdout,
			ErrOut:     os.Stderr,
			IsTerminal: term.IsTerminal(int(os.Stdin.Fd())),
		},
		CredentialStore: auth.NewProviderCredentialStore(secretStore),
	}
}
