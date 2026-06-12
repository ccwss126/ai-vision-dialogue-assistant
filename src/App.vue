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

async function sendToAi() {
  if (loading.value) return

  const question = questionText.value.trim()
  const imageBase64 = cameraView.value?.captureFrame()

  if (!imageBase64) {
    errorMessage.value = '请先开启摄像头并等待画面加载。'
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    const response = await fetch('/api/analyze', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        imageBase64,
        question,
      }),
    })

    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.error || 'AI 分析失败')
    }

    messages.value.push({
      question,
      image: imageBase64,
      answer: data.answer,
    })
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <CameraView ref="cameraView" />
  <VoiceInput @recognized="questionText = $event" />

  <p>识别到的问题：{{ questionText }}</p>
  <button type="button" :disabled="loading" @click="sendToAi">
    {{ loading ? '分析中...' : '发送给 AI 分析' }}
  </button>
  <p v-if="errorMessage">{{ errorMessage }}</p>

  <ChatHistory :messages="messages" />
</template>
