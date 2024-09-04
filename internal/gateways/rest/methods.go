package rest

import (
	"github.com/gin-gonic/gin"
	"wallet/internal/errs"
)

func (s *Server) pong() gin.HandlerFunc {
	return func(c *gin.Context) {
		Return(c, "pong", nil)
	}
}

func (s *Server) exists() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			Return(c, nil, errs.ErrInvalidParam)
			return
		}

		exists, err := s.walletExists.Execute(c.Request.Context(), id)
		if err != nil {
			Return(c, nil, err)
			return
		}

		Return(c, exists, nil)
	}
}
