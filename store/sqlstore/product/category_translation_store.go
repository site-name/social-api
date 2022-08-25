package product

import (
	"github.com/sitename/sitename/store"
)

type SqlCategoryTranslationStore struct {
	store.Store
}

func NewSqlCategoryTranslationStore(s store.Store) store.CategoryTranslationStore {
	return &SqlCategoryTranslationStore{s}
}
