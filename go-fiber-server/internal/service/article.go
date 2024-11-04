package service

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/data"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	"go-fiber-ent-web-layout/pkg/pool"
	"log/slog"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	topArticleCacheKey        = "BLOG:article:list:top"
	hotArticleCacheKey        = "BLOG:article:list:hot"
	archivesArticleCacheKey   = "BLOG:article:archives"
	articleInfoCacheKeyPrefix = "BLOG:article:info:"
)

type ArticleService struct {
	repo          usercase.IArticleRepo
	redisTemplate *data.RedisTemplate
}

func NewArticleService(repo usercase.IArticleRepo, redisTemplate *data.RedisTemplate) usercase.IArticleService {
	return &ArticleService{
		repo:          repo,
		redisTemplate: redisTemplate,
	}
}

func (self *ArticleService) SaveArticle(article *usercase.Article) error {
	if strings.TrimSpace(article.Content) == "" {
		return tools.FiberRequestError("文章内容不能为空")
	}
	if err := self.repo.Save(article); err != nil {
		slog.Error("保存博客文章失败", "error", err.Error())
		return tools.FiberServerError("保存失败")
	}
	// 如果文章状态为正常 那么需要删除归档数据缓存
	if article.Status != nil && *article.Status == 1 {
		self.deleteRedisArticleArchives()
	}
	// 如果是置顶文章 那么删除redis的置顶文章缓存
	if article.IsTop {
		self.deleteRedisTopArticle()
	}
	// 如果是热门文章 删除redis的热门文章缓存
	if article.IsHot {
		self.deleteRedisHotArticle()
	}
	// TODO 保存到redis 保存到 elasticsearch
	return nil
}

func (self *ArticleService) UpdateArticle(article *usercase.Article) error {
	if article.ArticleId == 0 {
		return tools.FiberRequestError("文章ID不能为空")
	}
	if err := self.repo.Update(article); err != nil {
		slog.Error("更新博客文章失败", "error", err.Error())
		return tools.FiberServerError("更新失败")
	}
	// 清除文章的redis缓存
	self.deleteRedisArticleArchives()
	self.deleteRedisHotArticle()
	self.deleteRedisTopArticle()
	return nil
}

func (self *ArticleService) UpdateSelectiveArticle(form *usercase.ArticleUpdateForm) error {
	if form.ArticleId == 0 {
		return tools.FiberRequestError("文章ID不能为空")
	}
	if err := self.repo.UpdateSelective(form); err != nil {
		slog.Error("快捷更新博客文章失败", "error", err.Error())
		return tools.FiberServerError("更新失败")
	}
	// 清除文章的redis缓存
	self.deleteRedisArticleArchives()
	self.deleteRedisHotArticle()
	self.deleteRedisTopArticle()
	return nil
}

func (self *ArticleService) Page(query *usercase.ArticleQueryForm) (*usercase.PageData[usercase.ArticleVo], error) {
	articles, total, err := self.repo.Page(query)
	if err != nil {
		slog.Error("分页查询博客文章失败", "error", err.Error())
		return nil, tools.FiberServerError("查询失败")
	}
	return usercase.NewPageData(articles, total, query.Page, query.Size), nil
}

func (self *ArticleService) ListTopArticle() ([]*usercase.Article, error) {
	topList, err := data.RedisGetSlice[*usercase.Article](context.Background(), topArticleCacheKey)
	if err == nil && len(topList) > 0 {
		return topList, nil
	}
	topList, err = self.repo.ListTopArticle()
	if err != nil {
		slog.Error("查询置顶博客文章失败", "err", err.Error())
		return nil, tools.FiberServerError("查询置顶文章失败")
	}
	pool.Go(func() {
		if setErr := self.redisTemplate.Set(context.Background(), topArticleCacheKey, topList, math.MaxInt64); setErr != nil {
			slog.Error("置顶文章添加redis缓存失败", "error", err)
		}
	})
	return topList, nil
}

