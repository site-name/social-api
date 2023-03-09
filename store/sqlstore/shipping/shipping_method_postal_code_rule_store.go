package shipping

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlShippingMethodPostalCodeRuleStore struct {
	store.Store
}

func NewSqlShippingMethodPostalCodeRuleStore(s store.Store) store.ShippingMethodPostalCodeRuleStore {
	return &SqlShippingMethodPostalCodeRuleStore{s}
}

func (s *SqlShippingMethodPostalCodeRuleStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
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
	query := s.GetQueryBuilder().Select("*").From(store.ShippingMethodPostalCodeRuleTableName)

	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.ShippingMethodID != nil {
		query = query.Where(options.ShippingMethodID)
	}

	queryStr, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.ShippingMethodPostalCodeRule
	err = s.GetReplicaX().Select(&res, queryStr, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find shipping method postal code rules by given options")
	}

	return res, nil
}
