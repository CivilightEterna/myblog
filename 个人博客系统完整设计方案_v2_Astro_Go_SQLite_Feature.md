# 个人博客系统完整设计方案 v2

> 方案定位：基于 **Astro** 自研博客前台，不套用现成主题；使用 **Go + Gin** 自研后台；使用 **SQLite** 作为后台工作台数据库；文章最终保存为 **Markdown + Frontmatter**；通过 **Nginx** 托管 Astro 构建后的静态页面。  
>
> 本版本新增重点：**可扩展主题系统、Feature 插件系统、写作热力图、Live2D/桌宠集成、后台功能开关、构建时静态数据生成、未来迁移升级策略**。

---

## 1. 项目目标

本项目目标是开发一个适合个人长期使用的博客系统，既要轻量，又要有较好的扩展性。

### 1.1 核心目标

1. 支持个人博客文章发布。
2. 支持 Markdown 写作。
3. 支持后台管理文章、图片、配置和主题。
4. 前台使用 Astro 自研，不套用 Mizuki、Fuwari 等现成主题。
5. 支持主题切换和布局扩展。
6. 支持写作热力图。
7. 支持桌宠 / Live2D 看板娘。
8. 支持后续扩展音乐播放器、评论、访问统计、AI 摘要等功能。
9. 适配 2 核 2G 小服务器。
10. 保证文章内容迁移方便，不被数据库锁死。

### 1.2 设计原则

```text
1. 前台访问尽量静态化。
2. 后台只服务管理员，不参与普通访客文章访问。
3. 数据库只作为后台工作台。
4. 已发布文章最终以 Markdown 文件保存。
5. 主题系统和功能插件系统分离。
6. 所有可选功能必须可关闭。
7. 交互功能尽量懒加载，不影响首屏性能。
8. 文章 URL 使用 slug，不使用数据库 id。
9. 图片资源统一放到 uploads 目录。
10. 支持备份、导出、迁移和回滚。
```

---

## 2. 总体架构

### 2.1 访问架构

```text
普通访客
   |
   v
Nginx
   |
   v
Astro 构建后的静态页面
   |
   v
HTML / CSS / JS / 图片 / JSON 静态数据
```

普通访客访问博客时，不访问 Go 后台，也不查询 SQLite。

这样即使后台服务挂了，博客前台仍然可以正常访问。

### 2.2 管理架构

```text
管理员
   |
   v
后台管理页面 Admin Web
   |
   v
Go Admin API
   |
   +------------------------+
   |                        |
   v                        v
SQLite 数据库           Markdown / 图片文件
   |                        |
   +-----------+------------+
               |
               v
        构建前生成配置与数据
               |
               v
        执行 Astro build
               |
               v
        生成 dist 静态文件
               |
               v
        切换 releases/current
```

### 2.3 系统模块图

```text
Blog System
├─ Astro Blog Frontend
│  ├─ Content Core
│  ├─ Theme Core
│  ├─ Feature Core
│  └─ Build Generated Data
│
├─ Admin Web
│  ├─ Article Editor
│  ├─ Upload Manager
│  ├─ Site Setting
│  ├─ Theme Setting
│  ├─ Feature Setting
│  ├─ Build Manager
│  └─ Backup Manager
│
├─ Admin API
│  ├─ Auth Module
│  ├─ Article Module
│  ├─ Markdown Module
│  ├─ Upload Module
│  ├─ Setting Module
│  ├─ Feature Module
│  ├─ Build Module
│  └─ Backup Module
│
├─ SQLite
│  ├─ admin_user
│  ├─ article_draft
│  ├─ upload_file
│  ├─ site_setting
│  ├─ theme_setting
│  ├─ feature_setting
│  ├─ writing_activity
│  ├─ build_record
│  └─ backup_record
│
└─ Deploy Runtime
   ├─ Nginx
   ├─ Go Admin API
   ├─ SQLite File
   ├─ Astro dist
   └─ releases/current
```

---

## 3. 技术选型

### 3.1 前台博客

| 技术 | 用途 |
|---|---|
| Astro | 静态博客前台框架 |
| TypeScript | 类型约束 |
| Markdown / MDX | 文章内容格式 |
| Content Collections | 文章内容集合与 schema 校验 |
| Tailwind CSS / UnoCSS | 样式开发 |
| CSS Variables | 主题变量 |
| Shiki / Expressive Code | 代码高亮 |
| KaTeX | 数学公式 |
| Pagefind | 静态搜索 |
| RSS | 订阅输出 |
| Sitemap | 站点地图 |
| Astro Islands | 桌宠、搜索框、音乐播放器等交互组件懒加载 |

推荐组合：

```text
Astro + TypeScript + Tailwind CSS + Content Collections + Pagefind
```

### 3.2 后台服务

| 技术 | 用途 |
|---|---|
| Go | 后台服务语言 |
| Gin | HTTP API 框架 |
| SQLite | 轻量数据库 |
| JWT / Session | 登录态 |
| bcrypt | 密码加密 |
| zap / zerolog | 日志 |
| Viper | 配置管理 |
| systemd | 服务守护 |

推荐组合：

```text
Go + Gin + SQLite + JWT + bcrypt
```

### 3.3 后台前端

