# yummycli

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23-blue.svg)](https://go.dev/)
[![npm version](https://img.shields.io/npm/v/@yummysource/yummycli.svg)](https://www.npmjs.com/package/@yummysource/yummycli)

[繁體中文](./README.zh-TW.md) | [简体中文](./README.zh.md) | [English](./README.md)

<img src="./assets/logo.png" alt="yummycli logo" width="120" />

An AI-friendly CLI for multimodal model providers — built for humans and AI Agents.

Supports image generation and editing, video generation, and text-to-speech synthesis via [Gemini](https://deepmind.google/technologies/gemini/), with more providers (Claude, OpenAI, Qwen) planned.

[Install](#installation) · [Auth](#authentication) · [Image Generation](#image-generation) · [Video Generation](#video-generation) · [Voice Generation](#voice-generation) · [Agent Skills](#agent-skills) · [Commands](#command-reference)

---

## Why yummycli?

- **Agent-Native Design** — Structured [Skills](./skills/) out of the box, designed so AI Agents can call image, video, and audio APIs with zero extra setup
- **Capability-First Architecture** — `image generate`, `video generate`, and `audio speak` are the stable automation contracts; `gemini nanobanana`, `gemini veo`, and `gemini speak` are human-friendly shortcuts on top
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

## Video Generation

yummycli supports video generation via Google Veo, with two equivalent entry points:

| Entry point | Intended for |
|-------------|-------------|
| `gemini veo` | Human use — Gemini Veo defaults pre-applied |
| `video generate --provider gemini` | Automation / scripting — explicit, stable contract |

### Quick start

```bash
# Step 1: configure Gemini credentials (one-time)
yummycli gemini init --api-key "AIza..."

# Step 2: generate a video from a text prompt
yummycli gemini veo --prompt "A golden retriever puppy chasing a red ball in a sunny park"
```

The generated video is saved to the current directory with an auto-generated filename:

```
veo_20260417_142301_047.mp4
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--prompt` | Video generation prompt (**required**) | — |
| `--output` | Output file path (must end in `.mp4`) | auto-generated |
| `--model` | Veo model | `veo-3.1-fast-generate-preview` |
| `--aspect-ratio` | Video aspect ratio | `16:9` |
| `--duration` | Duration in seconds | `8` |
| `--resolution` | Video resolution | `1080p` |
| `--input-image` | Input image for image-to-video (repeatable, up to 3) | — |

### Generation modes

`--input-image` can be repeated; the count determines the generation mode automatically:

| `--input-image` count | Mode |
|-----------------------|------|
| 0 | Text-to-video |
| 1 | Image-to-video — image used as the starting frame |
| 2–3 | Reference-guided — images used as ASSET reference inputs |

```bash
# Text-to-video
yummycli gemini veo --prompt "Timelapse of clouds moving over mountain peaks"

# Image-to-video (animate a still image)
yummycli gemini veo \
  --prompt "The dog starts running toward the camera" \
  --input-image ./dog.jpg

# Reference-guided (two images)
yummycli gemini veo \
  --prompt "Combine the character with this background environment" \
  --input-image ./character.png \
  --input-image ./background.jpg
```

### Model compatibility

**Duration** accepts only discrete values:

| Model | Valid durations (seconds) |
|-------|--------------------------|
| `veo-2.0-generate-001` | 5, 6, 7, 8 |
| `veo-3.0-*` | 4, 6, 8 |
| `veo-3.1-*` | 4, 6, 8 |

**Resolution** support per model:

| Model | Supported resolutions |
|-------|-----------------------|
| `veo-2.0-generate-001` | `720p` |
| `veo-3.0-*` | `720p`, `1080p` |
| `veo-3.1-*` | `720p`, `1080p`, `4k` |

Constraints: `1080p` and `4k` require `--duration 8`. `4k` requires a veo-3.1 model.

### JSON output

```json
{
  "provider": "gemini",
  "output": "veo_20260417_142301_047.mp4",
  "model": "veo-3.1-fast-generate-preview",
  "duration_seconds": 8,
  "aspect_ratio": "16:9",
  "resolution": "1080p",
  "elapsed_seconds": 73
}
```

---

## Voice Generation

yummycli supports text-to-speech synthesis via Google Gemini TTS, with two equivalent entry points:

| Entry point | Intended for |
|-------------|-------------|
| `gemini speak` | Human use — Gemini TTS defaults pre-applied |
| `audio speak --provider gemini` | Automation / scripting — explicit, stable contract |

### Quick start

```bash
# Generate speech from text (saves to auto-generated .wav file)
yummycli gemini speak --text "A golden retriever is the best dog in the world."

# Specify voice and output path
yummycli gemini speak \
  --text "Welcome to the future of AI-powered audio." \
  --voice Puck \
  --output welcome.wav
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--text` | Text to synthesise (**required**) | — |
| `--output` | Output file path (must end in `.wav`) | auto-generated |
| `--model` | TTS model | `gemini-3.1-flash-tts-preview` |
| `--voice` | Prebuilt voice name | `Aoede` |
| `--language` | BCP-47 language code (auto-detected if omitted) | — |
| `--speaker` | Multi-speaker mapping `Name:Voice` (repeatable, up to 2) | — |

`--voice` and `--speaker` are mutually exclusive.

### Single-speaker synthesis

```bash
# Default voice (Aoede, Breezy)
yummycli gemini speak --text "Hello, world!"

# Upbeat voice with explicit output path
yummycli gemini speak \
  --text "This is an exciting announcement!" \
  --voice Puck \
  --output announcement.wav

# Explicit language code
yummycli gemini speak \
  --text "你好，欢迎使用语音合成服务。" \
  --voice Aoede \
  --language zh-CN \
  --output greeting.wav
```

### Multi-speaker dialogue

Tag each speaker's lines with `[Name]:` in the text, then map each name to a voice with `--speaker`:

```bash
yummycli gemini speak \
  --text "[Alice]: Hello! Nice to meet you. [Bob]: Hi Alice, great to meet you too!" \
  --speaker Alice:Aoede \
  --speaker Bob:Kore \
  --output dialogue.wav
```

### List available voices

```bash
yummycli gemini voices
```

Returns a JSON array of all 30 prebuilt voices with their style descriptions.

### JSON output

```json
{
  "provider": "gemini",
  "output": "tts_20260420_142301_047.wav",
  "model": "gemini-3.1-flash-tts-preview",
  "voice": "Aoede",
  "elapsed_seconds": 3
}
```

For multi-speaker requests, `voice` is replaced by `speakers`:

```json
{
  "provider": "gemini",
  "output": "dialogue_20260420_143010_112.wav",
  "model": "gemini-3.1-flash-tts-preview",
  "speakers": [
    {"name": "Alice", "voice": "Aoede"},
    {"name": "Bob", "voice": "Kore"}
  ],
  "elapsed_seconds": 4
}
```

---

## Agent Skills

yummycli ships with Skills — structured instruction files that teach AI Agents how to use the CLI correctly.

| Skill | Description |
|-------|-------------|
| [`yummy-shared`](./skills/yummy-shared/SKILL.md) | Credential checks, output contract, and shared safety rules — loaded automatically by all other skills |
| [`yummy-gen-image`](./skills/yummy-gen-image/SKILL.md) | Text-to-image generation, single-image editing, and multi-image reference editing via Gemini |
| [`yummy-gen-video`](./skills/yummy-gen-video/SKILL.md) | Text-to-video, image-to-video, and reference-image-guided video generation via Gemini Veo |
| [`yummy-gen-voice`](./skills/yummy-gen-voice/SKILL.md) | Single-speaker TTS, multi-speaker dialogue synthesis, and voice listing via Gemini TTS |

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
│   ├── nanobanana                       Generate / edit images with Gemini
│   │     --prompt        (required)
│   │     --output
│   │     --model
│   │     --aspect-ratio
│   │     --image-size
│   │     --input-image   (repeatable)
│   ├── veo                              Generate videos with Gemini Veo
│   │     --prompt        (required)
│   │     --output
│   │     --model
│   │     --aspect-ratio
│   │     --duration
│   │     --resolution
│   │     --input-image   (repeatable, up to 3)
│   ├── speak                            Synthesise speech with Gemini TTS
│   │     --text          (required)
│   │     --output
│   │     --model
│   │     --voice
│   │     --language
│   │     --speaker       (repeatable, up to 2; mutually exclusive with --voice)
│   └── voices                           List available Gemini TTS voices
│
├── image
│   └── generate                         Provider-agnostic image generation
│         --provider      (required)
│         --prompt        (required)
│         --output
│         --model
│         --aspect-ratio
│         --image-size
│         --input-image   (repeatable)
│
├── video
│   └── generate                         Provider-agnostic video generation
│         --provider      (required)
│         --prompt        (required)
│         --output
│         --model
│         --aspect-ratio
│         --duration
│         --resolution
│         --input-image   (repeatable, up to 3)
│
└── audio
    ├── speak                            Provider-agnostic speech synthesis
    │     --provider      (required)
    │     --text          (required)
    │     --output
    │     --model
    │     --voice
    │     --language
    │     --speaker       (repeatable, up to 2; mutually exclusive with --voice)
    └── voices                           List available voices for a provider
          --provider      (required)
```

---

## Contributing

Contributions are welcome. If you find a bug or have a feature request, open an [Issue](https://github.com/yummysource/yummycli/issues) or [Pull Request](https://github.com/yummysource/yummycli/pulls).

For significant changes, open an Issue first to discuss the approach.

## License

[MIT](./LICENSE)
