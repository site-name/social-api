package channel

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlChannelShopStore struct {
	store.Store
}

func NewSqlChannelShopStore(s store.Store) store.ChannelShopStore {
	return &SqlChannelShopStore{s}
}

var channelShopModelFields = util.AnyArray[string]{
	"Id",
	"ChannelID",
	"CreateAt",
	"EndAt",
}

func (s *SqlChannelShopStore) ModelFields(prefix string) util.AnyArray[string] {
	if prefix == "" {
		return channelShopModelFields
	}
	return channelShopModelFields.Map(func(_ int, item string) string {
		return prefix + item
	})
}

func (s *SqlChannelShopStore) Save(relation *model.ChannelShopRelation) (*model.ChannelShopRelation, error) {
	relation.PreSave()
	if appErr := relation.IsValid(); appErr != nil {
		return nil, appErr
	}

	_, err := s.GetMaster().NamedExec("INSERT INTO "+model.ChannelShopRelationTableName+"("+s.ModelFields("").Join(",")+") VALUES ("+s.ModelFields(":").Join(",")+")", relation)
	if err != nil {
		if s.IsUniqueConstraintError(err, []string{"ChannelID", "channelshops_shopid_channelid_key"}) {
			return nil, store.NewErrInvalidInput(model.ChannelShopRelationTableName, "channelID / shopID", "")
		}
		return nil, errors.Wrapf(err, "failed to insert channel-shop relation with id=%s", relation.Id)
	}
	return relation, nil
}

func (s *SqlChannelShopStore) FilterByOptions(options *model.ChannelShopRelationFilterOptions) ([]*model.ChannelShopRelation, error) {
	query := s.GetQueryBuilder().
		Select(s.ModelFields(model.ChannelShopRelationTableName + ".")...).
		From(model.ChannelShopRelationTableName)

	if options == nil {
		options = new(model.ChannelShopRelationFilterOptions)
	}
	// parse
	if options.Id != nil {
		query = query.Where(options.Id)
	}

	if options.ChannelID != nil {
		query = query.Where(options.ChannelID)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOptions_ToSql")
	}

	var res []*model.ChannelShopRelation
	err = s.GetReplica().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find channel-shop relations by options")
	}

	return res, nil
}
