package product

import (
	"net/http"
	"strings"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/channel"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/util"
)

// getTotalDiscountFromRange Calculate the discount amount between two TaxedMoneyRange.
//
// Subtract two prices and return their total discount, if any.
// Otherwise, it returns None.
func getTotalDiscountFromRange(undiscounted *goprices.TaxedMoneyRange, discounted *goprices.TaxedMoneyRange) (*goprices.TaxedMoney, error) {
	return getTotalDiscount(undiscounted.Start, discounted.Start)
}

// getTotalDiscount  Calculate the discount amount between two TaxedMoney.
//
// Subtract two prices and return their total discount, if any.
// Otherwise, it returns None.
func getTotalDiscount(unDiscounted *goprices.TaxedMoney, discounted *goprices.TaxedMoney) (*goprices.TaxedMoney, error) {
	less, err := discounted.LessThan(unDiscounted)
	if err != nil {
		return nil, err
	}
	if less {
		return unDiscounted.Sub(discounted)
	}

	return nil, nil
}

// `PriceRangeLocal` and `DiscountLocalCurrency` can be of 2 format:
// (*MoneyRange, *MoneyRange) or (*TaxedMoneyRange, *TaxedMoneyRange)
type aStructType struct {
	PriceRangeLocal       interface{}
	DiscountLocalCurrency interface{}
}

