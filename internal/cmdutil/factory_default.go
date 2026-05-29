package cmdutil

import (
	"os"

	"github.com/yummysource/yummycli/internal/audio"
	"github.com/yummysource/yummycli/internal/auth"
	internalimage "github.com/yummysource/yummycli/internal/image"
	"github.com/yummysource/yummycli/internal/output"
	"github.com/yummysource/yummycli/internal/providers"
	"github.com/yummysource/yummycli/internal/video"
	"golang.org/x/term"
)

// NewDefault creates a production Factory with the process IO streams.
func NewDefault() *Factory {
	secretStore := auth.NewKeychainSecretStore()
	credentialStore := auth.NewProviderCredentialStore(secretStore)

	// Detect whether stdout is a real terminal so progress output can be enabled.
	isTerminal := term.IsTerminal(int(os.Stdout.Fd()))

	streams := &IOStreams{
		In:         os.Stdin,
		Out:        os.Stdout,
		ErrOut:     os.Stderr,
		IsTerminal: isTerminal,
	}
	imageGenerators := map[string]internalimage.ImageGenerator{
		providers.Gemini: internalimage.NewGeminiGenerator(credentialStore),
		providers.OpenAI: internalimage.NewOpenAIGenerator(credentialStore, ""),
	}
	imageGenerator := internalimage.NewMultiGenerator(imageGenerators, streams.ErrOut)

	return &Factory{
		IOStreams:       streams,
		CredentialStore: credentialStore,
		ImageGenerator:  imageGenerator,
		VideoGenerator:  video.NewGeminiGenerator(credentialStore),
		Speaker:         audio.NewGeminiSpeaker(credentialStore),
		Output:          output.New(streams.Out),
		Progress:        NewProgress(streams.ErrOut, isTerminal),
	}
}
