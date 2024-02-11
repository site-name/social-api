package model_helper

type GiftCardSettingsExpiryType string

const (
	NEVER_EXPIRE  GiftCardSettingsExpiryType = "never_expire"
	EXPIRY_PERIOD GiftCardSettingsExpiryType = "expiry_period"
)

func (g GiftCardSettingsExpiryType) IsValid() bool {
	switch g {
	case NEVER_EXPIRE, EXPIRY_PERIOD:
		return true
	default:
		return false
	}
}

type ShopStaffFilterOptions struct {
	CommonQueryOptions
}
