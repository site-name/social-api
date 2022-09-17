package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserAccessTokenSearchJson(t *testing.T) {
	userAccessTokenSearch := UserAccessTokenSearch{Term: NewId()}
	json := userAccessTokenSearch.ToJSON()
	ruserAccessTokenSearch := UserAccessTokenSearchFromJson(strings.NewReader(json))
	require.Equal(t, userAccessTokenSearch.Term, ruserAccessTokenSearch.Term, "Terms do not match")
}
