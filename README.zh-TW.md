# yummycli

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23-blue.svg)](https://go.dev/)
[![npm version](https://img.shields.io/npm/v/@yummysource/yummycli.svg)](https://www.npmjs.com/package/@yummysource/yummycli)

面向多模態模型供應商的 AI 友好命令列工具 —— 專為人類用戶和 AI Agent 設計。

目前支援透過 [Gemini](https://deepmind.google/technologies/gemini/) 進行圖像生成與編輯、影片生成和語音合成，Claude、OpenAI、Qwen 等供應商支援正在規劃中。

[繁體中文](./README.zh-TW.md) | [简体中文](./README.zh.md) | [English](./README.md)

<img src="./assets/logo.png" alt="yummycli logo" width="120" />

[安裝](#安裝) · [認證](#認證) · [圖像生成](#圖像生成) · [影片生成](#影片生成) · [語音合成](#語音合成) · [Agent Skills](#agent-skills) · [命令參考](#命令參考)

---

## 為什麼選擇 yummycli？

- **Agent 原生設計** —— 開箱即用的結構化 [Skills](./skills/)，AI Agent 無需額外設定即可呼叫圖像、影片與音訊 API
- **能力優先架構** —— `image generate`、`video generate` 和 `audio speak` 是穩定的自動化介面；`gemini nanobanana`、`gemini veo` 和 `gemini speak` 是其上層的人性化快捷方式
- **結構化 JSON 輸出** —— 每條命令將結果寫入 stdout，方便 Agent、腳本和其他工具直接使用
- **安全的憑證儲存** —— API Key 儲存於作業系統原生金鑰鏈（macOS Keychain、Linux Secret Service），從不以明文儲存
- **供應商無關** —— 統一的 CLI 介面，新增供應商時無需修改現有腳本

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

### 命令

| 命令 | 說明 |
|------|------|
| `auth init` | 儲存供應商的 API Key |
| `auth list` | 列出所有供應商及其憑證狀態 |
| `auth status` | 查看指定供應商的憑證狀態 |
| `auth remove` | 刪除指定供應商的憑證 |

### 範例

```bash
# 儲存 Gemini API Key
yummycli auth init --provider gemini --api-key "AIza..."

# 查看 Gemini 是否已設定（顯示遮罩預覽）
yummycli auth status --provider gemini

# 刪除 Gemini 憑證
yummycli auth remove --provider gemini
```

**`auth init` 輸出：**

```json
{"provider":"gemini","configured":true}
```

**`auth list` 輸出：**

```json
[{"provider":"gemini","configured":true,"apiKeyPreview":"AIza...xxxx"}]
```

**`auth status` 輸出：**

```json
{"provider":"gemini","configured":true,"apiKeyPreview":"AIza...xxxx"}
```

### Gemini 快捷方式

`gemini init` 是 `auth init --provider gemini` 的等價快捷命令：

```bash
yummycli gemini init --api-key "AIza..."
```

---

## 圖像生成

yummycli 提供兩個等價的圖像生成入口：

| 入口 | 適用場景 |
|------|----------|
| `gemini nanobanana` | 人工使用 —— 已預設 Gemini 預設參數 |
| `image generate --provider gemini` | 自動化/腳本 —— 明確、穩定的介面 |

兩者呼叫相同的底層實作，依使用場景選擇即可。

### 快速開始

```bash
# 第一步：設定 Gemini 憑證（一次性）
yummycli gemini init --api-key "AIza..."

# 第二步：根據文字提示生成圖像
yummycli gemini nanobanana --prompt "白色盤子上放著一根成熟的香蕉，工作室打光"
```

生成的圖像儲存在目前目錄，檔名自動產生：

```
gemini_20260410123456_789.png
```

### 參數說明

| 參數 | 說明 | 預設值（Gemini） |
|------|------|----------------|
| `--prompt` | 圖像生成提示詞（**必填**） | — |
| `--output` | 輸出檔案路徑 | 自動產生 |
| `--model` | Gemini 模型 | `gemini-3.1-flash-image-preview` |
| `--aspect-ratio` | 圖像寬高比 | `16:9` |
| `--image-size` | 輸出解析度 | `1K` |
| `--input-image` | 輸入圖像（可重複使用） | — |

> `image generate` 使用時需額外傳入 `--provider gemini`（必填）。Gemini 預設值同樣適用 —— 省略時 `--model`、`--aspect-ratio`、`--image-size` 均自動填入。

### 文字生成圖像

```bash
yummycli gemini nanobanana \
  --prompt "賽博龐克夜晚都市，霓虹燈倒映在濕漉漉的街道上"

# 指定輸出路徑和解析度
yummycli gemini nanobanana \
  --prompt "極簡主義 logo，扁平設計，白色背景" \
  --output logo.png \
  --image-size 4K
```

### 圖像編輯

透過 `--input-image` 傳入一張或多張參考圖像：

```bash
# 單圖編輯
yummycli gemini nanobanana \
  --prompt "將這張圖轉換為水彩插畫風格" \
  --input-image ./photo.png

# 多圖參考
yummycli gemini nanobanana \
  --prompt "將這兩張參考圖融合為一張精緻的海報插畫" \
  --input-image ./subject.png \
  --input-image ./background.jpg
```

支援的輸入格式：`.png`、`.jpg` / `.jpeg`、`.webp`。

### 寬高比

```bash
# 直式 —— 手機桌布、限時動態格式
yummycli gemini nanobanana --prompt "..." --aspect-ratio 9:16

# 正方形 —— 社群頭像、圖示
yummycli gemini nanobanana --prompt "..." --aspect-ratio 1:1

# 寬螢幕 —— 桌面桌布、橫幅
yummycli gemini nanobanana --prompt "..." --aspect-ratio 21:9
```

### 模型選擇

```bash
# Flash（預設）—— 速度較快，支援更多寬高比和較小尺寸
yummycli gemini nanobanana --prompt "..." --model gemini-3.1-flash-image-preview

# Pro —— 品質較高，支援的尺寸和寬高比選項較少
yummycli gemini nanobanana --prompt "..." --model gemini-3-pro-image-preview
```

### 模型相容性

**寬高比**

| 模型 | 支援的值 |
|------|---------|
| `gemini-3.1-flash-image-preview` | `1:1` `1:4` `1:8` `2:3` `3:2` `3:4` `4:1` `4:3` `4:5` `5:4` `8:1` `9:16` `16:9` `21:9` |
| `gemini-3-pro-image-preview` | `1:1` `2:3` `3:2` `3:4` `4:3` `4:5` `5:4` `9:16` `16:9` `21:9` |

`1:4`、`1:8`、`4:1`、`8:1` 僅 Flash 模型支援，Pro 模型不可用。

**圖像尺寸**

| 模型 | 支援的值 |
|------|---------|
| `gemini-3.1-flash-image-preview` | `512` `0.5K` `1K` `2K` `4K` |
| `gemini-3-pro-image-preview` | `1K` `2K` `4K` |

`512` 和 `0.5K` 僅 Flash 模型支援。尺寸值大小寫不敏感（`4k` 和 `4K` 均可）。

### JSON 輸出

每次成功生成後向 stdout 寫入結果：

```json
{
  "provider": "gemini",
  "output": "gemini_20260410123456_789.png",
  "model": "gemini-3.1-flash-image-preview",
  "inputImageCount": 0
}
```

使用 `output` 欄位定位生成的檔案。

### 直接使用 `image generate`

`image generate` 是供應商無關的穩定 API，接受相同的參數，但需要明確傳入 `--provider`：

```bash
yummycli image generate \
  --provider gemini \
  --prompt "日出時寧靜的山中湖泊" \
  --aspect-ratio 16:9 \
  --image-size 2K \
  --output landscape.png
```

推薦在腳本和 AI Agent 中使用此形式 —— 新增供應商後無需修改。

---

### 直接使用 `video generate`

`video generate` 是供應商無關的穩定 API，接受相同的參數，但需要明確傳入 `--provider`：

```bash
yummycli video generate \
  --provider gemini \
  --prompt "雲朵在山頂緩緩飄過的縮時攝影" \
  --resolution 1080p \
  --output timelapse.mp4
```

---

## 影片生成

yummycli 透過 Google Veo 支援影片生成，提供兩個等價入口：

| 入口 | 適用場景 |
|------|----------|
| `gemini veo` | 人工使用 —— 已預設 Gemini Veo 預設參數 |
| `video generate --provider gemini` | 自動化/腳本 —— 明確、穩定的介面 |

### 快速開始

```bash
# 第一步：設定 Gemini 憑證（一次性）
yummycli gemini init --api-key "AIza..."

# 第二步：根據文字提示生成影片
yummycli gemini veo --prompt "陽光明媚的公園裡，黃金獵犬追逐紅球"
```

生成的影片儲存在目前目錄，檔名自動產生：

```
veo_20260417_142301_047.mp4
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
| `--input-image` | 輸入圖像（可重複使用，最多 3 張） | — |

### 生成模式

`--input-image` 可重複使用，數量決定生成模式：

| `--input-image` 數量 | 模式 |
|----------------------|------|
| 0 | 文字生成影片 |
| 1 | 圖像生成影片 —— 圖像作為起始幀 |
| 2–3 | 參考圖引導 —— 圖像作為 ASSET 參考輸入 |

```bash
# 文字生成影片
yummycli gemini veo --prompt "山頂雲海的縮時攝影，黃金時段"

# 圖像生成影片（為靜態圖像添加動效）
yummycli gemini veo \
  --prompt "小狗向鏡頭跑來" \
  --input-image ./dog.jpg

# 參考圖引導（兩張圖）
yummycli gemini veo \
  --prompt "將角色融入這個背景環境中" \
  --input-image ./character.png \
  --input-image ./background.jpg
```

### 模型相容性

**時長**僅接受離散值：

| 模型 | 有效時長（秒） |
|------|--------------|
| `veo-2.0-generate-001` | 5、6、7、8 |
| `veo-3.0-*` | 4、6、8 |
| `veo-3.1-*` | 4、6、8 |

**解析度**支援情況：

| 模型 | 支援的解析度 |
|------|------------|
| `veo-2.0-generate-001` | `720p` |
| `veo-3.0-*` | `720p`、`1080p` |
| `veo-3.1-*` | `720p`、`1080p`、`4k` |

約束：`1080p` 和 `4k` 需要 `--duration 8`；`4k` 需要 veo-3.1 系列模型。

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

yummycli 透過 Google Gemini TTS 支援語音合成，提供兩個等價入口：

| 入口 | 適用場景 |
|------|----------|
| `gemini speak` | 人工使用 —— 已預設 Gemini TTS 預設參數 |
| `audio speak --provider gemini` | 自動化/腳本 —— 明確、穩定的介面 |

### 快速開始

```bash
# 合成語音（自動產生檔名）
yummycli gemini speak --text "你好，這是一段語音合成範例。"

# 指定聲音和輸出路徑
yummycli gemini speak \
  --text "歡迎使用 AI 語音合成服務。" \
  --voice Puck \
  --output welcome.wav
```

生成的音訊儲存在目前目錄，檔名自動產生：

```
tts_20260420_142301_047.wav
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
# 預設聲音（Aoede，輕盈風格）
yummycli gemini speak --text "今天天氣真好！"

# 指定聲音和輸出
yummycli gemini speak \
  --text "歡迎來到 AI 的世界。" \
  --voice Kore \
  --output greeting.wav
```

### 多說話人對話合成

在文字中以 `[名字]:` 標記每位說話人的台詞，再用 `--speaker` 映射聲音：

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

單說話人：

```json
{
  "provider": "gemini",
  "output": "tts_20260420_142301_047.wav",
  "model": "gemini-3.1-flash-tts-preview",
  "voice": "Aoede",
  "elapsed_seconds": 3
}
```

多說話人：

```json
{
  "provider": "gemini",
  "output": "dialogue_20260420_143010_112.wav",
  "model": "gemini-3.1-flash-tts-preview",
  "speakers": [
    {"name": "小明", "voice": "Aoede"},
    {"name": "小紅", "voice": "Kore"}
  ],
  "elapsed_seconds": 4
}
```

---

## Agent Skills

yummycli 內建 Skills —— 結構化的指令檔案，協助 AI Agent 正確使用 CLI。

| Skill | 說明 |
|-------|------|
| [`yummy-shared`](./skills/yummy-shared/SKILL.md) | 憑證檢查、輸出格式約定和共用安全規則 —— 所有其他 Skill 載入前自動使用 |
| [`yummy-gen-image`](./skills/yummy-gen-image/SKILL.md) | 透過 Gemini 進行文字生成圖像、單圖編輯和多圖參考編輯 |
| [`yummy-gen-video`](./skills/yummy-gen-video/SKILL.md) | 透過 Gemini Veo 進行文字生成影片、圖像生成影片和參考圖引導影片生成 |
| [`yummy-gen-voice`](./skills/yummy-gen-voice/SKILL.md) | 透過 Gemini TTS 進行單說話人語音合成、多說話人對話合成和聲音清單查詢 |

Skills 位於 [`./skills/`](./skills/) 目錄。

### 安裝

```bash
npx skills add yummysource/yummycli -y -g
```

使用任何其他 yummycli Skill 前，請先載入 `yummy-shared`。

---

## 命令參考

```
yummycli
├── version                              顯示 yummycli 版本
│
├── auth
│   ├── init    --provider  --api-key    儲存供應商 API Key
│   ├── list                             列出所有供應商及憑證狀態
│   ├── status  --provider               查看指定供應商的憑證狀態
│   └── remove  --provider               刪除指定供應商的憑證
│
├── gemini
│   ├── init  --api-key                  初始化 Gemini 憑證
│   ├── nanobanana                       使用 Gemini 生成 / 編輯圖像
│   │     --prompt        （必填）
│   │     --output
│   │     --model
│   │     --aspect-ratio
│   │     --image-size
│   │     --input-image   （可重複）
│   ├── veo                              使用 Gemini Veo 生成影片
│   │     --prompt        （必填）
│   │     --output
│   │     --model
│   │     --aspect-ratio
│   │     --duration
│   │     --resolution
│   │     --input-image   （可重複，最多 3 張）
│   ├── speak                            使用 Gemini TTS 合成語音
│   │     --text          （必填）
│   │     --output
│   │     --model
│   │     --voice
│   │     --language
│   │     --speaker       （可重複，最多 2 次；與 --voice 互斥）
│   └── voices                           列出 Gemini TTS 可用聲音
│
├── image
│   └── generate                         供應商無關的圖像生成介面
│         --provider      （必填）
│         --prompt        （必填）
│         --output
│         --model
│         --aspect-ratio
│         --image-size
│         --input-image   （可重複）
│
├── video
│   └── generate                         供應商無關的影片生成介面
│         --provider      （必填）
│         --prompt        （必填）
│         --output
│         --model
│         --aspect-ratio
│         --duration
│         --resolution
│         --input-image   （可重複，最多 3 張）
│
└── audio
    ├── speak                            供應商無關的語音合成介面
    │     --provider      （必填）
    │     --text          （必填）
    │     --output
    │     --model
    │     --voice
    │     --language
    │     --speaker       （可重複，最多 2 次；與 --voice 互斥）
    └── voices                           列出指定供應商的可用聲音
          --provider      （必填）
```

---

## 貢獻

歡迎社群貢獻。如發現 Bug 或有功能建議，請提交 [Issue](https://github.com/yummysource/yummycli/issues) 或 [Pull Request](https://github.com/yummysource/yummycli/pulls)。

重大改動建議先透過 Issue 與我們討論。

## 授權條款

[MIT](./LICENSE)
