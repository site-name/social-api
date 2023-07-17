package file

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
	"github.com/sitename/sitename/store"
	"gorm.io/gorm"
)

type SqlFileInfoStore struct {
	store.Store
	metrics     einterfaces.MetricsInterface
	queryFields util.AnyArray[string]
}

func (fs *SqlFileInfoStore) ModelFields(prefix string) util.AnyArray[string] {
	res := util.AnyArray[string]{
		"Id",
		"CreatorId",
		"ParentID",
		"CreateAt",
		"UpdateAt",
		"DeleteAt",
		"Path",
		"ThumbnailPath",
		"PreviewPath",
		"Name",
		"Extension",
		"Size",
		"MimeType",
		"Width",
		"Height",
		"HasPreviewImage",
		"MiniPreview",
		"Content",
		"RemoteId",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
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

func (fs *SqlFileInfoStore) Upsert(info *model.FileInfo) (*model.FileInfo, error) {
	var isSaving bool

	if !model.IsValidId(info.Id) {
		isSaving = true
		info.PreSave()
	}

	if err := info.IsValid(); err != nil {
		return nil, err
	}

	var (
		err        error
		numUpdated int64
	)

	if isSaving {
		query := "INSERT INTO " + model.FileInfoTableName + "(" + fs.ModelFields("").Join(",") + ") VALUES (" + fs.ModelFields(":").Join(",") + ")"
		_, err = fs.GetMasterX().NamedExec(query, info)

	} else {
		query := "UPDATE " + model.FileInfoTableName + " SET " + fs.
			ModelFields("").
			Map(func(_ int, s string) string {
				return s + "=:" + s
			}).
			Join(",") + " WHERE Id=:Id"

		var result sql.Result
		result, err = fs.GetMasterX().NamedExec(query, info)
		if err == nil && result != nil {
			numUpdated, _ = result.RowsAffected()
		}
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert given file info")
	}

	if numUpdated > 1 {
		return nil, errors.Errorf("updated %d fileinfos instead of 1", numUpdated)
	}

	return info, nil
}

func (fs *SqlFileInfoStore) GetByIds(ids []string) ([]*model.FileInfo, error) {
	var infos []*model.FileInfo
	query, args, err := fs.GetQueryBuilder().
		Select(fs.queryFields...).From(model.FileInfoTableName).
		Where(squirrel.Eq{"Id": ids}).
		Where("DeleteAt = 0").
		OrderBy("CreateAt DESC").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "GetByIds_ToSql")
	}
	if err := fs.GetReplicaX().Select(&infos, query, args...); err != nil {
		return nil, errors.Wrap(err, "failed to find FileInfos")
	}
	return infos, nil
}

func (fs *SqlFileInfoStore) get(id string, fromMaster bool) (*model.FileInfo, error) {
	info := &model.FileInfo{}

	query := fs.GetQueryBuilder().
		Select(fs.queryFields...).
		From(model.FileInfoTableName).
		Where(squirrel.Eq{"Id": id}).
		Where(squirrel.Eq{"DeleteAt": 0})
	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "file_info_tosql")
	}

	db := fs.GetReplicaX()
	if fromMaster {
		db = fs.GetMasterX()
	}

	if err := db.Get(info, queryString, args...); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound("FileInfos", id)
		}
		return nil, errors.Wrapf(err, "failed to get FileInfo with id=%s", id)
	}
	return info, nil
}

func (fs *SqlFileInfoStore) Get(id string) (*model.FileInfo, error) {
	return fs.get(id, false)
}

func (fs *SqlFileInfoStore) GetFromMaster(id string) (*model.FileInfo, error) {
	return fs.get(id, true)
}

