package plugins

import (
	"io"

	"github.com/sitename/sitename/model"
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

func PluginEventDataFromJson(data io.Reader) PluginEventData {
	var m PluginEventData

	model.ModelFromJson(&m, data)
	return m
}
