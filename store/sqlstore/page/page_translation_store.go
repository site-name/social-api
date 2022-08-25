package page

import (
	"github.com/sitename/sitename/store"
)

type SqlPageTranslationStore struct {
	store.Store
}

func NewSqlPageTranslationStore(s store.Store) store.PageTranslationStore {
	return &SqlPageTranslationStore{s}
}
