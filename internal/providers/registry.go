package providers

import (
	"fmt"
	"strings"
)

// Gemini identifies the Gemini provider.
const Gemini = "gemini"

var supported = map[string]struct{}{
	Gemini: {},
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
