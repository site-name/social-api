package sqlstore

import (
	"github.com/sitename/sitename/model/page"
	"github.com/sitename/sitename/store"
)

type SqlPageTypeStore struct {
	*SqlStore
}

func newSqlPageTypeStore(s *SqlStore) store.PageTypeStore {
	ps := &SqlPageTypeStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(page.PageType{}, "PageTypes").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Name").SetMaxSize(page.PAGE_TYPE_NAME_MAX_LENGTH)
		table.ColMap("Slug").SetMaxSize(page.PAGE_TYPE_SLUG_MAX_LENGTH)
	}
	return ps
}

func (ps *SqlPageTypeStore) createIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_page_types_name", "PageTypes", "Title")
	ps.CreateIndexIfNotExists("idx_page_types_slug", "PageTypes", "Slug")

	ps.CreateIndexIfNotExists("idx_page_types_name_lower_textpattern", "PageTypes", "lower(Name) text_pattern_ops")
}
