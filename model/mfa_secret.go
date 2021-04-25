package model

import (
	"io"

	"github.com/sitename/sitename/modules/json"
)

type MfaSecret struct {
	Secret string `json:"secret"`
	QRCode string `json:"qr_code"`
}

func (mfa *MfaSecret) ToJson() string {
	b, _ := json.JSON.Marshal(mfa)
	return string(b)
}

func MfaSecretFromJson(data io.Reader) *MfaSecret {
	var mfa *MfaSecret
	json.JSON.NewDecoder(data).Decode(&mfa)
	return mfa
}
