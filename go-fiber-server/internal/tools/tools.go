package tools

import (
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/usercase"
)

func FiberRequestError(message string) *fiber.Error {
	return fiber.NewError(fiber.StatusBadRequest, message)
}

func FiberAuthError(message string) *fiber.Error {
	return fiber.NewError(fiber.StatusUnauthorized, message)
}

func FiberServerError(message string) *fiber.Error {
	return fiber.NewError(fiber.StatusInternalServerError, message)
}

// BuilderTree 将数据列表格式化为树形结构
// 使用泛型 待格式化的数据需要实现 Tree 接口
func BuilderTree[K any, T usercase.Tree[K]](list []T) []T {
	if len(list) <= 1 {
		return list
	}
	cacheMap := make(map[any]T)
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
