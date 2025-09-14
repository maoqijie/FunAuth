package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"strconv"

	"github.com/Yeah114/g79client/utils"
	"github.com/Yeah114/g79client"
)

func TanLobbyLogin(ctx context.Context, cli *g79client.Client, p TanLobbyLoginParams) (TanLobbyLoginResult, error) {
	var result TanLobbyLoginResult

	roomInfo, err := cli.GetTransferRoomWithName(p.RoomID)
	if err != nil {
		return result, fmt.Errorf("get transfer room with name: %w", err)
	}
	if roomInfo.Code != 0 {
		return result, fmt.Errorf("get transfer room with name: %s(%d)", roomInfo.Message, roomInfo.Code)
	}
	if len(roomInfo.List) == 0 {
		return result, fmt.Errorf("room not found")
	}

	encryptedToken := utils.GetEncryptedToken(cli.UserToken)
	raknetRand := make([]byte, 16)
	_, err = rand.Read(raknetRand)
	if err != nil {
		return result, fmt.Errorf("rand read: %w", err)
	}
	raknetAESRand, err := utils.AesECBEncrypt(raknetRand, encryptedToken)
	if err != nil {
		return result, fmt.Errorf("aes encrypt: %w", err)
	}
	encryptKeyBytes := append(encryptedToken, raknetRand...)
	decryptKeyBytes := append(raknetRand, encryptedToken...)

	seed := make([]byte, 16)
	_, err = rand.Read(seed)
	if err != nil {
		return result, fmt.Errorf("rand read: %w", err)
	}

	ticket, err := utils.AesECBEncrypt(seed, []byte(cli.UserToken))
	if err != nil {
		return result, fmt.Errorf("aes encrypt: %w", err)
	}

	result.RoomOwnerID = uint32(roomInfo.List[0].HID.Int64())
	result.UserPlayerName = cli.UserDetail.Name
	userUniqueID, err := strconv.ParseInt(cli.UserID, 10, 64)
	if err != nil {
		return result, fmt.Errorf("parse int: %w", err)
	}
	result.UserUniqueID = uint32(userUniqueID)
	result.RaknetRand = raknetRand
	result.RaknetAESRand = raknetAESRand
	result.SignalingSeed = seed
	result.SignalingTicket = ticket
	result.EncryptKeyBytes = encryptKeyBytes
	result.DecryptKeyBytes = decryptKeyBytes

	return result, nil
}
