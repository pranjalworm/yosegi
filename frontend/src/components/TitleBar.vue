<script lang="ts" setup>
defineProps<{
  theme: "dark" | "light";
}>();

defineEmits<{
  toggleTheme: [];
}>();
</script>

<template>
  <div class="titlebar" style="--wails-draggable: drag">
    <button
      class="theme-btn"
      style="--wails-draggable: none"
      @click="$emit('toggleTheme')"
      :title="theme === 'dark' ? 'Switch to light theme' : 'Switch to dark theme'"
    >
      <!-- sun icon for dark mode (click to go light) -->
      <svg v-if="theme === 'dark'" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="12" cy="12" r="5" />
        <line x1="12" y1="1" x2="12" y2="3" />
        <line x1="12" y1="21" x2="12" y2="23" />
        <line x1="4.22" y1="4.22" x2="5.64" y2="5.64" />
        <line x1="18.36" y1="18.36" x2="19.78" y2="19.78" />
        <line x1="1" y1="12" x2="3" y2="12" />
        <line x1="21" y1="12" x2="23" y2="12" />
        <line x1="4.22" y1="19.78" x2="5.64" y2="18.36" />
        <line x1="18.36" y1="5.64" x2="19.78" y2="4.22" />
      </svg>
      <!-- moon icon for light mode (click to go dark) -->
      <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
      </svg>
    </button>
  </div>
</template>

<style scoped>
.titlebar {
  height: var(--drag-height);
  flex-shrink: 0;
  background: var(--bg-primary);
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding-right: 12px;
}

.theme-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  padding: 0;
  background: transparent;
  color: var(--text-secondary);
  border-radius: 4px;
}

.theme-btn:hover {
  color: var(--text-primary);
  background: var(--bg-elevated);
}
</style>
