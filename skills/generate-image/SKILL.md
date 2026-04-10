---
name: generate-image
version: 1.0.0
description: "Use when the user wants to generate or edit raster images with Gemini through yummycli, including prompt-only generation, single-image editing, and multi-image reference editing."
metadata:
  requires:
    bins: ["yummycli"]
---

# Generate Image

Create or edit images with `yummycli gemini nanobanana`.

> **Prerequisite:** Apply the `yummy-shared` skill first.

This skill uses one command for all Gemini image flows:

- Prompt-only generation
- Single-image editing
- Multi-image reference editing

## Command Contract

Always use:

```bash
yummycli gemini nanobanana --prompt "<prompt>"
```

Add one or more reference images when needed:

```bash
yummycli gemini nanobanana \
  --prompt "<prompt>" \
  --input-image ./source-a.png \
  --input-image ./source-b.png
```

Optional output controls:

```bash
--output <file>
--model <model>
--aspect-ratio <ratio>
--image-size <size>
```

If `--output` is omitted, `yummycli` generates a default filename in the current working directory.

Default values when omitted: `--aspect-ratio 16:9`, `--image-size 1K`, `--model gemini-3.1-flash-image-preview`.

## Execution Rules

- Use prompt-only generation when no reference images are provided.
- Use one `--input-image` flag per local image file.
- Preserve the user-provided order of reference images.
- Prefer the default model unless the user asks for a specific model.
- Use `--aspect-ratio` and `--image-size` only when the user specifies them or when a concrete output format is clearly required.

## Model Selection

Use the following model mapping rules:

- If the user explicitly says `gemini pro`, `pro model`, or simply `pro` in the model-selection context, use:

```bash
--model gemini-3-pro-image-preview
```

- If the user explicitly says `gemini flash`, `flash model`, or simply `flash` in the model-selection context, use:

```bash
--model gemini-3.1-flash-image-preview
```

- If the user does not explicitly request a model, omit `--model` and let `yummycli` use its default model.

Do not switch models implicitly from vague quality words alone. Only apply the `pro` or `flash` mapping when the user's wording clearly refers to model choice.

## Model and Aspect-Ratio Compatibility

Some aspect ratios are only valid for the flash model. When the pro model is selected, restrict aspect-ratio choices to the intersection supported by both models:

`1:1`, `2:3`, `3:2`, `3:4`, `4:3`, `4:5`, `5:4`, `9:16`, `16:9`, `21:9`

The following ratios are flash-only and must not be used with the pro model:

`1:4`, `1:8`, `4:1`, `8:1`

## Intent to Parameters

Translate clear user intent into CLI flags when the mapping is obvious.

Aspect-ratio guidance:

- Use `--aspect-ratio 9:16` for requests such as phone wallpaper, vertical poster, story format, or other clearly vertical mobile outputs.
- Use `--aspect-ratio 16:9` for requests such as desktop wallpaper, presentation cover, horizontal banner, or other clearly widescreen outputs.
- Use `--aspect-ratio 1:1` for square social images, avatars, or other explicitly square outputs.
- If the user already provides a specific ratio, pass it through directly.
- If the output shape is unclear, omit `--aspect-ratio` and let `yummycli` use its default.

Image-size guidance:

- Use `--image-size 4K` only when the user explicitly asks for 4K or a clearly print-grade / high-resolution deliverable.
- Use `--image-size 2K` when the user explicitly asks for 2K or a medium-resolution deliverable.
- Use `--image-size 1K` only when the user explicitly asks for 1K or a lightweight preview-sized result.
- If the user does not explicitly request output size, omit `--image-size`.

Do not guess `--image-size` from general quality adjectives alone.

## Examples

Prompt-only generation:

```bash
yummycli gemini nanobanana \
  --prompt "A single ripe banana on a white plate, studio lighting, realistic photo"
```

Single-image edit:

```bash
yummycli gemini nanobanana \
  --prompt "Turn this into a watercolor illustration" \
  --input-image ./source.png
```

Multi-image reference edit:

```bash
yummycli gemini nanobanana \
  --prompt "Blend these references into one polished poster illustration" \
  --input-image ./subject.png \
  --input-image ./background.jpg
```
