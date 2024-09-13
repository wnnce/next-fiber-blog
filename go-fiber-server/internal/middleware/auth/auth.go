package auth

import (
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"go-fiber-ent-web-layout/internal/tools"
	"time"
)

// LoginUser 登录用户接口
type LoginUser interface {
	GetUserId() uint64   // 获取用户Id
	GetUsername() string // 获取用户名
}

// CheckToken 检查是否存在Token且Token是否符合要求
func CheckToken(ctx fiber.Ctx) (*jwt.RegisteredClaims, error) {
	authorization := ctx.Get(fiber.HeaderAuthorization)
	if len(authorization) <= 7 {
		return nil, tools.FiberAuthError("The token does not exist")
	}
	claims, err := tools.VerifyToken(authorization[7:])
	now := time.Now()
	if err != nil || claims.NotBefore.After(now) {
		return nil, tools.FiberAuthError("Invalid token")
	}
	if claims.ExpiresAt.Before(now) {
		return nil, tools.FiberAuthError("The token has expired")
	}
	return claims, nil
}

// ManageAuth 管理端用户登录验证
// 如果Token验证成功那么将LoginUser存储到ctx.Locals中
func ManageAuth(ctx fiber.Ctx) error {
	claims, err := CheckToken(ctx)
	if err != nil {
		return err
	}
	token := claims.Subject
	loginUser := GetManageLoginUser(token)
	if loginUser == nil {
		return tools.FiberAuthError("Login user has expired")
	}
	fiber.Locals[ManageLoginUser](ctx, "loginUser", loginUser)
	// 请求处理完成后重设用户的过期时间
	defer ResetManageLoginUserExpire(token, ManageUserCacheExpireTime)
	return ctx.Next()
}

// ClassicAuth 博客用户登录验证
// 验证成功将用户登录信息保存到 Locals
func ClassicAuth(ctx fiber.Ctx) error {
	claims, err := CheckToken(ctx)
	if err != nil {
		return err
	}
	token := claims.Subject
	classicUser := GetClassicLoginUser(token)
	if classicUser == nil {
		return tools.FiberAuthError("Login user has expired")
	}
	fiber.Locals[ClassicLoginUser](ctx, "classicUser", classicUser)
	return ctx.Next()
}

// VerifyRoles 用户角色验证
func VerifyRoles(roles ...string) fiber.Handler {
	roleMap := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		roleMap[role] = struct{}{}
	}
	return func(ctx fiber.Ctx) error {
		if loginUser := fiber.Locals[ManageLoginUser](ctx, "loginUser"); loginUser != nil {
			for _, value := range loginUser.GetRoles() {
				if _, ok := roleMap[value]; ok || value == "admin" {
					return ctx.Next()
				}
			}
		}
		return fiber.NewError(fiber.StatusForbidden, "Current role has no permissions")
	}
}

// VerifyPermissions 用户权限验证
func VerifyPermissions(permissions ...string) fiber.Handler {
	permissionMap := make(map[string]struct{}, len(permissions))
	for _, permission := range permissions {
		permissionMap[permission] = struct{}{}
	}
	return func(ctx fiber.Ctx) error {
		if loginUser := fiber.Locals[ManageLoginUser](ctx, "loginUser"); loginUser != nil {
			for _, value := range loginUser.GetPermissions() {
				if _, ok := permissionMap[value]; ok || value == "all" {
					return ctx.Next()
				}
			}
		}
		return fiber.NewError(fiber.StatusForbidden, "No permission")
	}
}
