package rest

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
	"wallet/internal/config"
)

type Server struct {
	l          *zap.Logger
	router     *gin.Engine
	httpServer *http.Server
}

func New(cfg config.Config) *Server {
	r := gin.New()
	r.Use(gin.Recovery())

	return &Server{
		l:      zap.L(),
		router: r,
		httpServer: &http.Server{
			Addr:        cfg.Port,
			Handler:     r,
			ReadTimeout: time.Minute,
		},
	}
}

func (s *Server) Start() {

}
