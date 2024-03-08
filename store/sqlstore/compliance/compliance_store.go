package compliance

import (
	"database/sql"
	"fmt"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SqlComplianceStore struct {
	store.Store
}

func NewSqlComplianceStore(s store.Store) store.ComplianceStore {
	return &SqlComplianceStore{s}
}

func (s *SqlComplianceStore) Upsert(compliance model.Compliance) (*model.Compliance, error) {
	isSaving := compliance.ID == ""
	if isSaving {
		model_helper.CompliancePreSave(&compliance)
	} else {
		model_helper.ComplianceCommonPre(&compliance)
	}

	if err := model_helper.ComplianceIsValid(compliance); err != nil {
		return nil, err
	}

	var err error
	if isSaving {
		err = compliance.Insert(s.GetMaster(), boil.Infer())
	} else {
		_, err = compliance.Update(s.GetMaster(), boil.Blacklist(model.ComplianceColumns.CreatedAt))
	}

	if err != nil {
		return nil, err
	}

	return &compliance, nil
}

func (s *SqlComplianceStore) GetAll(offset, limit int) (model.ComplianceSlice, error) {
	return model.Compliances(
		qm.OrderBy(fmt.Sprintf("%s %s", model.ComplianceColumns.CreatedAt, model_helper.DESC)),
		qm.Offset(offset),
		qm.Limit(limit),
	).All(s.GetReplica())
}

func (s *SqlComplianceStore) Get(id string) (*model.Compliance, error) {
	record, err := model.FindCompliance(s.GetReplica(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound(model.TableNames.Compliances, id)
		}
		return nil, err
	}

	return record, nil
}

func (s *SqlComplianceStore) ComplianceExport(job model.Compliance, cursor model_helper.ComplianceExportCursor, limit int) ([]*model_helper.CompliancePost, model_helper.ComplianceExportCursor, error) {
	// 	keywordQuery := ""
	// 	var argsKeywords []any
	// 	keywords := strings.Fields(strings.TrimSpace(strings.ToLower(strings.Replace(job.Keywords, ",", " ", -1))))
	// 	if len(keywords) > 0 {
	// 		clauses := make([]string, len(keywords))

	// 		for i, keyword := range keywords {
	// 			keyword = store.SanitizeSearchTerm(keyword, "\\")
	// 			clauses[i] = "LOWER(Posts.Message) LIKE ?"
	// 			argsKeywords = append(argsKeywords, "%"+keyword+"%")
	// 		}

	// 		keywordQuery = "AND (" + strings.Join(clauses, " OR ") + ")"
	// 	}

	// 	emailQuery := ""
	// 	var argsEmails []any
	// 	emails := strings.Fields(strings.TrimSpace(strings.ToLower(strings.Replace(job.Emails, ",", " ", -1))))
	// 	if len(emails) > 0 {
	// 		clauses := make([]string, len(emails))

	// 		for i, email := range emails {
	// 			clauses[i] = "Users.Email = ?"
	// 			argsEmails = append(argsEmails, email)
	// 		}

	// 		emailQuery = "AND (" + strings.Join(clauses, " OR ") + ")"
	// 	}

	// 	// The idea is to first iterate over the channel posts, and then when we run out of those,
	// 	// start iterating over the direct message posts.

	// 	channelPosts := []*model.CompliancePost{}
	// 	channelsQuery := ""
	// 	var argsChannelsQuery []any
	// 	if !cursor.ChannelsQueryCompleted {
	// 		if cursor.LastChannelsQueryPostCreateAt == 0 {
	// 			cursor.LastChannelsQueryPostCreateAt = job.StartAt
	// 		}
	// 		// append the named parameters of SQL query in the correct order to argsChannelsQuery
	// 		argsChannelsQuery = append(argsChannelsQuery, cursor.LastChannelsQueryPostCreateAt, cursor.LastChannelsQueryPostCreateAt, cursor.LastChannelsQueryPostID, job.EndAt)
	// 		argsChannelsQuery = append(argsChannelsQuery, argsEmails...)
	// 		argsChannelsQuery = append(argsChannelsQuery, argsKeywords...)
	// 		argsChannelsQuery = append(argsChannelsQuery, limit)
	// 		channelsQuery = `
	// 		SELECT
	// 			Teams.Name AS TeamName,
	// 			Teams.DisplayName AS TeamDisplayName,
	// 			Channels.Name AS ChannelName,
	// 			Channels.DisplayName AS ChannelDisplayName,
	// 			Channels.Type AS ChannelType,
	// 			Users.Username AS UserUsername,
	// 			Users.Email AS UserEmail,
	// 			Users.Nickname AS UserNickname,
	// 			Posts.Id AS PostId,
	// 			Posts.CreateAt AS PostCreateAt,
	// 			Posts.UpdateAt AS PostUpdateAt,
	// 			Posts.DeleteAt AS PostDeleteAt,
	// 			Posts.RootId AS PostRootId,
	// 			Posts.OriginalId AS PostOriginalId,
	// 			Posts.Message AS PostMessage,
	// 			Posts.Type AS PostType,
	// 			Posts.Props AS PostProps,
	// 			Posts.Hashtags AS PostHashtags,
	// 			Posts.FileIds AS PostFileIds,
	// 			Bots.UserId IS NOT NULL AS IsBot
	// 		FROM
	// 			Teams,
	// 			Channels,
	// 			Users,
	// 			Posts
	// 		LEFT JOIN
	// 			Bots ON Bots.UserId = Posts.UserId
	// 		WHERE
	// 			Teams.Id = Channels.TeamId
	// 				AND Posts.ChannelId = Channels.Id
	// 				AND Posts.UserId = Users.Id
	// 				AND (
	// 					Posts.CreateAt > ?
	// 					OR (Posts.CreateAt = ? AND Posts.Id > ?)
	// 				)
	// 				AND Posts.CreateAt < ?
	// 				` + emailQuery + `
	// 				` + keywordQuery + `
	// 		ORDER BY Posts.CreateAt, Posts.Id
	// 		LIMIT ?`
	// 		if err := s.GetReplica().Raw(channelsQuery, argsChannelsQuery...).Scan(&channelPosts).Error; err != nil {
	// 			return nil, cursor, errors.Wrap(err, "unable to export compliance")
	// 		}
	// 		if len(channelPosts) < limit {
	// 			cursor.ChannelsQueryCompleted = true
	// 		} else {
	// 			cursor.LastChannelsQueryPostCreateAt = channelPosts[len(channelPosts)-1].PostCreateAt
	// 			cursor.LastChannelsQueryPostID = channelPosts[len(channelPosts)-1].PostId
	// 		}
	// 	}

	// 	directMessagePosts := []*model.CompliancePost{}
	// 	directMessagesQuery := ""
	// 	var argsDirectMessagesQuery []any
	// 	if !cursor.DirectMessagesQueryCompleted && len(channelPosts) < limit {
	// 		if cursor.LastDirectMessagesQueryPostCreateAt == 0 {
	// 			cursor.LastDirectMessagesQueryPostCreateAt = job.StartAt
	// 		}
	// 		// append the named parameters of SQL query in the correct order to argsDirectMessagesQuery
	// 		argsDirectMessagesQuery = append(argsDirectMessagesQuery, cursor.LastDirectMessagesQueryPostCreateAt, cursor.LastDirectMessagesQueryPostCreateAt, cursor.LastDirectMessagesQueryPostID, job.EndAt)
	// 		argsDirectMessagesQuery = append(argsDirectMessagesQuery, argsEmails...)
	// 		argsDirectMessagesQuery = append(argsDirectMessagesQuery, argsKeywords...)
	// 		argsDirectMessagesQuery = append(argsDirectMessagesQuery, limit-len(channelPosts))
	// 		directMessagesQuery = `
	// 		SELECT
	// 			'direct-messages' AS TeamName,
	// 			'Direct Messages' AS TeamDisplayName,
	// 			Channels.Name AS ChannelName,
	// 			Channels.DisplayName AS ChannelDisplayName,
	// 			Channels.Type AS ChannelType,
	// 			Users.Username AS UserUsername,
	// 			Users.Email AS UserEmail,
	// 			Users.Nickname AS UserNickname,
	// 			Posts.Id AS PostId,
	// 			Posts.CreateAt AS PostCreateAt,
	// 			Posts.UpdateAt AS PostUpdateAt,
	// 			Posts.DeleteAt AS PostDeleteAt,
	// 			Posts.RootId AS PostRootId,
	// 			Posts.OriginalId AS PostOriginalId,
	// 			Posts.Message AS PostMessage,
	// 			Posts.Type AS PostType,
	// 			Posts.Props AS PostProps,
	// 			Posts.Hashtags AS PostHashtags,
	// 			Posts.FileIds AS PostFileIds,
	// 			Bots.UserId IS NOT NULL AS IsBot
	// 		FROM
	// 			Channels,
	// 			Users,
	// 			Posts
	// 		LEFT JOIN
	// 			Bots ON Bots.UserId = Posts.UserId
	// 		WHERE
	// 			Channels.TeamId = ''
	// 				AND Posts.ChannelId = Channels.Id
	// 				AND Posts.UserId = Users.Id
	// 				AND (
	// 					Posts.CreateAt > ?
	// 					OR (Posts.CreateAt = ? AND Posts.Id > ?)
	// 				)
	// 				AND Posts.CreateAt < ?
	// 				` + emailQuery + `
	// 				` + keywordQuery + `
	// 		ORDER BY Posts.CreateAt, Posts.Id
	// 		LIMIT ?`

	// 		if err := s.GetReplica().Raw(directMessagesQuery, argsDirectMessagesQuery...).Scan(&directMessagePosts).Error; err != nil {
	// 			return nil, cursor, errors.Wrap(err, "unable to export compliance")
	// 		}
	// 		if len(directMessagePosts) < limit {
	// 			cursor.DirectMessagesQueryCompleted = true
	// 		} else {
	// 			cursor.LastDirectMessagesQueryPostCreateAt = directMessagePosts[len(directMessagePosts)-1].PostCreateAt
	// 			cursor.LastDirectMessagesQueryPostID = directMessagePosts[len(directMessagePosts)-1].PostId
	// 		}
	// 	}

	//		return append(channelPosts, directMessagePosts...), cursor, nil
	//	}
	panic("not implemented")
}

func (s *SqlComplianceStore) MessageExport(cursor model_helper.MessageExportCursor, limit int) ([]*model_helper.MessageExport, model_helper.MessageExportCursor, error) {
	panic("not implemented")
}
