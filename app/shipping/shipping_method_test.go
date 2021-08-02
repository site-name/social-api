package shipping

import (
	"fmt"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/site-name/decimal"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/store"
	"github.com/stretchr/testify/require"
)

func Test_applicablePriceBasedMethods(t *testing.T) {
	money, err := goprices.NewMoney(
		model.NewDecimal(decimal.NewFromFloat(34.678)),
		"USD",
	)
	require.NoError(t, err)
	option := applicablePriceBasedMethods(money, model.NewId())

	priceBasedMethodsFilterOption := &model.StringFilter{
		StringOption: &model.StringOption{
			Eq: shipping.PRICE_BASED,
		},
	}

	subQuery := squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("Id").
		From(store.ShippingMethodTableName).
		Where(priceBasedMethodsFilterOption.ToSquirrel("Type"))

	// squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Select("*").
		From(store.ShippingMethodChannelListingTableName).
		FromSelect(subQuery, "ALIAS").
		// Where(option.ShippingMethodID.ToSquirrel("ShippingMethodID")).
		Where(option.MaximumOrderPriceAmount.ToSquirrel("MaximumOrderPriceAmount")).
		Where(option.MinimumOrderPriceAmount.ToSquirrel("MinimumOrderPriceAmount"))

	queryString, args, err := query.ToSql()
	require.NoError(t, err)

	fmt.Println(queryString)
	fmt.Println(args...)
}
