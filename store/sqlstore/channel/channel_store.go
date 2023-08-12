package channel

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlChannelStore struct {
	store.Store
}

func NewSqlChannelStore(sqlStore store.Store) store.ChannelStore {
	return &SqlChannelStore{sqlStore}
}

func (cs *SqlChannelStore) ScanFields(ch *model.Channel) []interface{} {
	return []interface{}{
		&ch.Id,
		&ch.Name,
		&ch.IsActive,
		&ch.Slug,
		&ch.Currency,
		&ch.DefaultCountry,
	}
}

func (s *SqlChannelStore) Upsert(transaction *gorm.DB, channel *model.Channel) (*model.Channel, error) {
	err := transaction.Save(channel).Error
	if err != nil {
		if s.IsUniqueConstraintError(err, []string{"slug", "slug_unique_key", "idx_channels_slug_unique"}) {
			return nil, store.NewErrInvalidInput(model.ChannelTableName, "Slug", channel.Slug)
		}
		return nil, errors.Wrap(err, "failed to upsert channel")
	}
	return channel, nil
}

func (cs *SqlChannelStore) Get(id string) (*model.Channel, error) {
	var channel model.Channel

	err := cs.GetReplica().First(&channel, "Id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.ChannelTableName, id)
		}
		return nil, errors.Wrapf(err, "Failed to get Channel with ChannelID=%s", id)
	}

	return &channel, nil
}

// FilterByOption returns a list of channels with given option
func (cs *SqlChannelStore) FilterByOption(option *model.ChannelFilterOption) ([]*model.Channel, error) {
	selectFields := []string{model.ChannelTableName + ".*"}
	if option.AnnotateHasOrders {
		selectFields = append(selectFields, `EXISTS ( SELECT (1) AS "a" FROM Orders WHERE Orders.ChannelID = Channels.Id LIMIT 1 ) AS HasOrders`)
	}

	query := cs.GetQueryBuilder().
		Select(selectFields...).
		From(model.ChannelTableName).
		Where(option.Conditions)

	// parse options
	if option.ShippingZoneChannels_ShippingZoneID != nil {
		query = query.
			InnerJoin(model.ShippingZoneChannelTableName + " ON ShippingZoneChannels.ChannelID = Channels.Id").
			Where(option.ShippingZoneChannels_ShippingZoneID)
	}
	if option.VoucherChannelListing_VoucherID != nil {
		query = query.
			InnerJoin(model.VoucherChannelListingTableName + " ON VoucherChannelListings.ChannelID = Channels.Id").
			Where(option.VoucherChannelListing_VoucherID)
	}
	if option.Limit > 0 {
		query = query.Limit(uint64(option.Limit))
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res model.Channels

	rows, err := cs.GetReplica().Raw(queryString, args...).Rows()
	if err != nil {
		return nil, errors.Wrap(err, "failed to find channels with given option")
	}
	defer rows.Close()

	for rows.Next() {
		var (
			hasOrder   bool
			channel    model.Channel
			scanFields = cs.ScanFields(&channel)
		)
		if option.AnnotateHasOrders {
			scanFields = append(scanFields, &hasOrder)
		}

		err = rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan channel row")
		}

		if option.AnnotateHasOrders {
			channel.SetHasOrders(hasOrder)
		}
		res = append(res, &channel)
	}

	return res, nil
}

func (s *SqlChannelStore) DeleteChannels(transaction *gorm.DB, ids []string) error {
	if transaction == nil {
		transaction = s.GetMaster()
	}
	return transaction.Raw("DELETE FROM "+model.ChannelTableName+" WHERE Id IN ?", ids).Error
}
