package shipping

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlShippingMethodPostalCodeRuleStore struct {
	store.Store
}

func NewSqlShippingMethodPostalCodeRuleStore(s store.Store) store.ShippingMethodPostalCodeRuleStore {
	return &SqlShippingMethodPostalCodeRuleStore{s}
}

func (s *SqlShippingMethodPostalCodeRuleStore) Save(transaction boil.ContextTransactor, rules model.ShippingMethodPostalCodeRuleSlice) (model.ShippingMethodPostalCodeRuleSlice, error) {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	for _, rule := range rules {
		if rule == nil {
			continue
		}

		isSaving := rule.ID == ""
		if isSaving {
			model_helper.ShippingMethodPostalCodeRulePreSave(rule)
		}

		if err := model_helper.ShippingMethodPostalCodeRuleIsValid(*rule); err != nil {
			return nil, err
		}

		var err error
		if isSaving {
			err = rule.Insert(transaction, boil.Infer())
		} else {
			_, err = rule.Update(transaction, boil.Infer())
		}

		if err != nil {
			if s.IsUniqueConstraintError(err, []string{"shipping_method_postal_code_rules_shipping_method_id_start_end_key", model.ShippingMethodPostalCodeRuleColumns.End, model.ShippingMethodPostalCodeRuleColumns.Start, model.ShippingMethodPostalCodeRuleColumns.ShippingMethodID}) {
				return nil, store.NewErrInvalidInput(model.TableNames.ShippingMethodPostalCodeRules, "", "duplicate rule")
			}
			return nil, err
		}
	}

	return rules, nil
}

func (s *SqlShippingMethodPostalCodeRuleStore) FilterByOptions(options model_helper.ShippingMethodPostalCodeRuleFilterOptions) (model.ShippingMethodPostalCodeRuleSlice, error) {
	conds := options.Conditions
	return model.ShippingMethodPostalCodeRules(conds...).All(s.GetReplica())
}

func (s *SqlShippingMethodPostalCodeRuleStore) Delete(transaction boil.ContextTransactor, ids []string) error {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	_, err := model.ShippingMethodPostalCodeRules(
		model.ShippingMethodPostalCodeRuleWhere.ID.IN(ids),
	).DeleteAll(transaction)
	return err
}
