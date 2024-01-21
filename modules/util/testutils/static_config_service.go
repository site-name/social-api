package testutils

import (
	"crypto/ecdsa"

	"github.com/sitename/sitename/model_helper"
)

type StaticConfigService struct {
	Cfg *model_helper.Config
}

func (s StaticConfigService) Config() *model_helper.Config {
	return s.Cfg
}

func (StaticConfigService) AddConfigListener(func(old, current *model_helper.Config)) string {
	return ""
}

func (StaticConfigService) RemoveConfigListener(string) {

}

func (StaticConfigService) AsymmetricSigningKey() *ecdsa.PrivateKey {
	return &ecdsa.PrivateKey{}
}
func (StaticConfigService) PostActionCookieSecret() []byte {
	return make([]byte, 32)
}
