package rum

import (
	"github.com/gin-gonic/gin"
)

type Rum struct {
	*gin.Engine
	rootGrp *RumGroup
}

// Default 创建一个默认的 Engine
func Default() *Rum {
	r := gin.Default()
	return NewWithEngine(r)
}

// NewWithEngine 使用自定义 gin engine 创建
func NewWithEngine(e *gin.Engine) *Rum {
	rum := &Rum{
		Engine: e,
	}

	rum.initial()

	return rum
}

// initial 初始化 Rum
func (rum *Rum) initial() {
	if rum.rootGrp == nil {
		rum.rootGrp = baseRumGroup(rum, "/")
	}
}

// Run 启动 gin-rum server。
func (rum *Rum) Run() error {
	return rum.Engine.Run(":8089")
}

// Group 扩展路由组， 可以顺带增加几个控制器
func (rum *Rum) Group(group string, classes ...ClassController) *RumGroup {
	// 04.1. 注册路由组
	return rum.rootGrp.Group(group, classes...)
}

// func (rum *Rum) Group() {}

// BasePath 设置 Rum 的根路由
func (rum *Rum) BasePath(group string) *Rum {
	rum.rootGrp = baseRumGroup(rum, group)

	return rum
}

// Use 使用中间件
func (rum *Rum) Use(fairs ...Fairing) *Rum {
	rum.rootGrp.Use(fairs...)
	return rum
}
