package dataloaders

import "github.com/graph-gophers/dataloader"

func dataloaderKeysToStringSlice(keys dataloader.Keys) []string {
	res := make([]string, len(keys))
	for i, v := range keys {
		res[i] = v.String()
	}

	return res
}
