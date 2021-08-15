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
		table := db.AddTableWithName(page.PageTranslation{}, store.PageTranslationtableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Title").SetMaxSize(page.PAGE_TITLE_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(model.LANGUAGE_CODE_MAX_LENGTH)
		table.ColMap("PageID").SetMaxSize(store.UUID_MAX_LENGTH)

		table.SetUniqueTogether("LanguageCode", "PageID")
		s.CommonSeoMaxLength(table)
	}
	return ps
}

func (ps *SqlPageTranslationStore) CreateIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_pages_title", store.PageTranslationtableName, "Title")
	ps.CreateIndexIfNotExists("idx_pages_title_lower_textpattern", store.PageTranslationtableName, "lower(Title) text_pattern_ops")
	ps.CreateForeignKeyIfNotExists(store.PageTranslationtableName, "PageID", store.PageTableName, "Id", true)
}
