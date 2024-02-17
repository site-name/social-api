package file

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
)

type SqlFileInfoStore struct {
	store.Store
	metrics einterfaces.MetricsInterface
}

func (fs *SqlFileInfoStore) ClearCaches() {}

func NewSqlFileInfoStore(sqlStore store.Store, metrics einterfaces.MetricsInterface) store.FileInfoStore {
	return &SqlFileInfoStore{
		Store:   sqlStore,
		metrics: metrics,
	}
}

func (fs *SqlFileInfoStore) Upsert(info model.FileInfo) (*model.FileInfo, error) {
	isSaving := info.ID == ""
	if isSaving {
		model_helper.FileInfoPreSave(&info)
	} else {
		model_helper.FileInfoPreUpdate(&info)
	}

	if err := model_helper.FileInfoIsValid(info); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = info.Insert(fs.GetMaster(), boil.Infer())
	} else {
		_, err = info.Update(fs.GetMaster(), boil.Blacklist(model.FileInfoColumns.CreatedAt))
	}

	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (fs *SqlFileInfoStore) Get(id string, fromMaster bool) (*model.FileInfo, error) {
	db := fs.GetReplica()
	if fromMaster {
		db = fs.GetMaster()
	}

	fileInfo, err := model.FileInfos(
		model.FileInfoWhere.ID.EQ(id),
		model.FileInfoWhere.DeleteAt.EQ(model_types.NewNullInt64(0)),
	).One(db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.FileInfos, id)
		}
		return nil, err
	}

	return fileInfo, nil
}

func (fs *SqlFileInfoStore) GetWithOptions(options model_helper.FileInfoFilterOption) (model.FileInfoSlice, error) {
	return model.FileInfos(options.Conditions...).All(fs.GetReplica())
}

func (fs *SqlFileInfoStore) InvalidateFileInfosForPostCache(postId string, deleted bool) {
}

func (fs *SqlFileInfoStore) PermanentDelete(fileId string) error {
	_, err := model.FileInfos(model.FileInfoWhere.ID.EQ(fileId)).DeleteAll(fs.GetMaster())
	return err
}

func (fs *SqlFileInfoStore) PermanentDeleteBatch(endTime int64, limit int64) (int64, error) {
	result, err := queries.Raw(
		fmt.Sprintf(
			`DELETE FROM %[1]s WHERE %[2]s = any (array (SELECT %[2]s FROM %[1]s WHERE %[3]s < $1 LIMIT $2))`,
			model.TableNames.FileInfos,      // 1
			model.FileInfoColumns.ID,        // 2
			model.FileInfoColumns.CreatedAt, // 3
		),
		endTime,
		limit,
	).Exec(fs.GetMaster())
	if err != nil {
		return 0, errors.Wrap(err, "failed to delete FileInfos in batch")
	}

	return result.RowsAffected()
}

func (fs SqlFileInfoStore) PermanentDeleteByUser(userId string) (int64, error) {
	return model.FileInfos(model.FileInfoWhere.CreatorID.EQ(userId)).DeleteAll(fs.GetMaster())
}

func (fs *SqlFileInfoStore) CountAll() (int64, error) {
	return model.FileInfos(model.FileInfoWhere.DeleteAt.EQ(model_types.NewNullInt64(0))).Count(fs.GetReplica())
}
