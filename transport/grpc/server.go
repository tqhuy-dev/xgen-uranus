package grpc

import (
	"fmt"
	"net"

	"github.com/tqhuy-dev/xgen-uranus/common"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	*grpc.Server
	option       option
	healthServer *health.Server
}

func NewServer(opts ...IOptionGrpc) *Server {
	opt := option{}
	for _, o := range opts {
		o.Apply(&opt)
	}
	if opt.zapLog == nil {
		opt.zapLog, _ = zap.NewProduction()
	}
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(opt.unaryInterceptors...))
	return &Server{Server: server, option: opt}
}

func (s *Server) Register(registerFunc func(server *Server)) *Server {
	registerFunc(s)
	return s
}

func (s *Server) ApplyHealth() *Server {
	s.healthServer = health.NewServer()
	return s
}

func (s *Server) StartGrpcServer() {
	if s.option.zapLog == nil {
		s.option.zapLog, _ = zap.NewProduction()
	}
	if s.healthServer != nil {
		registerHealth(s, s.healthServer)
		s.SwitchHealthStatusGrpc(healthpb.HealthCheckResponse_SERVING)
	}

	if s.option.useReflection {
		reflection.Register(s.Server)
	}
	s.run()
}

func (s *Server) run() {
	port := s.option.port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	s.option.zapLog.Info("grpc listen on port", zap.String(common.LogKeyAppName, s.option.appName),
		zap.Int("port", port))
	if err = s.Serve(lis); err != nil {
		panic(err)
	}
}

func registerHealth(server *Server, hs *health.Server) {
	healthpb.RegisterHealthServer(server, hs)
}

func (s *Server) SwitchHealthStatusGrpc(servingStatus healthpb.HealthCheckResponse_ServingStatus) {
	s.healthServer.SetServingStatus(s.option.appName, servingStatus)
}
