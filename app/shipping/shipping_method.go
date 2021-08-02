package shipping

import (
	"github.com/Masterminds/squirrel"
	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/shipping"
	"github.com/sitename/sitename/store"
)

func applicableWeightBasedMethods(weight float32) *shipping.ShippingMethodFilterOption {
	return &shipping.ShippingMethodFilterOption{
		Type: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: shipping.WEIGHT_BASED,
			},
		},
		MinimumOrderWeight: &model.NumberFilter{
			Or: &model.NumberOption{
				LtE:  model.NewFloat64(float64(weight)),
				NULL: model.NewBool(true),
			},
		},
		MaximumOrderWeight: &model.NumberFilter{
			Or: &model.NumberOption{
				GtE:  model.NewFloat64(float64(weight)),
				NULL: model.NewBool(true),
			},
		},
	}
}

func applicablePriceBasedMethods(price *goprices.Money, channelID string) *shipping.ShippingMethodChannelListingFilterOption {
	float64PriceAmount, _ := price.Amount.Float64()

	priceBasedMethodsFilterOption := &model.StringFilter{
		StringOption: &model.StringOption{
			Eq: shipping.PRICE_BASED,
		},
	}

	subQuery, args, _ := squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("Id").
		From(store.ShippingMethodTableName).
		Where(priceBasedMethodsFilterOption.ToSquirrel("Type")).ToSql()

	return &shipping.ShippingMethodChannelListingFilterOption{
		ShippingMethodID: &model.StringFilter{
			StringOption: &model.StringOption{
				ExtraExpr: []squirrel.Sqlizer{
					squirrel.Expr("ShippingMethodID IN ("+subQuery+")", args...),
				},
			},
		},
		ChannelID: &model.StringFilter{
			StringOption: &model.StringOption{
				Eq: channelID,
			},
		},
		MinimumOrderPriceAmount: &model.NumberFilter{
			NumberOption: &model.NumberOption{
				LtE: model.NewFloat64(float64PriceAmount),
			},
		},
		MaximumOrderPriceAmount: &model.NumberFilter{
			Or: &model.NumberOption{
				NULL: model.NewBool(true),
				GtE:  model.NewFloat64(float64PriceAmount),
			},
		},
	}
}
