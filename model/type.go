package model

import (
	"github.com/shopspring/decimal"
)

type StringMap map[string]string

type Money struct {
	Amount   *decimal.Decimal
	Currency string
}

// // EncodeMsg implements msgp.Encodable
// func (z StringMap) EncodeMsg(en *msgp.Writer) (err error) {
// 	err = en.WriteMapHeader(uint32(len(z)))
// 	if err != nil {
// 		err = msgp.WrapError(err)
// 		return
// 	}
// 	for zb0004, zb0005 := range z {
// 		err = en.WriteString(zb0004)
// 		if err != nil {
// 			err = msgp.WrapError(err)
// 			return
// 		}
// 		err = en.WriteString(zb0005)
// 		if err != nil {
// 			err = msgp.WrapError(err, zb0004)
// 			return
// 		}
// 	}
// 	return
// }
