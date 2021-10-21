package classes

import (
	"context"

	"github.com/go-jarvis/gin-rum/cmd/demo/adaptors/db"
	"github.com/go-jarvis/gin-rum/httpx"
	"gorm.io/gorm"
)

// User 定义 gorm 的 User 模型
type User struct {
	gorm.Model
	UserId   int    `gorm:"index"`
	UserName string `gorm:"index"`
}

// GetUserByID class 控制器
type GetUserByID struct {
	httpx.MethodPost

	UserID int `uri:"id"`
}

func (user *GetUserByID) Path() string {
	return "/users/:id"
}

func (user *GetUserByID) Handler(ctx context.Context) (interface{}, error) {
	// 获取 ctx 中注入的 *gorm.DB 对象
	gorm := db.FromContextGormDB(ctx)

	userModel := &User{}
	tx := gorm.Where("user_id=?", user.UserID).First(&userModel)

	return userModel, tx.Error
}
