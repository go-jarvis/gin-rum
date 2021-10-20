package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserChecker struct {
	Name string `query:"name"`
}

func NewUserChecker() *UserChecker {
	return &UserChecker{}
}

// OnRequest 实现 Fairing 接口
// 这里是否应该使用 指针方法 呢？
//    即 `func (user User) OnRequest(c *gin.Context)`
func (user UserChecker) OnRequest(c *gin.Context) (err error) {

	validUser := "zhangsan"
	user.Name = c.Query("name")
	if user.Name != validUser {
		err = fmt.Errorf("中间件拦击， 非法用户: %s。 只允许: %s", user.Name, validUser)

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	return
}
