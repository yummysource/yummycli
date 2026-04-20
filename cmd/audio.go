package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	internalaudio "github.com/yummysource/yummycli/internal/audio"
	"github.com/yummysource/yummycli/internal/cmdutil"
	"github.com/yummysource/yummycli/internal/providers"
)

// audioSpeakOptions is the canonical options struct for audio speech synthesis.
// Shared by `audio speak` (capability layer) and `gemini speak` (shortcut).
type audioSpeakOptions struct {
	Provider     string
	Text         string
	Output       string
	Model        string
	Voice        string
	LanguageCode string
	// Speakers holds raw "Name:VoiceName" strings from repeated --speaker flags.
	// Parsed into []SpeakerConfig before the API call.
	Speakers []string
}

// NewCmdAudio creates the provider-agnostic audio command group.
func NewCmdAudio(f *cmdutil.Factory) *cobra.Command {
	command := &cobra.Command{
		Use:   "audio",
		Short: "Synthesise speech and list voices",
	}

	command.AddCommand(
		newCmdAudioSpeak(f),
		newCmdAudioVoices(f),
	)

	return command
}

// newCmdAudioSpeak creates the `audio speak` subcommand.
func newCmdAudioSpeak(f *cmdutil.Factory) *cobra.Command {
	opts := &audioSpeakOptions{}

	command := &cobra.Command{
		Use:   "speak",
		Short: "Synthesise text to speech and save as a WAV file",
		Annotations: map[string]string{
			"capability": "audio.speak",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAudioSpeak(f, opts)
		},
	}

	command.Flags().StringVar(&opts.Provider, "provider", "", "AI provider (e.g. gemini)")
	command.Flags().StringVar(&opts.Text, "text", "", "Text to synthesise")
	command.Flags().StringVar(&opts.Output, "output", "", "Output WAV file path (auto-generated if omitted)")
	command.Flags().StringVar(&opts.Model, "model", "", "TTS model to use")
	command.Flags().StringVar(&opts.Voice, "voice", "", "Prebuilt voice name (e.g. Aoede); mutually exclusive with --speaker")
	command.Flags().StringVar(&opts.LanguageCode, "language", "", "BCP-47 language code (auto-detected from text if omitted)")
	command.Flags().StringArrayVar(&opts.Speakers, "speaker", nil, "Multi-speaker mapping Name:Voice; repeat up to 2 (e.g. --speaker Alice:Aoede)")

	if err := command.MarkFlagRequired("provider"); err != nil {
		panic(err)
	}
	if err := command.MarkFlagRequired("text"); err != nil {
		panic(err)
	}

	return command
}

// newCmdAudioVoices creates the `audio voices` subcommand.
func newCmdAudioVoices(f *cmdutil.Factory) *cobra.Command {
	var provider string

	command := &cobra.Command{
		Use:   "voices",
		Short: "List available prebuilt voices for a provider",
		Annotations: map[string]string{
			"capability": "audio.voices",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAudioVoices(f, provider)
		},
	}

	command.Flags().StringVar(&provider, "provider", "", "AI provider (e.g. gemini)")
	if err := command.MarkFlagRequired("provider"); err != nil {
		panic(err)
	}

	return command
}

// runAudioSpeak is the canonical implementation for speech synthesis.
// Shared by `audio speak` and the `gemini speak` shortcut.
func runAudioSpeak(f *cmdutil.Factory, opts *audioSpeakOptions) error {
	if f.Speaker == nil {
		return fmt.Errorf("speech synthesiser is not configured")
	}

	if err := applyAudioDefaults(opts); err != nil {
		return err
	}

	speakers, err := parseSpeakers(opts.Speakers)
	if err != nil {
		return err
	}
	if len(speakers) > 0 && opts.Voice != "" {
		return fmt.Errorf("--voice and --speaker are mutually exclusive")
	}
	if len(speakers) > 2 {
		return fmt.Errorf("at most 2 speakers are supported (got %d)", len(speakers))
	}

	if opts.Output == "" {
		opts.Output = defaultAudioOutputPath()
	} else if !strings.HasSuffix(strings.ToLower(opts.Output), ".wav") {
		return fmt.Errorf("output file must have a .wav extension (got %q)", opts.Output)
	}

	if _, statErr := os.Stat(opts.Output); statErr == nil {
		return fmt.Errorf("output file already exists: %s (remove it or choose a different path)", opts.Output)
	}

	if dir := filepath.Dir(opts.Output); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("creating output directory %s: %w", dir, err)
		}
	}

	req := internalaudio.SpeakRequest{
		Provider:     opts.Provider,
		Text:         opts.Text,
		Output:       opts.Output,
		Model:        opts.Model,
		VoiceName:    opts.Voice,
		LanguageCode: opts.LanguageCode,
		Speakers:     speakers,
	}

	result, err := f.Speaker.Speak(context.Background(), req)
	if err != nil {
		return err
	}

	return f.Output.JSON(result)
}

// runAudioVoices lists prebuilt voices for the given provider.
func runAudioVoices(f *cmdutil.Factory, provider string) error {
	switch provider {
	case providers.Gemini:
		return f.Output.JSON(internalaudio.GeminiVoices())
	default:
		return fmt.Errorf("unsupported provider: %s", provider)
	}
}

// applyAudioDefaults fills in provider-specific default values for empty fields.
func applyAudioDefaults(opts *audioSpeakOptions) error {
	switch opts.Provider {
	case providers.Gemini:
		if opts.Model == "" {
			opts.Model = internalaudio.DefaultTTSModel
		}
		if opts.Voice == "" && len(opts.Speakers) == 0 {
			opts.Voice = internalaudio.DefaultVoice
		}
	default:
		return fmt.Errorf("unsupported provider: %s", opts.Provider)
	}
	return nil
}

// parseSpeakers parses "Name:VoiceName" strings into SpeakerConfig values.
func parseSpeakers(raw []string) ([]internalaudio.SpeakerConfig, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	result := make([]internalaudio.SpeakerConfig, 0, len(raw))
	for _, s := range raw {
		parts := strings.SplitN(s, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid --speaker value %q: expected Name:VoiceName", s)
		}
		result = append(result, internalaudio.SpeakerConfig{
			Name:      parts[0],
			VoiceName: parts[1],
		})
	}
	return result, nil
}

// defaultAudioOutputPath generates a timestamped default output filename.
func defaultAudioOutputPath() string {
	t := time.Now()
	return fmt.Sprintf("tts_%s_%03d.wav", t.Format("20060102_150405"), t.Nanosecond()/1e6)
}
