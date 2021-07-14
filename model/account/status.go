package account

import (
	"io"

	"github.com/sitename/sitename/model"
)

const (
	STATUS_OUT_OF_OFFICE   = "ooo"
	STATUS_OFFLINE         = "offline"
	STATUS_AWAY            = "away"
	STATUS_DND             = "dnd"
	STATUS_ONLINE          = "online"
	STATUS_CACHE_SIZE      = model.SESSION_CACHE_SIZE
	STATUS_CHANNEL_TIMEOUT = 20000  // 20 seconds
	STATUS_MIN_UPDATE_TIME = 120000 // 2 minutes
)

type Status struct {
	UserId         string `json:"user_id"`
	Status         string `json:"status"`
	Manual         bool   `json:"manual"`
	LastActivityAt int64  `json:"last_activity_at"`
	ActiveChannel  string `json:"active_channel,omitempty" db:"-"`
}

func (o *Status) ToJson() string {
	oCopy := *o
	oCopy.ActiveChannel = ""
	return model.ModelToJson(&oCopy)
}

func (o *Status) ToClusterJson() string {
	oCopy := *o
	return model.ModelToJson(&oCopy)
}

func StatusFromJson(data io.Reader) *Status {
	var o *Status
	model.ModelFromJson(&o, data)
	return o
}

func StatusListToJson(u []*Status) string {
	uCopy := make([]Status, len(u))
	for i, s := range u {
		sCopy := *s
		sCopy.ActiveChannel = ""
		uCopy[i] = sCopy
	}

	return model.ModelToJson(uCopy)
}

func StatusListFromJson(data io.Reader) []*Status {
	var statuses []*Status
	model.ModelFromJson(&statuses, data)
	return statuses
}

func StatusMapToInterfaceMap(statusMap map[string]*Status) map[string]interface{} {
	interfaceMap := map[string]interface{}{}
	for _, s := range statusMap {
		// Omitted statues mean offline
		if s.Status != STATUS_OFFLINE {
			interfaceMap[s.UserId] = s.Status
		}
	}
	return interfaceMap
}
