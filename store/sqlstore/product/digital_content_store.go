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
	args, err := store.BuildSqlizer(option.Conditions, "DigitalContent_GetByOption")
	if err != nil {
		return nil, err
	}

	var res model.DigitalContent
	err = ds.GetReplica().First(&res, args...).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound(model.DigitalContentTableName, "option")
		}
		return nil, errors.Wrap(err, "failed to find digital content with given option")
	}

	return &res, nil
}

func (ds *SqlDigitalContentStore) FilterByOption(option *model.DigitalContentFilterOption) (int64, []*model.DigitalContent, error) {
	query := ds.GetQueryBuilder().
		Select(model.DigitalContentTableName + ".*").
		From(model.DigitalContentTableName).
		Where(option.Conditions)

	var totalCount int64
	if option.CountTotal {
		countQuery, args, err := ds.GetQueryBuilder().Select("COUNT (*)").FromSelect(query, "subquery").ToSql()
		if err != nil {
			return 0, nil, errors.Wrap(err, "FilterByOptin_Count_ToSql")
		}
		err = ds.GetReplica().Raw(countQuery, args...).Scan(&totalCount).Error
		if err != nil {
			return 0, nil, errors.Wrap(err, "failed to count total number of digital content by given options")
		}
	}

	option.PaginationValues.AddPaginationToSelectBuilderIfNeeded(&query)

	queryStr, args, err := query.ToSql()
	if err != nil {
		return 0, nil, errors.Wrap(err, "FilterByOption_ToSql")
	}

	var res []*model.DigitalContent
	err = ds.GetReplica().Raw(queryStr, args...).Scan(&res).Error
	if err != nil {
		return 0, nil, errors.Wrap(err, "failed to find digital contents with given options")
	}

	return totalCount, res, nil
}

func (s *SqlDigitalContentStore) Delete(transaction *gorm.DB, options *model.DigitalContentFilterOption) error {
	if transaction == nil {
		transaction = s.GetMaster()
	}

	query, args, err := s.GetQueryBuilder().Delete(model.DigitalContentTableName).Where(options.Conditions).ToSql()
	if err != nil {
		return errors.Wrap(err, "Delete_ToSql")
	}

	err = transaction.Raw(query, args...).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete digital content")
	}

	return nil
}
