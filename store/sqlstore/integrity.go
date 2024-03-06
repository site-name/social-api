package sqlstore

import (
	"github.com/mattermost/squirrel"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
	"github.com/volatiletech/sqlboiler/v4/queries"
)

type relationalCheckConfig struct {
	parentName         string
	parentIdAttr       string
	childName          string
	childIdAttr        string
	canParentIdBeEmpty bool
	sortRecords        bool
	filter             interface{}
}

func getOrphanedRecords(ss *SqlStore, cfg relationalCheckConfig) ([]model_helper.OrphanedRecord, error) {
	var records []model_helper.OrphanedRecord

	sub := ss.GetQueryBuilder().
		Select("TRUE").
		From(cfg.parentName + " AS PT").
		Prefix("NOT Exists (").
		Suffix(")").
		Where("PT.id = CT." + cfg.parentIdAttr)

	main := ss.GetQueryBuilder().
		Select().
		Column("CT." + cfg.parentIdAttr + " AS ParentId").
		From(cfg.childName + " AS CT").
		Where(sub)

	if cfg.childIdAttr != "" {
		main = main.Column("CT." + cfg.childIdAttr + " AS ChildId")
	}

	if cfg.canParentIdBeEmpty {
		main = main.Where(squirrel.NotEq{"CT." + cfg.parentIdAttr: ""})
	}

	if cfg.filter != nil {
		main = main.Where(cfg.filter)
	}

	if cfg.sortRecords {
		main = main.OrderBy("CT." + cfg.parentIdAttr)
	}

	query, args, _ := main.ToSql()

	return records, queries.Raw(query, args...).Bind(ss.context, ss.GetMaster(), &records)
}

func checkParentChildIntegrity(ss *SqlStore, config relationalCheckConfig) model_helper.IntegrityCheckResult {
	var result model_helper.IntegrityCheckResult
	var data model_helper.RelationalIntegrityCheckData

	config.sortRecords = true
	data.Records, result.Err = getOrphanedRecords(ss, config)
	if result.Err != nil {
		slog.Error("Error while getting orphaned records", slog.Err(result.Err))
		return result
	}
	data.ParentName = config.parentName
	data.ChildName = config.childName
	data.ParentIdAttr = config.parentIdAttr
	data.ChildIdAttr = config.childIdAttr
	result.Data = data

	return result
}

func checkUsersAuditsIntegrity(ss *SqlStore) model_helper.IntegrityCheckResult {
	return checkParentChildIntegrity(ss, relationalCheckConfig{
		parentName:         "Users",
		parentIdAttr:       "UserId",
		childName:          "Audits",
		childIdAttr:        "Id",
		canParentIdBeEmpty: true,
	})
}
func checkUsersCommandWebhooksIntegrity(ss *SqlStore) model_helper.IntegrityCheckResult {
	return checkParentChildIntegrity(ss, relationalCheckConfig{
		parentName:   "Users",
		parentIdAttr: "UserId",
		childName:    "CommandWebhooks",
		childIdAttr:  "Id",
	})
}
func checkUsersChannelMemberHistoryIntegrity(ss *SqlStore) model_helper.IntegrityCheckResult {
	return checkParentChildIntegrity(ss, relationalCheckConfig{
		parentName:   "Users",
		parentIdAttr: "UserId",
		childName:    "ChannelMemberHistory",
		childIdAttr:  "",
	})
}

func checkUsersIntegrity(ss *SqlStore, results chan<- model_helper.IntegrityCheckResult) {
	results <- checkUsersAuditsIntegrity(ss)
	results <- checkUsersCommandWebhooksIntegrity(ss)
	results <- checkUsersChannelMemberHistoryIntegrity(ss)
	// results <- checkUsersChannelMembersIntegrity(ss)
	// results <- checkUsersChannelsIntegrity(ss)
	// results <- checkUsersCommandsIntegrity(ss)
	// results <- checkUsersCompliancesIntegrity(ss)
	// results <- checkUsersEmojiIntegrity(ss)
	// results <- checkUsersFileInfoIntegrity(ss)
	// results <- checkUsersIncomingWebhooksIntegrity(ss)
	// results <- checkUsersOAuthAccessDataIntegrity(ss)
	// results <- checkUsersOAuthAppsIntegrity(ss)
	// results <- checkUsersOAuthAuthDataIntegrity(ss)
	// results <- checkUsersOutgoingWebhooksIntegrity(ss)
	// results <- checkUsersPostsIntegrity(ss)
	// results <- checkUsersPreferencesIntegrity(ss)
	// results <- checkUsersReactionsIntegrity(ss)
	// results <- checkUsersSessionsIntegrity(ss)
	// results <- checkUsersStatusIntegrity(ss)
	// results <- checkUsersTeamMembersIntegrity(ss)
	// results <- checkUsersUserAccessTokensIntegrity(ss)
}

func CheckRelationalIntegrity(ss *SqlStore, results chan<- model_helper.IntegrityCheckResult) {
	slog.Info("Starting relational integrity checks...")
	// checkChannelsIntegrity(ss, results)
	// checkCommandsIntegrity(ss, results)
	// checkPostsIntegrity(ss, results)
	// checkSchemesIntegrity(ss, results)
	// checkSessionsIntegrity(ss, results)
	// checkTeamsIntegrity(ss, results)
	checkUsersIntegrity(ss, results)
	slog.Info("Done with relational integrity checks")
	close(results)
}
