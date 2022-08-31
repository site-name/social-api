package csv

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/api/gqlmodel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/product_and_discount"
	"github.com/sitename/sitename/store"
)

var (
	productFetchBatchSize uint64 = 10000
)

// ExportProducts is called by product export job, taks needed arguments then exports products
func (s *ServiceCsv) ExportProducts(input *gqlmodel.ExportProductsInput, delimeter string) *model.AppError {
	// if delimeter == "" {
	// 	delimeter = ";"
	// }

	// productFilterQuery := s.srv.Store.Product().AdvancedFilterQueryBuilder(input)

	// exportFields, fileHeaders, dataHeaders, appErr := s.GetExportFieldsAndHeadersInfo(*input.ExportInfo)
	// if appErr != nil {
	// 	return appErr
	// }

	// getFileName("product", strings.ToLower(string(input.FileType)))

	panic("not implt")
}

func (s *ServiceCsv) ExportProductsInBatches(productQuery squirrel.SelectBuilder, exportInfo gqlmodel.ExportInfoInput, exportFields []string, headers []string, delimiter string, fileType string) *model.AppError {
	var createAtGt int64 = 0

	for {
		prds, err := s.srv.Store.Product().FilterByQuery(productQuery.Where("Products.CreateAt > ?", createAtGt).Limit(productFetchBatchSize))
		if err != nil {
			return model.NewAppError("ExportProductsInBatches", "app.csv.error_finding_products_by_query.app_error", nil, err.Error(), http.StatusInternalServerError)
		}

		if len(prds) == 0 {
			break
		}

		// reset for later loop(s)
		createAtGt = prds[len(prds)-1].CreateAt

		products, appErr := s.srv.ProductService().ProductsByOption(&product_and_discount.ProductFilterOption{
			Id:                                       squirrel.Eq{store.ProductTableName + ".Id": prds.IDs()},
			PrefetchRelatedAssignedProductAttributes: true,
			PrefetchRelatedVariants:                  true,
			PrefetchRelatedCollections:               true,
			PrefetchRelatedMedia:                     true,
			PrefetchRelatedProductType:               true,
			PrefetchRelatedCategory:                  true,
		})
		if appErr != nil {
			if appErr.StatusCode == http.StatusInternalServerError {
				return appErr
			}
			break
		}

		s.GetProductsData(products, exportFields, exportInfo.Attributes, exportInfo.Warehouses, exportInfo.Channels)
	}

	panic("not implt")
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
