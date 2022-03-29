package file

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model/file"
	"github.com/sitename/sitename/store"
)

type SqlFileInfoStore struct {
	store.Store
	metrics     einterfaces.MetricsInterface
	queryFields []string
}

func (fs *SqlFileInfoStore) ClearCaches() {}

func NewSqlFileInfoStore(sqlStore store.Store, metrics einterfaces.MetricsInterface) store.FileInfoStore {
	s := &SqlFileInfoStore{
		Store:   sqlStore,
		metrics: metrics,
	}

	s.queryFields = []string{
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

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(file.FileInfo{}, store.FileInfoTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CreatorId").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ParentID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Path").SetMaxSize(512)
		table.ColMap("ThumbnailPath").SetMaxSize(512)
		table.ColMap("PreviewPath").SetMaxSize(512)
		table.ColMap("Name").SetMaxSize(256)
		table.ColMap("Content").SetMaxSize(0)
		table.ColMap("Extension").SetMaxSize(64)
		table.ColMap("MimeType").SetMaxSize(256)
		table.ColMap("RemoteId").SetMaxSize(26)
	}

	return s
}

func (fs *SqlFileInfoStore) CreateIndexesIfNotExists() {
	fs.CreateIndexIfNotExists("idx_fileinfo_update_at", store.FileInfoTableName, "UpdateAt")
	fs.CreateIndexIfNotExists("idx_fileinfo_create_at", store.FileInfoTableName, "CreateAt")
	fs.CreateIndexIfNotExists("idx_fileinfo_delete_at", store.FileInfoTableName, "DeleteAt")
	fs.CreateIndexIfNotExists("idx_fileinfo_extension_at", store.FileInfoTableName, "Extension")
	fs.CreateIndexIfNotExists("idx_fileinfo_parent_id", store.FileInfoTableName, "ParentID")

	fs.CreateFullTextIndexIfNotExists("idx_fileinfo_name_txt", store.FileInfoTableName, "Name")
	fs.CreateFullTextFuncIndexIfNotExists("idx_fileinfo_name_splitted", store.FileInfoTableName, "Translate(Name, '.,-', '   ')")
	fs.CreateFullTextIndexIfNotExists("idx_fileinfo_content_txt", store.FileInfoTableName, "Content")
}

func (fs *SqlFileInfoStore) Save(info *file.FileInfo) (*file.FileInfo, error) {
	info.PreSave()
	if err := info.IsValid(); err != nil {
		return nil, err
	}

	if err := fs.GetMaster().Insert(info); err != nil {
		return nil, errors.Wrap(err, "failed to save FileInfos")
	}
	return info, nil
}

func (fs *SqlFileInfoStore) GetByIds(ids []string) ([]*file.FileInfo, error) {
	query := fs.GetQueryBuilder().
		Select(fs.queryFields...).
		From(store.FileInfoTableName).
		Where(squirrel.Eq{"Id": ids}).
		Where(squirrel.Eq{"DeleteAt": 0}).
		OrderBy("CreateAt DESC")

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "file_info_tosql")
	}

	var infos []*file.FileInfo
	if _, err := fs.GetReplica().Select(&infos, queryString, args...); err != nil {
		return nil, errors.Wrap(err, "failed to find FileInfos")
	}
	return infos, nil
}

func (fs *SqlFileInfoStore) Upsert(info *file.FileInfo) (*file.FileInfo, error) {
	info.PreSave()
	if err := info.IsValid(); err != nil {
		return nil, err
	}

	n, err := fs.GetMaster().Update(info)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update FileInfos")
	}
	if n == 0 {
		if err = fs.GetMaster().Insert(info); err != nil {
			return nil, errors.Wrap(err, "failed to save FileInfos")
		}
	}
	return info, nil
}

func (fs *SqlFileInfoStore) get(id string, fromMaster bool) (*file.FileInfo, error) {
	info := &file.FileInfo{}

	query := fs.GetQueryBuilder().
		Select(fs.queryFields...).
		From(store.FileInfoTableName).
		Where(squirrel.Eq{"Id": id}).
		Where(squirrel.Eq{"DeleteAt": 0})
	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "file_info_tosql")
	}

	db := fs.GetReplica()
	if fromMaster {
		db = fs.GetMaster()
	}

	if err := db.SelectOne(info, queryString, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("FileInfos", id)
		}
		return nil, errors.Wrapf(err, "failed to get FileInfos with id=%s")
	}
	return info, nil
}

func (fs *SqlFileInfoStore) Get(id string) (*file.FileInfo, error) {
	return fs.get(id, false)
}

func (fs *SqlFileInfoStore) GetFromMaster(id string) (*file.FileInfo, error) {
	return fs.get(id, true)
}

