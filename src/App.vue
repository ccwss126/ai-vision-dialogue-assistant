<script setup>
import { ref } from 'vue'
import CameraView from './components/CameraView.vue'
import VoiceInput from './components/VoiceInput.vue'

const questionText = ref('')
</script>

<template>
  <CameraView />
  <VoiceInput @recognized="questionText = $event" />

  <p>识别到的问题：{{ questionText }}</p>
</template>
