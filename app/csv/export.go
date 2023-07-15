package csv

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
)

var (
	productFetchBatchSize uint64 = 10000
)

// ExportProducts is called by product export job, taks needed arguments then exports products
func (s *ServiceCsv) ExportProducts(input *model.ExportProductsFilterOptions, delimeter string) *model.AppError {
	// if delimeter == "" {
	// 	delimeter = ";"
	// }

	// productFilterQuery := s.srv.Store.Product().AdvancedFilterQueryBuilder(input)

	// exportFields, fileHeaders, dataHeaders, appErr := s.GetExportFieldsAndHeadersInfo(*input.ExportInfo)
	// if appErr != nil {
	// 	return appErr
	// }

	// getFileName("product", strings.ToLower(string(input.FileType)))

	panic("not implemented")
}

func (s *ServiceCsv) ExportProductsInBatches(
	productQuery squirrel.SelectBuilder,
	ExportInfo struct {
		Attributes []string
		Warehouses []string
		Channels   []string
		Fields     []string
	},
	exportFields []string,
	headers []string,
	delimiter string,
	fileType string,
) *model.AppError {
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

		products, appErr := s.srv.ProductService().ProductsByOption(&model.ProductFilterOption{
			Id:                                       squirrel.Eq{model.ProductTableName + ".Id": prds.IDs()},
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

		s.GetProductsData(products, exportFields, ExportInfo.Attributes, ExportInfo.Warehouses, ExportInfo.Channels)
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
