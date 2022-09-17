package compliance

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/store"
)

type SqlComplianceStore struct {
	store.Store
}

func (s *SqlComplianceStore) ModelFields(prefix string) model.AnyArray[string] {
	res := model.AnyArray[string]{
		"Id",
		"CreateAt",
		"UserId",
		"Status",
		"Count",
		"Desc",
		"Type",
		"StartAt",
		"EndAt",
		"Keywords",
		"Emails",
	}
	if prefix == "" {
		return res
	}

	return res.Map(func(_ int, s string) string {
		return prefix + s
	})
}

func NewSqlComplianceStore(s store.Store) store.ComplianceStore {
	return &SqlComplianceStore{s}
}

func (s *SqlComplianceStore) Save(compliance *model.Compliance) (*model.Compliance, error) {
	compliance.PreSave()
	if err := compliance.IsValid(); err != nil {
		return nil, err
	}

	query := "INSERT INTO " + store.ComplianceTableName + " (" + s.ModelFields("").Join(",") + ") VALUES (" + s.ModelFields(":").Join(",") + ")"

	if _, err := s.GetMasterX().NamedExec(query, compliance); err != nil {
		return nil, errors.Wrap(err, "failed to save Compliance")
	}
	return compliance, nil
}

func (s *SqlComplianceStore) Update(compliance *model.Compliance) (*model.Compliance, error) {
	compliance.PreSave()
	if err := compliance.IsValid(); err != nil {
		return nil, err
	}

	query := "UPDATE " + store.ComplianceTableName + " SET " + s.
		ModelFields("").
		Map(func(_ int, s string) string {
			return s + "=:" + s
		}).
		Join(",") + " WHERE Id=:Id"

	if _, err := s.GetMasterX().NamedExec(query, compliance); err != nil {
		return nil, errors.Wrap(err, "failed to update Compliance")
	}
	return compliance, nil
}

func (s *SqlComplianceStore) GetAll(offset, limit int) (model.Compliances, error) {
	query := "SELECT * FROM Compliances ORDER BY CreateAt DESC LIMIT ? OFFSET ?"

	var compliances model.Compliances
	if err := s.GetReplicaX().Select(&compliances, query, limit, offset); err != nil {
		return nil, errors.Wrap(err, "failed to find all Compliances")
	}
	return compliances, nil
}

func (s *SqlComplianceStore) Get(id string) (*model.Compliance, error) {
	var res model.Compliance

	err := s.GetReplicaX().Get(&res, "SELECT * FROM "+store.ComplianceTableName+" WHERE Id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(store.ComplianceTableName, id)
		}
		return nil, errors.Wrapf(err, "failed to get Compliance with id=%s", id)
	}
	return &res, nil
}

