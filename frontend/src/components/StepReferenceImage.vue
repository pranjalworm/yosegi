<script lang="ts" setup>
defineProps<{
  referencePath: string;
  referencePreview: string;
}>();

defineEmits<{
  select: [];
  back: [];
  next: [];
}>();

function fileName(path: string): string {
  return path.split("/").pop() ?? path;
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
        <rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
        <circle cx="8.5" cy="8.5" r="1.5" />
        <polyline points="21 15 16 10 5 21" />
      </svg>
    </div>
    <h2 class="step-heading">Choose your reference image</h2>
    <p class="step-description">
      This is the photo your mosaic will recreate — each section will be
      color-matched to a tile from your collection.
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
        <rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
        <circle cx="8.5" cy="8.5" r="1.5" />
        <polyline points="21 15 16 10 5 21" />
      </svg>
      <span v-if="referencePath" class="browse-value">{{
        fileName(referencePath)
      }}</span>
      <span v-else class="browse-placeholder">Browse for an image...</span>
    </button>

    <Transition name="preview-fade">
      <div v-if="referencePreview" class="image-preview">
        <img :src="referencePreview" alt="Reference image" />
      </div>
    </Transition>

    <div class="step-actions">
      <button class="btn-secondary" @click="$emit('back')">Back</button>
      <button
        class="btn-primary"
        :disabled="!referencePath"
        @click="$emit('next')"
      >
        Continue
      </button>
    </div>
  </div>
</template>

<style scoped>
.image-preview {
  width: 100%;
  max-height: 260px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.image-preview img {
  max-width: 100%;
  max-height: 260px;
  object-fit: contain;
  border-radius: var(--radius);
}

.preview-fade-enter-active {
  transition: opacity 0.4s ease, transform 0.4s ease;
}

.preview-fade-enter-from {
  opacity: 0;
  transform: scale(0.95);
}
</style>
