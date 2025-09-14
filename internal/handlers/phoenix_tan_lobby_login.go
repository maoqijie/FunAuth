package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Yeah114/FunAuth/auth"
	"github.com/Yeah114/g79client"
	"github.com/gin-gonic/gin"
)

func RegisterPhoenixTanLobbyLoginRoute(api *gin.RouterGroup) {
	api.POST("/phoenix/tan_lobby_login", func(c *gin.Context) {
		rawAuthorization := c.GetHeader("Authorization")
		authorization := strings.TrimPrefix(rawAuthorization, "Bearer ")
		if authorization == "" {
			c.JSON(http.StatusOK, gin.H{"success": false, "error_info": "Authorization: missing Bearer token"})
			return
		}

		var req TanLobbyLoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "error_info": fmt.Sprintf("RequestBindJSON: %v", err)})
			return
		}
		cookieStr := req.FBToken

		cli, err := g79client.NewClient()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "error_info": fmt.Sprintf("NewClient: %v", err)})
			return
		}

		if cookieStr != "" {
			if err := cli.AuthenticateWithCookie(cookieStr); err != nil {
				c.JSON(http.StatusOK, gin.H{"success": false, "error_info": fmt.Sprintf("AuthenticateWithCookie: %v", err)})
				return
			}
		} else if err := cli.AuthenticateWithCookie(fixedCookie); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "error_info": fmt.Sprintf("AuthenticateWithCookie: %v", err)})
			return
		}

		loginRes, err := auth.TanLobbyLogin(c.Request.Context(), cli, auth.TanLobbyLoginParams{
			RoomID:       req.RoomID,
		})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "error_info": fmt.Sprintf("TanLobbyLogin: %v", err)})
			return
		}
		
		enableSkin := true
		var skinInfo SkinInfo
		if enableSkin {
			authSkinInfo, err := auth.GetSkinInfo(cli)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"success": false, "error_info": fmt.Sprintf("GetSkinInfo: %v", err)})
				return
			}
			skinInfo = SkinInfo{
				ItemID:          authSkinInfo.ItemID,
				SkinDownloadURL: authSkinInfo.SkinDownloadURL,
				SkinIsSlim:      authSkinInfo.SkinIsSlim,
			}
		}

		c.JSON(http.StatusOK, TanLobbyLoginResponse{
			Success:   true,
			ErrorInfo: "",
			RoomOwnerID: loginRes.RoomOwnerID,
			UserUniqueID: loginRes.UserUniqueID,
			UserPlayerName: loginRes.UserPlayerName,
			RaknetRand: loginRes.RaknetRand,
			RaknetAESRand: loginRes.RaknetAESRand,
			EncryptKeyBytes: loginRes.EncryptKeyBytes,
			DecryptKeyBytes: loginRes.DecryptKeyBytes,
			SignalingSeed: loginRes.SignalingSeed,
			SignalingTicket: loginRes.SignalingTicket,
			BotLevel: int(cli.UserDetail.Level.Int64()),
			BotSkin: skinInfo,
		})
	})
}
