/** @type {import('tailwindcss').Config} */
export default {
  content: ["./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}"],
  darkMode: ["selector", '[data-theme="dark"]'],
  theme: {
    extend: {
      colors: {
        bg: "var(--color-bg)",
        "bg-soft": "var(--color-bg-soft)",
        text: "var(--color-text)",
        "text-muted": "var(--color-text-muted)",
        primary: "var(--color-primary)",
        card: "var(--color-card)",
        border: "var(--color-border)",
      },
      fontFamily: {
        body: "var(--font-body)",
        code: "var(--font-code)",
      },
      borderRadius: {
        sm: "var(--radius-sm)",
        md: "var(--radius-md)",
        lg: "var(--radius-lg)",
      },
      boxShadow: {
        card: "var(--shadow-card)",
      },
      typography: {
        DEFAULT: {
          css: {
            maxWidth: "none",
          },
        },
      },
    },
  },
  plugins: [],
};
