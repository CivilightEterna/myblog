<template>
  <div class="build-page">
    <div class="page-header">
      <h1>构建发布</h1>
      <el-button type="primary" :loading="building" @click="handleBuild">
        <el-icon><UploadFilled /></el-icon>
        {{ building ? "构建中..." : "手动构建" }}
      </el-button>
    </div>

    <el-card shadow="never" v-loading="loading">
      <el-table :data="builds" style="width: 100%" empty-text="暂无构建记录">
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag
              :type="row.status === 'success' ? 'success' : row.status === 'running' ? 'warning' : 'danger'"
              size="small"
            >
              {{ row.status === "success" ? "成功" : row.status === "running" ? "运行中" : "失败" }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="trigger_type" label="触发方式" width="100" />
        <el-table-column prop="release_path" label="发布路径" min-width="200" show-overflow-tooltip />
        <el-table-column prop="started_at" label="开始时间" width="170">
          <template #default="{ row }">
            {{ formatDate(row.started_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="finished_at" label="完成时间" width="170">
          <template #default="{ row }">
            {{ row.finished_at ? formatDate(row.finished_at) : "进行中..." }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button text type="primary" size="small" @click="viewLog(row)">
              查看日志
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper" v-if="total > pageSize">
        <el-pagination
          v-model:current-page="page"
          :total="total"
          :page-size="pageSize"
          layout="total, prev, pager, next"
          @current-change="loadBuilds"
        />
      </div>
    </el-card>

    <!-- Log Dialog -->
    <el-dialog v-model="logVisible" title="构建日志" width="700px" top="5vh">
      <pre class="build-log">{{ selectedLog }}</pre>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { getBuilds, triggerBuild } from "@/api/build";
import { UploadFilled } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";

const loading = ref(false);
const building = ref(false);
const builds = ref<any[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = ref(10);

const logVisible = ref(false);
const selectedLog = ref("");

function formatDate(dateStr: string) {
  if (!dateStr) return "—";
  return new Date(dateStr).toLocaleString("zh-CN");
}

async function loadBuilds() {
  loading.value = true;
  try {
    const res = await getBuilds({ page: page.value, page_size: pageSize.value });
    builds.value = res.data.items || [];
    total.value = res.data.total || 0;
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
}

async function handleBuild() {
  building.value = true;
  try {
    await triggerBuild();
    ElMessage.success("构建已触发，请等待完成");
    // Poll for updates
    setTimeout(loadBuilds, 3000);
    setTimeout(loadBuilds, 8000);
    setTimeout(loadBuilds, 15000);
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || "构建触发失败");
  } finally {
    building.value = false;
  }
}

function viewLog(row: any) {
  selectedLog.value = row.log || "暂无日志";
  logVisible.value = true;
}

onMounted(loadBuilds);
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

.pagination-wrapper {
  margin-top: 1rem;
  display: flex;
  justify-content: flex-end;
}

.build-log {
  background: #1f2937;
  color: #e5e7eb;
  padding: 1rem;
  border-radius: 6px;
  font-family: "JetBrains Mono", monospace;
  font-size: 0.8125rem;
  line-height: 1.6;
  max-height: 60vh;
  overflow-y: auto;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
