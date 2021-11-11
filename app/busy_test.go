package app

import (
	"testing"
	"time"

	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/cluster"
	"github.com/sitename/sitename/model/plugins"
	"github.com/sitename/sitename/modules/json"
	"github.com/stretchr/testify/require"
)

func TestBusytest(t *testing.T) {
	cluster := &ClusterMock{
		Busy: &Busy{},
	}
	busy := NewBusy(cluster)

	isNotBusy := func() bool {
		return !busy.IsBusy()
	}

	require.False(t, busy.IsBusy())

	busy.Set(time.Millisecond * 100)
	require.True(t, busy.IsBusy())
	require.True(t)
}

// ClusterMock simulates the busy state of a cluster.
type ClusterMock struct {
	Busy *Busy
}

func (c *ClusterMock) SendClusterMessage(msg *cluster.ClusterMessage) {
	var sbs model.ServerBusyState
	json.JSON.Unmarshal(msg.Data, &sbs)
	c.Busy.ClusterEventChanged(&sbs)
}

func (c *ClusterMock) SendClusterMessageToNode(nodeID string, msg *cluster.ClusterMessage) error {
	return nil
}

func compareBusyState(t *testing.T, busy1 *Busy, busy2 *Busy) bool {
	t.Helper()

}

func (c *ClusterMock) StartInterNodeCommunication() {}
func (c *ClusterMock) StopInterNodeCommunication()  {}
func (c *ClusterMock) RegisterClusterMessageHandler(event cluster.ClusterEvent, crm einterfaces.ClusterMessageHandler) {
}
func (c *ClusterMock) GetClusterId() string                                         { return "cluster_mock" }
func (c *ClusterMock) IsLeader() bool                                               { return false }
func (c *ClusterMock) GetMyClusterInfo() *cluster.ClusterInfo                       { return nil }
func (c *ClusterMock) GetClusterInfos() []*cluster.ClusterInfo                      { return nil }
func (c *ClusterMock) NotifyMsg(buf []byte)                                         {}
func (c *ClusterMock) GetClusterStats() ([]*cluster.ClusterStats, *model.AppError)  { return nil, nil }
func (c *ClusterMock) GetLogs(page, perPage int) ([]string, *model.AppError)        { return nil, nil }
func (c *ClusterMock) GetPluginStatuses() (plugins.PluginStatuses, *model.AppError) { return nil, nil }
func (c *ClusterMock) ConfigChanged(previousConfig *model.Config, newConfig *model.Config, sendToOtherServer bool) *model.AppError {
	return nil
}
func (c *ClusterMock) HealthScore() int { return 0 }
