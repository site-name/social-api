package discount

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/json"
)

type VoucherCustomer struct {
	Id            string `json:"id"`
	VoucherID     string `json:"voucher_id"`
	CustomerEmail string `json:"customer_email"`
}

func (vc *VoucherCustomer) createAppError(fieldName string) *model.AppError {
	id := fmt.Sprintf("model.voucher_customer.is_valid.%s.app_error", fieldName)
	var details string
	if !strings.EqualFold(fieldName, "id") {
		details = "voucher_customer_id=" + vc.Id
	}

	return model.NewAppError("VoucherCustomer.IsValid", id, nil, details, http.StatusBadRequest)
}

func (vc *VoucherCustomer) IsValid() *model.AppError {
	if vc.Id == "" {
		return vc.createAppError("id")
	}
	if vc.VoucherID == "" {
		return vc.createAppError("voucher_id")
	}
	if !model.IsValidEmail(vc.CustomerEmail) {
		return vc.createAppError("customer_email")
	}

	return nil
}

func (vc *VoucherCustomer) ToJson() string {
	b, _ := json.JSON.Marshal(vc)
	return string(b)
}

func VoucherCustomerFromJson(data io.Reader) *VoucherCustomer {
	var vc VoucherCustomer
	err := json.JSON.NewDecoder(data).Decode(&vc)
	if err != nil {
		return nil
	}
	return &vc
}
