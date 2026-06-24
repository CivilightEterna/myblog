<template>
  <div class="image-page">
    <div class="page-header">
      <h1>图床管理</h1>
      <el-upload
        :auto-upload="false"
        :show-file-list="false"
        :on-change="handleUpload"
        accept=".jpg,.jpeg,.png,.webp,.gif"
        :disabled="uploading"
      >
        <el-button type="primary" :loading="uploading">
          <el-icon><UploadFilled /></el-icon>
          {{ uploading ? "上传中..." : "上传图片" }}
        </el-button>
      </el-upload>
    </div>

    <!-- Upload dialog -->
    <el-dialog v-model="uploadDialog" title="上传图片" width="420px">
      <el-form label-position="top">
        <el-form-item label="分类">
          <el-select v-model="uploadCategory" style="width:100%">
            <el-option label="背景图 (backgrounds)" value="backgrounds" />
            <el-option label="文章配图 (articles)" value="articles" />
            <el-option label="头像 (avatars)" value="avatars" />
            <el-option label="通用 (common)" value="common" />
          </el-select>
        </el-form-item>
        <el-form-item label="文件">
          <span>{{ uploadFile?.name }}</span>
          <span class="file-size">{{ formatSize(uploadFile?.size || 0) }}</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="uploadDialog = false">取消</el-button>
        <el-button type="primary" :loading="uploading" @click="confirmUpload">确认上传</el-button>
      </template>
    </el-dialog>

    <!-- Filter bar -->
    <div class="filter-bar">
      <el-radio-group v-model="filterCategory" @change="loadImages" size="small">
        <el-radio-button value="">全部</el-radio-button>
        <el-radio-button value="backgrounds">背景图</el-radio-button>
        <el-radio-button value="articles">文章配图</el-radio-button>
        <el-radio-button value="avatars">头像</el-radio-button>
        <el-radio-button value="common">通用</el-radio-button>
      </el-radio-group>
      <el-input
        v-model="search"
        placeholder="搜索文件名..."
        clearable
        style="width:240px"
        size="small"
        @input="onSearch"
      />
      <span class="total-hint">共 {{ total }} 张</span>
    </div>

    <!-- Image grid -->
    <div v-loading="loading" class="image-grid">
      <div v-if="images.length === 0 && !loading" class="empty-hint">暂无图片</div>
      <div v-for="img in images" :key="img.id" class="image-item">
        <div class="image-preview" @click="previewImage = img.url">
          <img :src="img.url" :alt="img.original_name" loading="lazy" />
        </div>
        <div class="image-meta">
          <span class="img-name" :title="img.original_name">{{ img.original_name }}</span>
          <span class="img-cat">
            <el-tag size="small" type="info">{{ img.category }}</el-tag>
          </span>
          <span class="img-size">{{ formatSize(img.size) }}</span>
        </div>
        <div class="image-actions">
          <el-button text size="small" @click="copyUrl(img.url)">复制URL</el-button>
          <el-button text size="small" @click="copyMd(img)">复制MD</el-button>
          <el-button text size="small" type="success" @click="setBg(img)">设为背景</el-button>
          <el-button text size="small" type="warning" @click="setLanding(img)">设为封面</el-button>
          <el-button text size="small" @click="setAvatar(img)">设为头像</el-button>
          <el-popconfirm title="确认删除？" @confirm="handleDelete(img.id)">
            <template #reference>
              <el-button text size="small" type="danger">删除</el-button>
            </template>
          </el-popconfirm>
        </div>
      </div>
    </div>

    <!-- Pagination -->
    <div class="pagination-wrap" v-if="total > pageSize">
      <el-pagination
        v-model:current-page="page"
        :total="total"
        :page-size="pageSize"
        layout="total, prev, pager, next"
        @current-change="loadImages"
      />
    </div>

    <!-- Preview dialog -->
    <el-dialog v-model="previewVisible" title="预览" width="80%" top="5vh">
      <img :src="previewImage" style="width:100%" />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import client from "@/api/client";
import { getUploads } from "@/api/upload";
import { updateSettings } from "@/api/setting";
import { UploadFilled } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";

const loading = ref(false);
const uploading = ref(false);
const uploadDialog = ref(false);
const uploadCategory = ref("common");
const uploadFile = ref<File | null>(null);
const previewVisible = ref(false);
const previewImage = ref("");

