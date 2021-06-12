package page

import (
	"github.com/sitename/sitename/model/page"
	"github.com/sitename/sitename/store"
)

type SqlPageStore struct {
	store.Store
}

func NewSqlPageStore(s store.Store) store.PageStore {
	ps := &SqlPageStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(page.Page{}, "Pages").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("PageTypeID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Title").SetMaxSize(page.PAGE_TITLE_MAX_LENGTH).SetUnique(true)
		table.ColMap("Slug").SetMaxSize(page.PAGE_SLUG_MAX_LENGTH).SetUnique(true)

		s.CommonSeoMaxLength(table)
	}
	return ps
}

func (ps *SqlPageStore) CreateIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_pages_title", "Pages", "Title")
	ps.CreateIndexIfNotExists("idx_pages_slug", "Pages", "Slug")

	ps.CreateIndexIfNotExists("idx_pages_title_lower_textpattern", "Pages", "lower(Title) text_pattern_ops")
}
