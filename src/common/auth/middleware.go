package auth

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sportivaid/go-template/util/httputil"
)

const (
	AccessToken = "AccessToken"
)

type Middleware struct {
}

func (m *Middleware) AuthUserToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		authHeader := c.GetHeader("Authorization")

		if authHeader != AccessToken {
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, errors.New("Invalid Auth Token"))
			return
		}

		c.Next()
	}
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}
