import { defineCollection, z } from "astro:content";

const posts = defineCollection({
  type: "content",
  schema: z.object({
    contentVersion: z.number().default(1),
    title: z.string(),
    description: z.string(),
    date: z.string(),
    updated: z.string().optional(),
    category: z.string(),
    tags: z.array(z.string()).default([]),
    cover: z.string().optional(),
    draft: z.boolean().default(false),
    pinned: z.boolean().default(false),
    comment: z.boolean().default(true),
    toc: z.boolean().default(true),
  }),
});

export const collections = {
  posts,
};
