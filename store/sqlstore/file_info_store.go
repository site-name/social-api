package sqlstore

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlFileInfoStore struct {
	*SqlStore
	metrics     einterfaces.MetricsInterface
	queryFields []string
}

func (fs *SqlFileInfoStore) ClearCaches() {

}

func newSqlFileInfoStore(sqlStore *SqlStore, metrics einterfaces.MetricsInterface) store.FileInfoStore {
	s := &SqlFileInfoStore{
		SqlStore: sqlStore,
		metrics:  metrics,
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
		table := db.AddTableWithName(model.FileInfo{}, "FileInfos").SetKeys(false, "id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("CreatorId").SetMaxSize(UUID_MAX_LENGTH)
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

func (fs *SqlFileInfoStore) createIndexesIfNotExists() {
	fs.CreateIndexIfNotExists("idx_fileinfo_update_at", "FileInfo", "UpdateAt")
	fs.CreateIndexIfNotExists("idx_fileinfo_create_at", "FileInfo", "CreateAt")
	fs.CreateIndexIfNotExists("idx_fileinfo_delete_at", "FileInfo", "DeleteAt")
	fs.CreateIndexIfNotExists("idx_fileinfo_postid_at", "FileInfo", "PostId")
	fs.CreateIndexIfNotExists("idx_fileinfo_extension_at", "FileInfo", "Extension")
	fs.CreateFullTextIndexIfNotExists("idx_fileinfo_name_txt", "FileInfo", "Name")
	fs.CreateFullTextFuncIndexIfNotExists("idx_fileinfo_name_splitted", "FileInfo", "Translate(Name, '.,-', '   ')")
	fs.CreateFullTextIndexIfNotExists("idx_fileinfo_content_txt", "FileInfo", "Content")
}

func (fs SqlFileInfoStore) Save(info *model.FileInfo) (*model.FileInfo, error) {
	info.PreSave()
	if err := info.IsValid(); err != nil {
		return nil, err
	}

	if err := fs.GetMaster().Insert(info); err != nil {
		return nil, errors.Wrap(err, "failed to save FileInfo")
	}
	return info, nil
}

func (fs SqlFileInfoStore) GetByIds(ids []string) ([]*model.FileInfo, error) {
	query := fs.getQueryBuilder().
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

func (fs SqlFileInfoStore) Upsert(info *model.FileInfo) (*model.FileInfo, error) {
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

func (fs SqlFileInfoStore) get(id string, fromMaster bool) (*model.FileInfo, error) {
	info := &model.FileInfo{}

	query := fs.getQueryBuilder().
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

func (fs SqlFileInfoStore) Get(id string) (*model.FileInfo, error) {
	return fs.get(id, false)
}

func (fs SqlFileInfoStore) GetFromMaster(id string) (*model.FileInfo, error) {
	return fs.get(id, true)
}
