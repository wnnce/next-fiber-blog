package tools

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	res "go-fiber-ent-web-layout/internal/tools/res"
	"go-fiber-ent-web-layout/internal/usercase"
	"log/slog"
	"net/http"
)

func CustomStackTraceHandler(ctx fiber.Ctx, e interface{}) {
	trace := fmt.Sprintf("fiber application panic, StackTrace:%v, uri:%s, method:%s", e, ctx.OriginalURL(), ctx.Method())
	slog.Error(trace)
}

func CustomErrorHandler(ctx fiber.Ctx, err error) error {
	code, message := http.StatusInternalServerError, "server error"
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}
	result := res.Fail(code, message)
	return ctx.Status(code).JSON(result)
}

func FiberRequestError(message string) *fiber.Error {
	return fiber.NewError(http.StatusBadRequest, message)
}

func FiberAuthError(message string) *fiber.Error {
	return fiber.NewError(http.StatusUnauthorized, message)
}

func FiberServerError(message string) *fiber.Error {
	return fiber.NewError(http.StatusInternalServerError, message)
}

// SqlxRowsScan SqlxRows封装辅助函数
func SqlxRowsScan[T any](rows *sqlx.Rows, list []*T) []*T {
	for rows.Next() {
		var entity T
		if err := rows.StructScan(&entity); err != nil {
			slog.Error(fmt.Sprintf("sqlx rows scan error, message:%s", err))
		}
		list = append(list, &entity)
	}
	return list
}

// BuilderTree 将数据列表格式化为树形结构
// 使用泛型 待格式化的数据需要实现 Tree 接口
func BuilderTree[T usercase.Tree](list []T) []T {
	cacheMap := make(map[int64]T)
	roots := make([]T, 0)
	for _, v := range list {
		cacheMap[v.GetId()] = v
	}
	for _, v := range list {
		parent, ok := cacheMap[v.GetParentId()]
		if ok {
			parent.AppendChild(v)
		} else {
			roots = append(roots, v)
		}
	}
	return roots
}
