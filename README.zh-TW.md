# yummycli

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23-blue.svg)](https://go.dev/)
[![npm version](https://img.shields.io/npm/v/@yummysource/yummycli.svg)](https://www.npmjs.com/package/@yummysource/yummycli)

面向多模態模型供應商的 AI 友好命令列工具 —— 專為人類用戶和 AI Agent 設計。

支援透過 **Gemini** 和 **OpenAI** 進行圖像生成與編輯，以及透過 Gemini 進行影片生成和語音合成。

[繁體中文](./README.zh-TW.md) | [简体中文](./README.zh.md) | [English](./README.md)

<img src="./assets/logo.png" alt="yummycli logo" width="120" />

[安裝](#安裝) · [認證](#認證) · [圖像生成](#圖像生成) · [影片生成](#影片生成) · [語音合成](#語音合成) · [Agent Skills](#agent-skills) · [命令參考](#命令參考)

---

## 為什麼選擇 yummycli？

- **Agent 原生設計** —— 開箱即用的結構化 [Skills](./skills/)，AI Agent 無需額外設定即可呼叫圖像、影片與音訊 API
- **能力優先架構** —— `image generate`、`video generate` 和 `audio speak` 是穩定的自動化介面；`gemini nanobanana`、`gemini veo` 和 `gemini speak` 是其上層的人性化快捷方式
- **結構化 JSON 輸出** —— 每條命令將結果寫入 stdout，方便 Agent、腳本和其他工具直接使用
- **安全的憑證儲存** —— API Key 儲存於作業系統原生金鑰鏈（macOS Keychain、Linux Secret Service），從不以明文儲存
- **多供應商自動降級** —— 設定兩個供應商後，主供應商失敗時自動切換到另一個

---

## 安裝

### 環境需求

- Node.js 16+ 及 `npm`
- Go 1.23+ 和 `make`（僅從原始碼建置時需要）

### 透過 npm 安裝（推薦）

```bash
# 安裝 CLI
npm install -g @yummysource/yummycli

# 安裝 Agent Skills（AI Agent 使用必須）
npx skills add yummysource/yummycli -y -g
```

驗證安裝：

```bash
yummycli version
```

### 從原始碼建置

```bash
git clone https://github.com/yummysource/yummycli.git
cd yummycli
make install

# 安裝 Agent Skills（AI Agent 使用必須）
npx skills add yummysource/yummycli -y -g
```

---

## 認證

yummycli 將每個供應商的 API Key 儲存於作業系統金鑰鏈中，每個供應商只需設定一次。

### 供應商覆蓋範圍

| 能力 | Gemini | OpenAI |
|---|---|---|
| 圖像生成與編輯 | ✅ | ✅ |
| 影片生成 | ✅ | — |
| 語音合成 | ✅ | — |

### 快速設定

```bash
# 設定主供應商並設為預設
yummycli init --provider gemini --api-key "AIza..." --default

# 可選：新增第二個供應商作為降級備選
yummycli init --provider openai --api-key "sk-..."
```

設定兩個供應商後，主供應商失敗時自動使用另一個。

### 查看設定狀態

```bash
yummycli auth list
```

```json
[
  {"provider":"gemini","configured":true,"default":true,"apiKeyPreview":"AIza******g7Aw"},
  {"provider":"openai","configured":true,"default":false,"apiKeyPreview":"sk-p******mEAA"}
]
```

`"default": true` 表示省略 `--provider` 時使用的供應商。

### 切換預設供應商

```bash
# Key 已儲存時無需重新輸入
yummycli init --provider openai --default
```

### 其他認證命令

| 命令 | 說明 |
|------|------|
| `yummycli init --provider <name> --api-key <key> [--default]` | 儲存 API Key，可選設為預設 |
| `yummycli auth list` | 列出所有供應商、憑證狀態及預設值 |
| `yummycli auth status --provider <name>` | 查看指定供應商的憑證狀態 |
| `yummycli auth remove --provider <name>` | 刪除指定供應商的憑證 |

### Gemini 快捷方式

`gemini init` 是 `init --provider gemini` 的等價快捷命令：

```bash
yummycli gemini init --api-key "AIza..." --default
```

---

## 圖像生成

圖像生成支援 **Gemini** 和 **OpenAI**。供應商從設定中自動解析，無需每次指定。

| 供應商 | 預設模型 | 預設尺寸 |
|---|---|---|
| Gemini | `gemini-3.1-flash-image` | `16:9`，`1K` |
| OpenAI | `gpt-image-2` | `1536x864`（16:9） |

兩個等價入口：

| 入口 | 適用場景 |
|------|----------|
| `yummycli image generate` | 所有供應商 —— 使用設定的預設供應商 |
| `yummycli gemini nanobanana` | Gemini 快捷方式，預設參數已填入 |

### 快速開始

```bash
# 第一步：設定供應商（一次性）
yummycli init --provider gemini --api-key "AIza..." --default
yummycli init --provider openai --api-key "sk-..."   # 可選降級備選

# 第二步：生成圖像（使用預設供應商）
yummycli image generate --prompt "白色盤子上放著一根成熟的香蕉，工作室打光"
```

### 參數說明

**通用參數（兩個供應商均支援）：**

| 參數 | 說明 | 預設值 |
|------|------|--------|
| `--prompt` | 圖像生成或編輯提示詞（**必填**） | — |
| `--provider` | 覆蓋供應商（`gemini` 或 `openai`） | 從設定讀取 |
| `--model` | 模型名稱 | 供應商預設值 |
| `--output` | 輸出檔案路徑 | 自動產生 |
| `--input-image` | 編輯用輸入圖像（可重複） | — |

**Gemini 專屬參數：**

| 參數 | 說明 | 預設值 |
|------|------|--------|
| `--aspect-ratio` | 命名比例：`16:9`、`9:16`、`21:9`、`1:1` 等 | `16:9` |
| `--image-size` | 解析度：`512`、`0.5K`、`1K`、`2K`、`4K` | `1K` |

**OpenAI 專屬參數：**

| 參數 | 說明 | 預設值 |
|------|------|--------|
| `--image-size` | 尺寸（WxH 格式，兩邊均為 16 的倍數，比例 1:3 到 3:1） | `1536x864` |
| `--quality` | 渲染品質：`low`、`medium`、`high`、`auto` | API 自動 |
| `--output-format` | 輸出格式：`png`、`jpeg`、`webp` | `png` |

### 文字生成圖像

```bash
# 使用設定的預設供應商
yummycli image generate --prompt "賽博龐克夜晚都市，霓虹燈倒映在街道上"

# 指定 Gemini，4K 解析度
yummycli image generate --provider gemini \
  --prompt "極簡主義 logo，扁平設計" \
  --image-size 4K --output logo.png

# 指定 OpenAI，JPEG 輸出
yummycli image generate --provider openai \
  --prompt "木桌上的一隻白兔" \
  --output-format jpeg
```

### 圖像編輯

透過 `--input-image` 傳入參考圖像，兩個供應商均支援編輯。

```bash
# 單圖編輯（Gemini）
yummycli image generate --provider gemini \
  --prompt "將這張圖轉換為水彩插畫風格" \
  --input-image ./photo.png

# 單圖編輯（OpenAI）
yummycli image generate --provider openai \
  --prompt "在背景中加入一道彩虹" \
  --input-image ./photo.png

# 多圖合成（Gemini）
yummycli image generate --provider gemini \
  --prompt "將這兩張參考圖融合為一張海報" \
  --input-image ./subject.png \
  --input-image ./background.jpg
```

支援的輸入格式：`.png`、`.jpg` / `.jpeg`、`.webp`。

### 寬高比

兩個供應商預設均為 16:9。指定其他比例：

| 比例 | Gemini 參數 | OpenAI 參數 |
|---|---|---|
| 21:9 超寬 | `--aspect-ratio 21:9` | `--image-size 2016x864` |
| 16:9 寬螢幕 | `--aspect-ratio 16:9` | `--image-size 1536x864` |
| 1:1 正方形 | `--aspect-ratio 1:1` | `--image-size 1024x1024` |
| 9:16 直式 | `--aspect-ratio 9:16` | `--image-size 864x1536` |

### 模型選擇

**Gemini：**
```bash
# Flash（預設）—— 速度較快，支援更多比例和尺寸
yummycli image generate --provider gemini --prompt "..." --model gemini-3.1-flash-image

# Pro —— 品質較高
yummycli image generate --provider gemini --prompt "..." --model gemini-3-pro-image-preview
```

**OpenAI：**
```bash
# gpt-image-2（預設）
yummycli image generate --provider openai --prompt "..."

# gpt-5.5
yummycli image generate --provider openai --prompt "..." --model gpt-5.5
```

未知 OpenAI 模型名稱會觸發警告並自動回退到 `gpt-image-2`。

### Gemini 模型相容性

**寬高比：**

| 模型 | 支援的值 |
|------|---------|
| `gemini-3.1-flash-image` | `1:1` `1:4` `1:8` `2:3` `3:2` `3:4` `4:1` `4:3` `4:5` `5:4` `8:1` `9:16` `16:9` `21:9` |
| `gemini-3-pro-image-preview` | `1:1` `2:3` `3:2` `3:4` `4:3` `4:5` `5:4` `9:16` `16:9` `21:9` |

`1:4`、`1:8`、`4:1`、`8:1` 僅 Flash 模型支援。

**圖像尺寸：**

| 模型 | 支援的值 |
|------|---------|
| `gemini-3.1-flash-image` | `512` `0.5K` `1K` `2K` `4K` |
| `gemini-3-pro-image-preview` | `1K` `2K` `4K` |

`512` 和 `0.5K` 僅 Flash 模型支援。

### JSON 輸出

```json
{
  "provider": "gemini",
  "output": "gemini_20260410123456_789.png",
  "model": "gemini-3.1-flash-image",
  "inputImageCount": 0
}
```

---

## 影片生成

影片生成透過 Google Veo 實現，僅支援 Gemini。

| 入口 | 適用場景 |
|------|----------|
| `gemini veo` | 人工使用 —— 已預設 Gemini Veo 預設參數 |
| `video generate --provider gemini` | 自動化/腳本 |

### 快速開始

```bash
yummycli gemini veo --prompt "陽光明媚的公園裡，黃金獵犬追逐紅球"
```

### 參數說明

| 參數 | 說明 | 預設值 |
|------|------|--------|
| `--prompt` | 影片生成提示詞（**必填**） | — |
| `--output` | 輸出檔案路徑（須以 `.mp4` 結尾） | 自動產生 |
| `--model` | Veo 模型 | `veo-3.1-fast-generate-preview` |
| `--aspect-ratio` | 影片寬高比 | `16:9` |
| `--duration` | 時長（秒） | `8` |
| `--resolution` | 影片解析度 | `1080p` |
| `--input-image` | 輸入圖像（可重複，最多 3 張） | — |

### 生成模式

| `--input-image` 數量 | 模式 |
|----------------------|------|
| 0 | 文字生成影片 |
| 1 | 圖像生成影片 —— 圖像作為起始幀 |
| 2–3 | 參考圖引導 —— 圖像作為 ASSET 參考輸入 |

### 模型相容性

**時長**（僅接受離散值）：

| 模型 | 有效時長（秒） |
|------|--------------|
| `veo-2.0-generate-001` | 5、6、7、8 |
| `veo-3.0-*` | 4、6、8 |
| `veo-3.1-*` | 4、6、8 |

**解析度：**

| 模型 | 支援的解析度 |
|------|------------|
| `veo-2.0-generate-001` | `720p` |
| `veo-3.0-*` | `720p`、`1080p` |
| `veo-3.1-*` | `720p`、`1080p`、`4k` |

`1080p` 和 `4k` 需要 `--duration 8`；`4k` 需要 veo-3.1 系列模型。

### JSON 輸出

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

## 語音合成

文字轉語音透過 Google Gemini TTS 實現，僅支援 Gemini。

| 入口 | 適用場景 |
|------|----------|
| `gemini speak` | 人工使用 —— 已預設 Gemini TTS 預設參數 |
| `audio speak --provider gemini` | 自動化/腳本 |

### 快速開始

```bash
yummycli gemini speak --text "你好，這是一段語音合成範例。"
```

### 參數說明

| 參數 | 說明 | 預設值 |
|------|------|--------|
| `--text` | 要合成的文字（**必填**） | — |
| `--output` | 輸出檔案路徑（須以 `.wav` 結尾） | 自動產生 |
| `--model` | TTS 模型 | `gemini-3.1-flash-tts-preview` |
| `--voice` | 預設聲音名稱 | `Aoede` |
| `--language` | BCP-47 語言代碼（省略時自動偵測） | — |
| `--speaker` | 多說話人映射 `名字:聲音`（可重複，最多 2 次） | — |

`--voice` 與 `--speaker` 互斥。

### 單說話人合成

```bash
yummycli gemini speak --text "今天天氣真好！"

yummycli gemini speak \
  --text "歡迎來到 AI 的世界。" \
  --voice Kore --output greeting.wav
```

### 多說話人對話合成

```bash
yummycli gemini speak \
  --text "[小明]: 你好！今天天氣真好。 [小紅]: 是啊，我們去公園走走吧！" \
  --speaker 小明:Aoede \
  --speaker 小紅:Kore \
  --output dialogue.wav
```

### 列出可用聲音

```bash
yummycli gemini voices
```

傳回所有 30 個預設聲音及其風格描述的 JSON 清單。

### JSON 輸出

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

yummycli 內建 Skills —— 結構化的指令檔案，協助 AI Agent 正確使用 CLI。

| Skill | 說明 |
|-------|------|
| [`yummy-shared`](./skills/yummy-shared/SKILL.md) | 初次設定、供應商設定、輸出格式約定和安全規則 |
| [`yummy-gen-image`](./skills/yummy-gen-image/SKILL.md) | 圖像生成與編輯 —— Gemini 和 OpenAI |
| [`yummy-gen-video`](./skills/yummy-gen-video/SKILL.md) | 影片生成 —— Gemini Veo |
| [`yummy-gen-voice`](./skills/yummy-gen-voice/SKILL.md) | 語音合成 —— Gemini TTS |

### 安裝

```bash
npx skills add yummysource/yummycli -y -g
```

---

## 命令參考

```
yummycli
├── version                                     顯示 yummycli 版本
│
├── init  --provider  [--api-key]  [--default]  設定供應商 API Key
│         --provider   gemini 或 openai（必填）
│         --api-key    API Key（未儲存時必填）
│         --default    設為預設供應商
│
├── auth
│   ├── list                                    列出所有供應商、憑證狀態及預設值
│   ├── status  --provider                      查看指定供應商的憑證狀態
│   └── remove  --provider                      刪除指定供應商的憑證
│
├── gemini
│   ├── init  --api-key  [--default]            初始化 Gemini 憑證
│   ├── nanobanana                              使用 Gemini 生成 / 編輯圖像
│   │     --prompt         （必填）
│   │     --output
│   │     --model          預設：gemini-3.1-flash-image
│   │     --aspect-ratio   預設：16:9
│   │     --image-size     預設：1K
│   │     --input-image    （可重複）
│   ├── veo                                     使用 Gemini Veo 生成影片
│   │     --prompt         （必填）
│   │     --output
│   │     --model          預設：veo-3.1-fast-generate-preview
│   │     --aspect-ratio   預設：16:9
│   │     --duration       預設：8
│   │     --resolution     預設：1080p
│   │     --input-image    （可重複，最多 3 張）
│   ├── speak                                   使用 Gemini TTS 合成語音
│   │     --text           （必填）
│   │     --output
│   │     --model          預設：gemini-3.1-flash-tts-preview
│   │     --voice          預設：Aoede
│   │     --language
│   │     --speaker        （可重複，最多 2 次；與 --voice 互斥）
│   └── voices                                  列出 Gemini TTS 可用聲音
│
├── image
│   └── generate                                生成或編輯圖像
│         --prompt         （必填）
│         --provider       省略時從設定讀取
│         --output
│         --model
│         --aspect-ratio   僅 Gemini
│         --image-size     Gemini: 1K/2K/4K — OpenAI: WxH 如 1536x864
│         --quality        僅 OpenAI：low/medium/high/auto
│         --output-format  僅 OpenAI：png/jpeg/webp
│         --input-image    （可重複）
│
├── video
│   └── generate                                生成影片（僅 Gemini）
│         --provider       （必填）
│         --prompt         （必填）
│         --output
│         --model
│         --aspect-ratio
│         --duration
│         --resolution
│         --input-image    （可重複，最多 3 張）
│
└── audio
    ├── speak                                   合成語音（僅 Gemini）
    │     --provider       （必填）
    │     --text           （必填）
    │     --output
    │     --model
    │     --voice
    │     --language
    │     --speaker        （可重複，最多 2 次；與 --voice 互斥）
    └── voices                                  列出指定供應商的可用聲音
          --provider       （必填）
```

---

## 貢獻

歡迎社群貢獻。如發現 Bug 或有功能建議，請提交 [Issue](https://github.com/yummysource/yummycli/issues) 或 [Pull Request](https://github.com/yummysource/yummycli/pulls)。

重大改動建議先透過 Issue 與我們討論。

## 授權條款

[MIT](./LICENSE)
