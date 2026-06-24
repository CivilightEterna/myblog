<template>
  <div class="dashboard">
    <h1 class="page-title">仪表盘</h1>

    <!-- Stats Cards -->
    <el-row :gutter="16" class="stats-row">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-value">{{ stats.totalArticles }}</div>
          <div class="stat-label">文章总数</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-value">{{ stats.draftCount }}</div>
          <div class="stat-label">草稿</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-value">{{ stats.totalWords }}</div>
          <div class="stat-label">总字数</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-value">{{ stats.lastBuild }}</div>
          <div class="stat-label">最近构建</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Recent Posts -->
    <el-card class="section-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span>最近文章</span>
          <el-button type="primary" size="small" @click="$router.push('/articles/new')">
            写文章
          </el-button>
        </div>
      </template>

      <el-table :data="recentPosts" style="width: 100%" v-loading="loading" empty-text="暂无文章">
        <el-table-column prop="title" label="标题" min-width="200">
          <template #default="{ row }">
            <el-link type="primary" @click="$router.push(`/articles/${row.id}/edit`)">
              {{ row.title }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="category" label="分类" width="120" />
        <el-table-column prop="word_count" label="字数" width="80" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.draft ? 'info' : 'success'" size="small">
              {{ row.draft ? "草稿" : "已发布" }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="更新时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.updated_at) }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { getArticles } from "@/api/article";
import { getBuilds } from "@/api/build";

const loading = ref(false);
const recentPosts = ref<any[]>([]);
const stats = ref({
  totalArticles: 0,
  draftCount: 0,
  totalWords: 0,
  lastBuild: "—",
});

function formatDate(dateStr: string) {
  if (!dateStr) return "—";
  return new Date(dateStr).toLocaleString("zh-CN");
}

onMounted(async () => {
  loading.value = true;
  try {
    const [articlesRes, buildsRes] = await Promise.all([
      getArticles({ page_size: 5 }),
      getBuilds({ page_size: 1 }),
    ]);

    recentPosts.value = articlesRes.data.items || [];
    const allArticles = articlesRes.data.items || [];

    // Get total counts with a separate query for all articles
    const allRes = await getArticles({ page_size: 1000 });
    const items = allRes.data.items || [];
    stats.value.totalArticles = allRes.data.total || items.length;
    stats.value.draftCount = items.filter((a: any) => a.draft).length;
    stats.value.totalWords = items.reduce((sum: number, a: any) => sum + (a.word_count || 0), 0);

    const builds = buildsRes.data.items || [];
    if (builds.length > 0) {
      stats.value.lastBuild = builds[0].status === "success" ? "成功" : "失败";
    }
  } catch (e) {
    console.error("Failed to load dashboard:", e);
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped>
.page-title {
  font-size: 1.5rem;
  font-weight: 700;
  margin-bottom: 1.5rem;
  color: #1f2937;
}

.stats-row {
  margin-bottom: 1.5rem;
}

.stat-card {
  text-align: center;
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: #8b5cf6;
}

.stat-label {
  font-size: 0.875rem;
  color: #9ca3af;
  margin-top: 0.25rem;
}

.section-card {
  margin-bottom: 1.5rem;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>
