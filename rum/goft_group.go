package rum

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tangx/ginbinder"
)

type RumGroup struct {
	*gin.RouterGroup
}

// baseRumGroup 通过 Rum 返回一个根 RumGroup
func baseRumGroup(r *Rum, group string) *RumGroup {
	return &RumGroup{
		RouterGroup: r.RouterGroup.Group(group),
	}
}

// newRumGroup 通过 RumGroup 扩展新的 RumGroup
func newRumGroup(base *RumGroup, group string) *RumGroup {
	return &RumGroup{
		RouterGroup: base.RouterGroup.Group(group),
	}
}

// Group 在 RumGroup 上绑定/注册 控制器
func (grp *RumGroup) Group(group string, classes ...ClassController) *RumGroup {
	new := newRumGroup(grp, group)
	for _, class := range classes {
		new.Handle(class)
	}
	return new
}

// Use 绑定/注册 中间件
func (grp *RumGroup) Use(fairs ...Fairing) *RumGroup {
	return grp.use(fairs...)
}

func (grp *RumGroup) use(fairs ...Fairing) *RumGroup {
	for _, fair := range fairs {
		fair := fair

		// 创建一个临时中间件 handler
		handler := func(c *gin.Context) {

			// cc := c.Copy()
			// 这里不应该传入 cc 备份给 Middleware 处理。
			// 某些中间件可能就是需要修改 gin.Context 中的一些内容。
			// 如果要避免类似中间件读取 body， 而导致业务逻辑失效的话
			//    可以在 OnRequest 中自行使用 cc 副本
			// if err := fair.OnRequest(cc); err != nil {
			// 	// c.Abort()
			// 	c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{
			// 		"err": err.Error(),
			// 	})
			// 	return
			// }

			// 由于 rum 是一个框架， 不应该对任何已经放行的中间件做任何阻拦
			// 如果需要中断， 可以在业务实现的中间件本身中进行阻拦。
			_ = fair.OnRequest(c)
			c.Next()
		}

		// 使用 中间件
		grp.RouterGroup.Use(handler)
	}

	return grp
}

// Handle 重载 RumGroup 的 Handle 方法
func (grp *RumGroup) Handle(class ClassController) {

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

		// 执行业务逻辑，获取返回值
		v, err := handler()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		// 以 JSON 格式返回信息
		c.JSON(http.StatusOK, v)
	}

	// 调用 gin RouterGroup 的 Handle 方法注册路由
	grp.RouterGroup.Handle(m, p, handlerFunc)
}
