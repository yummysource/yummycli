// Package video defines the VideoGenerator interface and supporting types
// for AI-powered video generation.
//
// Design: mirrors internal/image exactly — a provider-agnostic interface
// with one concrete implementation per provider. The interface is single-step
// (GenerateVideo blocks until the video is ready), so callers do not need to
// manage polling or download logic.
//
// Phase 1 supports text-to-video via Google Gemini Veo.
// Future phases may add image-to-video and multi-provider support.
package video
