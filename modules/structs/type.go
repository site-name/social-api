package structs

type VisibleType int

const (
	// VisibleTypePublic Visible for everyone
	VisibleTypePublic VisibleType = iota

	// VisibleTypeLimited Visible for every connected user
	VisibleTypeLimited

	// VisibleTypePrivate Visible only for organization's members
	VisibleTypePrivate
)

// VisibilityModes is a map of org Visibility types
var VisibilityModes = map[string]VisibleType{
	"public":  VisibleTypePublic,
	"limited": VisibleTypeLimited,
	"private": VisibleTypePrivate,
}

// ExtractKeysFromMapString provides a slice of keys from map
func ExtractKeysFromMapString(in map[string]VisibleType) (keys []string) {
	for k := range in {
		keys = append(keys, k)
	}
	return
}
