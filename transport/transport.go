package transport

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/tqhuy-dev/xgen-uranus/common"
	"github.com/tqhuy-dev/xgen-uranus/transport/grpc"
	"github.com/tqhuy-dev/xgen-uranus/transport/http"
	"go.uber.org/zap"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type UranusServer struct {
	httpServer     *http.Server
	grpcServer     *grpc.Server
	shutdownOption common.GracefulShutdown
	zapLog         *zap.Logger
}

func NewUranusSever() *UranusServer {
	return &UranusServer{}
}

func (s *UranusServer) WithHttpServer(server *http.Server) *UranusServer {
	s.httpServer = server
	return s
}

func (s *UranusServer) WithGrpcServer(server *grpc.Server) *UranusServer {
	s.grpcServer = server
	return s
}

func (s *UranusServer) WithGracefulShutdown(shutdownOption common.GracefulShutdown) *UranusServer {
	s.shutdownOption = shutdownOption
	return s
}

func (s *UranusServer) WithZapLog(zapLog *zap.Logger) *UranusServer {
	s.zapLog = zapLog
	return s
}

func (s *UranusServer) Run() {
	if s.zapLog == nil {
		s.zapLog, _ = zap.NewProduction()
	}
	signalCtx, signalCtxStop := signal.NotifyContext(context.Background(), s.shutdownOption.Signal...)
	defer signalCtxStop()
	go func() {
		if s.grpcServer != nil {
			s.grpcServer.StartGrpcServer()
		}
	}()
	go func() {
		if s.httpServer != nil {
			s.httpServer.StartHttpServer()
		}
	}()
	<-signalCtx.Done()
	if s.httpServer != nil {
		s.httpServer.SwitchHealthCheck(false)
	}
	if s.grpcServer != nil {
		s.grpcServer.SwitchHealthStatusGrpc(healthpb.HealthCheckResponse_NOT_SERVING)
	}
	s.zapLog.Info(fmt.Sprintf("server stopped gracefully, waiting for shutdown for duration: %s", s.shutdownOption.Timeout.String()+""))
	time.Sleep(s.shutdownOption.Timeout)
	wg := &sync.WaitGroup{}

	wg.Go(func() {
		if s.grpcServer != nil {
			s.grpcServer.GracefulStop()
			s.zapLog.Info("grpc server stopped")
		}
	})

	wg.Go(func() {
		if s.httpServer != nil {
			err := s.httpServer.Shutdown(context.Background())
			if err != nil {
				s.zapLog.Error("http server shutdown error", zap.Error(err))
				time.Sleep(s.shutdownOption.HardStop)
			} else {
				s.zapLog.Info("http server stopped")
			}
		}
	})

	wg.Wait()
	s.zapLog.Info("shutdown complete")

	os.Exit(0)
}
