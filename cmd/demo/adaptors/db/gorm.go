package db

import (
	"context"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type contextGormDBType int

var contextGormDBKey contextGormDBType = 0

func WithGormDB(value interface{}) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, contextGormDBKey, value)
	}
}

func FromContextGormDB(ctx context.Context) *gorm.DB {
	return ctx.Value(contextGormDBKey).(*gorm.DB)
}

func NewGormDB() *gorm.DB {
	dsn := "root:Mysql12345@tcp(127.0.0.1:3306)/goftdemo?charset=utf8mb4&parseTime=True&loc=Local"
	gormdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return gormdb
}