| 技术 | 用途 |
|---|---|
| Vue 3 | 后台 SPA |
| Vite | 前端构建 |
| TypeScript | 类型约束 |
| Element Plus | 后台 UI |
| MdEditorV3 / Monaco | Markdown 编辑器 |
| Axios | API 请求 |

推荐组合：

```text
Vue 3 + Vite + TypeScript + Element Plus + MdEditorV3
```

### 3.4 部署

| 技术 | 用途 |
|---|---|
| Nginx | 静态资源托管与反向代理 |
| systemd | 后端进程守护 |
| pnpm | Astro 构建 |
| Git | 版本管理，可选 |
| Shell Scripts | 构建、备份、回滚脚本 |

---

## 4. 项目目录设计

推荐使用 monorepo。

```text
my-blog/
├─ apps/
│  ├─ blog-web/                 # Astro 博客前台
│  ├─ admin-web/                # Vue 后台管理页面
│  └─ admin-api/                # Go 后台 API
│
├─ data/
│  └─ blog_admin.db             # SQLite 数据库
│
├─ storage/
│  ├─ uploads/                  # 原始上传备份，可选
│  ├─ backups/                  # 备份文件
│  └─ logs/                     # 日志
│
├─ releases/                    # 每次发布后的静态版本
├─ current -> releases/xxx      # Nginx 当前访问版本
│
├─ scripts/
│  ├─ build.sh                  # 构建脚本
│  ├─ deploy.sh                 # 发布脚本
│  ├─ backup.sh                 # 备份脚本
│  ├─ rollback.sh               # 回滚脚本
│  └─ migrations/               # 内容迁移脚本
│
└─ README.md
```

---

## 5. Astro 前台目录结构

```text
apps/blog-web/
├─ src/
│  ├─ pages/
│  │  ├─ index.astro
│  │  ├─ posts/
│  │  │  └─ [...slug].astro
│  │  ├─ categories/
│  │  │  └─ [category].astro
│  │  ├─ tags/
│  │  │  └─ [tag].astro
│  │  ├─ archives.astro
│  │  ├─ about.astro
│  │  ├─ search.astro
│  │  ├─ 404.astro
│  │  └─ rss.xml.ts
│  │
│  ├─ content/
│  │  ├─ posts/
│  │  │  ├─ 2026/
│  │  │  │  └─ 06/
│  │  │  │     └─ redis-cache-breakdown.md
│  │  └─ config.ts
│  │
│  ├─ layouts/
│  │  ├─ BaseLayout.astro
│  │  ├─ PageLayout.astro
│  │  ├─ PostLayout.astro
│  │  ├─ home/
│  │  │  ├─ ClassicHome.astro
│  │  │  ├─ MinimalHome.astro
│  │  │  ├─ MagazineHome.astro
│  │  │  └─ TerminalHome.astro
│  │  └─ post/
│  │     ├─ SidebarPost.astro
│  │     ├─ CleanPost.astro
│  │     └─ WidePost.astro
│  │
│  ├─ components/
│  │  ├─ Header.astro
│  │  ├─ Footer.astro
│  │  ├─ PostCard.astro
│  │  ├─ PostMeta.astro
│  │  ├─ TagList.astro
│  │  ├─ CategoryBadge.astro
│  │  ├─ ThemeToggle.astro
│  │  ├─ Toc.astro
│  │  ├─ SearchBox.astro
│  │  ├─ Pagination.astro
│  │  └─ Icon.astro
│  │
│  ├─ features/
│  │  ├─ index.ts
│  │  ├─ FeatureHost.astro
│  │  ├─ heatmap/
│  │  │  ├─ Heatmap.astro
│  │  │  ├─ HeatmapIsland.vue
│  │  │  ├─ config.ts
│  │  │  └─ heatmap.css
│  │  ├─ pet/
│  │  │  ├─ PetIsland.vue
│  │  │  ├─ config.ts
│  │  │  └─ pet.css
│  │  ├─ comment/
│  │  ├─ music/
│  │  ├─ analytics/
│  │  └─ ai-summary/
│  │
│  ├─ config/
│  │  ├─ site.config.ts
│  │  ├─ site.generated.ts
│  │  ├─ theme.config.ts
│  │  ├─ theme.generated.ts
│  │  ├─ feature.config.ts
│  │  └─ feature.generated.ts
│  │
│  ├─ styles/
│  │  ├─ global.css
│  │  ├─ variables.css
│  │  ├─ markdown.css
│  │  └─ themes/
│  │     ├─ default.css
│  │     ├─ dark.css
│  │     ├─ minimal.css
│  │     └─ terminal.css
│  │
│  └─ utils/
│     ├─ post.ts
│     ├─ date.ts
│     ├─ slug.ts
│     ├─ reading-time.ts
│     └─ feature.ts
│
├─ public/
│  ├─ uploads/
│  │  └─ 2026/
│  │     └─ 06/
│  ├─ data/
│  │  ├─ writing-heatmap.json
│  │  ├─ site-stat.json
│  │  └─ feature-data.json
│  ├─ pet/
│  │  └─ models/
│  ├─ favicon.svg
│  └─ robots.txt
│
├─ astro.config.mjs
├─ package.json
└─ tsconfig.json
```

---

## 6. 内容系统设计

### 6.1 文章存储方式

已发布文章最终保存为 Markdown 文件。

