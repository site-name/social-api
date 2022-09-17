package model

import (
	"io"
)

type UserAccessTokenSearch struct {
	Term string `json:"term"`
}

// ToJson convert a UserAccessTokenSearch to json string
func (c *UserAccessTokenSearch) ToJSON() string {
	return ModelToJson(c)
}

// UserAccessTokenSearchJson decodes the input and returns a UserAccessTokenSearch
func UserAccessTokenSearchFromJson(data io.Reader) *UserAccessTokenSearch {
	var u *UserAccessTokenSearch
	ModelFromJson(&u, data)
	return u
}
