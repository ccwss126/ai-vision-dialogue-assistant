<script setup>
defineProps({
  messages: {
    type: Array,
    default: () => [],
  },
})
</script>

<template>
  <section class="chat-history">
    <div v-for="(message, index) in messages" :key="index" class="message">
      <p>用户问题：{{ message.question }}</p>
      <p>AI 回复：{{ message.answer }}</p>
    </div>
  </section>
</template>

<style scoped>
.chat-history {
  display: grid;
  gap: 12px;
  max-width: 720px;
  margin: 24px auto 0;
}

.message {
  padding: 12px;
  border: 1px solid #ddd;
}
</style>
