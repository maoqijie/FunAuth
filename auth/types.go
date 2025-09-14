package auth

// LoginParams 定义进入服务器验证所需的参数。
type LoginParams struct {
	ServerCode      string
	ServerPassword  string
	ClientPublicKey string
}

// LoginResult 为登录/进入服务器后的结果。
type LoginResult struct {
	UID       string
	ChainInfo string
	IP        string
	BotLevel  int
}

type SkinInfo struct {
	ItemID          string
	SkinDownloadURL string
	SkinIsSlim      bool
}

type TanLobbyLoginParams struct {
	RoomID       string
}

type TanLobbyLoginResult struct {
	RoomOwnerID    uint32
	UserUniqueID   uint32
	UserPlayerName string

	RaknetRand      []byte
	RaknetAESRand   []byte
	EncryptKeyBytes []byte
	DecryptKeyBytes []byte

	SignalingSeed   []byte
	SignalingTicket []byte
}
