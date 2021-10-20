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
	// g.AddGroup("/v1", classes.NewIndex())
	g.AddGroup("/v1").Handle(classes.NewIndex())

	{
		v2Router := g.AddGroup("/v2")
		// 子路由注册中间件
		// v2Router.Attach(middlewares.NewUser())

		v2Router.AddGroup("/v3", classes.NewIndex())

	}

	// 3. 启动 rum server
	g.Run()
}
