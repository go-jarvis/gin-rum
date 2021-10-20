# 业务控制器模型

控制器模型可以让 1. 为了保障 main 的干净； 2. 让业务逻辑和路由信息结合的更紧密

1. 新建 [/classes/](/classes) 目录存在所有控制器
2. 在目录中创建 [index.go](/classes/index.go) 作为首页控制器

其中， Index 需要内嵌 `*gin.Engine` 以方便控制路由信息。

```go
type Index struct {
	e *gin.Engine
}
```

其二， Index 需要具有 `Build()` 方法， 这里就是执行 **注册路由信息** 的地方。

```go
func (index *Index) Build() {
	index.e.Handle("GET", "/", handlerIndex)
}

func handlerIndex(c *gin.Context) {
	c.JSON(200, map[string]string{
		"hello": "gin-rum",
	})
}
```

最后， 在 [main.go](/cmd/rum/main.go) 中创建 Index 副本并注册路由

```go
func main() {
	r := gin.Default()

	// 3. 向 gin engine 注册路由信息
	classes.NewIndex(r).Build()

	if err := r.Run(":8089"); err != nil {
		panic(err)
	}
}
```

注意这里一定要调用 `Build()` 方法， 只有调用了 build 方法， 路由信息才会被注册到 gin 中。

