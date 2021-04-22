package model

import (
	"io"

	"github.com/sitename/sitename/modules/json"
)

type ClusterStats struct {
	Id                        string `json:"id"`
	TotalWebsocketConnections int    `json:"total_websocket_connections"`
	TotalReadDbConnections    int    `json:"total_read_db_connections"`
	TotalMasterDbConnections  int    `json:"total_master_db_connections"`
}

func (cs *ClusterStats) ToJson() string {
	b, _ := json.JSON.Marshal(cs)
	return string(b)
}

func ClusterStatsFromJson(data io.Reader) *ClusterStats {
	var cs *ClusterStats
	json.JSON.NewDecoder(data).Decode(&cs)
	return cs
}