// getProductPriceRange
//
// NOTE: `discounted`, `unDiscounted` both can be either *MoneyRange or *TaxedMoneyRange
func (a *AppProduct) getProductPriceRange(discounted interface{}, unDiscounted interface{}, localCurrency string) (*aStructType, *model.AppError) {

	// validate `discounted` and `unDiscounted` and `localCurrency`
	// are provided valid and have same currencies
	errorArguments := []string{}
	switch discounted.(type) {
	case *goprices.MoneyRange, *goprices.TaxedMoneyRange:
	default:
		errorArguments = append(errorArguments, "discounted")
	}
	switch unDiscounted.(type) {
	case *goprices.MoneyRange, *goprices.TaxedMoneyRange:
	default:
		errorArguments = append(errorArguments, "unDiscounted")
	}
	// validate they go in pair like:
	// (*MoneyRange, *MoneyRange) or (*TaxedMoneyRange, *TaxedMoneyRange)
	switch v := discounted.(type) {
	case *goprices.MoneyRange:
		if t, ok := unDiscounted.(*goprices.MoneyRange); !ok {
			errorArguments = append(errorArguments, "unDiscounted.(type) != discounted.(type)")
		} else if !strings.EqualFold(t.Currency, v.Currency) {
			errorArguments = append(errorArguments, "unDiscounted.Currency != discounted.Currency")
		}
	case *goprices.TaxedMoneyRange:
		if t, ok := unDiscounted.(*goprices.TaxedMoneyRange); !ok {
			errorArguments = append(errorArguments, "unDiscounted.(type) != discounted.(type)")
		} else if !strings.EqualFold(v.Currency, t.Currency) {
			errorArguments = append(errorArguments, "unDiscounted.Currency != discounted.Currency")
		}
	}

	if len(errorArguments) > 0 {
		return nil, model.NewAppError("getProductPriceRange", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": strings.Join(errorArguments, ", ")}, "", http.StatusBadRequest)
	}

	localCurrency = strings.ToUpper(localCurrency)

	var (
		priceRangeLocal       interface{}
		discountLocalCurrency interface{}
	)
	// validate provided currency is calid
	if goprices.CurrenciesMap[localCurrency] != "" {
		priceRangeLocal, appErr := a.ToLocalCurrency(discounted, localCurrency)
		if appErr != nil {
			return nil, appErr
		}
		unDiscountedLocal, appErr := a.ToLocalCurrency(unDiscounted, localCurrency)
		if appErr != nil {
			return nil, appErr
		}

		if unDiscountedLocal != nil {
			switch t := priceRangeLocal.(type) {
			case *goprices.MoneyRange:
				unDiscountedLocalValue := unDiscountedLocal.(*goprices.MoneyRange)
				if less, err := t.Start.LessThan(unDiscountedLocalValue.Start); err == nil && less {
					discountLocalCurrency, _ = unDiscountedLocalValue.Start.Sub(t.Start)
				}

			case *goprices.TaxedMoneyRange:
				unDiscountedLocalValue := unDiscountedLocal.(*goprices.TaxedMoneyRange)
				if less, err := t.Start.LessThan(unDiscountedLocalValue.Start); err == nil && less {
					discountLocalCurrency, _ = unDiscountedLocalValue.Start.Sub(t.Start)
				}
			}
		}
	}

	return &aStructType{
		PriceRangeLocal:       priceRangeLocal,
		DiscountLocalCurrency: discountLocalCurrency,
	}, nil
}

// GetVariantPrice
func (a *AppProduct) GetVariantPrice(
	variant *product_and_discount.ProductVariant,
	variantChannelListing *product_and_discount.ProductVariantChannelListing,
	product *product_and_discount.Product,
	collections []*product_and_discount.Collection,
	discounts []*product_and_discount.DiscountInfo,
	chanNel *channel.Channel,

) (*goprices.Money, *model.AppError) {

	variantChannelListing.PopulateNonDbFields() // must call this initially

	return a.DiscountApp().CalculateDiscountedPrice(
		product,
		variantChannelListing.Price,
		collections,
		discounts,
		chanNel,
	)
}

func (a *AppProduct) GetProductPriceRange(
	product *product_and_discount.Product,
	variants []*product_and_discount.ProductVariant,
	variantsChannelListing []*product_and_discount.ProductVariantChannelListing,
	collections []*product_and_discount.Collection,
	discounts []*product_and_discount.DiscountInfo,
	chanNel *channel.Channel,

) (*goprices.MoneyRange, *model.AppError) {

	// validate variantsChannelListing have same currency
	var currency string

	if len(variants) > 0 {
		variantChannelListingsMap := map[string]*product_and_discount.ProductVariantChannelListing{}
		for i, listing := range variantsChannelListing {
			if listing != nil {
				variantChannelListingsMap[listing.VariantID] = listing

				// compare or set currency for checking:
				if i == 0 {
					currency = listing.Currency
					continue
				}
				if !strings.EqualFold(currency, listing.Currency) {
					return nil, model.NewAppError("GetProductPriceRange", app.InvalidArgumentAppErrorID, map[string]interface{}{"Fields": "variantsChannelListing's currencies"}, "", http.StatusBadRequest)
				}
			}
		}

		prices := []*goprices.Money{}
		for _, variant := range variants {
			variantChannelListing := variantChannelListingsMap[variant.Id]
			if variantChannelListing != nil {
				price, appErr := a.GetVariantPrice(
					variant,
					variantChannelListing, // no need to populate non db fields, since GetVariantPrice() does that.
					product,
					collections,
					discounts,
					chanNel,
				)
				if appErr != nil {
					return nil, appErr
				}

				prices = append(prices, price)
			}
		}

		if len(prices) > 0 {
			min, max := util.MinMaxMoneyInMoneySlice(prices)
			return &goprices.MoneyRange{
				Start:    min,
				Stop:     max,
				Currency: min.Currency,
			}, nil
		}
	}

	return nil, nil
}

func (a *AppProduct) GetProductAvailability(
	product *product_and_discount.Product,
	productChannelListing *product_and_discount.ProductChannelListing,
	variants []*product_and_discount.ProductVariant,
	variantsChannelListing []*product_and_discount.ProductVariantChannelListing,
	collections []*product_and_discount.Collection,
	discounts []*product_and_discount.DiscountInfo,
	chanNel *channel.Channel,
	manager interface{},
	countryCode string, // can be empty
	localCurrency string, // can be empty

) (*product_and_discount.ProductAvailability, *model.AppError) {

	if countryCode == "" {
		countryCode = model.DEFAULT_COUNTRY
	}

	// discountedNetRange, appErr := a.GetProductPriceRange(
	// 	product,
	// 	variants,
	// 	variantsChannelListing,
	// 	collections,
	// 	discounts,
	// 	chanNel,
	// )
	// if appErr != nil {
	// 	return nil, appErr
	// }

	panic("not implemented")

}

func (a *AppProduct) GetVariantAvailability(
	variant *product_and_discount.ProductVariant,
	variantChannelListing *product_and_discount.ProductVariantChannelListing,
	product *product_and_discount.Product,
	productChannelListing *product_and_discount.ProductChannelListing,
	collections []*product_and_discount.Collection,
	discounts []*product_and_discount.DiscountInfo,
	chanNel *channel.Channel,
	plugins interface{},
	country string, // can be empty
	localCurrency string, // can be empty

) (*product_and_discount.VariantAvailability, *model.AppError) {
	panic("not implt")
}
