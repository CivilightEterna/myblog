<template>
  <div v-if="provider !== 'none'" class="comment-section">
    <h3 class="comment-title">评论</h3>
    <div v-if="provider === 'twikoo'" id="twikoo-comments"></div>
    <div v-else-if="provider === 'placeholder'" class="comment-placeholder">
      <p>{{ provider }} 评论系统尚未实现，已保留扩展配置结构。</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";

const provider = ref("none");
const twikooEnvId = ref("");
const twikooServerUrl = ref("");

onMounted(async () => {
  // Read comment config from public API (same source as background)
  try {
    const resp = await fetch("/api/public/site-config");
    if (!resp.ok) throw new Error("API unavailable");
    const cfg = await resp.json();

    if (cfg.comment_enabled !== "true") return;
    provider.value = cfg.comment_provider || "none";
    if (provider.value === "none") return;

    twikooEnvId.value = cfg.twikoo_env_id || "";
    twikooServerUrl.value = cfg.twikoo_server_url || "";

    // Load Twikoo
    if (provider.value === "twikoo" && twikooEnvId.value) {
      await loadTwikoo();
    }
  } catch (e) {
    console.log("Comment config not available, comments disabled");
  }
});

async function loadTwikoo() {
  try {
    // Twikoo CDN
    const script = document.createElement("script");
    script.src = "https://cdn.jsdelivr.net/npm/twikoo@1.6.41/dist/twikoo.all.min.js";
    script.onload = () => {
      const twikoo = (window as any).twikoo;
      if (twikoo) {
        const opts: any = {
          envId: twikooEnvId.value,
          el: "#twikoo-comments",
          lang: "zh-CN",
        };
        if (twikooServerUrl.value) {
          opts.serverUrl = twikooServerUrl.value;
        }
        twikoo.init(opts);
      }
    };
    document.body.appendChild(script);
  } catch (e) {
    console.error("Failed to load Twikoo:", e);
  }
}
</script>

<style scoped>
.comment-section {
  margin-top: 3rem;
  padding-top: 2rem;
  border-top: 1px solid var(--color-border);
}

.comment-title {
  font-size: 1.25rem;
  font-weight: 700;
  margin-bottom: 1.5rem;
  color: var(--color-text);
}

.comment-placeholder {
  padding: 2rem;
  text-align: center;
  color: var(--color-text-muted);
  background: var(--color-bg-soft);
  border-radius: var(--radius-md);
  font-size: 0.9375rem;
}
</style>
