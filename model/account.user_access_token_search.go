package model

type UserAccessTokenSearch struct {
	Term string `json:"term"`
}

// ToJson convert a UserAccessTokenSearch to json string
func (c *UserAccessTokenSearch) ToJSON() string {
	return ModelToJson(c)
}
