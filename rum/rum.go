package rum

import (
	"context"

	"github.com/gin-gonic/gin"
)

var _ IRumRouter = &Rum{}

type Rum struct {
	*gin.Engine
	rootGrp *RumGroup
}

// Default 创建一个默认的 Engine
func Default() *Rum {
	r := gin.Default()
	ctx := context.TODO()
	return NewWithEngine(ctx, r)
}

// NewWithEngine 使用自定义 gin engine 创建
func NewWithEngine(ctx context.Context, e *gin.Engine) *Rum {
	rum := &Rum{
		Engine: e,
	}

	if rum.rootGrp == nil {
		rum.rootGrp = baseRumGroup(ctx, rum, "/")
	}

	return rum
}

// Run 启动 gin-rum server。
func (rum *Rum) Run(addr ...string) error {
	return rum.Engine.Run(addr...)
}

// Group 扩展路由组， 可以顺带增加几个控制器
func (rum *Rum) Group(group string, classes ...ClassController) *RumGroup {
	// 04.1. 注册路由组
	return rum.rootGrp.Group(group, classes...)
}

// Use 使用中间件
func (rum *Rum) Use(fairs ...Fairing) IRumRoutes {
	rum.rootGrp.Use(fairs...)
	return rum
}

// Handle 绑定路由信息
func (rum *Rum) Handle(class ClassController) IRumRoutes {
	rum.rootGrp.Handle(class)
	return rum
}

func (rum *Rum) WithContext(fns ...ContextFunc) {
	rum.rootGrp.WithContext(fns...)
}
