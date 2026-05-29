package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yummysource/yummycli/internal/cmdutil"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(f, opts)
		},
	}

	command.Flags().StringVar(&opts.Provider, "provider", "", "provider name (gemini, openai)")
	command.Flags().StringVar(&opts.APIKey, "api-key", "", "API key for the provider")
	command.Flags().BoolVar(&opts.Default, "default", false, "set this provider as the default")

	if err := command.MarkFlagRequired("provider"); err != nil {
		panic(err)
	}
	if err := command.MarkFlagRequired("api-key"); err != nil {
		panic(err)
	}

	return command
}

func runInit(f *cmdutil.Factory, opts *initOptions) error {
	if f.Output == nil {
		return fmt.Errorf("output writer is not configured")
	}

	if err := f.CredentialStore.SaveAPIKey(opts.Provider, opts.APIKey); err != nil {
		return err
	}

	if opts.Default {
		if err := f.CredentialStore.SetDefaultProvider(opts.Provider); err != nil {
			return err
		}
	}

	return f.Output.JSON(initResult{
		Provider:   opts.Provider,
		Configured: true,
		Default:    opts.Default,
	})
}
