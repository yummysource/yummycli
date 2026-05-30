---
name: yummy-shared
version: 2.0.0
description: "Use when operating yummycli — covers first-time setup, provider configuration, JSON output parsing, and shared CLI safety rules."
always: true
metadata:
  requires:
    bins: ["yummycli"]
  openclaw:
    requires:
      bins: ["yummycli"]
  hermes:
    tags: [yummycli, shared, authentication]
install:
  - kind: node
    package: "@yummysource/yummycli"
    bins: ["yummycli"]
---

# yummycli Shared Rules

## Provider Coverage

| Capability | Gemini | OpenAI |
|---|---|---|
| Image generation & editing | ✅ | ✅ |
| Video generation | ✅ | — |
| Speech synthesis (TTS) | ✅ | — |

Video and speech require Gemini. Configure Gemini first if you need these capabilities.

## First-Time Setup

Check which providers are configured:

```bash
yummycli auth list
```

If no provider is configured, initialize one:

```bash
# Gemini
yummycli init --provider gemini --api-key "<key>" --default

# OpenAI
yummycli init --provider openai --api-key "<key>" --default
```

Add a second provider as fallback (omit --default to keep existing default):

```bash
yummycli init --provider openai --api-key "<key>"
```

## Output Contract

All `yummycli` generation commands return JSON on stdout. Read the response and use the `output` field as the generated file path.

```json
{
  "provider": "gemini",
  "output": "gemini_20260410123456_789.png",
  "model": "gemini-3.1-flash-image",
  "inputImageCount": 2
}
```

## Safety Rules

- Only use local image files explicitly provided by the user.
- Preserve the order of repeated `--input-image` flags.
- Do not overwrite a user-specified output path unless explicitly intended.
- If the command returns a validation error, fix the arguments before retrying.
- Report the final output path back to the user after a successful run.
