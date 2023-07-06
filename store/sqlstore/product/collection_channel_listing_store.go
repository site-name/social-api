package product

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/store_iface"
)

type SqlCollectionChannelListingStore struct {
	store.Store
}

func NewSqlCollectionChannelListingStore(s store.Store) store.CollectionChannelListingStore {
	return &SqlCollectionChannelListingStore{s}
}

func (s *SqlCollectionChannelListingStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id", "CreateAt", "CollectionID", "ChannelID", "PublicationDate", "IsPublished",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, item string) string {
		return prefix + item
	})
}

func (s *SqlCollectionChannelListingStore) FilterByOptions(options *model.CollectionChannelListingFilterOptions) ([]*model.CollectionChannelListing, error) {
	query := s.GetQueryBuilder().Select(s.ModelFields(store.CollectionChannelListingTableName + ".")...).From(store.CollectionChannelListingTableName)

	for _, opt := range []squirrel.Sqlizer{options.Id, options.ChannelID, options.CollectionID} {
		if opt != nil {
			query = query.Where(opt)
		}
	}

	queryStr, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.CollectionChannelListing
	err = s.GetReplicaX().Select(&res, queryStr, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find collection channel listings by given options")
	}

	return res, nil
}

func (s *SqlCollectionChannelListingStore) Delete(transaction store_iface.SqlxTxExecutor, options *model.CollectionChannelListingFilterOptions) error {
	query := s.GetQueryBuilder().Delete(store.CollectionChannelListingTableName)

	for _, opt := range []squirrel.Sqlizer{options.Id, options.ChannelID, options.CollectionID} {
		if opt != nil {
			query = query.Where(opt)
		}
	}

	queryStr, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "Delete_ToSql")
	}

	runner := s.GetMasterX()
	if transaction != nil {
		runner = transaction
	}

	_, err = runner.Exec(queryStr, args...)
	if err != nil {
		return errors.Wrap(err, "failed to delete collection channel listing relations")
	}

	return nil
}

func (s *SqlCollectionChannelListingStore) Upsert(transaction store_iface.SqlxTxExecutor, relations ...*model.CollectionChannelListing) ([]*model.CollectionChannelListing, error) {
	saveQuery := "INSERT INTO " + store.CollectionChannelListingTableName + "(" + s.ModelFields("").Join(",") + ") VALUES (" + s.ModelFields(":").Join(",") + ")"
	updateQuery := "UPDATE " + store.CollectionChannelListingTableName + " SET " + s.ModelFields("").Map(func(_ int, item string) string { return item + ":=" + item }).Join(",") + " WHERE Id=:Id"
	runner := s.GetMasterX()
	if transaction != nil {
		runner = transaction
	}

	for _, rel := range relations {
		isSaving := false

		if rel.Id == "" {
			rel.PreSave()
			isSaving = true
		}

		if err := rel.IsValid(); err != nil {
			return nil, err
		}

		var (
			result sql.Result
			err    error
		)
		if isSaving {
			result, err = runner.NamedExec(saveQuery, rel)
		} else {
			result, err = runner.NamedExec(updateQuery, rel)
		}

		if err != nil {
			if s.IsUniqueConstraintError(err, []string{"CollectionID", "ChannelID", "collectionchannellistings_collectionid_channelid_key"}) {
				return nil, store.NewErrInvalidInput("CollectionChannelListings", "collectionID/channelID", "duplicate")
			}
			return nil, errors.Wrap(err, "failed to upsert collection channel listing relation")
		}

		numUpserted, _ := result.RowsAffected()
		if numUpserted != 1 {
			return nil, errors.Errorf("%d relation upserted instead of 1", numUpserted)
		}
	}

	return relations, nil
}
