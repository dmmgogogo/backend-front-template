# Vben Admin 模板清理指南

> 基于 vben-admin (web-naive) 模板，清理模板默认数据，精简为干净的后台管理系统基础框架。

---

## 一、品牌命名替换（6处）

### 1. `apps/web-naive/.env`

```diff
- VITE_APP_TITLE=Vben Admin Naive
- VITE_APP_NAMESPACE=vben-web-naive
+ VITE_APP_TITLE=你的项目名
+ VITE_APP_NAMESPACE=你的项目命名空间
```

### 2. `apps/web-naive/index.html`

meta 标签替换：

```diff
- <meta name="description" content="A Modern Back-end Management System" />
- <meta name="keywords" content="Vben Admin Vue3 Vite" />
- <meta name="author" content="Vben" />
+ <meta name="description" content="你的项目描述" />
+ <meta name="keywords" content="你的关键词" />
+ <meta name="author" content="你的作者名" />
```

删除百度统计脚本（整段 `<script>` 删除）：

```diff
-    <script>
-      // 生产环境下注入百度统计
-      if (window._VBEN_ADMIN_PRO_APP_CONF_) {
-        var _hmt = _hmt || [];
-        (function () {
-          var hm = document.createElement('script');
-          hm.src = 'https://hm.baidu.com/hm.js?24bb3eb91dfe4ebfcbcee6952a107cb6';
-          var s = document.getElementsByTagName('script')[0];
-          s.parentNode.insertBefore(hm, s);
-        })();
-      }
-    </script>
```

### 3. `packages/@core/preferences/src/config.ts`

```diff
- name: 'Vben Admin',
+ name: '你的项目名',

- defaultHomePath: '/analytics',
+ defaultHomePath: '/workspace',

  copyright: {
-   companyName: 'Vben',
-   companySiteLink: 'https://www.vben.pro',
-   date: '2024',
+   companyName: '你的公司/项目名',
+   companySiteLink: '',
+   date: '2025',
  },
```

### 4. `packages/effects/common-ui/src/ui/about/about.vue`

```diff
- name: 'Vben Admin',
+ name: '你的项目名',
```

### 5. `packages/effects/layouts/src/basic/copyright/copyright.vue`

```diff
- companyName: 'Vben Admin',
+ companyName: '你的项目名',
```

### 6. `internal/vite-config/src/utils/env.ts`

```diff
- appTitle: getString(VITE_APP_TITLE, 'Vben Admin'),
+ appTitle: getString(VITE_APP_TITLE, '你的项目名'),
```

---

## 二、国际化文本更新（2处）

### 7. `packages/locales/src/langs/zh-CN/authentication.json`

```diff
- "pageTitle": "开箱即用的大型中后台管理系统",
- "pageDesc": "工程化、高性能、跨组件库的前端模版",
- "loginSubtitle": "请输入您的帐户信息以开始管理您的项目",
+ "pageTitle": "你的系统名称",
+ "pageDesc": "你的系统描述",
+ "loginSubtitle": "请输入您的账户信息以登录系统",
```

### 8. `packages/locales/src/langs/en-US/authentication.json`

```diff
- "pageTitle": "Plug-and-play Admin system",
- "pageDesc": "Efficient, versatile frontend template",
- "loginSubtitle": "Enter your account details to manage your projects",
+ "pageTitle": "Your System Name",
+ "pageDesc": "Your system description",
+ "loginSubtitle": "Enter your account details to login",
```

---

## 三、登录页精简（2处）

### 9. `apps/web-naive/src/views/_core/authentication/login.vue`

**删除内容：**

- `MOCK_USER_OPTIONS` 数组（Super/Admin/User 快速选择账号）
- `VbenSelect` 表单项（selectAccount 下拉框）
- username 字段的 `dependencies` 联动逻辑（选账号自动填充用户名密码）
- `import type { BasicOption } from '@vben/types'` 引用

**只保留三个表单项：**

```typescript
formSchema = [
  { component: 'VbenInput',         fieldName: 'username' },  // 账号
  { component: 'VbenInputPassword', fieldName: 'password' },  // 密码
  { component: SliderCaptcha,       fieldName: 'captcha'  },  // 滑块验证
]
```

**传入 props 隐藏多余功能：**

```vue
<AuthenticationLogin
  :form-schema="formSchema"
  :loading="authStore.loginLoading"
  :show-code-login="false"
  :show-forget-password="false"
  :show-qrcode-login="false"
  :show-register="false"
  :show-third-party-login="false"
  @submit="authStore.authLogin"
/>
```

> 这会隐藏：手机号登录按钮、扫码登录按钮、第三方登录图标、忘记密码链接、创建账号链接

### 10. `apps/web-naive/src/router/routes/core.ts`

**删除** 4 个认证子路由，只保留 `Login`：