func (s *SqlComplianceStore) ComplianceExport(job *model.Compliance, cursor model.ComplianceExportCursor, limit int) ([]*model.CompliancePost, model.ComplianceExportCursor, error) {
	keywordQuery := ""
	var argsKeywords []interface{}
	keywords := strings.Fields(strings.TrimSpace(strings.ToLower(strings.Replace(job.Keywords, ",", " ", -1))))
	if len(keywords) > 0 {
		clauses := make([]string, len(keywords))

		for i, keyword := range keywords {
			keyword = store.SanitizeSearchTerm(keyword, "\\")
			clauses[i] = "LOWER(Posts.Message) LIKE ?"
			argsKeywords = append(argsKeywords, "%"+keyword+"%")
		}

		keywordQuery = "AND (" + strings.Join(clauses, " OR ") + ")"
	}

	emailQuery := ""
	var argsEmails []interface{}
	emails := strings.Fields(strings.TrimSpace(strings.ToLower(strings.Replace(job.Emails, ",", " ", -1))))
	if len(emails) > 0 {
		clauses := make([]string, len(emails))

		for i, email := range emails {
			clauses[i] = "Users.Email = ?"
			argsEmails = append(argsEmails, email)
		}

		emailQuery = "AND (" + strings.Join(clauses, " OR ") + ")"
	}

	// The idea is to first iterate over the channel posts, and then when we run out of those,
	// start iterating over the direct message posts.

	channelPosts := []*model.CompliancePost{}
	channelsQuery := ""
	var argsChannelsQuery []interface{}
	if !cursor.ChannelsQueryCompleted {
		if cursor.LastChannelsQueryPostCreateAt == 0 {
			cursor.LastChannelsQueryPostCreateAt = job.StartAt
		}
		// append the named parameters of SQL query in the correct order to argsChannelsQuery
		argsChannelsQuery = append(argsChannelsQuery, cursor.LastChannelsQueryPostCreateAt, cursor.LastChannelsQueryPostCreateAt, cursor.LastChannelsQueryPostID, job.EndAt)
		argsChannelsQuery = append(argsChannelsQuery, argsEmails...)
		argsChannelsQuery = append(argsChannelsQuery, argsKeywords...)
		argsChannelsQuery = append(argsChannelsQuery, limit)
		channelsQuery = `
		SELECT
			Teams.Name AS TeamName,
			Teams.DisplayName AS TeamDisplayName,
			Channels.Name AS ChannelName,
			Channels.DisplayName AS ChannelDisplayName,
			Channels.Type AS ChannelType,
			Users.Username AS UserUsername,
			Users.Email AS UserEmail,
			Users.Nickname AS UserNickname,
			Posts.Id AS PostId,
			Posts.CreateAt AS PostCreateAt,
			Posts.UpdateAt AS PostUpdateAt,
			Posts.DeleteAt AS PostDeleteAt,
			Posts.RootId AS PostRootId,
			Posts.OriginalId AS PostOriginalId,
			Posts.Message AS PostMessage,
			Posts.Type AS PostType,
			Posts.Props AS PostProps,
			Posts.Hashtags AS PostHashtags,
			Posts.FileIds AS PostFileIds,
			Bots.UserId IS NOT NULL AS IsBot
		FROM
			Teams,
			Channels,
			Users,
			Posts
		LEFT JOIN
			Bots ON Bots.UserId = Posts.UserId
		WHERE
			Teams.Id = Channels.TeamId
				AND Posts.ChannelId = Channels.Id
				AND Posts.UserId = Users.Id
				AND (
					Posts.CreateAt > ?
					OR (Posts.CreateAt = ? AND Posts.Id > ?)
				)
				AND Posts.CreateAt < ?
				` + emailQuery + `
				` + keywordQuery + `
		ORDER BY Posts.CreateAt, Posts.Id
		LIMIT ?`
		if err := s.GetReplicaX().Select(&channelPosts, channelsQuery, argsChannelsQuery...); err != nil {
			return nil, cursor, errors.Wrap(err, "unable to export compliance")
		}
		if len(channelPosts) < limit {
			cursor.ChannelsQueryCompleted = true
		} else {
			cursor.LastChannelsQueryPostCreateAt = channelPosts[len(channelPosts)-1].PostCreateAt
			cursor.LastChannelsQueryPostID = channelPosts[len(channelPosts)-1].PostId
		}
	}

	directMessagePosts := []*model.CompliancePost{}
	directMessagesQuery := ""
	var argsDirectMessagesQuery []interface{}
	if !cursor.DirectMessagesQueryCompleted && len(channelPosts) < limit {
		if cursor.LastDirectMessagesQueryPostCreateAt == 0 {
			cursor.LastDirectMessagesQueryPostCreateAt = job.StartAt
		}
		// append the named parameters of SQL query in the correct order to argsDirectMessagesQuery
		argsDirectMessagesQuery = append(argsDirectMessagesQuery, cursor.LastDirectMessagesQueryPostCreateAt, cursor.LastDirectMessagesQueryPostCreateAt, cursor.LastDirectMessagesQueryPostID, job.EndAt)
		argsDirectMessagesQuery = append(argsDirectMessagesQuery, argsEmails...)
		argsDirectMessagesQuery = append(argsDirectMessagesQuery, argsKeywords...)
		argsDirectMessagesQuery = append(argsDirectMessagesQuery, limit-len(channelPosts))
		directMessagesQuery = `
		SELECT
			'direct-messages' AS TeamName,
			'Direct Messages' AS TeamDisplayName,
			Channels.Name AS ChannelName,
			Channels.DisplayName AS ChannelDisplayName,
			Channels.Type AS ChannelType,
			Users.Username AS UserUsername,
			Users.Email AS UserEmail,
			Users.Nickname AS UserNickname,
			Posts.Id AS PostId,
			Posts.CreateAt AS PostCreateAt,
			Posts.UpdateAt AS PostUpdateAt,
			Posts.DeleteAt AS PostDeleteAt,
			Posts.RootId AS PostRootId,
			Posts.OriginalId AS PostOriginalId,
			Posts.Message AS PostMessage,
			Posts.Type AS PostType,
			Posts.Props AS PostProps,
			Posts.Hashtags AS PostHashtags,
			Posts.FileIds AS PostFileIds,
			Bots.UserId IS NOT NULL AS IsBot
		FROM
			Channels,
			Users,
			Posts
		LEFT JOIN
			Bots ON Bots.UserId = Posts.UserId
		WHERE
			Channels.TeamId = ''
				AND Posts.ChannelId = Channels.Id
				AND Posts.UserId = Users.Id
				AND (
					Posts.CreateAt > ?
					OR (Posts.CreateAt = ? AND Posts.Id > ?)
				)
				AND Posts.CreateAt < ?
				` + emailQuery + `
				` + keywordQuery + `
		ORDER BY Posts.CreateAt, Posts.Id
		LIMIT ?`

		if err := s.GetReplicaX().Select(&directMessagePosts, directMessagesQuery, argsDirectMessagesQuery...); err != nil {
			return nil, cursor, errors.Wrap(err, "unable to export compliance")
		}
		if len(directMessagePosts) < limit {
			cursor.DirectMessagesQueryCompleted = true
		} else {
			cursor.LastDirectMessagesQueryPostCreateAt = directMessagePosts[len(directMessagePosts)-1].PostCreateAt
			cursor.LastDirectMessagesQueryPostID = directMessagePosts[len(directMessagePosts)-1].PostId
		}
	}

	return append(channelPosts, directMessagePosts...), cursor, nil
}

func (s *SqlComplianceStore) MessageExport(cursor model.MessageExportCursor, limit int) ([]*model.MessageExport, model.MessageExportCursor, error) {
	panic("not implemented")
}