const images = ref<any[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = ref(12);
const filterCategory = ref("");
const search = ref("");

let searchTimer: any = null;

function formatSize(bytes: number) {
  if (!bytes) return "—";
  if (bytes < 1024) return bytes + " B";
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB";
  return (bytes / (1024 * 1024)).toFixed(1) + " MB";
}

function handleUpload(uploadObj: any) {
  const file = uploadObj.raw as File;
  if (!file) return;

  if (file.size > 5 * 1024 * 1024) {
    ElMessage.error("文件大小不能超过 5MB");
    return;
  }
  const ext = file.name.split(".").pop()?.toLowerCase();
  if (!["jpg", "jpeg", "png", "webp", "gif"].includes(ext || "")) {
    ElMessage.error("只支持 jpg/png/webp/gif 格式");
    return;
  }

  uploadFile.value = file;
  uploadDialog.value = true;
}

async function confirmUpload() {
  if (!uploadFile.value) return;

  uploading.value = true;
  try {
    const fd = new FormData();
    fd.append("file", uploadFile.value);
    fd.append("category", uploadCategory.value);

    await client.post("/uploads", fd, {
      headers: { "Content-Type": "multipart/form-data" },
    });
    ElMessage.success("上传成功");
    uploadDialog.value = false;
    uploadFile.value = null;
    loadImages();
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || "上传失败");
  } finally {
    uploading.value = false;
  }
}

async function loadImages() {
  loading.value = true;
  try {
    const params: any = { page: page.value, page_size: pageSize.value };
    if (filterCategory.value) params.category = filterCategory.value;
    if (search.value) params.search = search.value;

    const res = await getUploads(params);
    images.value = res.data.items || [];
    total.value = res.data.total || 0;
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
}

function onSearch() {
  clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    page.value = 1;
    loadImages();
  }, 300);
}

async function copyUrl(url: string) {
  try {
    await navigator.clipboard.writeText(url);
  } catch {
    const input = document.createElement("input");
    input.value = url;
    document.body.appendChild(input);
    input.select();
    document.execCommand("copy");
    document.body.removeChild(input);
  }
  ElMessage.success("URL 已复制");
}

async function copyMd(img: any) {
  const md = `![${img.original_name || "image"}](${img.url})`;
  try {
    await navigator.clipboard.writeText(md);
  } catch {
    const input = document.createElement("input");
    input.value = md;
    document.body.appendChild(input);
    input.select();
    document.execCommand("copy");
    document.body.removeChild(input);
  }
  ElMessage.success("Markdown 已复制");
}

async function setBg(img: any) {
  try {
    await client.post(`/uploads/${img.id}/set-background`);
    ElMessage.success("已设置为前台背景图，刷新前台页面即可生效");
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || "设置失败");
  }
}

async function setLanding(img: any) {
  try {
    await client.put("/settings", {
      landing_cover_image: img.url,
      landing_enabled: "true",
    });
    ElMessage.success("已设置为首页封面图，刷新前台页面即可生效");
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || "设置失败");
  }
}

async function setAvatar(img: any) {
  try {
    await client.put("/settings", { site_avatar: img.url });
    ElMessage.success("已设置为头像，刷新前台页面即可生效");
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || "设置失败");
  }
}

async function handleDelete(id: number) {
  try {
    await client.delete(`/uploads/${id}`);
    ElMessage.success("已删除");
    loadImages();
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || "删除失败");
  }
}

onMounted(loadImages);
</script>

<style scoped>
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1rem;
}
.page-header h1 { font-size: 1.5rem; font-weight: 700; color: #1f2937; }

.filter-bar {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 1rem;
  flex-wrap: wrap;
}
.total-hint { font-size: 0.8125rem; color: #9ca3af; }

.image-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 1rem;
  min-height: 200px;
}
.empty-hint { grid-column: 1/-1; text-align: center; padding: 4rem; color: #9ca3af; }

.image-item {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  overflow: hidden;
  background: #fff;
  transition: box-shadow 0.2s;
}
.image-item:hover { box-shadow: 0 4px 12px rgba(0,0,0,0.1); }

.image-preview {
  aspect-ratio: 4/3;
  background: #f9fafb;
  cursor: pointer;
  overflow: hidden;
}
.image-preview img { width: 100%; height: 100%; object-fit: cover; }

.image-meta {
  padding: 0.5rem 0.75rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}
.img-name {
  flex: 1;
  font-size: 0.75rem;
  color: #374151;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.img-size { font-size: 0.6875rem; color: #9ca3af; white-space: nowrap; }

.image-actions {
  padding: 0 0.75rem 0.75rem;
  display: flex;
  flex-wrap: wrap;
  gap: 0;
}

.file-size { font-size: 0.75rem; color: #9ca3af; margin-left: 0.5rem; }
.pagination-wrap { margin-top: 1rem; display: flex; justify-content: flex-end; }
</style>