func (self *ArticleService) ListHotArticle() ([]usercase.SimpleArticleVo, error) {
	hots, err := data.RedisGetSlice[usercase.SimpleArticleVo](context.Background(), hotArticleCacheKey)
	if err == nil && len(hots) > 0 {
		return hots, nil
	}
	hots, err = self.repo.ListHotArticle()
	if err != nil {
		slog.Error("查询热门文章失败", "error", err)
		return nil, tools.FiberServerError("查询热门文章失败")
	}
	pool.Go(func() {
		// 缓存30分钟之内有效
		if setErr := self.redisTemplate.Set(context.Background(), hotArticleCacheKey, hots, 30*time.Minute); setErr != nil {
			slog.Error("热门文章添加reids缓存失败", "error", setErr)
		}
	})
	return hots, nil
}

func (self *ArticleService) PageByLabel(query *usercase.ArticleQueryForm) (*usercase.PageData[usercase.Article], error) {
	articles, total, err := self.repo.PageByLabel(query)
	if err != nil {
		slog.Error("查询分类,标签关联文章失败", "error", err.Error())
		return nil, tools.FiberServerError("查询失败")
	}
	return usercase.NewPageData(articles, total, query.Page, query.Size), nil
}

func (self *ArticleService) Archives() ([]usercase.ArticleArchive, error) {
	archives, err := data.RedisGetSlice[usercase.ArticleArchive](context.Background(), archivesArticleCacheKey)
	if err == nil && len(archives) > 0 {
		return archives, nil
	}
	archives, err = self.repo.Archives()
	if err != nil {
		slog.Error("获取文章归档信息失败", "err", err.Error())
		return nil, tools.FiberServerError("查询归档信息失败")
	}
	pool.Go(func() {
		if setErr := self.redisTemplate.Set(context.Background(), archivesArticleCacheKey, archives, math.MaxInt64); setErr != nil {
			slog.Error("文章归档数据添加redis缓存失败", "error", setErr)
		}
	})
	return archives, nil
}

func (self *ArticleService) SelectById(articleId uint64, isAdmin bool) (*usercase.ArticleVo, error) {
	article, err := self.repo.SelectById(articleId, isAdmin)
	if err != nil {
		slog.Error("查询博客文章详情失败", "error", err.Error())
		return nil, err
	}
	if article == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "文章不存在")
	}
	return article, err
}

func (self *ArticleService) DeleteArticleById(articleId uint64) error {
	if err := self.repo.DeleteById(articleId); err != nil {
		slog.Error("删除博客文章失败", "error", err.Error())
		return tools.FiberServerError("删除失败")
	}
	return nil
}

func (self *ArticleService) ArticleVoteUp(articleId uint64) error {
	if err := self.repo.VoteUp(articleId, 1); err != nil {
		slog.Error("更新文章点赞数失败", "error", err)
		return tools.FiberServerError("点赞失败")
	}
	return nil
}

func (self *ArticleService) SearchArticle(keyword string) ([]usercase.SimpleArticleVo, error) {
	result, err := self.repo.Search(keyword, 10)
	if err != nil {
		slog.Error("搜索文章错误", "error", err, "keyword", keyword)
		return nil, tools.FiberServerError("搜索异常")
	}
	for i := 0; i < len(result); i++ {
		item := &result[i]
		if index := strings.Index(item.Title, keyword); index >= 0 {
			item.Title = strings.ReplaceAll(item.Title, keyword, "<span>"+keyword+"</span>")
		}
		if index := strings.Index(item.Summary, keyword); index >= 0 {
			item.Summary = strings.ReplaceAll(item.Summary, keyword, "<span>"+keyword+"</span>")
		}
	}
	return result, nil
}

func (self *ArticleService) getArticleInfoCacheKey(articleId uint64) string {
	return articleInfoCacheKeyPrefix + strconv.FormatUint(articleId, 10)
}

func (self *ArticleService) deleteRedisHotArticle() {
	pool.Go(func() {
		if err := self.redisTemplate.Delete(context.Background(), hotArticleCacheKey); err != nil {
			slog.Error("删除redis热门文章缓存失败", "error", err)
		}
	})
}

func (self *ArticleService) deleteRedisTopArticle() {
	pool.Go(func() {
		if err := self.redisTemplate.Delete(context.Background(), topArticleCacheKey); err != nil {
			slog.Error("删除redis置顶文章缓存失败", "error", err)
		}
	})
}

func (self *ArticleService) deleteRedisArticleArchives() {
	pool.Go(func() {
		if err := self.redisTemplate.Delete(context.Background(), archivesArticleCacheKey); err != nil {
			slog.Error("删除redis文章归档缓存失败", "error", err)
		}
	})
}
