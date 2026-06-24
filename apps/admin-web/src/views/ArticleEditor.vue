<template>
  <div class="article-editor">
    <div class="page-header">
      <el-button text @click="$router.back()">
        <el-icon><ArrowLeft /></el-icon> 返回
      </el-button>
      <h1>{{ isNew ? "写文章" : "编辑文章" }}</h1>
      <div class="header-actions">
        <el-button @click="handleSave(true)">保存草稿</el-button>
        <el-button type="success" @click="handleSave(false)">保存并发布</el-button>
      </div>
    </div>

    <el-row :gutter="16">
      <!-- Form Panel -->
      <el-col :span="8">
        <el-card shadow="never" class="form-card">
          <el-form :model="form" label-position="top">
            <el-form-item label="标题" required>
              <el-input v-model="form.title" placeholder="文章标题" />
            </el-form-item>

            <el-form-item label="Slug" required>
              <el-input v-model="form.slug" placeholder="url-friendly-slug">
                <template #append>
                  <el-button @click="generateSlug">生成</el-button>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="摘要">
              <el-input v-model="form.description" type="textarea" :rows="2" placeholder="文章摘要（用于SEO和列表展示）" />
            </el-form-item>

            <el-form-item label="分类">
              <el-input v-model="form.category" placeholder="如：后端、前端、数据库" />
            </el-form-item>

            <el-form-item label="标签">
              <el-select
                v-model="form.tags"
                multiple
                filterable
                allow-create
                default-first-option
                placeholder="输入标签后回车"
                style="width: 100%"
              />
            </el-form-item>

            <el-form-item label="封面图 URL">
              <el-input v-model="form.cover" placeholder="如 /uploads/2026/06/image.png" />
            </el-form-item>

            <el-divider />

            <el-form-item label="选项">
              <el-checkbox v-model="form.pinned">置顶</el-checkbox>
              <el-checkbox v-model="form.toc_enabled">显示目录</el-checkbox>
              <el-checkbox v-model="form.comment_enabled">允许评论</el-checkbox>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <!-- Editor Panel -->
      <el-col :span="16">
        <el-card shadow="never" class="editor-card">
          <div class="editor-wrapper">
            <textarea
              v-model="form.content"
              class="markdown-input"
              placeholder="开始写作 Markdown..."
            ></textarea>
            <div class="markdown-preview markdown-body" v-html="renderedPreview"></div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { getArticle, createArticle, updateArticle, publishArticle } from "@/api/article";
import type { ArticleDraft } from "@/api/article";
import { ArrowLeft } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";

const route = useRoute();
const router = useRouter();

const isNew = computed(() => !route.params.id);

const form = reactive<ArticleDraft>({
  title: "",
  slug: "",
  description: "",
  content: "",
  cover: "",
  category: "",
  tags: [],
  draft: true,
  pinned: false,
  comment_enabled: true,
  toc_enabled: true,
});

import { marked } from "marked";

// Configure marked for safe rendering
marked.setOptions({
  breaks: true,       // Support GFM line breaks
  gfm: true,          // GitHub Flavored Markdown
});

const renderedPreview = computed(() => {
  if (!form.content) return "<p style='color:#9ca3af'>在左侧输入 Markdown，这里实时预览...</p>";
  return marked.parse(form.content) as string;
});

function generateSlug() {
  if (!form.title) return;
  form.slug = form.title
    .toLowerCase()
    .trim()
    .replace(/[^\w\s-]/g, "")
    .replace(/[\s_]+/g, "-")
    .replace(/-+/g, "-")
    .replace(/^-+|-+$/g, "");
}

async function handleSave(draft: boolean) {
  if (!form.title || !form.slug) {
    ElMessage.warning("请填写标题和 Slug");
    return;
  }
  if (!form.content) {
    ElMessage.warning("请填写文章内容");
    return;
  }

  form.draft = draft;

  try {
    if (isNew.value) {
      const res = await createArticle({ ...form });
      const id = res.data.id;
      if (!draft) {
        await publishArticle(id);
        ElMessage.success("发布成功！");
      } else {
        ElMessage.success("草稿已保存");
      }
      router.replace(`/articles/${id}/edit`);
    } else {
      const id = Number(route.params.id);
      await updateArticle(id, { ...form });
      if (!draft) {
        await publishArticle(id);
        ElMessage.success("发布成功！");
      } else {
        ElMessage.success("已保存");
      }
    }
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || "保存失败");
  }
}

