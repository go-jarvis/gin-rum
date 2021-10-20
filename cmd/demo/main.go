package main

import (
	"github.com/go-jarvis/gin-rum/cmd/demo/classes"
	"github.com/go-jarvis/gin-rum/rum"
)

func main() {

	// 1. 使用 rum 代替 gin
	g := rum.Default()
	// g.Attach(&middlewares.User{})

	g.Use(
	// &middlewares.CorsMid{},
	// middlewares.NewUserChecker(),
	)

	app := g.Group("demo")

	// 2. 注册多个路由组
	// g.AddGroup("/v1", classes.NewIndex())
	v1 := app.Group("/v1", classes.NewIndex())
	v1.Handle(&classes.User{})

	{
		v2Router := app.Group("/v2")
		// 子路由注册中间件
		// v2Router.Attach(middlewares.NewUser())

		v2Router.Group("/v3", classes.NewIndex())

	}

	// 3. 启动 rum server
	g.Run()
}
