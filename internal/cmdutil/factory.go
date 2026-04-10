package cmdutil

import (
	"github.com/yummysource/yummycli/internal/auth"
	internalimage "github.com/yummysource/yummycli/internal/image"
	"github.com/yummysource/yummycli/internal/output"
)

// Factory holds shared dependencies injected into commands.
type Factory struct {
	IOStreams        *IOStreams
	CredentialStore  *auth.ProviderCredentialStore
	ImageGenerator   internalimage.ImageGenerator
	Output           *output.Writer
}
