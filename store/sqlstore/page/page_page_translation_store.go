package page

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/page"
	"github.com/sitename/sitename/store"
)

type SqlPageTranslationStore struct {
	store.Store
}

func NewSqlPageTranslationStore(s store.Store) store.PageTranslationStore {
	ps := &SqlPageTranslationStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(page.PageTranslation{}, "PageTranslations").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("PageID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Title").SetMaxSize(page.PAGE_TITLE_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(model.LANGUAGE_CODE_MAX_LENGTH)

		s.CommonSeoMaxLength(table)
	}
	return ps
}

func (ps *SqlPageTranslationStore) CreateIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_pages_title", "PageTranslations", "Title")
	ps.CreateIndexIfNotExists("idx_pages_title_lower_textpattern", "PageTranslations", "lower(Title) text_pattern_ops")
}