onMounted(async () => {
  if (!isNew.value) {
    try {
      const id = Number(route.params.id);
      const res = await getArticle(id);
      const data = res.data;
      form.title = data.title;
      form.slug = data.slug;
      form.description = data.description || "";
      form.content = data.content;
      form.cover = data.cover || "";
      form.category = data.category || "";
      form.tags = data.tags_json ? JSON.parse(data.tags_json) : [];
      form.draft = data.draft;
      form.pinned = data.pinned;
      form.toc_enabled = data.toc_enabled;
      form.comment_enabled = data.comment_enabled;
    } catch (e: any) {
      ElMessage.error("加载文章失败");
      router.push("/articles");
    }
  }
});
</script>

<style scoped>
.page-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1rem;
}

.page-header h1 {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1f2937;
  flex: 1;
}

.header-actions {
  display: flex;
  gap: 0.5rem;
}

.form-card {
  margin-bottom: 1rem;
}

.editor-card {
  height: calc(100vh - 180px);
}

.editor-wrapper {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  height: calc(100vh - 240px);
}

.markdown-input {
  width: 100%;
  height: 100%;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  padding: 1rem;
  font-family: "JetBrains Mono", "Fira Code", monospace;
  font-size: 0.875rem;
  line-height: 1.7;
  resize: none;
  outline: none;
}

.markdown-input:focus {
  border-color: #8b5cf6;
}

.markdown-preview {
  padding: 1rem;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  overflow-y: auto;
  font-size: 0.9375rem;
  line-height: 1.8;
  color: #1f2937;
}

.markdown-preview :deep(h1) { font-size: 1.75rem; font-weight: 700; margin: 1.5rem 0 0.75rem; border-bottom: 1px solid #e5e7eb; padding-bottom: 0.5rem; }
.markdown-preview :deep(h2) { font-size: 1.5rem; font-weight: 700; margin: 1.25rem 0 0.625rem; }
.markdown-preview :deep(h3) { font-size: 1.25rem; font-weight: 700; margin: 1rem 0 0.5rem; }
.markdown-preview :deep(h4) { font-size: 1.1rem; font-weight: 700; margin: 0.875rem 0 0.5rem; }
.markdown-preview :deep(p) { margin-bottom: 0.75rem; line-height: 1.8; }
.markdown-preview :deep(strong) { font-weight: 700; }
.markdown-preview :deep(em) { font-style: italic; }
.markdown-preview :deep(a) { color: #8b5cf6; text-decoration: underline; }
.markdown-preview :deep(ul), .markdown-preview :deep(ol) { margin: 0.5rem 0 0.75rem 1.5rem; }
.markdown-preview :deep(li) { margin-bottom: 0.25rem; }
.markdown-preview :deep(blockquote) {
  margin: 0.75rem 0;
  padding: 0.5rem 1rem;
  border-left: 4px solid #8b5cf6;
  background: #f3f4f6;
  color: #6b7280;
}
.markdown-preview :deep(pre) { background: #1f2937; color: #e5e7eb; padding: 1rem; border-radius: 8px; overflow-x: auto; margin: 0.75rem 0; }
.markdown-preview :deep(pre code) { background: none; padding: 0; font-size: 0.8125rem; }
.markdown-preview :deep(code) { background: #f3f4f6; padding: 0.125rem 0.375rem; border-radius: 4px; font-size: 0.875rem; font-family: "JetBrains Mono", "Fira Code", monospace; }
.markdown-preview :deep(table) { width: 100%; border-collapse: collapse; margin: 0.75rem 0; }
.markdown-preview :deep(th) { background: #f3f4f6; font-weight: 600; padding: 0.5rem 0.75rem; border: 1px solid #e5e7eb; text-align: left; }
.markdown-preview :deep(td) { padding: 0.5rem 0.75rem; border: 1px solid #e5e7eb; }
.markdown-preview :deep(hr) { border: none; border-top: 1px solid #e5e7eb; margin: 1.5rem 0; }
.markdown-preview :deep(img) { max-width: 100%; border-radius: 6px; margin: 0.75rem 0; }
.markdown-preview :deep(input[type="checkbox"]) { margin-right: 0.375rem; }

@media (max-width: 1024px) {
  .editor-wrapper {
    grid-template-columns: 1fr;
  }
}
</style>
