package gqlmodel

// MapToGraphqlMetaDataItems converts a map of key-value into a slice of graphql MetadataItems
func MapToGraphqlMetaDataItems(m map[string]string) []*MetadataItem {
	if m == nil {
		return []*MetadataItem{}
	}

	res := make([]*MetadataItem, len(m))
	for key, value := range m {
		res = append(res, &MetadataItem{Key: key, Value: value})
	}

	return res
}

// MetaDataToStringMap converts a slice of *MetadataInput || *MetadataItem to map[string]string.
//
// Other types will result in an empty map
func MetaDataToStringMap(metaList interface{}) map[string]string {
	res := make(map[string]string)

	switch t := metaList.(type) {
	case []*MetadataInput:
		for _, input := range t {
			res[input.Key] = input.Value
		}
	case []*MetadataItem:
		for _, item := range t {
			res[item.Key] = item.Value
		}
	default:
		return res
	}

	return res
}
