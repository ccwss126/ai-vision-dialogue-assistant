<script setup>
import { onBeforeUnmount, ref } from 'vue'

const emit = defineEmits(['perception-change'])

const videoRef = ref(null)
const stream = ref(null)
const isSensing = ref(false)

async function startSensing() {
  if (stream.value) return

  stream.value = await navigator.mediaDevices.getUserMedia({
    video: true,
    audio: true,
  })
  isSensing.value = true
  emit('perception-change', true)

  if (videoRef.value) {
    videoRef.value.srcObject = stream.value
  }
}

function stopSensing() {
  stream.value?.getTracks().forEach((track) => track.stop())
  stream.value = null
  isSensing.value = false
  emit('perception-change', false)

  if (videoRef.value) {
    videoRef.value.srcObject = null
  }
}

function captureFrame() {
  const video = videoRef.value

  if (!video || !video.videoWidth || !video.videoHeight) {
    return ''
  }

  const maxWidth = 640
  const scale = Math.min(1, maxWidth / video.videoWidth)
  const canvas = document.createElement('canvas')
  canvas.width = Math.round(video.videoWidth * scale)
  canvas.height = Math.round(video.videoHeight * scale)

  const context = canvas.getContext('2d')
  context.drawImage(video, 0, 0, canvas.width, canvas.height)

  return canvas.toDataURL('image/jpeg', 0.82)
}

defineExpose({
  captureFrame,
})

onBeforeUnmount(stopSensing)
</script>

<template>
  <section class="camera-view">
    <div class="video-frame" :class="{ sensing: isSensing }">
      <video ref="videoRef" autoplay playsinline muted></video>
      <div v-if="!isSensing" class="video-placeholder">
        AI 将通过摄像头看画面，并通过麦克风听你说话
      </div>
      <div v-else class="sensing-badge">AI 正在看见画面，并已开启麦克风</div>
    </div>

    <div class="actions">
      <button type="button" class="primary-button" @click="startSensing">开启感知</button>
      <button type="button" class="secondary-button" @click="stopSensing">关闭感知</button>
    </div>
  </section>
</template>

<style scoped>
.camera-view {
  display: grid;
  gap: 16px;
}

.video-frame {
  position: relative;
  overflow: hidden;
  border: 1px solid var(--border);
  border-radius: 8px;
  background:
    linear-gradient(135deg, rgba(20, 184, 166, 0.16), transparent 38%),
    #0f172a;
  aspect-ratio: 16 / 10;
}

.video-frame.sensing {
  border-color: rgba(15, 118, 110, 0.4);
}

video {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
}

.video-placeholder {
  position: absolute;
  inset: 0;
  display: grid;
  place-items: center;
  padding: 24px;
  color: rgba(255, 255, 255, 0.82);
  text-align: center;
}

.sensing-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  padding: 6px 10px;
  border-radius: 999px;
  color: #ffffff;
  background: rgba(15, 118, 110, 0.88);
  font-size: 13px;
  font-weight: 700;
}

.actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}
</style>
