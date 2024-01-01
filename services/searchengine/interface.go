package searchengine

import (
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
)

type SearchEngineInterface interface {
	Start() *model_helper.AppError
	Stop() *model_helper.AppError
	GetFullVersion() string
	GetVersion() int
	GetPlugins() []string
	UpdateConfig(cfg *model_helper.Config)
	GetName() string
	IsActive() bool
	IsIndexingEnabled() bool
	IsSearchEnabled() bool
	IsAutocompletionEnabled() bool
	IsIndexingSync() bool
	// DeleteUserPosts(userID string) *model_helper.AppError
	IndexUser(user *model.User, teamsIds, channelsIds []string) *model_helper.AppError
	// DeleteFile(fileID string) *model_helper.AppError
	DeleteUser(user *model.User) *model_helper.AppError
	// DeletePostFiles(postID string) *model_helper.AppError
	// DeleteUserFiles(userID string) *model_helper.AppError
	// DeleteFilesBatch(endTime, limit int64) *model_helper.AppError
	TestConfig(cfg *model_helper.Config) *model_helper.AppError
	PurgeIndexes() *model_helper.AppError
	RefreshIndexes() *model_helper.AppError
	DataRetentionDeleteIndexes(cutoff time.Time) *model_helper.AppError

	// SearchUsersInChannel(teamId, channelId string, restrictedToChannels []string, term string, options *model_helper.UserSearchOptions) ([]string, []string, *model_helper.AppError)
	// SearchUsersInTeam(teamId string, restrictedToChannels []string, term string, options *model_helper.UserSearchOptions) ([]string, *model_helper.AppError)
	// IndexChannel(channel *model_helper.Channel) *model_helper.AppError
	// SearchChannels(teamId, term string) ([]string, *model_helper.AppError)
	// DeleteChannel(channel *model_helper.Channel) *model_helper.AppError
	// IndexPost(post *model_helper.Post, teamId string) *model_helper.AppError
	// SearchPosts(channels *model_helper.ChannelList, searchParams []*model_helper.SearchParams, page, perPage int) ([]string, model_helper.PostSearchMatches, *model_helper.AppError)
	// DeletePost(post *model_helper.Post) *model_helper.AppError
	// DeleteChannelPosts(channelID string) *model_helper.AppError
	// IndexFile(file *model_helper.FileInfo, channelId string) *model_helper.AppError
	// SearchFiles(channels *model_helper.ChannelList, searchParams []*model_helper.SearchParams, page, perPage int) ([]string, *model_helper.AppError)

}
