<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'

const emit = defineEmits(['update:modelValue'])
const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  isPerceptionActive: {
    type: Boolean,
    default: false,
  },
})

const status = ref('idle')
const recognition = ref(null)
const hasStartedOnce = ref(false)

const isRecording = computed(() => status.value === 'recording')
const inputText = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value),
})

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
    inputText.value = Array.from(event.results)
      .map((result) => result[0].transcript)
      .join('')
  }

  instance.onend = () => {
    status.value = 'idle'
  }

  instance.onerror = () => {
    status.value = 'idle'
  }

  return instance
}

function startListening() {
  if (!props.isPerceptionActive) {
    alert('请先开启感知，让 AI 获取摄像头和麦克风权限。')
    return
  }

  if (props.disabled) return
  if (isRecording.value) return

  if (!recognition.value) {
    recognition.value = getRecognition()
  }

  if (!recognition.value) return

  status.value = 'recording'
  recognition.value.start()
  hasStartedOnce.value = true
}

watch(
  () => props.isPerceptionActive,
  (active) => {
    if (active) {
      startListening()
      return
    }

    recognition.value?.stop()
    status.value = 'idle'
    hasStartedOnce.value = false
  },
)

onBeforeUnmount(() => {
  recognition.value?.stop()
})
</script>

<template>
  <section class="voice-input">
    <div class="composer">
      <textarea
        v-model="inputText"
        rows="1"
        :readonly="isRecording"
        :disabled="disabled || !isPerceptionActive"
        placeholder="开启感知后直接说话，识别后可在这里修改"
      ></textarea>

      <div v-if="isRecording" class="wave" aria-label="录音中">
        <span></span>
        <span></span>
        <span></span>
        <span></span>
      </div>

      <button
        type="button"
        class="mic-button"
        :class="{ active: isRecording, unavailable: !isPerceptionActive }"
        :disabled="disabled || isRecording"
        aria-label="重新识别语音"
        @click="startListening"
      >
        <span>{{ isRecording ? '●' : hasStartedOnce ? '↻' : '🎙' }}</span>
      </button>
    </div>
  </section>
</template>

<style scoped>
.voice-input {
  display: grid;
  gap: 10px;
}

.composer {
  display: grid;
  grid-template-columns: 1fr auto auto auto;
  gap: 8px;
  align-items: center;
  min-height: 64px;
  padding: 8px 8px 8px 18px;
  border: 1px solid var(--border);
  border-radius: 28px;
  background: var(--surface);
  box-shadow: 0 12px 34px rgba(28, 43, 38, 0.1);
}

textarea {
  width: 100%;
  min-height: 24px;
  max-height: 120px;
  padding: 10px 0;
  border: 0;
  resize: none;
  outline: none;
  color: var(--text);
  background: transparent;
  font: inherit;
  line-height: 1.45;
}

textarea::placeholder {
  color: var(--muted);
}

.mic-button {
  width: 52px;
  min-width: 52px;
  height: 52px;
  min-height: 52px;
  padding: 0;
  border-radius: 999px;
  color: #ffffff;
  border-color: transparent;
  background: var(--primary);
  font-size: 24px;
  line-height: 1;
}

.mic-button span {
  display: block;
  transform: translateY(1px);
}

.mic-button.active {
  color: #ffffff;
  background: var(--primary);
}

.mic-button.unavailable {
  opacity: 0.56;
}

.wave {
  display: flex;
  align-items: center;
  gap: 3px;
  height: 24px;
  padding: 0 2px;
}

.wave span {
  width: 3px;
  height: 8px;
  border-radius: 999px;
  background: var(--primary);
  animation: wave 0.9s ease-in-out infinite;
}

.wave span:nth-child(2) {
  animation-delay: 0.12s;
}

.wave span:nth-child(3) {
  animation-delay: 0.24s;
}

.wave span:nth-child(4) {
  animation-delay: 0.36s;
}

@keyframes wave {
  0%,
  100% {
    height: 8px;
  }

  50% {
    height: 22px;
  }
}

@media (max-width: 640px) {
  .composer {
    grid-template-columns: 1fr auto;
    border-radius: 22px;
  }

  .wave {
    grid-column: 1 / -1;
    justify-self: end;
  }
}
</style>
