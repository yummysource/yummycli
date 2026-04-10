package cmdutil

import (
	"os"

	"github.com/yummysource/yummycli/internal/auth"
	"github.com/yummysource/yummycli/internal/image"
	"github.com/yummysource/yummycli/internal/output"
	"golang.org/x/term"
)

// NewDefault creates a production Factory with the process IO streams.
func NewDefault() *Factory {
	secretStore := auth.NewKeychainSecretStore()
	credentialStore := auth.NewProviderCredentialStore(secretStore)
	streams := &IOStreams{
		In:         os.Stdin,
		Out:        os.Stdout,
		ErrOut:     os.Stderr,
		IsTerminal: term.IsTerminal(int(os.Stdin.Fd())),
	}
	return &Factory{
		IOStreams:       streams,
		CredentialStore: credentialStore,
		ImageGenerator:  image.NewGeminiGenerator(credentialStore),
		Output:          output.New(streams.Out),
	}
}
