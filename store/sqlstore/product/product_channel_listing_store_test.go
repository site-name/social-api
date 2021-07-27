package product

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/modules/util"
)

func Test_FilterByOption(t *testing.T) {
	var s = new(SqlProductChannelListingStore)
	sql, args, err := s.FilterByOption(&product_and_discount.ProductChannelListingFilterOption{
		ProductID: &product_and_discount.StringFilter{
			Eq: uuid.NewString(),
		},
		ChannelID: &product_and_discount.StringFilter{
			In: []string{"one", "two", "three"},
		},
		ChannelSlug:       model.NewString("this-is-a-slug"),
		VisibleInListings: model.NewBool(true),
		AvailableForPurchase: &model.TimeFilter{
			Gt:   util.NewTime(time.Now()),
			LtE:  util.NewTime(time.Now()),
			Full: false,
		},
		Currency: &product_and_discount.StringFilter{
			Eq: "VND",
		},
		ProductVariantsId: &product_and_discount.StringFilter{
			Eq: uuid.NewString(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(sql)
	fmt.Println(args)
}
