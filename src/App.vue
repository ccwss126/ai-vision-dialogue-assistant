<script setup>
import { ref } from 'vue'
import CameraView from './components/CameraView.vue'
import ChatHistory from './components/ChatHistory.vue'
import VoiceInput from './components/VoiceInput.vue'

const cameraView = ref(null)
const questionText = ref('')
const messages = ref([])
const loading = ref(false)
const errorMessage = ref('')
const aiCallCount = ref(0)
const isPerceptionActive = ref(false)

function speakAnswer(text) {
  if (!window.speechSynthesis || !text) return

  window.speechSynthesis.cancel()
  const utterance = new SpeechSynthesisUtterance(text)
  utterance.lang = 'zh-CN'
  window.speechSynthesis.speak(utterance)
}

async function sendToAi() {
  if (loading.value) return

  const finalQuestion = questionText.value.trim()

  if (!isPerceptionActive.value) {
    errorMessage.value = '请先开启感知，让 AI 获取摄像头和麦克风权限。'
    return
  }

  if (!finalQuestion) {
    errorMessage.value = '请先说出问题，或在输入框中输入问题。'
    return
  }

  const imageBase64 = cameraView.value?.captureFrame()

  if (!imageBase64) {
    errorMessage.value = '请先等待摄像头画面加载完成。'
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    const response = await fetch('http://localhost:8080/api/analyze', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        imageBase64,
        question: finalQuestion,
      }),
    })

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.error || 'AI 分析失败')
    }

    aiCallCount.value += 1
    messages.value.push({
      question: finalQuestion,
      answer: data.answer,
      mode: data.mode,
    })
    speakAnswer(data.answer)
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    loading.value = false
  }
}

function handlePerceptionChange(active) {
  isPerceptionActive.value = active
  if (active) {
    errorMessage.value = ''
  }
}
</script>

<template>
  <main class="app-shell">
    <header class="app-header">
      <div>
        <p class="eyebrow">AI Vision Assistant</p>
        <h1>摄像头与语音 AI 对话</h1>
      </div>
      <div class="status-pill">AI 调用 {{ aiCallCount }} 次</div>
    </header>

    <section class="workspace">
      <div class="preview-panel">
        <div class="panel-heading">
          <h2>摄像头画面</h2>
          <span>
            {{ isPerceptionActive ? 'AI 正在看见画面，并已开启麦克风' : '点击开启感知以授权摄像头和麦克风' }}
          </span>
        </div>
        <CameraView ref="cameraView" @perception-change="handlePerceptionChange" />
        <div class="perception-input">
          <VoiceInput
            v-model="questionText"
            :disabled="loading"
            :is-perception-active="isPerceptionActive"
          />
          <button class="send-button" type="button" :disabled="loading" @click="sendToAi">
            {{ loading ? '分析中...' : '分析当前画面' }}
          </button>
        </div>
      </div>

      <aside class="side-panel">
        <div class="panel-heading">
          <h2>语音问题</h2>
          <span>开启感知后自动识别，修改后点击分析</span>
        </div>

        <div class="question-box">
          <span>当前问题</span>
          <p>{{ questionText || '等待语音识别或手动输入...' }}</p>
        </div>
        <p v-if="errorMessage" class="error-message">{{ errorMessage }}</p>
      </aside>
    </section>

    <section class="history-panel">
      <div class="panel-heading">
        <h2>对话历史</h2>
        <span>{{ messages.length }} 条记录</span>
      </div>
      <ChatHistory :messages="messages" />
    </section>
  </main>
</template>
