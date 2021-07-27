package product

// import (
// 	"fmt"
// 	"testing"
// 	"time"

// 	"github.com/google/uuid"
// 	"github.com/sitename/sitename/app"
// 	"github.com/sitename/sitename/model"
// 	"github.com/sitename/sitename/model/product_and_discount"
// 	"github.com/sitename/sitename/modules/util"
// 	"github.com/sitename/sitename/store/sqlstore"
// 	"github.com/stretchr/testify/require"
// )

// func Test_FilterByOption(t *testing.T) {
// 	srv, err := app.NewServer()
// 	require.NoError(t, err)

// 	app := app.New(app.ServerConnector(srv))
// 	store := &SqlProductChannelListingStore{
// 		Store: sqlstore.New(app.Config().SqlSettings, app.Srv().Metrics),
// 	}

// 	res, err := store.FilterByOption(&product_and_discount.ProductChannelListingFilterOption{
// 		ProductID: &product_and_discount.StringFilter{
// 			Eq: uuid.NewString(),
// 		},
// 		ChannelID: &product_and_discount.StringFilter{
// 			In: []string{"one", "two", "three"},
// 		},
// 		ChannelSlug:       model.NewString("this-is-a-slug"),
// 		VisibleInListings: model.NewBool(true),
// 		AvailableForPurchase: &model.TimeFilter{
// 			Gt:              util.NewTime(time.Now()),
// 			LtE:             util.NewTime(time.Now()),
// 			CompareFullTime: true,
// 		},
// 		Currency: &product_and_discount.StringFilter{
// 			Eq: "VND",
// 		},
// 		ProductVariantsId: &product_and_discount.StringFilter{
// 			Eq: uuid.NewString(),
// 		},
// 		PublicationDate: &model.TimeFilter{
// 			Eq:              util.NewTime(time.Now()),
// 			CompareFullTime: false,
// 		},
// 		IsPublished: model.NewBool(true),
// 	})

// 	require.NoError(t, err)

// 	fmt.Println(res)
// }
