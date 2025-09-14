package auth

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	g79 "github.com/Yeah114/g79client"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// Login
func Login(ctx context.Context, cli *g79.Client, p LoginParams) (LoginResult, error) {
	var result LoginResult

	if cli == nil {
		return result, fmt.Errorf("nil client")
	}

	// 确保用户详情可用，用于昵称与等级
	if cli.UserDetail == nil {
		detail, err := cli.GetUserDetail()
		if err != nil {
			return result, fmt.Errorf("GetUserDetail: %w", err)
		}
		cli.UserDetail = &detail.Entity
	}
	if cli.UserDetail != nil && cli.UserDetail.Name == "" {
		name := fmt.Sprintf("SW%06d", random.Intn(1000000))
		if err := cli.UpdateNickname(name); err != nil {
			return result, fmt.Errorf("UpdateNickname: %w", err)
		}
	}

	// IP
	var ipAddress string
	// ChainInfo
	var chainInfoStr string

	if p.ServerCode == "" {
		return result, fmt.Errorf("server code is empty")
	}

	if after, ok := strings.CutPrefix(p.ServerCode, "LobbyGame:"); ok && after != "" {
		// 在线大厅
		roomCode := after

		// 获取房间信息
		roomInfo, err := cli.GetOnlineLobbyRoom(roomCode)
		if err != nil {
			return result, fmt.Errorf("GetOnlineLobbyRoom: %w", err)
		}
		if roomInfo.Code != 0 {
			return result, fmt.Errorf("GetOnlineLobbyRoom: %s(%d)", roomInfo.Message, roomInfo.Code)
		}

		// 购买房间地图
		roomMap, err := cli.PurchaseItem(roomInfo.Entity.ResID.String())
		if err != nil {
			return result, fmt.Errorf("PurchaseItem: %w", err)
		}
		if !(roomMap.Code == 0 || roomMap.Code == 502 || roomMap.Code == 44) {
			return result, fmt.Errorf("PurchaseItem: %s(%d)", roomMap.Message, roomMap.Code)
		}

		// 进入房间
		enterResp, err := cli.EnterOnlineLobbyRoom(roomCode, p.ServerPassword)
		if err != nil {
			return result, fmt.Errorf("EnterOnlineLobbyRoom: %w", err)
		}
		if enterResp.Code != 0 {
			return result, fmt.Errorf("EnterOnlineLobbyRoom: %s(%d)", enterResp.Message, enterResp.Code)
		}

		// 进入房间游戏
		gameEnter, err := cli.OnlineLobbyGameEnter()
		if err != nil {
			return result, fmt.Errorf("OnlineLobbyGameEnter: %w", err)
		}
		if gameEnter.Code != 0 {
			return result, fmt.Errorf("OnlineLobbyGameEnter: %s(%d)", gameEnter.Message, gameEnter.Code)
		}
		ipAddress = fmt.Sprintf("%s:%d", gameEnter.Entity.ServerHost, gameEnter.Entity.ServerPort.Int64())

		// 获取 ChainInfo
		authv2Data, err := cli.GenerateLobbyGameAuthV2(roomCode, p.ClientPublicKey)
		if err != nil {
			return result, fmt.Errorf("GenerateLobbyGameAuthV2: %w", err)
		}
		chainInfo, err := cli.SendAuthV2Request(authv2Data)
		if err != nil {
			return result, fmt.Errorf("SendAuthV2Request: %w", err)
		}
		chainInfoStr = string(chainInfo)
	} else if after, ok := strings.CutPrefix(p.ServerCode, "NetworkGame:"); ok && after != "" {
		gameCode := after

		// 获取网络游戏服务器地址
		serverAddress, err := cli.GetPeGameServerAddress(gameCode)
		if err != nil {
			return result, fmt.Errorf("GetPeGameServerAddress: %w", err)
		}
		if serverAddress.Code != 0 {
			return result, fmt.Errorf("GetPeGameServerAddress: %s(%d)", serverAddress.Message, serverAddress.Code)
		}
		ipAddress = fmt.Sprintf("%s:%d", serverAddress.Entity.IP, serverAddress.Entity.Port.Int64())
		
		// 生成网络游戏认证v2数据
		authv2Data, err := cli.GenerateNetworkGameAuthV2(gameCode, p.ClientPublicKey)
		if err != nil {
			return result, fmt.Errorf("GenerateNetworkGameAuthV2: %w", err)
		}
		chainInfo, err := cli.SendAuthV2Request(authv2Data)
		if err != nil {
			return result, fmt.Errorf("SendAuthV2Request: %w", err)
		}
		chainInfoStr = string(chainInfo)
	} else {
		// 租赁服
		serverCode := p.ServerCode

		// 搜索租赁服
		searchResp, err := cli.SearchRentalServerByName(serverCode)
		if err != nil {
			return result, fmt.Errorf("SearchRentalServerByName: %w", err)
		}
		if searchResp.Code != 0 {
			return result, fmt.Errorf("SearchRentalServerByName: %s(%d)", searchResp.Message, searchResp.Code)
		}
		if len(searchResp.Entities) == 0 {
			return result, fmt.Errorf("SearchRentalServerByName: 找不到服务器")
		}
		serverID := searchResp.Entities[0].EntityID

		// 进入租赁服世界
		enterResp, err := cli.EnterRentalServerWorld(serverID.String(), p.ServerPassword)
		if err != nil {
			return result, fmt.Errorf("EnterRentalServerWorld: %w", err)
		}
		if enterResp.Code != 0 {
			return result, fmt.Errorf("EnterRentalServerWorld: %s(%d)", enterResp.Message, enterResp.Code)
		}
		ipAddress = fmt.Sprintf("%s:%d", enterResp.Entity.McserverHost, enterResp.Entity.McserverPort.Int64())

		// 获取 ChainInfo
		authv2Data, err := cli.GenerateRentalGameAuthV2(serverID.String(), p.ClientPublicKey)
		if err != nil {
			return result, fmt.Errorf("GenerateRentalGameAuthV2: %w", err)
		}
		chainInfo, err := cli.SendAuthV2Request(authv2Data)
		if err != nil {
			return result, fmt.Errorf("SendAuthV2Request: %w", err)
		}
		chainInfoStr = string(chainInfo)
	}

	result.UID = cli.UserID
	result.ChainInfo = chainInfoStr
	result.IP = ipAddress
	result.BotLevel = int(cli.UserDetail.Level.Int64())
	return result, nil
}
