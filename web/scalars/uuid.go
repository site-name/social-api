package scalars

import (
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
)

func MarshalUUID(u uuid.UUID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(u.String()))
	})
}

func UnmarshalUUID(v interface{}) (uuid.UUID, error) {
	switch v := v.(type) {
	case string:
		return uuid.Parse(v)

	default:
		return uuid.Nil, fmt.Errorf("%T is not a uuid string", v)
	}
}
