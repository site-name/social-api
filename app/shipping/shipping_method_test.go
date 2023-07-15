package shipping

import (
	"fmt"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/sitename/sitename/model"
	"github.com/stretchr/testify/require"
)

func applicableShippingMethodsByChannel(channelID string) squirrel.SelectBuilder {
	subQuery := squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("A", "B", fmt.Sprintf("%s.PriceAmount", model.ShippingMethodChannelListingTableName)).
		From(model.ShippingMethodChannelListingTableName).
		Where(squirrel.And{
			squirrel.Eq{fmt.Sprintf("%s.ChannelID", model.ShippingMethodChannelListingTableName): channelID},
			squirrel.Eq{fmt.Sprintf("%s.ShippingMethodID", model.ShippingMethodChannelListingTableName): fmt.Sprintf("%s.Id", model.ShippingMethodTableName)},
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
