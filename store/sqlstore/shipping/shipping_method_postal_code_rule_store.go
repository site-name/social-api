package shipping

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlShippingMethodPostalCodeRuleStore struct {
	store.Store
}

func NewSqlShippingMethodPostalCodeRuleStore(s store.Store) store.ShippingMethodPostalCodeRuleStore {
	return &SqlShippingMethodPostalCodeRuleStore{s}
}

func (s *SqlShippingMethodPostalCodeRuleStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id",
		"ShippingMethodID",
		"Start",
		"End",
		"InclusionType",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}
