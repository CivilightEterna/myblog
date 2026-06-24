<template>
  <div class="search-container">
    <div class="search-input-wrap">
      <svg class="search-icon" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>
      </svg>
      <input
        ref="inputRef"
        v-model="query"
        type="search"
        class="search-input"
        placeholder="搜索文章..."
        autocomplete="off"
        @input="doSearch"
      />
      <button v-if="query" class="search-clear" @click="clear" aria-label="清除">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
      </button>
    </div>

    <div v-if="loading" class="search-status">搜索中...</div>

    <div v-if="results !== null && results.length === 0 && !loading && query" class="search-status">
      未找到与 "{{ query }}" 相关的结果
    </div>

    <ul v-if="results && results.length > 0" class="search-results">
      <li v-for="r in results" :key="r.url" class="search-result-item">
        <a :href="r.url" class="search-result-link">
          <h3 class="search-result-title" v-html="r.meta?.title || r.url"></h3>
          <p class="search-result-excerpt" v-html="r.excerpt"></p>
        </a>
      </li>
    </ul>

    <div v-if="query && results === null && !loading" class="search-hint">
      输入关键词搜索文章标题和内容
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from "vue";

const query = ref("");
const results = ref<any[] | null>(null);
const loading = ref(false);
const inputRef = ref<HTMLInputElement>();
const pagefind = ref<any>(null);

let debounceTimer: any = null;

onMounted(async () => {
  // Load Pagefind from the generated index
  try {
    const pf = await import("/pagefind/pagefind.js");
    pf.init();
    pagefind.value = pf;
  } catch (e) {
    console.warn("Pagefind not available (run pnpm build to index)");
  }

  // Focus input
  nextTick(() => inputRef.value?.focus());
});

async function doSearch() {
  if (debounceTimer) clearTimeout(debounceTimer);

  if (!query.value.trim()) {
    results.value = null;
    return;
  }

  debounceTimer = setTimeout(async () => {
    if (!pagefind.value || !query.value.trim()) return;

    loading.value = true;
    try {
      const search = await pagefind.value.search(query.value.trim());
      const items = await Promise.all(search.results.map((r: any) => r.data()));
      results.value = items;
    } catch (e) {
      console.error("Search error:", e);
      results.value = [];
    } finally {
      loading.value = false;
    }
  }, 200);
}

function clear() {
  query.value = "";
  results.value = null;
  nextTick(() => inputRef.value?.focus());
}
</script>

<style scoped>
.search-container {
  max-width: 640px;
  margin: 0 auto;
}

.search-input-wrap {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.875rem 1.25rem;
  background: var(--color-card);
  border: 2px solid var(--color-border);
  border-radius: 999px;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.search-input-wrap:focus-within {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 4px color-mix(in srgb, var(--color-primary) 15%, transparent);
}

.search-icon {
  color: var(--color-text-muted);
  flex-shrink: 0;
}

.search-input {
  flex: 1;
  border: none;
  background: none;
  font-size: 1.0625rem;
  font-family: var(--font-body);
  color: var(--color-text);
  outline: none;
}

.search-input::placeholder {
  color: var(--color-text-muted);
}

.search-clear {
  background: none;
  border: none;
  cursor: pointer;
  color: var(--color-text-muted);
  padding: 4px;
  border-radius: 50%;
  display: flex;
  align-items: center;
}

.search-clear:hover {
  color: var(--color-text);
  background: var(--color-bg-soft);
}

.search-status {
  text-align: center;
  color: var(--color-text-muted);
  padding: 2rem 0;
  font-size: 0.9375rem;
}

.search-hint {
  text-align: center;
  color: var(--color-text-muted);
  padding: 4rem 0;
  font-size: 1rem;
}

.search-results {
  list-style: none;
  padding: 0;
  margin: 2rem 0 0;
}

.search-result-item {
  margin-bottom: 1.5rem;
}

.search-result-link {
  display: block;
  padding: 1.25rem;
  background: var(--color-card);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  text-decoration: none;
  transition: all 0.2s;
}

.search-result-link:hover {
  border-color: var(--color-primary);
  box-shadow: var(--shadow-card);
}

.search-result-title {
  font-size: 1.0625rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0 0 0.375rem;
}

.search-result-link:hover .search-result-title {
  color: var(--color-primary);
}

.search-result-excerpt {
  font-size: 0.875rem;
  color: var(--color-text-muted);
  line-height: 1.6;
  margin: 0;
}

.search-result-excerpt :deep(mark) {
  background: color-mix(in srgb, var(--color-primary) 25%, transparent);
  color: var(--color-text);
  padding: 0.05em 0.2em;
  border-radius: 2px;
}
</style>
