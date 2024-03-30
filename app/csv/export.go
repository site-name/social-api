package csv

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mattermost/squirrel"
	"github.com/samber/lo"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
)

var (
	productFetchBatchSize uint64 = 10000
)

func (s *ServiceCsv) ExportProducts(input model_helper.ExportProductsFilterOptions, delimeter string) *model_helper.AppError {
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

// NOTE: ordering by created_at should be applied to `productQuery`
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
) *model_helper.AppError {
	var createAtGt int64 = 0

	for {
		products, err := s.srv.Store.
			Product().
			FilterByQuery(
				productQuery.
					Where(squirrel.Gt{
						model.ProductTableColumns.CreatedAt: createAtGt,
					}).
					Limit(productFetchBatchSize),
			)
		if err != nil {
			return model_helper.NewAppError("ExportProductsInBatches", "app.csv.error_finding_products_by_query.app_error", nil, err.Error(), http.StatusInternalServerError)
		}

		if len(products) == 0 {
			break
		}

		// reset for later loop(s)
		createAtGt = products[len(products)-1].CreatedAt

		productIDs := lo.Map(products, func(p *model.Product, _ int) string { return p.ID })
		products, appErr := s.srv.Product.ProductsByOption(model_helper.ProductFilterOption{
			// Conditions: squirrel.Eq{model.ProductTableName + ".Id": products.IDs()},
			// Preloads: []string{
			// 	"ProductMedias",
			// 	"ProductType",
			// 	"Category",
			// 	"Collections",
			// 	"ProductVariants",
			// 	"Attributes",
			// },
			CommonQueryOptions: model_helper.NewCommonQueryOptions(model.ProductWhere.ID.IN(productIDs)),
			Preloads: []string{
				model.ProductRels.ProductMedia,
				model.ProductRels.Category,
				model.ProductRels.ProductCollections + "." + model.ProductCollectionRels.Collection,
				model.ProductRels.ProductVariants,
				// model.ProductRels.AssignedProductAttributes + "." + model.AssignedProductAttributeRels.Attribute,
			},
		})
		// model.ProductCollection
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
		model_helper.NewRandomString(16),
		fileType,
	)
}
