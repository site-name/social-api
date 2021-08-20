package product

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlDigitalContentUrlStore struct {
	store.Store
}

func NewSqlDigitalContentUrlStore(s store.Store) store.DigitalContentUrlStore {
	dcs := &SqlDigitalContentUrlStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.DigitalContentUrl{}, store.ProductDigitalContentURLTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Token").SetMaxSize(product_and_discount.DIGITAL_CONTENT_URL_TOKEN_MAX_LENGTH).SetUnique(true)
		table.ColMap("ContentID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("LineID").SetMaxSize(store.UUID_MAX_LENGTH).SetUnique(true)
	}
	return dcs
}

func (ps *SqlDigitalContentUrlStore) CreateIndexesIfNotExists() {
	ps.CreateForeignKeyIfNotExists(store.ProductDigitalContentURLTableName, "ContentID", store.ProductDigitalContentTableName, "Id", true)
	ps.CreateForeignKeyIfNotExists(store.ProductDigitalContentURLTableName, "LineID", store.OrderLineTableName, "Id", true)
}

// Save insert given digital content url into database then returns it
func (ps *SqlDigitalContentUrlStore) Save(contentURL *product_and_discount.DigitalContentUrl) (*product_and_discount.DigitalContentUrl, error) {

	for {
		contentURL.PreSave()
		if err := contentURL.IsValid(); err != nil {
			return nil, err
		}

		err := ps.GetMaster().Insert(contentURL)
		if err != nil {
			if ps.IsUniqueConstraintError(err, []string{"Token", "digitalcontenturls_token_key"}) {
				contentURL.NewToken(true)
				continue
			}
			if ps.IsUniqueConstraintError(err, []string{"LineID", "digitalcontenturls_lineid_key"}) {
				return nil, store.NewErrInvalidInput(store.ProductDigitalContentURLTableName, "LinesID", contentURL.LineID)
			}
			return nil, errors.Wrapf(err, "failed to save digital content url with given id=%s", contentURL.Id)
		}

		return contentURL, nil
	}
}

// Get finds and returns a digital content url with given id
func (ps *SqlDigitalContentUrlStore) Get(id string) (*product_and_discount.DigitalContentUrl, error) {
	var res *product_and_discount.DigitalContentUrl
	err := ps.GetReplica().SelectOne(&res, "SELECT * FROM "+store.ProductDigitalContentURLTableName+" WHERE Id = :ID", map[string]interface{}{"ID": id})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ProductDigitalContentURLTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find digital content url with id=%s", id)
	}

	return res, nil
}
