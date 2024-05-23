package auth

import (
	"go-fiber-ent-web-layout/internal/tools/pool"
	"sync"
	"time"
)

// LoginUser 登录用户接口
type LoginUser interface {
	GetUserId() uint64
	GetUserName() string
	GetRoles() []string
	GetPermissions() []string
}

type ManagerUserCache interface {
	// AddUser 添加管理端登录用户
	AddUser(token string, user LoginUser, expire time.Duration)
	// ResetExpire 重置Token的过期时间
	ResetExpire(token string, expire time.Duration)
	// RemoveUser 删除管理端登录用户
	RemoveUser(token string)
	// GetUser 获取管理端登录用户
	GetUser(token string) LoginUser
}

func NewManagerUserCache() ManagerUserCache {
	return &inMemoryManagerUserCache{
		cache: make(map[string]*cacheNode),
		nodePool: &sync.Pool{
			New: func() any {
				return new(cacheNode)
			},
		},
	}
}

type cacheNode struct {
	expireTime int64
	value      LoginUser
}

func (cn *cacheNode) Reset() {
	cn.expireTime = 0
	cn.value = nil
}

type inMemoryManagerUserCache struct {
	cache    map[string]*cacheNode
	mutex    sync.RWMutex
	nodePool *sync.Pool
}

func (mc *inMemoryManagerUserCache) AddUser(token string, user LoginUser, expire time.Duration) {
	if user == nil {
		return
	}
	mc.mutex.Lock()
	node := mc.nodePool.Get().(*cacheNode)
	node.expireTime = time.Now().UnixMilli() + expire.Milliseconds()
	node.value = user
	mc.cache[token] = node
	mc.mutex.Unlock()
}

func (mc *inMemoryManagerUserCache) ResetExpire(token string, expire time.Duration) {
	mc.mutex.Lock()
	if node, ok := mc.cache[token]; ok {
		node.expireTime = time.Now().UnixMilli() + expire.Milliseconds()
	}
	mc.mutex.Unlock()
}

func (mc *inMemoryManagerUserCache) RemoveUser(token string) {
	mc.mutex.Lock()
	if node, ok := mc.cache[token]; ok {
		node.Reset()
		mc.nodePool.Put(node)
		delete(mc.cache, token)
	}
	mc.mutex.Unlock()
}

func (mc *inMemoryManagerUserCache) GetUser(token string) LoginUser {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	node, ok := mc.cache[token]
	if !ok {
		return nil
	}
	if node.expireTime <= time.Now().UnixMilli() {
		// 异步删除
		pool.Go(func() {
			mc.RemoveUser(token)
		})
		return nil
	}
	return node.value
}

var (
	defaultManagerUserCache ManagerUserCache
	// ManageUserCacheExpireTime 管理端登录用户的过期时间
	ManageUserCacheExpireTime = 30 * time.Minute
)

func init() {
	defaultManagerUserCache = NewManagerUserCache()
}

// AddManageLoginUser 添加管理端登录用户
// token 请求中的token参数
// user 管理端登录用户
// expire 过期时间
func AddManageLoginUser(token string, user LoginUser, expire time.Duration) {
	defaultManagerUserCache.AddUser(token, user, expire)
}

// ResetManageLoginUserExpire 重置管理端登录用户的过期时间
// token 请求携带的token
// expire 新的过期时间
func ResetManageLoginUserExpire(token string, expire time.Duration) {
	defaultManagerUserCache.ResetExpire(token, expire)
}

// RemoveManageLoginUser 删除管理端登录用户
func RemoveManageLoginUser(token string) {
	defaultManagerUserCache.RemoveUser(token)
}

// GetManageLoginUser 获取管理端登录用户
func GetManageLoginUser(token string) LoginUser {
	return defaultManagerUserCache.GetUser(token)
}
