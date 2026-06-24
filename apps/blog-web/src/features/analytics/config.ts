import type { BlogFeature } from "../index";

export const analyticsFeature: BlogFeature = {
  key: "analytics",
  name: "访问统计",
  enabled: false,
  position: "global-footer",
  component: "AnalyticsIsland",
  client: "visible",
  config: {
    provider: "umami",       // umami, plausible, google analytics
    websiteId: "",
    scriptUrl: "",
  },
};
