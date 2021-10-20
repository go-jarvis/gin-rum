package main

import (
	"github.com/go-jarvis/gin-rum/cmd/demo/classes"
	"github.com/go-jarvis/gin-rum/rum"
)

func main() {

	// 1. 使用 rum 代替 gin
	g := rum.Default()
	// g.Attach(&middlewares.User{})

	// 2. 设置 base Path
	g.BasePath("/demo")

	// 2. 注册多个路由组
	g.Mount("/v1", classes.NewIndex())

	{
		v2Router := g.Mount("/v2")
		// 子路由注册中间件
		// v2Router.Attach(middlewares.NewUser())

		v2Router.Mount("/v3", classes.NewIndex())

	}

	// 3. 启动 rum server
	g.Launch()
}
