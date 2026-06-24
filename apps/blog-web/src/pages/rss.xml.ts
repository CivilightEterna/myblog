import rss from "@astrojs/rss";
import { getCollection } from "astro:content";
import { getPublishedPosts } from "@utils/post";
import { siteConfig } from "@config/site.config";

export async function GET(context: { site: URL }) {
  const posts = await getCollection("posts");
  const published = getPublishedPosts(posts);

  return rss({
    title: siteConfig.title,
    description: siteConfig.description,
    site: context.site,
    items: published.map((post) => ({
      title: post.data.title,
      description: post.data.description,
      pubDate: new Date(post.data.date),
      link: `/posts/${post.slug}`,
      categories: [post.data.category, ...post.data.tags],
    })),
    customData: `<language>${siteConfig.lang}</language>`,
  });
}
