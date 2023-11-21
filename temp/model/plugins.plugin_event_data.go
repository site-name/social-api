package model

import "encoding/json"

// PluginEventData used to notify peers about plugin changes.
type PluginEventData struct {
	Id string `json:"id"`
}

func (p *PluginEventData) ToJSON() []byte {
	res, _ := json.Marshal(p)
	return res
}
