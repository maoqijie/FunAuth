package handlers

type LoginRequest struct {
	FBToken         string `json:"login_token,omitempty"`
	UserName        string `json:"username,omitempty"`
	Password        string `json:"password,omitempty"`
	ServerCode      string `json:"server_code"`
	ServerPassword  string `json:"server_passcode"`
	ClientPublicKey string `json:"client_public_key"`
}

type SkinInfo struct {
	ItemID          string `json:"entity_id"`
	SkinDownloadURL string `json:"res_url"`
	SkinIsSlim      bool   `json:"is_slim"`
}

type LoginResponse struct {
	SuccessStates  bool            `json:"success"`
	ServerMessage  string          `json:"server_msg,omitempty"`
	BotLevel       int             `json:"growth_level"`
	BotSkin        SkinInfo        `json:"skin_info,omitempty"`
	BotComponent   map[string]*int `json:"outfit_info,omitempty"`
	FBToken        string          `json:"token"`
	RentalServerIP string          `json:"ip_address"`
	ChainInfo      string          `json:"chainInfo"`
}

type TransferCheckNumRequest struct {
	Data string `json:"data"`
}

type TransferCheckNumResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Value   string `json:"value,omitempty"`
}

type TransferStartTypeQuery struct {
	Content string `form:"content"`
}

type TransferStartTypeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    string `json:"data,omitempty"`
}

// TanLobbyTransferServersResponse ..
type TanLobbyTransferServersResponse struct {
	Success          bool     `json:"success"`
	ErrorInfo        string   `json:"error_info"`
	RaknetServers    []string `json:"raknet_servers"`
	WebsocketServers []string `json:"websocket_servers"`
}

// TanLobbyLoginRequest ..
type TanLobbyLoginRequest struct {
	FBToken      string `json:"login_token"`
	RoomID       string `json:"room_id"`
}

// TanLobbyLoginResponse ..
type TanLobbyLoginResponse struct {
	Success   bool   `json:"success"`
	ErrorInfo string `json:"error_info"`

	RoomOwnerID    uint32 `json:"room_owner_id"`
	UserUniqueID   uint32 `json:"user_unique_id"`
	UserPlayerName string `json:"user_player_name"`

	RaknetRand      []byte `json:"raknet_rand"`
	RaknetAESRand   []byte `json:"raknet_aes_rand"`
	EncryptKeyBytes []byte `json:"encrypt_key_bytes"`
	DecryptKeyBytes []byte `json:"decrypt_key_bytes"`

	SignalingSeed   []byte `json:"signaling_seed"`
	SignalingTicket []byte `json:"signaling_ticket"`
	
	BotLevel       int             `json:"growth_level"`
	BotSkin        SkinInfo        `json:"skin_info,omitempty"`
	BotComponent   map[string]*int `json:"outfit_info,omitempty"`
}
