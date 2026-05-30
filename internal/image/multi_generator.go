package image

import (
	"context"
	"fmt"
	"io"
)

// MultiGenerator routes image generation requests to the appropriate provider generator.
// If the primary provider call fails and a Fallback provider is set in the request,
// it retries with the fallback generator and logs a one-line message to errLog.
type MultiGenerator struct {
	generators map[string]ImageGenerator
	errLog     io.Writer
}

// NewMultiGenerator creates a MultiGenerator.
// errLog receives fallback warning messages; pass nil to discard them.
func NewMultiGenerator(generators map[string]ImageGenerator, errLog io.Writer) *MultiGenerator {
	if errLog == nil {
		errLog = io.Discard
	}
	return &MultiGenerator{
		generators: generators,
		errLog:     errLog,
	}
}

// GenerateImage dispatches the request to the provider named in req.Provider.
// If that call fails and req.Fallback is set, it retries with the fallback provider.
func (m *MultiGenerator) GenerateImage(ctx context.Context, req GenerateImageRequest) error {
	gen, ok := m.generators[req.Provider]
	if !ok {
		return fmt.Errorf("no generator registered for provider: %s", req.Provider)
	}

	err := gen.GenerateImage(ctx, req)
	if err == nil {
		return nil
	}

	if req.Fallback == "" {
		return err
	}

	fallbackGen, ok := m.generators[req.Fallback]
	if !ok {
		return fmt.Errorf("primary provider %s failed (%v); no generator for fallback provider %s", req.Provider, err, req.Fallback)
	}

	fmt.Fprintf(m.errLog, "primary provider %s failed (%v), retrying with %s\n", req.Provider, err, req.Fallback)

	// Clear all provider-specific fields so the fallback uses its own defaults.
	fallbackReq := req
	fallbackReq.Provider = req.Fallback
	fallbackReq.Fallback = ""
	fallbackReq.Model = ""
	fallbackReq.ImageSize = ""
	fallbackReq.AspectRatio = ""
	fallbackReq.Quality = ""
	return fallbackGen.GenerateImage(ctx, fallbackReq)
}
