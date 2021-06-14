package model

import (
	"io"
)

type MfaSecret struct {
	Secret string `json:"secret"`
	QRCode string `json:"qr_code"`
}

func (mfa *MfaSecret) ToJson() string {
	return ModelToJson(mfa)
}

func MfaSecretFromJson(data io.Reader) *MfaSecret {
	var mfa *MfaSecret
	ModelFromJson(&mfa, data)
	return mfa
}
