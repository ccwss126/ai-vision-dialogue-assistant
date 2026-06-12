# AI Vision Dialogue Assistant

## 项目简介

AI Vision Dialogue Assistant 是一款 AI 视觉对话助手。用户点击“开启感知”后，浏览器会同时请求摄像头和麦克风权限；应用会展示摄像头画面，自动识别用户语音问题，并在用户点击“分析当前画面”后，把当前视频帧和问题发送给 Go 后端，由 Gemini Vision API 返回回答。

AI 回复会展示在对话历史中，并可通过浏览器语音播报功能自动朗读。

## 技术栈

前端：

* Vue3
* Vite
* getUserMedia
* Web Speech API
* SpeechSynthesis API

后端：

* Go
* 标准库 `net/http`
* Gemini Vision API

## 项目结构

```text
.
├── backend/
│   ├── main.go          # Go 后端，提供 /api/analyze
│   ├── go.mod
│   └── README.md
├── docs/
│   └── design.md        # 设计文档
├── src/
│   ├── components/
│   │   ├── CameraView.vue
│   │   ├── VoiceInput.vue
│   │   └── ChatHistory.vue
│   ├── App.vue
│   └── style.css
├── .env.example
├── README.md
├── package.json
└── vite.config.js
```

## 启动方式

安装前端依赖：

```bash
npm install
```

启动前端：

```bash
npm run dev
```

启动后端：

```bash
cd backend
go run main.go
```

默认地址：

* 前端：http://localhost:5173
* 后端：http://localhost:8080
* 分析接口：`POST http://localhost:8080/api/analyze`

如果 Windows 本地 Go 缓存目录没有权限，可以临时指定项目内缓存目录：

```bash
set GOCACHE=D:\jrx\Desktop\ai-vision-assistant\.gocache
cd backend
go run main.go
```

## 环境变量说明

复制 `.env.example` 为 `.env`：

```env
GEMINI_API_KEY=
PORT=8080
HTTP_PROXY=http://127.0.0.1:7890
HTTPS_PROXY=http://127.0.0.1:7890
```

变量说明：

* `GEMINI_API_KEY`：Gemini API Key。为空时自动进入 Demo Mode。
* `PORT`：Go 后端端口，默认 `8080`。
* `HTTP_PROXY` / `HTTPS_PROXY`：可选代理配置，用于本地网络访问 Gemini API。

不要提交真实 `.env`。

## Demo Mode 说明

以下情况会自动进入 Demo Mode：

* 未配置 `GEMINI_API_KEY`
* Gemini API 请求失败
* 网络超时或代理不可用

Demo Mode 会返回模拟回答，保证评委可以复现完整交互流程，不会因为云端配置或网络问题导致前端崩溃。

## 成本控制策略

已实现：

* 图片压缩到 640px
* 问题为空不发送
* 请求期间禁止重复点击
* AI 调用次数统计
* Gemini 失败自动 Demo Mode

设计取舍：

* 不上传连续视频流，只在用户点击分析时截取当前帧。
* 不上传原始麦克风音频，语音识别在浏览器端完成。
* 后端默认使用轻量 Gemini 模型，便于控制调用成本。

## 功能展示截图占位

截图 1：开启感知后的摄像头预览

```text
TODO: 添加截图链接或图片
```

截图 2：语音识别问题与 AI 回答

```text
TODO: 添加截图链接或图片
```

截图 3：语音播报开关与对话历史

```text
TODO: 添加截图链接或图片
```

## Demo 视频链接占位

```text
TODO: 添加 Demo 视频链接
```

## 设计文档

见 [docs/design.md](docs/design.md)。
