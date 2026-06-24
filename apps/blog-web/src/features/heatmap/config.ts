import type { BlogFeature } from "../index";

export const heatmapFeature: BlogFeature = {
  key: "writing-heatmap",
  name: "写作热力图",
  enabled: true,
  position: "home-sidebar",
  component: "Heatmap",
  client: "visible",
  config: {
    type: "publish",       // "publish" = from Markdown dates; "writing" = from DB activity (Phase 4)
    showCount: true,        // show post count in tooltip
    showWords: false,       // word count requires DB-backed writing_activity
    daysToShow: 365,        // how many days of history to display
  },
};