路径示例：

```text
apps/blog-web/src/content/posts/2026/06/redis-cache-breakdown.md
```

### 6.2 文章 Frontmatter 标准

```markdown
---
contentVersion: 1
title: "Redis 缓存击穿怎么解决"
slug: "redis-cache-breakdown"
description: "整理 Redis 缓存击穿、缓存穿透、缓存雪崩的区别和解决方案。"
date: "2026-06-24"
updated: "2026-06-24"
category: "后端"
tags:
  - Redis
  - 缓存
  - Java面试
cover: "/uploads/2026/06/redis-cover.png"
draft: false
pinned: false
comment: true
toc: true
---

# Redis 缓存击穿怎么解决

正文内容……
```

### 6.3 字段说明

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| contentVersion | number | 是 | 内容格式版本 |
| title | string | 是 | 文章标题 |
| slug | string | 是 | 文章唯一 URL 标识 |
| description | string | 是 | 文章摘要 |
| date | string | 是 | 发布时间 |
| updated | string | 否 | 更新时间 |
| category | string | 是 | 分类 |
| tags | string[] | 否 | 标签 |
| cover | string | 否 | 封面图 |
| draft | boolean | 是 | 是否草稿 |
| pinned | boolean | 否 | 是否置顶 |
| comment | boolean | 否 | 是否允许评论 |
| toc | boolean | 否 | 是否显示目录 |

### 6.4 URL 设计

推荐：

```text
/posts/redis-cache-breakdown
```

不推荐：

```text
/posts/123
```

原因：

1. slug 稳定。
2. 迁移时链接不变。
3. SEO 更好。
4. 不依赖数据库 id。
5. Markdown 文件天然适合 slug 命名。

---

## 7. Astro Content Collections

`src/content/config.ts` 示例：

```ts
import { defineCollection, z } from "astro:content";

const posts = defineCollection({
  type: "content",
  schema: z.object({
    contentVersion: z.number().default(1),
    title: z.string(),
    slug: z.string(),
    description: z.string(),
    date: z.string(),
    updated: z.string().optional(),
    category: z.string(),
    tags: z.array(z.string()).default([]),
    cover: z.string().optional(),
    draft: z.boolean().default(false),
    pinned: z.boolean().default(false),
    comment: z.boolean().default(true),
    toc: z.boolean().default(true),
  }),
});

export const collections = {
  posts,
};
```

作用：

1. 构建时校验文章字段。
2. 防止后台生成错误 Markdown。
3. 保证文章数据结构稳定。
4. 方便后续迁移。

---

## 8. 前台页面设计

### 8.1 首页 `/`

展示：

1. Banner。
2. 个人简介。
3. 置顶文章。
4. 最新文章。
5. 分类入口。
6. 标签云。
7. 写作热力图。
8. 侧边栏 Feature 插槽。
9. 页脚。

排序规则：

```text
过滤 draft = false
先按 pinned 排序
再按 date 倒序
```

### 8.2 文章页 `/posts/[slug]`

展示：

1. 标题。
2. 发布时间。
3. 更新时间。
4. 分类。
5. 标签。
6. 阅读时间。
7. 封面图。
8. TOC 目录。
9. Markdown 正文。
10. 上一篇 / 下一篇。
11. 评论区 Feature 插槽。
12. 文章底部 Feature 插槽。

### 8.3 分类页 `/categories/[category]`

展示某个分类下的文章。

### 8.4 标签页 `/tags/[tag]`

展示某个标签下的文章。

### 8.5 归档页 `/archives`

按年月分组：

```text
2026 年 6 月
- Redis 缓存击穿怎么解决
- MySQL B+ 树索引原理

2026 年 5 月
- Java 线程池拒绝策略
```

### 8.6 搜索页 `/search`

使用 Pagefind 静态搜索。

优点：

1. 不需要后端。
2. 不需要数据库。
3. 构建时生成索引。
4. 适合静态博客。

---

## 9. 主题系统设计

### 9.1 主题系统目标

主题系统负责视觉和布局，不负责具体功能插件。

支持：

1. 明暗模式。
2. 主题色。
3. 字体。
4. 圆角。
5. 阴影。
6. 首页布局。
7. 文章页布局。
8. Banner 风格。
9. 卡片风格。
10. 后台配置主题。

### 9.2 CSS 变量层

`src/styles/variables.css`：

```css
:root {
  --color-bg: #ffffff;
  --color-bg-soft: #f8fafc;
  --color-text: #111827;
  --color-text-muted: #6b7280;
  --color-primary: #3b82f6;
  --color-card: #ffffff;
  --color-border: #e5e7eb;

  --font-body: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  --font-code: "JetBrains Mono", "Fira Code", monospace;

  --radius-sm: 6px;
  --radius-md: 12px;
  --radius-lg: 18px;

  --shadow-card: 0 10px 30px rgba(15, 23, 42, 0.08);
}
```

暗色主题：

```css
[data-theme="dark"] {
  --color-bg: #020617;
  --color-bg-soft: #0f172a;
  --color-text: #e5e7eb;
  --color-text-muted: #94a3b8;
  --color-primary: #8b5cf6;
  --color-card: #111827;
  --color-border: #1f2937;

  --shadow-card: 0 10px 30px rgba(0, 0, 0, 0.35);
}
```

