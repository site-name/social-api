package compliance

import "github.com/sitename/sitename/model"

type MessageExport struct {
	TeamId          *string
	TeamName        *string
	TeamDisplayName *string

	ChannelId          *string
	ChannelName        *string
	ChannelDisplayName *string
	ChannelType        *string

	UserId    *string
	UserEmail *string
	Username  *string
	IsBot     bool

	PostId         *string
	PostCreateAt   *int64
	PostUpdateAt   *int64
	PostDeleteAt   *int64
	PostMessage    *string
	PostType       *string
	PostRootId     *string
	PostProps      *string
	PostOriginalId *string
	PostFileIds    model.StringArray
}

type MessageExportCursor struct {
	LastPostUpdateAt int64
	LastPostId       string
}
