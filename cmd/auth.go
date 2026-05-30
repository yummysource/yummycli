package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yummysource/yummycli/internal/cmdutil"
)

// NewCmdAuth creates the shared auth command group.
func NewCmdAuth(f *cmdutil.Factory) *cobra.Command {
	command := newGroupCommand(
		"auth",
		"Provider and credential management",
		nil,
		newCmdAuthInit(f),
		newCmdAuthList(f),
		newCmdAuthStatus(f),
		newCmdAuthRemove(f),
	)

	return command
}

type authStatusOptions struct {
	Provider string
}

type authStatusResult struct {
	Provider      string `json:"provider"`
	Configured    bool   `json:"configured"`
	APIKeyPreview string `json:"apiKeyPreview,omitempty"`
}

type authInitOptions struct {
	Provider string
	APIKey   string
}

type authInitResult struct {
	Provider   string `json:"provider"`
	Configured bool   `json:"configured"`
}

func newCmdAuthStatus(f *cmdutil.Factory) *cobra.Command {
	opts := &authStatusOptions{}
	command := &cobra.Command{
		Use:   "status",
		Short: "Show credentials status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus(f, opts)
		},
	}
	command.Flags().StringVar(&opts.Provider, "provider", "", "provider name")
	if err := command.MarkFlagRequired("provider"); err != nil {
		panic(err)
	}
	return command
}

func runStatus(f *cmdutil.Factory, opts *authStatusOptions) error {
	if f.Output == nil {
		return fmt.Errorf("output writer is not configured")
	}
	configured, err := f.CredentialStore.HasAPIKey(opts.Provider)
	if err != nil {
		return err
	}

	result := authStatusResult{
		Provider:   opts.Provider,
		Configured: configured,
	}

	if configured {
		preview, err := f.CredentialStore.APIKeyPreview(opts.Provider)
		if err != nil {
			return err
		}
		result.APIKeyPreview = preview
	}

	return f.Output.JSON(result)
}

func newCmdAuthInit(f *cmdutil.Factory) *cobra.Command {
	opts := &authInitOptions{}
	command := &cobra.Command{
		Use:   "init",
		Short: "Initialize the provider credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAuthInit(f, opts)
		},
	}

	command.Flags().StringVar(&opts.Provider, "provider", "", "provider name")
	command.Flags().StringVar(&opts.APIKey, "api-key", "", "api key")

	if err := command.MarkFlagRequired("provider"); err != nil {
		panic(err)
	}
	if err := command.MarkFlagRequired("api-key"); err != nil {
		panic(err)
	}

	return command
}

func runAuthInit(f *cmdutil.Factory, opts *authInitOptions) error {
	if f.Output == nil {
		return fmt.Errorf("output writer is not configured")
	}
	err := f.CredentialStore.SaveAPIKey(opts.Provider, opts.APIKey)
	if err != nil {
		return err
	}

	result := &authInitResult{
		Provider:   opts.Provider,
		Configured: true,
	}
	return f.Output.JSON(result)
}

// authListResult is a single entry in the auth list JSON output.
type authListResult struct {
	Provider      string `json:"provider"`
	Configured    bool   `json:"configured"`
	Default       bool   `json:"default"`
	APIKeyPreview string `json:"apiKeyPreview,omitempty"`
}

func newCmdAuthList(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all providers and their credential status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAuthList(f)
		},
	}
}

func runAuthList(f *cmdutil.Factory) error {
	if f.Output == nil {
		return fmt.Errorf("output writer is not configured")
	}

	statuses, err := f.CredentialStore.ListConfigured()
	if err != nil {
		return err
	}

	defaultProvider, err := f.CredentialStore.GetDefaultProvider()
	if err != nil {
		return err
	}

	results := make([]authListResult, 0, len(statuses))
	for _, s := range statuses {
		results = append(results, authListResult{
			Provider:      s.Provider,
			Configured:    s.Configured,
			Default:       s.Provider == defaultProvider,
			APIKeyPreview: s.APIKeyPreview,
		})
	}

	return f.Output.JSON(results)
}

type authRemoveOptions struct {
	Provider string
}

type authRemoveResult struct {
	Provider string `json:"provider"`
	Removed  bool   `json:"removed"`
}

func newCmdAuthRemove(f *cmdutil.Factory) *cobra.Command {
	opts := &authRemoveOptions{}
	command := &cobra.Command{
		Use:   "remove",
		Short: "Remove stored provider credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAuthRemove(f, opts)
		},
	}

	command.Flags().StringVar(&opts.Provider, "provider", "", "provider name")
	if err := command.MarkFlagRequired("provider"); err != nil {
		panic(err)
	}

	return command
}

func runAuthRemove(f *cmdutil.Factory, opts *authRemoveOptions) error {
	if f.Output == nil {
		return fmt.Errorf("output writer is not configured")
	}
	err := f.CredentialStore.RemoveAPIKey(opts.Provider)
	if err != nil {
		return err
	}
	result := &authRemoveResult{
		Provider: opts.Provider,
		Removed:  true,
	}

	return f.Output.JSON(result)
}
