import type { BlogFeature } from "../index";

export const musicFeature: BlogFeature = {
  key: "music",
  name: "音乐播放器",
  enabled: false,
  position: "global-floating",
  component: "MusicIsland",
  client: "idle",
  config: {
    provider: "aplayer",     // aplayer, meting, custom
    server: "netease",
    playlistId: "",
    autoPlay: false,
    defaultVolume: 0.4,
    showOnMobile: false,
  },
};
