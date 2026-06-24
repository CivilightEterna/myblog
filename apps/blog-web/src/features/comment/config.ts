import type { BlogFeature } from "../index";

export const commentFeature: BlogFeature = {
  key: "comment",
  name: "评论系统",
  enabled: false,
  position: "post-after-content",
  component: "CommentIsland",
  client: "visible",
  config: {
    provider: "giscus",      // giscus, utteranc, twikoo, etc.
    repo: "",
    repoId: "",
    category: "",
    categoryId: "",
  },
};
