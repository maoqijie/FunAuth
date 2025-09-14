package handlers

import (
	"net/http"
	"strconv"
	"strings"

	g79 "github.com/Yeah114/g79client"
	"github.com/gin-gonic/gin"
)

// extractCookieFromAuth 从 Authorization 头中提取 cookie:<cookie>
func extractCookieFromAuth(c *gin.Context) (string, bool) {
	raw := c.GetHeader("Authorization")
	if strings.HasPrefix(raw, "cookie:") {
		return strings.TrimPrefix(raw, "cookie:"), true
	}
	return "", false
}

// RegisterOpenAPIRoutes registers cookie-required helper endpoints under /api/open
func RegisterOpenAPIRoutes(api *gin.RouterGroup) {
	open := api.Group("/open")

	// 获取用户详情（需要 cookie）
	open.GET("/g79/user_detail", func(c *gin.Context) {
		cookie, ok := extractCookieFromAuth(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "authorization is required (use Authorization: cookie:<cookie>)"})
			return
		}
		cli, err := g79.NewClient()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "message": err.Error()})
			return
		}
		if err := cli.AuthenticateWithCookie(cookie); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
			return
		}
		if cli.UserDetail == nil {
			if detail, err := cli.GetUserDetail(); err == nil && detail != nil {
				cli.UserDetail = &detail.Entity
			}
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "user": cli.UserDetail})
	})

	// 租赁服按名称搜索（需要 cookie）
	open.GET("/g79/rental_search", func(c *gin.Context) {
		cookie, ok := extractCookieFromAuth(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "authorization is required (use Authorization: cookie:<cookie>)"})
			return
		}
		name := c.Query("name")
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "name is required"})
			return
		}
		cli, err := g79.NewClient()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "message": err.Error()})
			return
		}
		if err := cli.AuthenticateWithCookie(cookie); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
			return
		}
		resp, err := cli.SearchRentalServerByName(name)
		if err != nil || resp == nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "search failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "result": resp})
	})

	// 批量获取可用租赁服（需要 cookie）
	open.GET("/g79/rental_available", func(c *gin.Context) {
		cookie, ok := extractCookieFromAuth(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "authorization is required (use Authorization: cookie:<cookie>)"})
			return
		}
		sortType, _ := strconv.Atoi(c.DefaultQuery("sort_type", "0"))
		orderType, _ := strconv.Atoi(c.DefaultQuery("order_type", "0"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		cli, err := g79.NewClient()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "message": err.Error()})
			return
		}
		if err := cli.AuthenticateWithCookie(cookie); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
			return
		}
		list, err := cli.GetAvailableRentalServers(sortType, orderType, offset)
		if err != nil || list == nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "query failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "result": list})
	})

	// 租赁服详情（需要 cookie）
	open.GET("/g79/rental_details", func(c *gin.Context) {
		cookie, ok := extractCookieFromAuth(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "authorization is required (use Authorization: cookie:<cookie>)"})
			return
		}
		id := c.Query("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id is required"})
			return
		}
		cli, err := g79.NewClient()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "message": err.Error()})
			return
		}
		if err := cli.AuthenticateWithCookie(cookie); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
			return
		}
		details, err := cli.GetRentalServerDetails(id)
		if err != nil || details == nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "query failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "result": details})
	})

	// 用户设置（需要 cookie）
	open.GET("/g79/user_settings", func(c *gin.Context) {
		cookie, ok := extractCookieFromAuth(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "authorization is required (use Authorization: cookie:<cookie>)"})
			return
		}
		cli, err := g79.NewClient()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "message": err.Error()})
			return
		}
		if err := cli.AuthenticateWithCookie(cookie); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
			return
		}
		settings, err := cli.GetUserSettingList()
		if err != nil || settings == nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "query failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "result": settings})
	})

	// 用户搜索（需要 cookie）
	open.GET("/g79/user_search", func(c *gin.Context) {
		cookie, ok := extractCookieFromAuth(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "authorization is required (use Authorization: cookie:<cookie>)"})
			return
		}
		kw := c.Query("kw")
		if kw == "" {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "kw is required"})
			return
		}
		stype, _ := strconv.Atoi(c.DefaultQuery("type", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		cli, err := g79.NewClient()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "message": err.Error()})
			return
		}
		if err := cli.AuthenticateWithCookie(cookie); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
			return
		}
		res, err := cli.SearchUserByNameOrMail(kw, stype, limit)
		if err != nil || res == nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "search failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "result": res})
	})

	// 组件下载信息（需要 cookie）
	open.GET("/g79/download_info", func(c *gin.Context) {
		cookie, ok := extractCookieFromAuth(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "authorization is required (use Authorization: cookie:<cookie>)"})
			return
		}
		itemID := c.Query("item_id")
		if itemID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "item_id is required"})
			return
		}
		cli, err := g79.NewClient()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "message": err.Error()})
			return
		}
		if err := cli.AuthenticateWithCookie(cookie); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
			return
		}
		info, err := cli.GetDownloadInfo(itemID)
		if err != nil || info == nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "query failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "result": info})
	})

	// 在线大厅房间信息（保留，只读，需要 cookie）
	open.GET("/g79/lobby_room", func(c *gin.Context) {
		cookie, ok := extractCookieFromAuth(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "authorization is required (use Authorization: cookie:<cookie>)"})
			return
		}
		id := c.Query("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id is required"})
			return
		}
		cli, err := g79.NewClient()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "message": err.Error()})
			return
		}
		if err := cli.AuthenticateWithCookie(cookie); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
			return
		}
		roomInfo, err := cli.GetOnlineLobbyRoom(id)
		if err != nil || roomInfo == nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "query failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "result": roomInfo})
	})
}
