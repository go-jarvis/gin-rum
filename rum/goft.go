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

// Launch 启动 gin-rum server。
// 这里由于重载问题， 不能将启动方法命名为 Run
func (rum *Rum) Launch() error {
	return rum.Run(":8089")
}

// Mount 挂载控制器
// 03.1. 关联控制器与 rum
// 03.2. 返回 *RumGroup 是为了方便链式调用
func (rum *Rum) Mount(group string, classes ...ClassController) *RumGroup {
	// 04.1. 注册路由组
	return rum.rootGrp.Mount(group, classes...)
}

// BasePath 设置 Rum 的根路由
func (rum *Rum) BasePath(group string) *Rum {
	rum.rootGrp = baseRumGroup(rum, group)

	return rum
}

func (rum *Rum) Attach(fairs ...Fairing) {
	rum.rootGrp.Attach(fairs...)
}
