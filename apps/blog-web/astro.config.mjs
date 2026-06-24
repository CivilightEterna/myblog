import { defineConfig } from "astro/config";
import tailwindcss from "@tailwindcss/vite";
import sitemap from "@astrojs/sitemap";
import vue from "@astrojs/vue";
import { fileURLToPath } from "node:url";

// https://astro.build/config
export default defineConfig({
  site: "https://your-domain.com",
  integrations: [
    sitemap(),
    vue(),
  ],
  markdown: {
    shikiConfig: {
      theme: "github-dark",
      wrap: true,
    },
  },
  vite: {
    plugins: [tailwindcss()],
    build: {
      rollupOptions: {
        external: ["/pagefind/pagefind.js"],
      },
    },
    server: {
      proxy: {
        "/api": "http://localhost:8080",
        "/uploads": "http://localhost:8080",
      },
    },
    resolve: {
      alias: {
        "@": fileURLToPath(new URL("./src", import.meta.url)),
        "@components": fileURLToPath(new URL("./src/components", import.meta.url)),
        "@layouts": fileURLToPath(new URL("./src/layouts", import.meta.url)),
        "@utils": fileURLToPath(new URL("./src/utils", import.meta.url)),
        "@config": fileURLToPath(new URL("./src/config", import.meta.url)),
        "@features": fileURLToPath(new URL("./src/features", import.meta.url)),
      },
    },
  },
});
