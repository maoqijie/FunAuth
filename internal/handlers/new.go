package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var Authorizations = make(map[string]string)

func RegisterNewRoutes(rg *gin.RouterGroup) {
	rg.GET("/new", func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if strings.HasPrefix(authorization, "cookie:") {
			c.String(http.StatusOK, "ok")
			return
		}
		id := uuid.New()
		c.String(http.StatusOK, id.String())
	})
}
