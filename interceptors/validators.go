package interceptors

import (
	"context"

	"google.golang.org/grpc"
)

type IValidator interface {
	Validate() error
}

func Validators() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if reqValidator, ok := req.(IValidator); ok {
			if err := reqValidator.Validate(); err != nil {
				return nil, err
			}
		}
		resp, err := handler(ctx, req)
		return resp, err
	}
}
