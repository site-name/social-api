package storetest

import (
	"testing"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
	"github.com/stretchr/testify/require"
)

func TestProductStore(t *testing.T, ss store.Store, s SqlStore) {
	_, err := ss.Product().FilterByOption(&product_and_discount.ProductFilterOption{})
	require.NoError(t, err, "failed cleaning up test products")

	t.Run("Save", func(t *testing.T) { testSave(t, ss) })
	t.Run("FilterByOption", func(t *testing.T) { testFilterByOption(t, ss) })
}

func testSave(t *testing.T, ss store.Store) {
	id := model.NewId()

	product := product_and_discount.Product{
		Id: id,
	}

	_, err := ss.Product().Save(&product)
	require.NoError(t, err, "couldn't save product")
}

func testFilterByOption(t *testing.T, ss store.Store) {

}
