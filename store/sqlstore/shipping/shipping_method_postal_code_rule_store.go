package shipping

import (
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/store"
)

type SqlShippingMethodPostalCodeRuleStore struct {
	store.Store
}

func NewSqlShippingMethodPostalCodeRuleStore(s store.Store) store.ShippingMethodPostalCodeRuleStore {
	smls := &SqlShippingMethodPostalCodeRuleStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(shipping.ShippingMethodPostalCodeRule{}, store.ShippingMethodPostalCodeRuleTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingMethodID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Start").SetMaxSize(shipping.SHIPPING_METHOD_POSTAL_CODE_RULE_COMMON_MAX_LENGTH)
		table.ColMap("End").SetMaxSize(shipping.SHIPPING_METHOD_POSTAL_CODE_RULE_COMMON_MAX_LENGTH)
		table.ColMap("InclusionType").SetMaxSize(shipping.SHIPPING_METHOD_POSTAL_CODE_RULE_COMMON_MAX_LENGTH)

		table.SetUniqueTogether("ShippingMethodID", "Start", "End")
	}
	return smls
}

func (s *SqlShippingMethodPostalCodeRuleStore) CreateIndexesIfNotExists() {
	s.CreateForeignKeyIfNotExists(store.ShippingMethodPostalCodeRuleTableName, "ShippingMethodID", store.ShippingMethodTableName, "Id", true)
}

func (s *SqlShippingMethodPostalCodeRuleStore) ModelFields() []string {
	return []string{
		"ShippingMethodPostalCodeRules.Id",
		"ShippingMethodPostalCodeRules.ShippingMethodID",
		"ShippingMethodPostalCodeRules.Start",
		"ShippingMethodPostalCodeRules.End",
		"ShippingMethodPostalCodeRules.InclusionType",
	}
}
