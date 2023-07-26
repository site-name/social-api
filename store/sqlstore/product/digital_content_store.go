package product

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlDigitalContentStore struct {
	store.Store
}

func NewSqlDigitalContentStore(s store.Store) store.DigitalContentStore {
	return &SqlDigitalContentStore{s}
}

func (ds *SqlDigitalContentStore) ScanFields(content *model.DigitalContent) []interface{} {
	return []interface{}{
		&content.Id,
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
	err := ds.GetMaster().Create(content).Error
	if err != nil {
		return nil, errors.Wrapf(err, "failed to save digital content with id=%s", content.Id)
	}

	return content, nil
}

// GetByOption finds and returns 1 digital content filtered using given option
func (ds *SqlDigitalContentStore) GetByOption(option *model.DigitalContentFilterOption) (*model.DigitalContent, error) {
	var res model.DigitalContent
	err := ds.GetReplica().First(&res, store.BuildSqlizer(option.Conditions)...).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.DigitalContentTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to find digital content with given option")
	}

	return &res, nil
}

func (ds *SqlDigitalContentStore) FilterByOption(option *model.DigitalContentFilterOption) ([]*model.DigitalContent, error) {
	var res []*model.DigitalContent
	err := ds.GetReplica().Find(&res, store.BuildSqlizer(option.Conditions)...).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find digital contents with given options")
	}

	return res, nil
}