### 9.3 主题预设

`src/config/theme.config.ts`：

```ts
export const themePresets = {
  default: {
    label: "默认主题",
    mode: "light",
    primaryColor: "#3b82f6",
    homeLayout: "classic",
    postLayout: "sidebar",
    radius: "large",
  },

  minimal: {
    label: "极简主题",
    mode: "light",
    primaryColor: "#111827",
    homeLayout: "minimal",
    postLayout: "clean",
    radius: "small",
  },

  dark: {
    label: "暗黑主题",
    mode: "dark",
    primaryColor: "#8b5cf6",
    homeLayout: "classic",
    postLayout: "sidebar",
    radius: "large",
  },

  terminal: {
    label: "终端主题",
    mode: "dark",
    primaryColor: "#22c55e",
    homeLayout: "terminal",
    postLayout: "clean",
    radius: "small",
  },
};
```

### 9.4 布局扩展

首页布局：

```text
src/layouts/home/
├─ ClassicHome.astro
├─ MinimalHome.astro
├─ MagazineHome.astro
└─ TerminalHome.astro
```

文章页布局：

```text
src/layouts/post/
├─ SidebarPost.astro
├─ CleanPost.astro
└─ WidePost.astro
```

后台切换布局时，生成：

```ts
export const activeTheme = {
  name: "minimal",
  mode: "light",
  primaryColor: "#111827",
  homeLayout: "minimal",
  postLayout: "clean",
};
```

---

## 10. Feature 插件系统设计

### 10.1 为什么需要 Feature 系统

写作热力图、桌宠、音乐播放器、评论、访问统计、AI 摘要这些功能不应该写死在核心布局里。

否则后期会出现：

1. BaseLayout 越来越臃肿。
2. 每加一个功能都要改多个页面。
3. 主题和功能互相污染。
4. 功能无法关闭。
5. 交互组件影响首屏性能。
6. 后续迁移困难。

因此需要单独设计 Feature Core。

### 10.2 Feature Core 目标

Feature Core 负责：

1. 功能注册。
2. 功能启用/关闭。
3. 功能挂载位置。
4. 功能配置。
5. 功能静态数据。
6. 交互组件加载策略。
7. 后台配置管理。

### 10.3 Feature 目录结构

```text
src/features/
├─ index.ts
├─ FeatureHost.astro
├─ heatmap/
│  ├─ Heatmap.astro
│  ├─ HeatmapIsland.vue
│  ├─ config.ts
│  └─ heatmap.css
├─ pet/
│  ├─ PetIsland.vue
│  ├─ config.ts
│  └─ pet.css
├─ comment/
├─ music/
├─ analytics/
└─ ai-summary/
```

### 10.4 Feature 类型定义

```ts
export type FeaturePosition =
  | "global-header"
  | "global-footer"
  | "global-floating"
  | "home-before-list"
  | "home-sidebar"
  | "home-after-list"
  | "post-before-content"
  | "post-sidebar"
  | "post-after-content"
  | "post-footer";

export interface BlogFeature {
  key: string;
  name: string;
  enabled: boolean;
  position: FeaturePosition;
  component: string;
  client?: "load" | "idle" | "visible" | "only";
  config?: Record<string, any>;
}
```

### 10.5 Feature 配置示例

写作热力图：

```ts
{
  key: "writing-heatmap",
  name: "写作热力图",
  enabled: true,
  position: "home-sidebar",
  component: "Heatmap",
  config: {
    type: "writing",
    showWords: true,
    showPublishCount: true
  }
}
```

桌宠：

```ts
{
  key: "live2d-pet",
  name: "Live2D 桌宠",
  enabled: true,
  position: "global-floating",
  component: "PetIsland",
  client: "idle",
  config: {
    model: "/pet/models/hiyori/model3.json",
    width: 280,
    height: 380,
    showOnMobile: false
  }
}
```

### 10.6 Feature Host

布局中只放统一挂载点。

```astro
---
import FeatureHost from "@/features/FeatureHost.astro";
---

<html>
  <body>
    <slot />

    <FeatureHost position="global-floating" />
  </body>
</html>
```

首页侧边栏：

```astro
<aside>
  <FeatureHost position="home-sidebar" />
</aside>
```

文章底部：

```astro
<FeatureHost position="post-after-content" />
```

### 10.7 Feature 插槽列表

全局：

```text
global-header
global-footer
global-floating
```

首页：

```text
home-before-list
home-sidebar
home-after-list
```

文章页：

```text
post-before-content
post-sidebar
post-after-content
post-footer
```

---

## 11. 写作热力图设计

### 11.1 功能目标

写作热力图用于展示博客作者的写作活跃情况。

支持两种统计方式：

1. 发布热力图。
2. 真实写作热力图。

### 11.2 方案 A：发布热力图

数据来源：

```text
Markdown frontmatter date 字段
```

统计逻辑：

```text
扫描所有已发布 Markdown
按 date 分组
统计每天发布文章数量
生成 writing-heatmap.json
```

生成 JSON：

```json
[
  {
    "date": "2026-06-24",
    "count": 2,
    "words": 0,
    "type": "publish"
  },
  {
    "date": "2026-06-25",
    "count": 1,
    "words": 0,
    "type": "publish"
  }
]
```

优点：

