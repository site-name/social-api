package app

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model_helper"
	"github.com/stretchr/testify/require"
)

func TestBusySet(t *testing.T) {
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
	require.True(t, compareBusyState(t, busy, cluster.Busy))
	// should automatically expire after 100ms
	require.Eventually(t, isNotBusy, time.Second*15, time.Millisecond*20)
	// allow a moment for cluster to sync.
	require.Eventually(t, func() bool { return compareBusyState(t, busy, cluster.Busy) }, time.Second*15, time.Millisecond*20)

	// tes set after auto expiry
	busy.Set(time.Second * 30)
	require.True(t, busy.IsBusy())
	require.True(t, compareBusyState(t, busy, cluster.Busy))
	expire := busy.Expires()
	require.Greater(t, expire.Unix(), time.Now().Add(time.Second*10).Unix())

	// test extending existing expiry
	busy.Set(time.Minute * 5)
	require.True(t, busy.IsBusy())
	require.True(t, compareBusyState(t, busy, cluster.Busy))
	expire = busy.Expires()
	require.Greater(t, expire.Unix(), time.Now().Add(time.Minute*2).Unix())

	busy.Clear()
	require.False(t, busy.IsBusy())
	require.True(t, compareBusyState(t, busy, cluster.Busy))
}

// ClusterMock simulates the busy state of a cluster.
type ClusterMock struct {
	Busy *Busy
}

func (c *ClusterMock) SendClusterMessage(msg *model_helper.ClusterMessage) {
	var sbs model_helper.ServerBusyState
	json.Unmarshal(msg.Data, &sbs)
	c.Busy.ClusterEventChanged(&sbs)
}

func (c *ClusterMock) SendClusterMessageToNode(nodeID string, msg *model_helper.ClusterMessage) error {
	return nil
}

func TestBusyRace(t *testing.T) {
	cluster := &ClusterMock{Busy: &Busy{}}
	busy := NewBusy(cluster)

	busy.Set(500 * time.Millisecond)

	// we are sleeping in order to let the race trigger
	time.Sleep(time.Second)
}
func compareBusyState(t *testing.T, busy1 *Busy, busy2 *Busy) bool {
	t.Helper()
	if busy1.IsBusy() != busy2.IsBusy() || busy1.Expires().Unix() != busy2.Expires().Unix() {
		busy1JSON, _ := busy1.ToJSON()
		busy2JSON, _ := busy2.ToJSON()
		t.Logf("busy1:%s; busy:%s\n", busy1JSON, busy2JSON)
		return false
	}

	return true
}

func (c *ClusterMock) StartInterNodeCommunication() {}
func (c *ClusterMock) StopInterNodeCommunication()  {}
func (c *ClusterMock) RegisterClusterMessageHandler(event model_helper.ClusterEvent, crm einterfaces.ClusterMessageHandler) {
}
func (c *ClusterMock) GetClusterId() string                         { return "cluster_mock" }
func (c *ClusterMock) IsLeader() bool                               { return false }
func (c *ClusterMock) GetMyClusterInfo() *model_helper.ClusterInfo  { return nil }
func (c *ClusterMock) GetClusterInfos() []*model_helper.ClusterInfo { return nil }
func (c *ClusterMock) NotifyMsg(buf []byte)                         {}
func (c *ClusterMock) GetClusterStats() ([]*model_helper.ClusterStats, *model_helper.AppError) {
	return nil, nil
}
func (c *ClusterMock) GetLogs(page, perPage int) ([]string, *model_helper.AppError) { return nil, nil }
func (c *ClusterMock) GetPluginStatuses() (model_helper.PluginStatuses, *model_helper.AppError) {
	return nil, nil
}
func (c *ClusterMock) ConfigChanged(previousConfig *model_helper.Config, newConfig *model_helper.Config, sendToOtherServer bool) *model_helper.AppError {
	return nil
}
func (c *ClusterMock) HealthScore() int { return 0 }
