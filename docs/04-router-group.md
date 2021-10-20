# 注册路由组

在 gin 中有路由组的概念， 可以理解为路由的 prefix。

## rum 增加路由组

1. 在 Rum 中增加路由组 `rg *gin.RouterGroup`

```go
type Rum struct {
	*gin.Engine
	rg *gin.RouterGroup
}
```

2. 有了路由组字段之后， 就需要使用起来。 在 `Mount()` 方法中， 增加 `group name` 传参数

```go
// Mount 参数中增加了 group 的传参
func (rum *Rum) Mount(group string, classes ...ClassController) *Rum {

	// 04.1. 注册路由组
	rum.rg = rum.Group(group)

	for _, class := range classes {
		// 03.3. 将 rum 传入到控制器中
		class.Build(rum)
	}

	return rum
}
```

有了 group name 之后， 肯定是要将其注册到 rum engine 中。  

```go
rum.rg = rum.Group(group)
```

3. 为了能在不改变控制器的情况下使用 **路由组** 路径， 需要 **重载** rum 的 `Handle` 方法。

```go
// Handle 重载 gin.Engine 的 Handle 方法。
// 04.2. 这样子路由注册的时候， 就直接挂载到了 RouterGroup 上， 有了层级关系
func (rum *Rum) Handle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) {
	rum.rg.Handle(httpMethod, relativePath, handlers...)

    return rum
}
```

重载 Handle 方法之后， 控制器的子路由就被路由组分组了。


## 挂载路由组

在 [main.go](/cmd/rum/main.go) 中， 为 Mount 方法增加路由组 `v1`， 并添加了一个新的路由组 `v2`

```go
	// 2. 注册路由
	g.Mount("/v1",
		classes.NewIndex(),
	)
	// 04.2. 注册多个路由组。
	g.Mount("/v2",
		classes.NewIndex(),
	)
```

启动服务后，可以看到两组路由， v1 和 v2

```bash
# cd cmd/rum/ && go run .
[GIN-debug] GET    /v1/                      --> github.com/tangx-labs/gin-rum/classes.handlerIndex (3 handlers)
[GIN-debug] GET    /v2/                      --> github.com/tangx-labs/gin-rum/classes.handlerIndex (3 handlers)
[GIN-debug] Listening and serving HTTP on :8089
```

## 遗留问题

在 gin 中， RouterGroup 是可以一级一级往下扩展的。 但是在当前 rum 中所有的路由组都是挂载到 `gin.Engine` 上的， 所以就丢失了这个功能。
