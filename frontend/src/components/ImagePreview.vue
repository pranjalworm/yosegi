<script lang="ts" setup>
defineProps<{
  src: string;
  label: string;
}>();
</script>

<template>
  <div class="preview">
    <span class="preview-label">{{ label }}</span>
    <img :src="src" :alt="label" class="preview-img" />
  </div>
</template>

<style scoped>
.preview {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  width: 100%;
  height: 100%;
  min-height: 0;
  overflow: hidden;
}

.preview-label {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  flex-shrink: 0;
}

.preview-img {
  width: 100%;
  height: 100%;
  min-height: 0;
  border-radius: var(--radius);
  object-fit: contain;
}
</style>
