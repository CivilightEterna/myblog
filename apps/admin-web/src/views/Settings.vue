<template>
  <div class="settings-page">
    <div class="page-header">
      <h1>站点设置</h1>
      <el-button type="primary" :loading="saving" @click="handleSave">保存设置</el-button>
    </div>

    <el-card shadow="never" v-loading="loading">
      <el-form :model="settings" label-position="top" label-width="120px">
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="站点标题">
              <el-input v-model="settings.title" placeholder="如：InkSpace" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="站点语言">
              <el-input v-model="settings.lang" placeholder="如：zh-CN" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="站点描述">
          <el-input v-model="settings.description" type="textarea" :rows="2" placeholder="博客描述（用于SEO）" />
        </el-form-item>

        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="作者">
              <el-input v-model="settings.author" placeholder="作者名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="站点 URL">
              <el-input v-model="settings.url" placeholder="https://your-domain.com" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="每页文章数">
          <el-input-number v-model="settings.posts_per_page" :min="1" :max="50" />
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { getSettings, updateSettings } from "@/api/setting";
import { ElMessage } from "element-plus";

const loading = ref(false);
const saving = ref(false);

const settings = reactive<Record<string, any>>({
  title: "InkSpace",
  description: "",
  author: "",
  url: "",
  lang: "zh-CN",
  posts_per_page: "10",
});

onMounted(async () => {
  loading.value = true;
  try {
    const res = await getSettings();
    Object.assign(settings, res.data);
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
});

async function handleSave() {
  saving.value = true;
  try {
    await updateSettings(settings as Record<string, string>);
    ElMessage.success("设置已保存");
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || "保存失败");
  } finally {
    saving.value = false;
  }
}
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
</style>
