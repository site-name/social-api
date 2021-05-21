package sqlstore

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlDigitalContentStore struct {
	*SqlStore
}

func newSqlDigitalContentStore(s *SqlStore) store.DigitalContentStore {
	dcs := &SqlDigitalContentStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.DigitalContent{}, "DigitalContents").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ProductVariantID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("ContentType").SetMaxSize(product_and_discount.DIGITAL_CONTENT_CONTENT_TYPE_MAX_LENGTH).
			SetDefaultConstraint(model.NewString(product_and_discount.FILE))
	}
	return dcs
}

func (ps *SqlDigitalContentStore) createIndexesIfNotExists() {

}
