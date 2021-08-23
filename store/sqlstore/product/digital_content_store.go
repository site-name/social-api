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

// Save inserts given digital content into database then returns it
func (ds *SqlDigitalContentStore) Save(content *product_and_discount.DigitalContent) (*product_and_discount.DigitalContent, error) {
	content.PreSave()
	if err := content.IsValid(); err != nil {
		return nil, err
	}

	err := ds.GetMaster().Insert(content)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to save digital content with id=%s", content.Id)
	}

	return content, nil
}

// GetByOption finds and returns 1 digital content filtered using given option
func (ds *SqlDigitalContentStore) GetByOption(option *product_and_discount.DigitalContenetFilterOption) (*product_and_discount.DigitalContent, error) {
	query := ds.GetQueryBuilder().
		Select(ds.ModelFields()...).
		From(store.ProductDigitalContentTableName)

	// parse option
	if option.Id != nil {
		query = query.Where(option.Id.ToSquirrel("DigitalContents.Id"))
	}
	if option.ProductVariantID != nil {
		query = query.Where(option.ProductVariantID.ToSquirrel("DigitalContents.ProductVariantID"))
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetbyOption_ToSql")
	}

	var res product_and_discount.DigitalContent
	err = ds.GetReplica().SelectOne(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductDigitalContentTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to find digital content with given option")
	}

	return &res, nil
}
