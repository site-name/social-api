package page

import (
	"github.com/sitename/sitename/model/page"
	"github.com/sitename/sitename/store"
)

type SqlPageTypeStore struct {
	store.Store
}

func NewSqlPageTypeStore(s store.Store) store.PageTypeStore {
	ps := &SqlPageTypeStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(page.PageType{}, "PageTypes").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(page.PAGE_TYPE_NAME_MAX_LENGTH).SetUnique(true)
		table.ColMap("Slug").SetMaxSize(page.PAGE_TYPE_SLUG_MAX_LENGTH).SetUnique(true)
	}
	return ps
}

func (ps *SqlPageTypeStore) CreateIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_page_types_name", "PageTypes", "Name")
	ps.CreateIndexIfNotExists("idx_page_types_slug", "PageTypes", "Slug")

	ps.CreateIndexIfNotExists("idx_page_types_name_lower_textpattern", "PageTypes", "lower(Name) text_pattern_ops")
}