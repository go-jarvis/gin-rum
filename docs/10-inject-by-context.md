# 使用 context 上下文完成依赖注入

在开发过程中， 不可避免的会用到诸如 **数据库、redis** 等其他组件。 使用 **依赖注入** 的方式可以很好的对程序进行解耦。

## 选择 context 作为容器


之所以选择 Context 作为容器， 

**其一** ， context 具有很强的裂变性，不同 RumGroup 的 context 可以添加属于自己的内容； 

```go
// github.com/go-jarvis/gin-rum/rum/rumgroup.go

// baseRumGroup 通过 Rum 返回一个根 RumGroup
func baseRumGroup(ctx context.Context, r *Rum, group string) *RumGroup {
	return &RumGroup{
		RouterGroup: r.RouterGroup.Group(group),
		ctx:         ctx,
	}
}

// newRumGroup 通过 RumGroup 扩展新的 RumGroup
func newRumGroup(base *RumGroup, group string) *RumGroup {
	return &RumGroup{
		RouterGroup: base.RouterGroup.Group(group),
		ctx:         base.ctx,
	}
}

// WithContext 向 RumGoft 中注入任何内容
// 以向 class 控制器传递
func (grp *RumGroup) WithContext(fns ...ContextFunc) {
	for _, fn := range fns {
		grp.ctx = fn(grp.ctx)
	}
}
```

**其二** ， context 的包容性， **key ， value** 都是 `interface{}` 可以包容一切。

```go
// context/context.go
func WithValue(parent Context, key, val interface{}) Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	if key == nil {
		panic("nil key")
	}
	if !reflectlite.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}
	return &valueCtx{parent, key, val}
}
```

### rum 增加上下文支持

随后， 在 **RumGroup** 中增加 `context.Context` 字段，并添加相关方法。

```go
type RumGroup struct {
	*gin.RouterGroup
	ctx context.Context
}

type ContextFunc = func(context.Context) context.Context

// WithContext 向 RumGoft 中注入任何内容
// 以向 class 控制器传递
func (grp *RumGroup) WithContext(fns ...ContextFunc) {
	for _, fn := range fns {
		grp.ctx = fn(grp.ctx)
	}
}
```

`WithContext` 方法要求传入的是一个 **操作 Context** 的函数， 这个函数由用户自己实现。WithContext 方法对于 Context 的修改仅限于 RumGroup 自身与及其子 RumGroup。

另外， 在 ClassController 中， 也需要做响应的变更， 需要 `Handler(ctx context.Context)` 方法支持 context 作为参数传递。

```go
type ClassController interface {
	Method() string
	Path() string

	// Handler() (interface{}, error)  // 老方法
	Handler(context.Context) (interface{}, error)
}
```


## 为什么不用 `gin.Context`

虽然 **gin.Context** 也实现了 `context.Context` 的接口， 在在和我们常用的 **标准库** 还是还是有很多差别

**首先**， gin.Context 在 gin 初始化的时候会生成一个 **公共的祖先 gin.Context**， 随着程序的启动， 用户每次请求都将创建一个独立的 **gin.Context** 副本。 

由于 gin 并没有提供一个可以用户初始化的 **gin.Context** 的 API。

```go
func New() *Engine {
	debugPrintWARNINGNew()
	engine := &Engine{
		// ...省略
	}
	engine.RouterGroup.engine = engine
	// pool 是 engine 的私有字段
	engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}
	return engine
}

// allocateContext 也是 engine 的私有方法
func (engine *Engine) allocateContext() *Context {
	v := make(Params, 0, engine.maxParams)
	skippedNodes := make([]skippedNode, 0, engine.maxSections)
	return &Context{engine: engine, params: &v, skippedNodes: &skippedNodes}
}
```

因此用户只能写入到 **每次请求** 的 gin.Context 中。
而类似 **数据库连接池** 这样的句柄， 在程序启动的时候就初始化了，不在会改变。 如果写入到 gin.Context 中就造成了计算资源的浪费。


**其次**， 在 `gin.Context.Value()` 方法首先与标准库的实现不同， 有一个 **比较的致命问题**

```go
// github.com/gin-gonic/gin@v1.7.4/context.go
func (c *Context) Value(key interface{}) interface{} {
	if key == 0 {
		return c.Request
	}
	if keyAsString, ok := key.(string); ok {
		val, _ := c.Get(keyAsString)
		return val
	}
	return nil
}
```

可以看到， gin 中的 `Value()` 方法将 **key** 转为了字符串。 因为此失去了 **数据类型** 的支持， key 的唯一性概率就大大降低了， 很容易发生覆盖冲突。

这一点在标准库中的 `valueCtx` 就不存在这种情况， 因为 **key 不会被断言， 类型也是一部分**。

```go
func (c *valueCtx) Value(key interface{}) interface{} {
	if c.key == key {
		return c.val
	}
	return c.Context.Value(key)
}
```

测试代码参考 [context-Context-and-gin-Context](https://github.com/tangx-labs/context-Context-and-gin-Context)


### 遗留问题

由于使用了自建的 GroupGroup Context， 并且 gin.Context 没有交叉点。 因此 gin.Context 中与 **Cancel** 相关的方法也无法传递到 ClassController 中的 Handler 方法中。


## demo: 使用 context 传递

![gorm 数据库句柄](/cmd/demo/adaptors/db/gorm.go) 通过 Context 注入到 Rum 中。


### 封装 db adaptor

首先， 创建新的数据类型， 并用该类型创建 **唯一 key**。 

```go
type contextGormDBType int

var contextGormDBKey contextGormDBType = 0

```

在实践中， 不同的适配器可以创建自己的数据类型。 如此一来， 即使 **字面值** 相同也不会冲突、覆盖。

其次， 实现 **注入** 与 **提取** 函数。

```go
// WithGormDB 注入 *gorm.DB 到 context 中
func WithGormDB(vaule *gorm.DB) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, contextGormDBKey, value)
	}
}

// FromContextGormDB 从 context 中提取 *gorm.DB
func FromContextGormDB(ctx context.Context) *gorm.DB {
	return ctx.Value(contextGormDBKey).(*gorm.DB)
}
```

### 在 rum 中注入 adaptor

```go
func main() {

	// 1. 使用 rum 代替 gin
	g := rum.Default()
	g.WithContext(
		db.WithGormDB(db.NewGormDB()),
	)

// 省略
```

### 在 ClassController 中使用 adaptor

最后， 在 ClassController 实例中直接使用。 

> 注意， 在控制器定义的时候依旧保持 **干净**， 无任何依赖适配器的字段。

```go

// GetUserByID class 控制器
type GetUserByID struct {
	httpx.MethodPost

	UserID int `uri:"id"`
}

func (user *GetUserByID) Handler(ctx context.Context) (interface{}, error) {
	// 获取 ctx 中注入的 *gorm.DB 对象
	gorm := db.FromContextGormDB(ctx)

	userModel := &User{}
	tx := gorm.Where("user_id=?", user.UserID).First(&userModel)

	return userModel, tx.Error
}
```
