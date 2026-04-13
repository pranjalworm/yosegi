<script lang="ts" setup>
defineProps<{
  collectionPath: string;
  collectionCount: number;
}>();

defineEmits<{
  select: [];
  next: [];
}>();

function dirName(path: string): string {
  const parts = path.split("/");
  return parts[parts.length - 1] || parts[parts.length - 2] || path;
}
</script>

<template>
  <div class="step-body">
    <div class="step-icon">
      <svg
        width="48"
        height="48"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="1.5"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path
          d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"
        />
      </svg>
    </div>
    <h2 class="step-heading">Choose your photo collection</h2>
    <p class="step-description">
      Select a folder of images — these will become the individual tiles that
      make up your mosaic. The more variety you have, the better the color
      matching.
    </p>

    <button class="browse-btn" @click="$emit('select')">
      <svg
        width="18"
        height="18"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path
          d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"
        />
      </svg>
      <span v-if="collectionPath" class="browse-value">{{
        dirName(collectionPath)
      }}</span>
      <span v-else class="browse-placeholder">Browse for a folder...</span>
    </button>

    <p v-if="collectionCount > 0" class="collection-count">
      {{ collectionCount }} images found
    </p>

    <div class="step-actions">
      <div></div>
      <button
        class="btn-primary"
        :disabled="!collectionPath"
        @click="$emit('next')"
      >
        Continue
      </button>
    </div>
  </div>
</template>

<style scoped>
.collection-count {
  font-size: 14px;
  color: var(--accent);
  font-weight: 500;
  margin-top: -4px;
}
</style>
