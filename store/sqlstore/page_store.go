package sqlstore

import (
	"github.com/sitename/sitename/model/page"
	"github.com/sitename/sitename/store"
)

type SqlPageStore struct {
	*SqlStore
}

func newSqlPageStore(s *SqlStore) store.PageStore {
	ps := &SqlPageStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(page.Page{}, "Pages").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("PageTypeID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Title").SetMaxSize(page.PAGE_TITLE_MAX_LENGTH)
		table.ColMap("Slug").SetMaxSize(page.PAGE_SLUG_MAX_LENGTH)

		s.commonSeoMaxLength(table)
	}
	return ps
}

func (ps *SqlPageStore) createIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_pages_title", "Pages", "Title")
	ps.CreateIndexIfNotExists("idx_pages_slug", "Pages", "Slug")

	ps.CreateIndexIfNotExists("idx_pages_title_lower_textpattern", "Pages", "lower(Title) text_pattern_ops")
}