```diff
  children: [
    { name: 'Login', path: 'login', ... },           // ✅ 保留
-   { name: 'CodeLogin', path: 'code-login', ... },       // ❌ 删除
-   { name: 'QrCodeLogin', path: 'qrcode-login', ... },   // ❌ 删除
-   { name: 'ForgetPassword', path: 'forget-password', ... }, // ❌ 删除
-   { name: 'Register', path: 'register', ... },          // ❌ 删除
  ],
```

---

## 四、登录后页面精简（5处）

### 11. `apps/web-naive/src/preferences.ts`

添加默认首页路径覆盖：

```diff
  app: {
+   defaultHomePath: '/workspace',
    name: import.meta.env.VITE_APP_TITLE,
  },
```

### 12. `apps/web-naive/src/router/routes/modules/dashboard.ts`

**删除** Analytics（分析页）路由，只保留 Workspace（工作台），并给工作台加上 `affixTab: true`：

```diff
  children: [
-   { name: 'Analytics', path: '/analytics', meta: { affixTab: true, ... } },
    { name: 'Workspace', path: '/workspace', meta: { affixTab: true, ... } },
  ],
```

### 13. `apps/web-naive/src/views/dashboard/workspace/index.vue`

**只保留** 顶部 `WorkbenchHeader`，**删除**下方全部示例内容：

- `projectItems` 项目卡片数据 + `WorkbenchProject` 组件
- `quickNavItems` 快捷导航数据 + `WorkbenchQuickNav` 组件
- `todoItems` 待办事项数据 + `WorkbenchTodo` 组件
- `trendItems` 最新动态数据 + `WorkbenchTrends` 组件
- `AnalyticsVisitsSource` 访问来源图表
- `navTo` 导航方法
- 相关的所有 import

**精简后的完整文件：**

```vue
<script lang="ts" setup>
import { WorkbenchHeader } from '@vben/common-ui';
import { preferences } from '@vben/preferences';
import { useUserStore } from '@vben/stores';

const userStore = useUserStore();
</script>

<template>
  <div class="p-5">
    <WorkbenchHeader
      :avatar="userStore.userInfo?.avatar || preferences.app.defaultAvatar"
    >
      <template #title>
        早安, {{ userStore.userInfo?.realName }}, 开始您一天的工作吧！
      </template>
      <template #description> 今日晴，20℃ - 32℃！ </template>
    </WorkbenchHeader>
  </div>
</template>
```

### 14. `apps/web-naive/src/router/routes/modules/demos.ts`

**清空**路由数组（删除"演示"侧边栏菜单）：

```typescript
import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [];

export default routes;
```

### 15. `apps/web-naive/src/router/routes/modules/vben.ts`

**删除** VbenProject（"项目"菜单）和 VbenAbout（"关于"菜单），**只保留** Profile（个人资料页，隐藏菜单不显示）：

```typescript
import type { RouteRecordRaw } from 'vue-router';

import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    name: 'Profile',
    path: '/profile',
    component: () => import('#/views/_core/profile/index.vue'),
    meta: {
      icon: 'lucide:user',
      hideInMenu: true,
      title: $t('page.auth.profile'),
    },
  },
];

export default routes;
```

---

## 修改清单速查

| # | 文件 | 类别 | 操作 |
|---|------|------|------|
| 1 | `apps/web-naive/.env` | 品牌 | 改标题和命名空间 |
| 2 | `apps/web-naive/index.html` | 品牌 | 改 meta 标签，删百度统计 |
| 3 | `packages/@core/preferences/src/config.ts` | 品牌 | 改应用名、首页路径、版权 |
| 4 | `packages/effects/common-ui/src/ui/about/about.vue` | 品牌 | 改默认名 |
| 5 | `packages/effects/layouts/src/basic/copyright/copyright.vue` | 品牌 | 改默认公司名 |
| 6 | `internal/vite-config/src/utils/env.ts` | 品牌 | 改 fallback 标题 |
| 7 | `packages/locales/src/langs/zh-CN/authentication.json` | 国际化 | 改登录页中文文案 |
| 8 | `packages/locales/src/langs/en-US/authentication.json` | 国际化 | 改登录页英文文案 |
| 9 | `apps/web-naive/src/views/_core/authentication/login.vue` | 登录页 | 删 mock 选择器，传 props 隐藏多余功能 |
| 10 | `apps/web-naive/src/router/routes/core.ts` | 登录页 | 删多余认证路由 |
| 11 | `apps/web-naive/src/preferences.ts` | 页面 | 加 defaultHomePath |
| 12 | `apps/web-naive/src/router/routes/modules/dashboard.ts` | 页面 | 删分析页路由 |
| 13 | `apps/web-naive/src/views/dashboard/workspace/index.vue` | 页面 | 精简为只有 Header |
| 14 | `apps/web-naive/src/router/routes/modules/demos.ts` | 页面 | 清空演示路由 |
| 15 | `apps/web-naive/src/router/routes/modules/vben.ts` | 页面 | 删项目/关于，保留 Profile |
