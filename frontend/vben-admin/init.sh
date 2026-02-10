#!/bin/bash
#
# Vben Admin 管理后台一键初始化脚本
# 用法: bash init.sh [项目名称]
# 示例: bash init.sh
#

set -e

# ========== 参数 ==========
PROJECT_NAME="${1:-vben-admin}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# ========== 颜色 ==========
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

info()    { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[OK]${NC} $1"; }
warn()    { echo -e "${YELLOW}[WARN]${NC} $1"; }
error()   { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# ========== 环境检查 ==========
info "检查 Node.js 环境..."
if ! command -v node &> /dev/null; then
    error "未安装 Node.js，请先安装: https://nodejs.org/"
fi
NODE_VERSION=$(node --version)
success "Node.js: $NODE_VERSION"

if ! command -v pnpm &> /dev/null; then
    warn "未安装 pnpm，正在安装..."
    npm install -g pnpm || error "pnpm 安装失败"
fi
PNPM_VERSION=$(pnpm --version)
success "pnpm: $PNPM_VERSION"

# 检查 git
if ! command -v git &> /dev/null; then
    error "未安装 Git"
fi

# ========== 选择克隆源 ==========
echo ""
echo "========================================"
echo "  Vben Admin 管理后台初始化"
echo "========================================"
echo ""
echo "  请选择克隆源："
echo "  1) GitHub（推荐，最新代码）"
echo "  2) Gitee 镜像（国内加速，可能不是最新）"
echo ""
read -p "请选择 (1/2): " -n 1 -r CLONE_SOURCE
echo ""

case $CLONE_SOURCE in
    2)  REPO_URL="https://gitee.com/annsion/vue-vben-admin.git" ;;
    *)  REPO_URL="https://github.com/vbenjs/vue-vben-admin.git" ;;
esac

# ========== 克隆项目到临时目录 ==========
TEMP_DIR=$(mktemp -d)

info "正在克隆 $REPO_URL ..."
git clone --depth 1 "$REPO_URL" "$TEMP_DIR/vben" || {
    rm -rf "$TEMP_DIR"
    error "克隆失败，请检查网络连接"
}

# 将文件合并到目标目录（保留 CLAUDE.md 和 init.sh）
info "合并文件到目标目录..."
rm -rf "$TEMP_DIR/vben/.git"
cd "$TEMP_DIR/vben"
for item in * .[!.]*; do
    [ ! -e "$item" ] && continue
    [ "$item" = "CLAUDE.md" ] || [ "$item" = "init.sh" ] && continue
    cp -r "$item" "$SCRIPT_DIR/"
done

# 清理临时目录
rm -rf "$TEMP_DIR"

cd "$SCRIPT_DIR"
success "项目创建完成"

# ========== 查找主应用目录 ==========
# Vben5 monorepo 结构：apps/web-* 或 apps/backend
APP_DIR=""
if [ -d "apps/web-admin" ]; then
    APP_DIR="apps/web-admin"
elif [ -d "apps/web-antd" ]; then
    APP_DIR="apps/web-antd"
elif [ -d "apps/web-ele" ]; then
    APP_DIR="apps/web-ele"
elif [ -d "apps/web-naive" ]; then
    APP_DIR="apps/web-naive"
else
    # 尝试找第一个 apps/web-* 目录
    APP_DIR=$(find apps -maxdepth 1 -name "web-*" -type d 2>/dev/null | head -1)
fi

if [ -z "$APP_DIR" ] || [ ! -d "$APP_DIR" ]; then
    warn "未找到标准的 apps/web-* 目录，将在 src/ 下创建业务目录"
    APP_DIR="."
fi

info "主应用目录: $APP_DIR"

# ========== 创建业务目录结构 ==========
info "创建业务目录结构..."

SRC_DIR="$APP_DIR/src"

# API 目录
mkdir -p "$SRC_DIR/api/core"
mkdir -p "$SRC_DIR/api/system"
mkdir -p "$SRC_DIR/api/business"

# 页面目录
mkdir -p "$SRC_DIR/views/system/user"
mkdir -p "$SRC_DIR/views/system/role"
mkdir -p "$SRC_DIR/views/system/menu"
mkdir -p "$SRC_DIR/views/business"
mkdir -p "$SRC_DIR/views/dashboard"

success "业务目录创建完成"

# ========== 创建环境配置 ==========
info "创建环境配置文件..."

