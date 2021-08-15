package product

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlDigitalContentStore struct {
	store.Store
}

func NewSqlDigitalContentStore(s store.Store) store.DigitalContentStore {
	dcs := &SqlDigitalContentStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.DigitalContent{}, store.ProductDigitalContentTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductVariantID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ContentType").SetMaxSize(product_and_discount.DIGITAL_CONTENT_CONTENT_TYPE_MAX_LENGTH)
		table.ColMap("ContentFile").SetMaxSize(model.URL_LINK_MAX_LENGTH)
	}
	return dcs
}

func (ds *SqlDigitalContentStore) CreateIndexesIfNotExists() {
	ds.CreateForeignKeyIfNotExists(store.ProductDigitalContentTableName, "ProductVariantID", store.ProductVariantTableName, "Id", true)
}

func (ds *SqlDigitalContentStore) ModelFields() []string {
	return []string{
		"DigitalContents.Id",
		"DigitalContents.UseDefaultSettings",
		"DigitalContents.AutomaticFulfillment",
		"DigitalContents.ContentType",
		"DigitalContents.ProductVariantID",
		"DigitalContents.ContentFile",
		"DigitalContents.MaxDownloads",
		"DigitalContents.UrlValidDays",
		"DigitalContents.Metadata",
		"DigitalContents.PrivateMetadata",
	}
}

// GetByProductVariantID finds and returns 1 digital content that is related to given product variant
func (ds *SqlDigitalContentStore) GetByProductVariantID(variantID string) (*product_and_discount.DigitalContent, error) {
	var res product_and_discount.DigitalContent
	err := ds.GetReplica().SelectOne(&res, "SELECT * FROM "+store.ProductDigitalContentTableName+" WHERE ProductVariantID = :ID", map[string]interface{}{"ID": variantID})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductDigitalContentTableName, "productVariantID="+variantID)
		}
		return nil, errors.Wrapf(err, "failed to find digital content with product variant id=%s", variantID)
	}

	return &res, nil
}
