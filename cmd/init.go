package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yummysource/yummycli/internal/cmdutil"
	"github.com/yummysource/yummycli/internal/providers"
)

type initOptions struct {
	Provider string
	APIKey   string
	Default  bool
}

type initResult struct {
	Provider   string `json:"provider"`
	Configured bool   `json:"configured"`
	Default    bool   `json:"default"`
}

// NewCmdInit creates the top-level init command for configuring a provider API key.
func NewCmdInit(f *cmdutil.Factory) *cobra.Command {
	opts := &initOptions{}

	command := &cobra.Command{
		Use:   "init",
		Short: "Configure a provider API key",
		Long: `Save an API key for a provider and optionally set it as the default.

If the API key for the provider is already configured, --api-key can be omitted
and --default alone is sufficient to switch the default provider.

PROVIDERS
  gemini   Google Gemini (image, video, speech)
  openai   OpenAI (image only)

The --default flag sets this provider as the one used when --provider is
omitted from generation commands. If two providers are configured, the
non-default provider acts as automatic fallback when the primary fails.`,
		Example: `  # First-time setup: save key and set as default
  yummycli init --provider gemini --api-key <key> --default

  # Add a second provider as fallback (existing default unchanged)
  yummycli init --provider openai --api-key <key>

  # Switch default to a provider whose key is already saved
  yummycli init --provider openai --default`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(f, opts)
		},
	}

	command.Flags().StringVar(&opts.Provider, "provider", "", "provider name (gemini, openai)")
	command.Flags().StringVar(&opts.APIKey, "api-key", "", "API key for the provider (required if not already configured)")
	command.Flags().BoolVar(&opts.Default, "default", false, "set this provider as the default")

	if err := command.MarkFlagRequired("provider"); err != nil {
		panic(err)
	}

	return command
}

func runInit(f *cmdutil.Factory, opts *initOptions) error {
	if f.CredentialStore == nil {
		return fmt.Errorf("credential store is not configured")
	}
	if f.Output == nil {
		return fmt.Errorf("output writer is not configured")
	}

	normalized, err := providers.Normalize(opts.Provider)
	if err != nil {
		return err
	}

	if opts.APIKey != "" {
		if err := f.CredentialStore.SaveAPIKey(normalized, opts.APIKey); err != nil {
			return err
		}
	} else {
		// No key provided — require that one is already stored.
		configured, err := f.CredentialStore.HasAPIKey(normalized)
		if err != nil {
			return err
		}
		if !configured {
			return fmt.Errorf("no API key configured for %s; provide one with --api-key", normalized)
		}
	}

	if opts.Default {
		if err := f.CredentialStore.SetDefaultProvider(normalized); err != nil {
			return err
		}
	}

	return f.Output.JSON(initResult{
		Provider:   normalized,
		Configured: true,
		Default:    opts.Default,
	})
}
