package interceptors

import (
	"context"

	"github.com/tqhuy-dev/xgen-uranus/common"
	"github.com/tqhuy-dev/xgen/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func CorrelationTracing() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		correlationId := utilities.GenerateUUIDV7()

		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			if len(md[common.CorrelationIdKey]) > 0 && len(md[common.CorrelationIdKey][0]) > 0 {
				correlationId = md[common.CorrelationIdKey][0]
			}
		}
		ctx = context.WithValue(ctx, common.CorrelationIdKey, correlationId)
		resp, err := handler(ctx, req)
		return resp, err
	}
}
