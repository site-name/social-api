package account

import (
	"strings"
	"testing"

	"github.com/sitename/sitename/model"
	"github.com/stretchr/testify/require"
)

func TestUserAccessTokenSearchJson(t *testing.T) {
	userAccessTokenSearch := UserAccessTokenSearch{Term: model.NewId()}
	json := userAccessTokenSearch.ToJson()
	ruserAccessTokenSearch := UserAccessTokenSearchFromJson(strings.NewReader(json))
	require.Equal(t, userAccessTokenSearch.Term, ruserAccessTokenSearch.Term, "Terms do not match")
}
