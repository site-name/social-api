package einterfaces

import (
	"github.com/sitename/sitename/model_helper"
)

type DataRetentionInterface interface {
	// GetGlobalPolicy() (*model.GlobalRetentionPolicy, *model_helper.AppError)
	// GetPolicies(offset, limit int) (*model.RetentionPolicyWithTeamAndChannelCountsList, *model_helper.AppError)
	GetPoliciesCount() (int64, *model_helper.AppError)
	// GetPolicy(policyID string) (*model.RetentionPolicyWithTeamAndChannelCounts, *model_helper.AppError)
	// CreatePolicy(policy *model.RetentionPolicyWithTeamAndChannelIDs) (*model.RetentionPolicyWithTeamAndChannelCounts, *model_helper.AppError)
	// PatchPolicy(patch *model.RetentionPolicyWithTeamAndChannelIDs) (*model.RetentionPolicyWithTeamAndChannelCounts, *model_helper.AppError)
	DeletePolicy(policyID string) *model_helper.AppError
	// GetTeamsForPolicy(policyID string, offset, limit int) (*model.TeamsWithCount, *model_helper.AppError)
	AddTeamsToPolicy(policyID string, teamIDs []string) *model_helper.AppError
	RemoveTeamsFromPolicy(policyID string, teamIDs []string) *model_helper.AppError
	// GetChannelsForPolicy(policyID string, offset, limit int) (*model.ChannelsWithCount, *model_helper.AppError)
	AddChannelsToPolicy(policyID string, channelIDs []string) *model_helper.AppError
	RemoveChannelsFromPolicy(policyID string, channelIDs []string) *model_helper.AppError
	// GetTeamPoliciesForUser(userID string, offset, limit int) (*model.RetentionPolicyForTeamList, *model_helper.AppError)
	// GetChannelPoliciesForUser(userID string, offset, limit int) (*model.RetentionPolicyForChannelList, *model_helper.AppError)
}
