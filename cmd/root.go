package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/yummysource/yummycli/internal/build"
	"github.com/yummysource/yummycli/internal/cmdutil"
)

const rootLong = `yummycli - AI-friendly CLI for multimodal model providers.

  Current Phase 1 scope:
    - Gemini auth initialization
    - Gemini Nano Banana image generation
    - Generic image capability commands

  This CLI uses a provider-first command surface with a capability-first internal architecture.`

// Execute runs the root command and returns the process exit code.
func Execute() int {
	f := cmdutil.NewDefault()

	rootCmd := NewRootCommand(f)

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(f.IOStreams.ErrOut, "Error:", err)
		return 1
	}

	return 0
}

// NewRootCommand constructs the root Cobra command for yummycli.
func NewRootCommand(f *cmdutil.Factory) *cobra.Command {
	root := &cobra.Command{
		Use:   "yummycli",
		Short: "AI-friendly CLI for multimodal model providers",
		Long:  rootLong,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true
			return nil
		},
	}

	root.AddCommand(
		newVersionCommand(f),
		NewCmdGemini(f),
		NewCmdImage(f),
		NewCmdVideo(f),
		NewCmdAudio(f),
		NewCmdAuth(f),
	)

	return root
}

func newGroupCommand(use, short string, annotations map[string]string, subcommands ...*cobra.Command) *cobra.Command {
	command := &cobra.Command{
		Use:         use,
		Short:       short,
		Annotations: annotations,
	}

	command.AddCommand(subcommands...)
	return command
}

func newVersionCommand(f *cmdutil.Factory) *cobra.Command {
	command := &cobra.Command{
		Use:   "version",
		Short: "Show the yummycli version",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintln(f.IOStreams.Out, build.Version)
			return err
		},
	}

	return command
}

func assertSubcommands(t interface {
	Fatalf(string, ...any)
}, cmd *cobra.Command, expected []string,
) {
	got := make([]string, 0, len(cmd.Commands()))
	for _, sub := range cmd.Commands() {
		if !sub.IsAvailableCommand() || sub.IsAdditionalHelpTopicCommand() {
			continue
		}
		got = append(got, sub.Name())
	}

	sort.Strings(got)
	sortedExpected := append([]string(nil), expected...)
	sort.Strings(sortedExpected)

	if len(got) != len(sortedExpected) {
		t.Fatalf("subcommands for %s = %v, want %v", cmd.Name(), got, sortedExpected)
	}

	for i := range got {
		if got[i] != sortedExpected[i] {
			t.Fatalf("subcommands for %s = %v, want %v", cmd.Name(), got, sortedExpected)
		}
	}
}
