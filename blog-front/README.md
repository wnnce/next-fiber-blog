# 博客前端工程
技术栈`Next.js` + `Typescript` + `UnoCSS`，基于`RSS` + `SPA`混合开发，具有完善的`SEO`优化，对于搜索引擎爬虫良好。

[Preview](https://fiber-blog-front-three.vercel.app/)

**还处于开发阶段，当前仅供预览使用!**

## 说明
不依赖任何`UI`库，手动实现了一些常用的`UI`组件，例如轮播、分页、时间轴等。`Markdown`解析使用[markdown-it]()库，拥有良好的性能和完善的生态系统。
性能方面，客户端组件遵循`React`的最佳实践，使用`React.memo`和`useMemo`针对组件和数据进行缓存，减少组件重绘的性能开销。

## 待开发功能
- [ ] 文章详情页面完善
- [ ] 页面加载loading页
- [ ] 404以及其它错误页
- [ ] 其它小功能...

## 快速开始
`node`版本大于等于`20`，`clone`本项目后使用`pnpm`安装依赖，而后访问[http://localhost:3000](http://localhost:3000)
```bash
pnpm install

pnpm run dev
```
