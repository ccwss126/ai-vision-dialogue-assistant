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
    <p v-if="!messages.length" class="empty-state">还没有对话记录。</p>

    <article v-for="(message, index) in messages" :key="index" class="message">
      <div class="message-meta">第 {{ index + 1 }} 轮对话</div>
      <div class="message-row">
        <span>用户问题</span>
        <p>{{ message.question || '未提供问题' }}</p>
      </div>
      <div class="message-row answer">
        <span>AI 回复</span>
        <p>{{ message.answer }}</p>
      </div>
    </article>
  </section>
</template>

<style scoped>
.chat-history {
  display: grid;
  gap: 14px;
}

.empty-state {
  padding: 24px;
  border: 1px dashed var(--border);
  border-radius: 8px;
  color: var(--muted);
  text-align: center;
  background: var(--surface-soft);
}

.message {
  display: grid;
  gap: 12px;
  padding: 16px;
  border: 1px solid var(--border);
  border-radius: 8px;
  background: var(--surface);
}

.message-meta {
  color: var(--muted);
  font-size: 13px;
}

.message-row {
  display: grid;
  gap: 6px;
}

.message-row span {
  color: var(--muted);
  font-size: 13px;
}

.message-row p {
  color: var(--text);
}

.answer {
  padding-top: 10px;
  border-top: 1px solid var(--border);
}
</style>
