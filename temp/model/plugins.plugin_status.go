package model

import (
	"encoding/json"
)

const (
	PluginStateNotRunning          = 0
	PluginStateStarting            = 1 // unused by server
	PluginStateRunning             = 2
	PluginStateFailedToStart       = 3
	PluginStateFailedToStayRunning = 4
	PluginStateStopping            = 5 // unused by server
)

// PluginStatus provides a cluster-aware view of installed plugins.
type PluginStatus struct {
	PluginId    string `json:"plugin_id"`
	ClusterId   string `json:"cluster_id"`
	PluginPath  string `json:"plugin_path"`
	State       int    `json:"state"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type PluginStatuses []*PluginStatus

func (m *PluginStatuses) ToJSON() string {
	b, _ := json.Marshal(m)
	return string(b)
}