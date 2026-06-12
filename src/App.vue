<script setup>
import { ref } from 'vue'
import CameraView from './components/CameraView.vue'
import ChatHistory from './components/ChatHistory.vue'
import VoiceInput from './components/VoiceInput.vue'

const cameraView = ref(null)
const questionText = ref('')
const messages = ref([])

function sendToAi() {
  const question = questionText.value.trim()
  const image = cameraView.value?.captureFrame()

  messages.value.push({
    question,
    image,
    answer: '这是模拟 AI 回复：我已经看到了当前画面。',
  })
}
</script>

<template>
  <CameraView ref="cameraView" />
  <VoiceInput @recognized="questionText = $event" />

  <p>识别到的问题：{{ questionText }}</p>
  <button type="button" @click="sendToAi">发送给 AI 分析</button>

  <ChatHistory :messages="messages" />
</template>
