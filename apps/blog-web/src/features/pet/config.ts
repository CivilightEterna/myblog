import type { BlogFeature } from "../index";

export const petFeature: BlogFeature = {
  key: "live2d-pet",
  name: "Live2D 桌宠",
  enabled: false,                       // Disabled by default — enable when implemented
  position: "global-floating",
  component: "PetIsland",
  client: "idle",                       // Load when browser is idle — don't block page
  config: {
    model: "/pet/models/hiyori/model3.json",
    width: 280,
    height: 380,
    showOnMobile: false,
  },
};
