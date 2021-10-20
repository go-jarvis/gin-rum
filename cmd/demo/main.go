package main

import (
	"github.com/go-jarvis/gin-rum/cmd/demo/classes"
	"github.com/go-jarvis/gin-rum/cmd/demo/middlewares"
	"github.com/go-jarvis/gin-rum/rum"
)

func main() {

	// 1. 使用 rum 代替 gin
	g := rum.Default()
	g.Use(&middlewares.CorsMid{})

	app := g.Group("demo")

	// 2. 注册多个路由组
	// g.AddGroup("/v1", classes.NewIndex())
	v1 := app.Group("/v1", classes.NewIndex())
	v1.Handle(&classes.User{})

	{
		v2 := app.Group("/v2")
		v2.Use(&middlewares.UserChecker{})
		// 子路由注册中间件
		v2.Group("/v3", classes.NewIndex())

	}

	// 3. 启动 rum server
	g.Run(":8089")
}
