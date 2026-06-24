// ============================================
// Theme Presets — 4 built-in themes
// ============================================
export type ThemeMode = "light" | "dark";
export type HomeLayout = "classic" | "minimal" | "magazine" | "terminal";
export type PostLayout = "sidebar" | "clean" | "wide";
export type RadiusSize = "small" | "medium" | "large";

export interface ThemePreset {
  label: string;
  mode: ThemeMode;
  primaryColor: string;
  homeLayout: HomeLayout;
  postLayout: PostLayout;
  radius: RadiusSize;
  fontFamily: "sans" | "serif" | "mono";
}

export const themePresets: Record<string, ThemePreset> = {
  default: {
    label: "默认",
    mode: "light",
    primaryColor: "#3b82f6",
    homeLayout: "classic",
    postLayout: "sidebar",
    radius: "large",
    fontFamily: "sans",
  },
  minimal: {
    label: "极简",
    mode: "light",
    primaryColor: "#111827",
    homeLayout: "minimal",
    postLayout: "clean",
    radius: "small",
    fontFamily: "sans",
  },
  dark: {
    label: "暗黑",
    mode: "dark",
    primaryColor: "#8b5cf6",
    homeLayout: "classic",
    postLayout: "sidebar",
    radius: "large",
    fontFamily: "sans",
  },
  terminal: {
    label: "终端",
    mode: "dark",
    primaryColor: "#22c55e",
    homeLayout: "terminal",
    postLayout: "clean",
    radius: "small",
    fontFamily: "mono",
  },
};

/** Get active theme from localStorage or default */
export function getActiveTheme(): string {
  if (typeof window === "undefined") return "default";
  return localStorage.getItem("activeTheme") || "default";
}
