package product

import (
	"net/http"
	"strings"

	goprices "github.com/site-name/go-prices"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
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
	if discounted.LessThan(unDiscounted) {
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
func (a *ServiceProduct) getProductPriceRange(discounted interface{}, unDiscounted interface{}, localCurrency string) (*aStructType, *model.AppError) {

	// validate `discounted` and `unDiscounted` and `localCurrency` are valid and have same currencies
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
		priceRangeLocal, appErr := a.srv.ToLocalCurrency(discounted, localCurrency)
		if appErr != nil {
			return nil, appErr
		}
		unDiscountedLocal, appErr := a.srv.ToLocalCurrency(unDiscounted, localCurrency)
		if appErr != nil {
			return nil, appErr
		}

		if unDiscountedLocal != nil {
			switch t := priceRangeLocal.(type) {
			case *goprices.MoneyRange:
				unDiscountedLocalValue := unDiscountedLocal.(*goprices.MoneyRange)
				if t.Start.LessThan(unDiscountedLocalValue.Start) {
					discountLocalCurrency, _ = unDiscountedLocalValue.Start.Sub(t.Start)
				}

			case *goprices.TaxedMoneyRange:
				unDiscountedLocalValue := unDiscountedLocal.(*goprices.TaxedMoneyRange)
				if t.Start.LessThan(unDiscountedLocalValue.Start) {
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
func (a *ServiceProduct) GetVariantPrice(
	variant model.ProductVariant,
	variantChannelListing model.ProductVariantChannelListing,
	product model.Product,
	collections []*model.Collection,
	discounts []*model.DiscountInfo,
	chanNel model.Channel,

) (*goprices.Money, *model.AppError) {

	variantChannelListing.PopulateNonDbFields() // must call this initially

	return a.srv.DiscountService().CalculateDiscountedPrice(
		product,
		variantChannelListing.Price,
		collections,
		discounts,
		chanNel,
		variant.Id,
	)
}

func (a *ServiceProduct) GetProductPriceRange(
	product model.Product,
	variants model.ProductVariants,
	variantsChannelListing []*model.ProductVariantChannelListing,
	collections []*model.Collection,
	discounts []*model.DiscountInfo,
	chanNel model.Channel,

) (*goprices.MoneyRange, *model.AppError) {

	// validate variantsChannelListing have same currency
	var currency string

	if len(variants) > 0 {
		variantChannelListingsMap := map[string]*model.ProductVariantChannelListing{}
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
		for _, variant := range variants.FilterNils() {
			variantChannelListing := variantChannelListingsMap[variant.Id]
			if variantChannelListing != nil {
				price, appErr := a.GetVariantPrice(
					*variant,
					*variantChannelListing, // no need to populate non db fields, since GetVariantPrice() does that.
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

func (a *ServiceProduct) GetProductAvailability(
	product model.Product,
	productChannelListing *model.ProductChannelListing,
	variants []*model.ProductVariant,
	variantsChannelListing []*model.ProductVariantChannelListing,
	collections []*model.Collection,
	discounts []*model.DiscountInfo,
	chanNel model.Channel,
	manager interfaces.PluginManagerInterface,
	countryCode string, // can be empty
	localCurrency string, // can be empty

) (*model.ProductAvailability, *model.AppError) {

	if countryCode == "" {
		countryCode = model.DEFAULT_COUNTRY
	}

	var discounted *goprices.TaxedMoneyRange

	discountedNetRange, appErr := a.GetProductPriceRange(product, variants, variantsChannelListing, collections, discounts, chanNel)
	if appErr != nil {
		return nil, appErr
	}

	if discountedNetRange != nil {
		start, appErr := manager.ApplyTaxesToProduct(product, *discountedNetRange.Start, countryCode, chanNel.Slug)
		if appErr != nil {
			return nil, appErr
		}
		stop, appErr := manager.ApplyTaxesToProduct(product, *discountedNetRange.Stop, countryCode, chanNel.Slug)
		if appErr != nil {
			return nil, appErr
		}

		discounted = &goprices.TaxedMoneyRange{
			Start:    start,
			Stop:     stop,
			Currency: start.Currency,
		}
	}

	var undiscounted *goprices.TaxedMoneyRange
	undiscountedNetRange, appErr := a.GetProductPriceRange(product, variants, variantsChannelListing, collections, []*model.DiscountInfo{}, chanNel)
	if appErr != nil {
		return nil, appErr
	}

	if undiscountedNetRange != nil {
		start, appErr := manager.ApplyTaxesToProduct(product, *undiscountedNetRange.Start, countryCode, chanNel.Slug)
		if appErr != nil {
			return nil, appErr
		}

		stop, appErr := manager.ApplyTaxesToProduct(product, *undiscountedNetRange.Stop, countryCode, chanNel.Slug)
		if appErr != nil {
			return nil, appErr
		}

		undiscounted = &goprices.TaxedMoneyRange{
			Start:    start,
			Stop:     stop,
			Currency: start.Currency,
		}
	}

	var (
		discount              *goprices.TaxedMoney
		priceRangeLocal       *goprices.TaxedMoneyRange
		discountLocalCurrency *goprices.TaxedMoneyRange
	)
	if discountedNetRange != nil && undiscountedNetRange != nil {
		discount, _ = getTotalDiscountFromRange(undiscounted, discounted)

		aType, appErr := a.getProductPriceRange(discounted, undiscounted, localCurrency)
		if appErr != nil {
			return nil, appErr
		}

		priceRangeLocal = aType.PriceRangeLocal.(*goprices.TaxedMoneyRange)
		discountLocalCurrency = aType.DiscountLocalCurrency.(*goprices.TaxedMoneyRange)
	}

	isOnSale := productChannelListing != nil && productChannelListing.IsVisible() && discount != nil

	return &model.ProductAvailability{
		OnSale:                  isOnSale,
		PriceRange:              discounted,
		PriceRangeUnDiscounted:  undiscounted,
		Discount:                discount,
		PriceRangeLocalCurrency: priceRangeLocal,
		DiscountLocalCurrency:   discountLocalCurrency,
	}, nil
}

func (a *ServiceProduct) GetVariantAvailability(
	variant model.ProductVariant,
	variantChannelListing model.ProductVariantChannelListing,
	product model.Product,
	productChannelListing *model.ProductChannelListing,
	collections []*model.Collection,
	discounts []*model.DiscountInfo,
	chanNel model.Channel,
	plugins interfaces.PluginManagerInterface,
	country string, // can be empty
	localCurrency string, // can be empty

) (*model.VariantAvailability, *model.AppError) {

	variarntPrice, appErr := a.GetVariantPrice(variant, variantChannelListing, product, collections, discounts, chanNel)
	if appErr != nil {
		return nil, appErr
	}

	discounted, appErr := plugins.ApplyTaxesToProduct(product, *variarntPrice, country, chanNel.Id)
	if appErr != nil {
		return nil, appErr
	}

	variarntPrice, appErr = a.GetVariantPrice(variant, variantChannelListing, product, collections, []*model.DiscountInfo{}, chanNel)
	if appErr != nil {
		return nil, appErr
	}

	undiscounted, appErr := plugins.ApplyTaxesToProduct(product, *variarntPrice, country, chanNel.Id)
	if appErr != nil {
		return nil, appErr
	}

	discount, err := getTotalDiscount(undiscounted, discounted)
	if err != nil {
		return nil, model.NewAppError("GetVariantAvailability", app.ErrorCalculatingMoneyErrorID, nil, err.Error(), http.StatusInternalServerError)
	}

	var (
		priceLocalCurrency    *goprices.TaxedMoney
		discountLocalCurrency *goprices.TaxedMoney
	)
	if localCurrency != "" {
		iface1, appErr := a.srv.ToLocalCurrency(discounted, localCurrency)
		if appErr != nil {
			return nil, appErr
		}
		priceLocalCurrency = iface1.(*goprices.TaxedMoney)

		iface2, appErr := a.srv.ToLocalCurrency(discount, localCurrency)
		if appErr != nil {
			return nil, appErr
		}
		discountLocalCurrency = iface2.(*goprices.TaxedMoney)
	}

	isOnSale := (productChannelListing != nil && productChannelListing.IsVisible()) && discount != nil

	return &model.VariantAvailability{
		OnSale:                isOnSale,
		Price:                 *discounted,
		PriceUnDiscounted:     *undiscounted,
		Discount:              discount,
		PriceLocalCurrency:    priceLocalCurrency,
		DiscountLocalCurrency: discountLocalCurrency,
	}, nil
}
