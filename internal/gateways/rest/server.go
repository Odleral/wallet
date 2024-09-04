package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
	"wallet/internal/config"
	"wallet/internal/usecases/exists"
)

type Server struct {
	l          *zap.Logger
	router     *gin.Engine
	httpServer *http.Server

	walletExists *exists.UseCase
}

func New(cfg config.Config, we *exists.UseCase) *Server {
	r := gin.New()
	r.Use(gin.Recovery())

	s := Server{
		l:      zap.L(),
		router: r,
	}

	httpServer := &http.Server{
		Addr:              cfg.Port,
		Handler:           &s,
		ReadHeaderTimeout: time.Minute,
	}

	s.httpServer = httpServer

	s.endpoints()

	return &s
}

func (s *Server) Run() {
	s.l.Info("starting server", zap.String("port", s.httpServer.Addr))
	if err := s.httpServer.ListenAndServe(); err != nil {
		s.l.Error("server error", zap.Error(err))
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.l.Error("shutdown error", zap.Error(err))
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
