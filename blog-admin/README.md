# 博客系统管理后台前端

技术栈`Vue3` + `Pinia` + `Typescript` + `ArcoDesign`。手动实现了大部分开源快速开发框架的后台管理功能，动态菜单、动态权限、路由缓存、角色校验、请求缓存、国际化等。

采用模块化开发的方式，管理模块和业务模块解耦分离，代码目录结构清晰、可读性较强、有完善的类型提示，部署成功之后可以快速上手！

`Markdown`编辑器使用开源的`Vditor`编辑器，支持所见即所得模式，除了`GFM`标准之外还支持数学公式等等等。

[Preview](https://fiber-blog-admin.vercel.app/) `admin@admin123456`

## 待实现的功能

- [x] 个人中心
- [x] 评论管理
- [x] 博客端用户管理
- [ ] 实时日志
- [x] 首页数据统计
- [ ] 一些小优化...

## 快速开始

```bash
pnpm install
# dev模式
pnpm run dev
# 编译
pnpm run build
```