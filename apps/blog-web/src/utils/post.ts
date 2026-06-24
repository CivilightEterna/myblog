import type { CollectionEntry } from "astro:content";

type Post = CollectionEntry<"posts">;

/**
 * Filter out drafts and sort posts: pinned first, then by date descending.
 */
export function getPublishedPosts(posts: Post[]): Post[] {
  return posts
    .filter((post) => !post.data.draft)
    .sort((a, b) => {
      // Pinned posts come first
      if (a.data.pinned && !b.data.pinned) return -1;
      if (!a.data.pinned && b.data.pinned) return 1;
      // Then by date descending
      return new Date(b.data.date).getTime() - new Date(a.data.date).getTime();
    });
}

/**
 * Get unique categories from posts, with counts.
 */
export function getCategories(posts: Post[]): Map<string, number> {
  const categories = new Map<string, number>();
  for (const post of posts) {
    if (post.data.draft) continue;
    const cat = post.data.category;
    categories.set(cat, (categories.get(cat) || 0) + 1);
  }
  return categories;
}

/**
 * Get unique tags from posts, with counts.
 */
export function getTags(posts: Post[]): Map<string, number> {
  const tags = new Map<string, number>();
  for (const post of posts) {
    if (post.data.draft) continue;
    for (const tag of post.data.tags) {
      tags.set(tag, (tags.get(tag) || 0) + 1);
    }
  }
  // Sort by count descending
  return new Map([...tags.entries()].sort((a, b) => b[1] - a[1]));
}

/**
 * Get previous and next post for navigation.
 */
export function getPrevNext(
  posts: Post[],
  currentSlug: string
): { prev: Post | null; next: Post | null } {
  const published = getPublishedPosts(posts);
  const index = published.findIndex((p) => p.slug === currentSlug);

  return {
    prev: index > 0 ? published[index - 1] : null,
    next: index < published.length - 1 ? published[index + 1] : null,
  };
}

/**
 * Filter posts by category.
 */
export function getPostsByCategory(
  posts: Post[],
  category: string
): Post[] {
  return getPublishedPosts(posts).filter(
    (p) => p.data.category === category
  );
}

/**
 * Filter posts by tag.
 */
export function getPostsByTag(posts: Post[], tag: string): Post[] {
  return getPublishedPosts(posts).filter((p) => p.data.tags.includes(tag));
}
