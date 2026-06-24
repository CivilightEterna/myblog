<template>
  <div class="settings-page">
    <div class="page-header">
      <h1>站点设置</h1>
      <el-button type="primary" :loading="saving" @click="handleSave">保存设置</el-button>
    </div>

    <!-- Site Info -->
    <el-card shadow="never" v-loading="loading" class="settings-card">
      <template #header><span>站点基础信息</span></template>
      <el-form :model="settings" label-position="top" label-width="120px">
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="博客名称 (site_name)">
              <el-input v-model="settings.site_name" placeholder="如：InkSpace" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="作者名称 (site_author)">
              <el-input v-model="settings.site_author" placeholder="如：CivilightEterna" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="博客简介 (site_description)">
          <el-input v-model="settings.site_description" type="textarea" :rows="2" placeholder="一个安静写作的地方" />
        </el-form-item>
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="头像 URL (site_avatar)">
              <el-input v-model="settings.site_avatar" placeholder="输入 URL 或上传头像" />
              <div class="avatar-preview" v-if="settings.site_avatar">
                <img :src="settings.site_avatar" alt="头像预览" />
                <span class="preview-label">头像预览</span>
              </div>
              <div class="upload-row">
                <el-upload :auto-upload="false" :show-file-list="false" :on-change="uploadAvatar" accept=".jpg,.jpeg,.png,.webp,.gif">
                  <el-button size="small" text type="primary">上传头像</el-button>
                </el-upload>
              </div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Logo URL (site_logo)">
              <el-input v-model="settings.site_logo" placeholder="输入 URL 或上传 Logo（可选）" />
              <div class="avatar-preview" v-if="settings.site_logo">
                <img :src="settings.site_logo" alt="Logo 预览" />
                <span class="preview-label">Logo 预览</span>
              </div>
              <div class="upload-row">
                <el-upload :auto-upload="false" :show-file-list="false" :on-change="uploadLogo" accept=".jpg,.jpeg,.png,.webp,.gif">
                  <el-button size="small" text type="primary">上传 Logo</el-button>
                </el-upload>
              </div>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="站点语言">
              <el-input v-model="settings.lang" placeholder="zh-CN" />
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

        <!-- Legacy field sync -->
        <el-form-item v-if="settings.site_name" style="display:none">
          <el-input v-model="settings.title" />
          <el-input v-model="settings.description" />
          <el-input v-model="settings.author" />
        </el-form-item>
      </el-form>
    </el-card>

    <!-- Background Settings -->
    <el-card shadow="never" class="settings-card">
      <template #header>
        <div class="card-header">
          <span>前台背景图设置</span>
          <el-switch
            v-model="bgEnabled"
            active-text="启用"
            inactive-text="关闭"
            @change="onBgToggle"
          />
        </div>
      </template>

      <div v-if="bgEnabled" class="bg-settings">
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="上传背景图">
              <div class="bg-upload">
                <el-upload
                  :auto-upload="false"
                  :show-file-list="false"
                  :on-change="handleBgUpload"
                  accept=".jpg,.jpeg,.png,.webp"
                  :disabled="bgUploading"
                >
                  <el-button :loading="bgUploading" type="primary" plain>
                    {{ bgUploading ? "上传中..." : "选择图片上传" }}
                  </el-button>
                </el-upload>
                <span class="bg-hint">支持 jpg / png / webp，最大 5MB</span>
              </div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="或手动输入图片 URL">
              <el-input v-model="settings.site_background_image" placeholder="/uploads/backgrounds/xxx.jpg" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="背景透明度">
              <el-slider
                v-model="bgOpacity"
                :min="0.03"
                :max="0.6"
                :step="0.01"
                show-input
                :format-tooltip="(v: number) => `不透明度 ${Math.round(v * 100)}% 可见`"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="背景显示方式">
              <el-select v-model="settings.site_background_mode" style="width:100%">
                <el-option label="cover — 铺满并裁剪" value="cover" />
                <el-option label="contain — 完整显示" value="contain" />
                <el-option label="100% — 强制拉伸" value="100% 100%" />
                <el-option label="repeat — 平铺重复" value="repeat" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <!-- Preview -->
        <div v-if="settings.site_background_image" class="bg-preview-section">
          <span class="bg-preview-label">背景预览</span>
          <div class="bg-preview" :style="previewStyle">
            <div class="bg-preview-overlay">
              <span>预览效果</span>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="bg-disabled-hint">
        背景图已关闭，前台使用纯色背景。
      </div>
    </el-card>

    <!-- Comment Settings -->
    <el-card shadow="never" class="settings-card">
      <template #header>
        <div class="card-header">
          <span>评论系统设置</span>
          <el-switch
            v-model="commentEnabled"
            active-text="启用"
            inactive-text="关闭"
          />
        </div>
      </template>

      <div v-if="commentEnabled">
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="评论系统">
              <el-select v-model="settings.comment_provider" style="width:100%">
                <el-option label="关闭" value="none" />
                <el-option label="Twikoo（隐私优先）" value="twikoo" />
                <el-option label="Waline（即将支持）" value="waline" disabled />
                <el-option label="Giscus（即将支持）" value="giscus" disabled />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <!-- Twikoo config -->
        <div v-if="settings.comment_provider === 'twikoo'">
          <el-row :gutter="24">
            <el-col :span="12">
              <el-form-item label="Twikoo Env ID（必需）">
                <el-input
                  v-model="settings.twikoo_env_id"
                  placeholder="腾讯云环境 ID，如 https://xxx.ap-shanghai.tcb.qcloud.la"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="自部署地址（可选）">
                <el-input
                  v-model="settings.twikoo_server_url"
                  placeholder="留空则使用腾讯云默认地址"
                />
              </el-form-item>
            </el-col>
          </el-row>
          <p class="config-hint">
            <strong>如何获取 Env ID？</strong>
            前往 <a href="https://twikoo.js.org" target="_blank">Twikoo 文档</a>，
            使用腾讯云 CloudBase 免费部署，部署后在「环境管理」中复制环境 ID。<br/>
            自部署地址：仅当使用 Vercel/Netlify/Docker 自部署 Twikoo 后端时需要填写。
          </p>
        </div>

        <!-- Placeholder for future providers -->
        <div v-else-if="settings.comment_provider === 'waline'" class="provider-placeholder">
          Waline 评论系统支持即将推出。
        </div>
        <div v-else-if="settings.comment_provider === 'giscus'" class="provider-placeholder">
          Giscus 评论系统支持即将推出。
        </div>
      </div>

      <div v-else class="bg-disabled-hint">
        评论功能已关闭。
      </div>
    </el-card>

    <!-- Landing Page Settings -->
    <el-card shadow="never" class="settings-card">
      <template #header>
        <div class="card-header">
          <span>首页封面设置</span>
          <el-switch v-model="landingEnabled" active-text="启用" inactive-text="关闭" />
        </div>
      </template>

      <div v-if="landingEnabled">
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="封面背景图 URL">
              <el-input v-model="settings.landing_cover_image" placeholder="输入 URL 或从图床选择" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="封面遮罩透明度">
              <el-slider v-model="landingOverlay" :min="0" :max="0.8" :step="0.01" show-input />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="24">
          <el-col :span="8">
            <el-form-item label="主标题"><el-input v-model="settings.landing_title" /></el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="副标题"><el-input v-model="settings.landing_subtitle" /></el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="文字位置">
              <el-select v-model="settings.landing_text_position" style="width:100%">
                <el-option label="居中" value="center" />
                <el-option label="左对齐" value="left" />
                <el-option label="右对齐" value="right" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="描述"><el-input v-model="settings.landing_description" type="textarea" :rows="2" /></el-form-item>

        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="按钮文字"><el-input v-model="settings.landing_button_text" /></el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="按钮链接"><el-input v-model="settings.landing_button_link" /></el-form-item>
          </el-col>
        </el-row>

        <el-divider>封面格言</el-divider>
        <el-form-item>
          <el-checkbox v-model="quoteEnabled">启用封面格言</el-checkbox>
        </el-form-item>
        <div v-if="quoteEnabled">
          <el-form-item label="格言内容"><el-input v-model="settings.landing_quote_text" type="textarea" :rows="2" /></el-form-item>
          <el-row :gutter="24">
            <el-col :span="12">
              <el-form-item label="格言作者"><el-input v-model="settings.landing_quote_author" /></el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item>
                <el-checkbox v-model="typingEnabled">打字机效果</el-checkbox>
              </el-form-item>
            </el-col>
          </el-row>
        </div>

        <el-divider>动态效果</el-divider>
        <el-form-item>
          <el-checkbox v-model="animEnabled">启用动态效果</el-checkbox>
        </el-form-item>
        <div v-if="animEnabled">
          <el-row :gutter="24">
            <el-col :span="12">
              <el-form-item label="效果类型">
                <el-select v-model="settings.landing_animation_style" style="width:100%">
                  <el-option label="无" value="none" />
                  <el-option label="🌸 樱花飘落" value="sakura" />
                  <el-option label="✨ 星光粒子" value="particles" />
                  <el-option label="🔮 柔和光斑" value="glow" />
                  <el-option label="🌧 细雨" value="rain" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="效果强度">
                <el-select v-model="settings.landing_animation_intensity" style="width:100%">
                  <el-option label="低" value="low" />
                  <el-option label="中" value="medium" />
                  <el-option label="高" value="high" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
        </div>

        <el-divider>风格选项</el-divider>
        <el-row :gutter="24">
          <el-col :span="8">
            <el-form-item>
              <el-checkbox v-model="animeStyle">二次元风格增强</el-checkbox>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item>
              <el-checkbox v-model="blurCard">毛玻璃卡片</el-checkbox>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="强调色">
              <el-color-picker v-model="settings.landing_accent_color" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item v-if="blurCard" label="卡片透明度">
          <el-slider v-model="landingCardOpacity" :min="0.05" :max="0.6" :step="0.01" show-input />
        </el-form-item>
      </div>

      <div v-else class="bg-disabled-hint">
        首页封面已关闭，访问 / 将直接跳转到 /blog。
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from "vue";
import { getSettings, updateSettings } from "@/api/setting";
import client from "@/api/client";
import { ElMessage } from "element-plus";

