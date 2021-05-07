package wishlist

type Wishlish struct {
	Id       string  `json:"id"`
	Token    string  `json:"token"`
	UserID   *string `json:"user_id"`
	CreateAt int64   `json:"create_at"`
}
