package file

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
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
		"FileInfo.Id",
		"FileInfo.CreatorId",
		// "FileInfo.PostId",
		"FileInfo.CreateAt",
		"FileInfo.UpdateAt",
		"FileInfo.DeleteAt",
		"FileInfo.Path",
		"FileInfo.ThumbnailPath",
		"FileInfo.PreviewPath",
		"FileInfo.Name",
		"FileInfo.Extension",
		"FileInfo.Size",
		"FileInfo.MimeType",
		"FileInfo.Width",
		"FileInfo.Height",
		"FileInfo.HasPreviewImage",
		"FileInfo.MiniPreview",
		"Coalesce(FileInfo.Content, '') AS Content",
		"Coalesce(FileInfo.RemoteId, '') AS RemoteId",
	}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.FileInfo{}, "FileInfos").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("CreatorId").SetMaxSize(store.UUID_MAX_LENGTH)
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
	fs.CreateIndexIfNotExists("idx_fileinfo_update_at", "FileInfos", "UpdateAt")
	fs.CreateIndexIfNotExists("idx_fileinfo_create_at", "FileInfos", "CreateAt")
	fs.CreateIndexIfNotExists("idx_fileinfo_delete_at", "FileInfos", "DeleteAt")
	// fs.CreateIndexIfNotExists("idx_fileinfo_postid_at", "FileInfos", "PostId")
	fs.CreateIndexIfNotExists("idx_fileinfo_extension_at", "FileInfos", "Extension")
	fs.CreateFullTextIndexIfNotExists("idx_fileinfo_name_txt", "FileInfos", "Name")
	fs.CreateFullTextFuncIndexIfNotExists("idx_fileinfo_name_splitted", "FileInfos", "Translate(Name, '.,-', '   ')")
	fs.CreateFullTextIndexIfNotExists("idx_fileinfo_content_txt", "FileInfos", "Content")
}

func (fs *SqlFileInfoStore) Save(info *model.FileInfo) (*model.FileInfo, error) {
	info.PreSave()
	if err := info.IsValid(); err != nil {
		return nil, err
	}

	if err := fs.GetMaster().Insert(info); err != nil {
		return nil, errors.Wrap(err, "failed to save FileInfo")
	}
	return info, nil
}

func (fs *SqlFileInfoStore) GetByIds(ids []string) ([]*model.FileInfo, error) {
	query := fs.GetQueryBuilder().
		Select(fs.queryFields...).
		From("FileInfos").
		Where(squirrel.Eq{"Id": ids}).
		Where(squirrel.Eq{"DeleteAt": 0}).
		OrderBy("CreateAt DESC")

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "file_info_tosql")
	}

	var infos []*model.FileInfo
	if _, err := fs.GetReplica().Select(&infos, queryString, args...); err != nil {
		return nil, errors.Wrap(err, "failed to find FileInfos")
	}
	return infos, nil
}

func (fs *SqlFileInfoStore) Upsert(info *model.FileInfo) (*model.FileInfo, error) {
	info.PreSave()
	if err := info.IsValid(); err != nil {
		return nil, err
	}

	n, err := fs.GetMaster().Update(info)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update FileInfo")
	}
	if n == 0 {
		if err = fs.GetMaster().Insert(info); err != nil {
			return nil, errors.Wrap(err, "failed to save FileInfo")
		}
	}
	return info, nil
}

func (fs *SqlFileInfoStore) get(id string, fromMaster bool) (*model.FileInfo, error) {
	info := &model.FileInfo{}

	query := fs.GetQueryBuilder().
		Select(fs.queryFields...).
		From("FileInfo").
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
			return nil, store.NewErrNotFound("FileInfo", id)
		}
		return nil, errors.Wrapf(err, "failed to get FileInfo with id=%s")
	}
	return info, nil
}

func (fs *SqlFileInfoStore) Get(id string) (*model.FileInfo, error) {
	return fs.get(id, false)
}

func (fs *SqlFileInfoStore) GetFromMaster(id string) (*model.FileInfo, error) {
	return fs.get(id, true)
}

func (fs *SqlFileInfoStore) GetWithOptions(page, perPage int, opt *model.GetFileInfosOptions) ([]*model.FileInfo, error) {
	if perPage < 0 {
		return nil, store.NewErrLimitExceeded("perPage", perPage, "value used in pagination while getting FileInfos")
	} else if page < 0 {
		return nil, store.NewErrLimitExceeded("page", page, "value used in pagination while getting FileInfos")
	}
	if perPage == 0 {
		return nil, nil
	}

	if opt == nil {
		opt = &model.GetFileInfosOptions{}
	}

	query := fs.GetQueryBuilder().
		Select(fs.queryFields...).
		From("FileInfo")

	// if len(opt.ChannelIds) > 0 {
	// 	query = query.Join("Posts ON FileInfo.PostId = Posts.Id").
	// 		Where(sq.Eq{"Posts.ChannelId": opt.ChannelIds})
	// }

	if len(opt.UserIds) > 0 {
		query = query.Where(squirrel.Eq{"FileInfo.CreatorId": opt.UserIds})
	}

	if opt.Since > 0 {
		query = query.Where(squirrel.GtOrEq{"FileInfo.CreateAt": opt.Since})
	}

	if !opt.IncludeDeleted {
		query = query.Where("FileInfo.DeleteAt = 0")
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
		query = query.OrderBy("FileInfo.CreateAt " + sortDirection)
	case model.FILEINFO_SORT_BY_SIZE:
		query = query.OrderBy("FileInfo.Size " + sortDirection)
	default:
		return nil, store.NewErrInvalidInput("FileInfo", "<sortOption>", opt.SortBy)
	}

	query = query.OrderBy("FileInfo.Id ASC") // secondary sort for sort stability

	query = query.Limit(uint64(perPage)).Offset(uint64(perPage * page))

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "file_info_tosql")
	}
	var infos []*model.FileInfo
	if _, err := fs.GetReplica().Select(&infos, queryString, args...); err != nil {
		return nil, errors.Wrap(err, "failed to find FileInfos")
	}
	return infos, nil
}

func (fs *SqlFileInfoStore) GetByPath(path string) (*model.FileInfo, error) {
	info := new(model.FileInfo)

	query := fs.GetQueryBuilder().
		Select(fs.queryFields...).
		From("FileInfos").
		Where(squirrel.Eq{"Path": path}).
		Where(squirrel.Eq{"DeleteAt": 0}).
		Limit(1)

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "file_info_tosql")
	}

	if err := fs.GetReplica().SelectOne(info, queryString, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("FileInfo", fmt.Sprintf("path=%s", path))
		}

		return nil, errors.Wrapf(err, "failed to get FileInfo with path=%s", path)
	}
	return info, nil
}

func (fs *SqlFileInfoStore) InvalidateFileInfosForPostCache(postId string, deleted bool) {
}

func (fs *SqlFileInfoStore) GetForUser(userId string) ([]*model.FileInfo, error) {
	var infos []*model.FileInfo

	query := fs.GetQueryBuilder().
		Select(fs.queryFields...).
		From("FileInfos").
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
		Update("FileInfos").
		Set("Content", content).
		Where(squirrel.Eq{"Id": fileId})
	queryString, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "file_info_tosql")
	}
	_, err = fs.GetMaster().Exec(queryString, args...)
	if err != nil {
		return errors.Wrapf(err, "failed to update FileInfo content with id=%s", fileId)
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
		return errors.Wrapf(err, "failed to delete FileInfo with id=%s", fileId)
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
		return 0, errors.Wrapf(err, "failed to delete FileInfo with creatorId=%s", userId)
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
		From("FileInfo").
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
