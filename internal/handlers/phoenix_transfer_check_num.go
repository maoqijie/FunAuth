package handlers

import (
	"fmt"
	"net/http"

	auth "github.com/Yeah114/FunAuth/auth"
	"github.com/gin-gonic/gin"
)

func RegisterPhoenixTransferCheckNumRoute(api *gin.RouterGroup) {
	api.POST("/phoenix/transfer_check_num", func(c *gin.Context) {
		var req TransferCheckNumRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": fmt.Sprintf("RequestBindJSON: %v", err)})
			return
		}

		value, err := auth.TransferCheckNum(c.Request.Context(), req.Data)
		if err != nil {
			c.JSON(http.StatusOK, &TransferCheckNumResponse{Success: false, Message: fmt.Sprintf("TransferCheckNum: %v", err)})
			return
		}
		c.JSON(http.StatusOK, &TransferCheckNumResponse{Success: true, Value: value})
	})
}
