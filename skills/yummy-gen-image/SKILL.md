---
name: yummy-gen-image
version: 2.1.0
description: "Use when the user wants to generate or edit raster images through yummycli, including prompt-only generation, single-image editing, and multi-image reference editing. Supports both Gemini and OpenAI providers."
metadata:
  requires:
    bins: ["yummycli"]
    skills: ["yummy-shared"]
  openclaw:
    requires:
      bins: ["yummycli"]
    related_skills: ["yummy-shared"]
  hermes:
    tags: [image, generation, editing, multimodal, gemini, openai]
    related_skills: ["yummy-shared"]
install:
  - kind: node
    package: "@yummysource/yummycli"
    bins: ["yummycli"]
---

# Generate Image

Create or edit images with `yummycli image generate`.

Supported providers: **Gemini** and **OpenAI**. Provider is resolved from `yummycli` config.

## When to Use

Load this skill when the user asks to generate, create, or edit an image — including text-to-image, style transfer, image editing with a reference photo, or multi-image compositing.

> **Prerequisite:** Apply the `yummy-shared` skill first.

## Command

```bash
yummycli image generate --prompt "<prompt>"
```

The provider is resolved automatically from `yummycli` config. Omit `--provider` unless the user explicitly requests a specific provider.

Add reference images when needed (triggers image editing mode):

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
| `--aspect-ratio` | Gemini only: e.g. `16:9`, `9:16`, `1:1` |
| `--image-size` | Gemini: `1K` `2K` `4K` — OpenAI: `1536x864` `1024x1024` `1536x1024` `1024x1536` |
| `--quality` | OpenAI only: `low` `medium` `high` `auto` |
| `--output-format` | OpenAI only: `png` (default) `jpeg` `webp` |
| `--input-image` | Input image path for editing; repeat for multiple images |

## Provider-Specific Defaults

**Gemini** (when gemini is the active provider):
- Model: `gemini-3.1-flash-image`
- Aspect ratio: `16:9`
- Image size: `1K`

**OpenAI** (when openai is the active provider):
- Model: `gpt-image-2`
- Size: `1536x864` (16:9)
- Quality: API default (auto)
- Output format: `png`

## Model Selection

### Gemini

- User says `gemini pro` or `pro` → `--model gemini-3-pro-image-preview`
- User says `gemini flash` or `flash` → `--model gemini-3.1-flash-image`
- No explicit model request → omit `--model`, let yummycli use its default

### OpenAI

Supported models: `gpt-image-2` (default), `gpt-5.5`

- User says `gpt-image-2` or `gpt image 2` → `--model gpt-image-2`
- User says `gpt-5.5` or `gpt 5.5` → `--model gpt-5.5`
- User mentions any other model name → pass it through; yummycli will warn and fall back to `gpt-image-2` automatically
- No explicit model request → omit `--model`, let yummycli use its default (`gpt-image-2`)

## Intent to Parameters

**Aspect-ratio guidance (Gemini only):**
- Vertical / portrait / phone wallpaper → `--aspect-ratio 9:16`
- Widescreen / horizontal / desktop → `--aspect-ratio 16:9`
- Square / avatar → `--aspect-ratio 1:1`
- User provides a specific ratio → pass through directly
- Shape unclear → omit

**Image-size guidance (OpenAI):**
- Landscape / horizontal → `--image-size 1536x864` (default, 16:9)
- Portrait / vertical → `--image-size 1024x1536`
- Square → `--image-size 1024x1024`

**Output format guidance (OpenAI):**
- User asks for JPEG or smaller file size → `--output-format jpeg`
- User asks for WebP → `--output-format webp`
- Otherwise → omit, defaults to `png`

**Quality guidance (OpenAI):**
- User asks for high quality / best quality → `--quality high`
- User asks for draft / fast / low quality → `--quality low`
- Otherwise → omit, let API decide
