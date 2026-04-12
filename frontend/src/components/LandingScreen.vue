<script lang="ts" setup>
import appIcon from "../assets/appicon.png";

defineEmits<{
  start: [];
}>();
</script>

<template>
  <div class="landing">
    <div class="hero">
      <div class="icon">
        <img :src="appIcon" alt="Yosegi" width="100" height="100" />
      </div>

      <h1 class="title">Yosegi</h1>
      <p class="tagline">
        Turn any photo into a mosaic made from your own images
      </p>
      <div class="steps">
        <div class="step">
          <span class="step-number">1</span>
          <div>
            <p class="step-title">Pick a target image</p>
            <p class="step-desc">The photo you want to recreate as a mosaic</p>
          </div>
        </div>
        <div class="step">
          <span class="step-number">2</span>
          <div>
            <p class="step-title">Choose a collection</p>
            <p class="step-desc">
              A folder of images that will become the tiles
            </p>
          </div>
        </div>
        <div class="step">
          <span class="step-number">3</span>
          <div>
            <p class="step-title">Generate</p>
            <p class="step-desc">
              Each tile is matched by color to recreate your target
            </p>
          </div>
        </div>
      </div>

      <button class="start-btn" @click="$emit('start')">Let's Go</button>
      <p class="privacy-note">
        Everything runs locally on your device — your photos never leave your
        computer.
      </p>
    </div>
  </div>
</template>

<style scoped>
.landing {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-primary);
  padding: 40px;
}

.hero {
  display: flex;
  flex-direction: column;
  align-items: center;
  max-width: 420px;
  text-align: center;
}

.icon {
  color: var(--accent);
  margin-bottom: 20px;
}

.title {
  font-size: 36px;
  font-weight: 700;
  letter-spacing: -0.5px;
  margin-bottom: 8px;
}

.tagline {
  font-size: 15px;
  color: var(--text-secondary);
  line-height: 1.6;
  margin-bottom: 36px;
}

.privacy-note {
  font-size: 13px;
  color: var(--text-secondary);
  opacity: 0.75;
  margin-top: 24px;
}

.steps {
  display: flex;
  flex-direction: column;
  gap: 16px;
  width: 100%;
  margin-bottom: 36px;
}

.step {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  text-align: left;
  padding: 14px 16px;
  background: var(--bg-secondary);
  border-radius: var(--radius);
}

.step-number {
  flex-shrink: 0;
  width: 26px;
  height: 26px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--accent);
  color: #fff;
  border-radius: 50%;
  font-size: 13px;
  font-weight: 600;
}

.step-title {
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 2px;
}

.step-desc {
  font-size: 13px;
  color: var(--text-secondary);
}

.start-btn {
  background: var(--accent);
  color: #fff;
  font-weight: 500;
  font-size: 15px;
  padding: 12px 40px;
}

.start-btn:hover {
  background: var(--accent-hover);
}
</style>
