package classes

import (
	"fmt"

	"github.com/go-jarvis/gin-rum/httpx"
)

// Index
type Index struct {
	httpx.MethodGet
	Name string `query:"name"`
}

func NewIndex() *Index {
	return &Index{}
}

func (index *Index) Path() string {
	return "/index"
}

// wanted
func (index *Index) Handler() (interface{}, error) {
	if index.Name != "wangwu" {
		return nil, fmt.Errorf("invalid user: %s", index.Name)
	}

	return "hello, gin-rum", nil
}
