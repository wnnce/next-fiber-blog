package auth

import (
	"fmt"
	"go-fiber-ent-web-layout/internal/usercase"
	"testing"
)

func TestAddManageLoginUser(t *testing.T) {
	user := &usercase.SysUser{
		UserId:    1,
		Username:  "test",
		RoleNames: []string{"admin"},
	}
	AddManageLoginUser("1", user, ManageUserCacheExpireTime)
	loginUser := GetManageLoginUser("1")
	fmt.Printf("%v\n", loginUser)
	loginUser.SetUsername("demo")
	loginUser2 := GetManageLoginUser("1")
	fmt.Printf("%v\n", loginUser2)
}
