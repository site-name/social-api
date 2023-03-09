package product

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
)

type SqlDigitalContentStore struct {
	store.Store
}

func NewSqlDigitalContentStore(s store.Store) store.DigitalContentStore {
	return &SqlDigitalContentStore{s}
}

func (ds *SqlDigitalContentStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"ShopID",
		"UseDefaultSettings",
		"AutomaticFulfillment",
		"ContentType",
		"ProductVariantID",
		"ContentFile",
		"MaxDownloads",
		"UrlValidDays",
		"Metadata",
		"PrivateMetadata",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func (ds *SqlDigitalContentStore) ScanFields(content *model.DigitalContent) []interface{} {
	return []interface{}{
		&content.Id,
		&content.ShopID,
		&content.UseDefaultSettings,
		&content.AutomaticFulfillment,
		&content.ContentType,
		&content.ProductVariantID,
		&content.ContentFile,
		&content.MaxDownloads,
		&content.UrlValidDays,
		&content.Metadata,
		&content.PrivateMetadata,
	}
}

// Save inserts given digital content into database then returns it
func (ds *SqlDigitalContentStore) Save(content *model.DigitalContent) (*model.DigitalContent, error) {
	content.PreSave()
	if err := content.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.DigitalContentTableName + "(" + ds.ModelFields("").Join(",") + ") VALUES (" + ds.ModelFields(":").Join(",") + ")"
	_, err := ds.GetMasterX().NamedExec(query, content)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to save digital content with id=%s", content.Id)
	}

	return content, nil
}

func (ds *SqlDigitalContentStore) commonQueryBuilder(option *model.DigitalContentFilterOption) (string, []interface{}, error) {
	query := ds.GetQueryBuilder().
		Select(ds.ModelFields(store.DigitalContentTableName + ".")...).
		From(store.DigitalContentTableName)

	// parse option
	if option.Id != nil {
		query = query.Where(option.Id)
	}
	if option.ProductVariantID != nil {
		query = query.Where(option.ProductVariantID)
	}

	return query.ToSql()
}

// GetByOption finds and returns 1 digital content filtered using given option
func (ds *SqlDigitalContentStore) GetByOption(option *model.DigitalContentFilterOption) (*model.DigitalContent, error) {
	queryString, args, err := ds.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "GetbyOption_ToSql")
	}

	var res model.DigitalContent
	err = ds.GetReplicaX().Get(&res, queryString, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.DigitalContentTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to find digital content with given option")
	}

	return &res, nil
}

func (ds *SqlDigitalContentStore) FilterByOption(option *model.DigitalContentFilterOption) ([]*model.DigitalContent, error) {
	queryString, args, err := ds.commonQueryBuilder(option)
	if err != nil {
		return nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*model.DigitalContent
	err = ds.GetReplicaX().Select(&res, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find digital contents with given options")
	}

	return res, nil
}
