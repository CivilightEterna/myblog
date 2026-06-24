<template>
  <el-container class="admin-layout">
    <el-aside width="220px" class="admin-sidebar">
      <div class="sidebar-header">
        <h2>Blog Admin</h2>
      </div>
      <el-menu
        :default-active="currentRoute"
        router
        background-color="#1f2937"
        text-color="#9ca3af"
        active-text-color="#8b5cf6"
      >
        <el-menu-item index="/">
          <el-icon><DataAnalysis /></el-icon>
          <span>仪表盘</span>
        </el-menu-item>
        <el-menu-item index="/articles">
          <el-icon><Document /></el-icon>
          <span>文章管理</span>
        </el-menu-item>
        <el-menu-item index="/uploads">
          <el-icon><PictureFilled /></el-icon>
          <span>文件上传</span>
        </el-menu-item>
        <el-menu-item index="/settings">
          <el-icon><Setting /></el-icon>
          <span>站点设置</span>
        </el-menu-item>
        <el-menu-item index="/builds">
          <el-icon><UploadFilled /></el-icon>
          <span>构建发布</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="admin-header">
        <div class="header-left">
          <span class="header-greeting">欢迎，{{ authStore.nickname || authStore.username }}</span>
        </div>
        <div class="header-right">
          <el-button text @click="authStore.logout">
            <el-icon><SwitchButton /></el-icon>
            退出登录
          </el-button>
        </div>
      </el-header>

      <el-main class="admin-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useRoute } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import {
  DataAnalysis,
  Document,
  PictureFilled,
  Setting,
  UploadFilled,
  SwitchButton,
} from "@element-plus/icons-vue";

const route = useRoute();
const authStore = useAuthStore();

const currentRoute = computed(() => {
  if (route.path.startsWith("/articles")) return "/articles";
  return route.path;
});
</script>

<style scoped>
.admin-layout {
  min-height: 100vh;
  background: #f3f4f6;
}

.admin-sidebar {
  background: #1f2937;
  min-height: 100vh;
}

.sidebar-header {
  padding: 1.25rem;
  border-bottom: 1px solid #374151;
}

.sidebar-header h2 {
  color: #f9fafb;
  font-size: 1.125rem;
  font-weight: 700;
  margin: 0;
}

.el-menu {
  border-right: none;
}

.admin-header {
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #e5e7eb;
  padding: 0 1.5rem;
  height: 60px;
}

.header-greeting {
  color: #6b7280;
}

.admin-main {
  padding: 1.5rem;
  min-height: calc(100vh - 60px);
}
</style>