1. 不需要数据库。
2. 完全静态。
3. 简单稳定。
4. 适合第一版。

缺点：

1. 只能统计发布行为。
2. 无法统计每天写了多少字。

### 11.3 方案 B：真实写作热力图

数据来源：

```text
SQLite writing_activity 表
```

记录行为：

1. 保存草稿。
2. 更新草稿。
3. 发布文章。
4. 修改文章。
5. 删除文章。

新增表：

```sql
CREATE TABLE writing_activity (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    article_id INTEGER,
    article_slug TEXT,
    action TEXT NOT NULL,
    word_count INTEGER DEFAULT 0,
    delta_word_count INTEGER DEFAULT 0,
    activity_date DATE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

字段说明：

| 字段 | 说明 |
|---|---|
| article_id | 草稿文章 id |
| article_slug | 文章 slug |
| action | 行为类型 |
| word_count | 当前总字数 |
| delta_word_count | 本次新增或减少字数 |
| activity_date | 活跃日期 |
| created_at | 记录创建时间 |

行为类型：

```text
create_draft
save_draft
update_draft
publish_article
update_published
delete_article
```

生成 JSON：

```json
[
  {
    "date": "2026-06-24",
    "count": 3,
    "words": 1200,
    "publishCount": 1
  },
  {
    "date": "2026-06-25",
    "count": 1,
    "words": 800,
    "publishCount": 0
  }
]
```

### 11.4 前台组件设计

```text
src/features/heatmap/
├─ Heatmap.astro
├─ HeatmapIsland.vue
├─ config.ts
└─ heatmap.css
```

Astro 静态组件负责外壳：

```astro
<section class="feature-card heatmap-card">
  <h2>写作热力图</h2>
  <HeatmapIsland client:visible />
</section>
```

Vue Island 负责交互：

```text
读取 /data/writing-heatmap.json
渲染日历格子
hover 显示当天写作数据
```

### 11.5 热力图展示内容

每个日期格子显示：

1. 日期。
2. 写作次数。
3. 新增字数。
4. 发布文章数量。

tooltip 示例：

```text
2026-06-24
写作 3 次
新增 1200 字
发布 1 篇
```

---

## 12. 桌宠 / Live2D 集成设计

### 12.1 功能目标

桌宠用于增强博客个性化，可以放在右下角。

支持：

1. 开启 / 关闭。
2. 设置模型。
3. 设置位置。
4. 设置大小。
5. 移动端隐藏。
6. 空闲加载。
7. 不影响首屏性能。

### 12.2 推荐实现方式

作为 Feature 插件接入：

```text
src/features/pet/
├─ PetIsland.vue
├─ config.ts
└─ pet.css
```

挂载位置：

```text
global-floating
```

加载策略：

```text
client:idle
```

说明：

桌宠不属于首屏核心内容，应该在浏览器空闲时加载。

### 12.3 配置示例

```json
{
  "feature_key": "live2d-pet",
  "enabled": true,
  "position": "global-floating",
  "config": {
    "model": "/pet/models/hiyori/model3.json",
    "width": 280,
    "height": 380,
    "x": "right",
    "y": "bottom",
    "showOnMobile": false,
    "loadMode": "idle"
  }
}
```

### 12.4 模型文件目录

```text
apps/blog-web/public/pet/models/
├─ hiyori/
│  ├─ hiyori.model3.json
│  ├─ hiyori.moc3
│  ├─ textures/
│  └─ motions/
```

### 12.5 性能注意事项

1. 桌宠模型资源可能较大。
2. 不要首屏立即加载。
3. 移动端默认关闭。
4. 支持后台一键关闭。
5. 模型文件应放在静态目录。
6. 不要影响文章阅读。
7. 加载失败不能影响页面主体。

### 12.6 替代方案

如果 Live2D 太重，可以先做轻量桌宠：

1. GIF 动图桌宠。
2. Spine 动画。
3. Lottie 动画。
4. 简单 CSS 小宠物。

第一版建议：

```text
先做轻量桌宠开关，再升级 Live2D。
```

---

## 13. 构建时数据生成设计

### 13.1 为什么要构建时生成数据

博客前台应尽量保持静态。

因此：

```text
能构建前生成 JSON 的数据，就不要运行时请求数据库。
```

### 13.2 生成内容

构建前由 Go 后台或脚本生成：

```text
src/config/site.generated.ts
src/config/theme.generated.ts
src/config/feature.generated.ts
public/data/writing-heatmap.json
public/data/site-stat.json
```

### 13.3 生成流程

```text
SQLite site_setting
        |
        v
site.generated.ts

SQLite theme_setting
        |
        v
theme.generated.ts

SQLite feature_setting
        |
        v
feature.generated.ts

SQLite writing_activity / Markdown posts
        |
        v
