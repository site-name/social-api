package account

import (
	"testing"

	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/storetest"
	"github.com/stretchr/testify/require"
)

func TestUserStore(t *testing.T, ss store.Store, s storetest.SqlStore) {
	users, err := ss.User().GetAll()
	require.NoError(t, err, "failed cleaning up test users")

	for _, u := range users {
		err := ss.User().PermanentDelete(u.Id)
		require.NoError(t, err, "failed cleaning up test user %s", u.Username)
	}
}
