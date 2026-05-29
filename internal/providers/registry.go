package providers

import (
	"fmt"
	"sort"
	"strings"
)

// Gemini identifies the Gemini provider.
const Gemini = "gemini"

// OpenAI identifies the OpenAI provider.
const OpenAI = "openai"

var supported = map[string]struct{}{
	Gemini: {},
	OpenAI: {},
}

// All returns the names of all registered providers in sorted order.
func All() []string {
	names := make([]string, 0, len(supported))
	for name := range supported {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Normalize validates and normalizes a provider name.
func Normalize(name string) (string, error) {
	normalized := strings.TrimSpace(strings.ToLower(name))
	if normalized == "" {
		return "", fmt.Errorf("provider is required")
	}
	if _, ok := supported[normalized]; !ok {
		return "", fmt.Errorf("unsupported provider: %s", name)
	}

	return normalized, nil
}