const loading = ref(false);
const saving = ref(false);
const bgUploading = ref(false);

const settings = reactive<Record<string, any>>({
  title: "InkSpace",
  description: "",
  author: "",
  url: "",
  lang: "zh-CN",
  posts_per_page: "10",
  site_name: "InkSpace",
  site_description: "一个安静写作的地方",
  site_author: "",
  site_avatar: "",
  site_logo: "",
  site_background_enabled: "false",
  site_background_image: "",
  site_background_opacity: "0.18",
  site_background_mode: "cover",
  comment_enabled: "false",
  comment_provider: "none",
  twikoo_env_id: "",
  twikoo_server_url: "",
  landing_enabled: "true",
  landing_cover_image: "",
  landing_title: "InkSpace",
  landing_subtitle: "",
  landing_description: "",
  landing_button_text: "进入博客",
  landing_button_link: "/blog",
  landing_overlay_opacity: "0.35",
  landing_text_position: "center",
  landing_quote_enabled: "true",
  landing_quote_text: "",
  landing_quote_author: "",
  landing_quote_typing_effect: "true",
  landing_animation_enabled: "true",
  landing_animation_style: "sakura",
  landing_animation_intensity: "medium",
  landing_anime_style_enabled: "true",
  landing_blur_card_enabled: "true",
  landing_card_opacity: "0.35",
  landing_accent_color: "#8b5cf6",
});

