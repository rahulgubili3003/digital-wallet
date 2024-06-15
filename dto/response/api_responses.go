package response

type Wallets struct {
	WalletId uint    `json:"wallet_id"`
	Balance  float64 `json:"balance"`
}

type WalletResponse struct {
	Wallets []Wallets `json:"wallet_info"`
}
