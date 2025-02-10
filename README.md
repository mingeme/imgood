# Imgood

一个基于SvelteKit和Supabase构建的自托管图床服务。

## 技术栈

- **前端框架**: SvelteKit 2.x + Svelte 5
- **数据库**: Supabase
- **存储**: AWS S3
- **Node版本**: 20.x
- **包管理器**: pnpm

## 环境要求

1. Node.js 20.x (可通过项目根目录的`.nvmrc`自动选择正确版本)
2. pnpm 包管理器
3. Supabase项目和配置
4. AWS S3存储配置

## 开发环境设置

1. 安装依赖：

```bash
pnpm install
```

2. 配置环境变量：

复制`.env.example`文件并重命名为`.env`，然后填入必要的环境变量：

```bash
cp .env.example .env
```

3. 启动开发服务器：

```bash
pnpm dev
```

或者自动打开浏览器：

```bash
pnpm dev -- --open
```

## 部署

项目使用Vercel进行部署。确保在部署环境中设置了所有必要的环境变量。

## 计划列表
- [x] 基础图片上传功能
- [x] 图片预览和分享链接
- [ ] 基本的用户认证
- [ ] 图片搜索功能
- [ ] 图片压缩和优化
- [ ] 自定义图片处理选项
- [ ] 多语言支持
- [ ] API接口

## 项目灵感

sm.ms 自托管实现
