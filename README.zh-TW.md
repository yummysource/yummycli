# yummycli

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23-blue.svg)](https://go.dev/)
[![npm version](https://img.shields.io/npm/v/@yummysource/yummycli.svg)](https://www.npmjs.com/package/@yummysource/yummycli)

面向多模態模型供應商的 AI 友好命令列工具 —— 專為人類用戶和 AI Agent 設計。

目前支援透過 [Gemini](https://deepmind.google/technologies/gemini/) 進行圖像生成與編輯，Claude、OpenAI、Qwen 等供應商支援正在規劃中。

[繁體中文](./README.zh-TW.md) | [简体中文](./README.zh.md) | [English](./README.md)

<img src="./assets/logo.png" alt="yummycli logo" width="120" />

[安裝](#安裝) · [認證](#認證) · [圖像生成](#圖像生成) · [Agent Skills](#agent-skills) · [命令參考](#命令參考)

---

## 為什麼選擇 yummycli？

- **Agent 原生設計** —— 開箱即用的結構化 [Skills](./skills/)，AI Agent 無需額外設定即可呼叫圖像 API
- **能力優先架構** —— `image generate` 是穩定的自動化介面；`gemini nanobanana` 是其上層的人性化快捷方式
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

## Agent Skills

yummycli 內建 Skills —— 結構化的指令檔案，協助 AI Agent 正確使用 CLI。

| Skill | 說明 |
|-------|------|
| [`yummy-shared`](./skills/yummy-shared/SKILL.md) | 憑證檢查、輸出格式約定和共用安全規則 —— 所有其他 Skill 載入前自動使用 |
| [`generate-image`](./skills/generate-image/SKILL.md) | 透過 Gemini 進行文字生成圖像、單圖編輯和多圖參考編輯 |

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
│   └── nanobanana                       使用 Gemini 生成 / 編輯圖像
│         --prompt        （必填）
│         --output
│         --model
│         --aspect-ratio
│         --image-size
│         --input-image   （可重複）
│
└── image
    └── generate                         供應商無關的圖像生成介面
          --provider      （必填）
          --prompt        （必填）
          --output
          --model
          --aspect-ratio
          --image-size
          --input-image   （可重複）
```

---

## 貢獻

歡迎社群貢獻。如發現 Bug 或有功能建議，請提交 [Issue](https://github.com/yummysource/yummycli/issues) 或 [Pull Request](https://github.com/yummysource/yummycli/pulls)。

重大改動建議先透過 Issue 與我們討論。

## 授權條款

[MIT](./LICENSE)
