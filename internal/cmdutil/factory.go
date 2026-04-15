package cmdutil

import (
	"github.com/yummysource/yummycli/internal/auth"
	internalimage "github.com/yummysource/yummycli/internal/image"
	"github.com/yummysource/yummycli/internal/output"
	internalvideo "github.com/yummysource/yummycli/internal/video"
)

// Factory holds shared dependencies injected into commands.
// Each capability (image, video, …) has its own generator field so commands
// receive exactly the dependency they need, keeping them testable in isolation.
type Factory struct {
	// IOStreams provides process-level I/O (stdin, stdout, stderr) and terminal detection.
	IOStreams *IOStreams

	// CredentialStore manages API keys for all providers.
	CredentialStore *auth.ProviderCredentialStore

	// ImageGenerator generates images; injected so commands are testable.
	ImageGenerator internalimage.ImageGenerator

	// VideoGenerator generates videos; injected so commands are testable.
	VideoGenerator internalvideo.VideoGenerator

	// Output writes structured results (JSON) to stdout.
	Output *output.Writer

	// Progress displays status during long-running operations.
	// No-op in non-terminal environments so piped output is not polluted.
	Progress Progress
}
