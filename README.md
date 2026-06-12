# AI 视觉对话助手

一款摄像头与麦克风驱动的 AI 对话应用。用户点击“开启感知”后，浏览器会同时请求摄像头和麦克风权限；AI 在用户确认语音问题后获取当前摄像头画面截图，并结合语音识别文本给出回答。回答会展示在页面中，并通过浏览器 SpeechSynthesis 自动朗读。

## 功能

- 一键“开启感知”：使用 `getUserMedia({ video: true, audio: true })` 同时打开摄像头和麦克风权限。
- 摄像头实时预览：AI 在发送时看到当前画面截图。
- 中文语音提问：浏览器 Web Speech API 负责语音识别。
- ChatGPT 风格输入区：识别文本先预览，用户点击勾号确认后自动发送。
- Go 后端：`backend/main.go` 提供 `POST /api/analyze`。
- Gemini Vision：后端把 `imageBase64 + question` 发送给 Gemini。
- 朗读回答：前端用 SpeechSynthesis 播放 AI 回复。
- 成本控制：空问题不请求、发送中禁止重复请求、图片压缩到 640px、记录 AI 调用次数、Gemini 失败自动 Demo Mode。

## Demo

Demo 视频链接：待补充

## 环境变量

复制 `.env.example` 为 `.env`，并按需填写：

```env
GEMINI_API_KEY=your_gemini_api_key_here
PORT=8080
HTTP_PROXY=http://127.0.0.1:7890
HTTPS_PROXY=http://127.0.0.1:7890
```

不填写 `GEMINI_API_KEY` 时，Go 后端会自动进入 Demo Mode。Gemini 网络请求失败时也会返回 Demo Mode，不会让前端崩溃。

## 启动

前端：

```bash
npm install
npm run dev
```

后端：

```bash
cd backend
go run main.go
```

默认地址：

- 前端：http://localhost:5173
- 后端：http://localhost:8080
- 接口：`POST http://localhost:8080/api/analyze`

## 提交物

- GitHub/Gitee 仓库：当前仓库
- README：`README.md`
- 设计文档：`docs/design.md`
- Demo 视频：替换上方 Demo 链接占位
