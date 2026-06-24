<template>
  <div class="reading-progress-bar" :style="{ width: progress + '%' }"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";

const progress = ref(0);

let ticking = false;

function updateProgress() {
  const h = document.documentElement;
  const total = h.scrollHeight - h.clientHeight;
  if (total <= 0) {
    progress.value = 0;
    return;
  }
  progress.value = Math.min(100, (h.scrollTop / total) * 100);
}

function onScroll() {
  if (!ticking) {
    requestAnimationFrame(() => {
      updateProgress();
      ticking = false;
    });
    ticking = true;
  }
}

onMounted(() => {
  updateProgress();
  window.addEventListener("scroll", onScroll, { passive: true });
});

onUnmounted(() => {
  window.removeEventListener("scroll", onScroll);
});
</script>

<style scoped>
.reading-progress-bar {
  position: fixed;
  top: 0;
  left: 0;
  height: 3px;
  background: var(--color-primary);
  z-index: 200;
  transition: width 0.15s linear;
  border-radius: 0 2px 2px 0;
}
</style>
