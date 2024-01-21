package file

import (
	"github.com/pkg/errors"
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gorm.io/gorm"
)

type SqlFileInfoStore struct {
	store.Store
	metrics     einterfaces.MetricsInterface
	queryFields util.AnyArray[string]
}

func (fs *SqlFileInfoStore) ClearCaches() {}

func NewSqlFileInfoStore(sqlStore store.Store, metrics einterfaces.MetricsInterface) store.FileInfoStore {
	s := &SqlFileInfoStore{
		Store:   sqlStore,
		metrics: metrics,
	}

	s.queryFields = util.AnyArray[string]{
		"FileInfos.Id",
		"FileInfos.CreatorId",
		"FileInfos.ParentID",
		"FileInfos.CreateAt",
		"FileInfos.UpdateAt",
		"FileInfos.DeleteAt",
		"FileInfos.Path",
		"FileInfos.ThumbnailPath",
		"FileInfos.PreviewPath",
		"FileInfos.Name",
		"FileInfos.Extension",
		"FileInfos.Size",
		"FileInfos.MimeType",
		"FileInfos.Width",
		"FileInfos.Height",
		"FileInfos.HasPreviewImage",
		"FileInfos.MiniPreview",
		"Coalesce(FileInfos.Content, '') AS Content",
		"Coalesce(FileInfos.RemoteId, '') AS RemoteId",
	}

	return s
}

func (fs *SqlFileInfoStore) Upsert(info model.FileInfo) (*model.FileInfo, error) {
	err := fs.GetMaster().Save(info).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert file info")
	}
	return info, nil
}

func (fs *SqlFileInfoStore) Get(id string, fromMaster bool) (*model.FileInfo, error) {
	info := model.FileInfo{}
	var db *gorm.DB
	switch {
	case fromMaster:
		db = fs.GetMaster()
	default:
		db = fs.GetReplica()
	}

	err := db.First(&info, "Id = ? AND DeleteAt = 0").Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound("FileInfos", id)
		}
		return nil, errors.Wrapf(err, "failed to get FileInfo with id=%s", id)
	}

	return &info, nil
}

// GetWithOptions finds and returns fileinfos with given options.
// Leave page, perPage nil to get all result.
func (fs *SqlFileInfoStore) GetWithOptions(conds ...qm.QueryMod) (model.FileInfoSlice, error) {
	query := fs.GetQueryBuilder().
		Select(fs.queryFields...).
		From(model.FileInfoTableName).
		Where(opt.Conditions)

	if opt.Limit > 0 {
		query = query.Limit(opt.Limit)
	}
	if opt.Offset > 0 {
		query = query.Offset(opt.Offset)
	}
	if len(opt.OrderBy) > 0 {
		query = query.OrderBy(opt.OrderBy)
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "file_info_tosql")
	}
	var infos []*model.FileInfo
	if err := fs.GetReplica().Raw(queryString, args...).Scan(&infos).Error; err != nil {
		return nil, errors.Wrap(err, "failed to find FileInfos")
	}
	return infos, nil
}

func (fs *SqlFileInfoStore) InvalidateFileInfosForPostCache(postId string, deleted bool) {
}

func (fs *SqlFileInfoStore) PermanentDelete(fileId string) error {
	if err := fs.GetMaster().Raw(`DELETE FROM FileInfos WHERE Id = ?`, fileId).Error; err != nil {
		return errors.Wrapf(err, "failed to delete FileInfos with id=%s", fileId)
	}
	return nil
}

func (fs *SqlFileInfoStore) PermanentDeleteBatch(endTime int64, limit int64) (int64, error) {
	result := fs.GetMaster().Raw(
		`DELETE FROM 
			FileInfos 
		WHERE Id = any (
			array (
				SELECT 
					Id 
				FROM 
					FileInfos 
				WHERE 
					CreateAt < ? 
				LIMIT ?
			)
		)`,
		endTime,
		limit,
	)
	if result.Error != nil {
		return 0, errors.Wrap(result.Error, "failed to delete FileInfos in batch")
	}

	return result.RowsAffected, nil
}

func (fs SqlFileInfoStore) PermanentDeleteByUser(userId string) (int64, error) {
	sqlResult := fs.GetMaster().Raw("DELETE FROM FileInfos WHERE CreatorId = ?", userId)
	if sqlResult.Error != nil {
		return 0, errors.Wrapf(sqlResult.Error, "failed to delete FileInfos with creatorId=%s", userId)
	}

	return sqlResult.RowsAffected, nil
}

func (fs *SqlFileInfoStore) CountAll() (int64, error) {
	var count int64
	err := fs.GetReplica().Raw("SELECT COUNT(*) FROM " + model.FileInfoTableName + " WHERE DeleteAt = 0").Scan(&count).Error
	if err != nil {
		return int64(0), errors.Wrap(err, "failed to count Files")
	}
	return count, nil
}
