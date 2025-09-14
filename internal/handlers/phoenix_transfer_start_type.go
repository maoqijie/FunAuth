package handlers

import (
	"fmt"
	"net/http"
	"strings"

	auth "github.com/Yeah114/FunAuth/auth"
	"github.com/gin-gonic/gin"
)

func RegisterPhoenixTransferStartTypeRoute(api *gin.RouterGroup) {
	api.GET("/phoenix/transfer_start_type", func(c *gin.Context) {
		var q TransferStartTypeQuery
		if err := c.ShouldBindQuery(&q); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": fmt.Sprintf("RequestBindQuery: %v", err)})
			return
		}
		rawAuthorization := c.GetHeader("Authorization")
		token := strings.TrimPrefix(rawAuthorization, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Authorization: missing Bearer token"})
			return
		}
		var ok bool
		userID, ok := Authorizations[token]
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Authorization: invalid token"})
			return
		}
		enc, err := auth.TransferStartType(userID, q.Content)
		if err != nil {
			c.JSON(http.StatusOK, &TransferStartTypeResponse{Success: false, Message: fmt.Sprintf("TransferStartType: %v", err)})
			return
		}
		c.JSON(http.StatusOK, &TransferStartTypeResponse{Success: true, Data: enc})
		if !strings.HasPrefix(rawAuthorization, "cookie:") {
			delete(Authorizations, token)
		}
	})
}