const bgEnabled = computed({
  get: () => settings.site_background_enabled === "true",
  set: (v: boolean) => { settings.site_background_enabled = v ? "true" : "false"; }
});

const commentEnabled = computed({
  get: () => settings.comment_enabled === "true",
  set: (v: boolean) => { settings.comment_enabled = v ? "true" : "false"; }
});

const landingEnabled = computed({
  get: () => settings.landing_enabled === "true",
  set: (v: boolean) => { settings.landing_enabled = v ? "true" : "false"; }
});
const quoteEnabled = computed({
  get: () => settings.landing_quote_enabled === "true",
  set: (v: boolean) => { settings.landing_quote_enabled = v ? "true" : "false"; }
});
const typingEnabled = computed({
  get: () => settings.landing_quote_typing_effect === "true",
  set: (v: boolean) => { settings.landing_quote_typing_effect = v ? "true" : "false"; }
});
const animEnabled = computed({
  get: () => settings.landing_animation_enabled === "true",
  set: (v: boolean) => { settings.landing_animation_enabled = v ? "true" : "false"; }
});
const animeStyle = computed({
  get: () => settings.landing_anime_style_enabled === "true",
  set: (v: boolean) => { settings.landing_anime_style_enabled = v ? "true" : "false"; }
});
const blurCard = computed({
  get: () => settings.landing_blur_card_enabled === "true",
  set: (v: boolean) => { settings.landing_blur_card_enabled = v ? "true" : "false"; }
});
const landingOverlay = computed({
  get: () => parseFloat(settings.landing_overlay_opacity) || 0.35,
  set: (v: number) => { settings.landing_overlay_opacity = String(v); }
});
const landingCardOpacity = computed({
  get: () => parseFloat(settings.landing_card_opacity) || 0.35,
  set: (v: number) => { settings.landing_card_opacity = String(v); }
});

