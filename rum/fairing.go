package rum

import (
	"github.com/gin-gonic/gin"
)

// fairing 中间件封装， 可以理解 fairing 是中间件的控制器
// fairing 只是提供中间件应该使用的业务逻辑， 但他并不是中间件的 handler
type Fairing interface {
	// 这里使用 *gin.Context 作为参数， 为了方便以后在中间件处理的时候获取请求体中的信息
	OnRequest(c *gin.Context) error
}
