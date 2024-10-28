package service

import (
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/tools/region"
	"go-fiber-ent-web-layout/internal/usercase"
	"go-fiber-ent-web-layout/pkg/pool"
	"log/slog"
	"sync"
)

type CommentService struct {
	repo        usercase.ICommentRepo
	userService usercase.IUserService
}

func NewCommentService(repo usercase.ICommentRepo, userService usercase.IUserService) usercase.ICommentService {
	return &CommentService{
		repo:        repo,
		userService: userService,
	}
}

func (self *CommentService) SaveComment(comment *usercase.Comment) error {
	location := region.SearchLocation(comment.CommentIp)
	comment.Location = location
	if err := self.repo.Save(comment); err != nil {
		slog.Error("保存评论失败", "err", err.Error())
		return tools.FiberServerError("保存评论失败")
	}
	// 异步更新
	pool.Go(func() {
		_ = self.userService.UpdateUserExpertise(20, comment.UserId, 2)
	})
	return nil
}

func (self *CommentService) Page(query *usercase.CommentQueryForm) (*usercase.PageData[usercase.CommentVo], error) {
	page, err := self.repo.Page(query)
	if err != nil {
		slog.Error("分页查询评论列表失败", "err", err.Error(), "topicId", query.TopicId, "articleId", query.ArticleId, "commentType", query.CommentType)
		return nil, tools.FiberServerError("获取评论列表失败")
	}
	wg := sync.WaitGroup{}
	// 查询所有一级评论的子评论
	for i := 0; i < len(page.Records); i++ {
		wg.Add(1)
		record := page.Records[i]
		pool.Go(func() {
			childrenQuery := &usercase.CommentQueryForm{
				PageQueryForm: usercase.PageQueryForm{
					Page: 1,
					Size: 10,
				},
				Fid: record.CommentId,
			}
			childrenPage, childrenErr := self.repo.Page(childrenQuery)
			if childrenErr != nil {
				slog.Error("获取一级评论的子评论失败", "err", childrenErr.Error(), "fid", query.Fid)
				return
			}
			record.Children = childrenPage
			wg.Done()
		})
	}
	wg.Wait()
	return page, nil
}

func (self *CommentService) TotalComment(query *usercase.CommentQueryForm) (uint64, error) {
	total, err := self.repo.TotalComment(query)
	if err != nil {
		slog.Error("获取评论数量失败", "err", err.Error())
		return 0, tools.FiberServerError("查询评论数量失败")
	}
	return total, nil
}

func (self *CommentService) ManagePage(query *usercase.CommentQueryForm) (*usercase.PageData[usercase.CommentManageVo], error) {
	page, err := self.repo.ManagePage(query)
	if err != nil {
		slog.Error("管理端分页查询评论失败", "error", err)
		return nil, tools.FiberServerError("查询失败")
	}
	return page, nil
}

func (self *CommentService) UpdateSelectiveComment(form *usercase.CommentUpdateForm) error {
	if err := self.repo.UpdateSelective(form); err != nil {
		slog.Error("快捷更新评论失败", "error", err, "commentId", form.CommentId)
		return tools.FiberServerError("更新失败")
	}
	return nil
}

func (self *CommentService) Delete(commentId int64) error {
	if err := self.repo.DeleteById(commentId); err != nil {
		slog.Error("删除评论失败", "error", err)
		return tools.FiberServerError("删除失败")
	}
	return nil
}

func (self *CommentService) CommentVoteUp(commentId int64, userId uint64) error {
	if err := self.repo.CommentVoteUp(commentId, 1); err != nil {
		slog.Error("评论点赞失败", "error", err, "commentId", commentId)
		return tools.FiberServerError("点赞失败")
	}
	// 更新用户经验
	pool.Go(func() {
		_ = self.userService.UpdateUserExpertise(10, userId, 1)
	})
	return nil
}
