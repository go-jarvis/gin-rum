# 控制器优化 - 更好用的控制器模式

之前在 [业务控制器模型](./02-class-controller.md) 中实现了一个简单的控制器模型， struct 对对象具有 `Build()` 方法， 就可以在 rum 中注册路由。

```go
// 第一版注册起方法
type Index struct {
	Name string `query:""` 
}

func (index *Index) Build(rum *rum.RumGroup) {
	rum.Handle("GET", "/index", handlerIndex)
}

func handlerIndex(c *gin.Context) {
	name:=c.Query("name")
	c.JSON(200, map[string]string{
		"hello": "gin-rum, " + name,
	})
}
```


但是这样的控制器模型实现后， 在使用上还是有很多的不方便。

1. **Struct 对象** 和将要绑定的 **http 方法**、 **路由地址** 完全割裂， 具体怎么定义需要在 `Build()` 方法中手工实现。
2. 要在业务逻辑中实现 handler `gin.HandleFunc` 方法，逻辑臃肿。
3. 变量绑定行为冗余， 每个 handler 都需要自己绑定变量并处理错误。



## 控制器优化

那要如何在 struct 定义的时候， 一站式完成所有需求约束呢？ 要这么样 rum 框架帮忙处理这些冗余操作呢？

分析 rum.Handle 的参数， 分别是

1. http 行为方法
2. 相对路由地址
3. 处理业务行为的 handler

```go
rum.Handle("GET", "/index", handlerIndex)

func (group *RouterGroup) Handle(httpMethod, relativePath string, handlers ...HandlerFunc)
```

### 重新定义控制器接口

因此可以将 `ClassController` 结构定义成

1. `Method() string` 方法: 返回 http 方法的值
2. `Path() string` 方法: 返回相对路由地址
3. `Handler() (interface{}, error)` 方法: 业务逻辑行为， 最终将被被 rum 封装成 handler 执行。

```go
type ClassController interface {
	Method() string
	Path() string
	Handler() (interface{}, error)
}
```

### 重载 Handle 行为

既然控制器的接口已经变了， 那么响应的行为也需要进行变更。

```go
// Mount 在 RumGroup 上绑定/注册 控制器
func (gg *RumGroup) Mount(group string, classes ...ClassController) *RumGroup {
	grp := newRumGroup(gg, group)

	// 老方法
	// for _, class := range claess {
	// 	class.Build(grp)
	// }

	// 新方法
	for _, class := range classes {
		grp.Handle(class)
	}

	return grp
}
```

之前老方法是将 RumGroup 作为结构体的 Build 方法的参数传递到内部进行处理。新方法中为了能 **简化请求参数的绑定和错误处理**， 因此不能再将 RumGroup 传递下去了。

重载之后的的 RumGroup.Handle 方法以 **结构体对象** 作为参数
1. 通过 `Method()` 和 `Path()` 方法获取必要参数
2. 将 `Handler()` 方法封装成为 `gin.HandlerFunc`
    + 使用 [ginbinder](https://github.com/tangx/ginbinder) 库读取 request 中的所有需要参数
    + 根据 `Handler()` 方法的结构，封装 gin 的返回结果。

```go
// Handle 重载 RumGroup 的 Handle 方法
func (gg *RumGroup) Handle(class ClassController) {

	m := class.Method()
	p := class.Path()
	handler := class.Handler

	// 将业务逻辑封装成为 gin.HandlerFunc
	handlerFunc := func(c *gin.Context) {
		// 绑定参数到对象中
		err := ginbinder.ShouldBindRequest(c, class)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		// 执行业务逻辑，获取返回值， 并封装返回结果
		v, err := handler()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		// 以 JSON 格式返回信息
		c.JSON(http.StatusOK, v)
	}

	// 调用 gin RouterGroup 的 Handle 方法注册路由
	gg.RouterGroup.Handle(m, p, handlerFunc)
}
```

### Handler 方法的传参数想法

最开的的想法， 是需要将 `*gin.Context` 传入到 Handler 方法中的， 这样使用者的行为能更广泛。 但后来思索再三后放弃了， Handler 方法更应该着力于 **业务逻辑** 而非管理或控制 gin 的行为。

```go
type ClassController interface {
	Method() string
	Path() string
	// Handler(*gin.Context) (interface{}, error)
	Handler() (interface{}, error)
}
```

## RumGroup.Handle 重载引起的 Fairing 控制器行为不兼容

之前在实现 Fairing 控制器的时候，为了方便 `Rum` 和 `RumGroup` 加载中间件的行为，使用了 `gin.IRoutes` 进行参数约束。

```go
func attachFairings(iroute gin.IRoutes, fairs ...Fairing) {
	for _, fair := range fairs {
		fair := fair
		// 创建一个临时中间件 handler
		handler := func(c *gin.Context) {
			_ = fair.OnRequest(c)
			c.Next()
		}
		// 使用 中间件
		iroute.Use(handler)
	}
}
```

但是在重载 `RumGroup` 的 Handle 方法之后签名发生了改变, 由 `Handle(string, string, ...HandlerFunc) IRoutes` 变成了 `Handle(class ClassController)`。 因此引发了不兼容的情况。

```go
type Rum struct {
	*gin.Engine
	rootGrp *RumGroup
}
```

因此最 Attach 也做了响应的改造。 将 Fairing 的接收者直接设置为了 RumGroup 。

```go
func (gg *RumGroup) attach(fairs ...Fairing) {
	for _, fair := range fairs {
		fair := fair

		// 创建一个临时中间件 handler
		handler := func(c *gin.Context) {
			_ = fair.OnRequest(c)
			c.Next()
		}

		// 使用 中间件
		gg.Use(handler)
	}
}
```
