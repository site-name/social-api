package model

type ClusterStats struct {
	Id                        string `json:"id"`
	TotalWebsocketConnections int    `json:"total_websocket_connections"`
	TotalReadDbConnections    int    `json:"total_read_db_connections"`
	TotalMasterDbConnections  int    `json:"total_master_db_connections"`
}

func (cs *ClusterStats) ToJSON() string {
	return ModelToJson(cs)
}
