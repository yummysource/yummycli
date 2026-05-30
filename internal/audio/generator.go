// Package audio defines the Speaker interface and shared types used by all
// text-to-speech implementations. The interface mirrors the pattern used by
// internal/image and internal/video: a single-method interface with a request
// struct and a result struct, so commands stay simple and implementations are
// swappable for testing.
package audio

import "context"

// Speaker synthesises speech from text and writes the result to a WAV file.
// Implementations handle provider-specific API calls and audio encoding.
type Speaker interface {
	Speak(ctx context.Context, req SpeakRequest) (*SpeakResult, error)
}

// SpeakRequest holds all parameters for a single text-to-speech call.
// The caller must perform all validation before passing this to Speak.
type SpeakRequest struct {
	// Provider identifies the AI backend, e.g. "gemini".
	Provider string

	// Text is the input text to synthesise. Required.
	Text string

	// Output is the destination file path. Must end in ".wav".
	// If empty, the generator generates a timestamped filename.
	Output string

	// Model is the provider-specific model ID.
	// Default: gemini-3.1-flash-tts-preview
	Model string

	// VoiceName is the prebuilt voice to use for single-speaker synthesis.
	// Mutually exclusive with Speakers. Default: Aoede
	VoiceName string

	// LanguageCode is an optional BCP-47 / ISO 639-1 language code.
	// Leave empty for automatic language detection from the text.
	LanguageCode string

	// Speakers configures multi-speaker dialogue synthesis (up to 2 speakers).
	// Mutually exclusive with VoiceName.
	Speakers []SpeakerConfig
}

// SpeakerConfig maps a speaker name (as used in the prompt text) to a voice.
// Used for multi-speaker dialogue: set --speaker Alice:Aoede --speaker Bob:Kore.
type SpeakerConfig struct {
	// Name is the speaker identifier used in the prompt, e.g. "Alice".
	Name string `json:"name"`

	// VoiceName is the prebuilt voice assigned to this speaker, e.g. "Aoede".
	VoiceName string `json:"voice"`
}

// SpeakResult holds metadata about a successfully generated audio file.
// Serialised to JSON on stdout by the command layer.
type SpeakResult struct {
	// Provider is the AI provider used, e.g. "gemini".
	Provider string `json:"provider"`

	// Output is the path of the written WAV file.
	Output string `json:"output"`

	// Model is the model that synthesised the audio.
	Model string `json:"model"`

	// Voice is the prebuilt voice used for single-speaker requests.
	Voice string `json:"voice,omitempty"`

	// Speakers is set for multi-speaker requests.
	Speakers []SpeakerConfig `json:"speakers,omitempty"`

	// ElapsedSeconds is the wall-clock time from request to file written.
	ElapsedSeconds int `json:"elapsed_seconds"`
}

// VoiceInfo describes a single prebuilt voice available from a provider.
type VoiceInfo struct {
	// Name is the voice identifier passed to --voice.
	Name string `json:"name"`

	// Style is a one-word English description of the voice character.
	Style string `json:"style"`
}
