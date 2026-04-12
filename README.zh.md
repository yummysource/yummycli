# yummycli

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23-blue.svg)](https://go.dev/)
[![npm version](https://img.shields.io/npm/v/@yummysource/yummycli.svg)](https://www.npmjs.com/package/@yummysource/yummycli)

面向多模态模型供应商的 AI 友好命令行工具 —— 专为人类用户和 AI Agent 设计。

当前支持通过 [Gemini](https://deepmind.google/technologies/gemini/) 进行图像生成与编辑，Claude、OpenAI、Qwen 等供应商支持正在规划中。

[繁體中文](./README.zh-TW.md) | [简体中文](./README.zh.md) | [English](./README.md)

<img src="./assets/logo.png" alt="yummycli logo" width="120" />

[安装](#安装) · [认证](#认证) · [图像生成](#图像生成) · [Agent Skills](#agent-skills) · [命令参考](#命令参考)

---

## 为什么选择 yummycli？

- **Agent 原生设计** —— 开箱即用的结构化 [Skills](./skills/)，AI Agent 无需额外配置即可调用图像 API
- **能力优先架构** —— `image generate` 是稳定的自动化接口；`gemini nanobanana` 是其上层的人性化快捷方式
- **结构化 JSON 输出** —— 每条命令将结果写入 stdout，方便 Agent、脚本和其他工具直接消费
- **安全的凭证存储** —— API Key 存储在操作系统原生密钥链（macOS Keychain、Linux Secret Service），从不以明文保存
- **供应商无关** —— 统一的 CLI 接口，新增供应商时无需修改现有脚本

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

### 命令

| 命令 | 说明 |
|------|------|
| `auth init` | 保存供应商的 API Key |
| `auth list` | 列出所有供应商及其凭证状态 |
| `auth status` | 查看指定供应商的凭证状态 |
| `auth remove` | 删除指定供应商的凭证 |

### 示例

```bash
# 保存 Gemini API Key
yummycli auth init --provider gemini --api-key "AIza..."

# 查看 Gemini 是否已配置（显示脱敏预览）
yummycli auth status --provider gemini

# 删除 Gemini 凭证
yummycli auth remove --provider gemini
```

**`auth init` 输出：**

```json
{"provider":"gemini","configured":true}
```

**`auth list` 输出：**

```json
[{"provider":"gemini","configured":true,"apiKeyPreview":"AIza...xxxx"}]
```

**`auth status` 输出：**

```json
{"provider":"gemini","configured":true,"apiKeyPreview":"AIza...xxxx"}
```

### Gemini 快捷方式

`gemini init` 是 `auth init --provider gemini` 的等价快捷命令：

```bash
yummycli gemini init --api-key "AIza..."
```

---

## 图像生成

yummycli 提供两个等价的图像生成入口：

| 入口 | 适用场景 |
|------|----------|
| `gemini nanobanana` | 人工使用 —— 已预设 Gemini 默认参数 |
| `image generate --provider gemini` | 自动化/脚本 —— 显式、稳定的接口 |

两者调用相同的底层实现，按使用场景选择即可。

### 快速开始

```bash
# 第一步：配置 Gemini 凭证（一次性）
yummycli gemini init --api-key "AIza..."

# 第二步：根据文本提示生成图像
yummycli gemini nanobanana --prompt "白色盘子上放着一根成熟的香蕉，工作室打光"
```

生成的图像保存在当前目录，文件名自动生成：

```
gemini_20260410123456_789.png
```

### 参数说明

| 参数 | 说明 | 默认值（Gemini） |
|------|------|----------------|
| `--prompt` | 图像生成提示词（**必填**） | — |
| `--output` | 输出文件路径 | 自动生成 |
| `--model` | Gemini 模型 | `gemini-3.1-flash-image-preview` |
| `--aspect-ratio` | 图像宽高比 | `16:9` |
| `--image-size` | 输出分辨率 | `1K` |
| `--input-image` | 输入图像（可重复使用） | — |

> `image generate` 使用时需额外传入 `--provider gemini`（必填）。Gemini 默认值同样适用 —— 省略时 `--model`、`--aspect-ratio`、`--image-size` 均自动填充。

### 文本生成图像

```bash
yummycli gemini nanobanana \
  --prompt "赛博朋克夜晚都市，霓虹灯倒映在湿漉漉的街道上"

# 指定输出路径和分辨率
yummycli gemini nanobanana \
  --prompt "极简主义 logo，扁平设计，白色背景" \
  --output logo.png \
  --image-size 4K
```

### 图像编辑

通过 `--input-image` 传入一张或多张参考图像：

```bash
# 单图编辑
yummycli gemini nanobanana \
  --prompt "将这张图转换为水彩插画风格" \
  --input-image ./photo.png

# 多图参考
yummycli gemini nanobanana \
  --prompt "将这两张参考图融合为一张精致的海报插画" \
  --input-image ./subject.png \
  --input-image ./background.jpg
```

支持的输入格式：`.png`、`.jpg` / `.jpeg`、`.webp`。

### 宽高比

```bash
# 竖版 —— 手机壁纸、故事格式
yummycli gemini nanobanana --prompt "..." --aspect-ratio 9:16

# 正方形 —— 社交头像、图标
yummycli gemini nanobanana --prompt "..." --aspect-ratio 1:1

# 宽屏 —— 桌面壁纸、横幅
yummycli gemini nanobanana --prompt "..." --aspect-ratio 21:9
```

### 模型选择

```bash
# Flash（默认）—— 速度更快，支持更多宽高比和更小尺寸
yummycli gemini nanobanana --prompt "..." --model gemini-3.1-flash-image-preview

# Pro —— 质量更高，支持的尺寸和宽高比选项较少
yummycli gemini nanobanana --prompt "..." --model gemini-3-pro-image-preview
```

### 模型兼容性

**宽高比**

| 模型 | 支持的值 |
|------|---------|
| `gemini-3.1-flash-image-preview` | `1:1` `1:4` `1:8` `2:3` `3:2` `3:4` `4:1` `4:3` `4:5` `5:4` `8:1` `9:16` `16:9` `21:9` |
| `gemini-3-pro-image-preview` | `1:1` `2:3` `3:2` `3:4` `4:3` `4:5` `5:4` `9:16` `16:9` `21:9` |

`1:4`、`1:8`、`4:1`、`8:1` 仅 Flash 模型支持，Pro 模型不可用。

**图像尺寸**

| 模型 | 支持的值 |
|------|---------|
| `gemini-3.1-flash-image-preview` | `512` `0.5K` `1K` `2K` `4K` |
| `gemini-3-pro-image-preview` | `1K` `2K` `4K` |

`512` 和 `0.5K` 仅 Flash 模型支持。尺寸值大小写不敏感（`4k` 和 `4K` 均可）。

### JSON 输出

每次成功生成后向 stdout 写入结果：

```json
{
  "provider": "gemini",
  "output": "gemini_20260410123456_789.png",
  "model": "gemini-3.1-flash-image-preview",
  "inputImageCount": 0
}
```

使用 `output` 字段定位生成的文件。

### 直接使用 `image generate`

`image generate` 是供应商无关的稳定 API，接受相同的参数，但需要显式传入 `--provider`：

```bash
yummycli image generate \
  --provider gemini \
  --prompt "日出时宁静的山中湖泊" \
  --aspect-ratio 16:9 \
  --image-size 2K \
  --output landscape.png
```

推荐在脚本和 AI Agent 中使用此形式 —— 新增供应商后无需修改。

---

## Agent Skills

yummycli 内置 Skills —— 结构化的指令文件，帮助 AI Agent 正确使用 CLI。

| Skill | 说明 |
|-------|------|
| [`yummy-shared`](./skills/yummy-shared/SKILL.md) | 凭证检查、输出格式约定和共享安全规则 —— 所有其他 Skill 加载前自动使用 |
| [`yummy-gen-image`](./skills/yummy-gen-image/SKILL.md) | 通过 Gemini 进行文本生成图像、单图编辑和多图参考编辑 |

Skills 位于 [`./skills/`](./skills/) 目录。

### 安装

```bash
npx skills add yummysource/yummycli -y -g
```

使用任何其他 yummycli Skill 前，请先加载 `yummy-shared`。

---

## 命令参考

```
yummycli
├── version                              显示 yummycli 版本
│
├── auth
│   ├── init    --provider  --api-key    保存供应商 API Key
│   ├── list                             列出所有供应商及凭证状态
│   ├── status  --provider               查看指定供应商的凭证状态
│   └── remove  --provider               删除指定供应商的凭证
│
├── gemini
│   ├── init  --api-key                  初始化 Gemini 凭证
│   └── nanobanana                       使用 Gemini 生成 / 编辑图像
│         --prompt        （必填）
│         --output
│         --model
│         --aspect-ratio
│         --image-size
│         --input-image   （可重复）
│
└── image
    └── generate                         供应商无关的图像生成接口
          --provider      （必填）
          --prompt        （必填）
          --output
          --model
          --aspect-ratio
          --image-size
          --input-image   （可重复）
```

---

## 贡献

欢迎社区贡献。如发现 Bug 或有功能建议，请提交 [Issue](https://github.com/yummysource/yummycli/issues) 或 [Pull Request](https://github.com/yummysource/yummycli/pulls)。

重大改动建议先通过 Issue 与我们讨论。

## 许可证

[MIT](./LICENSE)