// GetWithOptions finds and returns fileinfos with given options.
// Leave page, perPage nil to get all result.
func (fs *SqlFileInfoStore) GetWithOptions(page, perPage *int, opt *file.GetFileInfosOptions) ([]*file.FileInfo, error) {
	if perPage != nil && *perPage < 0 {
		return nil, store.NewErrLimitExceeded("perPage", *perPage, "value used in pagination while getting FileInfos")
	} else if page != nil && *page < 0 {
		return nil, store.NewErrLimitExceeded("page", *page, "value used in pagination while getting FileInfos")
	}
	if *perPage == 0 {
		return nil, nil
	}

	if opt == nil {
		opt = &file.GetFileInfosOptions{}
	}

	query := fs.GetQueryBuilder().
		Select(fs.queryFields...).
		From(store.FileInfoTableName)

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
		query = query.Where("FileInfos.ParentID IN ?", opt.ParentID)
	}

	if opt.SortBy == "" {
		opt.SortBy = file.FILEINFO_SORT_BY_CREATED
	}
	sortDirection := "ASC"
	if opt.SortDescending {
		sortDirection = "DESC"
	}

	switch opt.SortBy {
	case file.FILEINFO_SORT_BY_CREATED:
		query = query.OrderBy("FileInfos.CreateAt " + sortDirection)
	case file.FILEINFO_SORT_BY_SIZE:
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
	var infos []*file.FileInfo
	if _, err := fs.GetReplica().Select(&infos, queryString, args...); err != nil {
		return nil, errors.Wrap(err, "failed to find FileInfos")
	}
	return infos, nil
}

func (fs *SqlFileInfoStore) GetByPath(path string) (*file.FileInfo, error) {
	info := new(file.FileInfo)

	query := fs.GetQueryBuilder().
		Select(fs.queryFields...).
		From(store.FileInfoTableName).
		Where(squirrel.Eq{"Path": path}).
		Where(squirrel.Eq{"DeleteAt": 0}).
		Limit(1)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "file_info_tosql")
	}

	if err := fs.GetReplica().SelectOne(info, queryString, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("FileInfos", fmt.Sprintf("path=%s", path))
		}

		return nil, errors.Wrapf(err, "failed to get FileInfos with path=%s", path)
	}
	return info, nil
}

func (fs *SqlFileInfoStore) InvalidateFileInfosForPostCache(postId string, deleted bool) {
}

func (fs *SqlFileInfoStore) GetForUser(userId string) ([]*file.FileInfo, error) {
	var infos []*file.FileInfo

	query := fs.GetQueryBuilder().
		Select(fs.queryFields...).
		From(store.FileInfoTableName).
		Where(squirrel.Eq{"CreatorId": userId}).
		Where(squirrel.Eq{"DeleteAt": 0}).
		OrderBy("CreateAt")
	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "file_info_tosql")
	}

	if _, err := fs.GetReplica().Select(&infos, queryString, args...); err != nil {
		return nil, errors.Wrapf(err, "failed to find FileInfos with creatorId=%s", userId)
	}

	return infos, nil
}

func (fs *SqlFileInfoStore) SetContent(fileId, content string) error {
	query := fs.GetQueryBuilder().
		Update(store.FileInfoTableName).
		Set("Content", content).
		Where(squirrel.Eq{"Id": fileId})
	queryString, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "file_info_tosql")
	}
	_, err = fs.GetMaster().Exec(queryString, args...)
	if err != nil {
		return errors.Wrapf(err, "failed to update FileInfos content with id=%s", fileId)
	}

	return nil
}

func (fs *SqlFileInfoStore) PermanentDelete(fileId string) error {
	if _, err := fs.GetMaster().Exec(
		`DELETE FROM
				FileInfos
			WHERE
				Id = :FileId`,
		map[string]interface{}{"FileId": fileId}); err != nil {
		return errors.Wrapf(err, "failed to delete FileInfos with id=%s", fileId)
	}
	return nil
}

func (fs *SqlFileInfoStore) PermanentDeleteBatch(endTime int64, limit int64) (int64, error) {
	sqlResult, err := fs.GetMaster().Exec(
		`DELETE FROM 
			FileInfos 
		WHERE Id = any (
			array (
				SELECT 
					Id 
				FROM 
					FileInfos 
				WHERE 
					CreateAt < :EndTime 
				LIMIT :Limit
			)
		)`,
		map[string]interface{}{
			"EndTime": endTime,
			"Limit":   limit,
		},
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
	query := "DELETE FROM FileInfos WHERE CreatorId = :CreatorId"
	sqlResult, err := fs.GetMaster().Exec(query, map[string]interface{}{
		"CreatorId": userId,
	})
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
		From("FileInfos").
		Where("DeleteAt = 0")

	queryString, args, err := query.ToSql()
	if err != nil {
		return int64(0), errors.Wrap(err, "count_tosql")
	}

	count, err := fs.GetReplica().SelectInt(queryString, args...)
	if err != nil {
		return int64(0), errors.Wrap(err, "failed to count Files")
	}
	return count, nil
}
