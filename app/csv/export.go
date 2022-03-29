package csv

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/graphql/gqlmodel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

var (
	productFetchBatchSize uint64 = 10000
)

// ExportProducts is called by product export job, taks needed arguments then exports products
func (s *ServiceCsv) ExportProducts(input *gqlmodel.ExportProductsInput, delimeter string) *model.AppError {
	if delimeter == "" {
		delimeter = ";"
	}

	productFilterQuery := s.srv.Store.Product().AdvancedFilterQueryBuilder(input)

	exportFields, fileHeaders, dataHeaders, appErr := s.GetExportFieldsAndHeadersInfo(*input.ExportInfo)
	if appErr != nil {
		return appErr
	}

	getFileName("product", strings.ToLower(string(input.FileType)))

}

func (s *ServiceCsv) ExportProductsInBatches(productQuery squirrel.SelectBuilder, exportInfo gqlmodel.ExportInfoInput, exportFields []string, headers []string, delimiter string, fileType string) *model.AppError {
	var createAtGt int64 = 0

	// fetch products in batches with size 10000 each to prevent memory leaking
	for {
		products, err := s.srv.Store.Product().FilterByQuery(productQuery, &product_and_discount.ProductFilterByQueryOptions{
			CreateAt:                                 squirrel.Gt{store.ProductTableName + ".CreateAt": createAtGt},
			Limit:                                    &productFetchBatchSize,
			PrefetchRelatedAssignedProductAttributes: true,
			PrefetchRelatedVariants:                  true,
			PrefetchRelatedCollections:               true,
			PrefetchRelatedMedia:                     true,
			PrefetchRelatedProductType:               true,
			PrefetchRelatedCategory:                  true,
		})
		if err != nil {
			return model.NewAppError("ExportProductsInBatches", "app.csv.error_finding_products_by_query.app_error", nil, err.Error(), http.StatusInternalServerError)
		}

		if len(products) == 0 {
			break
		}

		createAtGt = products[len(products)-1].CreateAt

		s.GetProductsData(products, exportFields, exportInfo.Attributes, exportInfo.Warehouses, exportInfo.Channels)
	}

}

// getFileName returns a file name for exported file
func getFileName(modelName string, fileType string) string {
	return fmt.Sprintf(
		"%s_data_%s_%s.%s",
		modelName,
		time.Now().UTC().Format("02_Jan_2006_15_04_05"),
		model.NewRandomString(16),
		fileType,
	)
}
