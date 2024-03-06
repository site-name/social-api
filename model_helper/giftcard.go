package model_helper

import (
	"net/http"

	"github.com/site-name/decimal"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func GiftcardEventIsValid(g model.GiftcardEvent) *AppError {
	if !IsValidId(g.ID) {
		return NewAppError("GiftcardEventIsValid", "model.giftcard_event.is_valid.id.app_error", nil, "invalid id", http.StatusBadRequest)
	}
	if !IsValidId(g.GiftcardID) {
		return NewAppError("GiftcardEventIsValid", "model.giftcard_event.is_valid.giftcard_id.app_error", nil, "invalid giftcard id", http.StatusBadRequest)
	}
	if g.Date < 0 {
		return NewAppError("GiftcardEventIsValid", "model.giftcard_event.is_valid.created_at.app_error", nil, "invalid created at", http.StatusBadRequest)
	}
	if g.Type.IsValid() != nil {
		return NewAppError("GiftcardEventIsValid", "model.giftcard_event.is_valid.type.app_error", nil, "please provide valid type", http.StatusBadRequest)
	}
	if !g.UserID.IsNil() && !IsValidId(*g.UserID.String) {
		return NewAppError("GiftcardEventIsValid", "model.giftcard_event.is_valid.user_id.app_error", nil, "invalid giftcard id", http.StatusBadRequest)
	}
	return nil
}

func GiftCardEventPreSave(g *model.GiftcardEvent) {
	if g.ID == "" {
		g.ID = NewId()
	}
	if g.Type.IsValid() != nil {
		g.Type = model.GiftcardEventTypeActivated
	}
}

type GiftCardEventFilterOption struct {
	CommonQueryOptions
}

var GiftcardAnnotationKeys = struct {
	RelatedProductName     string
	RelatedProductSlug     string
	RelatedUsedByFirstName string
	RelatedUsedBylastName  string
}{
	RelatedProductName:     "related_product_name",
	RelatedProductSlug:     "related_product_slug",
	RelatedUsedByFirstName: "related_used_by_first_name",
	RelatedUsedBylastName:  "related_used_by_last_name",
}

type GiftcardFilterOption struct {
	CommonQueryOptions
	CheckoutToken                       qm.QueryMod // INNER JOIN giftcard_checkouts ON ... WHERE giftcard_checkouts.checkout_id ...
	OrderID                             qm.QueryMod // INNER JOIN order_giftcards ON ... WHERE order_gifcards.order_id ...
	AnnotateRelatedProductNameAndSlug   bool
	AnnotateUsedByFirstNameAndLastNames bool
}

func GiftcardPreSave(g *model.Giftcard) {
	if g.ID == "" {
		g.ID = NewId()
	}
	if g.CreatedAt == 0 {
		g.CreatedAt = GetMillis()
	}
	GiftcardCommonPre(g)
}

func GiftcardCommonPre(g *model.Giftcard) {
	if g.Currency.IsValid() != nil {
		g.Currency = DEFAULT_CURRENCY
	}
	if g.CurrentBalanceAmount.IsNil() {
		g.CurrentBalanceAmount = model_types.NewNullDecimal(decimal.Zero)
	}
	if g.InitialBalanceAmount.IsNil() {
		g.InitialBalanceAmount = model_types.NewNullDecimal(decimal.Zero)
	}
	if g.IsActive.IsNil() {
		g.IsActive = model_types.NewNullBool(true)
	}
	if g.StartDate.IsNil() {
		g.StartDate = model_types.NewNullTime(GetTimeUTCNow())
	}
	if g.Code == "" {
		g.Code = NewPromoCode()
	}
}

func GiftcardIsValid(g model.Giftcard) *AppError {
	if !IsValidId(g.ID) {
		return NewAppError("GiftcardIsValid", "model.giftcard.is_valid.id.app_error", nil, "invalid id", http.StatusBadRequest)
	}
	if g.CreatedAt <= 0 {
		return NewAppError("GiftcardIsValid", "model.giftcard.is_valid.created_at.app_error", nil, "invalid created at", http.StatusBadRequest)
	}
	if g.Currency.IsValid() != nil {
		return NewAppError("GiftcardIsValid", "model.giftcard.is_valid.currency.app_error", nil, "invalid currency", http.StatusBadRequest)
	}
	if g.CurrentBalanceAmount.IsNil() {
		return NewAppError("GiftcardIsValid", "model.giftcard.is_valid.current_balance_amount.app_error", nil, "invalid current balance amount", http.StatusBadRequest)
	}
	if g.InitialBalanceAmount.IsNil() {
		return NewAppError("GiftcardIsValid", "model.giftcard.is_valid.initial_balance_amount.app_error", nil, "invalid initial balance amount", http.StatusBadRequest)
	}
	if g.IsActive.IsNil() {
		return NewAppError("GiftcardIsValid", "model.giftcard.is_valid.is_active.app_error", nil, "invalid is active", http.StatusBadRequest)
	}
	if g.StartDate.IsNil() {
		return NewAppError("GiftcardIsValid", "model.giftcard.is_valid.start_date.app_error", nil, "invalid start date", http.StatusBadRequest)
	}
	if !PromoCodeRegex.MatchString(g.Code) {
		return NewAppError("GiftcardIsValid", "model.giftcard.is_valid.code.app_error", nil, "invalid code", http.StatusBadRequest)
	}
	if !g.CreatedByID.IsNil() && !IsValidId(*g.CreatedByID.String) {
		return NewAppError("GiftcardIsValid", "model.giftcard.is_valid.created_by_id.app_error", nil, "invalid created by id", http.StatusBadRequest)
	}
	if !g.UsedByID.IsNil() && !IsValidId(*g.UsedByID.String) {
		return NewAppError("GiftcardIsValid", "model.giftcard.is_valid.used_by_id.app_error", nil, "invalid used by id", http.StatusBadRequest)
	}
	if !g.CreatedByEmail.IsNil() && !IsValidEmail(*g.CreatedByEmail.String) {
		return NewAppError("GiftcardIsValid", "model.giftcard.is_valid.created_by_email.app_error", nil, "invalid created by email", http.StatusBadRequest)
	}
	if !g.UsedByEmail.IsNil() && !IsValidEmail(*g.UsedByEmail.String) {
		return NewAppError("GiftcardIsValid", "model.giftcard.is_valid.used_by_email.app_error", nil, "invalid used by email", http.StatusBadRequest)
	}
	if !g.ProductID.IsNil() && !IsValidId(*g.ProductID.String) {
		return NewAppError("GiftcardIsValid", "model.giftcard.is_valid.product_id.app_error", nil, "invalid product id", http.StatusBadRequest)
	}
	if !g.LastUsedOn.IsNil() && *g.LastUsedOn.Int64 <= 0 {
		return NewAppError("GiftcardIsValid", "model.giftcard.is_valid.last_used_on.app_error", nil, "invalid last used on", http.StatusBadRequest)
	}
	if !g.StartDate.IsNil() && g.StartDate.Time.IsZero() {
		return NewAppError("GiftcardIsValid", "model.giftcard.is_valid.start_date.app_error", nil, "invalid start date", http.StatusBadRequest)
	}
	if !g.ExpiryDate.IsNil() && g.ExpiryDate.Time.IsZero() {
		return NewAppError("GiftcardIsValid", "model.giftcard.is_valid.expiry_date.app_error", nil, "invalid expiry date", http.StatusBadRequest)
	}
	return nil
}
