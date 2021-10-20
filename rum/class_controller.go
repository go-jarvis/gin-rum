package rum

type ClassController interface {
	Method() string
	Path() string
	Handler() (interface{}, error)
}

type HandlerFunc = func() (interface{}, error)
