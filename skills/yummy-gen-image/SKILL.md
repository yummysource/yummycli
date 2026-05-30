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
| `--aspect-ratio` | Gemini only: named ratio e.g. `16:9`, `21:9`, `9:16`, `1:1` |
| `--image-size` | Gemini: `1K` `2K` `4K` — OpenAI: WxH e.g. `1536x864` (both dims multiples of 16) |
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

### Aspect ratio

Both providers support aspect ratio — Gemini via `--aspect-ratio`, OpenAI via `--image-size` (WxH, both dims multiples of 16).

When the user states a ratio, use the **configured default provider** and translate accordingly:

| User intent | Gemini | OpenAI |
|---|---|---|
| Ultrawide / cinematic (21:9) | `--aspect-ratio 21:9` | `--image-size 2016x864` |
| Widescreen / landscape (16:9) | `--aspect-ratio 16:9` | `--image-size 1536x864` (default) |
| Standard landscape (4:3) | `--aspect-ratio 4:3` | `--image-size 1024x768` |
| Square (1:1) | `--aspect-ratio 1:1` | `--image-size 1024x1024` |
| Portrait (3:4) | `--aspect-ratio 3:4` | `--image-size 768x1024` |
| Vertical / portrait (9:16) | `--aspect-ratio 9:16` | `--image-size 864x1536` |

Do NOT switch providers just because a ratio was requested — both providers handle all common ratios.

**Image-size guidance (OpenAI, when user specifies a size rather than a ratio):**
- User says "larger" or "high-res" → `--image-size 2048x1152`
- User says "smaller" or "draft" → `--image-size 1024x576`
- Otherwise → omit, use default (`1536x864`)

**Output format guidance (OpenAI):**
- User asks for JPEG or smaller file size → `--output-format jpeg`
- User asks for WebP → `--output-format webp`
- Otherwise → omit, defaults to `png`

**Quality guidance (OpenAI):**
- User asks for high quality / best quality → `--quality high`
- User asks for draft / fast / low quality → `--quality low`
- Otherwise → omit, let API decide
