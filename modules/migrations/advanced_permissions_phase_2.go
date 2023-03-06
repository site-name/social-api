package migrations

import (
	"io"

	"github.com/sitename/sitename/model"
)

type AdvancedPermissionsPhase2Progress struct {
	CurrentTable  string `json:"current_table"`
	LastTeamId    string `json:"last_team_id"`
	LastChannelId string `json:"last_channel_id"`
	LastUserId    string `json:"last_user"`
}

func (p *AdvancedPermissionsPhase2Progress) ToJSON() string {
	return model.ModelToJson(p)
}

func AdvancedPermissionsPhase2ProgressFromJson(data io.Reader) *AdvancedPermissionsPhase2Progress {
	var o *AdvancedPermissionsPhase2Progress
	model.ModelFromJson(&o, data)
	return o
}

func (p *AdvancedPermissionsPhase2Progress) IsValid() bool {
	if !model.IsValidId(p.LastChannelId) {
		return false
	}

	if !model.IsValidId(p.LastTeamId) {
		return false
	}

	if !model.IsValidId(p.LastUserId) {
		return false
	}

	switch p.CurrentTable {
	case "TeamMembers":
	case "ChannelMembers":
	default:
		return false
	}

	return true
}

func (worker *Worker) runAdvancedPermissionsPhase2Migration(lastDone string) (bool, string, *model.AppError) {
	panic("not implemented")
}
