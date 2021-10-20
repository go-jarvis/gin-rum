# 初始化项目

```go
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, map[string]string{
			"hello": "gin-rum",
		})
	})

	if err := r.Run(":8089"); err != nil {
		panic(err)
	}
}
```