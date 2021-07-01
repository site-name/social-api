package product

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

type SqlProductVariantStore struct {
	store.Store
}

const (
	ProductVariantTableName = "ProductVariants"
)

func NewSqlProductVariantStore(s store.Store) store.ProductVariantStore {
	pvs := &SqlProductVariantStore{s}

	for _, db := range s.GetAllConns() {
		table := db.AddTableWithName(product_and_discount.ProductVariant{}, ProductVariantTableName).SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("ProductID").SetMaxSize(store.UUID_MAX_LENGTH)
		table.ColMap("Sku").SetMaxSize(product_and_discount.PRODUCT_VARIANT_SKU_MAX_LENGTH).SetUnique(true)
		table.ColMap("Name").SetMaxSize(product_and_discount.PRODUCT_VARIANT_NAME_MAX_LENGTH)
	}
	return pvs
}

func (ps *SqlProductVariantStore) CreateIndexesIfNotExists() {
	// ps.CreateIndexIfNotExists("idx_product_variants_name", ProductVariantTableName, "Name")
	// ps.CreateIndexIfNotExists("idx_product_variants_name_lower_textpattern", ProductVariantTableName, "lower(Name) text_pattern_ops")
	ps.CreateIndexIfNotExists("idx_product_variants_sku", ProductVariantTableName, "Sku")
}

func (ps *SqlProductVariantStore) Save(variant *product_and_discount.ProductVariant) (*product_and_discount.ProductVariant, error) {
	variant.PreSave()
	if err := variant.IsValid(); err != nil {
		return nil, err
	}

	if err := ps.GetMaster().Insert(variant); err != nil {
		if ps.IsUniqueConstraintError(err, []string{"Sku", "idx_productvariants_sku_unique", "productvariants_sku_key"}) {
			return nil, store.NewErrInvalidInput(ProductVariantTableName, "Sku", variant.Sku)
		}
		return nil, errors.Wrapf(err, "failed to save product variant with id=%s", variant.Id)
	}

	return variant, nil
}

func (ps *SqlProductVariantStore) Get(id string) (*product_and_discount.ProductVariant, error) {
	var variant product_and_discount.ProductVariant
	if err := ps.GetReplica().SelectOne(&variant, "SELECT * FROM "+ProductVariantTableName+" WHERE Id = :id",
		map[string]interface{}{"id": id}); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(ProductVariantTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to find product variant with id=%s", id)
	}

	return &variant, nil
}
