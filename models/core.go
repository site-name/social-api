package models

type SortableModel struct {
	ID        int64 `xorm:"pr autoincr"`
	SortOrder int   `xorm:"INDEX"`
}
