// Code generated by "make app-layers"
// DO NOT EDIT

package sub_app_iface

import (
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
)

// InvoiceService contains methods for working with invoices
type InvoiceService interface {
	// SendInvoice Send an invoice to user of related order with URL to download it
	SendInvoice(inVoice model.Invoice, staffUser model.User, _ any, manager interfaces.PluginManagerInterface) *model_helper.AppError
	FilterInvoicesByOptions(options model_helper.InvoiceFilterOption) (model.InvoiceSlice, *model_helper.AppError)
	GetInvoiceByOptions(options model_helper.InvoiceFilterOption) (*model.Invoice, *model_helper.AppError)
	UpsertInvoice(invoice model.Invoice) (*model.Invoice, *model_helper.AppError)
	UpsertInvoiceEvent(option model_helper.InvoiceEventCreationOptions) (*model.InvoiceEvent, *model_helper.AppError)
}
