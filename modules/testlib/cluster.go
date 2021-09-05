package testlib

import (
	"sync"

	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/cluster"
)

type FakeClusterInterface struct {
	clusterMessageHandler einterfaces.ClusterMessageHandler
	mut                   sync.RWMutex
	messages              []*cluster.ClusterMessage
}

func (c *FakeClusterInterface) StartInterNodeCommunication() {}

func (c *FakeClusterInterface) StopInterNodeCommunication() {}

func (c *FakeClusterInterface) RegisterClusterMessageHandler(event string, crm einterfaces.ClusterMessageHandler) {
	c.clusterMessageHandler = crm
}

func (c *FakeClusterInterface) HealthScore() int {
	return 0
}

func (c *FakeClusterInterface) GetClusterId() string { return "" }

func (c *FakeClusterInterface) IsLeader() bool { return false }

func (c *FakeClusterInterface) GetMyClusterInfo() *cluster.ClusterInfo { return nil }

func (c *FakeClusterInterface) GetClusterInfos() []*cluster.ClusterInfo { return nil }

func (c *FakeClusterInterface) SendClusterMessage(message *cluster.ClusterMessage) {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.messages = append(c.messages, message)
}

func (c *FakeClusterInterface) SendClusterMessageToNode(nodeID string, message *cluster.ClusterMessage) error {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.messages = append(c.messages, message)
	return nil
}

func (c *FakeClusterInterface) NotifyMsg(buf []byte) {}

func (c *FakeClusterInterface) GetClusterStats() ([]*cluster.ClusterStats, *model.AppError) {
	return nil, nil
}

func (c *FakeClusterInterface) GetLogs(page, perPage int) ([]string, *model.AppError) {
	return []string{}, nil
}

func (c *FakeClusterInterface) ConfigChanged(previousConfig *model.Config, newConfig *model.Config, sendToOtherServer bool) *model.AppError {
	return nil
}

func (c *FakeClusterInterface) SendClearRoleCacheMessage() {
	if c.clusterMessageHandler != nil {
		c.clusterMessageHandler(&cluster.ClusterMessage{
			Event: cluster.ClusterEventInvalidateCacheForRoles,
		})
	}
}

// func (c *FakeClusterInterface) GetPluginStatuses() (model.PluginStatuses, *model.AppError) {
// 	return nil, nil
// }

func (c *FakeClusterInterface) GetMessages() []*cluster.ClusterMessage {
	c.mut.RLock()
	defer c.mut.RUnlock()
	return c.messages
}

func (c *FakeClusterInterface) SelectMessages(filterCond func(message *cluster.ClusterMessage) bool) []*cluster.ClusterMessage {
	c.mut.RLock()
	defer c.mut.RUnlock()

	filteredMessages := []*cluster.ClusterMessage{}
	for _, msg := range c.messages {
		if filterCond(msg) {
			filteredMessages = append(filteredMessages, msg)
		}
	}
	return filteredMessages
}

func (c *FakeClusterInterface) ClearMessages() {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.messages = nil
}
