package main

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/tqhuy-dev/xgen-uranus/common"
	"github.com/tqhuy-dev/xgen-uranus/examples/pb/greet"
	"github.com/tqhuy-dev/xgen-uranus/interceptors"
	"github.com/tqhuy-dev/xgen-uranus/transport"
	"github.com/tqhuy-dev/xgen-uranus/transport/grpc"
	"github.com/tqhuy-dev/xgen-uranus/transport/http"
	"go.uber.org/zap"
)

func main() {
	zapLog, _ := zap.NewProduction()
	serverGrpc := grpc.NewServer(
		grpc.WithConnectorOption(grpc.ConnectorOption{Port: 10000, AppName: "simple-grpc-app"}),
		grpc.WithReflection(),
		grpc.WithZapLog(zapLog),
		grpc.WithUnaryInterceptors(
			interceptors.CorrelationTracing(),
			interceptors.ZapLogInterceptor(zapLog, interceptors.WithAppName("simple-grpc-app")),
			grpc_recovery.UnaryServerInterceptor(),
			interceptors.Validators(),
		),
	).Register(func(server *grpc.Server) {
		Register(server)
	}).ApplyHealth()

	serverHttp := http.NewServer(
		http.WithConnectorOption(1323, "simple-http-app"),
		http.WithInterceptors(
			interceptors.HttpCorrelationTracing(),
			interceptors.HttpZapLogMiddleware(zapLog, interceptors.WithHttpAppName("simple-http-app")),
			gin.Recovery(),
		),
	).
		HealthCheck("health").RegisterRouter(func(r *gin.Engine) {
		RegisterRouter(r)
	}).WithZapLog(zapLog)

	uranusApp := transport.NewUranusSever().WithGracefulShutdown(common.GracefulShutdown{
		Timeout:  5 * time.Second,
		HardStop: 10 * time.Second,
		Delay:    10 * time.Second,
		Signal:   common.SignalStopDefault,
	}).WithHttpServer(serverHttp).WithGrpcServer(serverGrpc).WithZapLog(zapLog)
	uranusApp.Run()
}

type GreetServer struct {
	greet.UnimplementedGreetServer
}

func (*GreetServer) Greet(ctx context.Context, req *greet.HelloRequest) (*greet.HelloReturn, error) {
	return &greet.HelloReturn{Reply: "Hello " + req.Name}, nil
}
