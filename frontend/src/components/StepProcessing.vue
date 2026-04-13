<script lang="ts" setup>
defineProps<{
  referencePreview: string;
  resultPreview: string;
  progressMessage: string;
  progressPercent: number;
  isProcessing: boolean;
  timeTaken: string;
  resultPath: string;
  errorMessage: string;
}>();

defineEmits<{
  startOver: [];
  openResult: [];
}>();
</script>

<template>
  <div class="step-body step-processing">
    <div class="image-display">
      <Transition name="fade" mode="out-in">
        <img
          v-if="resultPreview"
          :key="'result'"
          :src="resultPreview"
          alt="Generated mosaic"
          class="display-img"
        />
        <img
          v-else
          :key="'reference'"
          :src="referencePreview"
          alt="Reference image"
          class="display-img"
        />
      </Transition>
    </div>

    <div v-if="isProcessing" class="progress-section">
      <div class="progress-bar">
        <div
          class="progress-fill"
          :style="{ width: progressPercent + '%' }"
        ></div>
      </div>
      <p class="progress-text">{{ progressMessage }}</p>
    </div>

    <div v-if="errorMessage" class="error-section">
      <p class="step-error">{{ errorMessage }}</p>
      <button class="btn-primary" @click="$emit('startOver')">
        Start Over
      </button>
    </div>

    <div
      v-if="!isProcessing && resultPreview && !errorMessage"
      class="result-section"
    >
      <p class="time-label">Created in {{ timeTaken }}</p>
      <div class="result-actions">
        <button class="btn-primary" @click="$emit('openResult')">
          Open Image
        </button>
        <button class="btn-primary" @click="$emit('startOver')">
          Create Another
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.step-processing {
  max-width: 640px !important;
  flex: 1;
  min-height: 0;
  margin-top: 0;
  margin-bottom: 0;
  padding-bottom: 0;
}

.image-display {
  width: 100%;
  flex: 1;
  min-height: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius);
  overflow: hidden;
}

.display-img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.progress-section {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 10px;
  align-items: center;
}

.progress-bar {
  width: 100%;
  height: 6px;
  background: var(--bg-elevated);
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: var(--accent);
  border-radius: 3px;
  transition: width 0.4s ease;
}

.progress-text {
  font-size: 14px;
  color: var(--text-secondary);
}

.result-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.time-label {
  font-size: 14px;
  color: var(--text-secondary);
}

.result-actions {
  display: flex;
  gap: 12px;
  margin-top: 10px;
}

.error-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}
</style>
