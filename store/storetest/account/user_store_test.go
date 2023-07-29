package account

import (
	"testing"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
	"github.com/sitename/sitename/store/storetest"
	"github.com/stretchr/testify/require"
)

func TestUserStore(t *testing.T) {
	storetest.StoreTestWithSqlStore(t, func(t *testing.T, ss store.Store, s storetest.SqlStore) {
		users, err := ss.User().GetAllProfiles(&model.UserGetOptions{})
		require.NoError(t, err, "failed cleaning test users")

		for _, u := range users {
			err := ss.User().PermanentDelete(u.Id)
			require.NoError(t, err, "failed cleaning up test user %s", u.Username)
		}

		t.Run("IsEmpty", func(t *testing.T) { testIsEmpty(t, ss) })
	})
}

func testIsEmpty(t *testing.T, ss store.Store) {
	numOfUsers, err := ss.User().Count(model.UserCountOptions{})
	require.NoError(t, err)
	require.Equal(t, 0, numOfUsers, "expected 0 users, got %d", numOfUsers)

	u := &model.User{
		Email:    "leminhon2398@outlook.com",
		Username: model.NewId(),
	}
	u, err = ss.User().Save(u)
	require.NoError(t, err)

	numOfUsers, err = ss.User().Count(model.UserCountOptions{})
	require.NoError(t, err)
	require.Greater(t, numOfUsers, 0, "expected at least 1 user in database, got 0")

	err = ss.User().PermanentDelete(u.Id)
	require.NoError(t, err)

	numOfUsers, err = ss.User().Count(model.UserCountOptions{})
	require.NoError(t, err)
	require.Equal(t, 0, numOfUsers, "expected 0 users, got %d", numOfUsers)
}
