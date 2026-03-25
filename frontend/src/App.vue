<script lang="ts" setup>
import { ref, computed, onMounted, watch } from "vue";
import {
  SelectTargetImage,
  SelectCollectionDirectory,
  SelectOutputDirectory,
  CountCollectionImages,
  GetImagePreview,
  GenerateCollage,
  GetDefaultOutputPath,
} from "../wailsjs/go/main/App";
import { EventsOn } from "../wailsjs/runtime/runtime";
import TitleBar from "./components/TitleBar.vue";
import LandingScreen from "./components/LandingScreen.vue";
import FileInput from "./components/FileInput.vue";
import ImagePreview from "./components/ImagePreview.vue";

const screen = ref<"landing" | "editor">("landing");
const theme = ref<"dark" | "light">("light");

const targetPath = ref("");
const collectionPath = ref("");
const outputPath = ref("");
const divisionFactor = ref(100);
const collectionCount = ref(0);
const targetPreview = ref("");
const resultPreview = ref("");
const resultPath = ref("");
const timeTaken = ref("");
const progressMessage = ref("");
const isProcessing = ref(false);
const errorMessage = ref("");

function toggleTheme() {
  theme.value = theme.value === "dark" ? "light" : "dark";
}

watch(
  theme,
  (val) => {
    document.documentElement.setAttribute("data-theme", val);
  },
  { immediate: true },
);

EventsOn("progress", (msg: string) => {
  progressMessage.value = msg;
});

onMounted(async () => {
  const defaultPath = await GetDefaultOutputPath();
  if (defaultPath) {
    outputPath.value = defaultPath;
  }
});

const canGenerate = computed(() => {
  return (
    targetPath.value &&
    collectionPath.value &&
    outputPath.value &&
    divisionFactor.value >= 2 &&
    !isProcessing.value
  );
});

function fileName(path: string): string {
  return path.split("/").pop() ?? path;
}

function dirName(path: string): string {
  const parts = path.split("/");
  return parts[parts.length - 1] || parts[parts.length - 2] || path;
}

async function selectTarget() {
  try {
    const path = await SelectTargetImage();
    if (!path) return;
    targetPath.value = path;
    errorMessage.value = "";
    resultPath.value = "";
    resultPreview.value = "";
    const preview = await GetImagePreview(path);
    targetPreview.value = preview;
  } catch (e: any) {
    errorMessage.value = e;
  }
}

async function selectCollection() {
  try {
    const path = await SelectCollectionDirectory();
    if (!path) return;
    collectionPath.value = path;
    errorMessage.value = "";
    resultPath.value = "";
    resultPreview.value = "";
    const count = await CountCollectionImages(path);
    collectionCount.value = count;
  } catch (e: any) {
    errorMessage.value = e;
  }
}

async function selectOutput() {
  try {
    const path = await SelectOutputDirectory();
    if (!path) return;
    outputPath.value = path;
    errorMessage.value = "";
  } catch (e: any) {
    errorMessage.value = e;
  }
}

async function generate() {
  if (!canGenerate.value) return;
  isProcessing.value = true;
  errorMessage.value = "";
  resultPath.value = "";
  resultPreview.value = "";
  progressMessage.value = "Starting...";

  try {
    const result = await GenerateCollage(
      targetPath.value,
      collectionPath.value,
      outputPath.value,
      divisionFactor.value,
    );
    if (result) {
      resultPath.value = result.outputPath;
      timeTaken.value = result.timeTaken;
      const preview = await GetImagePreview(result.outputPath);
      resultPreview.value = preview;
    }
  } catch (e: any) {
    errorMessage.value = String(e);
  } finally {
    isProcessing.value = false;
    progressMessage.value = "";
  }
}
</script>

