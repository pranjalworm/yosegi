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
  OpenFile,
} from "../wailsjs/go/main/App";
import { EventsOn } from "../wailsjs/runtime/runtime";
import LandingScreen from "./components/LandingScreen.vue";
import StepCollection from "./components/StepCollection.vue";
import StepReferenceImage from "./components/StepReferenceImage.vue";
import StepSettings from "./components/StepSettings.vue";
import StepProcessing from "./components/StepProcessing.vue";

type Screen =
  | "landing"
  | "collection"
  | "reference"
  | "settings"
  | "processing";

const screen = ref<Screen>("landing");
const theme = ref<"dark" | "light">("light");

const targetPath = ref("");
const collectionPath = ref("");
const outputPath = ref("");
const divisionFactor = ref(100);
const tileSize = ref(150);
const collectionCount = ref(0);
const targetPreview = ref("");
const resultPreview = ref("");
const resultPath = ref("");
const timeTaken = ref("");
const progressMessage = ref("");
const progressPercent = ref(0);
const isProcessing = ref(false);
const errorMessage = ref("");

const stepLabels = [
  "Photo Collection",
  "Reference Image",
  "Adjust & Create",
  "Mosaic",
];

const stepIndex = computed(() => {
  const map: Record<string, number> = {
    collection: 0,
    reference: 1,
    settings: 2,
    processing: 3,
  };
  return map[screen.value] ?? -1;
});

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

const progressMap: Record<string, number> = {
  "Starting...": 5,
  "Reading reference image...": 15,
  "Analyzing collection images...": 35,
  "Matching tiles to reference...": 60,
  "Building collage...": 80,
  "Encoding output...": 92,
  "Preparing for the grand reveal...": 97,
};

EventsOn("progress", (msg: string) => {
  progressMessage.value = msg;
  if (msg in progressMap) {
    progressPercent.value = progressMap[msg];
  }
});

onMounted(async () => {
  const defaultPath = await GetDefaultOutputPath();
  if (defaultPath) outputPath.value = defaultPath;
});

async function selectCollection() {
  try {
    const path = await SelectCollectionDirectory();
    if (!path) return;
    collectionPath.value = path;
    errorMessage.value = "";
    const count = await CountCollectionImages(path);
    collectionCount.value = count;
  } catch (e: any) {
    errorMessage.value = String(e);
  }
}

async function selectTarget() {
  try {
    const path = await SelectTargetImage();
    if (!path) return;
    targetPath.value = path;
    errorMessage.value = "";
    const preview = await GetImagePreview(path);
    targetPreview.value = preview;
  } catch (e: any) {
    errorMessage.value = String(e);
  }
}

async function selectOutput() {
  try {
    const path = await SelectOutputDirectory();
    if (!path) return;
    outputPath.value = path;
    errorMessage.value = "";
  } catch (e: any) {
    errorMessage.value = String(e);
  }
}