# 只在文件不存在时创建
if [ ! -f "$APP_DIR/.env.development" ]; then
    cat > "$APP_DIR/.env.development" << 'ENV_EOF'
# 开发环境配置
VITE_GLOB_API_URL=http://localhost:8282
VITE_GLOB_API_URL_PREFIX=/api/admin
ENV_EOF
    success "创建 .env.development"
else
    warn ".env.development 已存在，跳过"
fi

if [ ! -f "$APP_DIR/.env.production" ]; then
    cat > "$APP_DIR/.env.production" << 'ENV_EOF'
# 生产环境配置
VITE_GLOB_API_URL=https://api.yourcompany.com
VITE_GLOB_API_URL_PREFIX=/api/admin
ENV_EOF
    success "创建 .env.production"
else
    warn ".env.production 已存在，跳过"
fi

# ========== 创建 API 模板文件 ==========
info "创建 API 模板文件..."

# system/user.ts
if [ ! -f "$SRC_DIR/api/system/user.ts" ]; then
    cat > "$SRC_DIR/api/system/user.ts" << 'TS_EOF'
/**
 * 管理员用户接口
 * 对接后端: /api/admin/user/*
 */

// TODO: 根据项目实际使用的请求库替换 import
// import { request } from '../core/request';

enum Api {
  Login = '/user/login',
  Logout = '/user/logout',
  UserInfo = '/user/userinfo',
  ChangePassword = '/user/change-password',
}

/** 管理员登录 */
export interface LoginParams {
  username: string;
  password: string;
}

/** 用户信息 */
export interface UserInfo {
  id: number;
  username: string;
  nickname: string;
  roles: string[];
  permissions: string[];
}

// export function loginApi(params: LoginParams) {
//   return request.post({ url: Api.Login, data: params });
// }
//
// export function logoutApi() {
//   return request.post({ url: Api.Logout });
// }
//
// export function getUserInfoApi() {
//   return request.get({ url: Api.UserInfo });
// }
TS_EOF
    success "创建 api/system/user.ts"
fi

# system/role.ts
if [ ! -f "$SRC_DIR/api/system/role.ts" ]; then
    cat > "$SRC_DIR/api/system/role.ts" << 'TS_EOF'
/**
 * 角色管理接口
 * 对接后端: /api/admin/role/*
 */

enum Api {
  List = '/role/list',
  Create = '/role/create',
  Update = '/role/update',
  Delete = '/role/delete',
}

// TODO: 实现接口调用
TS_EOF
    success "创建 api/system/role.ts"
fi

# system/menu.ts
if [ ! -f "$SRC_DIR/api/system/menu.ts" ]; then
    cat > "$SRC_DIR/api/system/menu.ts" << 'TS_EOF'
/**
 * 菜单/权限管理接口
 * 对接后端: /api/admin/permission/*
 */

enum Api {
  List = '/permission/list',
  Create = '/permission/create',
  Update = '/permission/update',
  Delete = '/permission/delete',
}

// TODO: 实现接口调用
TS_EOF
    success "创建 api/system/menu.ts"
fi

# ========== 安装依赖 ==========
if [ -f "pnpm-workspace.yaml" ] || [ -f "package.json" ]; then
    info "安装依赖..."
    pnpm install || warn "依赖安装失败，请手动执行 pnpm install"
    success "依赖安装完成"
fi

# ========== 完成 ==========
echo ""
echo "========================================"
echo -e "  ${GREEN}Vben Admin 初始化完成！${NC}"
echo "========================================"
echo ""
echo "  项目路径: $SCRIPT_DIR"
echo "  主应用:   $APP_DIR"
echo ""
echo "  下一步:"
echo "  1. 修改 $APP_DIR/.env.development 中的 API 地址"
echo "  2. 修改 $APP_DIR/.env.production 中的 API 地址"
echo "  3. 根据 Vben 版本适配 API 请求封装"
echo "  4. 启动: pnpm dev"
echo ""
echo "  业务目录说明:"
echo "  $SRC_DIR/"
echo "  ├── api/"
echo "  │   ├── core/       请求封装"
echo "  │   ├── system/     系统管理 API（user/role/menu）"
echo "  │   └── business/   业务 API（按模块添加）"
echo "  └── views/"
echo "      ├── dashboard/  仪表盘"
echo "      ├── system/     系统管理页面"
echo "      └── business/   业务页面（按模块添加）"
echo ""
echo "========================================"
