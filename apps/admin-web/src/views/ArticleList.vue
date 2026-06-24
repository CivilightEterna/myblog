<template>
  <div class="article-list-page">
    <div class="page-header">
      <h1>文章管理</h1>
      <el-button type="primary" @click="$router.push('/articles/new')">
        <el-icon><Plus /></el-icon> 写文章
      </el-button>
    </div>

    <el-card shadow="never">
      <div class="search-bar">
        <el-input
          v-model="search"
          placeholder="搜索文章标题..."
          clearable
          style="width: 300px"
          @input="onSearch"
        />
        <el-select v-model="filterDraft" placeholder="状态" clearable style="width: 140px" @change="loadArticles">
          <el-option label="全部" value="" />
          <el-option label="草稿" value="true" />
          <el-option label="已发布" value="false" />
        </el-select>
      </div>

      <el-table :data="articles" v-loading="loading" style="width: 100%">
        <el-table-column prop="title" label="标题" min-width="200">
          <template #default="{ row }">
            <el-link type="primary" @click="$router.push(`/articles/${row.id}/edit`)">
              {{ row.title }}
              <el-tag v-if="row.pinned" size="small" type="warning" style="margin-left: 4px">置顶</el-tag>
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="category" label="分类" width="100" />
        <el-table-column prop="slug" label="Slug" width="140" show-overflow-tooltip />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.draft ? 'info' : 'success'" size="small">
              {{ row.draft ? "草稿" : "已发布" }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="word_count" label="字数" width="70" />
        <el-table-column prop="updated_at" label="更新时间" width="170">
          <template #default="{ row }">
            {{ formatDate(row.updated_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" size="small" @click="$router.push(`/articles/${row.id}/edit`)">
              编辑
            </el-button>
            <el-button
              v-if="row.draft"
              text
              type="success"
              size="small"
              :loading="publishingId === row.id"
              @click="handlePublish(row)"
            >
              发布
            </el-button>
            <el-button
              v-else
              text
              type="warning"
              size="small"
              @click="handleUnpublish(row)"
            >
              取消发布
            </el-button>
            <el-popconfirm title="确认删除？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button text type="danger" size="small">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="page"
          :total="total"
          :page-size="pageSize"
          layout="total, prev, pager, next"
          @current-change="loadArticles"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { getArticles, deleteArticle, publishArticle, unpublishArticle } from "@/api/article";
import { Plus } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";

const loading = ref(false);
const articles = ref<any[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = ref(20);
const search = ref("");
const filterDraft = ref("");
const publishingId = ref<number | null>(null);

let searchTimer: any = null;

function formatDate(dateStr: string) {
  if (!dateStr) return "—";
  return new Date(dateStr).toLocaleString("zh-CN");
}

function onSearch() {
  clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    page.value = 1;
    loadArticles();
  }, 300);
}

async function loadArticles() {
  loading.value = true;
  try {
    const params: any = { page: page.value, page_size: pageSize.value };
    if (search.value) params.search = search.value;
    if (filterDraft.value) params.draft = filterDraft.value;

    const res = await getArticles(params);
    articles.value = res.data.items || [];
    total.value = res.data.total || 0;
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
}

async function handlePublish(row: any) {
  publishingId.value = row.id;
  try {
    await publishArticle(row.id);
    ElMessage.success("发布成功");
    loadArticles();
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || "发布失败");
  } finally {
    publishingId.value = null;
  }
}

async function handleUnpublish(row: any) {
  try {
    await unpublishArticle(row.id);
    ElMessage.success("已取消发布");
    loadArticles();
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || "操作失败");
  }
}

async function handleDelete(row: any) {
  try {
    await deleteArticle(row.id);
    ElMessage.success("已删除");
    loadArticles();
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || "删除失败");
  }
}

onMounted(loadArticles);
</script>

<style scoped>
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1rem;
}

.page-header h1 {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1f2937;
}

.search-bar {
  display: flex;
  gap: 0.75rem;
  margin-bottom: 1rem;
}

.pagination-wrapper {
  margin-top: 1rem;
  display: flex;
  justify-content: flex-end;
}
</style>