// GetWithOptions finds and returns fileinfos with given options.
// Leave page, perPage nil to get all result.
func (fs *SqlFileInfoStore) GetWithOptions(page, perPage *int, opt *model.GetFileInfosOptions) ([]*model.FileInfo, error) {
	if perPage != nil && *perPage < 0 {
		return nil, store.NewErrLimitExceeded("perPage", *perPage, "value used in pagination while getting FileInfos")
	} else if page != nil && *page < 0 {
		return nil, store.NewErrLimitExceeded("page", *page, "value used in pagination while getting FileInfos")
	}
	if *perPage == 0 {
		return nil, nil
	}

	if opt == nil {
		opt = &model.GetFileInfosOptions{}
	}

	query := fs.GetQueryBuilder().
		Select(fs.queryFields...).
		From(model.FileInfoTableName)

	if len(opt.UserIds) > 0 {
		query = query.Where(squirrel.Eq{"FileInfos.CreatorId": opt.UserIds})
	}

	if opt.Since > 0 {
		query = query.Where(squirrel.GtOrEq{"FileInfos.CreateAt": opt.Since})
	}

	if !opt.IncludeDeleted {
		query = query.Where("FileInfos.DeleteAt = 0")
	}
	if len(opt.ParentID) > 0 {
		query = query.Where(squirrel.Eq{"FileInfos.ParentID": opt.ParentID})
	}

	if opt.SortBy == "" {
		opt.SortBy = model.FILEINFO_SORT_BY_CREATED
	}
	sortDirection := "ASC"
	if opt.SortDescending {
		sortDirection = "DESC"
	}

	switch opt.SortBy {
	case model.FILEINFO_SORT_BY_CREATED:
		query = query.OrderBy("FileInfos.CreateAt " + sortDirection)
	case model.FILEINFO_SORT_BY_SIZE:
		query = query.OrderBy("FileInfos.Size " + sortDirection)
	default:
		return nil, store.NewErrInvalidInput("FileInfos", "<sortOption>", opt.SortBy)
	}

	query = query.OrderBy("FileInfos.Id ASC") // secondary sort for sort stability

	if perPage != nil && page != nil {
		query = query.Limit(uint64(*perPage)).Offset(uint64(*perPage * *page))
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "file_info_tosql")
	}
	var infos []*model.FileInfo
	if err := fs.GetReplicaX().Select(&infos, queryString, args...); err != nil {
		return nil, errors.Wrap(err, "failed to find FileInfos")
	}
	return infos, nil
}

func (fs *SqlFileInfoStore) GetByPath(path string) (*model.FileInfo, error) {
	var info model.FileInfo

	query := fs.GetQueryBuilder().
		Select(fs.queryFields...).
		From(model.FileInfoTableName).
		Where(squirrel.Eq{"Path": path}).
		Where(squirrel.Eq{"DeleteAt": 0}).
		Limit(1)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "file_info_tosql")
	}

	if err := fs.GetReplicaX().Get(info, queryString, args...); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.NewErrNotFound("FileInfos", fmt.Sprintf("path=%s", path))
		}

		return nil, errors.Wrapf(err, "failed to get FileInfos with path=%s", path)
	}
	return &info, nil
}

func (fs *SqlFileInfoStore) InvalidateFileInfosForPostCache(postId string, deleted bool) {
}

func (fs *SqlFileInfoStore) GetForUser(userId string) ([]*model.FileInfo, error) {
	var infos []*model.FileInfo

	query := fs.GetQueryBuilder().
		Select(fs.queryFields...).
		From(model.FileInfoTableName).
		Where(squirrel.Eq{"CreatorId": userId}).
		Where(squirrel.Eq{"DeleteAt": 0}).
		OrderBy("CreateAt")
	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "file_info_tosql")
	}

	if err := fs.GetReplicaX().Select(&infos, queryString, args...); err != nil {
		return nil, errors.Wrapf(err, "failed to find FileInfos with creatorId=%s", userId)
	}

	return infos, nil
}

func (fs *SqlFileInfoStore) SetContent(fileId, content string) error {
	_, err := fs.GetMasterX().Exec("UPDATE "+model.FileInfoTableName+" SET Content=? WHERE Id=?", content, fileId)
	if err != nil {
		return errors.Wrapf(err, "failed to update FileInfos content with id=%s", fileId)
	}

	return nil
}

func (fs *SqlFileInfoStore) PermanentDelete(fileId string) error {
	if _, err := fs.GetMasterX().Exec(
		`DELETE FROM
				FileInfos
			WHERE
				Id = ?`,
		fileId); err != nil {
		return errors.Wrapf(err, "failed to delete FileInfos with id=%s", fileId)
	}
	return nil
}

func (fs *SqlFileInfoStore) PermanentDeleteBatch(endTime int64, limit int64) (int64, error) {
	sqlResult, err := fs.GetMasterX().Exec(
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
	if err != nil {
		return 0, errors.Wrap(err, "failed to delete FileInfos in batch")
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return 0, errors.Wrapf(err, "unable to retrieve rows affected")
	}

	return rowsAffected, nil
}

func (fs SqlFileInfoStore) PermanentDeleteByUser(userId string) (int64, error) {
	sqlResult, err := fs.GetMasterX().Exec("DELETE FROM FileInfos WHERE CreatorId = ?", userId)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to delete FileInfos with creatorId=%s", userId)
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return 0, errors.Wrapf(err, "unable to retrieve rows affected")
	}

	return rowsAffected, nil
}

func (fs *SqlFileInfoStore) CountAll() (int64, error) {
	query := fs.GetQueryBuilder().
		Select("COUNT(*)").
		From(model.FileInfoTableName).
		Where("DeleteAt = 0")

	queryString, args, err := query.ToSql()
	if err != nil {
		return int64(0), errors.Wrap(err, "count_tosql")
	}

	var count int64
	err = fs.GetReplicaX().Get(&count, queryString, args...)
	if err != nil {
		return int64(0), errors.Wrap(err, "failed to count Files")
	}
	return count, nil
}
