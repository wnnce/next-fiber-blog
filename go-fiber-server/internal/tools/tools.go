package tools

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"go-fiber-ent-web-layout/internal/usercase"
	"log/slog"
	"net/http"
)

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

// ComputeOffset 计算数据库查询分页偏移量
func ComputeOffset(total int64, page, size int, safe bool) int64 {
	if page == 0 {
		return 0
	}
	offset := int64((page - 1) * size)
	if !safe || offset < total {
		return offset
	} else {
		int64Size := int64(size)
		pages := total / int64Size
		if pages > 0 {
			return (pages - 1) * int64Size
		}
		return 0
	}

}
