package migrations

import (
	"encoding/json"
	"io"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
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
	err := json.NewDecoder(data).Decode(&o)
	if err != nil {
		slog.Warn("Error decoding advanced permissions phase 2 progress", slog.Err(err))
	}
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
	slog.Debug("runAdvancedPermissionsPhase2Migration...")
	// TODO: consider remove me
	return false, "", nil
}
