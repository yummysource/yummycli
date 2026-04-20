package audio

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/genai"

	"github.com/yummysource/yummycli/internal/auth"
	"github.com/yummysource/yummycli/internal/providers"
)

// DefaultTTSModel is used when the caller does not specify a model.
const DefaultTTSModel = "gemini-3.1-flash-tts-preview"

// DefaultVoice is the prebuilt voice used when --voice is not specified.
const DefaultVoice = "Aoede"

// GeminiSpeaker synthesises speech using the Google Gemini TTS API.
// It calls GenerateContent with ResponseModalities=["AUDIO"] and writes
// the resulting PCM bytes to a WAV file after prepending the RIFF header.
type GeminiSpeaker struct {
	// credentialStore provides the API key for the Gemini provider.
	credentialStore *auth.ProviderCredentialStore
}

// NewGeminiSpeaker creates a GeminiSpeaker backed by the given credential store.
func NewGeminiSpeaker(credentialStore *auth.ProviderCredentialStore) *GeminiSpeaker {
	return &GeminiSpeaker{credentialStore: credentialStore}
}

// Speak synthesises req.Text to speech and writes the WAV file to req.Output.
// For multi-speaker requests (req.Speakers non-empty), MultiSpeakerVoiceConfig
// is used; otherwise a single PrebuiltVoiceConfig is applied.
func (g *GeminiSpeaker) Speak(ctx context.Context, req SpeakRequest) (*SpeakResult, error) {
	if req.Provider != providers.Gemini {
		return nil, fmt.Errorf("unsupported provider: %s", req.Provider)
	}

	apiKey, err := g.credentialStore.GetAPIKey(req.Provider)
	if err != nil {
		return nil, err
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		return nil, fmt.Errorf("creating Gemini client: %w", err)
	}

	config := &genai.GenerateContentConfig{
		ResponseModalities: []string{"AUDIO"},
		SpeechConfig:       buildSpeechConfig(req),
	}

	startTime := time.Now()
	result, err := client.Models.GenerateContent(ctx, req.Model, genai.Text(req.Text), config)
	if err != nil {
		return nil, fmt.Errorf("generating speech: %w", err)
	}

	// Collect all PCM audio chunks from the response parts.
	// The API may split the audio into multiple inline parts.
	pcm, err := extractPCM(result)
	if err != nil {
		return nil, err
	}

	wav := writeWAVHeader(pcm)
	if err := os.WriteFile(req.Output, wav, 0o644); err != nil {
		return nil, fmt.Errorf("writing audio to %s: %w", req.Output, err)
	}

	elapsed := int(time.Since(startTime).Seconds())
	out := &SpeakResult{
		Provider:       req.Provider,
		Output:         req.Output,
		Model:          req.Model,
		ElapsedSeconds: elapsed,
	}
	if len(req.Speakers) > 0 {
		out.Speakers = req.Speakers
	} else {
		out.Voice = req.VoiceName
	}
	return out, nil
}

// buildSpeechConfig constructs the SpeechConfig for either single-speaker or
// multi-speaker mode based on the presence of req.Speakers.
func buildSpeechConfig(req SpeakRequest) *genai.SpeechConfig {
	cfg := &genai.SpeechConfig{}

	if req.LanguageCode != "" {
		cfg.LanguageCode = req.LanguageCode
	}

	if len(req.Speakers) > 0 {
		// Multi-speaker: map each SpeakerConfig to a SpeakerVoiceConfig.
		svcs := make([]*genai.SpeakerVoiceConfig, 0, len(req.Speakers))
		for _, s := range req.Speakers {
			svcs = append(svcs, &genai.SpeakerVoiceConfig{
				Speaker: s.Name,
				VoiceConfig: &genai.VoiceConfig{
					PrebuiltVoiceConfig: &genai.PrebuiltVoiceConfig{
						VoiceName: s.VoiceName,
					},
				},
			})
		}
		cfg.MultiSpeakerVoiceConfig = &genai.MultiSpeakerVoiceConfig{
			SpeakerVoiceConfigs: svcs,
		}
	} else {
		cfg.VoiceConfig = &genai.VoiceConfig{
			PrebuiltVoiceConfig: &genai.PrebuiltVoiceConfig{
				VoiceName: req.VoiceName,
			},
		}
	}

	return cfg
}

// extractPCM collects raw PCM bytes from all InlineData parts in the response.
// The API typically returns a single audio part, but may return multiple chunks.
func extractPCM(result *genai.GenerateContentResponse) ([]byte, error) {
	if result == nil || len(result.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates in TTS response")
	}
	cand := result.Candidates[0]
	if cand.Content == nil {
		return nil, fmt.Errorf("empty content in TTS response")
	}

	var pcm []byte
	for _, part := range cand.Content.Parts {
		if part.InlineData != nil && len(part.InlineData.Data) > 0 {
			pcm = append(pcm, part.InlineData.Data...)
		}
	}
	if len(pcm) == 0 {
		return nil, fmt.Errorf("TTS response contained no audio data")
	}
	return pcm, nil
}

// GeminiVoices returns the 30 prebuilt voices available for Gemini TTS.
// This list is static and requires no API call.
func GeminiVoices() []VoiceInfo {
	return []VoiceInfo{
		{Name: "Zephyr", Style: "Bright"},
		{Name: "Puck", Style: "Upbeat"},
		{Name: "Charon", Style: "Informative"},
		{Name: "Kore", Style: "Firm"},
		{Name: "Fenrir", Style: "Excitable"},
		{Name: "Leda", Style: "Youthful"},
		{Name: "Orus", Style: "Firm"},
		{Name: "Aoede", Style: "Breezy"},
		{Name: "Callirrhoe", Style: "Easy-going"},
		{Name: "Autonoe", Style: "Bright"},
		{Name: "Enceladus", Style: "Breathy"},
		{Name: "Iapetus", Style: "Clear"},
		{Name: "Umbriel", Style: "Easy-going"},
		{Name: "Algieba", Style: "Smooth"},
		{Name: "Despina", Style: "Smooth"},
		{Name: "Erinome", Style: "Clear"},
		{Name: "Algenib", Style: "Gravelly"},
		{Name: "Rasalghul", Style: "Informative"},
		{Name: "Achird", Style: "Friendly"},
		{Name: "Zubenelgenubi", Style: "Casual"},
		{Name: "Vindemiatrix", Style: "Gentle"},
		{Name: "Sadachbia", Style: "Lively"},
		{Name: "Sadaltager", Style: "Knowledgeable"},
		{Name: "Sulafat", Style: "Warm"},
		{Name: "Schedar", Style: "Even"},
		{Name: "Gacrux", Style: "Mature"},
		{Name: "Pulcherrima", Style: "Forward"},
		{Name: "Laomedeia", Style: "Upbeat"},
		{Name: "Achernar", Style: "Soft"},
		{Name: "Alnilam", Style: "Firm"},
	}
}
