package model

type FailedPayment struct {
	CardBrand      string `json:"card_brand"`
	LastFour       int    `json:"last_four"`
	FailureMessage string `json:"failure_message"`
}
