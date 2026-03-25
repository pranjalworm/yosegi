<script lang="ts" setup>
defineProps<{
  label: string;
  value: string;
  placeholder: string;
}>();

defineEmits<{
  click: [];
}>();
</script>

<template>
  <div class="field">
    <label class="field-label">{{ label }}</label>
    <button class="file-btn" @click="$emit('click')">
      <span v-if="value" class="file-value">{{ value }}</span>
      <span v-else class="file-placeholder">{{ placeholder }}</span>
      <span class="file-icon">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z" />
        </svg>
      </span>
    </button>
  </div>
</template>

<style scoped>
.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field-label {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.file-btn {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: var(--bg-elevated);
  border: 1px solid var(--border);
  color: var(--text-primary);
  padding: 8px 12px;
  text-align: left;
}

.file-btn:hover {
  border-color: var(--accent);
}

.file-value {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  font-size: 13px;
}

.file-placeholder {
  color: var(--text-secondary);
  font-size: 13px;
}

.file-icon {
  flex-shrink: 0;
  color: var(--text-secondary);
  margin-left: 8px;
  display: flex;
}
</style>
