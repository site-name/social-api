package shipping

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
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

func (s *SqlShippingMethodPostalCodeRuleStore) Delete(transaction store_iface.SqlxTxExecutor, ids ...string) error {
	query, args, err := s.GetQueryBuilder().Delete(store.ShippingMethodPostalCodeRuleTableName).Where(squirrel.Eq{"Id": ids}).ToSql()
	if err != nil {
		return errors.Wrap(err, "Delete_ToSql")
	}

	runner := s.GetMasterX()
	if transaction != nil {
		runner = transaction
	}

	result, err := runner.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to delete shipping method postal code rules")
	}
	numDeleted, _ := result.RowsAffected()
	if int(numDeleted) != len(ids) {
		return errors.Errorf("%d records deleted instead of %d", numDeleted, len(ids))
	}

	return nil
}

func (s *SqlShippingMethodPostalCodeRuleStore) Save(transaction store_iface.SqlxTxExecutor, rules model.ShippingMethodPostalCodeRules) (model.ShippingMethodPostalCodeRules, error) {
	query := "INSERT INTO " + store.ShippingMethodPostalCodeRuleTableName + "(" + s.ModelFields("").Join(",") + ") VALUES (" + s.ModelFields(":").Join(",") + ")"

	runner := s.GetMasterX()
	if transaction != nil {
		runner = transaction
	}

	for _, rule := range rules {
		rule.PreSave()

		if err := rule.IsValid(); err != nil {
			return nil, err
		}

		_, err := runner.NamedExec(query, rule)
		if err != nil {
			if s.IsUniqueConstraintError(err, []string{"shippingmethodpostalcoderules_shippingmethodid_start_end_key", "Start", "End", "ShippingMethodID"}) {
				return nil, store.NewErrInvalidInput(store.ShippingMethodPostalCodeRuleTableName, "", "")
			}
			return nil, errors.Wrap(err, "failed to save shipping method postal code rule")
		}
	}

	return rules, nil
}
