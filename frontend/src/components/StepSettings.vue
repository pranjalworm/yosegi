<script lang="ts" setup>
defineProps<{
  divisionFactor: number;
  tileSize: number;
  outputPath: string;
}>();

const emit = defineEmits<{
  "update:divisionFactor": [value: number];
  "update:tileSize": [value: number];
  selectOutput: [];
  back: [];
  generate: [];
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
        <line x1="4" y1="21" x2="4" y2="14" />
        <line x1="4" y1="10" x2="4" y2="3" />
        <line x1="12" y1="21" x2="12" y2="12" />
        <line x1="12" y1="8" x2="12" y2="3" />
        <line x1="20" y1="21" x2="20" y2="16" />
        <line x1="20" y1="12" x2="20" y2="3" />
        <line x1="1" y1="14" x2="7" y2="14" />
        <line x1="9" y1="8" x2="15" y2="8" />
        <line x1="17" y1="16" x2="23" y2="16" />
      </svg>
    </div>
    <h2 class="step-heading">Configure your mosaic</h2>
    <p class="step-description">Fine-tune how your final mosaic will look.</p>

    <div class="settings-list">
      <div class="setting">
        <div class="setting-header">
          <label class="setting-label">Grid Size</label>
          <span class="setting-value">
            {{ divisionFactor }} &times; {{ divisionFactor }} =
            {{ (divisionFactor * divisionFactor).toLocaleString() }} tiles
          </span>
        </div>
        <input
          type="range"
          :value="divisionFactor"
          @input="
            emit(
              'update:divisionFactor',
              +($event.target as HTMLInputElement).value,
            )
          "
          min="5"
          max="100"
          step="5"
        />
        <p class="setting-hint">
          Controls how many tiles make up the mosaic. Higher values create more
          detail but require more collection images for accurate color matching.
        </p>
      </div>

      <div class="setting">
        <div class="setting-header">
          <label class="setting-label">Tile Quality</label>
          <span class="setting-value">{{ tileSize }}px per tile</span>
        </div>
        <input
          type="range"
          :value="tileSize"
          @input="
            emit('update:tileSize', +($event.target as HTMLInputElement).value)
          "
          min="30"
          max="300"
          step="10"
        />
        <p class="setting-hint">
          Sets the resolution of each individual tile. Higher values produce a
          sharper final image but increase file size and processing time.
        </p>
      </div>

      <div class="setting">
        <div class="setting-header">
          <label class="setting-label">Save to</label>
        </div>
        <button class="output-btn" @click="emit('selectOutput')">
          <svg
            width="14"
            height="14"
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
          <span>{{ dirName(outputPath) }}</span>
        </button>
      </div>
    </div>

    <div class="step-actions">
      <button class="btn-secondary" @click="emit('back')">Back</button>
      <button class="btn-primary" @click="emit('generate')">
        Create Mosaic
      </button>
    </div>
  </div>
</template>

<style scoped>
.settings-list {
  display: flex;
  flex-direction: column;
  gap: 28px;
  width: 100%;
}

.setting {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.setting-header {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
}

.setting-label {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.setting-value {
  font-size: 13px;
  color: var(--text-secondary);
  font-variant-numeric: tabular-nums;
}

.setting-hint {
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.55;
}

.output-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  background: var(--bg-elevated);
  border: 1px solid var(--border);
  color: var(--text-primary);
  padding: 10px 14px;
  font-size: 14px;
  text-align: left;
}

.output-btn:hover {
  border-color: var(--accent);
}
</style>
