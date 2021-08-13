package order

// import (
// 	"github.com/sitename/sitename/model"
// 	"github.com/sitename/sitename/model/account"
// 	"github.com/sitename/sitename/model/order"
// )

// func (a *AppOrder) OrderReturned(ord *order.Order, user *account.User, returnedLines []*QuantityOrderLine) *model.AppError {
// 	var userID *string
// 	if user == nil || !model.IsValidId(user.Id) {
// 		userID = nil
// 	} else {
// 		userID = &user.Id
// 	}

// 	_, appErr := a.CommonCreateOrderEvent(&order.OrderEventOption{
// 		OrderID: ord.Id,
// 		Type:    order.ORDER_EVENT_TYPE__FULFILLMENT_RETURNED,
// 		UserID:  userID,
// 		Parameters: &model.StringInterface{
// 			"lines": linesPerQuantityToLineObjectList(returnedLines),
// 		},
// 	})
// 	if appErr != nil {
// 		return appErr
// 	}
// }
