package account

import (
	"testing"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/storetest"
	"github.com/stretchr/testify/require"
)

func TestUserStore(t *testing.T) {
	storetest.StoreTestWithSqlStore(t, func(t *testing.T, ss store.Store, s storetest.SqlStore) {
		users, err := ss.User().GetAllProfiles(model_helper.UserGetOptions{
			Sort: model.UserTableColumns.Username + " ASC",
		})
		require.NoError(t, err, "failed cleaning up test users")

		for _, user := range users {
			err := ss.User().PermanentDelete(user.ID)
			require.NoError(t, err, "failed cleaning up test user %s", user.Username)
		}

		t.Run("IsEmpty", func(t *testing.T) { testIsEmpty(t, ss) })
	})
}

func testIsEmpty(t *testing.T, ss store.Store) {
	empty, err := ss.User().IsEmpty()
	require.NoError(t, err)
	require.True(t, empty)

	u := model.User{
		Email: storetest.MakeEmail(),
		ID:    model_helper.NewId(),
	}

	savedU, err := ss.User().Save(u)
	require.NoError(t, err)

	empty, err = ss.User().IsEmpty()
	require.NoError(t, err)
	require.False(t, empty)

	err = ss.User().PermanentDelete(savedU.ID)
	require.NoError(t, err)

	empty, err = ss.User().IsEmpty()
	require.NoError(t, err)
	require.True(t, empty)
}
