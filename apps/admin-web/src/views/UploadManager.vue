<template>
  <div class="upload-page">
    <h1 class="page-title">文件上传</h1>

    <el-card shadow="never" class="upload-card">
      <el-upload
        class="upload-area"
        drag
        multiple
        :auto-upload="false"
        :on-change="handleFileChange"
        :show-file-list="false"
        accept=".jpg,.jpeg,.png,.gif,.webp,.svg"
      >
        <el-icon class="upload-icon"><UploadFilled /></el-icon>
        <div class="upload-text">将文件拖到此处，或点击上传</div>
        <div class="upload-hint">支持 jpg / png / gif / webp / svg，单文件不超过 5MB</div>
      </el-upload>

      <div v-if="uploadingFiles.length" class="uploading-list">
        <div v-for="f in uploadingFiles" :key="f.name" class="uploading-item">
          <span>{{ f.name }}</span>
          <el-progress :percentage="f.progress" :status="f.status" />
        </div>
      </div>
    </el-card>

    <el-card shadow="never" class="files-card">
      <template #header>
        <span>已上传文件（{{ total }}）</span>
      </template>

      <div v-if="files.length === 0 && !loading" class="empty">暂无文件</div>

      <div class="file-grid">
        <div v-for="file in files" :key="file.id" class="file-item">
          <div class="file-preview">
            <img v-if="isImage(file.mime_type)" :src="file.url" :alt="file.original_name" loading="lazy" />
            <div v-else class="file-icon">
              <el-icon size="48"><Document /></el-icon>
            </div>
          </div>
          <div class="file-info">
            <span class="file-name" :title="file.original_name">{{ file.original_name }}</span>
            <span class="file-size">{{ formatSize(file.size) }}</span>
          </div>
          <div class="file-actions">
            <el-button size="small" @click="copyUrl(file.url)">复制URL</el-button>
            <el-popconfirm title="确认删除？" @confirm="handleDelete(file.id)">
              <template #reference>
                <el-button size="small" type="danger" text>删除</el-button>
              </template>
            </el-popconfirm>
          </div>
        </div>
      </div>

      <div class="pagination-wrapper" v-if="total > pageSize">
        <el-pagination
          v-model:current-page="page"
          :total="total"
          :page-size="pageSize"
          layout="total, prev, pager, next"
          @current-change="loadFiles"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { getUploads, uploadFile, deleteUpload } from "@/api/upload";
import { UploadFilled, Document } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";

const loading = ref(false);
const files = ref<any[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = ref(12);
const uploadingFiles = ref<any[]>([]);

function isImage(mimeType: string) {
  return mimeType && mimeType.startsWith("image/");
}

function formatSize(bytes: number) {
  if (!bytes) return "—";
  if (bytes < 1024) return bytes + " B";
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB";
  return (bytes / (1024 * 1024)).toFixed(1) + " MB";
}

async function handleFileChange(uploadFile: any) {
  const file = uploadFile.raw;
  if (!file) return;

  const item = { name: file.name, progress: 0, status: "" };
  uploadingFiles.value.push(item);

  try {
    const res = await uploadFile(file);
    item.progress = 100;
    item.status = "success";
    ElMessage.success(`${file.name} 上传成功`);
    loadFiles();
  } catch (e: any) {
    item.status = "exception";
    ElMessage.error(e.response?.data?.error || "上传失败");
  }

  // Remove from uploading list after a delay
  setTimeout(() => {
    const idx = uploadingFiles.value.indexOf(item);
    if (idx >= 0) uploadingFiles.value.splice(idx, 1);
  }, 2000);
}

async function loadFiles() {
  loading.value = true;
  try {
    const res = await getUploads({ page: page.value, page_size: pageSize.value });
    files.value = res.data.items || [];
    total.value = res.data.total || 0;
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
}

async function handleDelete(id: number) {
  try {
    await deleteUpload(id);
    ElMessage.success("已删除");
    loadFiles();
  } catch (e) {
    ElMessage.error("删除失败");
  }
}

function copyUrl(url: string) {
  navigator.clipboard.writeText(url).then(() => {
    ElMessage.success("URL 已复制: " + url);
  }).catch(() => {
    // Fallback
    const input = document.createElement("input");
    input.value = url;
    document.body.appendChild(input);
    input.select();
    document.execCommand("copy");
    document.body.removeChild(input);
    ElMessage.success("URL 已复制");
  });
}

onMounted(loadFiles);
</script>

<style scoped>
.page-title {
  font-size: 1.5rem;
  font-weight: 700;
  margin-bottom: 1rem;
  color: #1f2937;
}

.upload-card {
  margin-bottom: 1.5rem;
}

.upload-area {
  width: 100%;
}

.upload-icon {
  font-size: 3rem;
  color: #8b5cf6;
}

.upload-text {
  font-size: 1rem;
  color: #374151;
  margin-top: 0.5rem;
}

.upload-hint {
  font-size: 0.8125rem;
  color: #9ca3af;
  margin-top: 0.25rem;
}

.uploading-list {
  margin-top: 1rem;
}

.uploading-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.5rem 0;
}

.file-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 1rem;
}

.file-item {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  overflow: hidden;
  transition: box-shadow 0.2s;
}

.file-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.file-preview {
  aspect-ratio: 16 / 10;
  background: #f9fafb;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.file-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.file-info {
  padding: 0.5rem 0.75rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.file-name {
  font-size: 0.8125rem;
  color: #374151;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

.file-size {
  font-size: 0.75rem;
  color: #9ca3af;
  margin-left: 0.5rem;
}

.file-actions {
  padding: 0 0.75rem 0.75rem;
  display: flex;
  gap: 0.25rem;
}

.pagination-wrapper {
  margin-top: 1rem;
  display: flex;
  justify-content: flex-end;
}

.empty {
  text-align: center;
  color: #9ca3af;
  padding: 3rem 0;
}
</style>
