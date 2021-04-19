package sqlstore

import (
	"github.com/Masterminds/squirrel"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/log"
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

func getOrphanedRecords(ss *SqlStore, cfg relationalCheckConfig) ([]model.OrphanedRecord, error) {
	var records []model.OrphanedRecord

	sub := ss.getQueryBuilder().
		Select("TRUE").
		From(cfg.parentName + " AS PT").
		Prefix("NOT Exists (").
		Suffix(")").
		Where("PT.id = CT." + cfg.parentIdAttr)

	main := ss.getQueryBuilder().
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
	_, err := ss.GetMaster().Select(&records, query, args...)

	return records, err
}

func checkParentChildIntegrity(ss *SqlStore, config relationalCheckConfig) model.IntegrityCheckResult {
	var result model.IntegrityCheckResult
	var data model.RelationalIntegrityCheckData

	config.sortRecords = true
	data.Records, result.Err = getOrphanedRecords(ss, config)
	if result.Err != nil {
		log.Error("Error while getting orphaned records", result.Err)
		return result
	}
	data.ParentName = config.parentName
	data.ChildName = config.childName
	data.ParentIdAttr = config.parentIdAttr
	data.ChildIdAttr = config.childIdAttr
	result.Data = data

	return result
}

func checkUsersAuditsIntegrity(ss *SqlStore) model.IntegrityCheckResult {
	return checkParentChildIntegrity(ss, relationalCheckConfig{
		parentName:         "Users",
		parentIdAttr:       "UserId",
		childName:          "Audits",
		childIdAttr:        "Id",
		canParentIdBeEmpty: true,
	})
}
func checkUsersCommandWebhooksIntegrity(ss *SqlStore) model.IntegrityCheckResult {
	return checkParentChildIntegrity(ss, relationalCheckConfig{
		parentName:   "Users",
		parentIdAttr: "UserId",
		childName:    "CommandWebhooks",
		childIdAttr:  "Id",
	})
}
func checkUsersChannelMemberHistoryIntegrity(ss *SqlStore) model.IntegrityCheckResult {
	return checkParentChildIntegrity(ss, relationalCheckConfig{
		parentName:   "Users",
		parentIdAttr: "UserId",
		childName:    "ChannelMemberHistory",
		childIdAttr:  "",
	})
}

func checkUsersIntegrity(ss *SqlStore, results chan<- model.IntegrityCheckResult) {
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

func CheckRelationalIntegrity(ss *SqlStore, results chan<- model.IntegrityCheckResult) {
	log.Info("Starting relational integrity checks...")

	log.Info("Done with relational integrity checks")
	close(results)
}
