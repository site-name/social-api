package plugin

import (
	"fmt"
)

func stringify(objects []any) []string {
	stringified := make([]string, len(objects))
	for i, object := range objects {
		stringified[i] = fmt.Sprintf("%+v", object)
	}
	return stringified
}

func toObjects(strings []string) []any {
	if strings == nil {
		return nil
	}
	objects := make([]any, len(strings))
	for i, str := range strings {
		objects[i] = str
	}
	return objects
}

func stringifyToObjects(objects []any) []any {
	return toObjects(stringify(objects))
}
