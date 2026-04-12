# yummycli

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23-blue.svg)](https://go.dev/)
[![npm version](https://img.shields.io/npm/v/@yummysource/yummycli.svg)](https://www.npmjs.com/package/@yummysource/yummycli)

[繁體中文](./README.zh-TW.md) | [简体中文](./README.zh.md) | [English](./README.md)

<img src="./assets/logo.png" alt="yummycli logo" width="120" />

An AI-friendly CLI for multimodal model providers — built for humans and AI Agents.

Supports image generation and editing via [Gemini](https://deepmind.google/technologies/gemini/), with more providers (Claude, OpenAI, Qwen) planned.

[Install](#installation) · [Auth](#authentication) · [Image Generation](#image-generation) · [Agent Skills](#agent-skills) · [Commands](#command-reference)

---

## Why yummycli?

- **Agent-Native Design** — Structured [Skills](./skills/) out of the box, designed so AI Agents can call image APIs with zero extra setup
- **Capability-First Architecture** — `image generate` is the stable automation contract; `gemini nanobanana` is a human-friendly shortcut on top
- **Structured JSON Output** — Every command writes JSON to stdout, making it trivial to pipe into agents, scripts, or other tools
- **Secure Credential Storage** — API keys stored in the OS-native keychain (macOS Keychain, Linux Secret Service), never in plain text
- **Provider-Agnostic** — One CLI surface across providers; add a new provider without changing your scripts

---

## Installation

### Requirements

- Node.js 16+ with `npm`
- Go 1.23+ and `make` (only required when building from source)

### From npm (recommended)

```bash
# Install CLI
npm install -g @yummysource/yummycli

# Install Agent Skills (required for AI Agent usage)
npx skills add yummysource/yummycli -y -g
```

Verify the install:

```bash
yummycli version
```

### From source

```bash
git clone https://github.com/yummysource/yummycli.git
cd yummycli
make install

# Install Agent Skills (required for AI Agent usage)
npx skills add yummysource/yummycli -y -g
```

---

## Authentication

yummycli stores API keys per provider in the OS keychain. You only need to do this once per provider.

### Commands

| Command | Description |
|---------|-------------|
| `auth init` | Save an API key for a provider |
| `auth list` | List all providers and their credential status |
| `auth status` | Show credential status for a specific provider |
| `auth remove` | Delete stored credentials for a provider |

### Examples

```bash
# Save a Gemini API key
yummycli auth init --provider gemini --api-key "AIza..."

# Check if Gemini is configured (shows masked key preview)
yummycli auth status --provider gemini

# Remove Gemini credentials
yummycli auth remove --provider gemini
```

**Output from `auth init`:**

```json
{"provider":"gemini","configured":true}
```

**Output from `auth list`:**

```json
[{"provider":"gemini","configured":true,"apiKeyPreview":"AIza...xxxx"}]
```

**Output from `auth status`:**

```json
{"provider":"gemini","configured":true,"apiKeyPreview":"AIza...xxxx"}
```

### Gemini shortcut

`gemini init` is a provider-scoped alias for `auth init --provider gemini`:

```bash
yummycli gemini init --api-key "AIza..."
```

---

## Image Generation

yummycli supports two equivalent entry points for image generation:

| Entry point | Intended for |
|-------------|-------------|
| `gemini nanobanana` | Human use — Gemini-specific defaults pre-applied |
| `image generate --provider gemini` | Automation / scripting — explicit, stable contract |

Both call the same underlying implementation. Use whichever fits your context.

### Quick start

```bash
# Step 1: configure Gemini credentials (one-time)
yummycli gemini init --api-key "AIza..."

# Step 2: generate an image from a text prompt
yummycli gemini nanobanana --prompt "A ripe banana on a white plate, studio lighting"
```

The generated image is saved to the current directory with an auto-generated filename:

```
gemini_20260410123456_789.png
```

### Flags

| Flag | Description | Default (Gemini) |
|------|-------------|-----------------|
| `--prompt` | Image generation prompt (**required**) | — |
| `--output` | Output file path | auto-generated |
| `--model` | Gemini model | `gemini-3.1-flash-image-preview` |
| `--aspect-ratio` | Image aspect ratio | `16:9` |
| `--image-size` | Output resolution | `1K` |
| `--input-image` | Input image for editing (repeatable) | — |

> For `image generate`, also pass `--provider gemini` (required). The same Gemini defaults apply — `--model`, `--aspect-ratio`, and `--image-size` are filled in automatically when omitted.

### Text-to-image generation

```bash
yummycli gemini nanobanana \
  --prompt "A cyberpunk cityscape at night, neon reflections on wet streets"

# Specify output path and resolution
yummycli gemini nanobanana \
  --prompt "A minimalist logo, flat design, white background" \
  --output logo.png \
  --image-size 4K
```

### Image editing

Pass one or more reference images with `--input-image`:

```bash
# Single-image edit
yummycli gemini nanobanana \
  --prompt "Turn this into a watercolor illustration" \
  --input-image ./photo.png

# Multi-image reference
yummycli gemini nanobanana \
  --prompt "Blend these two references into a single polished poster" \
  --input-image ./subject.png \
  --input-image ./background.jpg
```

Supported input formats: `.png`, `.jpg` / `.jpeg`, `.webp`.

### Aspect ratio

```bash
# Vertical — phone wallpaper, story format
yummycli gemini nanobanana --prompt "..." --aspect-ratio 9:16

# Square — social avatar, icon
yummycli gemini nanobanana --prompt "..." --aspect-ratio 1:1

# Widescreen — desktop wallpaper, presentation banner
yummycli gemini nanobanana --prompt "..." --aspect-ratio 21:9
```

### Model selection

```bash
# Flash (default) — faster, supports more aspect ratios and smaller sizes
yummycli gemini nanobanana --prompt "..." --model gemini-3.1-flash-image-preview

# Pro — higher quality, fewer size/ratio options
yummycli gemini nanobanana --prompt "..." --model gemini-3-pro-image-preview
```

### Model compatibility

**Aspect ratio**

| Model | Supported values |
|-------|-----------------|
| `gemini-3.1-flash-image-preview` | `1:1` `1:4` `1:8` `2:3` `3:2` `3:4` `4:1` `4:3` `4:5` `5:4` `8:1` `9:16` `16:9` `21:9` |
| `gemini-3-pro-image-preview` | `1:1` `2:3` `3:2` `3:4` `4:3` `4:5` `5:4` `9:16` `16:9` `21:9` |

`1:4`, `1:8`, `4:1`, `8:1` are Flash-only and not supported by the Pro model.

**Image size**

| Model | Supported values |
|-------|-----------------|
| `gemini-3.1-flash-image-preview` | `512` `0.5K` `1K` `2K` `4K` |
| `gemini-3-pro-image-preview` | `1K` `2K` `4K` |

`512` and `0.5K` are Flash-only. Size values are case-insensitive (`4k` and `4K` both work).

### JSON output

Every successful generation writes a result to stdout:

```json
{
  "provider": "gemini",
  "output": "gemini_20260410123456_789.png",
  "model": "gemini-3.1-flash-image-preview",
  "inputImageCount": 0
}
```

Use the `output` field to locate the generated file.

### Using `image generate` directly

`image generate` is the provider-agnostic stable API. It accepts the same flags but requires an explicit `--provider`:

```bash
yummycli image generate \
  --provider gemini \
  --prompt "A serene mountain lake at sunrise" \
  --aspect-ratio 16:9 \
  --image-size 2K \
  --output landscape.png
```

This form is recommended for scripts and AI Agents — it will continue to work unchanged as new providers are added.

---

## Agent Skills

yummycli ships with Skills — structured instruction files that teach AI Agents how to use the CLI correctly.

| Skill | Description |
|-------|-------------|
| [`yummy-shared`](./skills/yummy-shared/SKILL.md) | Credential checks, output contract, and shared safety rules — loaded automatically by all other skills |
| [`yummy-gen-image`](./skills/yummy-gen-image/SKILL.md) | Text-to-image generation, single-image editing, and multi-image reference editing via Gemini |

Skills are located in [`./skills/`](./skills/).

### Installation

```bash
npx skills add yummysource/yummycli -y -g
```

Load `yummy-shared` before any other yummycli skill.

---

## Command Reference

```
yummycli
├── version                              Show the yummycli version
│
├── auth
│   ├── init    --provider  --api-key    Save API key for a provider
│   ├── list                             List all providers and credential status
│   ├── status  --provider               Show credential status for a provider
│   └── remove  --provider               Delete stored credentials
│
├── gemini
│   ├── init  --api-key                  Initialize Gemini credentials
│   └── nanobanana                       Generate / edit images with Gemini
│         --prompt        (required)
│         --output
│         --model
│         --aspect-ratio
│         --image-size
│         --input-image   (repeatable)
│
└── image
    └── generate                         Provider-agnostic image generation
          --provider      (required)
          --prompt        (required)
          --output
          --model
          --aspect-ratio
          --image-size
          --input-image   (repeatable)
```

---

## Contributing

Contributions are welcome. If you find a bug or have a feature request, open an [Issue](https://github.com/yummysource/yummycli/issues) or [Pull Request](https://github.com/yummysource/yummycli/pulls).

For significant changes, open an Issue first to discuss the approach.

## License

[MIT](./LICENSE)
