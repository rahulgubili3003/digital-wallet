package request

type TransferAmount struct {
	UserId          uint    `json:"user_id"`
	Amount          float64 `json:"amount"`
	RecipientUserId uint    `json:"recipient_user_id"`
}
