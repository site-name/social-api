package scalars

import (
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
)

// MarshalUUID turns an uuid to string
func MarshalUUID(u uuid.UUID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(u.String()))
	})
}

// UnmarshalUUID turns string into an uuid
func UnmarshalUUID(v interface{}) (uuid.UUID, error) {
	switch v := v.(type) {
	case string:
		return uuid.Parse(v)

	default:
		return uuid.Nil, fmt.Errorf("%T is not a uuid string", v)
	}
}
