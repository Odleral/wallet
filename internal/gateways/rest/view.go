package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"wallet/internal/errs"
)

const (
	success = "Success"
	failure = "Failure"
)

type R struct {
	Status    string      `json:"status"`
	ErrorNote string      `json:"error_note"`
	Data      interface{} `json:"data"`
}

func Return(c *gin.Context, data any, err error) {
	switch {
	case err == nil:
		c.JSON(http.StatusOK, R{
			Status:    success,
			ErrorNote: "",
			Data:      data,
		})
	case errors.Is(err, errs.ErrNotFound):
		c.JSON(http.StatusNotFound, R{
			Status:    failure,
			ErrorNote: err.Error(),
		})
	}
}