writing-heatmap.json
```

### 13.4 feature.generated.ts 示例

```ts
export const features = {
  "writing-heatmap": {
    enabled: true,
    position: "home-sidebar",
    component: "Heatmap",
    config: {
      showWords: true,
      showPublishCount: true
    }
  },

  "live2d-pet": {
    enabled: true,
    position: "global-floating",
    component: "PetIsland",
    client: "idle",
    config: {
      model: "/pet/models/hiyori/hiyori.model3.json",
      showOnMobile: false,
      width: 280,
      height: 380
    }
  }
};
```

### 13.5 site-stat.json 示例

```json
{
  "postCount": 42,
  "categoryCount": 8,
  "tagCount": 36,
  "wordCount": 120000,
  "lastUpdated": "2026-06-24"
}
```

---

## 14. 后台管理系统设计

### 14.1 后台页面

```text
/admin/login
/admin/dashboard
/admin/articles
/admin/articles/new
/admin/articles/edit/:id
/admin/uploads
/admin/settings
/admin/themes
/admin/features
/admin/builds
/admin/backups
```

### 14.2 后台功能

核心：

1. 登录。
2. 仪表盘。
3. 文章草稿 CRUD。
4. Markdown 编辑器。
5. 图片上传。
6. 分类管理。
7. 标签管理。
8. 站点配置。
9. 主题配置。
10. Feature 配置。
11. 写作热力图统计。
12. 构建发布。
13. 构建日志。
14. 备份导出。
15. 备份导入。
16. 发布回滚。

### 14.3 仪表盘展示

```text
文章数量
草稿数量
总字数
本月写作字数
最近发布文章
最近构建状态
写作热力图预览
功能开关状态
```

### 14.4 Feature 管理页

功能列表：

| 功能 | 状态 | 挂载位置 | 操作 |
|---|---|---|---|
| 写作热力图 | 开启 | home-sidebar | 配置 / 关闭 |
| 桌宠 | 开启 | global-floating | 配置 / 关闭 |
| 评论 | 关闭 | post-after-content | 配置 / 开启 |
| 音乐播放器 | 关闭 | global-floating | 配置 / 开启 |
| 访问统计 | 关闭 | global-footer | 配置 / 开启 |

---

## 15. 后台 API 设计

### 15.1 认证接口

```http
POST /api/admin/login
POST /api/admin/logout
GET  /api/admin/profile
```

### 15.2 文章接口

```http
GET    /api/admin/articles
GET    /api/admin/articles/:id
POST   /api/admin/articles
PUT    /api/admin/articles/:id
DELETE /api/admin/articles/:id
POST   /api/admin/articles/:id/publish
POST   /api/admin/articles/:id/unpublish
```

### 15.3 上传接口

```http
POST   /api/admin/uploads
GET    /api/admin/uploads
DELETE /api/admin/uploads/:id
```

### 15.4 设置接口

```http
GET /api/admin/settings
PUT /api/admin/settings
```

### 15.5 主题接口

```http
GET /api/admin/themes
PUT /api/admin/themes/active
```

### 15.6 Feature 接口

```http
GET /api/admin/features
GET /api/admin/features/:key
PUT /api/admin/features/:key
POST /api/admin/features/:key/enable
POST /api/admin/features/:key/disable
```

### 15.7 构建接口

```http
POST /api/admin/build
GET  /api/admin/builds
GET  /api/admin/builds/:id
POST /api/admin/rollback
```

### 15.8 备份接口

```http
POST /api/admin/backups
GET  /api/admin/backups
POST /api/admin/backups/import
GET  /api/admin/backups/:id/download
```

---

## 16. SQLite 数据库设计

### 16.1 admin_user

```sql
CREATE TABLE admin_user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    nickname TEXT,
    avatar TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 16.2 article_draft

