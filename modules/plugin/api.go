package plugin

import (
	"io"
	"net/http"

	plugin "github.com/hashicorp/go-plugin"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
)

type PluginIface interface {
	SetAPI(api API)
	SetDriver(driver Driver)
}

// The API can be used to retrieve data or perform actions on behalf of the plugin. Most methods
// have direct counterparts in the REST API and very similar behavior.
//
// Plugins obtain access to the API by embedding SitenamePlugin and accessing the API member
// directly.
type API interface {
	// plugins
	// payment.PaymentInterface
	// LoadPluginConfiguration loads the plugin's configuration. dest should be a pointer to a
	// struct that the configuration JSON can be unmarshalled to.
	//
	// @tag Plugin
	// Minimum server version: 5.2
	LoadPluginConfiguration(dest interface{}) error

	// RegisterCommand registers a custom slash command. When the command is triggered, your plugin
	// can fulfill it via the ExecuteCommand hook.
	//
	// @tag Command
	// Minimum server version: 5.2
	// RegisterCommand(command *model.Command) error

	// UnregisterCommand unregisters a command previously registered via RegisterCommand.
	//
	// @tag Command
	// Minimum server version: 5.2
	// UnregisterCommand(teamID, trigger string) error

	// ExecuteSlashCommand executes a slash command with the given parameters.
	//
	// @tag Command
	// Minimum server version: 5.26
	// ExecuteSlashCommand(commandArgs *model.CommandArgs) (*model.CommandResponse, error)

	// GetSession returns the session object for the Session ID
	//
	// Minimum server version: 5.2
	GetSession(sessionID string) (*model.Session, *model_helper.AppError)

	// GetConfig fetches the currently persisted config
	//
	// @tag Configuration
	// Minimum server version: 5.2
	GetConfig() *model_helper.Config

	// GetUnsanitizedConfig fetches the currently persisted config without removing secrets.
	//
	// @tag Configuration
	// Minimum server version: 5.16
	GetUnsanitizedConfig() *model_helper.Config

	// SaveConfig sets the given config and persists the changes
	//
	// @tag Configuration
	// Minimum server version: 5.2
	SaveConfig(config *model_helper.Config) *model_helper.AppError

	// GetPluginConfig fetches the currently persisted config of plugin
	//
	// @tag Plugin
	// Minimum server version: 5.6
	GetPluginConfig() map[string]interface{}

	// SavePluginConfig sets the given config for plugin and persists the changes
	//
	// @tag Plugin
	// Minimum server version: 5.6
	SavePluginConfig(config map[string]interface{}) *model_helper.AppError

	// GetBundlePath returns the absolute path where the plugin's bundle was unpacked.
	//
	// @tag Plugin
	// Minimum server version: 5.10
	GetBundlePath() (string, error)

	// GetLicense returns the current license used by the Mattermost server. Returns nil if
	// the server does not have a license.
	//
	// @tag Server
	// Minimum server version: 5.10
	// GetLicense() *model.License

	// GetServerVersion return the current Mattermost server version
	//
	// @tag Server
	// Minimum server version: 5.4
	GetServerVersion() string

	// GetSystemInstallDate returns the time that Mattermost was first installed and ran.
	//
	// @tag Server
	// Minimum server version: 5.10
	GetSystemInstallDate() (int64, *model_helper.AppError)

	// GetDiagnosticId returns a unique identifier used by the server for diagnostic reports.
	//
	// @tag Server
	// Minimum server version: 5.10
	// GetDiagnosticId() string

	// GetTelemetryId returns a unique identifier used by the server for telemetry reports.
	//
	// @tag Server
	// Minimum server version: 5.28
	// GetTelemetryId() string

	// CreateUser creates a user.
	//
	// @tag User
	// Minimum server version: 5.2
	CreateUser(user *model.User) (*model.User, *model_helper.AppError)

	// DeleteUser deletes a user.
	//
	// @tag User
	// Minimum server version: 5.2
	DeleteUser(userID string) *model_helper.AppError

	// GetUsers a list of users based on search options.
	//
	// @tag User
	// Minimum server version: 5.10
	GetUsers(options *model_helper.UserGetOptions) ([]*model.User, *model_helper.AppError)

	// GetUser gets a user.
	//
	// @tag User
	// Minimum server version: 5.2
	GetUser(userID string) (*model.User, *model_helper.AppError)

	// GetUserByEmail gets a user by their email address.
	//
	// @tag User
	// Minimum server version: 5.2
	GetUserByEmail(email string) (*model.User, *model_helper.AppError)

	// GetUserByUsername gets a user by their username.
	//
	// @tag User
	// Minimum server version: 5.2
	GetUserByUsername(name string) (*model.User, *model_helper.AppError)

	// GetUsersByUsernames gets users by their usernames.
	//
	// @tag User
	// Minimum server version: 5.6
	GetUsersByUsernames(usernames []string) ([]*model.User, *model_helper.AppError)

	// GetUsersInTeam gets users in team.
	//
	// @tag User
	// @tag Team
	// Minimum server version: 5.6
	// GetUsersInTeam(teamID string, page int, perPage int) ([]*model.User, *model_helper.AppError)

	// GetPreferencesForUser gets a user's preferences.
	//
	// @tag User
	// @tag Preference
	// Minimum server version: 5.26
	GetPreferencesForUser(userID string) ([]model.Preference, *model_helper.AppError)

	// UpdatePreferencesForUser updates a user's preferences.
	//
	// @tag User
	// @tag Preference
	// Minimum server version: 5.26
	UpdatePreferencesForUser(userID string, preferences []model.Preference) *model_helper.AppError

	// DeletePreferencesForUser deletes a user's preferences.
	//
	// @tag User
	// @tag Preference
	// Minimum server version: 5.26
	DeletePreferencesForUser(userID string, preferences []model.Preference) *model_helper.AppError

	// CreateUserAccessToken creates a new access token.
	// @tag User
	CreateUserAccessToken(token *model.UserAccessToken) (*model.UserAccessToken, *model_helper.AppError)

	// RevokeUserAccessToken revokes an existing access token.
	// @tag User
	RevokeUserAccessToken(tokenID string) *model_helper.AppError

	// CreateOAuthApp creates a new OAuth App.
	//
	// @tag OAuth
	// Minimum server version: 5.38
	// CreateOAuthApp(app *model.OAuthApp) (*model.OAuthApp, *model_helper.AppError)

	// // GetOAuthApp gets an existing OAuth App by id.
	// //
	// // @tag OAuth
	// // Minimum server version: 5.38
	// GetOAuthApp(appID string) (*model.OAuthApp, *model_helper.AppError)

	// // UpdateOAuthApp updates an existing OAuth App.
	// //
	// // @tag OAuth
	// // Minimum server version: 5.38
	// UpdateOAuthApp(app *model.OAuthApp) (*model.OAuthApp, *model_helper.AppError)

	// // DeleteOAuthApp deletes an existing OAuth App by id.
	// //
	// // @tag OAuth
	// // Minimum server version: 5.38
	// DeleteOAuthApp(appID string) *model_helper.AppError

	// GetTeamIcon gets the team icon.
	//
	// @tag Team
	// Minimum server version: 5.6
	// GetTeamIcon(teamID string) ([]byte, *model_helper.AppError)

	// SetTeamIcon sets the team icon.
	//
	// @tag Team
	// Minimum server version: 5.6
	// SetTeamIcon(teamID string, data []byte) *model_helper.AppError

	// RemoveTeamIcon removes the team icon.
	//
	// @tag Team
	// Minimum server version: 5.6
	// RemoveTeamIcon(teamID string) *model_helper.AppError

	// UpdateUser updates a user.
	//
	// @tag User
	// Minimum server version: 5.2
	UpdateUser(user *model.User) (*model.User, *model_helper.AppError)

	// GetUserStatus will get a user's status.
	//
	// @tag User
	// Minimum server version: 5.2
	GetUserStatus(userID string) (*model.Status, *model_helper.AppError)

	// GetUserStatusesByIds will return a list of user statuses based on the provided slice of user IDs.
	//
	// @tag User
	// Minimum server version: 5.2
	GetUserStatusesByIds(userIds []string) ([]*model.Status, *model_helper.AppError)

	// UpdateUserStatus will set a user's status until the user, or another integration/plugin, sets it back to online.
	// The status parameter can be: "online", "away", "dnd", or "offline".
	//
	// @tag User
	// Minimum server version: 5.2
	UpdateUserStatus(userID, status string) (*model.Status, *model_helper.AppError)

	// SetUserStatusTimedDND will set a user's status to dnd for given time until the user,
	// or another integration/plugin, sets it back to online.
	// @tag User
	// Minimum server version: 5.35
	// SetUserStatusTimedDND(userId string, endtime int64) (*model.Status, *model_helper.AppError)

	// UpdateUserActive deactivates or reactivates an user.
	//
	// @tag User
	// Minimum server version: 5.8
	UpdateUserActive(userID string, active bool) *model_helper.AppError

	// GetUsersInChannel returns a page of users in a channel. Page counting starts at 0.
	// The sortBy parameter can be: "username" or "status".
	//
	// @tag User
	// @tag Channel
	// Minimum server version: 5.6
	// GetUsersInChannel(channelID, sortBy string, page, perPage int) ([]*model.User, *model_helper.AppError)

	// GetLDAPUserAttributes will return LDAP attributes for a user.
	// The attributes parameter should be a list of attributes to pull.
	// Returns a map with attribute names as keys and the user's attributes as values.
	// Requires an enterprise license, LDAP to be configured and for the user to use LDAP as an authentication method.
	//
	// @tag User
	// Minimum server version: 5.3
	GetLDAPUserAttributes(userID string, attributes []string) (map[string]string, *model_helper.AppError)

	// CreateTeam creates a team.
	//
	// @tag Team
	// Minimum server version: 5.2
	// CreateTeam(team *model.Team) (*model.Team, *model_helper.AppError)

	// DeleteTeam deletes a team.
	//
	// @tag Team
	// Minimum server version: 5.2
	// DeleteTeam(teamID string) *model_helper.AppError

	// GetTeam gets all teams.
	//
	// @tag Team
	// Minimum server version: 5.2
	// GetTeams() ([]*model.Team, *model_helper.AppError)

	// GetTeam gets a team.
	//
	// @tag Team
	// Minimum server version: 5.2
	// GetTeam(teamID string) (*model.Team, *model_helper.AppError)

	// GetTeamByName gets a team by its name.
	//
	// @tag Team
	// Minimum server version: 5.2
	// GetTeamByName(name string) (*model.Team, *model_helper.AppError)

	// GetTeamsUnreadForUser gets the unread message and mention counts for each team to which the given user belongs.
	//
	// @tag Team
	// @tag User
	// Minimum server version: 5.6
	// GetTeamsUnreadForUser(userID string) ([]*model.TeamUnread, *model_helper.AppError)

	// UpdateTeam updates a team.
	//
	// @tag Team
	// Minimum server version: 5.2
	// UpdateTeam(team *model.Team) (*model.Team, *model_helper.AppError)

	// SearchTeams search a team.
	//
	// @tag Team
	// Minimum server version: 5.8
	// SearchTeams(term string) ([]*model.Team, *model_helper.AppError)

	// GetTeamsForUser returns list of teams of given user ID.
	//
	// @tag Team
	// @tag User
	// Minimum server version: 5.6
	// GetTeamsForUser(userID string) ([]*model.Team, *model_helper.AppError)

	// CreateTeamMember creates a team membership.
	//
	// @tag Team
	// @tag User
	// Minimum server version: 5.2
	// CreateTeamMember(teamID, userID string) (*model.TeamMember, *model_helper.AppError)

	// CreateTeamMembers creates a team membership for all provided user ids.
	//
	// @tag Team
	// @tag User
	// Minimum server version: 5.2
	// CreateTeamMembers(teamID string, userIds []string, requestorId string) ([]*model.TeamMember, *model_helper.AppError)

	// CreateTeamMembersGracefully creates a team membership for all provided user ids and reports the users that were not added.
	//
	// @tag Team
	// @tag User
	// Minimum server version: 5.20
	// CreateTeamMembersGracefully(teamID string, userIds []string, requestorId string) ([]*model.TeamMemberWithError, *model_helper.AppError)

	// DeleteTeamMember deletes a team membership.
	//
	// @tag Team
	// @tag User
	// Minimum server version: 5.2
	// DeleteTeamMember(teamID, userID, requestorId string) *model_helper.AppError

	// GetTeamMembers returns the memberships of a specific team.
	//
	// @tag Team
	// @tag User
	// Minimum server version: 5.2
	// GetTeamMembers(teamID string, page, perPage int) ([]*model.TeamMember, *model_helper.AppError)

	// GetTeamMember returns a specific membership.
	//
	// @tag Team
	// @tag User
	// Minimum server version: 5.2
	// GetTeamMember(teamID, userID string) (*model.TeamMember, *model_helper.AppError)

	// GetTeamMembersForUser returns all team memberships for a user.
	//
	// @tag Team
	// @tag User
	// Minimum server version: 5.10
	// GetTeamMembersForUser(userID string, page int, perPage int) ([]*model.TeamMember, *model_helper.AppError)

	// UpdateTeamMemberRoles updates the role for a team membership.
	//
	// @tag Team
	// @tag User
	// Minimum server version: 5.2
	// UpdateTeamMemberRoles(teamID, userID, newRoles string) (*model.TeamMember, *model_helper.AppError)

	// CreateChannel creates a channel.
	//
	// @tag Channel
	// Minimum server version: 5.2
	// CreateChannel(channel *model.Channel) (*model.Channel, *model_helper.AppError)

	// DeleteChannel deletes a channel.
	//
	// @tag Channel
	// Minimum server version: 5.2
	// DeleteChannel(channelId string) *model_helper.AppError

	// GetPublicChannelsForTeam gets a list of all channels.
	//
	// @tag Channel
	// @tag Team
	// Minimum server version: 5.2
	// GetPublicChannelsForTeam(teamID string, page, perPage int) ([]*model.Channel, *model_helper.AppError)

	// GetChannel gets a channel.
	//
	// @tag Channel
	// Minimum server version: 5.2
	// GetChannel(channelId string) (*model.Channel, *model_helper.AppError)

	// GetChannelByName gets a channel by its name, given a team id.
	//
	// @tag Channel
	// Minimum server version: 5.2
	// GetChannelByName(teamID, name string, includeDeleted bool) (*model.Channel, *model_helper.AppError)

	// GetChannelByNameForTeamName gets a channel by its name, given a team name.
	//
	// @tag Channel
	// @tag Team
	// Minimum server version: 5.2
	// GetChannelByNameForTeamName(teamName, channelName string, includeDeleted bool) (*model.Channel, *model_helper.AppError)

	// GetChannelsForTeamForUser gets a list of channels for given user ID in given team ID.
	//
	// @tag Channel
	// @tag Team
	// @tag User
	// Minimum server version: 5.6
	// GetChannelsForTeamForUser(teamID, userID string, includeDeleted bool) ([]*model.Channel, *model_helper.AppError)

	// GetChannelStats gets statistics for a channel.
	//
	// @tag Channel
	// Minimum server version: 5.6
	// GetChannelStats(channelId string) (*model.ChannelStats, *model_helper.AppError)

	// GetDirectChannel gets a direct message channel.
	// If the channel does not exist it will create it.
	//
	// @tag Channel
	// @tag User
	// Minimum server version: 5.2
	// GetDirectChannel(userId1, userId2 string) (*model.Channel, *model_helper.AppError)

	// GetGroupChannel gets a group message channel.
	// If the channel does not exist it will create it.
	//
	// @tag Channel
	// @tag User
	// Minimum server version: 5.2
	// GetGroupChannel(userIds []string) (*model.Channel, *model_helper.AppError)

	// UpdateChannel updates a channel.
	//
	// @tag Channel
	// Minimum server version: 5.2
	// UpdateChannel(channel *model.Channel) (*model.Channel, *model_helper.AppError)

	// SearchChannels returns the channels on a team matching the provided search term.
	//
	// @tag Channel
	// Minimum server version: 5.6
	// SearchChannels(teamID string, term string) ([]*model.Channel, *model_helper.AppError)

	// CreateChannelSidebarCategory creates a new sidebar category for a set of channels.
	//
	// @tag ChannelSidebar
	// Minimum server version: 5.37
	// CreateChannelSidebarCategory(userID, teamID string, newCategory *model.SidebarCategoryWithChannels) (*model.SidebarCategoryWithChannels, *model_helper.AppError)

	// GetChannelSidebarCategories returns sidebar categories.
	//
	// @tag ChannelSidebar
	// Minimum server version: 5.37
	// GetChannelSidebarCategories(userID, teamID string) (*model.OrderedSidebarCategories, *model_helper.AppError)

	// UpdateChannelSidebarCategories updates the channel sidebar categories.
	//
	// @tag ChannelSidebar
	// Minimum server version: 5.37
	// UpdateChannelSidebarCategories(userID, teamID string, categories []*model.SidebarCategoryWithChannels) ([]*model.SidebarCategoryWithChannels, *model_helper.AppError)

	// SearchUsers returns a list of users based on some search criteria.
	//
	// @tag User
	// Minimum server version: 5.6
	SearchUsers(search *model_helper.UserSearch) ([]*model.User, *model_helper.AppError)

	// SearchPostsInTeam returns a list of posts in a specific team that match the given params.
	//
	// @tag Post
	// @tag Team
	// Minimum server version: 5.10
	// SearchPostsInTeam(teamID string, paramsList []*model.SearchParams) ([]*model.Post, *model_helper.AppError)

	// SearchPostsInTeamForUser returns a list of posts by team and user that match the given
	// search parameters.
	// @tag Post
	// Minimum server version: 5.26
	// SearchPostsInTeamForUser(teamID string, userID string, searchParams model.SearchParameter) (*model.PostSearchResults, *model_helper.AppError)

	// AddChannelMember joins a user to a channel (as if they joined themselves)
	// This means the user will not receive notifications for joining the channel.
	//
	// @tag Channel
	// @tag User
	// Minimum server version: 5.2
	// AddChannelMember(channelId, userID string) (*model.ChannelMember, *model_helper.AppError)

	// AddUserToChannel adds a user to a channel as if the specified user had invited them.
	// This means the user will receive the regular notifications for being added to the channel.
	//
	// @tag User
	// @tag Channel
	// Minimum server version: 5.18
	// AddUserToChannel(channelId, userID, asUserId string) (*model.ChannelMember, *model_helper.AppError)

	// GetChannelMember gets a channel membership for a user.
	//
	// @tag Channel
	// @tag User
	// Minimum server version: 5.2
	// GetChannelMember(channelId, userID string) (*model.ChannelMember, *model_helper.AppError)

	// GetChannelMembers gets a channel membership for all users.
	//
	// @tag Channel
	// @tag User
	// Minimum server version: 5.6
	// GetChannelMembers(channelId string, page, perPage int) (*model.ChannelMembers, *model_helper.AppError)

	// GetChannelMembersByIds gets a channel membership for a particular User
	//
	// @tag Channel
	// @tag User
	// Minimum server version: 5.6
	// GetChannelMembersByIds(channelId string, userIds []string) (*model.ChannelMembers, *model_helper.AppError)

	// GetChannelMembersForUser returns all channel memberships on a team for a user.
	//
	// @tag Channel
	// @tag User
	// Minimum server version: 5.10
	// GetChannelMembersForUser(teamID, userID string, page, perPage int) ([]*model.ChannelMember, *model_helper.AppError)

	// UpdateChannelMemberRoles updates a user's roles for a channel.
	//
	// @tag Channel
	// @tag User
	// Minimum server version: 5.2
	// UpdateChannelMemberRoles(channelId, userID, newRoles string) (*model.ChannelMember, *model_helper.AppError)

	// UpdateChannelMemberNotifications updates a user's notification properties for a channel.
	//
	// @tag Channel
	// @tag User
	// Minimum server version: 5.2
	// UpdateChannelMemberNotifications(channelId, userID string, notifications map[string]string) (*model.ChannelMember, *model_helper.AppError)

	// GetGroup gets a group by ID.
	//
	// @tag Group
	// Minimum server version: 5.18
	// GetGroup(groupId string) (*model.Group, *model_helper.AppError)

	// GetGroupByName gets a group by name.
	//
	// @tag Group
	// Minimum server version: 5.18
	// GetGroupByName(name string) (*model.Group, *model_helper.AppError)

	// GetGroupMemberUsers gets a page of users belonging to the given group.
	//
	// @tag Group
	// Minimum server version: 5.35
	// GetGroupMemberUsers(groupID string, page, perPage int) ([]*model.User, *model_helper.AppError)

	// GetGroupsBySource gets a list of all groups for the given source.
	//
	// @tag Group
	// Minimum server version: 5.35
	// GetGroupsBySource(groupSource model.GroupSource) ([]*model.Group, *model_helper.AppError)

	// GetGroupsForUser gets the groups a user is in.
	//
	// @tag Group
	// @tag User
	// Minimum server version: 5.18
	// GetGroupsForUser(userID string) ([]*model.Group, *model_helper.AppError)

	// DeleteChannelMember deletes a channel membership for a user.
	//
	// @tag Channel
	// @tag User
	// Minimum server version: 5.2
	// DeleteChannelMember(channelId, userID string) *model_helper.AppError

	// CreatePost creates a post.
	//
	// @tag Post
	// Minimum server version: 5.2
	// CreatePost(post *model.Post) (*model.Post, *model_helper.AppError)

	// AddReaction add a reaction to a post.
	//
	// @tag Post
	// Minimum server version: 5.3
	// AddReaction(reaction *model.Reaction) (*model.Reaction, *model_helper.AppError)

	// RemoveReaction remove a reaction from a post.
	//
	// @tag Post
	// Minimum server version: 5.3
	// RemoveReaction(reaction *model.Reaction) *model_helper.AppError

	// GetReaction get the reactions of a post.
	//
	// @tag Post
	// Minimum server version: 5.3
	// GetReactions(postId string) ([]*model.Reaction, *model_helper.AppError)

	// SendEphemeralPost creates an ephemeral post.
	//
	// @tag Post
	// Minimum server version: 5.2
	// SendEphemeralPost(userID string, post *model.Post) *model.Post

	// UpdateEphemeralPost updates an ephemeral message previously sent to the user.
	// EXPERIMENTAL: This API is experimental and can be changed without advance notice.
	//
	// @tag Post
	// Minimum server version: 5.2
	// UpdateEphemeralPost(userID string, post *model.Post) *model.Post

	// DeleteEphemeralPost deletes an ephemeral message previously sent to the user.
	// EXPERIMENTAL: This API is experimental and can be changed without advance notice.
	//
	// @tag Post
	// Minimum server version: 5.2
	// DeleteEphemeralPost(userID, postId string)

	// DeletePost deletes a post.
	//
	// @tag Post
	// Minimum server version: 5.2
	// DeletePost(postId string) *model_helper.AppError

	// GetPostThread gets a post with all the other posts in the same thread.
	//
	// @tag Post
	// Minimum server version: 5.6
	// GetPostThread(postId string) (*model.PostList, *model_helper.AppError)

	// GetPost gets a post.
	//
	// @tag Post
	// Minimum server version: 5.2
	// GetPost(postId string) (*model.Post, *model_helper.AppError)

	// GetPostsSince gets posts created after a specified time as Unix time in milliseconds.
	//
	// @tag Post
	// @tag Channel
	// Minimum server version: 5.6
	// GetPostsSince(channelId string, time int64) (*model.PostList, *model_helper.AppError)

	// GetPostsAfter gets a page of posts that were posted after the post provided.
	//
	// @tag Post
	// @tag Channel
	// Minimum server version: 5.6
	// GetPostsAfter(channelId, postId string, page, perPage int) (*model.PostList, *model_helper.AppError)

	// GetPostsBefore gets a page of posts that were posted before the post provided.
	//
	// @tag Post
	// @tag Channel
	// Minimum server version: 5.6
	// GetPostsBefore(channelId, postId string, page, perPage int) (*model.PostList, *model_helper.AppError)

	// GetPostsForChannel gets a list of posts for a channel.
	//
	// @tag Post
	// @tag Channel
	// Minimum server version: 5.6
	// GetPostsForChannel(channelId string, page, perPage int) (*model.PostList, *model_helper.AppError)

	// GetTeamStats gets a team's statistics
	//
	// @tag Team
	// Minimum server version: 5.8
	// GetTeamStats(teamID string) (*model.TeamStats, *model_helper.AppError)

	// UpdatePost updates a post.
	//
	// @tag Post
	// Minimum server version: 5.2
	// UpdatePost(post *model.Post) (*model.Post, *model_helper.AppError)

	// GetProfileImage gets user's profile image.
	//
	// @tag User
	// Minimum server version: 5.6
	GetProfileImage(userID string) ([]byte, *model_helper.AppError)

	// SetProfileImage sets a user's profile image.
	//
	// @tag User
	// Minimum server version: 5.6
	SetProfileImage(userID string, data []byte) *model_helper.AppError

	// GetEmojiList returns a page of custom emoji on the system.
	//
	// The sortBy parameter can be: "name".
	//
	// @tag Emoji
	// Minimum server version: 5.6
	// GetEmojiList(sortBy string, page, perPage int) ([]*model.Emoji, *model_helper.AppError)

	// GetEmojiByName gets an emoji by it's name.
	//
	// @tag Emoji
	// Minimum server version: 5.6
	// GetEmojiByName(name string) (*model.Emoji, *model_helper.AppError)

	// GetEmoji returns a custom emoji based on the emojiId string.
	//
	// @tag Emoji
	// Minimum server version: 5.6
	// GetEmoji(emojiId string) (*model.Emoji, *model_helper.AppError)

	// CopyFileInfos duplicates the FileInfo objects referenced by the given file ids,
	// recording the given user id as the new creator and returning the new set of file ids.
	//
	// The duplicate FileInfo objects are not initially linked to a post, but may now be passed
	// to CreatePost. Use this API to duplicate a post and its file attachments without
	// actually duplicating the uploaded files.
	//
	// @tag File
	// @tag User
	// Minimum server version: 5.2
	CopyFileInfos(userID string, fileIds []string) ([]string, *model_helper.AppError)

	// GetFileInfo gets a File Info for a specific fileId
	//
	// @tag File
	// Minimum server version: 5.3
	GetFileInfo(fileId string) (*model.FileInfo, *model_helper.AppError)

	// GetFileInfos gets File Infos with options
	//
	// @tag File
	// Minimum server version: 5.22
	GetFileInfos(page, perPage int, opt *model_helper.FileInfoFilterOption) ([]*model.FileInfo, *model_helper.AppError)

	// GetFile gets content of a file by it's ID
	//
	// @tag File
	// Minimum server version: 5.8
	GetFile(fileId string) ([]byte, *model_helper.AppError)

	// GetFileLink gets the public link to a file by fileId.
	//
	// @tag File
	// Minimum server version: 5.6
	GetFileLink(fileId string) (string, *model_helper.AppError)

	// ReadFile reads the file from the backend for a specific path
	//
	// @tag File
	// Minimum server version: 5.3
	ReadFile(path string) ([]byte, *model_helper.AppError)

	// GetEmojiImage returns the emoji image.
	//
	// @tag Emoji
	// Minimum server version: 5.6
	// GetEmojiImage(emojiId string) ([]byte, string, *model_helper.AppError)

	// UploadFile will upload a file to a channel using a multipart request, to be later attached to a post.
	//
	// @tag File
	// @tag Channel
	// Minimum server version: 5.6
	// UploadFile(data []byte, channelId string, filename string) (*model.FileInfo, *model_helper.AppError)

	// OpenInteractiveDialog will open an interactive dialog on a user's client that
	// generated the trigger ID. Used with interactive message buttons, menus
	// and slash commands.
	//
	// Minimum server version: 5.6
	// OpenInteractiveDialog(dialog model.OpenDialogRequest) *model_helper.AppError

	// Plugin Section

	// GetPlugins will return a list of plugin manifests for currently active plugin.
	//
	// @tag Plugin
	// Minimum server version: 5.6
	GetPlugins() ([]*model_helper.Manifest, *model_helper.AppError)

	// EnablePlugin will enable an plugin installed.
	//
	// @tag Plugin
	// Minimum server version: 5.6
	EnablePlugin(id string) *model_helper.AppError

	// DisablePlugin will disable an enabled plugin.
	//
	// @tag Plugin
	// Minimum server version: 5.6
	DisablePlugin(id string) *model_helper.AppError

	// RemovePlugin will disable and delete a plugin.
	//
	// @tag Plugin
	// Minimum server version: 5.6
	RemovePlugin(id string) *model_helper.AppError

	// GetPluginStatus will return the status of a plugin.
	//
	// @tag Plugin
	// Minimum server version: 5.6
	GetPluginStatus(id string) (*model_helper.PluginStatus, *model_helper.AppError)

	// InstallPlugin will upload another plugin with tar.gz model.
	// Previous version will be replaced on replace true.
	//
	// @tag Plugin
	// Minimum server version: 5.18
	InstallPlugin(file io.Reader, replace bool) (*model_helper.Manifest, *model_helper.AppError)

	// KV Store Section

	// KVSet stores a key-value pair, unique per plugin.
	// Provided helper functions and internal plugin code will use the prefix `mmi_` before keys. Do not use this prefix.
	//
	// @tag KeyValueStore
	// Minimum server version: 5.2
	KVSet(key string, value []byte) *model_helper.AppError

	// KVCompareAndSet updates a key-value pair, unique per plugin, but only if the current value matches the given oldValue.
	// Inserts a new key if oldValue == nil.
	// Returns (false, err) if DB error occurred
	// Returns (false, nil) if current value != oldValue or key already exists when inserting
	// Returns (true, nil) if current value == oldValue or new key is inserted
	//
	// @tag KeyValueStore
	// Minimum server version: 5.12
	KVCompareAndSet(key string, oldValue, newValue []byte) (bool, *model_helper.AppError)

	// KVCompareAndDelete deletes a key-value pair, unique per plugin, but only if the current value matches the given oldValue.
	// Returns (false, err) if DB error occurred
	// Returns (false, nil) if current value != oldValue or key does not exist when deleting
	// Returns (true, nil) if current value == oldValue and the key was deleted
	//
	// @tag KeyValueStore
	// Minimum server version: 5.16
	KVCompareAndDelete(key string, oldValue []byte) (bool, *model_helper.AppError)

	// KVSetWithOptions stores a key-value pair, unique per plugin, according to the given options.
	// Returns (false, err) if DB error occurred
	// Returns (false, nil) if the value was not set
	// Returns (true, nil) if the value was set
	//
	// Minimum server version: 5.20
	KVSetWithOptions(key string, value []byte, options model_helper.PluginKVSetOptions) (bool, *model_helper.AppError)

	// KVSet stores a key-value pair with an expiry time, unique per plugin.
	//
	// @tag KeyValueStore
	// Minimum server version: 5.6
	KVSetWithExpiry(key string, value []byte, expireInSeconds int64) *model_helper.AppError

	// KVGet retrieves a value based on the key, unique per plugin. Returns nil for non-existent keys.
	//
	// @tag KeyValueStore
	// Minimum server version: 5.2
	KVGet(key string) ([]byte, *model_helper.AppError)

	// KVDelete removes a key-value pair, unique per plugin. Returns nil for non-existent keys.
	//
	// @tag KeyValueStore
	// Minimum server version: 5.2
	KVDelete(key string) *model_helper.AppError

	// KVDeleteAll removes all key-value pairs for a plugin.
	//
	// @tag KeyValueStore
	// Minimum server version: 5.6
	KVDeleteAll() *model_helper.AppError

	// KVList lists all keys for a plugin.
	//
	// @tag KeyValueStore
	// Minimum server version: 5.6
	KVList(page, perPage int) ([]string, *model_helper.AppError)

	// PublishWebSocketEvent sends an event to WebSocket connections.
	// event is the type and will be prepended with "custom_<pluginid>_".
	// payload is the data sent with the event. Interface values must be primitive Go types or mattermost-server/model types.
	// broadcast determines to which users to send the event.
	//
	// Minimum server version: 5.2
	PublishWebSocketEvent(event string, payload map[string]interface{}, broadcast *model_helper.WebsocketBroadcast)

	// HasPermissionTo check if the user has the permission at system scope.
	//
	// @tag User
	// Minimum server version: 5.3
	HasPermissionTo(userID string, permission *model_helper.Permission) bool

	// HasPermissionToTeam check if the user has the permission at team scope.
	//
	// @tag User
	// @tag Team
	// Minimum server version: 5.3
	// HasPermissionToTeam(userID, teamID string, permission *model.Permission) bool

	// HasPermissionToChannel check if the user has the permission at channel scope.
	//
	// @tag User
	// @tag Channel
	// Minimum server version: 5.3
	// HasPermissionToChannel(userID, channelId string, permission *model.Permission) bool

	// LogDebug writes a log message to the Mattermost server log model.
	// Appropriate context such as the plugin name will already be added as fields so plugins
	// do not need to add that info.
	//
	// @tag Logging
	// Minimum server version: 5.2
	LogDebug(msg string, keyValuePairs ...interface{})

	// LogInfo writes a log message to the Mattermost server log model.
	// Appropriate context such as the plugin name will already be added as fields so plugins
	// do not need to add that info.
	//
	// @tag Logging
	// Minimum server version: 5.2
	LogInfo(msg string, keyValuePairs ...interface{})

	// LogError writes a log message to the Mattermost server log model.
	// Appropriate context such as the plugin name will already be added as fields so plugins
	// do not need to add that info.
	//
	// @tag Logging
	// Minimum server version: 5.2
	LogError(msg string, keyValuePairs ...interface{})

	// LogWarn writes a log message to the Mattermost server log model.
	// Appropriate context such as the plugin name will already be added as fields so plugins
	// do not need to add that info.
	//
	// @tag Logging
	// Minimum server version: 5.2
	LogWarn(msg string, keyValuePairs ...interface{})

	// SendMail sends an email to a specific address
	//
	// Minimum server version: 5.7
	SendMail(to, subject, htmlBody string) *model_helper.AppError

	// CreateBot creates the given bot and corresponding user.
	//
	// @tag Bot
	// Minimum server version: 5.10
	// CreateBot(bot *model.Bot) (*model.Bot, *model_helper.AppError)

	// PatchBot applies the given patch to the bot and corresponding user.
	//
	// @tag Bot
	// Minimum server version: 5.10
	// PatchBot(botUserId string, botPatch *model.BotPatch) (*model.Bot, *model_helper.AppError)

	// GetBot returns the given bot.
	//
	// @tag Bot
	// Minimum server version: 5.10
	// GetBot(botUserId string, includeDeleted bool) (*model.Bot, *model_helper.AppError)

	// GetBots returns the requested page of bots.
	//
	// @tag Bot
	// Minimum server version: 5.10
	// GetBots(options *model.BotGetOptions) ([]*model.Bot, *model_helper.AppError)

	// UpdateBotActive marks a bot as active or inactive, along with its corresponding user.
	//
	// @tag Bot
	// Minimum server version: 5.10
	// UpdateBotActive(botUserId string, active bool) (*model.Bot, *model_helper.AppError)

	// PermanentDeleteBot permanently deletes a bot and its corresponding user.
	//
	// @tag Bot
	// Minimum server version: 5.10
	// PermanentDeleteBot(botUserId string) *model_helper.AppError

	// GetBotIconImage gets LHS bot icon image.
	//
	// @tag Bot
	// Minimum server version: 5.14
	// GetBotIconImage(botUserId string) ([]byte, *model_helper.AppError)

	// SetBotIconImage sets LHS bot icon image.
	// Icon image must be SVG format, all other formats are rejected.
	//
	// @tag Bot
	// Minimum server version: 5.14
	// SetBotIconImage(botUserId string, data []byte) *model_helper.AppError

	// DeleteBotIconImage deletes LHS bot icon image.
	//
	// @tag Bot
	// Minimum server version: 5.14
	// DeleteBotIconImage(botUserId string) *model_helper.AppError

	// PluginHTTP allows inter-plugin requests to plugin APIs.
	//
	// Minimum server version: 5.18
	PluginHTTP(request *http.Request) *http.Response

	// PublishUserTyping publishes a user is typing WebSocket event.
	// The parentId parameter may be an empty string, the other parameters are required.
	//
	// @tag User
	// Minimum server version: 5.26
	// PublishUserTyping(userID, channelId, parentId string) *model_helper.AppError

	// CreateCommand creates a server-owned slash command that is not handled by the plugin
	// itself, and which will persist past the life of the plugin. The command will have its
	// CreatorId set to "" and its PluginId set to the id of the plugin that created it.
	//
	// @tag SlashCommand
	// Minimum server version: 5.28
	// CreateCommand(cmd *model.Command) (*model.Command, error)

	// ListCommands returns the list of all slash commands for teamID. E.g., custom commands
	// (those created through the integrations menu, the REST api, or the plugin api CreateCommand),
	// plugin commands (those created with plugin api RegisterCommand), and builtin commands
	// (those added internally through RegisterCommandProvider).
	//
	// @tag SlashCommand
	// Minimum server version: 5.28
	// ListCommands(teamID string) ([]*model.Command, error)

	// ListCustomCommands returns the list of slash commands for teamID that where created
	// through the integrations menu, the REST api, or the plugin api CreateCommand.
	//
	// @tag SlashCommand
	// Minimum server version: 5.28
	// ListCustomCommands(teamID string) ([]*model.Command, error)

	// ListPluginCommands returns the list of slash commands for teamID that were created
	// with the plugin api RegisterCommand.
	//
	// @tag SlashCommand
	// Minimum server version: 5.28
	// ListPluginCommands(teamID string) ([]*model.Command, error)

	// ListBuiltInCommands returns the list of slash commands that are builtin commands
	// (those added internally through RegisterCommandProvider).
	//
	// @tag SlashCommand
	// Minimum server version: 5.28
	// ListBuiltInCommands() ([]*model.Command, error)

	// GetCommand returns the command definition based on a command id string.
	//
	// @tag SlashCommand
	// Minimum server version: 5.28
	// GetCommand(commandID string) (*model.Command, error)

	// UpdateCommand updates a single command (commandID) with the information provided in the
	// updatedCmd model.Command struct. The following fields in the command cannot be updated:
	// Id, Token, CreateAt, DeleteAt, and PluginId. If updatedCmd.TeamId is blank, it
	// will be set to commandID's TeamId.
	//
	// @tag SlashCommand
	// Minimum server version: 5.28
	// UpdateCommand(commandID string, updatedCmd *model.Command) (*model.Command, error)

	// DeleteCommand deletes a slash command (commandID).
	//
	// @tag SlashCommand
	// Minimum server version: 5.28
	// DeleteCommand(commandID string) error

	// PublishPluginClusterEvent broadcasts a plugin event to all other running instances of
	// the calling plugin that are present in the cluster.
	//
	// This method is used to allow plugin communication in a High-Availability cluster.
	// The receiving side should implement the OnPluginClusterEvent hook
	// to receive events sent through this method.
	//
	// Minimum server version: 5.36
	// PublishPluginClusterEvent(ev model.PluginClusterEvent, opts model.PluginClusterEventSendOptions) error

	// RequestTrialLicense requests a trial license and installs it in the server
	//
	// Minimum server version: 5.36
	// RequestTrialLicense(requesterID string, users int, termsAccepted bool, receiveEmailsAccepted bool) *model_helper.AppError
}

var handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "SITENAME_PLUGIN",
	MagicCookieValue: "Securely message teams, anywhere.",
}
