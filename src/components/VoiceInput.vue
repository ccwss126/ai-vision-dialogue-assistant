<script setup>
import { onBeforeUnmount, ref } from 'vue'

const emit = defineEmits(['recognized'])

const text = ref('')
const recognition = ref(null)

function getRecognition() {
  const SpeechRecognition =
    window.SpeechRecognition || window.webkitSpeechRecognition

  if (!SpeechRecognition) {
    alert('当前浏览器不支持语音识别')
    return null
  }

  const instance = new SpeechRecognition()
  instance.lang = 'zh-CN'
  instance.interimResults = true
  instance.continuous = false

  instance.onresult = (event) => {
    text.value = Array.from(event.results)
      .map((result) => result[0].transcript)
      .join('')
  }

  instance.onend = () => {
    if (text.value.trim()) {
      emit('recognized', text.value.trim())
    }
  }

  return instance
}

function startListening() {
  if (!recognition.value) {
    recognition.value = getRecognition()
  }

  recognition.value?.start()
}

function stopListening() {
  recognition.value?.stop()
}

onBeforeUnmount(() => {
  recognition.value?.stop()
})
</script>

<template>
  <section class="voice-input">
    <div class="actions">
      <button type="button" @click="startListening">开始说话</button>
      <button type="button" @click="stopListening">停止识别</button>
    </div>

    <p>当前识别文本：{{ text }}</p>
  </section>
</template>

<style scoped>
.voice-input {
  display: grid;
  gap: 12px;
  max-width: 720px;
  margin: 24px auto 0;
}

.actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}
</style>
