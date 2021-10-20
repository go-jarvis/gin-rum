package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CorsMid struct {
}

func (cors *CorsMid) OnRequest(c *gin.Context) error {

	method := c.Request.Method
	origin := "*"
	if method != "" {
		c.Header("Access-Control-Allow-Origin", origin) // 可将将 * 替换为指定的域名
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization,X-Token")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
	}

	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}

	return nil
}
