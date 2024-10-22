package bleveengine

import (
	"strings"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/services/searchengine"
)

type BLVChannel struct {
	Id          string
	TeamId      []string
	NameSuggest []string
}

type BLVUser struct {
	Id                         string
	SuggestionsWithFullname    []string
	SuggestionsWithoutFullname []string
	TeamsIds                   []string
	ChannelsIds                []string
}

type BLVPost struct {
	Id          string
	TeamId      string
	ChannelId   string
	UserId      string
	CreateAt    int64
	Message     string
	Type        string
	Hashtags    []string
	Attachments string
}

type BLVFile struct {
	Id        string
	CreatorId string
	ChannelId string
	CreateAt  int64
	Name      string
	Content   string
	Extension string
}

// func BLVChannelFromChannel(channel *model.Channel) *BLVChannel {
// 	displayNameInputs := searchengine.GetSuggestionInputsSplitBy(channel.DisplayName, " ")
// 	nameInputs := searchengine.GetSuggestionInputsSplitByMultiple(channel.Name, []string{"-", "_"})

// 	return &BLVChannel{
// 		Id:          channel.Id,
// 		TeamId:      []string{channel.TeamId},
// 		NameSuggest: append(displayNameInputs, nameInputs...),
// 	}
// }

func BLVUserFromUserAndTeams(user *model.User, teamsIds, channelsIds []string) *BLVUser {
	usernameSuggestions := searchengine.GetSuggestionInputsSplitByMultiple(user.Username, []string{".", "-", "_"})

	fullnameStrings := []string{}
	if user.FirstName != "" {
		fullnameStrings = append(fullnameStrings, user.FirstName)
	}
	if user.LastName != "" {
		fullnameStrings = append(fullnameStrings, user.LastName)
	}

	fullnameSuggestions := []string{}
	if len(fullnameStrings) > 0 {
		fullname := strings.Join(fullnameStrings, " ")
		fullnameSuggestions = searchengine.GetSuggestionInputsSplitBy(fullname, " ")
	}

	nicknameSuggesitons := []string{}
	if user.Nickname != "" {
		nicknameSuggesitons = searchengine.GetSuggestionInputsSplitBy(user.Nickname, " ")
	}

	usernameAndNicknameSuggestions := append(usernameSuggestions, nicknameSuggesitons...)

	return &BLVUser{
		Id:                         user.ID,
		SuggestionsWithFullname:    append(usernameAndNicknameSuggestions, fullnameSuggestions...),
		SuggestionsWithoutFullname: usernameAndNicknameSuggestions,
		TeamsIds:                   teamsIds,
		ChannelsIds:                channelsIds,
	}
}

func BLVUserFromUserForIndexing(userForIndexing *model_helper.UserForIndexing) *BLVUser {
	user := &model.User{
		ID:        userForIndexing.Id,
		Username:  userForIndexing.Username,
		Nickname:  userForIndexing.Nickname,
		FirstName: userForIndexing.FirstName,
		LastName:  userForIndexing.LastName,
		CreatedAt: userForIndexing.CreatedAt,
		DeleteAt:  userForIndexing.DeletedAt,
	}

	// return BLVUserFromUserAndTeams(user, userForIndexing.TeamsIds, userForIndexing.ChannelsIds)
	return BLVUserFromUserAndTeams(user, []string{}, []string{})

}

// func BLVPostFromPost(post *model.Post, teamId string) *BLVPost {
// 	p := &model.PostForIndexing{
// 		TeamId: teamId,
// 	}
// 	post.ShallowCopy(&p.Post)
// 	return BLVPostFromPostForIndexing(p)
// }

// func BLVPostFromPostForIndexing(post *model.PostForIndexing) *BLVPost {
// 	return &BLVPost{
// 		Id:        post.Id,
// 		TeamId:    post.TeamId,
// 		ChannelId: post.ChannelId,
// 		UserId:    post.UserId,
// 		CreateAt:  post.CreateAt,
// 		Message:   post.Message,
// 		Type:      post.Type,
// 		Hashtags:  strings.Fields(post.Hashtags),
// 	}
// }

func splitFilenameWords(name string) string {
	result := name
	result = strings.ReplaceAll(result, "-", " ")
	result = strings.ReplaceAll(result, ".", " ")
	return result
}

// func BLVFileFromFileInfo(fileInfo *model.FileInfo, channelId string) *BLVFile {
// 	return &BLVFile{
// 		Id:        fileInfo.Id,
// 		ChannelId: channelId,
// 		CreatorId: fileInfo.CreatorId,
// 		CreateAt:  fileInfo.CreateAt,
// 		Content:   fileInfo.Content,
// 		Extension: fileInfo.Extension,
// 		Name:      fileInfo.Name + " " + splitFilenameWords(fileInfo.Name),
// 	}
// }

// func BLVFileFromFileForIndexing(file *model.FileForIndexing) *BLVFile {
// 	return &BLVFile{
// 		Id:        file.Id,
// 		ChannelId: file.ChannelId,
// 		CreatorId: file.CreatorId,
// 		CreateAt:  file.CreateAt,
// 		Content:   file.Content,
// 		Extension: file.Extension,
// 		Name:      file.Name + " " + splitFilenameWords(file.Name),
// 	}
// }
