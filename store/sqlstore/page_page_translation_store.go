package sqlstore

import (
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/page"
	"github.com/sitename/sitename/store"
)

type SqlPageTranslationStore struct {
	*SqlStore
}

func newSqlPageTranslationStore(s *SqlStore) store.PageTranslationStore {
	ps := &SqlPageTranslationStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(page.PageTranslation{}, "PageTranslations").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("PageID").SetMaxSize(UUID_MAX_LENGTH)
		table.ColMap("Title").SetMaxSize(page.PAGE_TITLE_MAX_LENGTH)
		table.ColMap("LanguageCode").SetMaxSize(model.LANGUAGE_CODE_MAX_LENGTH).
			SetDefaultConstraint(model.NewString(model.DEFAULT_LANGUAGE_CODE))

		s.commonSeoMaxLength(table)
	}
	return ps
}

func (ps *SqlPageTranslationStore) createIndexesIfNotExists() {
	ps.CreateIndexIfNotExists("idx_pages_title", "PageTranslations", "Title")
	ps.CreateIndexIfNotExists("idx_pages_title_lower_textpattern", "PageTranslations", "lower(Title) text_pattern_ops")
}
