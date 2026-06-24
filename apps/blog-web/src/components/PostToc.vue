<template>
  <!-- Expanded TOC -->
  <nav
    v-if="tocItems.length > 0 && !collapsed"
    class="toc"
    aria-label="目录"
  >
    <div class="toc-header">
      <h2 class="toc-title">目录</h2>
      <div class="toc-header-right">
        <span class="toc-progress">{{ Math.round(progress) }}%</span>
        <button class="toc-close-btn" title="收起目录" @click="toggle">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="9 18 15 12 9 6" />
          </svg>
        </button>
      </div>
    </div>
    <ul class="toc-list">
      <li
        v-for="(h, i) in tocItems"
        :key="i"
        class="toc-item"
        :class="{
          'toc-item--depth-3': h.depth === 3,
          'toc-item--active': activeIndex === i,
        }"
      >
        <a
          :href="'#' + h.slug"
          class="toc-link"
          :class="{ 'toc-link--active': activeIndex === i }"
          :title="h.text"
          @click.prevent="scrollTo(h.slug)"
        >
          <span class="toc-dot"></span>
          {{ h.text }}
        </a>
      </li>
    </ul>
  </nav>

  <!-- Collapsed: floating toggle button only -->
  <button
    v-if="tocItems.length > 0 && collapsed"
    class="toc-float-btn"
    title="展开目录"
    @click="toggle"
  >
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <line x1="3" y1="5" x2="21" y2="5" />
      <line x1="3" y1="12" x2="21" y2="12" />
      <line x1="3" y1="19" x2="21" y2="19" />
    </svg>
  </button>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from "vue";

export interface Heading {
  depth: number;
  slug: string;
  text: string;
}

const props = defineProps<{ headings: Heading[] }>();

const tocItems = computed(() =>
  props.headings.filter((h) => h.depth >= 2 && h.depth <= 3)
);

const activeIndex = ref(-1);
const progress = ref(0);
const collapsed = ref(false);

function toggle() {
  collapsed.value = !collapsed.value;
  try { localStorage.setItem("post_toc_collapsed", String(collapsed.value)); } catch { /* */ }
}

watch(collapsed, (val) => {
  if (typeof document === "undefined") return;
  if (val) document.documentElement.setAttribute("data-toc-collapsed", "true");
  else document.documentElement.removeAttribute("data-toc-collapsed");
});

let observer: IntersectionObserver | null = null;
let scrollHandler: (() => void) | null = null;

onMounted(() => {
  try { collapsed.value = localStorage.getItem("post_toc_collapsed") === "true"; } catch { /* */ }
  if (collapsed.value) document.documentElement.setAttribute("data-toc-collapsed", "true");

  const headingEls = tocItems.value
    .map((h) => document.getElementById(h.slug))
    .filter(Boolean) as HTMLElement[];

  if (headingEls.length > 0) {
    observer = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) {
          if (entry.isIntersecting) {
            const idx = headingEls.indexOf(entry.target as HTMLElement);
            if (idx >= 0) activeIndex.value = idx;
          }
        }
        if (!entries.some((e) => e.isIntersecting)) {
          let lastAbove = -1;
          for (let i = 0; i < headingEls.length; i++) {
            if (headingEls[i].getBoundingClientRect().top < window.innerHeight * 0.3) lastAbove = i;
          }
          activeIndex.value = lastAbove;
        }
      },
      { rootMargin: "-80px 0px -60% 0px", threshold: 0 }
    );
    headingEls.forEach((el) => observer!.observe(el));
  }

  const updateProgress = () => {
    const h2 = document.documentElement;
    const total = h2.scrollHeight - h2.clientHeight;
    progress.value = total <= 0 ? 0 : Math.min(100, Math.round((h2.scrollTop / total) * 100));
  };
  updateProgress();
  scrollHandler = () => updateProgress();
  window.addEventListener("scroll", scrollHandler, { passive: true });
});

onUnmounted(() => {
  observer?.disconnect();
  if (scrollHandler) window.removeEventListener("scroll", scrollHandler);
});

function scrollTo(slug: string) {
  const el = document.getElementById(slug);
  if (!el) return;
  window.scrollTo({ top: el.getBoundingClientRect().top + window.scrollY - 80, behavior: "smooth" });
}
</script>

<style scoped>
/* ===== Expanded TOC ===== */
.toc {
  background: var(--color-bg-soft);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 1rem 1.25rem;
}

.toc-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.75rem;
}

.toc-header-right {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.toc-title {
  font-size: 0.9375rem;
  font-weight: 700;
  color: var(--color-text);
  margin: 0;
}

.toc-progress {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--color-primary);
  background: color-mix(in srgb, var(--color-primary) 12%, transparent);
  padding: 0.125rem 0.5rem;
  border-radius: 999px;
}

.toc-close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: none;
  color: var(--color-text-muted);
  cursor: pointer;
  transition: all 0.2s;
}
.toc-close-btn:hover {
  color: var(--color-primary);
  border-color: var(--color-primary);
}

.toc-list {
  list-style: none;
  padding: 0; margin: 0;
}

.toc-item { margin-bottom: 0.25rem; }
.toc-item:last-child { margin-bottom: 0; }
.toc-item--depth-3 { padding-left: 1.125rem; }

.toc-link {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.3rem 0.6rem;
  font-size: 0.8125rem;
  color: var(--color-text-muted);
  text-decoration: none;
  border-radius: 4px;
  transition: all 0.2s;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.toc-link:hover {
  color: var(--color-primary);
  background: color-mix(in srgb, var(--color-primary) 6%, transparent);
}
.toc-link--active {
  color: var(--color-primary);
  font-weight: 600;
  background: color-mix(in srgb, var(--color-primary) 10%, transparent);
}

.toc-dot {
  flex-shrink: 0;
  width: 5px; height: 5px;
  border-radius: 50%;
  background: transparent;
  transition: background 0.2s;
}
.toc-link--active .toc-dot { background: var(--color-primary); }

/* ===== Collapsed: floating toggle button ===== */
.toc-float-btn {
  position: sticky;
  top: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: 1px solid var(--color-border);
  border-radius: 50%;
  background: var(--color-bg-soft);
  color: var(--color-text-muted);
  cursor: pointer;
  transition: all 0.2s;
  margin: 0 auto;
}
.toc-float-btn:hover {
  color: var(--color-primary);
  border-color: var(--color-primary);
  background: var(--color-card);
}

/* Mobile: hide */
@media (max-width: 1024px) {
  .toc-float-btn { display: none; }
}
</style>
