/**
 * 统一数据返回对象
 */
export interface Result<T> {
  code: number;
  message: string;
  timestamp: number;
  data: T,
}

/**
 * 分页数据
 */
export interface Page<T> {
  // 当前页码
  current: number;
  // 总页数
  pages: number;
  // 总记录数
  total: number;
  // 每页记录数
  size: number;
  // 记录列表
  records: T[];
}

/**
 * 博客标签数据
 */
export interface Tag {
  tagId: number;
  tagName: string;
  coverUrl: string;
  viewNum: number;
  color: string;
  articleNum: number;
  createTime: string;
}

/**
 * 博客分类数据
 */
export interface Category {
  categoryId: number;
  categoryName: string;
  description: string;
  coverUrl: string;
  viewNum: number;
  parentId: number;
  isHot: boolean;
  isTop: boolean;
  articleNum: number;
  createTime: string;
  children?: Category[]
}

/**
 * 联系方式
 */
export interface Concat {
  concatId: number;
  name: string;
  iconSvg: string;
  targetUrl: string;
  isMain: boolean;
}

/**
 * 通知公告数据
 */
export interface Notice {
  noticeId: number;
  title: string;
  message: string;
  level: number; // 通知级别 1:info 2:warn 3:error
  noticeType: number;
}

/**
 * 文章关联分类
 */
export interface ArticleCategory {
  categoryId: number;
  categoryName: string;
}

/**
 * 文章关联标签
 */
export interface ArticleTag {
  tagId: number;
  tagName: string;
  color: string;
}

/**
 * 博客文章数据
 */
export interface Article {
  articleId: number;
  title: string;
  summary: string;
  coverUrl: string;
  categoryIds: number[];
  tagIds: number[];
  viewNum: number;
  shareNum: number;
  wordCount: number;
  voteUp: number;
  content: string;
  protocol: string;
  tips: string;
  password: string;
  isHot: boolean;
  isTop: boolean;
  isComment: boolean;
  isPrivate: boolean;
  createTime: string;
  updateTime: string;
  sort: number;
  status: number;
  // 评论数量
  commentNum: number;
  // 分类列表
  categories: ArticleCategory[];
  // 标签列表
  tags: ArticleTag[];
}

/**
 * 文章归档信息
 */
export interface ArticleArchive {
  month: string;
  total: number;
}

/**
 * 友情链接数据
 */
export interface Link {
  linkId: number;
  name: string;
  summary: string;
  coverUrl: string;
  targetUrl: string;
  clickNum: number;
}

/**
 * 站点统计数据
 */
export interface SiteStats {
  articleCount: number;
  categoryCount: number;
  tagCount: number;
  commentCount: number;
  visitorCount: number;
  accessCount: number;
  wordTotal: number;
}

/**
 * 站点动态数据
 */
export interface Topic {
  topicId: number;
  content: string;
  imageUrls: string[];
  location: string;
  isHot: boolean;
  isTop: boolean;
  voteUp: number;
  mode: number;
  createTime: string
}

/**
 * 博客端用户数据
 */
export interface User {
  userId: number;
  nickname: string;
  summary: string;
  avatar: string;
  email: string;
  link: string;
  username: string;
  labels: string[];
  level: number;
  expertise: number;
  registerIp: string;
  registerLocation: string;
}

/**
 * 评论数据
 */
export interface Comment {
  commentId: number;
  content: string;
  userId: number;
  fid: number;
  rid: number;
  location: string;
  commentUa: string;
  voteUp: number;
  commentType: number;
  isHot: boolean;
  isTop: boolean;
  isColl: boolean;
  createTime: string;
  user: User;
  parentUser?: User;
  children: Page<Comment>
}