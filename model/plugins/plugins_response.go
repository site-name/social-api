package plugins

import (
	"encoding/json"
	"io"

	"github.com/sitename/sitename/model"
)

type PluginInfo struct {
	model.Manifest
}

type PluginsResponse struct {
	Active   []*PluginInfo `json:"active"`
	Inactive []*PluginInfo `json:"inactive"`
}

func (m *PluginsResponse) ToJson() string {
	b, _ := json.Marshal(m)
	return string(b)
}

func PluginsResponseFromJson(data io.Reader) *PluginsResponse {
	var m *PluginsResponse
	json.NewDecoder(data).Decode(&m)
	return m
}
