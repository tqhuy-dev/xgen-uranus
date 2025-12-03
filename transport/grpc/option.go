package grpc

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type IOptionGrpc interface {
	Apply(*option)
}

type optionFunc func(*option)

func (f optionFunc) Apply(o *option) { f(o) }

type ConnectorOption struct {
	Port    int
	AppName string
}

func WithConnectorOption(opt ConnectorOption) IOptionGrpc {
	return optionFunc(func(o *option) {
		o.port = opt.Port
		o.appName = opt.AppName
	})
}

func WithReflection() IOptionGrpc {
	return optionFunc(func(o *option) {
		o.useReflection = true
	})
}

func WithZapLog(zapLog *zap.Logger) IOptionGrpc {
	return optionFunc(func(o *option) {
		o.zapLog = zapLog
	})
}

func WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) IOptionGrpc {
	return optionFunc(func(o *option) {
		o.unaryInterceptors = interceptors
	})
}

type option struct {
	port              int
	useReflection     bool
	appName           string
	zapLog            *zap.Logger
	unaryInterceptors []grpc.UnaryServerInterceptor
}
