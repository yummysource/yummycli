---
name: yummy-gen-image
version: 2.0.0
description: "Use when the user wants to generate or edit raster images through yummycli, including prompt-only generation, single-image editing, and multi-image reference editing. Provider is resolved from yummycli config."
metadata:
  requires:
    bins: ["yummycli"]
    skills: ["yummy-shared"]
  openclaw:
    requires:
      bins: ["yummycli"]
    related_skills: ["yummy-shared"]
  hermes:
    tags: [image, generation, editing, multimodal]
    related_skills: ["yummy-shared"]
install:
  - kind: node
    package: "@yummysource/yummycli"
    bins: ["yummycli"]
---

# Generate Image

Create or edit images with `yummycli image generate`.

## When to Use

Load this skill when the user asks to generate, create, or edit an image — including text-to-image, style transfer, image editing with a reference photo, or multi-image compositing.

> **Prerequisite:** Apply the `yummy-shared` skill first.

## Command

```bash
yummycli image generate --prompt "<prompt>"
```

The provider is resolved automatically from `yummycli` config. Omit `--provider` unless the user explicitly requests a specific provider.

Add reference images when needed:

```bash
yummycli image generate \
  --prompt "<prompt>" \
  --input-image ./source-a.png \
  --input-image ./source-b.png
```

## Optional Flags

| Flag | Description |
|------|-------------|
| `--provider` | Override the configured default (`gemini`, `openai`) |
| `--model` | Specific model name |
| `--output` | Output file path (auto-generated if omitted) |
| `--aspect-ratio` | Gemini: e.g. `16:9`, `9:16`, `1:1` |
| `--image-size` | Gemini: `1K` `2K` `4K` — OpenAI: `1024x1024` `1024x1792` `1792x1024` |
| `--quality` | OpenAI only: `standard` (default) or `hd` |
| `--style` | OpenAI only: `vivid` (default) or `natural` |

## Provider-Specific Defaults

**Gemini** (when gemini is the active provider):
- Model: `gemini-3.1-flash-image-preview`
- Aspect ratio: `16:9`
- Image size: `1K`

**OpenAI** (when openai is the active provider):
- Model: `dall-e-3`
- Size: `1024x1024`
- Quality: `standard`
- Style: `vivid`

## Model Selection (Gemini)

- User says `gemini pro` or `pro` → `--model gemini-3-pro-image-preview`
- User says `gemini flash` or `flash` → `--model gemini-3.1-flash-image-preview`
- No explicit model request → omit `--model` and let yummycli use its default

## Intent to Parameters

Aspect-ratio guidance (Gemini):
- Vertical/portrait/phone → `--aspect-ratio 9:16`
- Widescreen/horizontal/desktop → `--aspect-ratio 16:9`
- Square/avatar → `--aspect-ratio 1:1`
- User provides a specific ratio → pass through directly
- Shape unclear → omit

Image-size guidance:
- Explicit 4K → `--image-size 4K` (Gemini) or not applicable for OpenAI (use default)
- Explicit HD quality → `--quality hd` (OpenAI only)
