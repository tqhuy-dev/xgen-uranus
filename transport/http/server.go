package http

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/tqhuy-dev/xgen-uranus/common"
	"go.uber.org/zap"
)

type Server struct {
	*http.Server
	option            option
	ginEngine         *gin.Engine
	toggleHealthCheck atomic.Bool
	zapLog            *zap.Logger
}

func NewServer(opts ...IOptionGrpc) *Server {
	opt := option{}
	for _, o := range opts {
		o.Apply(&opt)
	}

	r := gin.Default()
	r.Use(opt.interceptors...)
	return &Server{Server: &http.Server{
		Addr:    fmt.Sprintf(":%d", opt.port),
		Handler: r,
	}, option: opt, ginEngine: r}
}

func (s *Server) HealthCheck(path string) *Server {
	s.ginEngine.GET(path, func(c *gin.Context) {
		if s.toggleHealthCheck.Load() {
			c.JSON(200, gin.H{"status": "ok"})
		} else {
			c.JSON(503, gin.H{"status": "not ok"})
		}
	})
	return s
}

func (s *Server) RegisterRouter(registerFunc func(r *gin.Engine)) *Server {
	registerFunc(s.ginEngine)
	return s
}

func (s *Server) WithZapLog(zapLog *zap.Logger) *Server {
	s.zapLog = zapLog
	return s
}

func (s *Server) StartHttpServer() {
	if s.zapLog == nil {
		s.zapLog, _ = zap.NewProduction()
	}
	s.SwitchHealthCheck(true)
	s.zapLog.Info("http listen on port",
		zap.String(common.LogKeyAppName, s.option.appName),
		zap.Int("port", s.option.port))
	_ = s.ListenAndServe()
}

func (s *Server) SwitchHealthCheck(status bool) {
	s.toggleHealthCheck.Store(status)
}
