# fairing 中间件控制器



## 中间件控制器

在 [fairing.go](/rum/fairing.go) 定义接口 Fairing

可以认为 fairing 是中间件的控制器， 只需要对象实现 Fairing 接口， 实现具体的的中间件处理逻辑。 而具体的 `gin.HandleFunc` 有 rum 进行生成和管理。

```go
type Fairing interface {
	// 这里使用 *gin.Context 作为参数， 为了方便以后在中间件处理的时候获取请求体中的信息
	OnRequest(c *gin.Context) error
}
```

既然 Fairing 需要对请求进行 **预处理** ， 那么接口方法就需要使用 `*gin.Context` 作为参数端的。


```go
func attachFairings(iroute gin.IRoutes, fairs ...Fairing) {
	for _, fair := range fairs {
		fair := fair

		// 创建一个临时中间件 handler
		handler := func(c *gin.Context) {
			// 不拦截错误
			_ = fair.OnRequest(c)
			c.Next()
		}

		// 使用 中间件
		iroute.Use(handler)
	}
}
```

1. 创建 handler 的时候， 传入的是 **原始 `*gin.Context`** 而非 `c.Copy()` 生成的副本对象。
    1. 这里不应该传入 cc 备份给 Middleware 处理。
    2. 某些中间件 **功能本身** 需要修改 gin.Context 中的一些内容。
    3. 如果要避免类似中间件读取 body， 而导致业务逻辑失效的话
2. 在 rum 作为框架， 在生成的 handler 中逻辑中 **不拦截错误**。 如果要对请求进行中断处理， 应该是中间件控制器本身的行为逻辑。 而中间件控制器本身有 `gin.Context` 作为参数， 本身也能完成终端。
3. 中间件是绑定到 **路由** 上面的， 在 gin 中， 实现路由的接口 `gin.IRoutes` ， 因此在 `attachFairing` 函数中作为了第一个参数传入。 这里可以是 `gin.Engine` 也可以是 `gin.RouterGroup`

```go
func (rum *Rum) Attach(fairs ...Fairing) {
	attachFairings(rum, fairs...)
}

// Attach 绑定/注册 中间件
func (gg *RumGroup) Attach(fairs ...Fairing) {
	attachFairings(gg, fairs...)
}
```



### 实现中间件控制器

在 [middlewares/user.go](/cmd/demo/middlewares/user.go) 中实现了一个检测是否为合法用户的中间件。 业务逻辑很简单。

如果 query 参数中获取的用户姓名不合法， 则中断请求。 

```go
type User struct {
	Name string `query:"name"`
}

func (user User) OnRequest(c *gin.Context) (err error) {

	user.Name = c.Query("name")
	if user.Name != "zhangsan" {
		err = fmt.Errorf("非法用户: %s", user.Name)

		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	return
}
```


## 使用中间件控制器

这里就很简单了， 只需要在想使用中间件的地方， `Attach` 对应的控制器就好了。

```go
package main

import (
	"github.com/tangx-labs/gin-rum/cmd/demo/classes"
	"github.com/tangx-labs/gin-rum/cmd/demo/middlewares"
	"github.com/tangx-labs/gin-rum/rum"
)

func main() {

	g := rum.Default()

	// 全局中间件
	g.Attach(&middlewares.User{})
	// 省略

	{
		v2Router := g.Mount("/v2")

		// 子路由注册中间件
		v2Router.Attach(middlewares.NewUser())
		v2Router.Mount("/v3", classes.NewIndex())
	}

	// ...省略
}
```

## 遗留问题

两种方法都能实现功能。

但是中间件控制器使用 **指针方法** 是否会在高并发的时候出现数据紊乱？

```go
// 值方法
func (user User) OnRequest(c *gin.Context) (err error) {
	user.Name = c.Query("name")
	// ...省略
	return
}

// 指针方法
func (user *User) OnRequest(c *gin.Context) (err error) {
	user.Name = c.Query("name")
	// ...省略
	return
}
```
