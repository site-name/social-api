package page

const (
	PAGE_TYPE_NAME_MAX_LENGTH = 250
	PAGE_TYPE_SLUG_MAX_LENGTH = 255
)

type PageType struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"alug"`
}
