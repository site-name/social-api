package product

func stringPointerSliceToStringSlice(sl []*string) []string {
	res := []string{}

	if sl == nil {
		return res
	}

	for _, item := range sl {
		if item != nil {
			res = append(res, *item)
		}
	}

	return res
}