const bgOpacity = computed({
  get: () => parseFloat(settings.site_background_opacity) || 0.18,
  set: (v: number) => { settings.site_background_opacity = String(v); }
});

const previewStyle = computed(() => ({
  backgroundImage: `url(${settings.site_background_image})`,
  backgroundSize: settings.site_background_mode || "cover",
  backgroundPosition: "center",
  opacity: bgOpacity.value * 4, // Amplify for preview visibility
}));

function onBgToggle(enabled: boolean) {
  settings.site_background_enabled = enabled ? "true" : "false";
}

async function handleBgUpload(uploadFileObj: any) {
  const file = uploadFileObj.raw;
  if (!file) return;

  // Size check (5MB)
  if (file.size > 5 * 1024 * 1024) {
    ElMessage.error("文件大小不能超过 5MB");
    return;
  }

  // Type check
  const ext = file.name.split(".").pop()?.toLowerCase();
  if (!["jpg", "jpeg", "png", "webp"].includes(ext || "")) {
    ElMessage.error("只支持 jpg、png、webp 格式");
    return;
  }

  bgUploading.value = true;
  try {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("type", "background");

    const uploadRes = await client.post("/uploads", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });
    settings.site_background_image = uploadRes.data.url;
    ElMessage.success("背景图上传成功");
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || "上传失败");
  } finally {
    bgUploading.value = false;
  }
}

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

async function uploadAvatar(uploadObj: any) {
  const file = uploadObj.raw; if (!file) return;
  await uploadFileTo('avatars', file, 'site_avatar', '头像');
}
async function uploadLogo(uploadObj: any) {
  const file = uploadObj.raw; if (!file) return;
  await uploadFileTo('avatars', file, 'site_logo', 'Logo');
}
async function uploadFileTo(category: string, file: File, settingKey: string, label: string) {
  if (file.size > 2 * 1024 * 1024) { ElMessage.error(label + '文件不能超过 2MB'); return; }
  const ext = file.name.split('.').pop()?.toLowerCase();
  if (!['jpg','jpeg','png','webp','gif'].includes(ext||'')) { ElMessage.error('只支持 jpg/png/webp/gif'); return; }
  try {
    const fd = new FormData(); fd.append('file', file); fd.append('category', category);
    const res = await client.post('/uploads', fd, { headers: { 'Content-Type': 'multipart/form-data' } });
    (settings as any)[settingKey] = res.data.url;
    ElMessage.success(label + '上传成功');
  } catch (e: any) { ElMessage.error(e.response?.data?.error || '上传失败'); }
}

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

.settings-card {
  margin-bottom: 1.5rem;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.bg-settings {
  margin-top: 0.5rem;
}

.bg-upload {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.bg-hint {
  font-size: 0.75rem;
  color: #9ca3af;
}

.bg-preview-section {
  margin-top: 1rem;
}

.bg-preview-label {
  display: block;
  font-size: 0.875rem;
  color: #6b7280;
  margin-bottom: 0.5rem;
}

.bg-preview {
  width: 100%;
  height: 160px;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  overflow: hidden;
}

.bg-preview-overlay {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.7);
  color: #374151;
  font-size: 0.875rem;
}

.bg-disabled-hint {
  text-align: center;
  color: #9ca3af;
  padding: 2rem 0;
  font-size: 0.9375rem;
}

.config-hint {
  font-size: 0.8125rem;
  color: #6b7280;
  line-height: 1.7;
  background: #f9fafb;
  padding: 0.75rem 1rem;
  border-radius: 6px;
  margin-top: 0.5rem;
}
.config-hint a { color: #8b5cf6; }

.provider-placeholder {
  text-align: center;
  padding: 2rem 0;
  color: #9ca3af;
  font-size: 0.9375rem;
}
</style>
