package discount

import (
	"database/sql"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SqlVoucherTranslationStore struct {
	store.Store
}

func NewSqlVoucherTranslationStore(sqlStore store.Store) store.VoucherTranslationStore {
	return &SqlVoucherTranslationStore{sqlStore}
}

func (vts *SqlVoucherTranslationStore) Upsert(translation model.VoucherTranslation) (*model.VoucherTranslation, error) {
	isSaving := false
	if translation.ID == "" {
		isSaving = true
		model_helper.VoucherTranslationPreSave(&translation)
	} else {
		model_helper.VoucherTranslationCommonPre(&translation)
	}

	if err := model_helper.VoucherTranslationIsValid(translation); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = translation.Insert(vts.GetMaster(), boil.Infer())
	} else {
		_, err = translation.Update(vts.GetMaster(), boil.Blacklist(model.VoucherTranslationColumns.CreatedAt))
	}

	if err != nil {
		if vts.IsUniqueConstraintError(err, []string{model.VoucherTranslationColumns.VoucherID, "voucher_translations_language_code_voucher_id_key"}) {
			return nil, store.NewErrInvalidInput(model.TableNames.VoucherTranslations, model.VoucherTranslationColumns.VoucherID+"/"+model.VoucherTranslationColumns.LanguageCode, "unique")
		}
		return nil, err
	}

	return &translation, nil
}

func (vts *SqlVoucherTranslationStore) Get(id string) (*model.VoucherTranslation, error) {
	translation, err := model.FindVoucherTranslation(vts.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.VoucherTranslations, id)
		}
		return nil, err
	}

	return translation, nil
}

func (vts *SqlVoucherTranslationStore) FilterByOption(option model_helper.VoucherTranslationFilterOption) (model.VoucherTranslationSlice, error) {
	return model.VoucherTranslations(option.Conditions...).All(vts.GetReplica())
}

func (vts *SqlVoucherTranslationStore) GetByOption(option model_helper.VoucherTranslationFilterOption) (*model.VoucherTranslation, error) {
	translation, err := model.VoucherTranslations(option.Conditions...).One(vts.GetReplica())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.VoucherTranslations, "options")
		}
		return nil, err
	}

	return translation, nil
}
