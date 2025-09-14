package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Yeah114/FunAuth/auth"
	"github.com/Yeah114/g79client"
	"github.com/gin-gonic/gin"
)

const fixedCookie = `{"sauth_json":"{\"gameid\":\"x19\",\"login_channel\":\"netease\",\"app_channel\":\"netease\",\"platform\":\"pc\",\"sdkuid\":\"aibgraaesciluppl\",\"sessionid\":\"1-eyJzIjogImRscmdoa2RnaTh1eXF6ZmcyZDdrM3UwbXduNWtzNTg0IiwgIm9kaSI6ICJhbWF3cmFhYWF3cjV0Mm9lLWQiLCAic2kiOiAiYTA0NzFiYTY4MjEzZmUyZGZlMDMwZWRmZmQ0NTQyNDljNGY1Mjk4NyIsICJ1IjogImFpYmdyYWFlc2NpbHVwcGwiLCAidCI6IDIsICJnX2kiOiAiYWVjZnJ4b2R5cWFhYWFqcCJ9\",\"sdk_version\":\"3.9.0\",\"udid\":\"sznjy5jkn80387y93rsc1wm3z23iws3q\",\"deviceid\":\"amawraaaawr5t2oe-d\",\"aim_info\":\"{\\\"aim\\\":\\\"127.0.0.1\\\",\\\"country\\\":\\\"CN\\\",\\\"tz\\\":\\\"+0800\\\",\\\"tzid\\\":\\\"\\\"}\",\"client_login_sn\":\"C14DB363E5934FE0F529E6642EBA4D0E\",\"gas_token\":\"\",\"source_platform\":\"pc\",\"ip\":\"127.0.0.1\"}"}`

func RegisterPhoenixLoginRoute(api *gin.RouterGroup) {
	api.POST("/phoenix/login", func(c *gin.Context) {
		rawAuthorization := c.GetHeader("Authorization")
		authorization := strings.TrimPrefix(rawAuthorization, "Bearer ")
		if authorization == "" {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "Authorization: missing Bearer token"})
			return
		}

		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": fmt.Sprintf("RequestBindJSON: %v", err)})
			return
		}
		cookieStr := req.FBToken

		cli, err := g79client.NewClient()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": fmt.Sprintf("NewClient: %v", err)})
			return
		}

		if cookieStr != "" {
			if err := cli.AuthenticateWithCookie(cookieStr); err != nil {
				c.JSON(http.StatusOK, gin.H{"success": false, "message": fmt.Sprintf("AuthenticateWithCookie: %v", err)})
				return
			}
		} else if err := cli.AuthenticateWithCookie(fixedCookie); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": fmt.Sprintf("AuthenticateWithCookie: %v", err)})
			return
		}

		loginRes, err := auth.Login(c.Request.Context(), cli, auth.LoginParams{
			ServerCode:      req.ServerCode,
			ServerPassword:  req.ServerPassword,
			ClientPublicKey: req.ClientPublicKey,
		})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": fmt.Sprintf("LoginFlow: %v", err)})
			return
		}

		enableSkin := true
		var skinInfo SkinInfo
		if enableSkin {
			authSkinInfo, err := auth.GetSkinInfo(cli)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"success": false, "message": fmt.Sprintf("GetSkinInfo: %v", err)})
				return
			}
			skinInfo = SkinInfo{
				ItemID:          authSkinInfo.ItemID,
				SkinDownloadURL: authSkinInfo.SkinDownloadURL,
				SkinIsSlim:      authSkinInfo.SkinIsSlim,
			}
		}

		resp := &LoginResponse{
			SuccessStates:  true,
			BotLevel:       loginRes.BotLevel,
			FBToken:        req.FBToken,
			RentalServerIP: loginRes.IP,
			ChainInfo:      loginRes.ChainInfo,
			BotSkin:        skinInfo,
		}
		c.JSON(http.StatusOK, resp)
		Authorizations[authorization] = cli.UserID
	})
}
