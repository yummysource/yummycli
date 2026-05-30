# yummycli

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23-blue.svg)](https://go.dev/)
[![npm version](https://img.shields.io/npm/v/@yummysource/yummycli.svg)](https://www.npmjs.com/package/@yummysource/yummycli)

面向多模态模型供应商的 AI 友好命令行工具 —— 专为人类用户和 AI Agent 设计。

支持通过 **Gemini** 和 **OpenAI** 进行图像生成与编辑，以及通过 Gemini 进行视频生成和语音合成。

[繁體中文](./README.zh-TW.md) | [简体中文](./README.zh.md) | [English](./README.md)

<img src="./assets/logo.png" alt="yummycli logo" width="120" />

[安装](#安装) · [认证](#认证) · [图像生成](#图像生成) · [视频生成](#视频生成) · [语音合成](#语音合成) · [Agent Skills](#agent-skills) · [命令参考](#命令参考)

---

## 为什么选择 yummycli？

- **Agent 原生设计** —— 开箱即用的结构化 [Skills](./skills/)，AI Agent 无需额外配置即可调用图像、视频与音频 API
- **能力优先架构** —— `image generate`、`video generate` 和 `audio speak` 是稳定的自动化接口；`gemini nanobanana`、`gemini veo` 和 `gemini speak` 是其上层的人性化快捷方式
- **结构化 JSON 输出** —— 每条命令将结果写入 stdout，方便 Agent、脚本和其他工具直接消费
- **安全的凭证存储** —— API Key 存储在操作系统原生密钥链（macOS Keychain、Linux Secret Service），从不以明文保存
- **多供应商自动降级** —— 配置两个供应商后，主供应商失败时自动切换到另一个

---

## 安装

### 环境要求

- Node.js 16+ 及 `npm`
- Go 1.23+ 和 `make`（仅从源码构建时需要）

### 通过 npm 安装（推荐）

```bash
# 安装 CLI
npm install -g @yummysource/yummycli

# 安装 Agent Skills（AI Agent 使用必须）
npx skills add yummysource/yummycli -y -g
```

验证安装：

```bash
yummycli version
```

### 从源码构建

```bash
git clone https://github.com/yummysource/yummycli.git
cd yummycli
make install

# 安装 Agent Skills（AI Agent 使用必须）
npx skills add yummysource/yummycli -y -g
```

---

## 认证

yummycli 将每个供应商的 API Key 存储在操作系统密钥链中，每个供应商只需配置一次。

### 供应商覆盖范围

| 能力 | Gemini | OpenAI |
|---|---|---|
| 图像生成与编辑 | ✅ | ✅ |
| 视频生成 | ✅ | — |
| 语音合成 | ✅ | — |

### 快速配置

```bash
# 配置主供应商并设为默认
yummycli init --provider gemini --api-key "AIza..." --default

# 可选：添加第二个供应商作为降级备选
yummycli init --provider openai --api-key "sk-..."
```

配置两个供应商后，主供应商失败时自动使用另一个。

### 查看配置状态

```bash
yummycli auth list
```

```json
[
  {"provider":"gemini","configured":true,"default":true,"apiKeyPreview":"AIza******g7Aw"},
  {"provider":"openai","configured":true,"default":false,"apiKeyPreview":"sk-p******mEAA"}
]
```

`"default": true` 表示省略 `--provider` 时使用的供应商。

### 切换默认供应商

```bash
# Key 已存储时无需重新输入
yummycli init --provider openai --default
```

### 其他认证命令

| 命令 | 说明 |
|------|------|
| `yummycli init --provider <name> --api-key <key> [--default]` | 保存 API Key，可选设为默认 |
| `yummycli auth list` | 列出所有供应商、凭证状态及默认值 |
| `yummycli auth status --provider <name>` | 查看指定供应商的凭证状态 |
| `yummycli auth remove --provider <name>` | 删除指定供应商的凭证 |

### Gemini 快捷方式

`gemini init` 是 `init --provider gemini` 的等价快捷命令：

```bash
yummycli gemini init --api-key "AIza..." --default
```

---

## 图像生成

图像生成支持 **Gemini** 和 **OpenAI**。供应商从配置中自动解析，无需每次指定。

| 供应商 | 默认模型 | 默认尺寸 |
|---|---|---|
| Gemini | `gemini-3.1-flash-image` | `16:9`，`1K` |
| OpenAI | `gpt-image-2` | `1536x864`（16:9） |

两个等价入口：

| 入口 | 适用场景 |
|------|----------|
| `yummycli image generate` | 所有供应商 —— 使用配置的默认供应商 |
| `yummycli gemini nanobanana` | Gemini 快捷方式，预设默认参数 |

### 快速开始

```bash
# 第一步：配置供应商（一次性）
yummycli init --provider gemini --api-key "AIza..." --default
yummycli init --provider openai --api-key "sk-..."   # 可选降级备选

# 第二步：生成图像（使用默认供应商）
yummycli image generate --prompt "白色盘子上放着一根成熟的香蕉，工作室打光"
```

### 参数说明

**通用参数（两个供应商均支持）：**

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--prompt` | 图像生成或编辑提示词（**必填**） | — |
| `--provider` | 覆盖供应商（`gemini` 或 `openai`） | 从配置读取 |
| `--model` | 模型名称 | 供应商默认值 |
| `--output` | 输出文件路径 | 自动生成 |
| `--input-image` | 编辑用输入图像（可重复） | — |

**Gemini 专属参数：**

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--aspect-ratio` | 命名比例：`16:9`、`9:16`、`21:9`、`1:1` 等 | `16:9` |
| `--image-size` | 分辨率：`512`、`0.5K`、`1K`、`2K`、`4K` | `1K` |

**OpenAI 专属参数：**

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--image-size` | 尺寸（WxH 格式，两边均为 16 的倍数，比例 1:3 到 3:1） | `1536x864` |
| `--quality` | 渲染质量：`low`、`medium`、`high`、`auto` | API 自动 |
| `--output-format` | 输出格式：`png`、`jpeg`、`webp` | `png` |

### 文本生成图像

```bash
# 使用配置的默认供应商
yummycli image generate --prompt "赛博朋克夜晚都市，霓虹灯倒映在街道上"

# 指定 Gemini，4K 分辨率
yummycli image generate --provider gemini \
  --prompt "极简主义 logo，扁平设计" \
  --image-size 4K --output logo.png

# 指定 OpenAI，JPEG 输出
yummycli image generate --provider openai \
  --prompt "木桌上的一只白兔" \
  --output-format jpeg
```

### 图像编辑

通过 `--input-image` 传入参考图像，两个供应商均支持编辑。

```bash
# 单图编辑（Gemini）
yummycli image generate --provider gemini \
  --prompt "将这张图转换为水彩插画风格" \
  --input-image ./photo.png

# 单图编辑（OpenAI）
yummycli image generate --provider openai \
  --prompt "在背景中添加一道彩虹" \
  --input-image ./photo.png

# 多图合成（Gemini）
yummycli image generate --provider gemini \
  --prompt "将这两张参考图融合为一张海报" \
  --input-image ./subject.png \
  --input-image ./background.jpg
```

支持的输入格式：`.png`、`.jpg` / `.jpeg`、`.webp`。

### 宽高比

两个供应商默认均为 16:9。指定其他比例：

| 比例 | Gemini 参数 | OpenAI 参数 |
|---|---|---|
| 21:9 超宽 | `--aspect-ratio 21:9` | `--image-size 2016x864` |
| 16:9 宽屏 | `--aspect-ratio 16:9` | `--image-size 1536x864` |
| 1:1 正方形 | `--aspect-ratio 1:1` | `--image-size 1024x1024` |
| 9:16 竖版 | `--aspect-ratio 9:16` | `--image-size 864x1536` |

### 模型选择

**Gemini：**
```bash
# Flash（默认）—— 速度更快，支持更多比例和尺寸
yummycli image generate --provider gemini --prompt "..." --model gemini-3.1-flash-image

# Pro —— 质量更高
yummycli image generate --provider gemini --prompt "..." --model gemini-3-pro-image-preview
```

**OpenAI：**
```bash
# gpt-image-2（默认）
yummycli image generate --provider openai --prompt "..."

# gpt-5.5
yummycli image generate --provider openai --prompt "..." --model gpt-5.5
```

未知 OpenAI 模型名称会触发警告并自动回退到 `gpt-image-2`。

### Gemini 模型兼容性

**宽高比：**

| 模型 | 支持的值 |
|------|---------|
| `gemini-3.1-flash-image` | `1:1` `1:4` `1:8` `2:3` `3:2` `3:4` `4:1` `4:3` `4:5` `5:4` `8:1` `9:16` `16:9` `21:9` |
| `gemini-3-pro-image-preview` | `1:1` `2:3` `3:2` `3:4` `4:3` `4:5` `5:4` `9:16` `16:9` `21:9` |

`1:4`、`1:8`、`4:1`、`8:1` 仅 Flash 模型支持。

**图像尺寸：**

| 模型 | 支持的值 |
|------|---------|
| `gemini-3.1-flash-image` | `512` `0.5K` `1K` `2K` `4K` |
| `gemini-3-pro-image-preview` | `1K` `2K` `4K` |

`512` 和 `0.5K` 仅 Flash 模型支持。

### JSON 输出

```json
{
  "provider": "gemini",
  "output": "gemini_20260410123456_789.png",
  "model": "gemini-3.1-flash-image",
  "inputImageCount": 0
}
```

---

## 视频生成

视频生成通过 Google Veo 实现，仅支持 Gemini。

| 入口 | 适用场景 |
|------|----------|
| `gemini veo` | 人工使用 —— 已预设 Gemini Veo 默认参数 |
| `video generate --provider gemini` | 自动化/脚本 |

### 快速开始

```bash
yummycli gemini veo --prompt "阳光明媚的公园里，金毛犬追逐红球"
```

### 参数说明

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--prompt` | 视频生成提示词（**必填**） | — |
| `--output` | 输出文件路径（须以 `.mp4` 结尾） | 自动生成 |
| `--model` | Veo 模型 | `veo-3.1-fast-generate-preview` |
| `--aspect-ratio` | 视频宽高比 | `16:9` |
| `--duration` | 时长（秒） | `8` |
| `--resolution` | 视频分辨率 | `1080p` |
| `--input-image` | 输入图像（可重复，最多 3 张） | — |

### 生成模式

| `--input-image` 数量 | 模式 |
|----------------------|------|
| 0 | 文本生成视频 |
| 1 | 图像生成视频 —— 图像作为起始帧 |
| 2–3 | 参考图引导 —— 图像作为 ASSET 参考输入 |

### 模型兼容性

**时长**（仅接受离散值）：

| 模型 | 有效时长（秒） |
|------|--------------|
| `veo-2.0-generate-001` | 5、6、7、8 |
| `veo-3.0-*` | 4、6、8 |
| `veo-3.1-*` | 4、6、8 |

**分辨率：**

| 模型 | 支持的分辨率 |
|------|------------|
| `veo-2.0-generate-001` | `720p` |
| `veo-3.0-*` | `720p`、`1080p` |
| `veo-3.1-*` | `720p`、`1080p`、`4k` |

`1080p` 和 `4k` 需要 `--duration 8`；`4k` 需要 veo-3.1 系列模型。

### JSON 输出

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

## 语音合成

文本转语音通过 Google Gemini TTS 实现，仅支持 Gemini。

| 入口 | 适用场景 |
|------|----------|
| `gemini speak` | 人工使用 —— 已预设 Gemini TTS 默认参数 |
| `audio speak --provider gemini` | 自动化/脚本 |

### 快速开始

```bash
yummycli gemini speak --text "你好，这是一段语音合成示例。"
```

### 参数说明

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--text` | 要合成的文本（**必填**） | — |
| `--output` | 输出文件路径（须以 `.wav` 结尾） | 自动生成 |
| `--model` | TTS 模型 | `gemini-3.1-flash-tts-preview` |
| `--voice` | 预设声音名称 | `Aoede` |
| `--language` | BCP-47 语言代码（省略时自动检测） | — |
| `--speaker` | 多说话人映射 `名字:声音`（可重复，最多 2 次） | — |

`--voice` 与 `--speaker` 互斥。

### 单说话人合成

```bash
yummycli gemini speak --text "今天天气真好！"

yummycli gemini speak \
  --text "欢迎来到 AI 的世界。" \
  --voice Kore --output greeting.wav
```

### 多说话人对话合成

```bash
yummycli gemini speak \
  --text "[小明]: 你好！今天天气真好。 [小红]: 是啊，我们去公园走走吧！" \
  --speaker 小明:Aoede \
  --speaker 小红:Kore \
  --output dialogue.wav
```

### 列出可用声音

```bash
yummycli gemini voices
```

返回所有 30 个预设声音及其风格描述的 JSON 列表。

### JSON 输出

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

yummycli 内置 Skills —— 结构化的指令文件，帮助 AI Agent 正确使用 CLI。

| Skill | 说明 |
|-------|------|
| [`yummy-shared`](./skills/yummy-shared/SKILL.md) | 初次配置、供应商设置、输出格式约定和安全规则 |
| [`yummy-gen-image`](./skills/yummy-gen-image/SKILL.md) | 图像生成与编辑 —— Gemini 和 OpenAI |
| [`yummy-gen-video`](./skills/yummy-gen-video/SKILL.md) | 视频生成 —— Gemini Veo |
| [`yummy-gen-voice`](./skills/yummy-gen-voice/SKILL.md) | 语音合成 —— Gemini TTS |

### 安装

```bash
npx skills add yummysource/yummycli -y -g
```

---

## 命令参考

```
yummycli
├── version                                     显示 yummycli 版本
│
├── init  --provider  [--api-key]  [--default]  配置供应商 API Key
│         --provider   gemini 或 openai（必填）
│         --api-key    API Key（未存储时必填）
│         --default    设置为默认供应商
│
├── auth
│   ├── list                                    列出所有供应商、凭证状态及默认值
│   ├── status  --provider                      查看指定供应商的凭证状态
│   └── remove  --provider                      删除指定供应商的凭证
│
├── gemini
│   ├── init  --api-key  [--default]            初始化 Gemini 凭证
│   ├── nanobanana                              使用 Gemini 生成 / 编辑图像
│   │     --prompt         （必填）
│   │     --output
│   │     --model          默认：gemini-3.1-flash-image
│   │     --aspect-ratio   默认：16:9
│   │     --image-size     默认：1K
│   │     --input-image    （可重复）
│   ├── veo                                     使用 Gemini Veo 生成视频
│   │     --prompt         （必填）
│   │     --output
│   │     --model          默认：veo-3.1-fast-generate-preview
│   │     --aspect-ratio   默认：16:9
│   │     --duration       默认：8
│   │     --resolution     默认：1080p
│   │     --input-image    （可重复，最多 3 张）
│   ├── speak                                   使用 Gemini TTS 合成语音
│   │     --text           （必填）
│   │     --output
│   │     --model          默认：gemini-3.1-flash-tts-preview
│   │     --voice          默认：Aoede
│   │     --language
│   │     --speaker        （可重复，最多 2 次；与 --voice 互斥）
│   └── voices                                  列出 Gemini TTS 可用声音
│
├── image
│   └── generate                                生成或编辑图像
│         --prompt         （必填）
│         --provider       省略时从配置读取
│         --output
│         --model
│         --aspect-ratio   仅 Gemini
│         --image-size     Gemini: 1K/2K/4K — OpenAI: WxH 如 1536x864
│         --quality        仅 OpenAI：low/medium/high/auto
│         --output-format  仅 OpenAI：png/jpeg/webp
│         --input-image    （可重复）
│
├── video
│   └── generate                                生成视频（仅 Gemini）
│         --provider       （必填）
│         --prompt         （必填）
│         --output
│         --model
│         --aspect-ratio
│         --duration
│         --resolution
│         --input-image    （可重复，最多 3 张）
│
└── audio
    ├── speak                                   合成语音（仅 Gemini）
    │     --provider       （必填）
    │     --text           （必填）
    │     --output
    │     --model
    │     --voice
    │     --language
    │     --speaker        （可重复，最多 2 次；与 --voice 互斥）
    └── voices                                  列出指定供应商的可用声音
          --provider       （必填）
```

---

## 贡献

欢迎社区贡献。如发现 Bug 或有功能建议，请提交 [Issue](https://github.com/yummysource/yummycli/issues) 或 [Pull Request](https://github.com/yummysource/yummycli/pulls)。

重大改动建议先通过 Issue 与我们讨论。

## 许可证

[MIT](./LICENSE)
