package shipping

import (
	"fmt"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/sitename/sitename/store"
	"github.com/stretchr/testify/require"
)

func applicableShippingMethodsByChannel(channelID string) squirrel.SelectBuilder {
	subQuery := squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("A", "B", fmt.Sprintf("%s.PriceAmount", store.ShippingMethodChannelListingTableName)).
		From(store.ShippingMethodChannelListingTableName).
		Where(squirrel.And{
			squirrel.Eq{fmt.Sprintf("%s.ChannelID", store.ShippingMethodChannelListingTableName): channelID},
			squirrel.Eq{fmt.Sprintf("%s.ShippingMethodID", store.ShippingMethodChannelListingTableName): fmt.Sprintf("%s.Id", store.ShippingMethodTableName)},
		})

	return subQuery
}

func Test_applicablePriceBasedMethods(t *testing.T) {

	query := applicableShippingMethodsByChannel(uuid.NewString())

	str, args, err := query.ToSql()
	require.NoError(t, err)

	fmt.Println(str)
	fmt.Println(args...)
}
