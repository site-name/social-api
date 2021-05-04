package searchengine

import (
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
)

type SearchEngineInterface interface {
	Start() *model.AppError
	Stop() *model.AppError
	GetFullVersion() string
	GetVersion() int
	GetPlugins() []string
	UpdateConfig(cfg *model.Config)
	GetName() string
	IsActive() bool
	IsIndexingEnabled() bool
	IsSearchEnabled() bool
	IsAutocompletionEnabled() bool
	IsIndexingSync() bool
	// DeleteUserPosts(userID string) *model.AppError
	IndexUser(user *account.User, teamsIds, channelsIds []string) *model.AppError
	// DeleteFile(fileID string) *model.AppError
	DeleteUser(user *account.User) *model.AppError
	// DeletePostFiles(postID string) *model.AppError
	// DeleteUserFiles(userID string) *model.AppError
	// DeleteFilesBatch(endTime, limit int64) *model.AppError
	TestConfig(cfg *model.Config) *model.AppError
	PurgeIndexes() *model.AppError
	RefreshIndexes() *model.AppError
	DataRetentionDeleteIndexes(cutoff time.Time) *model.AppError

	// SearchUsersInChannel(teamId, channelId string, restrictedToChannels []string, term string, options *model.UserSearchOptions) ([]string, []string, *model.AppError)
	// SearchUsersInTeam(teamId string, restrictedToChannels []string, term string, options *model.UserSearchOptions) ([]string, *model.AppError)
	// IndexChannel(channel *model.Channel) *model.AppError
	// SearchChannels(teamId, term string) ([]string, *model.AppError)
	// DeleteChannel(channel *model.Channel) *model.AppError
	// IndexPost(post *model.Post, teamId string) *model.AppError
	// SearchPosts(channels *model.ChannelList, searchParams []*model.SearchParams, page, perPage int) ([]string, model.PostSearchMatches, *model.AppError)
	// DeletePost(post *model.Post) *model.AppError
	// DeleteChannelPosts(channelID string) *model.AppError
	// IndexFile(file *model.FileInfo, channelId string) *model.AppError
	// SearchFiles(channels *model.ChannelList, searchParams []*model.SearchParams, page, perPage int) ([]string, *model.AppError)

}
