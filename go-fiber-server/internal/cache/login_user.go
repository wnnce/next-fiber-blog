package cache

import (
	"sync"
)

// LoginUser 登录用户接口
type LoginUser interface {
	GetUserId() int64
	GetUserName() string
	GetRoles() []string
	GetPermissions() []string
}

type LoginUserCache interface {
	AddLoginUser(user LoginUser)
	RemoveLoginUser(userId int64)
	GetLoginUser(userId int64) LoginUser
}

var (
	requestUserCache = make(map[int64]LoginUser)
	requestUserLock  sync.RWMutex
)

type InMemoryLoginUserCache struct {
	cache map[int64]LoginUser
	mutex sync.RWMutex
}

func NewLoginUserCache() LoginUserCache {
	return &InMemoryLoginUserCache{
		cache: make(map[int64]LoginUser),
	}
}

// AddLoginUser 添加登录用户
func (i *InMemoryLoginUserCache) AddLoginUser(user LoginUser) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.cache[user.GetUserId()] = user
}

// RemoveLoginUser 删除登录用户
func (i *InMemoryLoginUserCache) RemoveLoginUser(userId int64) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	delete(i.cache, userId)
}

// GetLoginUser 查询登录用户
func (i *InMemoryLoginUserCache) GetLoginUser(userId int64) LoginUser {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	return i.cache[userId]
}

// SetRequestUser 设置当前请求的用户
func SetRequestUser(requestId int64, user LoginUser) {
	requestUserLock.Lock()
	defer requestUserLock.Unlock()
	requestUserCache[requestId] = user
}

// GetRequestUser 获取当前请求的用户
func GetRequestUser(requestId int64) LoginUser {
	requestUserLock.RLock()
	defer requestUserLock.RUnlock()
	return requestUserCache[requestId]
}

// ClearRequestUser 请求处理完成后，清除当前请求的登录用户
func ClearRequestUser(requestId int64) {
	requestUserLock.Lock()
	defer requestUserLock.Unlock()
	delete(requestUserCache, requestId)
}
