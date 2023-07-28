package shipping

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlShippingMethodPostalCodeRuleStore struct {
	store.Store
}

func NewSqlShippingMethodPostalCodeRuleStore(s store.Store) store.ShippingMethodPostalCodeRuleStore {
	return &SqlShippingMethodPostalCodeRuleStore{s}
}

func (s *SqlShippingMethodPostalCodeRuleStore) ScanFields(rule *model.ShippingMethodPostalCodeRule) []interface{} {
	return []interface{}{
		&rule.Id,
		&rule.ShippingMethodID,
		&rule.Start,
		&rule.End,
		&rule.InclusionType,
	}
}

func (s *SqlShippingMethodPostalCodeRuleStore) FilterByOptions(options *model.ShippingMethodPostalCodeRuleFilterOptions) ([]*model.ShippingMethodPostalCodeRule, error) {
	var res []*model.ShippingMethodPostalCodeRule
	err := s.GetReplica().Find(&res, store.BuildSqlizer(options.Conditions)...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find shipping method postal code rules by given options")
	}

	return res, nil
}

func (s *SqlShippingMethodPostalCodeRuleStore) Delete(transaction *gorm.DB, ids ...string) error {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	err := transaction.Raw("DELETE FROM "+model.ShippingMethodPostalCodeRuleTableName+" WHERE Id IN ?", ids).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete shipping method postal code rules")
	}

	return nil
}

func (s *SqlShippingMethodPostalCodeRuleStore) Save(transaction *gorm.DB, rules model.ShippingMethodPostalCodeRules) (model.ShippingMethodPostalCodeRules, error) {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	for _, rule := range rules {
		err := transaction.Create(rule).Error
		if err != nil {
			if s.IsUniqueConstraintError(err, []string{"shippingmethodid_start_end_key", "Start", "End", "ShippingMethodID"}) {
				return nil, store.NewErrInvalidInput(model.ShippingMethodPostalCodeRuleTableName, "", "")
			}
			return nil, errors.Wrap(err, "failed to save shipping method postal code rule")
		}
	}

	return rules, nil
}
