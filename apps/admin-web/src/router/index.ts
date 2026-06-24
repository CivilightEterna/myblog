import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "@/stores/auth";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/login",
      name: "Login",
      component: () => import("@/views/Login.vue"),
      meta: { public: true },
    },
    {
      path: "/",
      component: () => import("@/components/AdminLayout.vue"),
      children: [
        {
          path: "",
          name: "Dashboard",
          component: () => import("@/views/Dashboard.vue"),
        },
        {
          path: "articles",
          name: "ArticleList",
          component: () => import("@/views/ArticleList.vue"),
        },
        {
          path: "articles/new",
          name: "ArticleNew",
          component: () => import("@/views/ArticleEditor.vue"),
        },
        {
          path: "articles/:id/edit",
          name: "ArticleEdit",
          component: () => import("@/views/ArticleEditor.vue"),
          props: true,
        },
        {
          path: "uploads",
          name: "UploadManager",
          component: () => import("@/views/UploadManager.vue"),
        },
        {
          path: "settings",
          name: "Settings",
          component: () => import("@/views/Settings.vue"),
        },
        {
          path: "images",
          name: "ImageManager",
          component: () => import("@/views/ImageManager.vue"),
        },
        {
          path: "builds",
          name: "BuildManager",
          component: () => import("@/views/BuildManager.vue"),
        },
      ],
    },
    {
      path: "/:pathMatch(.*)*",
      name: "NotFound",
      component: () => import("@/views/NotFound.vue"),
    },
  ],
});

router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore();
  if (to.meta.public || authStore.token) {
    next();
  } else {
    next("/login");
  }
});

export default router;
