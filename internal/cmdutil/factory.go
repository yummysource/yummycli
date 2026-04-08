package cmdutil

import (
	"github.com/yummysource/yummycli/internal/auth"
	internalimage "github.com/yummysource/yummycli/internal/image"
)

// Factory holds shared dependencies injected into commands.
type Factory struct {
	IOStreams       *IOStreams
	CredentialStore *auth.ProviderCredentialStore
	ImageGenerator  internalimage.ImageGenerator
}
