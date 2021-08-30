package storetest

import (
	"testing"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/store"
	"github.com/stretchr/testify/require"
)

const (
	DayMilliseconds   = 24 * 60 * 60 * 1000
	MonthMilliseconds = 31 * DayMilliseconds
)

func cleanupStatusStore(t *testing.T, s SqlStore) {
	_, execerr := s.GetMaster().ExecNoTimeout(` DELETE FROM Status `)
	require.NoError(t, execerr)
}

func TestUserStore(t *testing.T, ss store.Store, s SqlStore) {
	users, err := ss.User().GetAll()
	require.NoError(t, err, "failed cleaning up test users")

	for _, u := range users {
		err := ss.User().PermanentDelete(u.Id)
		require.NoError(t, err, "failed cleaning up test user %s", u.Username)
	}

	t.Run("Count", func(t *testing.T) { testCount(t, ss) })

}

func testCount(t *testing.T, ss store.Store) {
	// Regular
	// teamId := model.NewId()
	// channelId := model.NewId()
	regularUser := &account.User{}
	regularUser.Email = MakeEmail()
	regularUser.Roles = model.SystemUserRoleId
	_, err := ss.User().Save(regularUser)
	require.NoError(t, err)
	defer func() { require.NoError(t, ss.User().PermanentDelete(regularUser.Id)) }()

	// _, nErr := ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: regularUser.Id, SchemeAdmin: false, SchemeUser: true}, -1)
	// require.NoError(t, nErr)
	// _, nErr = ss.Channel().SaveMember(&model.ChannelMember{UserId: regularUser.Id, ChannelId: channelId, SchemeAdmin: false, SchemeUser: true, NotifyProps: model.GetDefaultChannelNotifyProps()})
	// require.NoError(t, nErr)

	guestUser := &account.User{}
	guestUser.Email = MakeEmail()
	guestUser.Roles = model.SystemGuestRoleId
	_, err = ss.User().Save(guestUser)
	require.NoError(t, err)
	defer func() { require.NoError(t, ss.User().PermanentDelete(guestUser.Id)) }()

	// _, nErr = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: guestUser.Id, SchemeAdmin: false, SchemeUser: false, SchemeGuest: true}, -1)
	// require.NoError(t, nErr)
	// _, nErr = ss.Channel().SaveMember(&model.ChannelMember{UserId: guestUser.Id, ChannelId: channelId, SchemeAdmin: false, SchemeUser: false, SchemeGuest: true, NotifyProps: model.GetDefaultChannelNotifyProps()})
	// require.NoError(t, nErr)

	teamAdmin := &account.User{}
	teamAdmin.Email = MakeEmail()
	teamAdmin.Roles = model.SystemUserRoleId
	_, err = ss.User().Save(teamAdmin)
	require.NoError(t, err)
	defer func() { require.NoError(t, ss.User().PermanentDelete(teamAdmin.Id)) }()

	// _, nErr = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: teamAdmin.Id, SchemeAdmin: true, SchemeUser: true}, -1)
	// require.NoError(t, nErr)
	// _, nErr = ss.Channel().SaveMember(&model.ChannelMember{UserId: teamAdmin.Id, ChannelId: channelId, SchemeAdmin: true, SchemeUser: true, NotifyProps: model.GetDefaultChannelNotifyProps()})
	// require.NoError(t, nErr)

	sysAdmin := &account.User{}
	sysAdmin.Email = MakeEmail()
	sysAdmin.Roles = model.SystemAdminRoleId + " " + model.SystemUserRoleId
	_, err = ss.User().Save(sysAdmin)
	require.NoError(t, err)
	defer func() { require.NoError(t, ss.User().PermanentDelete(sysAdmin.Id)) }()

	// _, nErr = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: sysAdmin.Id, SchemeAdmin: false, SchemeUser: true}, -1)
	// require.NoError(t, nErr)
	// _, nErr = ss.Channel().SaveMember(&model.ChannelMember{UserId: sysAdmin.Id, ChannelId: channelId, SchemeAdmin: true, SchemeUser: true, NotifyProps: model.GetDefaultChannelNotifyProps()})
	// require.NoError(t, nErr)

	// Deleted
	deletedUser := &account.User{}
	deletedUser.Email = MakeEmail()
	deletedUser.DeleteAt = model.GetMillis()
	_, err = ss.User().Save(deletedUser)
	require.NoError(t, err)
	defer func() { require.NoError(t, ss.User().PermanentDelete(deletedUser.Id)) }()

	// Bot
	botUser, err := ss.User().Save(&account.User{
		Email: MakeEmail(),
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, ss.User().PermanentDelete(botUser.Id)) }()

	// _, nErr = ss.Bot().Save(&model.Bot{
	// 	UserId:   botUser.Id,
	// 	Username: botUser.Username,
	// 	OwnerId:  regularUser.Id,
	// })
	// require.NoError(t, nErr)
	// botUser.IsBot = true
	// defer func() { require.NoError(t, ss.Bot().PermanentDelete(botUser.Id)) }()

	testCases := []struct {
		Description string
		Options     account.UserCountOptions
		Expected    int64
	}{
		{
			"No bot accounts no deleted accounts and no team id",
			account.UserCountOptions{
				// IncludeBotAccounts: false,
				IncludeDeleted: false,
				// TeamId:             "",
			},
			4,
		},
		{
			"Include bot accounts no deleted accounts and no team id",
			account.UserCountOptions{
				// IncludeBotAccounts: true,
				IncludeDeleted: false,
				// TeamId:             "",
			},
			5,
		},
		{
			"Include delete accounts no bots and no team id",
			account.UserCountOptions{
				// IncludeBotAccounts: false,
				IncludeDeleted: true,
				// TeamId:             "",
			},
			5,
		},
		{
			"Include bot accounts and deleted accounts and no team id",
			account.UserCountOptions{
				// IncludeBotAccounts: true,
				IncludeDeleted: true,
				// TeamId:             "",
			},
			6,
		},
		{
			"Include bot accounts, deleted accounts, exclude regular users with no team id",
			account.UserCountOptions{
				// IncludeBotAccounts:  true,
				IncludeDeleted:      true,
				ExcludeRegularUsers: true,
				// TeamId:              "",
			},
			1,
		},
		{
			"Include bot accounts and deleted accounts with existing team id",
			account.UserCountOptions{
				// IncludeBotAccounts: true,
				IncludeDeleted: true,
				// TeamId:             teamId,
			},
			4,
		},
		{
			"Include bot accounts and deleted accounts with fake team id",
			account.UserCountOptions{
				// IncludeBotAccounts: true,
				IncludeDeleted: true,
				// TeamId:             model.NewId(),
			},
			0,
		},
		{
			"Include bot accounts and deleted accounts with existing team id and view restrictions allowing team",
			account.UserCountOptions{
				// IncludeBotAccounts: true,
				IncludeDeleted: true,
				// TeamId:             teamId,
				// ViewRestrictions:   &model.ViewUsersRestrictions{Teams: []string{teamId}},
			},
			4,
		},
		{
			"Include bot accounts and deleted accounts with existing team id and view restrictions not allowing current team",
			account.UserCountOptions{
				// IncludeBotAccounts: true,
				IncludeDeleted: true,
				// TeamId:             teamId,
				// ViewRestrictions:   &model.ViewUsersRestrictions{Teams: []string{model.NewId()}},
			},
			0,
		},
		{
			"Filter by system admins only",
			account.UserCountOptions{
				// TeamId: teamId,
				Roles: []string{model.SystemAdminRoleId},
			},
			1,
		},
		{
			"Filter by system users only",
			account.UserCountOptions{
				// TeamId: teamId,
				Roles: []string{model.SystemUserRoleId},
			},
			2,
		},
		{
			"Filter by system guests only",
			account.UserCountOptions{
				// TeamId: teamId,
				Roles: []string{model.SystemGuestRoleId},
			},
			1,
		},
		{
			"Filter by system admins and system users",
			account.UserCountOptions{
				// TeamId: teamId,
				Roles: []string{model.SystemAdminRoleId, model.SystemUserRoleId},
			},
			3,
		},
		{
			"Filter by system admins, system user and system guests",
			account.UserCountOptions{
				// TeamId: teamId,
				Roles: []string{model.SystemAdminRoleId, model.SystemUserRoleId, model.SystemGuestRoleId},
			},
			4,
		},
		// {
		// 	"Filter by team admins",
		// 	account.UserCountOptions{
		// 		TeamId:    teamId,
		// 		TeamRoles: []string{model.TeamAdminRoleId},
		// 	},
		// 	1,
		// },
		// {
		// 	"Filter by team members",
		// 	account.UserCountOptions{
		// 		TeamId:    teamId,
		// 		TeamRoles: []string{model.TeamUserRoleId},
		// 	},
		// 	1,
		// },
		// {
		// 	"Filter by team guests",
		// 	model.UserCountOptions{
		// 		TeamId:    teamId,
		// 		TeamRoles: []string{model.TeamGuestRoleId},
		// 	},
		// 	1,
		// },
		// {
		// 	"Filter by team guests and any system role",
		// 	model.UserCountOptions{
		// 		TeamId:    teamId,
		// 		TeamRoles: []string{model.TeamGuestRoleId},
		// 		Roles:     []string{model.SystemAdminRoleId},
		// 	},
		// 	2,
		// },
		// {
		// 	"Filter by channel members",
		// 	model.UserCountOptions{
		// 		ChannelId:    channelId,
		// 		ChannelRoles: []string{model.ChannelUserRoleId},
		// 	},
		// 	1,
		// },
		// {
		// 	"Filter by channel members and system admins",
		// 	model.UserCountOptions{
		// 		ChannelId:    channelId,
		// 		Roles:        []string{model.SystemAdminRoleId},
		// 		ChannelRoles: []string{model.ChannelUserRoleId},
		// 	},
		// 	2,
		// },
		// {
		// 	"Filter by channel members and system admins and channel admins",
		// 	model.UserCountOptions{
		// 		ChannelId:    channelId,
		// 		Roles:        []string{model.SystemAdminRoleId},
		// 		ChannelRoles: []string{model.ChannelUserRoleId, model.ChannelAdminRoleId},
		// 	},
		// 	3,
		// },
		// {
		// 	"Filter by channel guests",
		// 	model.UserCountOptions{
		// 		ChannelId:    channelId,
		// 		ChannelRoles: []string{model.ChannelGuestRoleId},
		// 	},
		// 	1,
		// },
		// {
		// 	"Filter by channel guests and any system role",
		// 	model.UserCountOptions{
		// 		ChannelId:    channelId,
		// 		ChannelRoles: []string{model.ChannelGuestRoleId},
		// 		Roles:        []string{model.SystemAdminRoleId},
		// 	},
		// 	2,
		// },
	}
	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			count, err := ss.User().Count(testCase.Options)
			require.NoError(t, err)
			require.Equal(t, testCase.Expected, count)
		})
	}
}
