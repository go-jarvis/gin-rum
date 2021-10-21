package rum

import "context"

type ClassController interface {
	Method() string
	Path() string
	Handler(context.Context) (interface{}, error)
}

type HandlerFunc = func() (interface{}, error)
