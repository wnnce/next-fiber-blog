package service

import (
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/data"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	"log/slog"
	"strings"
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
	// TODO 清空redis缓存
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

func (self *ArticleService) SelectById(articleId uint64, checkStatus bool) (*usercase.ArticleVo, error) {
	article, err := self.repo.SelectById(articleId, checkStatus)
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