```sql
CREATE TABLE article_draft (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    description TEXT,
    content TEXT NOT NULL,
    cover TEXT,
    category TEXT,
    tags_json TEXT,
    draft INTEGER DEFAULT 1,
    pinned INTEGER DEFAULT 0,
    comment_enabled INTEGER DEFAULT 1,
    toc_enabled INTEGER DEFAULT 1,
    content_version INTEGER DEFAULT 1,
    word_count INTEGER DEFAULT 0,
    published_at DATETIME,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 16.3 article_publish_record

```sql
CREATE TABLE article_publish_record (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    draft_id INTEGER,
    title TEXT NOT NULL,
    slug TEXT NOT NULL,
    markdown_path TEXT NOT NULL,
    published_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    build_record_id INTEGER,
    FOREIGN KEY (draft_id) REFERENCES article_draft(id)
);
```

### 16.4 upload_file

```sql
CREATE TABLE upload_file (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    filename TEXT NOT NULL,
    original_name TEXT,
    path TEXT NOT NULL,
    url TEXT NOT NULL,
    mime_type TEXT,
    size INTEGER,
    hash TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 16.5 site_setting

```sql
CREATE TABLE site_setting (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    setting_key TEXT NOT NULL UNIQUE,
    setting_value TEXT NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 16.6 theme_setting

```sql
CREATE TABLE theme_setting (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    theme_key TEXT NOT NULL UNIQUE,
    enabled INTEGER DEFAULT 1,
    config_json TEXT NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 16.7 feature_setting

```sql
CREATE TABLE feature_setting (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    feature_key TEXT NOT NULL UNIQUE,
    enabled INTEGER DEFAULT 1,
    position TEXT,
    config_json TEXT,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 16.8 writing_activity

```sql
CREATE TABLE writing_activity (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    article_id INTEGER,
    article_slug TEXT,
    action TEXT NOT NULL,
    word_count INTEGER DEFAULT 0,
    delta_word_count INTEGER DEFAULT 0,
    activity_date DATE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 16.9 build_record

```sql
CREATE TABLE build_record (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    status TEXT NOT NULL,
    trigger_type TEXT,
    log TEXT,
    release_path TEXT,
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    finished_at DATETIME
);
```

### 16.10 backup_record

```sql
CREATE TABLE backup_record (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    filename TEXT NOT NULL,
    path TEXT NOT NULL,
    size INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

---

## 17. 发布构建流程

### 17.1 完整发布流程

```text
点击发布
   |
   v
保存文章草稿
   |
   v
记录 writing_activity
   |
   v
生成 Markdown 文件
   |
   v
生成 site.generated.ts
   |
   v
生成 theme.generated.ts
   |
   v
生成 feature.generated.ts
   |
   v
生成 writing-heatmap.json
   |
   v
执行 pnpm build
   |
   v
构建成功
   |
   v
复制 dist 到 releases/yyyyMMdd_HHmmss
   |
   v
切换 current 软链接
   |
   v
记录 build_record
```

### 17.2 build.sh

```bash
#!/usr/bin/env bash

set -e

ROOT_DIR="/opt/blog"
APP_DIR="$ROOT_DIR/apps/blog-web"
RELEASES_DIR="$ROOT_DIR/releases"
CURRENT_LINK="$ROOT_DIR/current"
TIME=$(date +%Y%m%d_%H%M%S)
RELEASE_DIR="$RELEASES_DIR/$TIME"

cd "$APP_DIR"

pnpm install --frozen-lockfile
pnpm build

mkdir -p "$RELEASE_DIR"
cp -r dist/* "$RELEASE_DIR"

ln -sfn "$RELEASE_DIR" "$CURRENT_LINK"

echo "Deploy success: $RELEASE_DIR"
```

### 17.3 release 机制优势

1. 构建失败不影响线上。
2. 发布成功才切换 current。
3. 支持快速回滚。
4. 每次构建可追踪。
5. 可以保留最近多个历史版本。

---

## 18. Nginx 配置

```nginx
server {
    listen 80;
    server_name your-domain.com;

    root /opt/blog/current;
    index index.html;

    location / {
        try_files $uri $uri/ /404.html;
    }

    location /admin/ {
        alias /opt/blog/apps/admin-web/dist/;
        try_files $uri $uri/ /admin/index.html;
    }

    location /admin-api/ {
        proxy_pass http://127.0.0.1:8080/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    location /uploads/ {
        alias /opt/blog/apps/blog-web/public/uploads/;
    }

    location /data/ {
        alias /opt/blog/apps/blog-web/public/data/;
    }
}
```

---

## 19. 备份与迁移设计

### 19.1 需要备份的内容

必须备份：

```text
apps/blog-web/src/content/posts/
apps/blog-web/public/uploads/
apps/blog-web/public/pet/
apps/blog-web/src/config/*.generated.ts
apps/blog-web/public/data/
data/blog_admin.db
```

可选备份：

```text
releases/
storage/logs/
```

不需要备份：

```text
node_modules/
dist/
临时构建文件
```

### 19.2 备份包结构

```text
backup-2026-06-24.zip
├─ posts/
├─ uploads/
├─ pet/
├─ data/
│  ├─ writing-heatmap.json
│  └─ site-stat.json
├─ database/
│  └─ blog_admin.db
├─ settings.json
├─ themes.json
├─ features.json
├─ drafts.json
└─ manifest.json
```

### 19.3 manifest.json

```json
{
  "app": "InkSpace",
  "version": "2.0.0",
  "contentVersion": 1,
  "exportedAt": "2026-06-24T12:00:00+09:00",
  "postCount": 25,
  "draftCount": 3,
  "uploadCount": 40,
  "featureCount": 5
}
```

### 19.4 迁移策略

升级系统时执行：

```text
1. 创建完整备份。
2. 停止后台写入。
3. 拉取新代码。
4. 执行数据库 migration。
5. 执行内容 migration。
6. 重新生成 generated 配置。
7. 重新生成 public/data。
8. pnpm build。
9. 构建成功后切换 current。
10. 检查首页、文章、图片、热力图、桌宠、RSS、搜索。
```

### 19.5 内容迁移脚本

目录：

```text
scripts/migrations/
├─ 001_add_content_version.ts
├─ 002_published_to_date.ts
├─ 003_normalize_tags.ts
├─ 004_fix_upload_paths.ts
├─ 005_generate_heatmap_data.ts
└─ 006_feature_config_migration.ts
```

---

## 20. 安全设计

### 20.1 后台登录

1. 密码使用 bcrypt。
2. JWT 设置过期时间。
3. 登录失败次数限制。
4. 后台接口统一鉴权。
5. 开启 HTTPS。
6. 不暴露敏感错误信息。

### 20.2 文件上传

1. 限制文件大小。
2. 限制文件类型。
3. 重命名文件。
4. 禁止上传可执行文件。
5. SVG 谨慎处理。
6. 上传目录不允许执行脚本。
7. 记录文件 hash。

### 20.3 构建命令

正确：

```text
后端只执行固定脚本 scripts/build.sh
```

错误：

```text
前端传 command，后端执行 command
```

后者会导致命令注入风险。

### 20.4 Feature 安全

1. 第三方脚本要可关闭。
2. 桌宠模型路径要限制在 public/pet 目录。
3. 不允许后台输入任意远程脚本 URL。
4. 评论和统计类 Feature 要防 XSS。
5. HTML 插入必须做过滤。

---

## 21. 性能设计

### 21.1 服务器资源规划

服务器：2 核 2G。

常驻进程：

| 进程 | 预计内存 |
|---|---|
| Nginx | 20MB - 50MB |
| Go Admin API | 30MB - 100MB |
| SQLite | 不单独常驻 |
| Node/pnpm | 只在构建时运行 |

### 21.2 建议

1. 开启 1G swap。
2. 不要频繁自动构建。
3. 发布按钮手动触发。
4. 桌宠默认移动端关闭。
5. 桌宠使用 idle 加载。
6. 热力图数据使用静态 JSON。
7. 搜索使用 Pagefind。
8. 只保留最近 5-10 个 release。
9. 上传图片压缩。
10. 不要第一版上 MySQL、Redis、Java 服务。

---

## 22. 开发阶段规划

### 22.1 第一阶段：Astro 前台 MVP

目标：手写 Markdown，前台能跑。

功能：

1. 首页。
2. 文章详情。
3. 分类。
4. 标签。
5. 归档。
6. 关于页。
7. Markdown 渲染。
8. 代码高亮。
9. 明暗主题。
10. RSS。
11. Sitemap。
12. 基础响应式。

完成标准：

```text
手写 3 篇文章后，博客可以正常构建和访问。
```

### 22.2 第二阶段：Go 后台 MVP

目标：可以在线写文章。

功能：

1. 管理员登录。
2. 文章草稿 CRUD。
3. Markdown 编辑器。
4. 图片上传。
5. 发布为 Markdown。
6. 触发 Astro build。
7. 查看构建日志。

完成标准：

```text
不用手动写 Markdown，可以在后台写文章并发布。
```

### 22.3 第三阶段：主题系统

目标：后台可配置主题。

功能：

1. 明暗切换。
2. 主题色切换。
3. 首页布局切换。
4. 文章页布局切换。
5. Banner 配置。
6. 导航配置。
7. 社交链接配置。

完成标准：

```text
后台切换主题配置后，重新构建即可改变博客风格。
```

### 22.4 第四阶段：Feature 系统

目标：建立可扩展插件能力。

功能：

1. FeatureHost。
2. Feature slots。
3. feature_setting 表。
4. feature.generated.ts。
5. 后台 Feature 管理页面。
6. 写作热力图 Feature。
7. 桌宠 Feature。

完成标准：

```text
写作热力图和桌宠都可以通过后台开启、关闭和配置。
```

### 22.5 第五阶段：增强功能

后续可以做：

1. Pagefind 搜索。
2. Giscus 评论。
3. 音乐播放器。
4. 访问统计。
5. AI 摘要。
6. AI 自动打标签。
7. Git 自动提交。
8. 备份导入导出。
9. release 回滚。
10. 内容迁移脚本。

---

## 23. 扩展性边界

本方案能保证：

```text
1. 新功能不污染文章系统。
2. 新功能不污染主题系统。
3. 新功能可以后台开关。
4. 新功能可以选择挂载位置。
5. 交互功能可以懒加载。
6. 静态数据可以构建时生成。
7. 文章内容方便迁移。
8. 小服务器可以承受。
```

本方案不能保证：

```text
1. 所有第三方功能都零代码接入。
2. 所有主题都自动适配所有 Feature。
3. 所有桌宠模型都无兼容问题。
4. 所有功能都完全静态化。
5. 未来所有需求都不需要改架构。
```

合理理解：

```text
它保证的是低耦合、可插拔、可关闭、可迁移。
不是保证任何功能都能无脑接入。
```

---

## 24. 最终方案总结

最终架构：

```text
Astro Blog Core
├─ Content Core
│  ├─ posts
│  ├─ categories
│  ├─ tags
│  └─ archives
│
├─ Theme Core
│  ├─ CSS variables
│  ├─ theme presets
│  └─ layout components
│
├─ Feature Core
│  ├─ feature registry
│  ├─ feature slots
│  ├─ feature config
│  ├─ feature generated data
│  ├─ writing heatmap
│  └─ live2d pet
│
└─ Build Core
   ├─ generate markdown
   ├─ generate site config
   ├─ generate theme config
   ├─ generate feature config
   ├─ generate heatmap data
   └─ pnpm build
```

最终技术定型：

```text
前台：
Astro + TypeScript + Tailwind CSS + Content Collections

后台页面：
Vue 3 + Vite + TypeScript + Element Plus

后台服务：
Go + Gin

数据库：
SQLite

内容源：
Markdown + Frontmatter

图片：
public/uploads/yyyy/mm/

功能扩展：
Feature Registry + FeatureHost + Feature Slots

发布：
生成 Markdown -> 生成配置与数据 -> pnpm build -> releases -> current

部署：
Nginx + systemd

迁移：
posts + uploads + database + generated config + migration scripts
```

一句话总结：

> 这个博客系统的核心不是“做一个能发文章的网站”，而是做一个 **长期可维护、可迁移、可主题化、可挂插件的个人内容系统**。  
> Astro 负责静态展示，Go 后台负责管理和构建，SQLite 负责后台工作台，Markdown 和 uploads 才是最重要的长期资产。
