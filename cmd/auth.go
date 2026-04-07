package cmd

import (
	"encoding/json"

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
		newSimpleCommand("list", "List configured provider credentials", nil),
		newCmdAuthStatus(f),
		newCmdAuthRemove(f),
		newSimpleCommand("clear", "Clear all provider credentials", nil),
	)

	return command
}

type authStatusOptions struct {
	Provider   string `json:"provider"`
	Configured bool   `json:"configured"`
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
	return command
}

func runStatus(f *cmdutil.Factory, opts *authStatusOptions) error {
	configured, err := f.CredentialStore.HasAPIKey(opts.Provider)
	if err != nil {
		return err
	}

	result := authStatusOptions{
		Provider:   opts.Provider,
		Configured: configured,
	}

	return json.NewEncoder(f.IOStreams.Out).Encode(result)
}

func newCmdAuthInit(f *cmdutil.Factory) *cobra.Command {
	opts := &authInitOptions{}
	command := &cobra.Command{
		Use:   "init",
		Short: "Initilize the provider credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAuthInit(f, opts)
		},
	}

	command.Flags().StringVar(&opts.Provider, "provider", "", "provider name")
	command.Flags().StringVar(&opts.APIKey, "api-key", "", "api key")

	return command
}

func runAuthInit(f *cmdutil.Factory, opts *authInitOptions) error {
	err := f.CredentialStore.SaveAPIKey(opts.Provider, opts.APIKey)
	if err != nil {
		return err
	}

	result := &authInitResult{
		Provider:   opts.Provider,
		Configured: true,
	}
	return json.NewEncoder(f.IOStreams.Out).Encode(result)
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
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAuthRemove(f, opts)
		},
	}

	command.Flags().StringVar(&opts.Provider, "provider", "", "provider name")

	return command
}

func runAuthRemove(f *cmdutil.Factory, opts *authRemoveOptions) error {
	err := f.CredentialStore.RemoveAPIKey(opts.Provider)
	if err != nil {
		return err
	}
	result := &authRemoveResult{
		Provider: opts.Provider,
		Removed:  true,
	}

	return json.NewEncoder(f.IOStreams.Out).Encode(result)
}
