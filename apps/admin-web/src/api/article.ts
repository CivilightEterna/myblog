import client from "./client";

export interface ArticleDraft {
  id?: number;
  title: string;
  slug: string;
  description: string;
  content: string;
  cover: string;
  category: string;
  tags: string[];
  draft: boolean;
  pinned: boolean;
  comment_enabled: boolean;
  toc_enabled: boolean;
  word_count?: number;
  published_at?: string;
  updated_at?: string;
  created_at?: string;
}

export function getArticles(params?: Record<string, any>) {
  return client.get("/articles", { params });
}

export function getArticle(id: number) {
  return client.get(`/articles/${id}`);
}

export function createArticle(data: ArticleDraft) {
  return client.post("/articles", data);
}

export function updateArticle(id: number, data: ArticleDraft) {
  return client.put(`/articles/${id}`, data);
}

export function deleteArticle(id: number) {
  return client.delete(`/articles/${id}`);
}

export function publishArticle(id: number) {
  return client.post(`/articles/${id}/publish`);
}

export function unpublishArticle(id: number) {
  return client.post(`/articles/${id}/unpublish`);
}
