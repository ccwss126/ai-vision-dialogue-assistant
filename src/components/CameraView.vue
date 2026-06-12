<script setup>
import { onBeforeUnmount, ref } from 'vue'

const videoRef = ref(null)
const stream = ref(null)

async function startCamera() {
  if (stream.value) return

  stream.value = await navigator.mediaDevices.getUserMedia({ video: true })

  if (videoRef.value) {
    videoRef.value.srcObject = stream.value
  }
}

function stopCamera() {
  stream.value?.getTracks().forEach((track) => track.stop())
  stream.value = null

  if (videoRef.value) {
    videoRef.value.srcObject = null
  }
}

onBeforeUnmount(stopCamera)
</script>

<template>
  <section class="camera-view">
    <video ref="videoRef" autoplay playsinline></video>

    <div class="actions">
      <button type="button" @click="startCamera">开启摄像头</button>
      <button type="button" @click="stopCamera">关闭摄像头</button>
    </div>
  </section>
</template>

<style scoped>
.camera-view {
  display: grid;
  gap: 16px;
  max-width: 720px;
  margin: 0 auto;
}

video {
  width: 100%;
  min-height: 360px;
  background: #111;
}

.actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}
</style>
