package http

import "github.com/gin-gonic/gin"

type IOptionGrpc interface {
	Apply(*option)
}

type optionFunc func(*option)

func (f optionFunc) Apply(o *option) { f(o) }

func WithConnectorOption(port int, appName string) IOptionGrpc {
	return optionFunc(func(o *option) {
		o.port = port
		o.appName = appName
	})
}

type option struct {
	port         int
	appName      string
	interceptors []gin.HandlerFunc
}

func WithInterceptors(interceptors ...gin.HandlerFunc) IOptionGrpc {
	return optionFunc(func(o *option) {
		o.interceptors = interceptors
	})
}
