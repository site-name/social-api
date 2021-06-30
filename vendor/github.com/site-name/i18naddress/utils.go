package i18naddress

import "regexp"

// stringInSlice checks if given string presents in given slice
func stringInSlice(s string, slice *[]string) bool {
	for _, str := range *slice {
		if s == str {
			return true
		}
	}

	return false
}

// RegexesToStrings convert a slice of *Regexp(s) to a pointer to slice of string
func RegexesToStrings(in []*regexp.Regexp) *[]string {
	res := []string{}
	for _, rg := range in {
		res = append(res, rg.String())
	}

	return &res
}

// filterDuplicate filter all item(s) that appear(s) >= 2 times in given slice
func filterDuplicate(slice *[]string) *[]string {
	meetMap := make(map[string]bool)
	res := []string{}
	for _, str := range *slice {
		if _, met := meetMap[str]; !met {
			res = append(res, str)
			meetMap[str] = true
		}
	}

	return &res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// filterSlice filter from slice item(s) that does not satify given filter func
func filterSlice(slice []string, filter func(s string) bool) []string {
	res := []string{}
	for _, str := range slice {
		if filter(str) {
			res = append(res, str)
		}
	}

	return res
}
