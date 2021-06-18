package scalars

import (
	"io"
)

type PlaceHolder byte

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (y *PlaceHolder) UnmarshalGQL(_ interface{}) error {
	*y = '1'
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (y PlaceHolder) MarshalGQL(w io.Writer) {
	w.Write([]byte{'.'})
}
