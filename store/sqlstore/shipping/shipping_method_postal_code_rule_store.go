package shipping

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/store"
)

type SqlShippingMethodPostalCodeRuleStore struct {
	store.Store
}

func NewSqlShippingMethodPostalCodeRuleStore(s store.Store) store.ShippingMethodPostalCodeRuleStore {
	smls := &SqlShippingMethodPostalCodeRuleStore{s}
	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(shipping.ShippingMethodPostalCodeRule{}, "ShippingMethodPostalCodeRules").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ShippingMethodID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Start").SetMaxSize(shipping.SHIPPING_METHOD_POSTAL_CODE_RULE_COMMON_MAX_LENGTH)
		table.ColMap("End").SetMaxSize(shipping.SHIPPING_METHOD_POSTAL_CODE_RULE_COMMON_MAX_LENGTH)
		table.ColMap("InclusionType").SetMaxSize(shipping.SHIPPING_METHOD_POSTAL_CODE_RULE_COMMON_MAX_LENGTH).
			SetDefaultConstraint(model.NewString(shipping.EXCLUDE))

		table.SetUniqueTogether("ShippingMethodID", "Start", "End")
	}
	return smls
}

func (s *SqlShippingMethodPostalCodeRuleStore) CreateIndexesIfNotExists() {
	s.CreateIndexIfNotExists("idx_shipping_method_postal_code_rules_start", "ShippingMethodPostalCodeRules", "Start")
	// s.CreateIndexIfNotExists("idx_shipping_method_postal_code_rules_end", "ShippingMethodPostalCodeRules", "End")
	s.CreateIndexIfNotExists("idx_shipping_method_postal_code_rules_inclusion_type", "ShippingMethodPostalCodeRules", "InclusionType")
}