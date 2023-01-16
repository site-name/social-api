package discount

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlVoucherProductVariantStore struct {
	store.Store
}

func NewSqlVoucherProductVariantStore(s store.Store) store.VoucherProductVariantStore {
	return &SqlVoucherProductVariantStore{s}
}

func (s *SqlVoucherProductVariantStore) FilterByOptions(options *model.VoucherProductVariantFilterOption) ([]*model.VoucherProductVariant, error) {
	query := s.GetQueryBuilder().Select("*").From(store.VoucherProductVariantTableName)

	if options.Id != nil {
		query = query.Where(options.Id)
	}
	if options.ProductVariantID != nil {
		query = query.Where(options.ProductVariantID)
	}
	if options.VoucherID != nil {
		query = query.Where(options.VoucherID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.VoucherProductVariant
	err = s.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find voucher product variant relations with given options")
	}

	return res, nil
}
