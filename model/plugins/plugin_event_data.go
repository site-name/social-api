package plugins

import (
	"github.com/sitename/sitename/modules/json"
)

// PluginEventData used to notify peers about plugin changes.
type PluginEventData struct {
	Id string `json:"id"`
}

func (p *PluginEventData) ToJSON() []byte {
	res, _ := json.JSON.Marshal(p)
	return res
}
