package classes

import "github.com/go-jarvis/gin-rum/httpx"

type User struct {
	httpx.MethodPost

	Name string `uri:"name"`
}

func (user *User) Path() string {
	return "/users/:name"
}

func (user *User) Handler() (interface{}, error) {
	return user.Name, nil
}
