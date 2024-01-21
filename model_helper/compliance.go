package model_helper

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util"
)

type CompliancePost struct {

	// From Team
	TeamName        string
	TeamDisplayName string

	// From Channel
	ChannelName        string
	ChannelDisplayName string
	ChannelType        string

	// From User
	UserUsername string
	UserEmail    string
	UserNickname string

	// From Post
	PostId         string
	PostCreateAt   int64
	PostUpdateAt   int64
	PostDeleteAt   int64
	PostRootId     string
	PostParentId   string
	PostOriginalId string
	PostMessage    string
	PostType       string
	PostProps      string
	PostHashtags   string
	PostFileIds    string

	IsBot bool
}

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
	PostFileIds    util.AnyArray[string]
}

type MessageExportCursor struct {
	LastPostUpdateAt int64
	LastPostId       string
}

func CompliancePostHeader() []string {
	return []string{
		"TeamName",
		"TeamDisplayName",

		"ChannelName",
		"ChannelDisplayName",
		"ChannelType",

		"UserUsername",
		"UserEmail",
		"UserNickname",
		"UserType",

		"PostId",
		"PostCreateAt",
		"PostUpdateAt",
		"PostDeleteAt",
		"PostRootId",
		"PostParentId",
		"PostOriginalId",
		"PostMessage",
		"PostType",
		"PostProps",
		"PostHashtags",
		"PostFileIds",
	}
}

func cleanComplianceStrings(in string) string {
	if matched, _ := regexp.MatchString("^\\s*(=|\\+|\\-)", in); matched {
		return "'" + in
	}
	return in
}

func (cp *CompliancePost) Row() []string {

	postDeleteAt := ""
	if cp.PostDeleteAt > 0 {
		postDeleteAt = time.Unix(0, cp.PostDeleteAt*int64(1000*1000)).Format(time.RFC3339Nano)
	}

	postUpdateAt := ""
	if cp.PostUpdateAt != cp.PostCreateAt {
		postUpdateAt = time.Unix(0, cp.PostUpdateAt*int64(1000*1000)).Format(time.RFC3339Nano)
	}

	userType := "user"
	if cp.IsBot {
		userType = "bot"
	}

	return []string{
		cleanComplianceStrings(cp.TeamName),
		cleanComplianceStrings(cp.TeamDisplayName),

		cleanComplianceStrings(cp.ChannelName),
		cleanComplianceStrings(cp.ChannelDisplayName),
		cleanComplianceStrings(cp.ChannelType),

		cleanComplianceStrings(cp.UserUsername),
		cleanComplianceStrings(cp.UserEmail),
		cleanComplianceStrings(cp.UserNickname),
		userType,

		cp.PostId,
		time.Unix(0, cp.PostCreateAt*int64(1000*1000)).Format(time.RFC3339Nano),
		postUpdateAt,
		postDeleteAt,

		cp.PostRootId,
		cp.PostParentId,
		cp.PostOriginalId,
		cleanComplianceStrings(cp.PostMessage),
		cp.PostType,
		cp.PostProps,
		cp.PostHashtags,
		cp.PostFileIds,
	}
}

// ComplianceExportCursor is used for paginated iteration of posts
// for compliance export.
// We need to keep track of the last post ID in addition to the last post
// CreateAt to break ties when two posts have the same CreateAt.
type ComplianceExportCursor struct {
	LastChannelsQueryPostCreateAt       int64
	LastChannelsQueryPostID             string
	ChannelsQueryCompleted              bool
	LastDirectMessagesQueryPostCreateAt int64
	LastDirectMessagesQueryPostID       string
	DirectMessagesQueryCompleted        bool
}

func ComplianceCommonPre(c *model.Compliance) {
	if c.Status.IsValid() != nil {
		c.Status = model.ComplianceStatusCreated
	}
	c.Emails = NormalizeEmail(c.Emails)
	c.Keywords = strings.ToLower(c.Keywords)
}

func ComplianceJobNName(c *model.Compliance) string {
	jobName := c.Type.String()
	if c.Type == model.ComplianceTypeDaily {
		jobName += "-" + c.Desc
	}
	jobName += "-" + c.ID
	return jobName
}

func ComplianceIsValid(c *model.Compliance) *AppError {
	if c.Desc == "" {
		return NewAppError("Compliance.IsValid", "model.compliance.is_valid.desc.app_error", nil, "please provide valid desc", http.StatusBadRequest)
	}
	if c.StartAt == 0 {
		return NewAppError("Compliance.IsValid", "model.compliance.is_valid.start_at.app_error", nil, "please provide valid start at", http.StatusBadRequest)
	}
	if c.EndAt == 0 || c.EndAt < c.StartAt {
		return NewAppError("Compliance.IsValid", "model.compliance.is_valid.end_at.app_error", nil, "please provide valid end at", http.StatusBadRequest)
	}

	return nil
}