<template>
  <TitleBar :theme="theme" @toggle-theme="toggleTheme" />

  <LandingScreen v-if="screen === 'landing'" @start="screen = 'editor'" />

  <main v-else class="content">
    <div class="panel">
      <h1 class="title">Stitch</h1>
      <p class="subtitle">Create photo mosaics from your image collection</p>

      <form class="form" @submit.prevent="generate">
        <FileInput
          label="Target Image"
          :value="targetPath ? fileName(targetPath) : ''"
          placeholder="Select an image..."
          @click="selectTarget"
        />

        <FileInput
          label="Collection Folder"
          :value="
            collectionPath
              ? `${dirName(collectionPath)} (${collectionCount} images)`
              : ''
          "
          placeholder="Select a folder..."
          @click="selectCollection"
        />

        <FileInput
          label="Output Folder"
          :value="outputPath ? dirName(outputPath) : ''"
          placeholder="Select output location..."
          @click="selectOutput"
        />

        <div class="field">
          <label class="field-label">Grid Size</label>
          <div class="grid-input">
            <input
              v-model.number="divisionFactor"
              type="number"
              min="2"
              max="100"
            />
            <span class="grid-hint"
              >{{ divisionFactor }} x {{ divisionFactor }} =
              {{ divisionFactor * divisionFactor }} tiles</span
            >
          </div>
        </div>

        <button type="submit" class="generate-btn" :disabled="!canGenerate">
          <template v-if="isProcessing">
            <span class="spinner"></span>
            {{ progressMessage }}
          </template>
          <template v-else>Generate Mosaic</template>
        </button>

        <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
      </form>
    </div>

    <div class="preview-panel">
      <ImagePreview v-if="resultPreview" :src="resultPreview" label="Result" />
      <ImagePreview
        v-else-if="targetPreview"
        :src="targetPreview"
        label="Target"
      />
      <div v-else class="preview-empty">
        <p>Select a target image to preview</p>
      </div>

      <div v-if="resultPath" class="result-info">
        <p class="result-label">Mosaic saved</p>
        <p class="result-path">{{ resultPath }}</p>
        <p class="result-time">Completed in {{ timeTaken }}</p>
      </div>
    </div>
  </main>
</template>

<style scoped>
.content {
  flex: 1;
  display: flex;
  gap: 1px;
  background: var(--border);
  overflow: hidden;
}

.panel {
  flex: 0 0 360px;
  background: var(--bg-primary);
  padding: 32px 28px;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.title {
  font-size: 22px;
  font-weight: 600;
  letter-spacing: -0.3px;
  margin-bottom: 4px;
}

.subtitle {
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 28px;
}

.form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

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

.grid-input {
  display: flex;
  align-items: center;
  gap: 12px;
}

.grid-input input {
  width: 72px;
  text-align: center;
}

.grid-hint {
  font-size: 12px;
  color: var(--text-secondary);
}

.generate-btn {
  margin-top: 8px;
  background: var(--accent);
  color: #fff;
  font-weight: 500;
  padding: 10px 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.generate-btn:hover:not(:disabled) {
  background: var(--accent-hover);
}

.spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.error {
  font-size: 13px;
  color: var(--error);
}

.preview-panel {
  flex: 1;
  min-width: 0;
  background: var(--bg-secondary);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 24px;
  gap: 16px;
  overflow: hidden;
}

.preview-panel > :deep(.preview) {
  flex: 1;
}

.preview-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.preview-empty p {
  color: var(--text-secondary);
  font-size: 14px;
}

.result-info {
  text-align: center;
  padding: 16px 20px;
  background: var(--bg-elevated);
  border-radius: var(--radius);
  width: 100%;
  max-width: 480px;
}

.result-label {
  font-weight: 500;
  color: var(--success);
  margin-bottom: 4px;
}

.result-path {
  font-size: 12px;
  color: var(--text-secondary);
  word-break: break-all;
  margin-bottom: 4px;
}

.result-time {
  font-size: 12px;
  color: var(--text-secondary);
}
</style>
