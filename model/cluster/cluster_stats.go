package cluster

import (
	"io"

	"github.com/sitename/sitename/model"
)

type ClusterStats struct {
	Id                        string `json:"id"`
	TotalWebsocketConnections int    `json:"total_websocket_connections"`
	TotalReadDbConnections    int    `json:"total_read_db_connections"`
	TotalMasterDbConnections  int    `json:"total_master_db_connections"`
}

func (cs *ClusterStats) ToJson() string {
	return model.ModelToJson(cs)
}

func ClusterStatsFromJson(data io.Reader) *ClusterStats {
	var cs *ClusterStats
	model.ModelFromJson(&cs, data)
	return cs
}
