package auth

import (
	"go-fiber-ent-web-layout/pkg/pool"
	"sync"
	"time"
)

// ManageLoginUser 管理端登录用户
type ManageLoginUser interface {
	LoginUser
	GetRoles() []string                  // 获取用户角色
	GetPermissions() []string            // 获取用户权限
	SetUsername(username string)         // 重设用户名称
	SetRoles(roles []string)             // 重设用户角色
	SetPermissions(permissions []string) // 重设用户权限
}

type ManageUserCache interface {
	// AddUser 添加管理端登录用户
	AddUser(token string, user ManageLoginUser, expire time.Duration)
	// ResetExpire 重置Token的过期时间
	ResetExpire(token string, expire time.Duration)
	// RemoveUser 删除管理端登录用户
	RemoveUser(token string)
	// RemoveUserById 通过用户id移除登录用户
	RemoveUserById(userId uint64)
	// GetUser 获取管理端登录用户
	GetUser(token string) ManageLoginUser
	// GetUserById 通过用户Id获取登录用户
	GetUserById(userId uint64) ManageLoginUser
}

func NewManageUserCache() ManageUserCache {
	return &inMemoryManageUserCache{
		userMap:  make(map[uint64]*cacheNode),
		tokenMap: make(map[string]uint64),
		nodePool: &sync.Pool{
			New: func() any {
				return new(cacheNode)
			},
		},
	}
}

type cacheNode struct {
	expireTime int64
	token      string
	value      ManageLoginUser
}

func (cn *cacheNode) Reset() {
	cn.expireTime = 0
	cn.value = nil
	cn.token = ""
}

type inMemoryManageUserCache struct {
	userMap  map[uint64]*cacheNode
	tokenMap map[string]uint64
	mutex    sync.RWMutex
	nodePool *sync.Pool
}

func (mc *inMemoryManageUserCache) AddUser(token string, user ManageLoginUser, expire time.Duration) {
	if user == nil {
		return
	}
	mc.mutex.Lock()
	node := mc.nodePool.Get().(*cacheNode)
	node.expireTime = time.Now().UnixMilli() + expire.Milliseconds()
	node.token = token
	node.value = user
	mc.userMap[user.GetUserId()] = node
	mc.tokenMap[token] = user.GetUserId()
	mc.mutex.Unlock()
}

func (mc *inMemoryManageUserCache) ResetExpire(token string, expire time.Duration) {
	mc.mutex.Lock()
	if userId, ok := mc.tokenMap[token]; ok {
		node := mc.userMap[userId]
		node.expireTime = time.Now().UnixMilli() + expire.Milliseconds()
	}
	mc.mutex.Unlock()
}

func (mc *inMemoryManageUserCache) RemoveUser(token string) {
	mc.mutex.Lock()
	if userId, ok := mc.tokenMap[token]; ok {
		node := mc.userMap[userId]
		delete(mc.userMap, userId)
		delete(mc.tokenMap, token)
		node.Reset()
		mc.nodePool.Put(node)
	}
	mc.mutex.Unlock()
}

// RemoveUserById 通过用户Id删除用户
// 由于map的key为token 所以只能遍历删除
func (mc *inMemoryManageUserCache) RemoveUserById(userId uint64) {
	mc.mutex.Lock()
	if node, ok := mc.userMap[userId]; ok {
		delete(mc.tokenMap, node.token)
		delete(mc.userMap, userId)
		node.Reset()
		mc.nodePool.Put(node)
	}
	mc.mutex.Unlock()
}

func (mc *inMemoryManageUserCache) GetUser(token string) ManageLoginUser {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	userId, ok := mc.tokenMap[token]
	if !ok {
		return nil
	}
	node := mc.userMap[userId]
	if node.expireTime <= time.Now().UnixMilli() {
		// 异步删除
		pool.Go(func() {
			mc.RemoveUser(token)
		})
		return nil
	}
	return node.value
}

func (mc *inMemoryManageUserCache) GetUserById(userId uint64) ManageLoginUser {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	node, ok := mc.userMap[userId]
	if !ok {
		return nil
	}
	if node.expireTime <= time.Now().UnixMilli() {
		// 异步删除
		pool.Go(func() {
			mc.RemoveUserById(userId)
		})
		return nil
	}
	return node.value
}

var defaultManagerUserCache ManageUserCache

// ManageUserCacheExpireTime 管理端登录用户的过期时间
const ManageUserCacheExpireTime = 30 * time.Minute

func init() {
	defaultManagerUserCache = NewManageUserCache()
}

// AddManageLoginUser 添加管理端登录用户
// token 请求中的token参数
// user 管理端登录用户
// expire 过期时间
func AddManageLoginUser(token string, user ManageLoginUser, expire time.Duration) {
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

// RemoveManageLoginUserById 通过用户Id删除管理端登录用户
func RemoveManageLoginUserById(userId uint64) {
	defaultManagerUserCache.RemoveUserById(userId)
}

// GetManageLoginUser 获取管理端登录用户
func GetManageLoginUser(token string) ManageLoginUser {
	return defaultManagerUserCache.GetUser(token)
}

// GetManageLoginUserById 通过用户Id获取管理端登录用户
func GetManageLoginUserById(userId uint64) ManageLoginUser {
	return defaultManagerUserCache.GetUserById(userId)
}
