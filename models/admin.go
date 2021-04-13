package models

import (
	"fmt"

	"github.com/sitename/sitename/modules/log"
	"github.com/sitename/sitename/modules/storage"
	"github.com/sitename/sitename/modules/timeutil"
	"github.com/sitename/sitename/modules/util"
)

//NoticeType describes the notice type
type NoticeType int

const (
	//NoticeRepository type
	NoticeRepository NoticeType = iota + 1
	// NoticeTask type
	NoticeTask
)

// Notice represents a system notice for admin.
type Notice struct {
	ID          int64 `xorm:"pk autoincr"`
	Type        NoticeType
	Description string             `xorm:"TEXT"`
	CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
}

// TrStr returns a translation format string.
func (n *Notice) TrStr() string {
	return fmt.Sprintf("admin.notices.type_%d", n.Type)
}

// CreateNotice creates new system notice.
func CreateNotice(tp NoticeType, desc string, args ...interface{}) error {
	return createNotice(x, tp, desc, args...)
}

func createNotice(e Engine, tp NoticeType, desc string, args ...interface{}) error {
	if len(args) > 0 {
		desc = fmt.Sprintf(desc, args...)
	}
	n := &Notice{
		Type:        tp,
		Description: desc,
	}
	_, err := e.Insert(n)
	return err
}

// CreateRepositoryNotice creates new system notice with type NoticeRepository.
func CreateRepositoryNotice(desc string, args ...interface{}) error {
	return createNotice(x, NoticeRepository, desc, args...)
}

// RemoveAllWithNotice removes all directories in given path and
// creates a system notice when error occurs.
func RemoveAllWithNotice(title, path string) {
	removeAllWithNotice(x, title, path)
}

// RemoveStorageWithNotice removes a file from the storage and
// creates a system notice when error occurs.
func RemoveStorageWithNotice(bucket storage.ObjectStorage, title, path string) {
	removeStorageWithNotice(x, bucket, title, path)
}

func removeStorageWithNotice(e Engine, bucket storage.ObjectStorage, title, path string) {
	if err := bucket.Delete(path); err != nil {
		desc := fmt.Sprintf("%s [%s]: %v", title, path, err)
		log.Warn(title+" [%s]: %v", path, err)
		if err = createNotice(e, NoticeRepository, desc); err != nil {
			log.Error("CreateRepositoryNotice: %v", err)
		}
	}
}

func removeAllWithNotice(e Engine, title, path string) {
	if err := util.RemoveAll(path); err != nil {
		desc := fmt.Sprintf("%s [%s]: %v", title, path, err)
		log.Warn(title+" [%s]: %v", path, err)
		if err = createNotice(e, NoticeRepository, desc); err != nil {
			log.Error("CreateRepositoryNotice: %v", err)
		}
	}
}

// CountNotices returns number of notices.
func CountNotices() int64 {
	count, _ := x.Count(new(Notice))
	return count
}

// Notices returns notices in given page.
func Notices(page, pageSize int) ([]*Notice, error) {
	notices := make([]*Notice, 0, pageSize)
	return notices, x.
		Limit(pageSize, (page-1)*pageSize).
		Desc("id").
		Find(&notices)
}

// DeleteNotice deletes a system notice by given ID.
func DeleteNotice(id int64) error {
	_, err := x.ID(id).Delete(new(Notice))
	return err
}

// DeleteNotices deletes all notices with ID from start to end (inclusive).
func DeleteNotices(start, end int64) error {
	sess := x.Where("id >= ?", start)
	if end > 0 {
		sess.And("id <= ?", end)
	}
	_, err := sess.Delete(new(Notice))
	return err
}

// DeleteNoticesByIDs deletes notices by given IDs.
func DeleteNoticesByIDs(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	_, err := x.
		In("id", ids).
		Delete(new(Notice))
	return err
}

// GetAdminUser returns the first administrator
func GetAdminUser() (*User, error) {
	var admin User
	has, err := x.Where("is_admin=?", true).Get(&admin)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrUserNotExist{}
	}

	return &admin, nil
}
