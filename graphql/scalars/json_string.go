package scalars

import (
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/json"
)

func MarshalJSONString(v model.StringInterface) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if v == nil {
			w.Write([]byte(""))
			return
		}

		data, err := json.JSON.Marshal(v)
		if err != nil {
			w.Write([]byte(""))
			return
		}

		w.Write(data)
	})
}

func UnmarshalJSONString(v interface{}) (model.StringInterface, error) {
	switch value := v.(type) {
	case string:
		var res model.StringInterface
		err := json.JSON.Unmarshal([]byte(value), &res)
		if err != nil {
			return nil, err
		}

		return res, nil

	default:
		return model.StringInterface{}, fmt.Errorf("unknown type of given value: %T", v)
	}
}
