# yummycli

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23-blue.svg)](https://go.dev/)
[![npm version](https://img.shields.io/npm/v/@yummysource/yummycli.svg)](https://www.npmjs.com/package/@yummysource/yummycli)

[繁體中文](./README.zh-TW.md) | [简体中文](./README.zh.md) | [English](./README.md)

<img src="./assets/logo.png" alt="yummycli logo" width="120" />

An AI-friendly CLI for multimodal model providers — built for humans and AI Agents.

Supports image generation and editing via **Gemini** and **OpenAI**, plus video generation and text-to-speech synthesis via Gemini.

[Install](#installation) · [Auth](#authentication) · [Image Generation](#image-generation) · [Video Generation](#video-generation) · [Voice Generation](#voice-generation) · [Agent Skills](#agent-skills) · [Commands](#command-reference)

---

## Why yummycli?

- **Agent-Native Design** — Structured [Skills](./skills/) out of the box, designed so AI Agents can call image, video, and audio APIs with zero extra setup
- **Capability-First Architecture** — `image generate`, `video generate`, and `audio speak` are the stable automation contracts; `gemini nanobanana`, `gemini veo`, and `gemini speak` are human-friendly shortcuts on top
- **Structured JSON Output** — Every command writes JSON to stdout, making it trivial to pipe into agents, scripts, or other tools
- **Secure Credential Storage** — API keys stored in the OS-native keychain (macOS Keychain, Linux Secret Service), never in plain text
- **Multi-Provider with Automatic Fallback** — Configure two providers; if the primary fails, the other is tried automatically

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

yummycli stores API keys per provider in the OS keychain. Configure each provider once.

### Provider Coverage

| Capability | Gemini | OpenAI |
|---|---|---|
| Image generation & editing | ✅ | ✅ |
| Video generation | ✅ | — |
| Speech synthesis | ✅ | — |

### Quick Setup

```bash
# Configure primary provider (set as default)
yummycli init --provider gemini --api-key "AIza..." --default

# Optionally add a second provider as fallback
yummycli init --provider openai --api-key "sk-..."
```

If both providers are configured, the non-default acts as an automatic fallback when the primary fails.

### Check configuration

```bash
yummycli auth list
```

```json
[
  {"provider":"gemini","configured":true,"default":true,"apiKeyPreview":"AIza******g7Aw"},
  {"provider":"openai","configured":true,"default":false,"apiKeyPreview":"sk-p******mEAA"}
]
```

`"default": true` identifies the provider used when `--provider` is omitted.

### Switch default provider

```bash
# Key already stored — no need to re-enter it
yummycli init --provider openai --default
```

### Other auth commands

| Command | Description |
|---------|-------------|
| `yummycli init --provider <name> --api-key <key> [--default]` | Save API key, optionally set as default |
| `yummycli auth list` | List all providers, credential status, and default |
| `yummycli auth status --provider <name>` | Show credential status for a specific provider |
| `yummycli auth remove --provider <name>` | Delete stored credentials |

### Gemini shortcut

`gemini init` is an alias for `init --provider gemini`:

```bash
yummycli gemini init --api-key "AIza..." --default
```

---

## Image Generation

Image generation supports both **Gemini** and **OpenAI**. The provider is resolved from your config — no need to specify it every time.

| Provider | Default model | Default size |
|---|---|---|
| Gemini | `gemini-3.1-flash-image` | `16:9`, `1K` |
| OpenAI | `gpt-image-2` | `1536x864` (16:9) |

Two equivalent entry points are available:

| Entry point | Intended for |
|-------------|-------------|
| `yummycli image generate` | All providers — uses configured default |
| `yummycli gemini nanobanana` | Gemini shortcut with defaults pre-applied |

### Quick start

```bash
# Step 1: configure providers (one-time)
yummycli init --provider gemini --api-key "AIza..." --default
yummycli init --provider openai --api-key "sk-..."   # optional fallback

# Step 2: generate an image (uses default provider)
yummycli image generate --prompt "A ripe banana on a white plate, studio lighting"
```

### Flags

**Common flags (both providers):**

| Flag | Description | Default |
|------|-------------|---------|
| `--prompt` | Image generation or editing prompt (**required**) | — |
| `--provider` | Override provider (`gemini` or `openai`) | from config |
| `--model` | Model name | provider default |
| `--output` | Output file path | auto-generated |
| `--input-image` | Input image for editing (repeatable) | — |

**Gemini-specific flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--aspect-ratio` | Named ratio: `16:9`, `9:16`, `21:9`, `1:1`, etc. | `16:9` |
| `--image-size` | Resolution: `512`, `0.5K`, `1K`, `2K`, `4K` | `1K` |

**OpenAI-specific flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--image-size` | Dimensions in `WxH` format (multiples of 16, ratio 1:3–3:1) | `1536x864` |
| `--quality` | Rendering quality: `low`, `medium`, `high`, `auto` | API auto |
| `--output-format` | Output format: `png`, `jpeg`, `webp` | `png` |

### Text-to-image generation

```bash
# Uses configured default provider
yummycli image generate --prompt "A cyberpunk cityscape at night"

# Explicit Gemini with 4K resolution
yummycli image generate --provider gemini \
  --prompt "A minimalist logo, flat design" \
  --image-size 4K --output logo.png

# Explicit OpenAI with JPEG output
yummycli image generate --provider openai \
  --prompt "A rabbit on a wooden table" \
  --output-format jpeg
```

### Image editing

Pass one or more reference images with `--input-image`. Both providers support editing.

```bash
# Single-image edit (Gemini)
yummycli image generate --provider gemini \
  --prompt "Turn this into a watercolor illustration" \
  --input-image ./photo.png

# Single-image edit (OpenAI)
yummycli image generate --provider openai \
  --prompt "Add a rainbow in the background" \
  --input-image ./photo.png

# Multi-image compositing (Gemini)
yummycli image generate --provider gemini \
  --prompt "Blend these references into one polished poster" \
  --input-image ./subject.png \
  --input-image ./background.jpg
```

Supported input formats: `.png`, `.jpg` / `.jpeg`, `.webp`.

### Aspect ratio

Both providers support 16:9 by default. To request a different ratio:

| Ratio | Gemini flag | OpenAI flag |
|---|---|---|
| 21:9 ultrawide | `--aspect-ratio 21:9` | `--image-size 2016x864` |
| 16:9 widescreen | `--aspect-ratio 16:9` | `--image-size 1536x864` |
| 1:1 square | `--aspect-ratio 1:1` | `--image-size 1024x1024` |
| 9:16 portrait | `--aspect-ratio 9:16` | `--image-size 864x1536` |

### Model selection

**Gemini:**
```bash
# Flash (default) — faster, supports more aspect ratios and sizes
yummycli image generate --provider gemini --prompt "..." --model gemini-3.1-flash-image

# Pro — higher quality
yummycli image generate --provider gemini --prompt "..." --model gemini-3-pro-image-preview
```

**OpenAI:**
```bash
# gpt-image-2 (default)
yummycli image generate --provider openai --prompt "..."

# gpt-5.5
yummycli image generate --provider openai --prompt "..." --model gpt-5.5
```

Unknown OpenAI model names trigger a warning and fall back to `gpt-image-2` automatically.

### Gemini model compatibility

**Aspect ratio:**

| Model | Supported values |
|-------|-----------------|
| `gemini-3.1-flash-image` | `1:1` `1:4` `1:8` `2:3` `3:2` `3:4` `4:1` `4:3` `4:5` `5:4` `8:1` `9:16` `16:9` `21:9` |
| `gemini-3-pro-image-preview` | `1:1` `2:3` `3:2` `3:4` `4:3` `4:5` `5:4` `9:16` `16:9` `21:9` |

`1:4`, `1:8`, `4:1`, `8:1` are Flash-only.

**Image size:**

| Model | Supported values |
|-------|-----------------|
| `gemini-3.1-flash-image` | `512` `0.5K` `1K` `2K` `4K` |
| `gemini-3-pro-image-preview` | `1K` `2K` `4K` |

`512` and `0.5K` are Flash-only.

### JSON output

```json
{
  "provider": "gemini",
  "output": "gemini_20260410123456_789.png",
  "model": "gemini-3.1-flash-image",
  "inputImageCount": 0
}
```

---

## Video Generation

Video generation uses Google Veo (Gemini only).

| Entry point | Intended for |
|-------------|-------------|
| `gemini veo` | Human use — Gemini Veo defaults pre-applied |
| `video generate --provider gemini` | Automation / scripting |

### Quick start

```bash
yummycli gemini veo --prompt "A golden retriever puppy chasing a red ball in a sunny park"
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
| `--input-image` | Input image (repeatable, up to 3) | — |

### Generation modes

| `--input-image` count | Mode |
|-----------------------|------|
| 0 | Text-to-video |
| 1 | Image-to-video — image used as starting frame |
| 2–3 | Reference-guided — images used as ASSET reference inputs |

```bash
# Text-to-video
yummycli gemini veo --prompt "Timelapse of clouds over mountain peaks"

# Image-to-video
yummycli gemini veo --prompt "The dog starts running" --input-image ./dog.jpg

# Reference-guided
yummycli gemini veo --prompt "Combine character with background" \
  --input-image ./character.png --input-image ./background.jpg
```

### Model compatibility

**Duration** (discrete values only):

| Model | Valid durations (seconds) |
|-------|--------------------------|
| `veo-2.0-generate-001` | 5, 6, 7, 8 |
| `veo-3.0-*` | 4, 6, 8 |
| `veo-3.1-*` | 4, 6, 8 |

**Resolution:**

| Model | Supported |
|-------|-----------------------|
| `veo-2.0-generate-001` | `720p` |
| `veo-3.0-*` | `720p`, `1080p` |
| `veo-3.1-*` | `720p`, `1080p`, `4k` |

`1080p` and `4k` require `--duration 8`. `4k` requires a veo-3.1 model.

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

Text-to-speech synthesis uses Google Gemini TTS (Gemini only).

| Entry point | Intended for |
|-------------|-------------|
| `gemini speak` | Human use — Gemini TTS defaults pre-applied |
| `audio speak --provider gemini` | Automation / scripting |

### Quick start

```bash
yummycli gemini speak --text "A golden retriever is the best dog in the world."
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
yummycli gemini speak --text "Hello, world!"

yummycli gemini speak \
  --text "你好，欢迎使用语音合成服务。" \
  --voice Aoede --language zh-CN --output greeting.wav
```

### Multi-speaker dialogue

```bash
yummycli gemini speak \
  --text "[Alice]: Hello! [Bob]: Hi Alice, great to meet you!" \
  --speaker Alice:Aoede --speaker Bob:Kore \
  --output dialogue.wav
```

### List available voices

```bash
yummycli gemini voices
```

Returns a JSON array of all 30 prebuilt voices with style descriptions.

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

---

## Agent Skills

yummycli ships with Skills — structured instruction files that teach AI Agents how to use the CLI correctly.

| Skill | Description |
|-------|-------------|
| [`yummy-shared`](./skills/yummy-shared/SKILL.md) | First-time setup, provider config, output contract, safety rules |
| [`yummy-gen-image`](./skills/yummy-gen-image/SKILL.md) | Image generation and editing — Gemini and OpenAI |
| [`yummy-gen-video`](./skills/yummy-gen-video/SKILL.md) | Text-to-video, image-to-video, reference-guided video — Gemini Veo |
| [`yummy-gen-voice`](./skills/yummy-gen-voice/SKILL.md) | Single-speaker TTS, multi-speaker dialogue, voice listing — Gemini TTS |

### Installation

```bash
npx skills add yummysource/yummycli -y -g
```

---

## Command Reference

```
yummycli
├── version                                     Show the yummycli version
│
├── init  --provider  [--api-key]  [--default]  Configure a provider API key
│         --provider   gemini or openai (required)
│         --api-key    API key (required if not already stored)
│         --default    Set this provider as the default
│
├── auth
│   ├── list                                    List providers, status, and default
│   ├── status  --provider                      Show credential status
│   └── remove  --provider                      Delete stored credentials
│
├── gemini
│   ├── init  --api-key  [--default]            Initialize Gemini credentials
│   ├── nanobanana                              Generate / edit images with Gemini
│   │     --prompt         (required)
│   │     --output
│   │     --model          default: gemini-3.1-flash-image
│   │     --aspect-ratio   default: 16:9
│   │     --image-size     default: 1K
│   │     --input-image    (repeatable)
│   ├── veo                                     Generate videos with Gemini Veo
│   │     --prompt         (required)
│   │     --output
│   │     --model          default: veo-3.1-fast-generate-preview
│   │     --aspect-ratio   default: 16:9
│   │     --duration       default: 8
│   │     --resolution     default: 1080p
│   │     --input-image    (repeatable, up to 3)
│   ├── speak                                   Synthesise speech with Gemini TTS
│   │     --text           (required)
│   │     --output
│   │     --model          default: gemini-3.1-flash-tts-preview
│   │     --voice          default: Aoede
│   │     --language
│   │     --speaker        (repeatable, up to 2; mutually exclusive with --voice)
│   └── voices                                  List available Gemini TTS voices
│
├── image
│   └── generate                                Generate or edit an image
│         --prompt         (required)
│         --provider       uses config default if omitted
│         --output
│         --model
│         --aspect-ratio   Gemini only
│         --image-size     Gemini: 1K/2K/4K — OpenAI: WxH e.g. 1536x864
│         --quality        OpenAI only: low/medium/high/auto
│         --output-format  OpenAI only: png/jpeg/webp
│         --input-image    (repeatable)
│
├── video
│   └── generate                                Generate a video (Gemini only)
│         --provider       (required)
│         --prompt         (required)
│         --output
│         --model
│         --aspect-ratio
│         --duration
│         --resolution
│         --input-image    (repeatable, up to 3)
│
└── audio
    ├── speak                                   Synthesise speech (Gemini only)
    │     --provider       (required)
    │     --text           (required)
    │     --output
    │     --model
    │     --voice
    │     --language
    │     --speaker        (repeatable, up to 2; mutually exclusive with --voice)
    └── voices                                  List available voices
          --provider       (required)
```

---

## Contributing

Contributions are welcome. If you find a bug or have a feature request, open an [Issue](https://github.com/yummysource/yummycli/issues) or [Pull Request](https://github.com/yummysource/yummycli/pulls).

For significant changes, open an Issue first to discuss the approach.

## License

[MIT](./LICENSE)
