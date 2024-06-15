package request

type TopUpRequest struct {
	UserId      uint    `json:"user_id"`
	TopUpAmount float64 `json:"top_up_amount"`
}