async function generate() {
  screen.value = "processing";
  isProcessing.value = true;
  errorMessage.value = "";
  resultPath.value = "";
  resultPreview.value = "";
  progressMessage.value = "Starting...";
  progressPercent.value = 5;

  try {
    const result = await GenerateCollage(
      targetPath.value,
      collectionPath.value,
      outputPath.value,
      divisionFactor.value,
      tileSize.value,
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
  }
}

function startOver() {
  screen.value = "collection";
  targetPath.value = "";
  targetPreview.value = "";
  resultPreview.value = "";
  resultPath.value = "";
  timeTaken.value = "";
  progressMessage.value = "";
  progressPercent.value = 0;
  errorMessage.value = "";
}
</script>

<template>
  <div class="titlebar-drag" style="--wails-draggable: drag"></div>

  <button
    class="theme-toggle"
    style="--wails-draggable: none"
    @click="toggleTheme"
    :title="theme === 'dark' ? 'Switch to light theme' : 'Switch to dark theme'"
  >
    <svg
      v-if="theme === 'dark'"
      width="16"
      height="16"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
    >
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
    <svg
      v-else
      width="16"
      height="16"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
    >
      <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
    </svg>
  </button>

  <LandingScreen v-if="screen === 'landing'" @start="screen = 'collection'" />

  <main v-else class="wizard">
    <div class="step-indicator">
      <template v-for="(label, i) in stepLabels" :key="i">
        <div
          v-if="i > 0"
          class="indicator-line"
          :class="{ active: stepIndex >= i }"
        />
        <div
          class="indicator-step"
          :class="{ active: stepIndex >= i, current: stepIndex === i }"
        >
          <div class="indicator-number">{{ i + 1 }}</div>
          <span class="indicator-label">{{ label }}</span>
        </div>
      </template>
    </div>

    <div class="step-container">
      <Transition name="step-fade" mode="out-in">
        <StepCollection
          v-if="screen === 'collection'"
          :key="'collection'"
          :collection-path="collectionPath"
          :collection-count="collectionCount"
          @select="selectCollection"
          @next="screen = 'reference'"
        />
        <StepReferenceImage
          v-else-if="screen === 'reference'"
          :key="'reference'"
          :reference-path="targetPath"
          :reference-preview="targetPreview"
          @select="selectTarget"
          @back="screen = 'collection'"
          @next="screen = 'settings'"
        />
        <StepSettings
          v-else-if="screen === 'settings'"
          :key="'settings'"
          :division-factor="divisionFactor"
          :tile-size="tileSize"
          :output-path="outputPath"
          @update:division-factor="divisionFactor = $event"
          @update:tile-size="tileSize = $event"
          @select-output="selectOutput"
          @back="screen = 'reference'"
          @generate="generate"
        />
        <StepProcessing
          v-else-if="screen === 'processing'"
          :key="'processing'"
          :reference-preview="targetPreview"
          :result-preview="resultPreview"
          :progress-message="progressMessage"
          :progress-percent="progressPercent"
          :is-processing="isProcessing"
          :time-taken="timeTaken"
          :result-path="resultPath"
          :error-message="errorMessage"
          @start-over="startOver"
          @open-result="OpenFile(resultPath)"
        />
      </Transition>
    </div>
  </main>
</template>

<style scoped>
.wizard {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
  overflow: hidden;
}

.step-indicator {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px 40px;
}

.theme-toggle {
  position: fixed;
  top: 44px;
  right: 40px;
  z-index: 10;
  width: 30px;
  height: 30px;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-secondary);
  border: 2px solid var(--border);
  color: var(--text-secondary);
  border-radius: 50%;
}

.theme-toggle:hover {
  color: var(--text-primary);
  border-color: var(--accent);
}

.indicator-step {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
}

.indicator-number {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 600;
  border: 2px solid var(--border);
  color: var(--text-secondary);
  background: var(--bg-primary);
  transition: all 0.3s ease;
}

.indicator-step.active .indicator-number {
  border-color: var(--accent);
  background: var(--accent);
  color: #fff;
}

.indicator-label {
  font-size: 12px;
  color: var(--text-secondary);
  transition: color 0.3s ease;
}

.indicator-step.current .indicator-label {
  color: var(--text-primary);
  font-weight: 500;
}

.indicator-line {
  width: 60px;
  height: 2px;
  background: var(--border);
  margin: 0 8px;
  margin-bottom: 20px;
  transition: background 0.3s ease;
}

.indicator-line.active {
  background: var(--accent);
}

.step-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  overflow-y: auto;
  padding: 24px 40px 32px;
}

.step-fade-enter-active,
.step-fade-leave-active {
  transition: opacity 0.2s ease;
}

.step-fade-enter-from,
.step-fade-leave-to {
  opacity: 0;
}

.titlebar-drag {
  height: var(--drag-height);
  flex-shrink: 0;
  background: var(--bg-primary);
}
</style>
